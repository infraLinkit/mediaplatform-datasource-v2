package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Redis nil-safe wrappers harus tidak panic kalau client nil.
// Test ini paksa nil input — degraded mode #4.

func TestSafeRedisGet_NilClient(t *testing.T) {
	val, ok := SafeRedisGet(nil, "any-key")
	assert.False(t, ok)
	assert.Empty(t, val)
}

func TestSafeRedisSet_NilClient(t *testing.T) {
	ok := SafeRedisSet(nil, "k", "v", 30*time.Second)
	assert.False(t, ok)
}

func TestSafeRedisDel_NilClient_NoPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		SafeRedisDel(nil, "k1", "k2")
	})
}

func TestSafeRedisDel_EmptyKeys(t *testing.T) {
	assert.NotPanics(t, func() {
		SafeRedisDel(nil)
	})
}

func TestSafeRueidisJSONGet_NilClient(t *testing.T) {
	val, ok := SafeRueidisJSONGet(nil, context.Background(), "k", "$")
	assert.False(t, ok)
	assert.Empty(t, val)
}

func TestSafeRueidisExists_NilClient(t *testing.T) {
	exists := SafeRueidisExists(nil, context.Background(), "k")
	assert.False(t, exists)
}
