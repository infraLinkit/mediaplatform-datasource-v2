package messaging

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type (
	RabbitDeclare struct {
		ExchangeName string
		ExchangeType string
		QueueName    string
		RoutingKey   string
		Durable      bool
	}

	RabbitManager struct {
		URL       string
		Conn      *amqp.Connection
		ChPool    chan *amqp.Channel
		PoolSize  int
		Qos       int
		Configs   []RabbitDeclare
		Mu        sync.RWMutex
		IsReady   bool
		IsClosing bool
		Wg        sync.WaitGroup
		L         *logrus.Logger
	}
)

func NewRabbitManager(l *logrus.Logger, url string, poolSize int, configs []RabbitDeclare, qos int) (*RabbitManager, error) {
	rm := &RabbitManager{
		URL:      url,
		PoolSize: poolSize,
		Configs:  configs,
		ChPool:   make(chan *amqp.Channel, poolSize),
		Qos:      qos,
		L:        l,
	}

	if err := rm.connect(); err != nil {
		return nil, err
	}

	go rm.handleReconnect()
	return rm, nil
}

func (rm *RabbitManager) connect() error {
	rm.Mu.Lock()
	defer rm.Mu.Unlock()

	var conn *amqp.Connection
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(rm.URL)
		if err == nil {
			break
		}
		rm.L.Warnf("⚠️ RabbitMQ connection attempt %d/%d failed: %v. Retrying in 2s...", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("could not establish connection to RabbitMQ at %s after %d retries: %w", rm.URL, maxRetries, err)
	}
	rm.Conn = conn

	// 1. Declare Infrastructure using a temporary channel
	setupCh, err := conn.Channel()
	if err != nil {
		return err
	}
	for _, cfg := range rm.Configs {
		_ = setupCh.ExchangeDeclare(cfg.ExchangeName, cfg.ExchangeType, cfg.Durable, false, false, false, nil)
		args := amqp.Table{"x-queue-type": "classic", "x-queue-version": 2, "x-queue-mode": "lazy"}
		_, _ = setupCh.QueueDeclare(cfg.QueueName, cfg.Durable, false, false, false, args)
		_ = setupCh.QueueBind(cfg.QueueName, cfg.RoutingKey, cfg.ExchangeName, false, nil)
	}
	setupCh.Close()

	// 2. Clear and Refill Channel Pool
	for len(rm.ChPool) > 0 {
		<-rm.ChPool
	}
	for i := 0; i < rm.PoolSize; i++ {
		ch, err := conn.Channel()
		if err != nil {
			rm.L.Errorf("Re-Initiate channel RabbitMQ Failed : %#v", err)
			return err
		}

		// SET QOS HERE (Speed limit for consumers)
		_ = ch.Qos(rm.Qos, 0, false)
		_ = ch.Confirm(false) // Enable for publisher reliability
		rm.ChPool <- ch
	}

	rm.IsReady = true
	rm.L.Info("RabbitMQ: Infrastructure Ready and Pool Populated")
	return nil
}

func (rm *RabbitManager) handleReconnect() {
	for {
		if rm.IsClosing {
			return
		}

		rm.Mu.RLock()
		connection := rm.Conn
		rm.Mu.RUnlock()

		// If connection is nil or closed, try to connect
		if connection == nil || connection.IsClosed() {
			rm.L.Info("RabbitMQ connection lost. Retrying...")
			if err := rm.connect(); err != nil {
				time.Sleep(5 * time.Second)
				continue // Keep trying
			}
		}

		// --- THE CRITICAL PART ---
		// Create a channel to listen for the NEXT closure
		notify := rm.Conn.NotifyClose(make(chan *amqp.Error))

		// This blocks here until the connection drops
		err := <-notify
		if err != nil {
			rm.L.Warnf("RabbitMQ connection closed: %v", err)
			rm.Mu.Lock()
			rm.IsReady = false
			rm.Mu.Unlock()
		}
		// Loop repeats and goes back to the 'connect()' logic
	}
}

func (rm *RabbitManager) Publish(ctx context.Context, exchange, routingKey string, body []byte, corid string) error {
	rm.Mu.RLock()
	ready := rm.IsReady
	rm.Mu.RUnlock()
	if !ready {
		rm.L.Info("Publisher not ready")
		return fmt.Errorf("publisher not ready")
	}

	ch := <-rm.ChPool
	defer func() { rm.ChPool <- ch }()

	// Persistent ensures message survives RabbitMQ restart
	msg := amqp.Publishing{
		DeliveryMode:  amqp.Persistent,
		ContentType:   "application/json",
		Body:          body,
		Priority:      0,
		CorrelationId: corid,
	}

	conf, err := ch.PublishWithDeferredConfirmWithContext(ctx, exchange, routingKey, false, false, msg)
	if err != nil {
		rm.L.Errorf("publish err %#v", err)
		return err
	}

	if !conf.Wait() {
		rm.L.Warn("message nack'd by broker (check disk/memory alarms)")
		return fmt.Errorf("message nack'd by broker (check disk/memory alarms)")
	}
	return nil
}

func (rm *RabbitManager) PublishWithRetry(ctx context.Context, exchange, routingKey string, body []byte, corid string) error {
	var lastErr error

	// Attempt 3 retries
	for i := 0; i < 3; i++ {
		err := rm.Publish(ctx, exchange, routingKey, body, corid)
		if err == nil {
			return nil // Success!
		}

		lastErr = err
		rm.L.Warnf("Nack received or Publish failed: %#v. Retrying in %d ms...", err, (i+1)*500)

		// Wait before trying again (Exponential backoff)
		time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
	}

	if lastErr != nil {
		rm.L.Info("CRITICAL: Failed to publish after retries. Connection might be stale.")

		// Optional: Close the connection to force the BACKGROUND handleReconnect to wake up
		rm.Mu.Lock()
		if rm.Conn != nil {
			rm.Conn.Close()
		}
		rm.Mu.Unlock()
	}

	return fmt.Errorf("message dropped after 3 retries: %w", lastErr)
}

func (rm *RabbitManager) StartConsuming(queueName string, handler func(ch *amqp.Delivery, channel *amqp.Channel)) {
	workers := rm.PoolSize
	if workers < 1 {
		workers = 1
	}
	for i := 0; i < workers; i++ {
		rm.Wg.Add(1)
		go func(workerID int) {
			defer rm.Wg.Done()
			for {
				if rm.IsClosing {
					return
				}

				rm.Mu.RLock()
				conn := rm.Conn
				rm.Mu.RUnlock()

				if conn == nil || conn.IsClosed() {
					time.Sleep(2 * time.Second)
					continue
				}

				ch, err := conn.Channel()
				if err != nil {
					time.Sleep(2 * time.Second)
					continue
				}
				_ = ch.Qos(rm.Qos, 0, false)

				msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
				if err != nil {
					ch.Close()
					time.Sleep(2 * time.Second)
					continue
				}

				for d := range msgs {
					handler(&d, ch)
				}

				ch.Close()
				if rm.IsClosing {
					break
				}
			}
		}(i)
	}
}

func (rm *RabbitManager) DirectReplyTo(ctx context.Context, exchange, routingKey string, body []byte, corid string) ([]byte, error) {
	rm.Mu.RLock()
	ready := rm.IsReady
	rm.Mu.RUnlock()
	if !ready {
		return nil, fmt.Errorf("publisher not ready")
	}

	ch := <-rm.ChPool
	defer func() { rm.ChPool <- ch }()

	// 1. Consume from amq.rabbitmq.reply-to
	// Each channel has its own private reply-to pseudo-queue
	consumerTag := uuid.New().String()
	deliveries, err := ch.Consume(
		"amq.rabbitmq.reply-to", // queue
		consumerTag,             // consumer
		true,                    // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume from reply-to: %w", err)
	}

	// Ensure we cancel the consumer when done
	defer ch.Cancel(consumerTag, false)

	corrId := corid
	if corrId == "" {
		corrId = uuid.New().String()
	}

	// 2. Publish the request
	msg := amqp.Publishing{
		ContentType:   "application/json",
		ReplyTo:       "amq.rabbitmq.reply-to",
		CorrelationId: corrId,
		Body:          body,
	}

	err = ch.PublishWithContext(ctx, exchange, routingKey, false, false, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to publish request: %w", err)
	}

	// 3. Wait for response or timeout/cancellation
	for {
		select {
		case d, ok := <-deliveries:
			if !ok {
				return nil, fmt.Errorf("response channel closed")
			}
			if d.CorrelationId == corrId {
				return d.Body, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (rm *RabbitManager) Shutdown() {
	rm.L.Info("Shutting down RabbitMQ gracefully...")
	rm.IsClosing = true
	rm.Wg.Wait() // Wait for all handlers to finish
	rm.Conn.Close()
	rm.L.Info("RabbitMQ clean exit.")
}

func (rm *RabbitManager) DirectReplyToWithRetry(ctx context.Context, exchange, routingKey string, body []byte, corid string) ([]byte, error) {
	var lastErr error

	// Attempt 3 retries
	for i := 0; i < 3; i++ {
		res, err := rm.DirectReplyTo(ctx, exchange, routingKey, body, corid)
		if err == nil {
			return res, nil // Success!
		}

		lastErr = err
		rm.L.Warnf("DirectReplyTo failed: %v. Retrying in %d ms... (Attempt %d/3)", err, (i+1)*500, i+1)

		// Wait before trying again (Exponential backoff)
		select {
		case <-time.After(time.Duration(i+1) * 500 * time.Millisecond):
			continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("DirectReplyTo failed after 3 retries: %w", lastErr)
}

type RabbitConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Vhost    string
	PoolSize int
	Qos      int
	Declares []RabbitDeclare
}

func InitMessageBroker(cfg RabbitConfig, l *logrus.Logger) *RabbitManager {
	url := amqp.URI{
		Scheme:   "amqp",
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.User,
		Password: cfg.Password,
		Vhost:    cfg.Vhost,
	}.String()
	redactedURL := fmt.Sprintf("amqp://%s:****@%s:%d%s", cfg.User, cfg.Host, cfg.Port, cfg.Vhost)
	l.Infof("🔌 Attempting to connect to RabbitMQ: %s", redactedURL)

	rm, err := NewRabbitManager(l, url, cfg.PoolSize, cfg.Declares, cfg.Qos)
	if err != nil {
		l.Errorf("Fatal Error Connection RabbitMQ : %#v\n", err)
		panic(err)
	}

	return rm
}
