package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/interfaces/http/handler"
)

func RegisterReport(grp fiber.Router, h *handler.IncomingHandler) {
	grp.Get("/pinreport", h.DisplayPinReport).Name("Pin Report Summary FE")
	grp.Get("/datasentapiperformance", h.DisplayPinPerformanceReport).Name("Pin Performance Api Report Summary FE")
	grp.Get("/cpareportlist", h.DisplayCPAReport).Name("Receive Pin CPA Report Transactional")
	grp.Get("/costreport/:v", h.DisplayCostReport).Name("Receive Pin Cost Report / detail Transactional")
	grp.Get("/conversionlog", h.DisplayConversionLogReport).Name("Conversion Log Report")
	grp.Get("/campaign-monitoring-summary", h.DisplayCampaignSummary).Name("Campaign Summary")
	grp.Get("/campaign-monitoring-summary/chart", h.DisplayCampaignSummaryChart).Name("Campaign Summary Chart")
	grp.Get("/alertreport/:v", h.DisplayAlertReportAll).Name("All Alert Report list/")
	grp.Get("/trafficreport", h.DisplayTrafficReport).Name("Traffic Report list")
	grp.Get("/trafficreport/chart", h.GetTrafficReportChart).Name("Traffic Report chart")
	grp.Get("/mainstreamreport", h.DisplayMainstreamReport).Name("Mainstream Report list")
	grp.Get("/google-traffic-report", h.DisplayGoogleTrafficReport).Name("Google Traffic Report")
	grp.Get("/budgetmonitoring", h.DisplayBudgetMonitoring).Name("Budget Monitoring list")
	grp.Get("/performance-report", h.DisplayPerformanceReport).Name("Performance Report list")
	grp.Get("/revenuemonitoring", h.DisplayRevenueMonitoring).Name("Revenue Monitoring list")
	grp.Get("/revenuemonitoring/chart", h.DisplayRevenueMonitoringChart).Name("Revenue Monitoring chart")
	grp.Get("/defaultinput/", h.DisplayDefaultInput).Name("Default Input for cpa n mainstream")
	grp.Get("/redirectiontime", h.DisplayRedirectionTime).Name("Redirection Time")
	grp.Post("/resend-data", h.ResendData).Name("Resend Data")
	grp.Post("/resend-data-apireport", h.ResendDataAPIReport).Name("Resend Data API Report")
	grp.Get("/ioreport", h.DisplaySummaryBudgetIO).Name("IO Report")
	grp.Put("/ioreport/update", h.UpdateSummaryBudgetIO).Name("IO Report Update")
	grp.Post("/campaign-monitoring-summary/edit-target-budget", h.EditTargetBudgetLevel).Name("Edit Target Budget")
	grp.Post("/campaign-monitoring-summary/edit-target-budget-batch", h.EditTargetBudgetBatch).Name("Edit Target Budget Batch")
	grp.Get("/campaign-spending-channel", h.DisplayCampaignSpendingChannel).Name("Campaign Spending Channel")
	grp.Get("/campaign-spending-channel/country-children", h.DisplayCampaignSpendingChannelCountryChildren).Name("Campaign Spending Channel Country Children")
}
