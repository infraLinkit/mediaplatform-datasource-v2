package cmd

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/interfaces/http"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/spf13/cobra"
)

type rmqChannelDef struct {
	exchange string
	queue    string
}

// rmqpReconnectWatcher: blocks on connection close notify, then reconnect + re-setup channels.
// Loops forever; jangan exit kecuali process down.
func rmqpReconnectWatcher(c *config.Setup, cfg *config.Cfg, channels []rmqChannelDef) {
	for {
		if c.Rmqp.Connection == nil {
			c.Logs.Warn("[rmq] no connection, waiting 5s before reconnect attempt")
			time.Sleep(5 * time.Second)
			if err := c.Rmqp.SetupConnectionAmqpAndReconnect(); err != nil {
				c.Logs.Errorf("[rmq] reconnect failed: %v", err)
				continue
			}
			for _, ch := range channels {
				if err := c.Rmqp.SetUpChannel("direct", true, ch.exchange, true, ch.queue); err != nil {
					c.Logs.Errorf("[rmq] re-setup channel %s/%s failed: %v", ch.exchange, ch.queue, err)
				}
			}
			c.Logs.Info("[rmq] reconnect + channel re-setup OK")
			continue
		}

		closeCh := make(chan *amqp.Error, 1)
		c.Rmqp.Connection.NotifyClose(closeCh)

		// Block sampai connection close
		err := <-closeCh
		c.Logs.Warnf("[rmq] connection closed: %v, attempting reconnect", err)
		c.Rmqp.Connection = nil
	}
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Webserver CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.InitCfg()
		c, err := cfg.Initiate("api")
		if err != nil {
			log.Fatalf("init setup failed: %v", err)
		}

		// Schema migration dikelola via cmd/migrate.go — server tidak menjalankan AutoMigrate.

		channels := []rmqChannelDef{
			{cfg.RabbitMQPixelStorageExchangeName, cfg.RabbitMQPixelStorageQueueName},
			{cfg.RabbitMQClickStorageExchangeName, cfg.RabbitMQClickStorageQueueName},
			{cfg.RabbitMQRatioExchangeName, cfg.RabbitMQRatioQueueName},
			{cfg.RabbitMQCampaignManagementExchangeName, cfg.RabbitMQCampaignManagementQueueName},
			{"E_RESENDCAMPAIGNDATA", "Q_RESENDCAMPAIGNDATA"},
		}
		for _, ch := range channels {
			if err := c.Rmqp.SetUpChannel("direct", true, ch.exchange, true, ch.queue); err != nil {
				log.Fatalf("rmq setup channel %s/%s failed: %v", ch.exchange, ch.queue, err)
			}
		}

		// Reconnect watcher: monitor connection close, retry + re-setup channels.
		go rmqpReconnectWatcher(c, cfg, channels)

		router := http.MapUrls(http.App3rdParty{
			Config: cfg,
			Logs:   c.Logs,
			DB:     c.DB,
			R:      c.R,
			RCP:    c.RCP,
			Rmqp:   c.Rmqp,
			GS:     c.GS,
		})

		log.Fatal(router.Listen(":" + c.Config.AppApiPort))
	},
}
