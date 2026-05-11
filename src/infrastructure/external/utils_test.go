package external

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// === Encrypt/Decrypt ===

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	t.Setenv("AES_SECRET_KEY", "0123456789abcdef0123456789abcdef") // 32 byte
	plain := "hello-world-rahasia-12345"

	ct, err := Encrypt(plain)
	require.NoError(t, err)
	require.NotEqual(t, plain, ct, "ciphertext should differ from plaintext")

	pt, err := Decrypt(ct)
	require.NoError(t, err)
	assert.Equal(t, plain, pt)
}

func TestEncrypt_NonceRandomness(t *testing.T) {
	t.Setenv("AES_SECRET_KEY", "0123456789abcdef0123456789abcdef")
	plain := "same-input"

	ct1, err := Encrypt(plain)
	require.NoError(t, err)
	ct2, err := Encrypt(plain)
	require.NoError(t, err)

	assert.NotEqual(t, ct1, ct2, "two encrypts of same plaintext must differ (random nonce)")
}

func TestEncrypt_KeyMissing(t *testing.T) {
	os.Unsetenv("AES_SECRET_KEY")
	_, err := Encrypt("anything")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "AES_SECRET_KEY")
}

func TestEncrypt_KeyInvalidLength(t *testing.T) {
	t.Setenv("AES_SECRET_KEY", "tooshort")
	_, err := Encrypt("anything")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "16/24/32")
}

func TestEncrypt_AcceptedKeyLengths(t *testing.T) {
	cases := []struct {
		name string
		key  string
	}{
		{"AES-128", "0123456789abcdef"},                                 // 16
		{"AES-192", "0123456789abcdef01234567"},                         // 24
		{"AES-256", "0123456789abcdef0123456789abcdef"},                 // 32
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("AES_SECRET_KEY", tc.key)
			ct, err := Encrypt("ok")
			require.NoError(t, err)
			pt, err := Decrypt(ct)
			require.NoError(t, err)
			assert.Equal(t, "ok", pt)
		})
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	t.Setenv("AES_SECRET_KEY", "0123456789abcdef0123456789abcdef")
	_, err := Decrypt("not!!!base64!!!")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "base64")
}

func TestDecrypt_TooShort(t *testing.T) {
	t.Setenv("AES_SECRET_KEY", "0123456789abcdef0123456789abcdef")
	_, err := Decrypt("YQ==") // valid base64 but too short for nonce
	require.Error(t, err)
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	t.Setenv("AES_SECRET_KEY", "0123456789abcdef0123456789abcdef")
	ct, err := Encrypt("data")
	require.NoError(t, err)
	// flip last char
	tampered := ct[:len(ct)-1] + "A"
	if tampered == ct {
		tampered = ct[:len(ct)-1] + "B"
	}
	_, err = Decrypt(tampered)
	assert.Error(t, err, "tampered ciphertext should fail GCM auth")
}

// === Concat ===

func TestConcat(t *testing.T) {
	cases := []struct {
		name      string
		splitcode string
		args      []string
		want      string
	}{
		{"empty args", ",", []string{}, ""},
		{"single", ",", []string{"a"}, "a"},
		{"multi", ",", []string{"a", "b", "c"}, "a,b,c"},
		{"pipe", "|", []string{"x", "y"}, "x|y"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Concat(tc.splitcode, tc.args...)
			assert.Equal(t, tc.want, got)
		})
	}
}

// === Time helpers ===

func TestGetFormatTime(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	got := GetFormatTime(loc, "2006-01-02")
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}$`, got)
}

func TestGetCurrentTime(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	// GetCurrentTime parse `now.Format(RFC3339)` pakai val layout.
	// Pakai layout RFC3339 supaya parse sukses → non-zero.
	got := GetCurrentTime(loc, time.RFC3339)
	assert.NotZero(t, got)
}

// === InArray ===

func TestInArray(t *testing.T) {
	haystack := []string{"a", "b", "c"}
	assert.True(t, InArray("b", haystack))
	assert.False(t, InArray("z", haystack))
	assert.False(t, InArray("", haystack))
}
