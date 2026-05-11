package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/sirupsen/logrus"
)

// Smoke test url_mapping — pastikan MapUrls tidak panic + return *fiber.App valid.
// Tidak hit DB/Redis (tidak passing). Cuma cek route registration berjalan.
func TestMapUrls_BuildsApp(t *testing.T) {
	cfg := &config.Cfg{
		LogEnv:   "Test",
		LogPath:  "/tmp",
		LogLevel: "DEBUG",
	}
	logger := logrus.New()

	assert.NotPanics(t, func() {
		f := MapUrls(App3rdParty{
			Config: cfg,
			Logs:   logger,
		})
		require.NotNil(t, f)
	})
}

func TestAuthEnforceDefault_Off(t *testing.T) {
	t.Setenv("AUTH_ENFORCE_DEFAULT", "")
	assert.False(t, authEnforceDefault())
}

func TestAuthEnforceDefault_On(t *testing.T) {
	t.Setenv("AUTH_ENFORCE_DEFAULT", "true")
	assert.True(t, authEnforceDefault())
}

func TestAuthEnforceDefault_InvalidFallback(t *testing.T) {
	t.Setenv("AUTH_ENFORCE_DEFAULT", "yes-please")
	assert.False(t, authEnforceDefault(), "invalid bool harus fallback false")
}

// Route smoke test: hit endpoint public yang tidak butuh DB.
// Pastikan route /v1/postback ter-register (404 = not found, bukan 500).
func TestPublicRoute_NotFoundIsNot500(t *testing.T) {
	cfg := &config.Cfg{LogEnv: "Test", LogPath: "/tmp", LogLevel: "DEBUG"}
	app := MapUrls(App3rdParty{Config: cfg, Logs: logrus.New()})

	req := httptest.NewRequest(http.MethodGet, "/route-yg-tidak-ada", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	// Fiber default 404 untuk unknown route
	assert.Equal(t, 404, resp.StatusCode)
}

func TestPostbackRoute_PathRegistered(t *testing.T) {
	cfg := &config.Cfg{LogEnv: "Test", LogPath: "/tmp", LogLevel: "DEBUG"}
	app := MapUrls(App3rdParty{Config: cfg, Logs: logrus.New()})

	// /v1/postback ada — handler bakal panic karena nil DB, tapi route ke-resolve
	req := httptest.NewRequest(http.MethodGet, "/v1/postback?a=1", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	// Bisa 500 (handler error karena nil DB) atau 200 — yg penting BUKAN 404
	assert.NotEqual(t, 404, resp.StatusCode, "postback path harus ter-register")
	_ = strings.TrimSpace
}
