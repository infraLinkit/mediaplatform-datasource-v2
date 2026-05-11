package entity

/* import (
	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
	"github.com/sirupsen/logrus"
)

type (
	Traffic struct {
		Date          string
		Key           string
		KeyCfg        string
		KeyCounter    string
		KeyDataMining string
		DataTraffic   DataTraffic
	}

	ReceiveTraffic struct {
		URLServiceKey string `form:"urlservicekey" json:"urlservicekey" xml:"urlservicekey"`
		Aff_Sub       string `form:"aff_sub" json:"aff_sub" xml:"aff_sub"`
		PubId         string `form:"pubid" json:"pubid" xml:"pubid"`
	}
)

func NewInstanceTraffic(cfg *config.Cfg, o DataTraffic) *Traffic {

	date := external.GetFormatTime(cfg.TZ, "20060102")
	//key := external.Concat("-", o.URLServiceKey)

	return &Traffic{
		Date:       date,
		Key:        o.URLServiceKey,
		KeyCfg:     external.Concat("-", o.URLServiceKey, "configIdx"),
		KeyCounter: external.Concat("-", o.URLServiceKey, "counterIdx"),
		//KeyDataMining: external.Concat("-", date, key, "dataminingIdx"),
		DataTraffic: o,
	}
}

func (t *Traffic) ValidateParams(Logs *logrus.Logger, traffic *Traffic) GlobalResponse {

	if traffic.DataTraffic.URLServiceKey == "" {
		Logs.Debug("Receive traffic keyaccess param is empty ...\n")

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "parameters is not complete"}
	} else if traffic.DataTraffic.Aff_Sub == "" {
		Logs.Debug("Receive traffic aff_sub param is empty ...\n")

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "parameters is not complete"}
	} else {
		Logs.Debug("All traffic service is valid ...\n")

		return GlobalResponse{Code: fiber.StatusOK, Message: ""}
	}
}
*/
