package http

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/interfaces/http/routes"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/interfaces/http/handler"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
	_ "github.com/lib/pq"
	"github.com/mikhail-bigun/fiberlogrus"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm"
)

// authEnforceDefault: kalau true, semua route protected groups (dashboard/report/internal/management)
// wajib auth middleware. Default false untuk back-compat sementara FE migrate.
func authEnforceDefault() bool {
	v := os.Getenv("AUTH_ENFORCE_DEFAULT")
	if v == "" {
		return false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}

type App3rdParty struct {
	Config *config.Cfg
	Logs   *logrus.Logger
	DB     *gorm.DB
	R      *rueidis.Storage
	RCP    *redis.Client
	Rmqp   rmqp.AMQP
	GS     *sheets.Service
}

func MapUrls(obj App3rdParty) *fiber.App {

	f := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		BodyLimit:   100 * 1024 * 1024,
	})

	f.Use(
		fiberlogrus.New(
			fiberlogrus.Config{
				Logger: external.MakeLogger(
					external.Setup{
						Env:     obj.Config.LogEnv,
						Logname: obj.Config.LogPath + "/access_log",
						Display: true,
						Level:   obj.Config.LogLevel,
					}),
				Tags: []string{
					fiberlogrus.TagIP,
					fiberlogrus.TagIPs,
					fiberlogrus.TagProtocol,
					fiberlogrus.TagHost,
					fiberlogrus.TagPort,
					fiberlogrus.TagMethod,
					fiberlogrus.TagPath,
					fiberlogrus.TagURL,
					fiberlogrus.TagUA,
					fiberlogrus.TagBody,
					fiberlogrus.TagRoute,
					fiberlogrus.TagQueryStringParams,
					fiberlogrus.TagStatus,
					fiberlogrus.TagPid,
					fiberlogrus.TagReferer,
					fiberlogrus.TagLatency,
				},
			}))

	h := handler.NewIncomingHandler(handler.IncomingHandler{
		Config: obj.Config,
		Logs:   obj.Logs,
		R:      obj.R,
		RCP:    obj.RCP,
		DB:     obj.DB,
		Rmqp:   obj.Rmqp,
		GS:     obj.GS,
	})

	// Group-level auth ditarget kalau AUTH_ENFORCE_DEFAULT=true.
	authMW := func(c *fiber.Ctx) error { return c.Next() } // no-op fallback
	if authEnforceDefault() {
		authMW = h.AuthMiddleware
	}

	v1 := f.Group("/v1")

	// Public
	routes.RegisterPostback(v1, h)
	v1.Group("/ext") // External API placeholder

	// Auth-protected
	routes.RegisterDashboard(f.Group("/dashboard", authMW), h)
	routes.RegisterReport(v1.Group("/report", authMW), h)
	routes.RegisterInternal(v1.Group("/int", authMW), h)
	routes.RegisterManagement(v1.Group("/management", authMW), h)

	return f
}
