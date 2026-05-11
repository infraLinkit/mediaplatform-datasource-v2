package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnvIntDefault_Unset(t *testing.T) {
	t.Setenv("MY_TEST_VAR_X", "")
	assert.Equal(t, 42, envIntDefault("MY_TEST_VAR_X", 42))
}

func TestEnvIntDefault_Valid(t *testing.T) {
	t.Setenv("MY_TEST_VAR_X", "100")
	assert.Equal(t, 100, envIntDefault("MY_TEST_VAR_X", 42))
}

func TestEnvIntDefault_InvalidFallback(t *testing.T) {
	t.Setenv("MY_TEST_VAR_X", "not-a-number")
	assert.Equal(t, 42, envIntDefault("MY_TEST_VAR_X", 42), "invalid value harus fallback default")
}

// === Cfg.Redacted ===

func TestRedacted_MaskSecrets(t *testing.T) {
	c := &Cfg{
		PSQLPassword:     "supersecret",
		RedisPwd:         "redispass",
		RabbitMQPassword: "rmqpass",
		ARPUUsername:     "arpuser",
		ARPUPassword:     "arpsecret",
		GSPrivateKey:     "-----BEGIN PRIVATE KEY-----\nDATA\n-----END",
		GSPrivateKeyID:   "key-id-123",
		GSClientID:       "client-id-456",
		// public fields tetap visible
		AppHost:   "example.com",
		PSQLHost:  "db.host",
		RedisHost: "redis.host",
	}

	red := c.Redacted()

	// secrets masked
	assert.Equal(t, "***REDACTED***", red.PSQLPassword)
	assert.Equal(t, "***REDACTED***", red.RedisPwd)
	assert.Equal(t, "***REDACTED***", red.RabbitMQPassword)
	assert.Equal(t, "***REDACTED***", red.ARPUUsername)
	assert.Equal(t, "***REDACTED***", red.ARPUPassword)
	assert.Equal(t, "***REDACTED***", red.GSPrivateKey)
	assert.Equal(t, "***REDACTED***", red.GSPrivateKeyID)
	assert.Equal(t, "***REDACTED***", red.GSClientID)

	// public visible
	assert.Equal(t, "example.com", red.AppHost)
	assert.Equal(t, "db.host", red.PSQLHost)
	assert.Equal(t, "redis.host", red.RedisHost)

	// original tidak berubah
	assert.Equal(t, "supersecret", c.PSQLPassword)
}

func TestRedacted_EmptyFieldsStaysEmpty(t *testing.T) {
	c := &Cfg{} // all zero
	red := c.Redacted()
	// empty stays empty (mask hanya kalau non-empty)
	assert.Empty(t, red.PSQLPassword)
	assert.Empty(t, red.RedisPwd)
}

// === InitCfg env defaults ===

func TestInitCfg_DBPoolDefaults(t *testing.T) {
	t.Setenv("DB_MAX_IDLE_CONNS", "")
	t.Setenv("DB_MAX_OPEN_CONNS", "")
	t.Setenv("DB_CONN_MAX_LIFETIME_MIN", "")
	t.Setenv("DB_CONN_MAX_IDLE_TIME_MIN", "")

	cfg := InitCfg()
	assert.Equal(t, 10, cfg.DBMaxIdleConns)
	assert.Equal(t, 100, cfg.DBMaxOpenConns)
	assert.Equal(t, 30*time.Minute, cfg.DBConnMaxLifetime)
	assert.Equal(t, 10*time.Minute, cfg.DBConnMaxIdleTime)
}

func TestInitCfg_DBPoolFromEnv(t *testing.T) {
	t.Setenv("DB_MAX_IDLE_CONNS", "5")
	t.Setenv("DB_MAX_OPEN_CONNS", "50")
	t.Setenv("DB_CONN_MAX_LIFETIME_MIN", "60")
	t.Setenv("DB_CONN_MAX_IDLE_TIME_MIN", "20")

	cfg := InitCfg()
	assert.Equal(t, 5, cfg.DBMaxIdleConns)
	assert.Equal(t, 50, cfg.DBMaxOpenConns)
	assert.Equal(t, 60*time.Minute, cfg.DBConnMaxLifetime)
	assert.Equal(t, 20*time.Minute, cfg.DBConnMaxIdleTime)
}

func TestInitCfg_RedisRequiredDefault(t *testing.T) {
	t.Setenv("REDIS_REQUIRED", "")
	cfg := InitCfg()
	assert.True(t, cfg.RedisRequired, "default harus true (fail-fast)")
}

func TestInitCfg_RedisRequiredFalse(t *testing.T) {
	t.Setenv("REDIS_REQUIRED", "false")
	cfg := InitCfg()
	assert.False(t, cfg.RedisRequired)
}
