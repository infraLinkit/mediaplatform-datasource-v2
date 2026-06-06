package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
)

func (h *IncomingHandler) Postback(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	h.Logs.Debug(fmt.Sprintf("Receive request postback %#v ...\n", c.AllParams()))

	// Parse Postback Data
	p := entity.NewDataPostback(c)
	p.URLServiceKey = c.Params("urlservicekey")

	// Validate Parameters
	if v := p.ValidateParams(h.Logs); v.Code == fiber.StatusBadRequest {

		return c.Status(v.Code).JSON(entity.GlobalResponse{Code: v.Code, Message: v.Message})

	} else {

		if c.Cookies(p.CookieKey) != "" {

			return c.Status(fiber.StatusForbidden).JSON(entity.GlobalResponse{Code: fiber.StatusForbidden, Message: "forbidden access"})

		} else {
			// Setup cookie if double requested within n-hour
			c.Cookie(&fiber.Cookie{
				Name:     p.CookieKey,
				Value:    "1",
				Expires:  time.Now().Add(3 * time.Second),
				HTTPOnly: true,
				SameSite: "lax",
			})

			if dc, err := h.DS.GetDataConfig(external.Concat("-", p.URLServiceKey, "configIdx"), "$"); err == nil {

				pxData := entity.PixelStorage{
					URLServiceKey: p.URLServiceKey, Pixel: p.AffSub}

				var (
					px   entity.PixelStorage
					isPX bool
				)

				if dc.PostbackMethod == "ADNETCODE" {
					px, isPX = h.DS.GetByAdnetCode(pxData)
				} else if dc.PostbackMethod == "TOKEN" {
					px, isPX = h.DS.GetToken(pxData)
				} else if dc.PostbackMethod == "JSON-MSISDN" || dc.PostbackMethod == "XML-MSISDN" || dc.PostbackMethod == "HTML-MSISDN" {
					px, isPX = h.DS.GetPxByMsisdn(pxData)
				} else {
					px, isPX = h.DS.GetPx(pxData)
				}

				if !isPX {

					return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Pixel not found or duplicate used, pixel : " + p.AffSub})

				} else {

					if px.IsUsed {

						return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusNotFound, Message: "NOK - Pixel already used", Data: entity.PixelStorageRsp{
							URLServiceKey: dc.URLServiceKey,
							Adnet:         dc.Adnet,
							IsBillable:    dc.IsBillable,
							Pixel:         px.Pixel,
							Browser:       px.Browser,
							OS:            px.OS,
							Handset:       px.UserAgent,
							PubId:         px.PubId,
							PixelUsedDate: px.PixelUsedDate.Format(time.RFC3339),
						}})

					} else {

						px.PixelUsedDate = external.GetCurrentTime(h.Config.TZ, time.RFC3339)

						bodyReq, _ := json.Marshal(px)

						corId := "RTO" + external.GetUniqId(h.Config.TZ)

						pubCtx, pubCancel := context.WithTimeout(c.UserContext(), time.Duration(h.Config.RabbitMQCtxTimeout)*time.Second)
						defer pubCancel()
						if err := h.RM.PublishWithRetry(pubCtx, h.Config.RabbitMQRatioExchangeName, h.Config.RabbitMQRatioQueueName, bodyReq, corId); err != nil {
							h.Logs.Debug(fmt.Sprintf("[x] Failed published: %s, Data: %s ...", corId, string(bodyReq)))
						} else {
							h.Logs.Debug(fmt.Sprintf("[v] Published: %s, Data: %s ...", corId, string(bodyReq)))
						}

						return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusOK, Message: "OK", Data: entity.PixelStorageRsp{
							URLServiceKey: dc.URLServiceKey,
							Adnet:         dc.Adnet,
							IsBillable:    dc.IsBillable,
							Pixel:         px.Pixel,
							Browser:       px.Browser,
							OS:            px.OS,
							Handset:       px.UserAgent,
							PubId:         px.PubId,
							PixelUsedDate: external.GetFormatTime(h.Config.TZ, time.RFC3339),
						}})
					}
				}

			} else {

				return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Campaign ID not found"})

			}
		}
	}
}

func (h *IncomingHandler) Postback2(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	h.Logs.Debug(fmt.Sprintf("Receive request postback %#v ...\n", c.AllParams()))

	// Parse Postback Data
	p := entity.NewDataPostbackV2(c)
	//p.URLServiceKey = c.Params("urlservicekey")

	// Validate Parameters
	if v := p.ValidateParamsV2(h.Logs); v.Code == fiber.StatusBadRequest {

		return c.Status(v.Code).JSON(entity.GlobalResponse{Code: v.Code, Message: v.Message})

	} else {

		if c.Cookies(p.CookieKey) != "" {

			return c.Status(fiber.StatusForbidden).JSON(entity.GlobalResponse{Code: fiber.StatusForbidden, Message: "forbidden access"})

		} else {
			// Setup cookie if double requested within n-hour
			c.Cookie(&fiber.Cookie{
				Name:     p.CookieKey,
				Value:    "1",
				Expires:  time.Now().Add(3 * time.Second),
				HTTPOnly: true,
				SameSite: "lax",
			})

			if !strings.Contains(p.AffSub, "-") {

				return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Invalid pixel format, pixel : " + p.AffSub})

			} else {

				dataraw := strings.Split(p.AffSub, "-")
				p.URLServiceKey = external.Concat("-", dataraw[0], dataraw[1])

				if dc, err := h.DS.GetDataConfig(external.Concat("-", p.URLServiceKey, "configIdx"), "$"); err == nil {

					pxData := entity.PixelStorage{
						URLServiceKey:  p.URLServiceKey,
						Pxdate:         external.GetCurrentTime(h.Config.TZ, time.RFC3339),
						Pixel:          p.AffSub,
						PostbackMethod: p.Method,
					}

					var px entity.PixelStorage

					isPX := false

					switch p.Method {
					case "ADNETCODE":
						px, isPX = h.DS.GetByAdnetCode(pxData)
					case "TOKEN":
						px, isPX = h.DS.GetToken(pxData)
					case "JSON-MSISDN", "XML-MSISDN", "HTML-MSISDN":
						px, isPX = h.DS.GetPxByMsisdn(pxData)
					case "PIXEL":
						px, isPX = h.DS.GetPx(pxData)
					case "SPC-MVLS":

						campIdRemover := strings.NewReplacer(dc.URLServiceKey+"-", "")
						msisdn := campIdRemover.Replace(p.AffSub)

						isPX = true

						px = entity.PixelStorage{
							CampaignDetailId:  dc.Id,
							Pxdate:            external.GetCurrentTime(h.Config.TZ, time.RFC3339),
							URLServiceKey:     dc.URLServiceKey,
							CampaignId:        dc.CampaignId,
							Country:           dc.Country,
							Partner:           dc.Partner,
							Operator:          dc.Operator,
							Aggregator:        dc.Aggregator,
							Service:           dc.Service,
							ShortCode:         dc.ShortCode,
							Adnet:             dc.Adnet,
							Keyword:           dc.Keyword,
							Subkeyword:        dc.SubKeyword,
							IsBillable:        dc.IsBillable,
							Plan:              dc.Plan,
							URL:               dc.APIURL,
							URLType:           dc.URLType,
							Pixel:             "NA",
							Msisdn:            msisdn,
							TrxId:             "NA",
							Token:             "NA",
							IsUsed:            true,
							Browser:           "NA",
							OS:                "NA",
							Ip:                strings.Join(c.IPs(), ", "),
							ISP:               "NA",
							ReferralURL:       "NA",
							PubId:             dc.PubId,
							UserAgent:         "NA",
							TrafficSource:     false,
							TrafficSourceData: "NA",
							UserRejected:      false,
							UniqueClick:       false,
							UserDuplicated:    false,
							Handset:           "NA",
							HandsetCode:       "NA",
							HandsetType:       "NA",
							URLLanding:        dc.URLLanding,
							URLWarpLanding:    dc.URLWarpLanding,
							URLService:        dc.URLService,
							URLTFCORSmartlink: dc.URLTFCSmartlink,
							StatusCapping:     dc.StatusCapping,
							StatusRatio:       dc.StatusRatio,
							PO:                dc.PO,
							Cost:              dc.Cost,
							CampaignObjective: dc.Objective,
							Channel:           dc.Channel,
							Currency:          dc.Currency,
							PostbackMethod:    dc.PostbackMethod,
							LandingTime:       external.GetCurrentTime(h.Config.TZ, time.RFC3339),
							LandedTime:        float64(0),
							HttpStatus:        200,
							IsOperator:        false,
							CreatedAt:         external.GetCurrentTime(h.Config.TZ, time.RFC3339),
							UpdatedAt:         external.GetCurrentTime(h.Config.TZ, time.RFC3339),
						}

						h.DS.NewPixel(px)
					}

					if !isPX && p.Method == "" {

						return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Pixel not found or duplicate used and parameter should have a method parameter, pixel : " + p.AffSub})

					} else {

						if !isPX {

							return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Pixel not found or duplicate used and parameter should have a method parameter, pixel : " + p.AffSub})

						} else {

							if px.IsUsed {

								return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusNotFound, Message: "NOK - Pixel already used", Data: entity.PixelStorageRsp{
									URLServiceKey: dc.URLServiceKey,
									Adnet:         dc.Adnet,
									IsBillable:    dc.IsBillable,
									Pixel:         px.Pixel,
									Browser:       px.Browser,
									OS:            px.OS,
									Handset:       px.UserAgent,
									PubId:         px.PubId,
									PixelUsedDate: px.PixelUsedDate.Format(time.RFC3339),
								}})

							} else {

								px.PixelUsedDate = external.GetCurrentTime(h.Config.TZ, time.RFC3339)

								bodyReq, _ := json.Marshal(px)

								corId := "RTO" + external.GetUniqId(h.Config.TZ)

								pubCtx2, pubCancel2 := context.WithTimeout(c.UserContext(), time.Duration(h.Config.RabbitMQCtxTimeout)*time.Second)
								defer pubCancel2()
								if err := h.RM.PublishWithRetry(pubCtx2, h.Config.RabbitMQRatioExchangeName, h.Config.RabbitMQRatioQueueName, bodyReq, corId); err != nil {
									h.Logs.Debug(fmt.Sprintf("[x] Failed published: %s, Data: %s ...", corId, string(bodyReq)))
								} else {
									h.Logs.Debug(fmt.Sprintf("[v] Published: %s, Data: %s ...", corId, string(bodyReq)))
								}

								return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusOK, Message: "OK", Data: entity.PixelStorageRsp{
									URLServiceKey: dc.URLServiceKey,
									Adnet:         dc.Adnet,
									IsBillable:    dc.IsBillable,
									Pixel:         px.Pixel,
									Browser:       px.Browser,
									OS:            px.OS,
									Handset:       px.UserAgent,
									PubId:         px.PubId,
									PixelUsedDate: external.GetFormatTime(h.Config.TZ, time.RFC3339),
								}})
							}
						}
					}

				} else {

					return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Campaign ID not found"})

				}
			}
		}
	}
}

func (h *IncomingHandler) PostbackV3(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	h.Logs.Debug(fmt.Sprintf("Receive request postback %#v ...\n", c.AllParams()))

	// Parse Postback Data
	p := entity.NewDataPostbackV2(c)
	//p.URLServiceKey = c.Params("urlservicekey")

	// Validate Parameters
	if v := p.ValidateParamsV2(h.Logs); v.Code == fiber.StatusBadRequest {

		return c.Status(v.Code).JSON(entity.GlobalResponse{Code: v.Code, Message: v.Message})

	} else {

		if c.Cookies(p.CookieKey) != "" {

			return c.Status(fiber.StatusForbidden).JSON(entity.GlobalResponse{Code: fiber.StatusForbidden, Message: "forbidden access"})

		} else {
			// Setup cookie if double requested within n-hour
			c.Cookie(&fiber.Cookie{
				Name:     p.CookieKey,
				Value:    "1",
				Expires:  time.Now().Add(3 * time.Second),
				HTTPOnly: true,
				SameSite: "lax",
			})

			if !strings.Contains(p.AffSub, "-") {

				return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Invalid pixel format, pixel : " + p.AffSub})

			} else {

				dataraw := strings.Split(p.AffSub, "-")
				p.URLServiceKey = external.Concat("-", dataraw[0], dataraw[1])

				if dc, err := h.DS.GetDataConfig(external.Concat("-", p.URLServiceKey, "configIdx"), "$"); err == nil {

					if !dc.IsActive {
						return c.Status(fiber.StatusForbidden).JSON(entity.GlobalResponse{
							Code:    fiber.StatusForbidden,
							Message: "Campaign is currently inactive",
						})
					}

					var (
						//px_byte []byte
						px entity.PixelStorage
					)
					isPX := false

					pxData := entity.PixelStorage{
						URLServiceKey:        p.URLServiceKey,
						Pxdate:               external.GetCurrentTime(h.Config.TZ, time.RFC3339),
						Pixel:                p.AffSub,
						PostbackMethod:       p.Method,
						StatusBillable:       p.Status,
						StatusCodeBillable:   p.StatusCode,
						ReasonStatusBillable: p.StatusDetail,
					}

					switch p.Method {
					case "ADNETCODE":

						/* var key string

						breaking := false
						iter := h.RCP.Scan(0, p.URLServiceKey+"*", 0).Iterator()

						if err := iter.Err(); err == nil {

							for iter.Next() {
								key = iter.Val()
								//fmt.Println("keys", key)
								break
							}

							px_byte = []byte(h.RCP.Get(key).Val())
							isPX = true

							if err = json.Unmarshal(px_byte, &px); err == nil {
								h.RCP.Del(px.Pixel)
								breaking = true
							}
						}

						if !breaking { */
						px, isPX = h.DS.GetByAdnetCode(pxData)

						/* if !isPX {
								return c.Status(fiber.StatusNotFound).JSON(
									entity.GlobalResponse{
										Code: fiber.StatusNotFound,
										Message: "Invalid pixel format or this pixel not found, pixel : " + p.AffSub,
									},
								)
							}
						} */

					case "TOKEN":
						px, isPX = h.DS.GetToken(pxData)
					case "JSON-MSISDN", "XML-MSISDN", "HTML-MSISDN":
						px, isPX = h.DS.GetPxByMsisdn(pxData)
					case "PIXEL":
						if g := h.RCP.Get(p.AffSub); g.Val() != "" {

							isPX = true

							if err = json.Unmarshal([]byte(g.Val()), &px); err != nil {

								return c.Status(fiber.StatusNotAcceptable).JSON(entity.GlobalResponse{Code: fiber.StatusNotAcceptable, Message: "Invalid pixel format or this pixel not found, pixel : " + p.AffSub})
							}

							h.RCP.Del(p.AffSub)

						} else {

							switch dc.Partner {
							case "ID-XLSMART-LINKIT":
								px, isPX = h.DS.SpecialGetPx(pxData, h.Config.StartGetIntervalDatePXS, h.Config.EndGetIntervalDatePXS)
							default:
								px, isPX = h.DS.GetPx(pxData)
							}

						}
					case "SPC-MVLS", "SPC-TFCS", "SPC":

						//campIdRemover := strings.NewReplacer(dc.URLServiceKey+"-", "")
						//msisdn := campIdRemover.Replace(p.AffSub)

						isPX = true

						px = entity.PixelStorage{
							CampaignDetailId:  dc.Id,
							Pxdate:            external.GetCurrentTime(h.Config.TZ, time.RFC3339),
							URLServiceKey:     dc.URLServiceKey,
							CampaignId:        dc.CampaignId,
							Country:           dc.Country,
							Partner:           dc.Partner,
							Operator:          dc.Operator,
							Aggregator:        dc.Aggregator,
							Service:           dc.Service,
							ShortCode:         dc.ShortCode,
							Adnet:             dc.Adnet,
							Keyword:           dc.Keyword,
							Subkeyword:        p.SubKeyword,
							IsBillable:        dc.IsBillable,
							Plan:              dc.Plan,
							URL:               dc.APIURL,
							URLType:           dc.URLType,
							Pixel:             "NA",
							Msisdn:            p.Msisdn,
							TrxId:             p.Trxid,
							Token:             "NA",
							IsUsed:            false,
							Browser:           "NA",
							OS:                "NA",
							Ip:                strings.Join(c.IPs(), ", "),
							ISP:               "NA",
							ReferralURL:       "NA",
							PubId:             dc.PubId,
							UserAgent:         "NA",
							TrafficSource:     false,
							TrafficSourceData: "NA",
							UserRejected:      false,
							UniqueClick:       false,
							UserDuplicated:    false,
							Handset:           "NA",
							HandsetCode:       "NA",
							HandsetType:       "NA",
							URLLanding:        dc.URLLanding,
							URLWarpLanding:    dc.URLWarpLanding,
							URLService:        dc.URLService,
							URLTFCORSmartlink: dc.URLTFCSmartlink,
							StatusCapping:     dc.StatusCapping,
							StatusRatio:       dc.StatusRatio,
							PO:                dc.PO,
							Cost:              dc.Cost,
							CampaignObjective: dc.Objective,
							Channel:           dc.Channel,
							Currency:          dc.Currency,
							PostbackMethod:    p.Method,
							LandingTime:       external.GetCurrentTime(h.Config.TZ, time.RFC3339),
							LandedTime:        float64(0),
							HttpStatus:        200,
							IsOperator:        false,
							CreatedAt:         external.GetCurrentTime(h.Config.TZ, time.RFC3339),
							UpdatedAt:         external.GetCurrentTime(h.Config.TZ, time.RFC3339),
							IsUnique:          false,
						}

						px.ID = h.DS.NewPixel(px)

						h.DS.UpdateSummaryFromLandingPixelStorage(
							entity.IncSummaryCampaign{
								SummaryDate:   px.Pxdate,
								URLServiceKey: px.URLServiceKey,
								Country:       px.Country,
								Operator:      px.Operator,
								Partner:       px.Partner,
								Service:       px.Service,
								Adnet:         px.Adnet,
								CampaignId:    px.CampaignId,
							})

						h.DS.UpdateSummaryFromLandingPixelStorageHour(
							entity.IncSummaryCampaignHour{
								SummaryDateHour: px.Pxdate,
								URLServiceKey:   px.URLServiceKey,
								Country:         px.Country,
								Operator:        px.Operator,
								Partner:         px.Partner,
								Service:         px.Service,
								Adnet:           px.Adnet,
								CampaignId:      px.CampaignId,
							})

					}

					if !isPX && p.Method == "" {

						return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Pixel not found or duplicate used and parameter should have a method parameter, pixel : " + p.AffSub})

					} else {

						if !isPX {

							return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Pixel not found or duplicate used and parameter should have a method parameter, pixel : " + p.AffSub})

						} else {

							if px.IsUsed {

								return c.Status(fiber.StatusConflict).JSON(entity.GlobalResponseWithData{Code: fiber.StatusConflict, Message: "NOK - Pixel already used", Data: entity.PixelStorageRsp{
									URLServiceKey: dc.URLServiceKey,
									Adnet:         dc.Adnet,
									IsBillable:    dc.IsBillable,
									Pixel:         px.Pixel,
									Browser:       px.Browser,
									OS:            px.OS,
									Handset:       px.UserAgent,
									PubId:         px.PubId,
									PixelUsedDate: px.PixelUsedDate.Format(time.RFC3339),
								}})

							} else {

								px.PixelUsedDate = external.GetCurrentTime(h.Config.TZ, time.RFC3339)

								bodyReq, _ := json.Marshal(px)

								corId := "RTO" + external.GetUniqId(h.Config.TZ)

								pubCtx3, pubCancel3 := context.WithTimeout(c.UserContext(), time.Duration(h.Config.RabbitMQCtxTimeout)*time.Second)
								defer pubCancel3()
								if err := h.RM.PublishWithRetry(pubCtx3, h.Config.RabbitMQRatioExchangeName, h.Config.RabbitMQRatioQueueName, bodyReq, corId); err != nil {
									h.Logs.Debug(fmt.Sprintf("[x] Failed published: %s, Data: %s ...", corId, string(bodyReq)))
								} else {
									h.Logs.Debug(fmt.Sprintf("[v] Published: %s, Data: %s ...", corId, string(bodyReq)))
								}

								return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusOK, Message: "OK", Data: entity.PixelStorageRsp{
									URLServiceKey: dc.URLServiceKey,
									Adnet:         dc.Adnet,
									IsBillable:    dc.IsBillable,
									Pixel:         px.Pixel,
									Browser:       px.Browser,
									OS:            px.OS,
									Handset:       px.UserAgent,
									PubId:         px.PubId,
									PixelUsedDate: external.GetFormatTime(h.Config.TZ, time.RFC3339),
								}})
							}
						}
					}

				} else {

					return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Campaign ID not found"})

				}
			}
		}
	}
}

func (h *IncomingHandler) PostbackBilled(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	h.Logs.Debug(fmt.Sprintf("Receive request postback billed %#v ...\n", c.AllParams()))

	// Parse Postback Data
	p := entity.NewDataPostbackV2(c)

	// Validate Parameters
	if v := p.ValidateParamsV2(h.Logs); v.Code == fiber.StatusBadRequest {
		return c.Status(v.Code).JSON(entity.GlobalResponse{
			Code:    v.Code,
			Message: v.Message,
		})
	}

	if c.Cookies(p.CookieKey) != "" {
		return c.Status(fiber.StatusForbidden).JSON(entity.GlobalResponse{
			Code:    fiber.StatusForbidden,
			Message: "forbidden access",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     p.CookieKey,
		Value:    "1",
		Expires:  time.Now().Add(3 * time.Second),
		HTTPOnly: true,
		SameSite: "lax",
	})

	if !strings.Contains(p.AffSub, "-") {
		return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{
			Code:    fiber.StatusNotFound,
			Message: "Invalid pixel format, pixel : " + p.AffSub,
		})
	}

	dataraw := strings.Split(p.AffSub, "-")
	p.URLServiceKey = external.Concat("-", dataraw[0], dataraw[1])

	dc, err := h.DS.GetDataConfig(external.Concat("-", p.URLServiceKey, "configIdx"), "$")
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{
			Code:    fiber.StatusNotFound,
			Message: "Campaign ID not found",
		})
	}

	// full untuk DB
	pixelDB := p.AffSub

	// clean untuk Google Sheet
	apiurl := strings.NewReplacer(p.URLServiceKey+"-", "")
	pixelGS := apiurl.Replace(p.AffSub)

	now := external.GetCurrentTime(h.Config.TZ, time.RFC3339)

	pixelStorage := entity.PixelStorage{
		CampaignId:    dc.CampaignId,
		GoogleSheet:   dc.GoogleSheetBillable,
		Pixel:         pixelDB,
		PixelUsedDate: now,
		Msisdn:        p.Msisdn,
		URLServiceKey: p.URLServiceKey,
		MStatusCharge: strings.TrimSpace(strings.ToLower(p.Status)) == "success",
	}

	if err := h.DS.UpdatePixelBilled(pixelStorage); err != nil {
		h.Logs.Error(fmt.Sprintf("failed update pixel billed: %#v", err))
	}

	h.DS.UpdateGoogleSheetPixel(
		h.GS,
		entity.PixelStorage{
			CampaignId:    dc.CampaignId,
			GoogleSheet:   dc.GoogleSheetBillable,
			Pixel:         pixelGS,
			PixelUsedDate: now,
			Msisdn:        p.Msisdn,
		},
		"Billed",
		p.Desc,
	)

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponse{
		Code:    fiber.StatusOK,
		Message: "OK",
	})
}

func (h *IncomingHandler) InquiryCampID(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	h.Logs.Debug(fmt.Sprintf("Inquiry Camp ID %#v ...\n", c.AllParams()))

	request := new(entity.InquiryCampID)

	if err := c.QueryParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{Code: fiber.StatusBadRequest, Message: "check mandatory param : urlservicekey"})
	} else {

		CookieKey := external.Concat("-", external.GetIpAddress(c), request.URLServiceKey)

		if c.Cookies(CookieKey) != "" {

			return c.Status(fiber.StatusForbidden).JSON(entity.GlobalResponse{Code: fiber.StatusForbidden, Message: "forbidden access"})

		} else {
			// Setup cookie if double requested within n-hour
			c.Cookie(&fiber.Cookie{
				Name:     CookieKey,
				Value:    "1",
				Expires:  time.Now().Add(30 * time.Minute),
				HTTPOnly: true,
				SameSite: "lax",
			})

			if dc, err := h.DS.GetDataConfig(external.Concat("-", request.URLServiceKey, "configIdx"), "$"); err == nil {

				return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusOK, Message: "OK", Data: entity.InquiryCampID{
					URLServiceKey: dc.URLServiceKey,
					Country:       dc.Country,
					Operator:      dc.Operator,
					Partner:       dc.Partner,
					Aggregator:    dc.Aggregator,
					Adnet:         dc.Adnet,
					Service:       dc.Service,
					URLLanding:    dc.URLWarpLanding,
					URLService:    dc.URLService,
				}})

			} else {

				return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "url service key / campaign id not found"})
			}
		}
	}
}

func (h *IncomingHandler) InquiryAPICampID(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	h.Logs.Debug(fmt.Sprintf("Inquiry API Camp ID By Params %#v ...\n", c.AllParams()))

	request := new(entity.InquiryAPICampID)

	if err := c.QueryParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{Code: fiber.StatusBadRequest, Message: "check mandatory param : country, operator, adnet"})
	}

	request.Country = strings.ToUpper(request.Country)
	request.Operator = strings.ToUpper(request.Operator)
	request.Service = strings.ToUpper(request.Service)
	request.Adnet = strings.ToUpper(request.Adnet)

	if request.Country == "" || request.Operator == "" || request.Service == "" || request.Adnet == "" {
		return c.Status(fiber.StatusBadRequest).JSON(entity.GlobalResponse{Code: fiber.StatusBadRequest, Message: "mandatory params missing: country, operator, service, adnet"})
	}

	results, err := h.DS.GetAPICampaignDetails(request.Country, request.Operator, request.Service, request.Adnet)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entity.GlobalResponse{Code: fiber.StatusInternalServerError, Message: "failed to retrieve configs"})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "no campaign found for given params"})
	}

	return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusOK, Message: "OK", Data: results})
}

func (h *IncomingHandler) PostbackDirectReply(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	h.Logs.Debug(fmt.Sprintf("Receive request postback %#v ...\n", c.AllParams()))

	// Parse Postback Data
	p := entity.NewDataPostbackV2(c)

	// Validate Parameters
	if v := p.ValidateParamsV2(h.Logs); v.Code == fiber.StatusBadRequest {
		return c.Status(v.Code).JSON(entity.GlobalResponse{Code: v.Code, Message: v.Message})
	} else {
		if c.Cookies(p.CookieKey) != "" {
			return c.Status(fiber.StatusForbidden).JSON(entity.GlobalResponse{Code: fiber.StatusForbidden, Message: "forbidden access"})
		} else {
			// Setup cookie if double requested within n-hour
			c.Cookie(&fiber.Cookie{
				Name:     p.CookieKey,
				Value:    "1",
				Expires:  time.Now().Add(3 * time.Second),
				HTTPOnly: true,
				SameSite: "lax",
			})

			if !strings.Contains(p.AffSub, "-") {
				return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Invalid pixel format, pixel : " + p.AffSub})
			} else {
				dataraw := strings.Split(p.AffSub, "-")
				p.URLServiceKey = external.Concat("-", dataraw[0], dataraw[1])

				if dc, err := h.DS.GetDataConfig(external.Concat("-", p.URLServiceKey, "configIdx"), "$"); err == nil {
					if !dc.IsActive {
						return c.Status(fiber.StatusForbidden).JSON(entity.GlobalResponse{
							Code:    fiber.StatusForbidden,
							Message: "Campaign is currently inactive",
						})
					}

					var px entity.PixelStorage
					isPX := false

					pxData := entity.PixelStorage{
						URLServiceKey:        p.URLServiceKey,
						Pxdate:               external.GetCurrentTime(h.Config.TZ, time.RFC3339),
						Pixel:                p.AffSub,
						PostbackMethod:       p.Method,
						StatusBillable:       p.Status,
						StatusCodeBillable:   p.StatusCode,
						ReasonStatusBillable: p.StatusDetail,
					}

					switch p.Method {
					case "ADNETCODE":
						px, isPX = h.DS.GetByAdnetCode(pxData)
					case "TOKEN":
						px, isPX = h.DS.GetToken(pxData)
					case "JSON-MSISDN", "XML-MSISDN", "HTML-MSISDN":
						px, isPX = h.DS.GetPxByMsisdn(pxData)
					case "PIXEL":
						if g := h.RCP.Get(p.AffSub); g.Val() != "" {
							isPX = true
							if err = json.Unmarshal([]byte(g.Val()), &px); err != nil {
								return c.Status(fiber.StatusNotAcceptable).JSON(entity.GlobalResponse{Code: fiber.StatusNotAcceptable, Message: "Invalid pixel format or this pixel not found, pixel : " + p.AffSub})
							}
							h.RCP.Del(p.AffSub)
						} else {
							switch dc.Partner {
							case "ID-XLSMART-LINKIT":
								px, isPX = h.DS.SpecialGetPx(pxData, h.Config.StartGetIntervalDatePXS, h.Config.EndGetIntervalDatePXS)
							default:
								px, isPX = h.DS.GetPx(pxData)
							}
						}
					case "SPC-MVLS", "SPC-TFCS", "SPC":
						isPX = true
						px = entity.PixelStorage{
							CampaignDetailId:  dc.Id,
							Pxdate:            external.GetCurrentTime(h.Config.TZ, time.RFC3339),
							URLServiceKey:     dc.URLServiceKey,
							CampaignId:        dc.CampaignId,
							Country:           dc.Country,
							Partner:           dc.Partner,
							Operator:          dc.Operator,
							Aggregator:        dc.Aggregator,
							Service:           dc.Service,
							ShortCode:         dc.ShortCode,
							Adnet:             dc.Adnet,
							Keyword:           dc.Keyword,
							Subkeyword:        p.SubKeyword,
							IsBillable:        dc.IsBillable,
							Plan:              dc.Plan,
							URL:               dc.APIURL,
							URLType:           dc.URLType,
							Pixel:             "NA",
							Msisdn:            p.Msisdn,
							TrxId:             p.Trxid,
							Token:             "NA",
							IsUsed:            false,
							Browser:           "NA",
							OS:                "NA",
							Ip:                strings.Join(c.IPs(), ", "),
							ISP:               "NA",
							ReferralURL:       "NA",
							PubId:             dc.PubId,
							UserAgent:         "NA",
							TrafficSource:     false,
							TrafficSourceData: "NA",
							UserRejected:      false,
							UniqueClick:       false,
							UserDuplicated:    false,
							Handset:           "NA",
							HandsetCode:       "NA",
							HandsetType:       "NA",
							URLLanding:        dc.URLLanding,
							URLWarpLanding:    dc.URLWarpLanding,
							URLService:        dc.URLService,
							URLTFCORSmartlink: dc.URLTFCSmartlink,
							StatusCapping:     dc.StatusCapping,
							StatusRatio:       dc.StatusRatio,
							PO:                dc.PO,
							Cost:              dc.Cost,
							CampaignObjective: dc.Objective,
							Channel:           dc.Channel,
							Currency:          dc.Currency,
							PostbackMethod:    p.Method,
							LandingTime:       external.GetCurrentTime(h.Config.TZ, time.RFC3339),
							LandedTime:        float64(0),
							HttpStatus:        200,
							IsOperator:        false,
							CreatedAt:         external.GetCurrentTime(h.Config.TZ, time.RFC3339),
							UpdatedAt:         external.GetCurrentTime(h.Config.TZ, time.RFC3339),
							IsUnique:          false,
						}

						px.ID = h.DS.NewPixel(px)

						h.DS.UpdateSummaryFromLandingPixelStorage(
							entity.IncSummaryCampaign{
								SummaryDate:   px.Pxdate,
								URLServiceKey: px.URLServiceKey,
								Country:       px.Country,
								Operator:      px.Operator,
								Partner:       px.Partner,
								Service:       px.Service,
								Adnet:         px.Adnet,
								CampaignId:    px.CampaignId,
							})

						h.DS.UpdateSummaryFromLandingPixelStorageHour(
							entity.IncSummaryCampaignHour{
								SummaryDateHour: px.Pxdate,
								URLServiceKey:   px.URLServiceKey,
								Country:         px.Country,
								Operator:        px.Operator,
								Partner:         px.Partner,
								Service:         px.Service,
								Adnet:           px.Adnet,
								CampaignId:      px.CampaignId,
							})
					}

					if !isPX && p.Method == "" {
						return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Pixel not found or duplicate used and parameter should have a method parameter, pixel : " + p.AffSub})
					} else {
						if !isPX {
							return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Pixel not found or duplicate used and parameter should have a method parameter, pixel : " + p.AffSub})
						} else {
							if px.IsUsed {
								return c.Status(fiber.StatusConflict).JSON(entity.GlobalResponseWithData{Code: fiber.StatusConflict, Message: "NOK - Pixel already used", Data: entity.PixelStorageRsp{
									URLServiceKey: dc.URLServiceKey,
									Adnet:         dc.Adnet,
									IsBillable:    dc.IsBillable,
									Pixel:         px.Pixel,
									Browser:       px.Browser,
									OS:            px.OS,
									Handset:       px.UserAgent,
									PubId:         px.PubId,
									PixelUsedDate: px.PixelUsedDate.Format(time.RFC3339),
								}})
							} else {
								px.PixelUsedDate = external.GetCurrentTime(h.Config.TZ, time.RFC3339)
								bodyReq, _ := json.Marshal(px)
								corId := "RTD" + external.GetUniqId(h.Config.TZ)

								reply := "NOTSHAVED"
								ctx, cancel := context.WithTimeout(c.UserContext(), time.Duration(h.Config.RabbitMQCtxTimeout)*time.Second)
								defer cancel()

								if respBody, pubErr := h.RM.DirectReplyToWithRetry(ctx, h.Config.RabbitMQRatioExchangeName, h.Config.RabbitMQRatioQueueName, bodyReq, corId); pubErr != nil {
									h.Logs.Debug(fmt.Sprintf("[x] Failed published, Error: %v, Data: %s ...", pubErr, string(bodyReq)))
								} else {
									h.Logs.Debug(fmt.Sprintf("[v] Published, Data: %s, Response: %s ...", string(bodyReq), string(respBody)))
									reply = string(respBody)
								}

								return c.Status(fiber.StatusOK).JSON(entity.GlobalResponseWithData{Code: fiber.StatusOK, Message: reply, Data: entity.PixelStorageRsp{
									URLServiceKey: dc.URLServiceKey,
									Adnet:         dc.Adnet,
									IsBillable:    dc.IsBillable,
									Pixel:         px.Pixel,
									Browser:       px.Browser,
									OS:            px.OS,
									Handset:       px.UserAgent,
									PubId:         px.PubId,
									PixelUsedDate: external.GetFormatTime(h.Config.TZ, time.RFC3339),
								}})
							}
						}
					}
				} else {
					return c.Status(fiber.StatusNotFound).JSON(entity.GlobalResponse{Code: fiber.StatusNotFound, Message: "Campaign ID not found"})
				}
			}
		}
	}
}

