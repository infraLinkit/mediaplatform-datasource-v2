package handler

import (
	"github.com/go-redis/redis"
	"github.com/gofiber/storage/rueidis"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/messaging"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/persistence"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm"
)

type (
	IncomingHandler struct {
		Config *config.Cfg
		Logs   *logrus.Logger
		DB     *gorm.DB
		Rmqp   rmqp.AMQP
		RM     *messaging.RabbitManager
		R      *rueidis.Storage
		RCP    *redis.Client
		DS     *persistence.BaseModel
		GS     *sheets.Service
	}
)

func NewIncomingHandler(obj IncomingHandler) *IncomingHandler {

	b := persistence.NewBaseModel(persistence.BaseModel{
		Config: obj.Config,
		Logs:   obj.Logs,
		DB:     obj.DB,
		R:      obj.R,
		RCP:    obj.RCP,
	})

	return &IncomingHandler{
		Config: obj.Config,
		Logs:   obj.Logs,
		DB:     obj.DB,
		R:      obj.R,
		RCP:    obj.RCP,
		Rmqp:   obj.Rmqp,
		RM:     obj.RM,
		DS:     b,
		GS:     obj.GS,
	}
}
