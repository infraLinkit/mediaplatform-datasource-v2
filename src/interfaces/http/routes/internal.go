package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/interfaces/http/handler"
)

func RegisterInternal(grp fiber.Router, h *handler.IncomingHandler) {
	grp.Put("/setdata/:v/", h.SetData).Name("SetTargetDailyBudget")
	grp.Put("/updatedata/:v/", h.UpdateAgencyFeeAndCostConversion).Name("UpdateAgencyFeeAndCostConversion")
	grp.Put("/updateratio/:v/", h.UpdateRatio).Name("Update Ratio Transactional")
	grp.Put("/updatepostback/:v/", h.UpdatePostback).Name("Update Postback Transactional")
	grp.Put("/updateagencycost/:v", h.UpdateAgencyCost).Name("Update Agency fee and cost per conversion in db")
	grp.Put("/updatestatusalert/:v", h.UpdateStatusAlert).Name("Update Status Alert in db")
	grp.Get("/datasentapipinreport/", h.TrxPinReport).Name("Receive Pin Report Transactional")
	grp.Post("/pinreport/editpayout", h.EditPayoutAPIReport).Name("Edit payout api report")
	grp.Get("/datasentapiperformance/", h.TrxPerformancePinReport).Name("Receive Pin API Performance Report Transactional")
	grp.Post("/pinperformance/editcpa", h.EditCpaAPIPerformanceReport).Name("Edit cpa api performance report")
	grp.Post("/pinperformance/editarpu", h.EditArpuAPIPerformanceReport).Name("Edit arpu api performance report")
	grp.Get("/exportcpa/", h.ExportCpaButton).Name("Export CPA-Report Button")
	grp.Get("/exportcost/", h.ExportCostButton).Name("Export Cost-Report Button")
	grp.Get("/exportcostdetail/", h.ExportCostDetailButton).Name("Export Cost-Report-Detail Button")
	grp.Get("/pinperformance", h.PinPerformance).Name("Receive Pin Performance Report Transactional")
	grp.Post("/uploadexcel", h.UploadExcel).Name("Upload Excal SMS Campaign")
	grp.Put("/updateexcel/:id", h.UpdateExcel).Name("Update Excal SMS Campaign")
	grp.Put("/upsertexcel/", h.UpsertExcel).Name("Upsert Excal SMS Campaign")
	grp.Get("/getdataarpu/", h.GetDataArpu).Name("Get Data ARPU")
	grp.Get("/get_urlservice_in_summarylanding", h.GetURLServiceInSummaryLanding).Name("Get URL service for total load time")
	grp.Put("/update_response_url_service_in_summarylanding", h.UpdateResponseURLServiceInSummaryLanding).Name("Get URL service for total load time")
}
