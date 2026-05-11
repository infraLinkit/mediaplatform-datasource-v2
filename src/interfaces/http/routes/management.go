package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/interfaces/http/handler"
)

// RegisterManagement: daftar semua sub-group dibawah /v1/management.
// Auth middleware sudah diapply di group level oleh caller (MapUrls).
func RegisterManagement(grp fiber.Router, h *handler.IncomingHandler) {
	registerCampaign(grp.Group("/campaign"), h)
	registerCampaignSetting(grp.Group("/campaign-setting"), h)
	registerMenu(grp.Group("/menu"), h)
	registerRole(grp.Group("/role"), h)
	registerUser(grp.Group("/user"), h)
	registerUserLog(grp.Group("/userlog"), h)
	registerBudgetIO(grp.Group("/budget-io"), h)
	registerCountryService(grp.Group("/country-service"), h)
	registerIPRange(grp.Group("/ipranges"), h)
}

func registerCampaign(grp fiber.Router, h *handler.IncomingHandler) {
	grp.Get("/", h.DisplayCampaignManagement).Name("Campaign Management FE")
	grp.Get("/campaigncounts", h.GetCampaignCounts).Name("Campaign Management Campaign Counts FE")
	grp.Get("/:v", h.DisplayCampaignManagement).Name("Campaign Management Detail FE")
	grp.Post("/send", h.SendCampaignHandler).Name("Campaign Management Send FE")
	grp.Post("/updatestatus", h.UpdateStatusCampaign).Name("Update status campaign on campaign_details")
	grp.Post("/editcampaign", h.EditCampaign).Name("Edit capping campaign on campaign_details")
	grp.Post("/editmocapping", h.EditCampaignMOCapping).Name("Edit mocapping campaign on campaign_details")
	grp.Post("/editratio", h.EditCampaignRatio).Name("Edit ratio campaign on campaign_details")
	grp.Post("/editpo", h.EditCampaignPO).Name("Edit postback campaign on campaign_details")
	grp.Post("/delcampaign", h.DelCampaign).Name("Edit capping campaign on campaign_details")
	grp.Post("/updatekeymainstream", h.UpdateKeyMainstream).Name("Update key mainstream campaign on campaign_details")
	grp.Post("/updategooglesheet", h.UpdateGoogleSheet).Name("Update google sheet campaign on campaign_details")
	grp.Post("/updategooglesheetbillable", h.UpdateGoogleSheetBillable).Name("Update google sheet billable campaign on campaign_details")
	grp.Post("/editmocappingservices2s", h.EditMOCappingServiceS2S).Name("Update mocappingservices2s campaign on campaign_details")
	grp.Post("/editpoaf", h.EditPOAF).Name("Edit poaf")
	grp.Post("/editcampaignmanagementdetail", h.EditCampaignManagementDetail).Name("Edit campaign on campaign_details")
	grp.Post("/updatecampaignmanagement", h.UpdateCampaign).Name("Edit campaign on form")
}

func registerCampaignSetting(grp fiber.Router, h *handler.IncomingHandler) {
	grp.Post("/editratio", h.EditCampaignSettingRatio).Name("Edit ratio campaign on campaign_setting")
	grp.Post("/editpo", h.EditCampaignSettingPO).Name("Edit po campaign on campaign_setting")
	grp.Post("/editmocapping", h.EditCampaignSettingMOCapping).Name("Edit mocapping campaign on campaign_setting")
}

func registerMenu(grp fiber.Router, h *handler.IncomingHandler) {
	grp.Post("/", h.CreateMenu).Name("Menu Management Create FE")
	grp.Get("/", h.GetAllMenus).Name("Menu Management FE")
	grp.Get("/:id", h.GetMenuByID).Name("Menu Management Edit FE")
	grp.Put("/:id", h.UpdateMenu).Name("Menu Management Update FE")
	grp.Delete("/:id", h.DeleteMenu).Name("Menu Management Delete FE")
}

func registerRole(grp fiber.Router, h *handler.IncomingHandler) {
	grp.Post("/", h.CreateRole).Name("Role Management Create FE")
	grp.Get("/", h.GetRoleTable).Name("Role Management FE")
	grp.Put("/:id", h.UpdateRole).Name("Role Management Update FE")
	grp.Delete("/:id", h.DeleteRole).Name("Role Management Delete FE")
}

func registerUser(grp fiber.Router, h *handler.IncomingHandler) {
	grp.Post("/", h.CreateUser).Name("User Management Create FE")
	grp.Get("/", h.GetUserTable).Name("User Management FE")
	grp.Get("/usercounts", h.GetUserCounts).Name("User Management User Counts FE")
	grp.Put("/:id", h.UpdateUser).Name("User Management Update FE")
	grp.Put("/assignservice/:id", h.AssignService).Name("User Management Assign Service & Adnet FE")
	grp.Put("/updatestatus/:id", h.UpdateUserStatus).Name("User Management Update Status FE")
	grp.Delete("/:id", h.DeleteUser).Name("User Management Delete FE")
	grp.Get("/approvalrequest", h.GetUserApplovalRequestTable).Name("User Management Approval Request FE")
	grp.Put("/approveuser/:id", h.ApproveUser).Name("User Management Approve User FE")
}

func registerUserLog(grp fiber.Router, h *handler.IncomingHandler) {
	grp.Get("/ip", h.GetIpAddress).Name("User Log IP")
	grp.Post("/", h.CreateUserLog).Name(" Save User Log List")
	grp.Get("/", h.DisplayUserLogList).Name(" Display User Log List")
	grp.Get("/:id", h.DisplayUserLogHistory).Name(" Display User Log History")
}

func registerBudgetIO(grp fiber.Router, h *handler.IncomingHandler) {
	grp.Post("/", h.CreateBudgetIO).Name("Create BudgetIO")
	grp.Get("/budgetiolist", h.DisplayBudgetIO).Name("Budget IO List")
	grp.Get("/budgetiolistall", h.DisplayBudgetIOAll).Name("Budget IO List All")
	grp.Get("/budgetioapproved", h.DisplayBudgetIOApproved).Name("Budget IO List All")
	grp.Get("/budgetioapprovedall", h.DisplayBudgetIOApprovedAll).Name("Budget IO List All")
}

func registerCountryService(grp fiber.Router, h *handler.IncomingHandler) {
	// Email
	grp.Get("/email", h.DisplayEmail).Name("Display Email")
	grp.Get("/email/:id", h.DisplayEmailByID).Name("Display Email By ID")
	grp.Post("/email", h.CreateEmail).Name("Create Email")
	grp.Put("/email/:id", h.UpdateEmail).Name("Update Email")
	grp.Delete("/email/:id", h.DeleteEmail).Name("Delete Email")
	// Country
	grp.Get("/country", h.DisplayCountry).Name("Create Country")
	grp.Get("/country/:code", h.DisplayCountryInfo).Name("Country Information")
	grp.Post("/country", h.CreateCountry).Name("Create Country")
	grp.Put("/country/:id", h.UpdateCountry).Name("Update Country")
	grp.Delete("/country/:id", h.DeleteCountry).Name("Delete Country")
	grp.Get("/continent", h.DisplayCountry).Name("Display Continent")
	// Company
	grp.Get("/company", h.DisplayCompany).Name("Create Company")
	grp.Post("/company", h.CreateCompany).Name("Create Company")
	grp.Put("/company/:id", h.UpdateCompany).Name("Update Company")
	grp.Delete("/company/:id", h.DeleteCompany).Name("Delete Company")
	grp.Get("/company-group", h.DisplayCompanyGroup).Name("Display Company Group")
	grp.Post("/company-group", h.CreateCompanyGroup).Name("Create Company Group")
	grp.Put("/company-group/:id", h.UpdateCompanyGroup).Name("Update Company Group")
	grp.Delete("/company-group/:id", h.DeleteCompanyGroup).Name("Delete Company Group")
	// Domain
	grp.Get("/domain", h.DisplayDomain).Name("Create Domain")
	grp.Post("/domain", h.CreateDomain).Name("Create Domain")
	grp.Put("/domain/:id", h.UpdateDomain).Name("Update Domain")
	grp.Delete("/domain/:id", h.DeleteDomain).Name("Delete Domain")
	grp.Get("/domain-service", h.DisplayDomainService).Name("Show Domain Service")
	// Operator
	grp.Get("/operator", h.DisplayOperator).Name("Create Operator")
	grp.Post("/operator", h.CreateOperator).Name("Create Operator")
	grp.Put("/operator/:id", h.UpdateOperator).Name("Update Operator")
	grp.Delete("/operator/:id", h.DeleteOperator).Name("Delete Operator")
	grp.Get("/api-operator-list", h.DisplayAPIOperatorList).Name("Show API Operator List")
	// Partner
	grp.Get("/partner", h.DisplayPartner).Name("Create Partner")
	grp.Post("/partner", h.CreatePartner).Name("Create Partner")
	grp.Put("/partner/:id", h.UpdatePartner).Name("Update Partner")
	grp.Delete("/partner/:id", h.DeletePartner).Name("Delete Partner")
	// Service
	grp.Get("/service", h.DisplayService).Name("Create Service")
	grp.Post("/service", h.CreateService).Name("Create Service")
	grp.Put("/service/:id", h.UpdateService).Name("Update Service")
	grp.Delete("/service/:id", h.DeleteService).Name("Delete Service")
	grp.Get("/api-service-list", h.DisplayAPIServiceList).Name("Show API Service List")
	// Adnet
	grp.Get("/adnet-list", h.DisplayAdnetList).Name("Create AdnetList")
	grp.Get("/api-adnet-list", h.DisplayAPIAdnetList).Name("Show API AdnetList")
	grp.Post("/adnet-list", h.CreateAdnetList).Name("Create AdnetList")
	grp.Put("/adnet-list/:id", h.UpdateAdnetList).Name("Update AdnetList")
	grp.Delete("/adnet-list/:id", h.DeleteAdnetList).Name("Delete AdnetList")
	grp.Post("/edit-adnet-dsp-status", h.UpdateDSPAdnetStatus).Name("Edit DSP Status")
	// Agency
	grp.Get("/agency", h.DisplayAgency).Name("Show Agency")
	grp.Post("/agency", h.CreateAgency).Name("Create Agency")
	grp.Put("/agency/:id", h.UpdateAgency).Name("Update Agency")
	grp.Delete("/agency/:id", h.DeleteAgency).Name("Delete Agency")
	// Channel
	grp.Get("/channel", h.DisplayChannel).Name("Show Channel")
	grp.Post("/channel", h.CreateChannel).Name("Create Channel")
	grp.Put("/channel/:id", h.UpdateChannel).Name("Update Channel")
	grp.Delete("/channel/:id", h.DeleteChannel).Name("Delete Channel")
	// Mainstream Group
	grp.Get("/mainstream-group", h.DisplayMainstreamGroup).Name("Show MainStreamGroup")
	grp.Post("/mainstream-group", h.CreateMainstreamGroup).Name("Create MainStreamGroup")
	grp.Put("/mainstream-group/:id", h.UpdateMainstreamGroup).Name("Update MainStreamGroup")
	grp.Delete("/mainstream-group/:id", h.DeleteMainstreamGroup).Name("Delete MainStreamGroup")
}

func registerIPRange(grp fiber.Router, h *handler.IncomingHandler) {
	grp.Get("/", h.GetIPRangeFiles).Name(" Display IP Ranges List List")
	grp.Post("/upload", h.UploadIPRangeRows).Name("Upload IP Ranges CSV")
	grp.Post("/implement", h.ImplementIPRange).Name("Implement IP Ranges")
	grp.Post("/download", h.DownloadIPRangeCSV).Name("Download IP Ranges")
}
