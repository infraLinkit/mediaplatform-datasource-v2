package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/config"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/infrastructure/external"
)

func (h *IncomingHandler) DisplayCampaignManagement(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Accepts("application/x-www-form-urlencoded")
	c.AcceptsCharsets("utf-8", "iso-8859-1")

	m := c.Queries()
	page, _ := strconv.Atoi(m["page"])
	draw, _ := strconv.Atoi(m["draw"])

	fe := entity.DisplayCampaignManagement{
		Adnet:         m["adnet"],
		Country:       m["country"],
		Service:       m["service"],
		Operator:      m["operator"],
		Partner:       m["partner"],
		Status:        m["status"],
		CampaignName:  m["campaign_name"],
		CampaignType:  m["campaign_type"],
		CampaignId:    m["campaign_id"],
		Page:          page,
		Draw:          draw,
		Action:        m["action"],
		URLServiceKey: m["url_service_key"],
		OrderColumn:   m["order_column"],
		OrderDir:      m["order_dir"],
		CreatedDateBefore: m["created_date_before"],
		CreatedDateAfter:  m["created_date_after"],
	}

	v := c.Params("v")

	if v == "detail" {
		r := h.DisplayCampaignManagementDetail(c, fe)
		return c.Status(r.HttpStatus).JSON(r.Rsp)
	}

	r := h.DisplayCampaignManagementExtra(c, fe, draw)
	return c.Status(r.HttpStatus).JSON(r.Rsp)
}

func (h *IncomingHandler) DisplayCampaignManagementExtra(c *fiber.Ctx, fe entity.DisplayCampaignManagement, draw int) entity.ReturnResponse {
	key := "temp_key_api_campaign_management_" + strings.ReplaceAll(external.GetIpAddress(c), ".", "_")

	var (
		err                       error
		isempty                   bool
		campaignmanagement        []entity.CampaignManagementData
		displaycampaignmanagement []entity.CampaignManagementData
	)

	if fe.Action != "" {
		campaignmanagement, _, err = h.DS.GetCampaignManagement(fe)
	} else {
		if campaignmanagement, isempty = h.DS.RGetCampaignManagement(key, "$"); isempty {
			campaignmanagement, _, err = h.DS.GetCampaignManagement(fe)
			s, _ := json.Marshal(campaignmanagement)
			h.DS.SetData(key, "$", string(s))
			h.DS.SetExpireData(key, 60)
		}
	}

	if err == nil {
		totalRecords := len(campaignmanagement)
		pagesize := PAGESIZE
		page := 1
		if fe.Page > 0 {
			page = fe.Page
		}
		start := (page - 1) * pagesize
		end := start + pagesize

		if start < totalRecords {
			if end > totalRecords {
				end = totalRecords
			}
			displaycampaignmanagement = campaignmanagement[start:end]
		}

		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:            fe.Draw,
				Code:            fiber.StatusOK,
				Message:         config.OK_DESC,
				Data:            displaycampaignmanagement,
				RecordsTotal:    totalRecords,
				RecordsFiltered: totalRecords,
			},
		}
	}

	return entity.ReturnResponse{
		HttpStatus: fiber.StatusNotFound,
		Rsp: entity.GlobalResponse{
			Code:    fiber.StatusNotFound,
			Message: "empty",
		},
	}
}

func (h *IncomingHandler) DisplayCampaignManagementDetail(c *fiber.Ctx, fe entity.DisplayCampaignManagement) entity.ReturnResponse {
	var (
		err                             error
		campaignmanagementdetail        []entity.CampaignManagementDataDetail
		displaycampaignmanagementdetail []entity.CampaignManagementDataDetail
	)

	campaignmanagementdetail, err = h.DS.GetCampaignManagementDetail(fe)

	if err == nil {
		displaycampaignmanagementdetail = campaignmanagementdetail
		return entity.ReturnResponse{
			HttpStatus: fiber.StatusOK,
			Rsp: entity.GlobalResponseWithDataTable{
				Draw:    fe.Draw,
				Code:    fiber.StatusOK,
				Message: config.OK_DESC,
				Data:    displaycampaignmanagementdetail,
			},
		}
	}

	return entity.ReturnResponse{
		HttpStatus: fiber.StatusNotFound,
		Rsp: entity.GlobalResponseWithData{
			Code:    fiber.StatusNotFound,
			Message: "empty",
		},
	}
}

func (h *IncomingHandler) SendCampaignHandler(c *fiber.Ctx) error {
	var campaignData map[string]interface{}

	if err := c.BodyParser(&campaignData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var bodyReq bytes.Buffer
	encoder := json.NewEncoder(&bodyReq)
	encoder.SetEscapeHTML(false) // Prevents encoding & as \u0026

	if err := encoder.Encode(campaignData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to serialize data"})
	}

	pubCtx, pubCancel := context.WithTimeout(c.UserContext(), time.Duration(h.Config.RabbitMQCtxTimeout)*time.Second)
	defer pubCancel()
	if err := h.RM.PublishWithRetry(pubCtx, h.Config.RabbitMQCampaignManagementExchangeName, h.Config.RabbitMQCampaignManagementQueueName, bodyReq.Bytes(), ""); err != nil {
		h.Logs.Debug(fmt.Sprintf("[x] Failed published: Data: %s ...", bodyReq.String()))
	} else {
		h.Logs.Debug(fmt.Sprintf("[v] Published: Data: %s ...", bodyReq.String()))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Campaign data sent to RabbitMQ"})
}

func (h *IncomingHandler) UpdateStatusCampaign(c *fiber.Ctx) error {

	o := new(entity.CampaignDetail)

	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	} else {

		h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

		// Update to redis with key
		cfgRediskey := external.Concat("-", o.URLServiceKey, "configIdx")
		cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")
		cfgCmp.IsActive = o.IsActive

		if o.IsActive {
			cfgCmp.StatusCapping = false
		}

		cfgDataConfig, _ := json.Marshal(cfgCmp)

		h.DS.SetData(cfgRediskey, "$", string(cfgDataConfig))

		// Update to database
		h.DS.UpdateStatusCounterMOCampaignDetail(entity.CampaignDetail{
			URLServiceKey: o.URLServiceKey,
			IsActive:      o.IsActive,
		})

		h.DS.UpdateSummaryCampaign(entity.SummaryCampaign{
			SummaryDate:   external.GetCurrentTime(h.Config.TZ, time.RFC3339),
			Status:        o.IsActive,
			URLServiceKey: o.URLServiceKey,
			Country:       cfgCmp.Country,
			Operator:      cfgCmp.Operator,
			Partner:       cfgCmp.Partner,
			Adnet:         cfgCmp.Adnet,
			Service:       cfgCmp.Service,
			CampaignId:    cfgCmp.CampaignId,
		})

		return c.Status(fiber.StatusOK).Send([]byte("OK"))
	}
}

func (h *IncomingHandler) EditCampaign(c *fiber.Ctx) error {
	o := new(entity.CampaignDetail)
	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	s := new(entity.SummaryCampaign)
	_ = c.BodyParser(&s)

	h.Logs.Debug(fmt.Sprintf("campaign detail : %#v ...", o))
	h.Logs.Debug(fmt.Sprintf("summary campaign : %#v ...", s))

	cfgRediskey := external.Concat("-", o.URLServiceKey, "configIdx")
	cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")
	mocappingChanged := o.MOCapping != cfgCmp.MOCapping

	cfgCmp.PO = o.PO
	cfgCmp.RatioSend = o.RatioSend
	cfgCmp.RatioReceive = o.RatioReceive
	cfgCmp.MOCapping = o.MOCapping
	cfgCmp.LastUpdate = external.GetFormatTime(h.Config.TZ, time.RFC3339)
	if mocappingChanged {
		cfgCmp.StatusCapping = false
	}
	cfgDataConfig, _ := json.Marshal(cfgCmp)
	h.DS.SetData(cfgRediskey, "$", string(cfgDataConfig))

	now := time.Now().In(h.Config.TZ)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, h.Config.TZ)

	summaryLocal := s.SummaryDate.In(h.Config.TZ)
	summaryLocal = time.Date(
		summaryLocal.Year(),
		summaryLocal.Month(),
		summaryLocal.Day(),
		0, 0, 0, 0,
		h.Config.TZ,
	)

	if s.SummaryDate.IsZero() || summaryLocal.Equal(today) {
		h.DS.EditSettingCampaignDetail(entity.CampaignDetail{
			PO:            o.PO,
			MOCapping:     o.MOCapping,
			RatioSend:     o.RatioSend,
			RatioReceive:  o.RatioReceive,
			LastUpdate:    external.GetCurrentTime(h.Config.TZ, time.RFC3339),
			URLServiceKey: o.URLServiceKey,
			Country:       cfgCmp.Country,
			Operator:      cfgCmp.Operator,
			Partner:       cfgCmp.Partner,
			Adnet:         cfgCmp.Adnet,
			Service:       cfgCmp.Service,
			CampaignId:    o.CampaignId,
			StatusCapping: bool(mocappingChanged),
		})
	}

	pos, _ := strconv.ParseFloat(strings.TrimSpace(o.PO), 64)
	h.DS.EditSettingSummaryCampaign(entity.SummaryCampaign{
		SummaryDate:   s.SummaryDate,
		PO:            pos,
		MOLimit:       o.MOCapping,
		RatioSend:     o.RatioSend,
		RatioReceive:  o.RatioReceive,
		URLServiceKey: o.URLServiceKey,
		Country:       cfgCmp.Country,
		Operator:      cfgCmp.Operator,
		Partner:       cfgCmp.Partner,
		Adnet:         cfgCmp.Adnet,
		Service:       cfgCmp.Service,
		CampaignId:    o.CampaignId,
	})

	return c.Status(fiber.StatusOK).Send([]byte("OK"))
}

func (h *IncomingHandler) EditCampaignMOCapping(c *fiber.Ctx) error {
	o := new(entity.CampaignDetail)
	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	s := new(entity.SummaryCampaign)
	_ = c.BodyParser(&s)

	cfgRediskey := external.Concat("-", o.URLServiceKey, "configIdx")
	cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")


	cfgCmp.MOCapping = o.MOCapping
	cfgCmp.StatusCapping = false
	cfgCmp.LastUpdate = external.GetFormatTime(h.Config.TZ, time.RFC3339)

	cfgData, _ := json.Marshal(cfgCmp)

	now := time.Now().In(h.Config.TZ)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, h.Config.TZ)

	summaryLocal := s.SummaryDate.In(h.Config.TZ)
	summaryLocal = time.Date(
		summaryLocal.Year(),
		summaryLocal.Month(),
		summaryLocal.Day(),
		0, 0, 0, 0,
		h.Config.TZ,
	)

	if s.SummaryDate.IsZero() || summaryLocal.Equal(today) {
		h.DS.SetData(cfgRediskey, "$", string(cfgData))
		h.DS.UpdateCampaignMOCapping(entity.CampaignDetail{
			MOCapping:     o.MOCapping,
			StatusCapping: false,
			LastUpdate:    external.GetCurrentTime(h.Config.TZ, time.RFC3339),
			URLServiceKey: o.URLServiceKey,
		})
	}

	h.DS.UpdateSummaryMOCapping(entity.SummaryCampaign{
		SummaryDate:   s.SummaryDate,
		MOLimit:       o.MOCapping,
		URLServiceKey: o.URLServiceKey,
	})

	return c.SendStatus(fiber.StatusOK)
}

func (h *IncomingHandler) EditCampaignRatio(c *fiber.Ctx) error {
	o := new(entity.CampaignDetail)
	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	s := new(entity.SummaryCampaign)
	_ = c.BodyParser(&s)

	cfgRediskey := external.Concat("-", o.URLServiceKey, "configIdx")
	cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")


	cfgCmp.RatioSend = o.RatioSend
	cfgCmp.RatioReceive = o.RatioReceive
	cfgCmp.LastUpdate = external.GetFormatTime(h.Config.TZ, time.RFC3339)

	cfgData, _ := json.Marshal(cfgCmp)

	now := time.Now().In(h.Config.TZ)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, h.Config.TZ)

	summaryLocal := s.SummaryDate.In(h.Config.TZ)
	summaryLocal = time.Date(
		summaryLocal.Year(),
		summaryLocal.Month(),
		summaryLocal.Day(),
		0, 0, 0, 0,
		h.Config.TZ,
	)

	if s.SummaryDate.IsZero() || summaryLocal.Equal(today) {
		h.DS.SetData(cfgRediskey, "$", string(cfgData))
		h.DS.UpdateCampaignRatio(entity.CampaignDetail{
			RatioSend:     o.RatioSend,
			RatioReceive:  o.RatioReceive,
			URLServiceKey: o.URLServiceKey,
		})
	}

	h.DS.UpdateSummaryRatio(entity.SummaryCampaign{
		SummaryDate:   s.SummaryDate,
		RatioSend:     o.RatioSend,
		RatioReceive:  o.RatioReceive,
		URLServiceKey: o.URLServiceKey,
	})

	return c.SendStatus(fiber.StatusOK)
}

func (h *IncomingHandler) EditCampaignPO(c *fiber.Ctx) error {
	o := new(entity.CampaignDetail)
	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	s := new(entity.SummaryCampaign)
	_ = c.BodyParser(&s)

	cfgRediskey := external.Concat("-", o.URLServiceKey, "configIdx")
	cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")


	cfgCmp.PO = o.PO
	cfgCmp.LastUpdate = external.GetFormatTime(h.Config.TZ, time.RFC3339)

	cfgData, _ := json.Marshal(cfgCmp)

	now := time.Now().In(h.Config.TZ)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, h.Config.TZ)

	summaryLocal := s.SummaryDate.In(h.Config.TZ)
	summaryLocal = time.Date(
		summaryLocal.Year(),
		summaryLocal.Month(),
		summaryLocal.Day(),
		0, 0, 0, 0,
		h.Config.TZ,
	)

	if s.SummaryDate.IsZero() || summaryLocal.Equal(today) {
		h.DS.SetData(cfgRediskey, "$", string(cfgData))
		h.DS.UpdateCampaignPO(entity.CampaignDetail{
			PO:            o.PO,
			URLServiceKey: o.URLServiceKey,
		})
	}

	pos, _ := strconv.ParseFloat(strings.TrimSpace(o.PO), 64)

	/* h.DS.UpdateSummaryPO(entity.SummaryCampaign{
		SummaryDate:   s.SummaryDate,
		PO:            pos,
		URLServiceKey: o.URLServiceKey,
		Country:       cfgCmp.Country,
		Operator:      cfgCmp.Operator,
		Partner:       cfgCmp.Partner,
		Adnet:         cfgCmp.Adnet,
		Service:       cfgCmp.Service,
		CampaignId:    o.CampaignId,
	}) */

	updatedPO := pos

	if sum, isOK := h.DS.GetSummaryCampaign(entity.SummaryCampaign{
		SummaryDate:   s.SummaryDate,
		URLServiceKey: o.URLServiceKey,
		Country:       cfgCmp.Country,
		Operator:      cfgCmp.Operator,
		Partner:       cfgCmp.Partner,
		Adnet:         cfgCmp.Adnet,
		Service:       cfgCmp.Service,
		CampaignId:    o.CampaignId,
	}); isOK {

		sum.PO = updatedPO

		calculated := h.DS.FormulaCPA(sum)

		calculated.PO = updatedPO

		// Re-calculate summary CPA
		h.DS.ReCalculateSummaryCampaign(calculated)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *IncomingHandler) EditCampaignSettingMOCapping(c *fiber.Ctx) error {
	req := new(entity.CampaignDetail)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cfgRedisKey := external.Concat("-", req.URLServiceKey, "configIdx")
	cfgCmp, err := h.DS.GetDataConfig(cfgRedisKey, "$")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cfgCmp.MOCapping = req.MOCapping
	cfgCmp.StatusCapping = false
	cfgCmp.LastUpdate = external.GetFormatTime(h.Config.TZ, time.RFC3339)

	cfgData, _ := json.Marshal(cfgCmp)
	h.DS.SetData(cfgRedisKey, "$", string(cfgData))

	if err := h.DS.UpdateCampaignMOCapping(entity.CampaignDetail{
		MOCapping:      req.MOCapping,
		StatusCapping:  false,
		LastUpdate:     external.GetCurrentTime(h.Config.TZ, time.RFC3339),
		URLServiceKey:  req.URLServiceKey,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *IncomingHandler) EditCampaignSettingRatio(c *fiber.Ctx) error {
	req := new(entity.CampaignDetail)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cfgRedisKey := external.Concat("-", req.URLServiceKey, "configIdx")
	cfgCmp, err := h.DS.GetDataConfig(cfgRedisKey, "$")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cfgCmp.RatioSend = req.RatioSend
	cfgCmp.RatioReceive = req.RatioReceive
	cfgCmp.LastUpdate = external.GetFormatTime(h.Config.TZ, time.RFC3339)

	cfgData, _ := json.Marshal(cfgCmp)
	h.DS.SetData(cfgRedisKey, "$", string(cfgData))

	if err := h.DS.UpdateCampaignRatio(entity.CampaignDetail{
		RatioSend:     req.RatioSend,
		RatioReceive:  req.RatioReceive,
		URLServiceKey: req.URLServiceKey,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *IncomingHandler) EditCampaignSettingPO(c *fiber.Ctx) error {
	req := new(entity.CampaignDetail)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cfgRedisKey := external.Concat("-", req.URLServiceKey, "configIdx")
	cfgCmp, err := h.DS.GetDataConfig(cfgRedisKey, "$")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cfgCmp.PO = req.PO
	cfgCmp.LastUpdate = external.GetFormatTime(h.Config.TZ, time.RFC3339)

	cfgData, _ := json.Marshal(cfgCmp)
	h.DS.SetData(cfgRedisKey, "$", string(cfgData))

	if err := h.DS.UpdateCampaignPO(entity.CampaignDetail{
		PO:            req.PO,
		URLServiceKey: req.URLServiceKey,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *IncomingHandler) DelCampaign(c *fiber.Ctx) error {

	o := new(entity.CampaignDetail)

	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	} else {

		h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

		// DELETE REDIS KEY
		cfgRediskey := external.Concat("-", o.URLServiceKey, "configIdx")
		cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")
		h.DS.DelData(cfgRediskey, "$")

		ctrRedisKey := external.Concat("-", o.URLServiceKey, "counterIdx")

		h.DS.DelData(ctrRedisKey, "$")

		sumRedisKey := external.Concat("-", o.URLServiceKey, "summary")

		h.DS.DelData(sumRedisKey, "$")

		h.DS.DelCampaignDetail(entity.CampaignDetail{
			URLServiceKey: o.URLServiceKey,
			Country:       cfgCmp.Country,
			Operator:      cfgCmp.Operator,
			Partner:       cfgCmp.Partner,
			Adnet:         cfgCmp.Adnet,
			Service:       cfgCmp.Service,
			CampaignId:    o.CampaignId,
		})

		h.DS.DelSummaryCampaign(entity.SummaryCampaign{
			SummaryDate:   external.GetCurrentTime(h.Config.TZ, time.RFC3339),
			URLServiceKey: o.URLServiceKey,
			Country:       cfgCmp.Country,
			Operator:      cfgCmp.Operator,
			Partner:       cfgCmp.Partner,
			Adnet:         cfgCmp.Adnet,
			Service:       cfgCmp.Service,
			CampaignId:    o.CampaignId,
		})

		count, err := h.DS.CountCampaignDetailByCampaignID(entity.CampaignDetail{CampaignId: o.CampaignId})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count campaign details"})
		}

		if count == 0 {
			err = h.DS.DelCampaign(entity.Campaign{CampaignId: o.CampaignId})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete campaign"})
			}
		}

		return c.Status(fiber.StatusOK).Send([]byte("OK"))
	}
}

func (h *IncomingHandler) GetCampaignCounts(c *fiber.Ctx) error {
	var input entity.DisplayCampaignManagement

	if err := c.QueryParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Ambil hanya counts dari GetCampaignManagement
	_, counts, err := h.DS.GetCampaignManagement(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total_campaign":          counts.TotalCampaigns,
		"total_active_campaign":   counts.TotalActiveCampaigns,
		"total_inactive_campaign": counts.TotalNonActiveCampaigns,
	})
}

func (h *IncomingHandler) UpdateKeyMainstream(c *fiber.Ctx) error {

	o := new(entity.CampaignDetail)

	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	} else {

		h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

		// Update to redis with key
		cfgRediskey := external.Concat("-", o.URLServiceKey, "configIdx")
		cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")
		cfgCmp.StatusSubmitKeyMainstream = o.StatusSubmitKeyMainstream
		cfgCmp.KeyMainstream = o.KeyMainstream

		cfgDataConfig, _ := json.Marshal(cfgCmp)

		h.DS.SetData(cfgRediskey, "$", string(cfgDataConfig))

		err = h.DS.UpdateKeyMainstreamCampaignDetail(entity.CampaignDetail{
			StatusSubmitKeyMainstream: o.StatusSubmitKeyMainstream,
			KeyMainstream:             o.KeyMainstream,
			URLServiceKey:             o.URLServiceKey,
			CampaignId:                o.CampaignId,
		})

		if err != nil {
			h.Logs.Error(fmt.Sprintf("failed updating db: %v", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed save to db"})
		}

		return c.Status(fiber.StatusOK).Send([]byte("OK"))
	}
}

func (h *IncomingHandler) UpdateGoogleSheet(c *fiber.Ctx) error {

	o := new(entity.CampaignDetail)

	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	} else {

		h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

		// Update to redis with key
		cfgRediskey := external.Concat("-", o.URLServiceKey, "configIdx")
		cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")
		cfgCmp.GoogleSheet = o.GoogleSheet

		cfgDataConfig, _ := json.Marshal(cfgCmp)

		h.DS.SetData(cfgRediskey, "$", string(cfgDataConfig))

		err = h.DS.UpdateGoogleSheetCampaignDetail(entity.CampaignDetail{
			GoogleSheet:   o.GoogleSheet,
			URLServiceKey: o.URLServiceKey,
			CampaignId:    o.CampaignId,
		})

		if err != nil {
			h.Logs.Error(fmt.Sprintf("failed updating db: %v", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed save to db"})
		}

		return c.Status(fiber.StatusOK).Send([]byte("OK"))
	}
}

func (h *IncomingHandler) UpdateGoogleSheetBillable(c *fiber.Ctx) error {

	o := new(entity.CampaignDetail)

	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	} else {

		h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

		// Update to redis with key
		cfgRediskey := external.Concat("-", o.URLServiceKey, "configIdx")
		cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")
		cfgCmp.GoogleSheetBillable = o.GoogleSheetBillable

		cfgDataConfig, _ := json.Marshal(cfgCmp)

		h.DS.SetData(cfgRediskey, "$", string(cfgDataConfig))

		err = h.DS.UpdateGoogleSheetBillableCampaignDetail(entity.CampaignDetail{
			GoogleSheetBillable: o.GoogleSheetBillable,
			URLServiceKey:       o.URLServiceKey,
			CampaignId:          o.CampaignId,
		})

		if err != nil {
			h.Logs.Error(fmt.Sprintf("failed updating db: %v", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed save to db"})
		}

		return c.Status(fiber.StatusOK).Send([]byte("OK"))
	}
}

func (h *IncomingHandler) EditMOCappingServiceS2S(c *fiber.Ctx) error {
	o := new(entity.CampaignDetail)
	if err := c.BodyParser(o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

	keys, err := h.DS.GetUrlServiceKeyByService(*o)
	if err != nil || len(keys) == 0 {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": "Data not found"})
	}

	updated := false
	updatedCount := 0

	for _, key := range keys {
		data, err := h.DS.GetDataJSON(key)
		if err != nil {
			continue
		}

		if data["country"] == o.Country &&
			data["operator"] == o.Operator &&
			data["partner"] == o.Partner &&
			data["service"] == o.Service {

			h.DS.SetData(key, "$.mo_capping_service", strconv.Itoa(o.MOCappingService))
			h.DS.SetData(key, "$.status_capping", strconv.FormatBool(false))
			updated = true
			updatedCount++
			h.Logs.Debug(fmt.Sprintf("Updated config key: %s", key))
		}
	}

	if !updated {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data not found"})
	}

	h.Logs.Debug(fmt.Sprintf("Total configs updated: %d", updatedCount))

	err = h.DS.UpdateMOCappingS2S(entity.CampaignDetail{
		MOCappingService: o.MOCappingService,
		LastUpdate:       external.GetCurrentTime(h.Config.TZ, time.RFC3339),
		Country:          o.Country,
		Operator:         o.Operator,
		Partner:          o.Partner,
		Service:          o.Service,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update DB"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       fmt.Sprintf("Success update %d configs", updatedCount),
		"updated_count": updatedCount,
	})
}

func (h *IncomingHandler) EditPOAF(c *fiber.Ctx) error {

	o := new(entity.SummaryCampaign)

	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

	h.DS.EditPOAFIncSummaryCampaign(entity.IncSummaryCampaign{
		SummaryDate:   o.SummaryDate,
		POAF:          o.POAF,
		URLServiceKey: o.URLServiceKey,
	})

	if sum, isOK := h.DS.GetSummaryCampaign(entity.SummaryCampaign{
		SummaryDate:   o.SummaryDate,
		URLServiceKey: o.URLServiceKey,
	}); isOK {

		sum.POAF = o.POAF

		calculated := h.DS.FormulaCPA(sum)

		calculated.POAF = o.POAF
		h.DS.ReCalculateSummaryCampaign(calculated)
	}

	return c.Status(fiber.StatusOK).Send([]byte("OK"))
}

func (h *IncomingHandler) EditCampaignManagementDetail(c *fiber.Ctx) error {

	o := new(entity.CampaignDetail)

	if err := c.BodyParser(&o); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	} else {

		h.Logs.Debug(fmt.Sprintf("data : %#v ...", o))

		// Update to redis with key
		cfgRediskey := external.Concat("-", o.URLServiceKey, "configIdx")
		cfgCmp, _ := h.DS.GetDataConfig(cfgRediskey, "$")

		cfgCmp.APIURL = o.APIURL

		cfgDataConfig, _ := json.Marshal(cfgCmp)

		h.DS.SetData(cfgRediskey, "$", string(cfgDataConfig))

		h.DS.EditCampaignManagementDetail(entity.CampaignDetail{
			APIURL:        cfgCmp.APIURL,
			URLServiceKey: o.URLServiceKey,
			CampaignId:    o.CampaignId,
		})

		return c.Status(fiber.StatusOK).Send([]byte("OK"))
	}
}

func CheckAndAdd(s string, length int) string {
	for len(s) < length {
		s += "0"
	}
	return s
}

func (h *IncomingHandler) CampaignLandingId(campaignId int, obj entity.DataConfig) string {

	return strings.ToUpper(
		external.Concat("-",
			CheckAndAdd(obj.Country, 2)[0:2],
			/* CheckAndAdd(obj.Operator, 3)[0:3],
			CheckAndAdd(obj.Partner, 3)[0:3],
			CheckAndAdd(obj.Service, 3)[0:3],
			CheckAndAdd(obj.Adnet, 3)[0:3], */
			strconv.Itoa(campaignId),
		),
	)
}

func (h *IncomingHandler) UpdateCampaign(c *fiber.Ctx) error {
	var obj entity.DataCampaignAction

	if err := c.BodyParser(&obj); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if obj.Action != "ADD_UPDATE" && obj.Action != "UPDATE" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid action"})
	}

	tx := h.DB.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "failed to begin transaction"})
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	gs, err := h.DS.GetDataConfig("global_setting", "$")
	if err != nil {
		gset, _ := json.Marshal(entity.GlobalSetting{
			CostPerConversion: "0.5",
			AgencyFee:         "0.5",
			TargetDailyBudget: "0.5",
			TechnicalFee:      "0.5",
		})
		h.DS.SetData("global_setting", "$", string(gset))
		gs, _ = h.DS.GetDataConfig("global_setting", "$")
	} else {
		var gset []byte
		if gs.CPCR == "" {
			gs.CPCR = "0.05"
			gset, _ = json.Marshal(entity.GlobalSetting{
				CostPerConversion: gs.CPCR,
				AgencyFee:         gs.AgencyFee,
				TargetDailyBudget: gs.TargetDailyBudget,
				TechnicalFee:      gs.TechnicalFee,
			})
			h.DS.SetData("global_setting", "$", string(gset))
		}
		if gs.AgencyFee == "" {
			gs.AgencyFee = "0.05"
			gset, _ = json.Marshal(entity.GlobalSetting{
				CostPerConversion: gs.CPCR,
				AgencyFee:         gs.AgencyFee,
				TargetDailyBudget: gs.TargetDailyBudget,
				TechnicalFee:      gs.TechnicalFee,
			})
			h.DS.SetData("global_setting", "$", string(gset))
		}
		if gs.TargetDailyBudget == "" {
			gs.TargetDailyBudget = "0.05"
			gset, _ = json.Marshal(entity.GlobalSetting{
				CostPerConversion: gs.CPCR,
				AgencyFee:         gs.AgencyFee,
				TargetDailyBudget: gs.TargetDailyBudget,
				TechnicalFee:      gs.TechnicalFee,
			})
			h.DS.SetData("global_setting", "$", string(gset))
		}
		if gs.TechnicalFee == "" {
			gs.TechnicalFee = "0.05"
			gset, _ = json.Marshal(entity.GlobalSetting{
				CostPerConversion: gs.CPCR,
				AgencyFee:         gs.AgencyFee,
				TargetDailyBudget: gs.TargetDailyBudget,
				TechnicalFee:      gs.TechnicalFee,
			})
			h.DS.SetData("global_setting", "$", string(gset))
		}
	}

	h.Logs.Info(fmt.Sprintf("Parsing data for update: %#v\n", obj))

	h.DS.UpdateCampaign(entity.Campaign{
		CampaignId:        obj.CampaignId,
		Name:              obj.CampaignName,
		CampaignObjective: obj.Objective,
		Country:           obj.Country,
		Advertiser:        obj.Advertiser,
		IsBypassIpRange:   obj.IsBypassIpRange,
	})

	var singleURLData entity.DataSingleURL
	if obj.Objective == "SINGLE URL MAINSTREAM" || obj.Objective == "SINGLE URL S2S" {
		singleURLData = entity.DataSingleURL{
			SingleURLServiceKey: external.Concat("-", "S", obj.CampaignId),
			Country:             obj.Country,
			Objective:           obj.Objective,
			MCC:                 obj.MCC,
			IsBypassIpRange:     obj.IsBypassIpRange,
			Data:                []entity.URLServiceKeyChilds{},
		}
	}

	for _, dc := range obj.DataConfig {
		country := strings.ToUpper(dc.Country)
		operator := strings.ToUpper(dc.Operator)
		service := strings.ToUpper(dc.Service)
		partner := strings.ToUpper(dc.Partner)
		adnet := strings.ToUpper(dc.Adnet)

		var campaign_detail_id int
		if dc.URLServiceKey == "NEW_KEY" {

			var dummy int
			if err := tx.Raw(`
				SELECT 1 
				FROM campaign_details
				WHERE country = ?
				  AND operator = ?
				  AND service = ?
				  AND partner = ?
				  AND adnet = ?
				FOR UPDATE
			`, country, operator, service, partner, adnet).Scan(&dummy).Error; err != nil {
				tx.Rollback()
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			if err := tx.Raw(
				"SELECT nextval('public.campaign_details_id_seq')",
			).Scan(&campaign_detail_id).Error; err != nil {
				tx.Rollback()
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			dc.URLServiceKey = h.CampaignLandingId(
				campaign_detail_id,
				entity.DataConfig{Country: country},
			)
		}

		if obj.Objective == "SINGLE URL MAINSTREAM" || obj.Objective == "SINGLE URL S2S" {
			singleURLData.Data = append(singleURLData.Data, entity.URLServiceKeyChilds{
				URLServiceKey: dc.URLServiceKey,
				Operator:      dc.Operator,
			})
		}

		cfgRediskey := external.Concat("-", dc.URLServiceKey, "configIdx")

		globalSettingTDBKey := strings.ToLower(external.Concat("_", "global_setting_tdb", CheckAndAdd(dc.Country, 2)[0:2], CheckAndAdd(dc.Operator, 3)[0:3]))
		gsettdb, _ := json.Marshal(entity.GlobalSetting{
			TargetDailyBudget: gs.TargetDailyBudget,
		})
		h.DS.SetData(globalSettingTDBKey, "$", string(gsettdb))

		if dc.LastUpdate == "" {
			dc.LastUpdate = external.GetFormatTime(h.Config.TZ, time.RFC3339)
		}
		if dc.LastUpdateCapping == "" {
			dc.LastUpdateCapping = external.GetFormatTime(h.Config.TZ, time.RFC3339)
		}

		urlpostback := h.Config.AppApi + "/v1/postback/?aff_sub={pixel}"

		if !strings.Contains(dc.URLService, "aff_sub") {
			if strings.Contains(dc.URLService, "?") {
				dc.URLService = dc.URLService + "&aff_sub={pixel}"
			} else {
				dc.URLService = dc.URLService + "/?aff_sub={pixel}"
			}
		}

		r := strings.NewReplacer(
			"{adnet}", adnet,
			"{ADNET}", adnet,
			"[adnet]", adnet,
			"[ADNET]", adnet,
		)
		dc.URLService = r.Replace(dc.URLService)

		dc.URLLanding = h.Config.AppHost + "/lp/" + dc.URLServiceKey + "/?aff_sub={pixel}&pubid={pubid}"

		expected := "/lp/" + dc.URLServiceKey + "/?aff_sub={pixel}&pubid={pubid}"
		if !strings.HasSuffix(dc.URLWarpLanding, expected) {
			dc.URLWarpLanding = strings.TrimRight(dc.URLWarpLanding, "/") + expected
		}

		if dc.PostbackMethod == "" {
			dc.PostbackMethod = "PIXEL"
		}

		dataConfig := entity.DataConfig{
			Id:                        dc.Id,
			URLServiceKey:             dc.URLServiceKey,
			CampaignId:                obj.CampaignId,
			CampaignName:              dc.CampaignName,
			Objective:                 obj.Objective,
			Country:                   country,
			Advertiser:                obj.Advertiser,
			Operator:                  operator,
			Partner:                   partner,
			Aggregator:                dc.Aggregator,
			Adnet:                     adnet,
			Service:                   service,
			Keyword:                   dc.Keyword,
			SubKeyword:                dc.SubKeyword,
			IsBillable:                dc.IsBillable,
			Plan:                      dc.Plan,
			PO:                        dc.PO,
			Cost:                      dc.Cost,
			PubId:                     dc.PubId,
			ShortCode:                 dc.ShortCode,
			DeviceType:                dc.DeviceType,
			OS:                        dc.OS,
			URLType:                   dc.URLType,
			ClickType:                 dc.ClickType,
			ClickDelay:                dc.ClickDelay,
			ClientType:                dc.ClientType,
			TrafficSource:             dc.TrafficSource,
			UniqueClick:               dc.UniqueClick,
			URLLanding:                dc.URLLanding,
			URLWarpLanding:            dc.URLWarpLanding,
			URLService:                dc.URLService,
			URLTFCSmartlink:           dc.URLTFCSmartlink,
			GlobPost:                  dc.GlobPost,
			URLGlobPost:               dc.URLGlobPost,
			CustomIntegration:         dc.CustomIntegration,
			IPAddress:                 dc.IPAddress,
			ISP:                       dc.ISP,
			IsActive:                  dc.IsActive,
			MOCapping:                 dc.MOCapping,
			MOCappingService:          dc.MOCappingService,
			CounterMOCapping:          dc.CounterMOCapping,
			CounterMOCappingService:   dc.CounterMOCappingService,
			StatusCapping:             dc.StatusCapping,
			KPIUpperLimitCapping:      dc.KPIUpperLimitCapping,
			IsMachineLearningCapping:  dc.IsMachineLearningCapping,
			RatioSend:                 dc.RatioSend,
			RatioReceive:              dc.RatioReceive,
			CounterMORatio:            dc.CounterMORatio,
			StatusRatio:               dc.StatusRatio,
			KPIUpperLimitRatioSend:    dc.KPIUpperLimitRatioSend,
			KPIUpperLimitRatioReceive: dc.KPIUpperLimitRatioReceive,
			IsMachineLearningRatio:    dc.IsMachineLearningRatio,
			APIURL:                    dc.APIURL,
			LastUpdate:                external.GetFormatTime(h.Config.TZ, time.RFC3339),
			LastUpdateCapping:         external.GetFormatTime(h.Config.TZ, time.RFC3339),
			CPCR:                      gs.CPCR,
			AgencyFee:                 gs.AgencyFee,
			TargetDailyBudget:         gs.TargetDailyBudget,
			TechnicalFee:              gs.TechnicalFee,
			URLPostback:               urlpostback,
			PostbackMethod:            dc.PostbackMethod,
			MainstreamLpType:          dc.MainstreamLpType,
			Title:                     dc.Title,
			TitleOriginal:             dc.TitleOriginal,
			TitleColor:                dc.TitleColor,
			TitleStyle:                dc.TitleStyle,
			TitlePageType:             dc.TitlePageType,
			TitleFontSize:             dc.TitleFontSize,
			SubTitle:                  dc.SubTitle,
			SubTitleOriginal:          dc.SubTitleOriginal,
			SubTitleColor:             dc.SubTitleColor,
			SubTitleStyle:             dc.SubTitleStyle,
			SubTitlePageType:          dc.SubTitlePageType,
			SubTitleFontSize:          dc.SubTitleFontSize,
			BackgroundURL:             dc.BackgroundURL,
			BackgroundColor:           dc.BackgroundColor,
			LogoURL:                   dc.LogoURL,
			URLBanner:                 dc.URLBanner,
			URLBannerOriginal:         dc.URLBannerOriginal,
			Tnc:                       dc.Tnc,
			TncOriginal:               dc.TncOriginal,
			TncColor:                  dc.TncColor,
			TncStyle:                  dc.TncStyle,
			TncPageType:               dc.TncPageType,
			TncFontSize:               dc.TncFontSize,
			ButtonSubscribe:           dc.ButtonSubscribe,
			ButtonSubscribeOriginal:   dc.ButtonSubscribeOriginal,
			ButtonSubscribeColor:      dc.ButtonSubscribeColor,
			StatusSubmitKeyMainstream: dc.StatusSubmitKeyMainstream,
			KeyMainstream:             dc.KeyMainstream,
			Channel:                   dc.Channel,
			GoogleSheet:               dc.GoogleSheet,
			GoogleSheetBillable:       dc.GoogleSheetBillable,
			Currency:                  dc.Currency,
			MCC:                       dc.MCC,
			ClickableAnywhere:         dc.ClickableAnywhere,
			NonTargetURL:              dc.NonTargetURL,
			EnableIpRanges:            dc.EnableIpRanges,
			ConversionName:            dc.ConversionName,
			DomainService:             dc.DomainService,
			CampaignDetailName:        dc.CampaignDetailName,
			Prefix:                    dc.Prefix,
			CountryDialingCode:        dc.CountryDialingCode,
			UnusedTrafficRedirectType: dc.UnusedTrafficRedirectType,
			CompanyLegalName: dc.CompanyLegalName,
			CompanyAddress: dc.CompanyAddress,
			CompanyEmail: dc.CompanyEmail,
			CompanyPhone: dc.CompanyPhone,
			ServicePrice: dc.ServicePrice,
			PortalURL: dc.PortalURL,
			IsEvina: dc.IsEvina,
			EvinaRedirectFraudURL: dc.EvinaRedirectFraudURL,
			IsBypassIpRange: dc.IsBypassIpRange,
		}

		if cd, _ := h.DS.GetCampaignByCampaignDetailId(entity.CampaignDetail{
			URLServiceKey: dc.URLServiceKey, Country: country, Operator: operator, Service: service, Partner: partner, Adnet: adnet,
		}); cd.ID > 0 {
			dataConfig.Id = cd.ID
			cfgDataConfig, _ := json.Marshal(dataConfig)

			cd.Country = country
			cd.Operator = operator
			cd.Partner = partner
			cd.Aggregator = dc.Aggregator
			cd.Adnet = adnet
			cd.Service = service
			cd.Keyword = dc.Keyword
			cd.IsBillable = dc.IsBillable
			cd.Plan = dc.Plan
			cd.PO = dc.PO
			cd.Cost = dc.Cost
			cd.PubId = dc.PubId
			cd.ShortCode = dc.ShortCode
			cd.DeviceType = dc.DeviceType
			cd.OS = dc.OS
			cd.URLType = dc.URLType
			cd.ClickType = dc.ClickType
			cd.ClickDelay = dc.ClickDelay
			cd.ClientType = dc.ClientType
			cd.TrafficSource = dc.TrafficSource
			cd.UniqueClick = dc.UniqueClick
			cd.URLLanding = dc.URLLanding
			cd.URLWarpLanding = dc.URLWarpLanding
			cd.URLService = dc.URLService
			cd.URLTFCORSmartlink = dc.URLTFCSmartlink
			cd.CustomIntegration = dc.CustomIntegration
			cd.APIURL = dc.APIURL
			cd.MainstreamLpType = dc.MainstreamLpType
			cd.Title = dc.Title
			cd.TitleOriginal = dc.TitleOriginal
			cd.TitleColor = dc.TitleColor
			cd.TitleStyle = dc.TitleStyle
			cd.TitlePageType = dc.TitlePageType
			cd.TitleFontSize = dc.TitleFontSize
			cd.SubTitle = dc.SubTitle
			cd.SubTitleOriginal = dc.SubTitleOriginal
			cd.SubTitleColor = dc.SubTitleColor
			cd.SubTitleStyle = dc.SubTitleStyle
			cd.SubTitlePageType = dc.SubTitlePageType
			cd.SubTitleFontSize = dc.SubTitleFontSize
			cd.BackgroundURL = dc.BackgroundURL
			cd.BackgroundColor = dc.BackgroundColor
			cd.LogoURL = dc.LogoURL
			cd.URLBanner = dc.URLBanner
			cd.URLBannerOriginal = dc.URLBannerOriginal
			cd.Tnc = dc.Tnc
			cd.TncOriginal = dc.TncOriginal
			cd.TncColor = dc.TncColor
			cd.TncStyle = dc.TncStyle
			cd.TncPageType = dc.TncPageType
			cd.TncFontSize = dc.TncFontSize
			cd.ButtonSubscribe = dc.ButtonSubscribe
			cd.ButtonSubscribeOriginal = dc.ButtonSubscribeOriginal
			cd.ButtonSubscribeColor = dc.ButtonSubscribeColor
			cd.StatusSubmitKeyMainstream = dc.StatusSubmitKeyMainstream
			cd.KeyMainstream = dc.KeyMainstream
			cd.Channel = dc.Channel
			cd.GoogleSheet = dc.GoogleSheet
			cd.GoogleSheetBillable = dc.GoogleSheetBillable
			cd.Currency = dc.Currency
			cd.MCC = dc.MCC
			cd.ClickableAnywhere = dc.ClickableAnywhere
			cd.NonTargetURL = dc.NonTargetURL
			cd.EnableIpRanges = dc.EnableIpRanges
			cd.ConversionName = dc.ConversionName
			cd.DomainService = dc.DomainService
			cd.CampaignDetailName = dc.CampaignDetailName
			cd.Prefix = dc.Prefix
			cd.CountryDialingCode = dc.CountryDialingCode
			cd.UnusedTrafficRedirectType = dc.UnusedTrafficRedirectType
			cd.CompanyLegalName = dc.CompanyLegalName
			cd.CompanyAddress = dc.CompanyAddress
			cd.CompanyEmail = dc.CompanyEmail
			cd.CompanyPhone = dc.CompanyPhone
			cd.ServicePrice = dc.ServicePrice
			cd.PortalURL = dc.PortalURL
			cd.IsEvina = dc.IsEvina
			cd.EvinaRedirectFraudURL = dc.EvinaRedirectFraudURL

			h.DS.SaveCampaignDetail(cd)
			h.DS.SetData(cfgRediskey, "$", string(cfgDataConfig))
		} else {
			var ip_address []string
			for _, v := range dc.IPAddress {
				ip_address = append(ip_address, strconv.Itoa(int(v)))
			}
			cpcr, _ := strconv.ParseFloat(strings.TrimSpace(dc.CPCR), 64)
			agencyfee, _ := strconv.ParseFloat(strings.TrimSpace(gs.AgencyFee), 64)
			tdb, _ := strconv.ParseFloat(strings.TrimSpace(gs.TargetDailyBudget), 64)
			technical_fee, _ := strconv.ParseFloat(strings.TrimSpace(gs.TechnicalFee), 64)

			h.DS.NewCampaignDetail(entity.CampaignDetail{
				ID:                        campaign_detail_id,
				URLServiceKey:             dc.URLServiceKey,
				CampaignId:                obj.CampaignId,
				Country:                   country,
				Operator:                  operator,
				Partner:                   partner,
				Aggregator:                dc.Aggregator,
				Adnet:                     adnet,
				Service:                   service,
				Keyword:                   dc.Keyword,
				Subkeyword:                dc.SubKeyword,
				IsBillable:                dc.IsBillable,
				Plan:                      dc.Plan,
				PO:                        dc.PO,
				Cost:                      dc.Cost,
				PubId:                     dc.PubId,
				ShortCode:                 dc.ShortCode,
				DeviceType:                dc.DeviceType,
				OS:                        dc.OS,
				URLType:                   dc.URLType,
				ClickType:                 dc.ClickType,
				ClickDelay:                dc.ClickDelay,
				ClientType:                dc.ClientType,
				TrafficSource:             dc.TrafficSource,
				UniqueClick:               dc.UniqueClick,
				URLLanding:                dc.URLLanding,
				URLWarpLanding:            dc.URLWarpLanding,
				URLService:                dc.URLService,
				URLTFCORSmartlink:         dc.URLTFCSmartlink,
				GlobPost:                  dc.GlobPost,
				URLGlobPost:               dc.URLGlobPost,
				CustomIntegration:         dc.CustomIntegration,
				IpAddress:                 ip_address,
				IsActive:                  dc.IsActive,
				MOCapping:                 dc.MOCapping,
				MOCappingService:          dc.MOCappingService,
				CounterMOCapping:          dc.CounterMOCapping,
				CounterMOCappingService:   dc.CounterMOCappingService,
				StatusCapping:             dc.StatusCapping,
				KPIUpperLimitCapping:      dc.KPIUpperLimitCapping,
				IsMachineLearningCapping:  dc.IsMachineLearningCapping,
				RatioSend:                 dc.RatioSend,
				RatioReceive:              dc.RatioReceive,
				CounterMORatio:            dc.CounterMORatio,
				StatusRatio:               dc.StatusRatio,
				KPIUpperLimitRatioSend:    dc.KPIUpperLimitRatioSend,
				KPIUpperLimitRatioReceive: dc.KPIUpperLimitRatioReceive,
				IsMachineLearningRatio:    dc.IsMachineLearningRatio,
				APIURL:                    dc.APIURL,
				LastUpdate:                external.GetCurrentTime(h.Config.TZ, time.RFC3339),
				LastUpdateCapping:         external.GetCurrentTime(h.Config.TZ, time.RFC3339),
				CostPerConversion:         cpcr,
				AgencyFee:                 agencyfee,
				TargetDailyBudget:         tdb,
				TechnicalFee:              technical_fee,
				URLPostback:               urlpostback,
				MainstreamLpType:          dc.MainstreamLpType,
				Title:                     dc.Title,
				TitleOriginal:             dc.TitleOriginal,
				TitleColor:                dc.TitleColor,
				TitleStyle:                dc.TitleStyle,
				TitlePageType:             dc.TitlePageType,
				TitleFontSize:             dc.TitleFontSize,
				SubTitle:                  dc.SubTitle,
				SubTitleOriginal:          dc.SubTitleOriginal,
				SubTitleColor:             dc.SubTitleColor,
				SubTitleStyle:             dc.SubTitleStyle,
				SubTitlePageType:          dc.SubTitlePageType,
				SubTitleFontSize:          dc.SubTitleFontSize,
				BackgroundURL:             dc.BackgroundURL,
				BackgroundColor:           dc.BackgroundColor,
				LogoURL:                   dc.LogoURL,
				URLBanner:                 dc.URLBanner,
				URLBannerOriginal:         dc.URLBannerOriginal,
				Tnc:                       dc.Tnc,
				TncOriginal:               dc.TncOriginal,
				TncColor:                  dc.TncColor,
				TncStyle:                  dc.TncStyle,
				TncPageType:               dc.TncPageType,
				TncFontSize:               dc.TncFontSize,
				ButtonSubscribe:           dc.ButtonSubscribe,
				ButtonSubscribeOriginal:   dc.ButtonSubscribeOriginal,
				ButtonSubscribeColor:      dc.ButtonSubscribeColor,
				StatusSubmitKeyMainstream: dc.StatusSubmitKeyMainstream,
				KeyMainstream:             dc.KeyMainstream,
				Channel:                   dc.Channel,
				GoogleSheet:               dc.GoogleSheet,
				GoogleSheetBillable:       dc.GoogleSheetBillable,
				Currency:                  dc.Currency,
				MCC:                       dc.MCC,
				ClickableAnywhere:         dc.ClickableAnywhere,
				NonTargetURL:              dc.NonTargetURL,
				EnableIpRanges:            dc.EnableIpRanges,
				ConversionName:            dc.ConversionName,
				DomainService:             dc.DomainService,
				CampaignDetailName:        dc.CampaignDetailName,
				Prefix:                    dc.Prefix,
				CountryDialingCode:        dc.CountryDialingCode,
				UnusedTrafficRedirectType: dc.UnusedTrafficRedirectType,
				CompanyLegalName: dc.CompanyLegalName,
				CompanyAddress: dc.CompanyAddress,
				CompanyEmail: dc.CompanyEmail,
				CompanyPhone: dc.CompanyPhone,
				ServicePrice: dc.ServicePrice,
				PortalURL: dc.PortalURL,
				IsEvina: dc.IsEvina,
				EvinaRedirectFraudURL: dc.EvinaRedirectFraudURL,
			})

			dataConfig.Id = campaign_detail_id
			cfgDataConfig, _ := json.Marshal(dataConfig)

			h.DS.SetData(cfgRediskey, "$", string(cfgDataConfig))
		}
	}

	if obj.Objective == "SINGLE URL MAINSTREAM" || obj.Objective == "SINGLE URL S2S" {
		singleURLDataJSON, _ := json.Marshal(singleURLData)
		h.DS.SetData(singleURLData.SingleURLServiceKey, "$", string(singleURLDataJSON))
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Campaign updated successfully"})
}
