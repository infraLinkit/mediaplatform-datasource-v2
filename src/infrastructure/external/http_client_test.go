package external

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// === HttpClient TLS config ===

func TestHttpClient_DefaultTLSSecure(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	c := HttpClient(PHttp{
		DialTimeout:        1,
		Timeout:            5,
		KeepAlive:          1,
		IsDisableKeepAlive: true,
		MaxIdleConns:       10,
		IdleConnTimeout:    1 * time.Second,
		InsecureSkipVerify: false,
	})
	require.NotNil(t, c)
	tr, ok := c.Transport.(*http.Transport)
	require.True(t, ok)
	assert.False(t, tr.TLSClientConfig.InsecureSkipVerify, "default harus secure")
	assert.GreaterOrEqual(t, tr.TLSClientConfig.MinVersion, uint16(0x0303), "min TLS 1.2")
}

func TestHttpClient_OptInSkipVerify_NonProd(t *testing.T) {
	t.Setenv("APP_ENV", "dev")
	c := HttpClient(PHttp{InsecureSkipVerify: true, Timeout: 5})
	tr := c.Transport.(*http.Transport)
	assert.True(t, tr.TLSClientConfig.InsecureSkipVerify, "non-prod opt-in skip verify allowed")
}

func TestHttpClient_OptInSkipVerify_BlockedInProd(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	c := HttpClient(PHttp{InsecureSkipVerify: true, Timeout: 5})
	tr := c.Transport.(*http.Transport)
	assert.False(t, tr.TLSClientConfig.InsecureSkipVerify, "prod tidak boleh skip verify walau opt-in")
}

// === Get / Post HTTP ===

func TestGet_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	body, status, code, _, err := Get(srv.URL, nil, PHttp{Timeout: 5, DialTimeout: 1, MaxIdleConns: 1, IdleConnTimeout: 1 * time.Second})
	require.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Contains(t, status, "200")
	assert.Equal(t, `{"ok":true}`, string(body))
}

func TestGet_Non200_ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, _, code, _, err := Get(srv.URL, nil, PHttp{Timeout: 5, DialTimeout: 1})
	assert.Error(t, err, "non-OK status harus return error (post-fix #2)")
	assert.Equal(t, 500, code)
}

func TestGet_UnreachableHost_NoNilPanic(t *testing.T) {
	// Hit reserved test address yang gak dijalanin server. Pastikan tidak panic walau response nil.
	assert.NotPanics(t, func() {
		_, _, code, _, err := Get("http://127.0.0.1:1/never", nil, PHttp{Timeout: 1, DialTimeout: 1})
		assert.Error(t, err)
		assert.Equal(t, 0, code)
	})
}

func TestPost_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"received":true}`))
	}))
	defer srv.Close()

	body, _, code, _, err := Post(srv.URL, map[string]string{"X-Test": "1"}, []byte(`{"a":1}`), PHttp{Timeout: 5, DialTimeout: 1})
	require.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Contains(t, string(body), "received")
}

func TestPost_UnreachableHost_NoNilPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		_, _, code, _, err := Post("http://127.0.0.1:1/never", nil, []byte(`{}`), PHttp{Timeout: 1, DialTimeout: 1})
		assert.Error(t, err)
		assert.Equal(t, 0, code)
	})
}

// === Block try/catch ===

func TestBlock_TryNoError_FinallyRuns(t *testing.T) {
	finallyRan := false
	caughtVal := ""

	Block{
		Try: func() {
			// no panic
		},
		Catch: func(e Exception) {
			caughtVal = "should-not-run"
		},
		Finally: func() {
			finallyRan = true
		},
	}.Do()

	assert.True(t, finallyRan)
	assert.Empty(t, caughtVal)
}

func TestBlock_PanicCaught(t *testing.T) {
	caught := false
	Block{
		Try: func() {
			panic("boom")
		},
		Catch: func(e Exception) {
			caught = true
			assert.Equal(t, "boom", e)
		},
	}.Do()
	assert.True(t, caught, "Catch harus run pas Try panic")
}
