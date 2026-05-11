package persistence

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofiber/storage/rueidis"
)

// SafeRedisGet: nil-safe Get. Returns (value, found).
// found=false jika client nil ATAU key miss ATAU error → caller fallback ke DB.
func SafeRedisGet(c *redis.Client, key string) (string, bool) {
	if c == nil {
		return "", false
	}
	res := c.Get(key)
	if res == nil || res.Err() != nil {
		return "", false
	}
	v := res.Val()
	if v == "" {
		return "", false
	}
	return v, true
}

// SafeRedisSet: nil-safe Set. Returns true jika sukses, false kalau client nil/error.
func SafeRedisSet(c *redis.Client, key, value string, ttl time.Duration) bool {
	if c == nil {
		return false
	}
	return c.Set(key, value, ttl).Err() == nil
}

// SafeRedisDel: nil-safe Del. No-op kalau client nil.
func SafeRedisDel(c *redis.Client, keys ...string) {
	if c == nil || len(keys) == 0 {
		return
	}
	c.Del(keys...)
}

// SafeRueidisJSONGet: nil-safe rueidis JSON.GET. Returns (raw, found).
func SafeRueidisJSONGet(s *rueidis.Storage, ctx context.Context, key, path string) (string, bool) {
	if s == nil {
		return "", false
	}
	conn := s.Conn()
	if conn == nil {
		return "", false
	}
	res, err := conn.Do(ctx, conn.B().JsonGet().Key(key).Path(path).Build()).ToString()
	if err != nil || res == "" {
		return "", false
	}
	return res, true
}

// SafeRueidisExists: nil-safe EXISTS check.
func SafeRueidisExists(s *rueidis.Storage, ctx context.Context, key string) bool {
	if s == nil {
		return false
	}
	conn := s.Conn()
	if conn == nil {
		return false
	}
	n, err := conn.Do(ctx, conn.B().Exists().Key(key).Build()).ToInt64()
	if err != nil {
		return false
	}
	return n > 0
}
