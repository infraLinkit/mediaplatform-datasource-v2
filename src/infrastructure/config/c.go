package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/messaging"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var APP_PATH = "./" // default relatif; override via env APPPATH

func envIntDefault(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

const (
	//APP_PATH       = "/app"
	EVENT_TRAFFIC  = "Traffic"
	EVENT_LANDING  = "Landing"
	EVENT_CLICK    = "Click Landing"
	EVENT_REDIRECT = "Redirect"

	OK_DESC          = "OK"
	BAD_REQUEST_DESC = "Bad Request"
)

type (
	Cfg struct {
		AppDockerVer                           string
		AppHost                                string
		AppHostPort                            string
		AppApi                                 string
		AppApiPort                             string
		TrustedProxyHeader                     bool
		ReduceMemoryUsage                      bool
		ConcurrencyConnection                  int
		RedisHost                              string
		RedisPort                              int
		RedisDBIndex                           int
		RedisCachePixel                        int
		RedisPwd                               string
		RedisKeyExpiration                     int64
		PSQLHost                               string
		PSQLUsername                           string
		PSQLPassword                           string
		PSQLPort                               string
		PSQLDB                                 string
		RabbitMQHost                           string
		RabbitMQPort                           int
		RabbitMQUsername                       string
		RabbitMQPassword                       string
		RabbitMQVHost                          string
		RabbitMQDataType                       string
		RabbitMQCtxTimeout                     int
		RabbitMQPixelStorageExchangeName       string
		RabbitMQPixelStorageQueueName          string
		RabbitMQClickStorageExchangeName       string
		RabbitMQClickStorageQueueName          string
		RabbitMQRedisCounterExchangeName       string
		RabbitMQRedisCounterQueueName          string
		RabbitMQRatioExchangeName              string
		RabbitMQRatioQueueName                 string
		RabbitMQRatioQueueThreshold            int
		RabbitMQPostbackAdnetExchangeName      string
		RabbitMQPostbackAdnetQueueName         string
		RabbitMQCampaignManagementExchangeName string
		RabbitMQCampaignManagementQueueName    string
		RabbitMQAlertManagementExchangeName    string
		RabbitMQAlertManagementQueueName       string
		LogEnv                                 string
		LogPath                                string
		LogLevel                               string
		TZ                                     *time.Location
		APIARPU                                string
		ARPUUsername                           string
		ARPUPassword                           string
		APILINKITDashboard                     string
		SendToLinkitDashboard                  bool
		GetDataArpu                            bool
		GetSuccessFP                           bool
		GSType                                 string
		GSProjectID                            string
		GSPrivateKeyID                         string
		GSPrivateKey                           string
		GSClientEmail                          string
		GSClientID                             string
		GSAuthURI                              string
		GSTokenURI                             string
		GSAuthProvider                         string
		GSClient                               string
		GSUniversalDomain                      string
		CronResetCapping                       string
		StartGetIntervalDatePXS                int
		EndGetIntervalDatePXS                  int
		DBMaxIdleConns                         int
		DBMaxOpenConns                         int
		DBConnMaxLifetime                      time.Duration
		DBConnMaxIdleTime                      time.Duration
		RedisRequired                          bool
	}

	Setup struct {
		Config         *Cfg
		Logs           *logrus.Logger
		R              *rueidis.Storage
		RCP            *redis.Client
		DB             *gorm.DB
		RM             *messaging.RabbitManager
		GS             *sheets.Service
		RedisAvailable bool
	}
)

func InitCfg() *Cfg {

	if v := os.Getenv("APPPATH"); v != "" {
		APP_PATH = v
	} else {
		fmt.Println("[!] WARN: APPPATH env not set, fallback to './'")
	}

	loc, _ := time.LoadLocation(os.Getenv("TZ"))
	time.Local = loc // Set global location

	rabbitmq_port, _ := strconv.Atoi(os.Getenv("RABBITMQPORT"))
	redis_dbindex, _ := strconv.Atoi(os.Getenv("REDISDBINDEX"))
	redis_cache_pixel, _ := strconv.Atoi(os.Getenv("REDISCACHEPIXEL"))
	redis_port, _ := strconv.Atoi(os.Getenv("REDISPORT"))
	redis_exp, _ := strconv.Atoi(os.Getenv("REDISKEYEXPIRE"))
	ratio_queue_threshold, _ := strconv.Atoi(os.Getenv("RABBITMQRATIOQUEUETHRESHOLD"))
	trusted_proxy_header, _ := strconv.ParseBool(os.Getenv("TRUSTED_PROXY_HEADER"))
	reduce_memory_usage, _ := strconv.ParseBool(os.Getenv("REDUCE_MEMORY_USAGE"))
	concurrency_connection, _ := strconv.Atoi(os.Getenv("CONCURRENCY_CONNECTION"))
	send_to_linkit_dashboard, _ := strconv.ParseBool(os.Getenv("SEND_TO_LINKIT_DASHBOARD"))
	get_data_arpu, _ := strconv.ParseBool(os.Getenv("GET_DATA_ARPU"))
	get_success_fp, _ := strconv.ParseBool(os.Getenv("GET_SUCCESS_FP"))
	start_get_interval_date_pxs, _ := strconv.Atoi(os.Getenv("STARTGETINTERVALDATEPXS"))
	end_get_interval_date_pxs, _ := strconv.Atoi(os.Getenv("ENDGETINTERVALDATEPXS"))

	db_max_idle := envIntDefault("DB_MAX_IDLE_CONNS", 10)
	db_max_open := envIntDefault("DB_MAX_OPEN_CONNS", 100)
	db_conn_lifetime := time.Duration(envIntDefault("DB_CONN_MAX_LIFETIME_MIN", 30)) * time.Minute
	db_conn_idletime := time.Duration(envIntDefault("DB_CONN_MAX_IDLE_TIME_MIN", 10)) * time.Minute
	rmqpctxtimeout := envIntDefault("RABBITMQCONTEXTTIMEOUT", 30)

	redis_required := true // default: required (fail-fast)
	if v := os.Getenv("REDIS_REQUIRED"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			redis_required = b
		}
	}

	cfg := &Cfg{
		AppDockerVer:                           os.Getenv("APP_DOCKER_VER"),
		AppHost:                                os.Getenv("APPHOST"),
		AppHostPort:                            os.Getenv("APPHOSTPORT"),
		AppApi:                                 os.Getenv("APPAPI"),
		AppApiPort:                             os.Getenv("APPAPIPORT"),
		TrustedProxyHeader:                     trusted_proxy_header,
		ReduceMemoryUsage:                      reduce_memory_usage,
		ConcurrencyConnection:                  concurrency_connection,
		RedisHost:                              os.Getenv("REDISHOST"),
		RedisPort:                              redis_port,
		RedisDBIndex:                           redis_dbindex,
		RedisCachePixel:                        redis_cache_pixel,
		RedisPwd:                               os.Getenv("REDISPASSWORD"),
		RedisKeyExpiration:                     int64(redis_exp),
		PSQLHost:                               os.Getenv("DB_HOST"),
		PSQLUsername:                           os.Getenv("DB_USERNAME"),
		PSQLPassword:                           os.Getenv("DB_PASSWORD"),
		PSQLPort:                               os.Getenv("DB_PORT"),
		PSQLDB:                                 os.Getenv("DB_DATABASE"),
		RabbitMQHost:                           os.Getenv("RABBITMQHOST"),
		RabbitMQPort:                           rabbitmq_port,
		RabbitMQUsername:                       os.Getenv("RABBITMQUSERNAME"),
		RabbitMQPassword:                       os.Getenv("RABBITMQPASSWORD"),
		RabbitMQVHost:                          os.Getenv("RABBITMQVHOST"),
		RabbitMQDataType:                       "application/json",
		RabbitMQCtxTimeout:                     rmqpctxtimeout,
		RabbitMQPixelStorageExchangeName:       os.Getenv("RABBITMQPIXELSTORAGEEXCHANGENAME"),
		RabbitMQPixelStorageQueueName:          os.Getenv("RABBITMQPIXELSTORAGEQUEUENAME"),
		RabbitMQClickStorageExchangeName:       os.Getenv("RABBITMQCLICKSTORAGEEXCHANGENAME"),
		RabbitMQClickStorageQueueName:          os.Getenv("RABBITMQCLICKSTORAGEQUEUENAME"),
		RabbitMQRedisCounterExchangeName:       os.Getenv("RABBITMQREDISCOUNTEREXCHANGENAME"),
		RabbitMQRedisCounterQueueName:          os.Getenv("RABBITMQREDISCOUNTERQUEUENAME"),
		RabbitMQRatioExchangeName:              os.Getenv("RABBITMQRATIOEXCHANGENAME"),
		RabbitMQRatioQueueName:                 os.Getenv("RABBITMQRATIOQUEUENAME"),
		RabbitMQRatioQueueThreshold:            ratio_queue_threshold,
		RabbitMQPostbackAdnetExchangeName:      os.Getenv("RABBITMQPOSTBACKADNETEXCHANGENAME"),
		RabbitMQPostbackAdnetQueueName:         os.Getenv("RABBITMQPOSTBACKADNETQUEUENAME"),
		RabbitMQCampaignManagementExchangeName: os.Getenv("RABBITMQCAMPAIGNMANAGEMENTEXCHANGENAME"),
		RabbitMQCampaignManagementQueueName:    os.Getenv("RABBITMQCAMPAIGNMANAGEMENTQUEUENAME"),
		RabbitMQAlertManagementExchangeName:    os.Getenv("RABBITMQALERTMANAGEMENTEXCHANGENAME"),
		RabbitMQAlertManagementQueueName:       os.Getenv("RABBITMQALERTMANAGEMENTQUEUENAME"),
		LogEnv:                                 os.Getenv("LOGENV"),
		LogPath:                                os.Getenv("LOGPATH"),
		LogLevel:                               os.Getenv("LOGLEVEL"),
		TZ:                                     loc,
		APIARPU:                                os.Getenv("APIARPU"),
		ARPUUsername:                           os.Getenv("ARPUUsername"),
		ARPUPassword:                           os.Getenv("ARPUPassword"),
		APILINKITDashboard:                     os.Getenv("APILINKITDashboard"),
		SendToLinkitDashboard:                  send_to_linkit_dashboard,
		GetDataArpu:                            get_data_arpu,
		GetSuccessFP:                           get_success_fp,
		GSType:                                 os.Getenv("GSTYPE"),
		GSProjectID:                            os.Getenv("GSPROJECT_ID"),
		GSPrivateKeyID:                         os.Getenv("GSPRIVATE_KEY_ID"),
		GSPrivateKey:                           os.Getenv("GSPRIVATE_KEY"),
		GSClientEmail:                          os.Getenv("GSCLIENT_EMAIL"),
		GSClientID:                             os.Getenv("GSCLIENT_ID"),
		GSAuthURI:                              os.Getenv("GSAUTH_URI"),
		GSTokenURI:                             os.Getenv("GSTOKEN_URI"),
		GSAuthProvider:                         os.Getenv("GSAUTH_PROVIDER"),
		GSClient:                               os.Getenv("GSCLIENT"),
		GSUniversalDomain:                      os.Getenv("GSUNIVERSAL_DOMAIN"),
		CronResetCapping:                       os.Getenv("CRONRESETCAPPING"),
		StartGetIntervalDatePXS:                start_get_interval_date_pxs,
		EndGetIntervalDatePXS:                  end_get_interval_date_pxs,
		DBMaxIdleConns:                         db_max_idle,
		DBMaxOpenConns:                         db_max_open,
		DBConnMaxLifetime:                      db_conn_lifetime,
		DBConnMaxIdleTime:                      db_conn_idletime,
		RedisRequired:                          redis_required,
	}

	return cfg
}

// Redacted: copy Cfg dgn secrets disensor untuk safe logging.
func (c *Cfg) Redacted() Cfg {
	cp := *c
	mask := func(s string) string {
		if s == "" {
			return ""
		}
		return "***REDACTED***"
	}
	cp.PSQLPassword = mask(cp.PSQLPassword)
	cp.RedisPwd = mask(cp.RedisPwd)
	cp.RabbitMQPassword = mask(cp.RabbitMQPassword)
	cp.ARPUUsername = mask(cp.ARPUUsername)
	cp.ARPUPassword = mask(cp.ARPUPassword)
	cp.GSPrivateKey = mask(cp.GSPrivateKey)
	cp.GSPrivateKeyID = mask(cp.GSPrivateKeyID)
	cp.GSClientID = mask(cp.GSClientID)
	return cp
}

func (c *Cfg) Initiate(logname string) (*Setup, error) {

	l := external.MakeLogger(
		external.Setup{Env: c.LogEnv, Logname: c.LogPath + "/" + logname, Display: true, Level: c.LogLevel})
	red := c.Redacted()
	l.Info(fmt.Sprintf("Config Loaded : %#v\n", &red))

	redisAvailable := true
	rj, errRJ := c.InitRedisJSON(l, c.RedisDBIndex)
	rcp, errRCP := c.InitRedis(l, c.RedisCachePixel)

	if errRJ != nil || errRCP != nil {
		if c.RedisRequired {
			if errRJ != nil {
				return nil, fmt.Errorf("init redis json: %w", errRJ)
			}
			return nil, fmt.Errorf("init redis: %w", errRCP)
		}
		// Degraded mode: log + continue dengan nil client(s)
		redisAvailable = false
		l.Warnf("[!] DEGRADED MODE: Redis unavailable, app jalan tanpa cache. errRJ=%v errRCP=%v", errRJ, errRCP)
		// rj/rcp bisa nil — caller wajib cek RedisAvailable
	}

	rmDeclares := []messaging.RabbitDeclare{
		{ExchangeName: c.RabbitMQPixelStorageExchangeName, ExchangeType: "direct", QueueName: c.RabbitMQPixelStorageQueueName, RoutingKey: c.RabbitMQPixelStorageQueueName, Durable: true},
		{ExchangeName: c.RabbitMQClickStorageExchangeName, ExchangeType: "direct", QueueName: c.RabbitMQClickStorageQueueName, RoutingKey: c.RabbitMQClickStorageQueueName, Durable: true},
		{ExchangeName: c.RabbitMQRedisCounterExchangeName, ExchangeType: "direct", QueueName: c.RabbitMQRedisCounterQueueName, RoutingKey: c.RabbitMQRedisCounterQueueName, Durable: true},
		{ExchangeName: c.RabbitMQRatioExchangeName, ExchangeType: "direct", QueueName: c.RabbitMQRatioQueueName, RoutingKey: c.RabbitMQRatioQueueName, Durable: true},
		{ExchangeName: c.RabbitMQPostbackAdnetExchangeName, ExchangeType: "direct", QueueName: c.RabbitMQPostbackAdnetQueueName, RoutingKey: c.RabbitMQPostbackAdnetQueueName, Durable: true},
		{ExchangeName: c.RabbitMQCampaignManagementExchangeName, ExchangeType: "direct", QueueName: c.RabbitMQCampaignManagementQueueName, RoutingKey: c.RabbitMQCampaignManagementQueueName, Durable: true},
		{ExchangeName: c.RabbitMQAlertManagementExchangeName, ExchangeType: "direct", QueueName: c.RabbitMQAlertManagementQueueName, RoutingKey: c.RabbitMQAlertManagementQueueName, Durable: true},
	}

	rm := messaging.InitMessageBroker(messaging.RabbitConfig{
		Host:     c.RabbitMQHost,
		Port:     c.RabbitMQPort,
		User:     c.RabbitMQUsername,
		Password: c.RabbitMQPassword,
		Vhost:    c.RabbitMQVHost,
		PoolSize: 5,
		Qos:      10,
		Declares: rmDeclares,
	}, l)

	return &Setup{
		Config:         c,
		Logs:           l,
		R:              rj,
		RCP:            rcp,
		DB:             c.InitGormPgx(l),
		RM:             rm,
		GS:             c.InitGoogleSheet(l),
		RedisAvailable: redisAvailable,
	}, nil

}

// retryWithBackoff: 5 attempts, exp backoff 1s/2s/4s/8s/16s
func retryWithBackoff(l *logrus.Logger, name string, fn func() error) error {
	var err error
	delay := 1 * time.Second
	for i := 1; i <= 5; i++ {
		if err = fn(); err == nil {
			return nil
		}
		l.Warnf("[x] %s attempt %d/5 failed: %v, retry in %s", name, i, err, delay)
		if i < 5 {
			time.Sleep(delay)
			delay *= 2
		}
	}
	return fmt.Errorf("%s: all retries exhausted: %w", name, err)
}

func (c *Cfg) InitRedis(l *logrus.Logger, dbindex int) (*redis.Client, error) {
	var r *redis.Client
	err := retryWithBackoff(l, "redis-go", func() error {
		r = redis.NewClient(&redis.Options{
			Addr:     c.RedisHost + ":" + strconv.Itoa(c.RedisPort),
			Password: c.RedisPwd,
			DB:       dbindex,
		})
		pong, e := r.Ping().Result()
		if e != nil {
			return e
		}
		if pong != "PONG" {
			return fmt.Errorf("unexpected ping response: %q", pong)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	l.Info("[v] redis conn established")
	return r, nil
}

func (c *Cfg) InitRedisJSON(l *logrus.Logger, dbindex int) (*rueidis.Storage, error) {
	port := strconv.Itoa(c.RedisPort)
	var r *rueidis.Storage

	err := retryWithBackoff(l, "redis-rueidis", func() (e error) {
		// rueidis.New panics internally on conn fail; recover to make retryable
		defer func() {
			if rec := recover(); rec != nil {
				e = fmt.Errorf("rueidis.New panic: %v", rec)
			}
		}()
		r = rueidis.New(rueidis.Config{
			InitAddress: []string{c.RedisHost + ":" + port},
			Username:    "",
			Password:    c.RedisPwd,
			SelectDB:    dbindex,
			Reset:       false,
			TLSConfig:   nil,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	l.Info("[v] rueidis conn established")
	return r, nil
}

func (c *Cfg) InitPsql(l *logrus.Logger) *sql.DB {

	db, err := sql.Open("postgres", "postgresql://"+c.PSQLUsername+":"+c.PSQLPassword+"@"+c.PSQLHost+":"+c.PSQLPort+"/"+c.PSQLDB+"?sslmode=disable")
	if err != nil {

		// panic the function then hard exit
		l.Info(fmt.Sprintf("[x] An Error occured when establishing of the database : %#v\n", err))

		panic(err)

	} else {

		l.Info("[v] Database successful established\n")
	}

	return db
}

func (c *Cfg) InitGormPgx(l *logrus.Logger) *gorm.DB {

	/* dsn := "host=" + c.PSQLHost + " user=" + c.PSQLUsername + " password=" + c.PSQLPassword + " dbname=" + c.PSQLDB + " port=" + c.PSQLPort + " sslmode=disable TimeZone=" + c.TZ.String()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) */

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=" + c.PSQLHost + " user=" + c.PSQLUsername + " password=" + c.PSQLPassword + " dbname=" + c.PSQLDB + " port=" + c.PSQLPort + " sslmode=disable TimeZone=" + c.TZ.String(),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {

		// panic the function then hard exit
		l.Info(fmt.Sprintf("[x] An Error occured when establishing of the database : %#v\n", err))

		panic(err)

	} else {

		l.Info("[v] Database GORM successful established\n")

		sqlDB, sqlErr := db.DB()
		if sqlErr != nil {
			l.Errorf("[x] failed to get sql.DB: %v", sqlErr)
			panic(sqlErr)
		}

		sqlDB.SetMaxIdleConns(c.DBMaxIdleConns)
		sqlDB.SetMaxOpenConns(c.DBMaxOpenConns)
		sqlDB.SetConnMaxLifetime(c.DBConnMaxLifetime)
		sqlDB.SetConnMaxIdleTime(c.DBConnMaxIdleTime)

		l.Infof("[v] DB pool: maxOpen=%d, maxIdle=%d, lifetime=%s, idleTime=%s",
			c.DBMaxOpenConns, c.DBMaxIdleConns, c.DBConnMaxLifetime, c.DBConnMaxIdleTime)
	}

	return db
}


func (c *Cfg) InitGoogleSheet(l *logrus.Logger) *sheets.Service {
	type KeyGS struct {
		Type            string `json:"type"`
		ProjectID       string `json:"project_id"`
		PrivateKeyID    string `json:"private_key_id"`
		PrivateKey      string `json:"private_key"`
		ClientEmail     string `json:"client_email"`
		ClientID        string `json:"client_id"`
		AuthURI         string `json:"auth_uri"`
		TokenURI        string `json:"token_uri"`
		AuthProvider    string `json:"auth_provider_x509_cert_url"`
		Client          string `json:"client_x509_cert_url"`
		UniversalDomain string `json:"universal_domain"`
	}

	var (
		sheetKey = KeyGS{
			Type:            c.GSType,
			ProjectID:       c.GSProjectID,
			PrivateKeyID:    c.GSPrivateKeyID,
			PrivateKey:      strings.ReplaceAll(c.GSPrivateKey, `\n`, "\n"),
			ClientEmail:     c.GSClientEmail,
			ClientID:        c.GSClientID,
			AuthURI:         c.GSAuthURI,
			TokenURI:        c.GSTokenURI,
			AuthProvider:    c.GSAuthProvider,
			Client:          c.GSClient,
			UniversalDomain: c.GSUniversalDomain,
		}
	)

	credential, err := json.Marshal(sheetKey)
	if err != nil {
		l.Fatalf("Google Sheet Failed to get key: %v", err)
	}

	srv, err := sheets.NewService(context.Background(), option.WithCredentialsJSON(credential))
	if err != nil {
		l.Fatalf("Unable to retrieve Google Sheets client: %v", err)
	}

	l.Info("[v] Google sheet connection successful established\n")

	return srv
}
