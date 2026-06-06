package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/interfaces/http/handler"
)

// Postback INTENTIONALLY PUBLIC: adnet/operator callbacks. Auth via signed params/IP allowlist.
func RegisterPostback(grp fiber.Router, h *handler.IncomingHandler) {
	grp.Get("/postback/:urlservicekey/", h.Postback)
	grp.Get("/postback", h.PostbackV3)
	grp.Get("/postback_billed", h.PostbackBilled)
	grp.Get("/postback_sync", h.PostbackDirectReply)
	grp.Get("/inquire/campid", h.InquiryCampID)
	grp.Get("/inquire/api-campid", h.InquiryAPICampID)
}
