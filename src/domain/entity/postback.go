package entity

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
	"github.com/sirupsen/logrus"
)

type (
	PostbackReceive struct {
		CookieKey     string
		URLServiceKey string `json:"urlservicekey"`
		AffSub        string `json:"aff_sub"`
		Method        string `json:"method"`
		Msisdn        string `json:"msisdn"`
		Trxid         string `json:"trxid"`
		SubKeyword    string `json:"subkeyword"`
		Status        string `json:"status"`
		StatusCode    string `json:"statuscode"`
		StatusDetail  string `json:"statusdetail"`
		Desc          string `json:"desc"`
	}

	PostbackData struct {
		CmpDetail DataConfig
		Pxs       PixelStorage
	}
)

func NewDataPostback(c *fiber.Ctx) *PostbackReceive {

	m := c.Queries()

	CookieKey := external.Concat("-", external.GetIpAddress(c), m["urlservicekey"], m["aff_sub"])

	return &PostbackReceive{
		CookieKey:     CookieKey,
		URLServiceKey: m["urlservicekey"],
		AffSub:        m["aff_sub"],
		Msisdn:        m["msisdn"],
		Trxid:         m["trxid"],
		SubKeyword:    m["subkeyword"],
		Status:        m["status"],
		StatusCode:    m["statuscode"],
		StatusDetail:  m["statusdetail"],
	}
}

func (p *PostbackReceive) ValidateParams(Logs *logrus.Logger) GlobalResponse {

	if p.URLServiceKey == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "urlservicekey empty or not found"}

	} else if p.AffSub == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "pixel empty"}

	} else {
		Logs.Debug("All traffic service is valid ...\n")

		return GlobalResponse{Code: fiber.StatusOK, Message: ""}
	}
}

func NewDataPostbackV2(c *fiber.Ctx) *PostbackReceive {

	m := c.Queries()

	CookieKey := external.Concat("-", external.GetIpAddress(c), m["aff_sub"])

	return &PostbackReceive{
		CookieKey:    CookieKey,
		AffSub:       m["aff_sub"],
		Method:       strings.ToUpper(m["method"]),
		Msisdn:       m["msisdn"],
		Trxid:        m["trxid"],
		SubKeyword:   m["subkeyword"],
		Status:       m["status"],
		StatusCode:   m["statuscode"],
		StatusDetail: m["statusdetail"],
		Desc:         m["desc"],
	}
}

func (p *PostbackReceive) ValidateParamsV2(Logs *logrus.Logger) GlobalResponse {

	if p.AffSub == "" {

		return GlobalResponse{Code: fiber.StatusBadRequest, Message: "pixel empty"}

	} else {

		if strings.Contains(p.AffSub, "{pixel}") || strings.Contains(p.AffSub, "[pixel]") {
			return GlobalResponse{Code: fiber.StatusBadRequest, Message: "pixel invalid"}
		} else {
			Logs.Debug("All traffic service is valid ...\n")

			return GlobalResponse{Code: fiber.StatusOK, Message: ""}
		}
	}
}
