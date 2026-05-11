package handler

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

// === audienceMatches helper ===

func TestAudienceMatches_String(t *testing.T) {
	assert.True(t, audienceMatches("mediaplatform", "mediaplatform"))
	assert.False(t, audienceMatches("other", "mediaplatform"))
	assert.False(t, audienceMatches("", "mediaplatform"))
}

func TestAudienceMatches_Array(t *testing.T) {
	aud := []interface{}{"mediaplatform", "other"}
	assert.True(t, audienceMatches(aud, "mediaplatform"))
	assert.True(t, audienceMatches(aud, "other"))
	assert.False(t, audienceMatches(aud, "third"))
}

func TestAudienceMatches_Nil(t *testing.T) {
	assert.False(t, audienceMatches(nil, "mediaplatform"))
}

func TestAudienceMatches_WrongType(t *testing.T) {
	assert.False(t, audienceMatches(123, "mediaplatform"))
}

// === Sign sample JWT (helper untuk test) ===

func signTestJWT(secret string, claims jwt.MapClaims) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(secret))
}

func TestSignTestJWT_WorksAndParses(t *testing.T) {
	tokStr, err := signTestJWT("mysecret", jwt.MapClaims{
		"sub":  1,
		"jti":  "abc",
		"type": "access",
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, tokStr)

	tok, err := jwt.Parse(tokStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("mysecret"), nil
	})
	assert.NoError(t, err)
	assert.True(t, tok.Valid)
}

// === RevokeJWT (no-op kalau RCP nil — degraded mode) ===

func TestRevokeJWT_NilRedis_NoError(t *testing.T) {
	h := &IncomingHandler{RCP: nil}
	err := h.RevokeJWT("some-jti", 5*time.Minute)
	assert.NoError(t, err, "harus no-op kalau Redis nil (degraded)")
}

func TestRevokeJWT_EmptyJTI_NoOp(t *testing.T) {
	h := &IncomingHandler{RCP: nil}
	err := h.RevokeJWT("", 5*time.Minute)
	assert.NoError(t, err)
}

// NOTE: Full AuthMiddleware integration test memerlukan mock fiber.Ctx + Redis + DB.
// Test di atas cover unit-level helper (audienceMatches, RevokeJWT degraded path,
// JWT signing roundtrip). Untuk full integration: butuh test infra dengan
// httptest.Server + mocked DB + miniredis (TODO terpisah, lihat TESTING.md).
