package entity

import (
	"time"

	"github.com/lib/pq"
)

type (
	DisplayCampaignManagement struct {
		Country      string `form:"country" json:"country"`
		Adnet        string `form:"adnet" json:"adnet"`
		Operator     string `form:"operator" json:"operator"`
		Service      string `form:"service" json:"service"`
		Status       string `form:"status" json:"status"`
		Partner      string `form:"partner" json:"partner"`
		CampaignName string `form:"campaign_name" json:"campaign_name"`
		CampaignType string `form:"campaign_type" json:"campaign_type"`
		CampaignId   string `form:"campaign_id" json:"campaign_id"`
		Page         int    `form:"page" json:"page"`
		Draw         int    `form:"draw" json:"draw"`
		Action       string `form:"action" json:"action"`
		URLServiceKey     string `form:"url_service_key" json:"url_service_key"`
		OrderColumn       string `form:"order_column" json:"order_column"`
		OrderDir          string `form:"order_dir" json:"order_dir"`
		CreatedDateBefore string `form:"created_date_before" json:"created_date_before"`
		CreatedDateAfter  string `form:"created_date_after" json:"created_date_after"`
	}

	CampaignManagementData struct {
		ID                pq.Int64Array  `gorm:"type:bigint[]" json:"id"`
		CampaignID        string         `json:"campaign_id"`
		CampaignName      string         `json:"campaign_name"`
		CampaignObjective string         `json:"campaign_objective"`
		Country           string         `json:"country"`
		Partner           string         `json:"partner"`
		TotalOperator     int            `json:"total_operator"`
		Service           int            `json:"service"`
		TotalAdnet        int            `json:"total_adnet"`
		ShortCode         int            `json:"short_code"`
		IsActive          bool           `json:"is_active"`
		URLServiceKey     pq.StringArray `gorm:"type:text[]" json:"url_service_key"`
		CreatedAt         time.Time      `json:"created_at"`
	}

	CampaignManagementDetail struct {
		ID                int            `json:"id"`
		CampaignID        string         `json:"campaign_id"`
		CampaignName      string         `json:"campaign_name"`
		CampaignObjective string         `json:"campaign_objective"`
		Country           string         `json:"country"`
		Operator          string         `json:"operator"`
		Service           string         `json:"service"`
		Adnet             string         `json:"adnet"`
		Partner           string         `json:"partner"`
		ShortCode         string         `json:"short_code"`
		MOLimit           int            `json:"mo_limit"`
		Payout            string         `json:"po"`
		RatioSend         int            `json:"ratio_send"`
		RatioReceive      int            `json:"ratio_receive"`
		URLPostback       string         `json:"url_postback"`
		URLService        string         `json:"url_service"`
		URLanding         string         `json:"url_landing"`
		URLWarpLanding    string         `json:"url_warp_landing"`
		APIURL            string         `json:"api_url"`
		IsActive          bool           `json:"is_active"`
		UrlServiceKey     string         `json:"url_service_key"`
		URLType           string         `json:"url_type"`
		DeviceType        string         `json:"device_type"`
		Channel           string         `json:"channel"`
		CCEmail           pq.StringArray `json:"cc_email"`
		IsBillable        bool           `json:"is_billable"`
		CreatedAt         time.Time      `json:"created_at"`
	}

	CampaignManagementDataDetail struct {
		Operator string                     `json:"operator"`
		Service  string                     `json:"service"`
		Details  []CampaignManagementDetail `json:"details"`
	}

	CampaignCounts struct {
		TotalCampaigns          int `json:"total_campaign"`
		TotalActiveCampaigns    int `json:"total_active_campaign"`
		TotalNonActiveCampaigns int `json:"total_inactive_campaign"`
	}
)
