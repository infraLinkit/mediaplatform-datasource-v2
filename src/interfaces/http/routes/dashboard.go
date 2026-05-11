package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/interfaces/http/handler"
)

func RegisterDashboard(grp fiber.Router, h *handler.IncomingHandler) {
	grp.Get("get-data", h.DisplayDashboardData)
	grp.Get("get-top-campaign", h.DisplayDashboardTopCampaign)
	grp.Get("get-report-list", h.DisplayDashboardReport)
}
