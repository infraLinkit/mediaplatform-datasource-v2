package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/interfaces/http/handler"
)

func RegisterDashboard(grp fiber.Router, h *handler.IncomingHandler) {
	grp.Get("get-data", h.DisplayDashboardData)
	grp.Get("get-top-campaign", h.DisplayDashboardTopCampaign)
	grp.Get("get-report-list", h.DisplayDashboardReport)
	grp.Get("get-country-stats", h.DisplayCountryStats)
	grp.Get("get-ops-stats", h.DisplayOpsStats)
	grp.Get("get-alerts", h.DisplayAlerts)
	grp.Get("get-rollup", h.DisplayRollup)
	grp.Get("get-campaign-hierarchy", h.DisplayCampaignHierarchy)
	grp.Get("get-adnet-stats", h.DisplayAdnetStats)
	grp.Get("get-heatmap", h.DisplayHeatmap)
	grp.Get("get-campaign-daily", h.DisplayCampaignDaily)
	grp.Get("get-service-daily", h.DisplayServiceDaily)
	grp.Get("get-filter-options", h.DisplayFilterOptions)
}
