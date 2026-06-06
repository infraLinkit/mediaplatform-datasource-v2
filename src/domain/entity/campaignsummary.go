package entity

import (
	"time"
)

type (
	ParamsCampaignSummary struct {
		DataType      string `form:"data-type" json:"data_type"`
		ChartType     string `form:"chart-type" json:"chart_type"`
		ReportType    string `form:"report-type" json:"report_type"`
		Country       string `form:"country" json:"country"`
		Operator      string `form:"operator" json:"operator"`
		PartnerName   string `form:"partner-name" json:"partner-name"`
		UrlServiceKey string `form:"url_service_key" json:"url_service_key"`
		// CampaignName         string   `form:"partner-name" json:"campaign-name"`
		Adnet                string   `form:"adnet" json:"adnet"`
		Service              string   `form:"service" json:"service"`
		CampaignName         string   `form:"campaign-name" json:"campaign_name"`
		DataIndicators       []string `form:"data-indicators" json:"data-indicators"`
		DataBasedOn          string   `form:"data-based-on" json:"data-based-on"`
		DataBasedOnIndicator string   `form:"data-based-on-indicator" json:"data-based-on-indicator"`
		DateRange            string   `form:"date-range" json:"date-range"`
		DateStart            string   `form:"date-start" json:"date-start"`
		DateEnd              string   `form:"date-end" json:"date-end"`
		DateCustomRange      string   `form:"date-custom-range" json:"date-custom-range"`
		All                  string   `form:"custom-range" json:"all"`
	}

	DisplayCPASummary struct {
		Country    string `form:"country" json:"country"`
		Adnet      string `form:"adnet" json:"adnet"`
		Operator   string `form:"operator" json:"operator"`
		CampaignId string `form:"campaign_id" json:"campaign_id"`
		Page       int    `form:"page" json:"page"`
		DateRange  string `form:"date_range" json:"date_range"`
		DateBefore string `form:"date_before" json:"date_before"`
		DateAfter  string `form:"date_after" json:"date_after"`
		Action     string `form:"action" json:"action"`
		SortBy     string `form:"action" json:"sort_by"`
		Draw       int    `form:"draw" json:"draw"`
		PageSize   int    `form:"page_size" json:"page_size"`
	}

	CPASummary struct {
		ID                 int       `gorm:"primaryKey;autoIncrement" json:"id"`
		Status             bool      `gorm:"not null;size:50" json:"status"`
		SummaryDate        time.Time `gorm:"type:date" json:"summary_date"`
		CampaignId         string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		CampaignName       string    `gorm:"not null;size:100" json:"campaign_name"`
		Country            string    `gorm:"not null;size:50" json:"country"`
		Operator           string    `gorm:"not null;size:50" json:"operator"`
		Partner            string    `gorm:"not null;size:50" json:"partner"`
		Adnet              string    `gorm:"not null;size:50" json:"adnet"`
		Service            string    `gorm:"not null;size:50" json:"service"`
		Traffic            int       `gorm:"not null;length:20;default:0" json:"traffic"`
		MoReceived         int       `gorm:"not null;length:20;default:0" json:"mo_received"`
		CPA                float64   `gorm:"type:double precision" json:"cpa"`
		AgencyFee          float64   `gorm:"type:double precision" json:"agency_fee"`
		TargetDailyBudget  float64   `gorm:"type:double precision" json:"target_daily_budget"`
		CrMO               float64   `gorm:"type:double precision" json:"cr_mo"`
		CrPostback         float64   `gorm:"type:double precision" json:"cr_postback"`
		BudgetUsage        float64   `gorm:"type:double precision" json:"budget_usage"`
		WakiRevenue        float64   `gorm:"type:double precision;not null;length:20;default:0" json:"waki_revenue"`
		Fp                 float64   `gorm:"type:double precision;not null;length:20;default:0" json:"fp"`
		MoSent             float64   `gorm:"type:double precision;not null;length:20;default:0" json:"mo_sent"`
		SpendingToAdnets   float64   `gorm:"type:double precision;not null;length:20;default:0" json:"spending_to_adnets"`
		TotalSpending      float64   `gorm:"type:double precision;not null;length:20;default:0" json:"total_spending"`
		TotalWakiAgencyFee float64   `gorm:"type:double precision" json:"total_waki_agency_fee"`
	}

	ReportSummary struct {
		DataIndicators []string
		Total          map[string]float64
		Avg            map[string]float64
		TmoEnd         map[string]float64
		Days           map[string]map[string]map[string]interface{}
	}

	CampaignSummaryMonitoring struct {
		ID                 int       `form:"id" json:"id"`
		UrlServiceKey      string    `form:"url-service-key" json:"url_service_key"`
		Status             bool      `form:"status" json:"status"`
		SummaryDate        time.Time `form:"summary-date" json:"summary_date"`
		CampaignId         string    `form:"campaign-id" json:"campaign_id"`
		CampaignName       string    `form:"campaign-name" json:"campaign_name"`
		Country            string    `form:"country" json:"country"`
		Operator           string    `form:"operator" json:"operator"`
		Partner            string    `form:"partner" json:"partner"`
		Adnet              string    `form:"adnet" json:"adnet"`
		Service            string    `form:"service" json:"service"`
		Traffic            int       `form:"traffic" json:"traffic"`
		MoReceived         int       `form:"mo-received" json:"mo_received"`
		Cpa                float64   `form:"cpa" json:"cpa"`
		AgencyFee          float64   `form:"agency-fee" json:"agency_fee"`
		TargetDailyBudget  float64   `form:"target-daily-budget" json:"target_daily_budget"`
		TargetBudget       float64   `form:"target-budget" json:"target_budget"`
		CrMO               float64   `form:"cr-mo" json:"cr_mo"`
		CrPostback         float64   `form:"cr-postback" json:"cr_postback"`
		BudgetUsage        float64   `form:"budget-usage" json:"budget_usage"`
		WakiRevenue        float64   `form:"waki-revenue" json:"waki_revenue"`
		FirstPush          float64   `form:"fp" json:"fp"`
		MoSent             float64   `form:"mo-sent" json:"mo_sent"`
		SpendingToAdnets   float64   `form:"spending-to-adnets" json:"spending_to_adnets"`
		TotalSpending      float64   `form:"total-spending" json:"total_spending"`
		TotalWakiAgencyFee float64   `form:"total-waki-agency-fee" json:"total_waki_agency_fee"`
		Spending           float64   `form:"spending" json:"spending"`
		Budget             float64   `form:"budget" json:"budget"`
		Mo                 float64   `form:"mo" json:"mo"`
		Cr                 float64   `form:"cr" json:"cr"`
		Revenue            float64   `form:"revenue" json:"revenue"`
	}
	CampaignSummaryChart struct {
		SummaryDate string  `form:"summary-date" json:"summary_date"`
		Mo          float64 `form:"mo" json:"mo"`
		Cr          float64 `form:"cr" json:"cr"`
		Spending    float64 `form:"spending" json:"spending"`
	}

	CampaingCPASummary struct {
		ID                 int       `gorm:"primaryKey;autoIncrement" json:"id"`
		Status             bool      `gorm:"not null;size:50" json:"status"`
		SummaryDate        time.Time `gorm:"type:date" json:"summary_date"`
		CampaignId         string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		CampaignName       string    `gorm:"not null;size:100" json:"campaign_name"`
		Country            string    `gorm:"not null;size:50" json:"country"`
		Operator           string    `gorm:"not null;size:50" json:"operator"`
		Partner            string    `gorm:"not null;size:50" json:"partner"`
		Adnet              string    `gorm:"not null;size:50" json:"adnet"`
		Service            string    `gorm:"not null;size:50" json:"service"`
		Traffic            int       `gorm:"not null;length:20;default:0" json:"traffic"`
		MoReceived         int       `gorm:"not null;length:20;default:0" json:"mo_received"`
		CPA                float64   `gorm:"type:double precision" json:"cpa"`
		AgencyFee          float64   `gorm:"type:double precision" json:"agency_fee"`
		TargetDailyBudget  float64   `gorm:"type:double precision" json:"target_daily_budget"`
		CrMO               float64   `gorm:"type:double precision" json:"cr_mo"`
		CrPostback         float64   `gorm:"type:double precision" json:"cr_postback"`
		BudgetUsage        float64   `gorm:"type:double precision" json:"budget_usage"`
		WakiRevenue        float64   `gorm:"type:double precision;not null;length:20;default:0" json:"waki_revenue"`
		Fp                 float64   `gorm:"type:double precision;not null;length:20;default:0" json:"fp"`
		MoSent             float64   `gorm:"type:double precision;not null;length:20;default:0" json:"mo_sent"`
		SpendingToAdnets   float64   `gorm:"type:double precision;not null;length:20;default:0" json:"spending_to_adnets"`
		TotalSpending      float64   `gorm:"type:double precision;not null;length:20;default:0" json:"total_spending"`
		TotalWakiAgencyFee float64   `gorm:"type:double precision" json:"total_waki_agency_fee"`
	}

	ParamsRevenueMonitoring struct {
		DataType    string `form:"data-type" json:"data_type"`
		ChartType   string `form:"chart-type" json:"chart_type"`
		ReportType  string `form:"report-type" json:"report_type"`
		Country     string `form:"country" json:"country"`
		Operator    string `form:"operator" json:"operator"`
		PartnerName string `form:"partner-name" json:"partner-name"`
		// CampaignName         string   `form:"partner-name" json:"campaign-name"`
		Adnet                string   `form:"adnet" json:"adnet"`
		UrlServiceKey        string   `form:"url_service_key" json:"url_service_key"`
		Service              string   `form:"service" json:"service"`
		CampaignName         string   `form:"campaign-name" json:"campaign_name"`
		CampaignId           string   `form:"campaign-id" json:"campaign_id"`
		TypeData             string   `form:"type-data" json:"type_data"`
		DataIndicators       []string `form:"data-indicators" json:"data-indicators"`
		DataBasedOn          string   `form:"data-based-on" json:"data-based-on"`
		DataBasedOnIndicator string   `form:"data-based-on-indicator" json:"data-based-on-indicator"`
		DateRange            string   `form:"date-range" json:"date-range"`
		DateStart            string   `form:"date-start" json:"date-start"`
		DateEnd              string   `form:"date-end" json:"date-end"`
		DateCustomRange      string   `form:"date-custom-range" json:"date-custom-range"`
		All                  string   `form:"custom-range" json:"all"`
	}

	RevenueMonitoringMonitoring struct {
		ID                 int       `form:"id" json:"id"`
		Status             bool      `form:"status" json:"status"`
		SummaryDate        time.Time `form:"summary-date" json:"summary_date"`
		CampaignId         string    `form:"campaign-id" json:"campaign_id"`
		CampaignName       string    `form:"campaign-name" json:"campaign_name"`
		Country            string    `form:"country" json:"country"`
		Operator           string    `form:"operator" json:"operator"`
		Partner            string    `form:"partner" json:"partner"`
		Adnet              string    `form:"adnet" json:"adnet"`
		Service            string    `form:"service" json:"service"`
		Traffic            int       `form:"traffic" json:"traffic"`
		MoReceived         int       `form:"mo-received" json:"mo_received"`
		Cpa                float64   `form:"cpa" json:"cpa"`
		AgencyFee          float64   `form:"agency-fee" json:"agency_fee"`
		TargetDailyBudget  float64   `form:"target-daily-budget" json:"target_daily_budget"`
		CrMO               float64   `form:"cr-mo" json:"cr_mo"`
		CrPostback         float64   `form:"cr-postback" json:"cr_postback"`
		BudgetUsage        float64   `form:"budget-usage" json:"budget_usage"`
		WakiRevenue        float64   `form:"waki-revenue" json:"waki_revenue"`
		FirstPush          float64   `form:"fp" json:"fp"`
		MoSent             float64   `form:"mo-sent" json:"mo_sent"`
		SpendingToAdnets   float64   `form:"spending-to-adnets" json:"spending_to_adnets"`
		TotalSpending      float64   `form:"total-spending" json:"total_spending"`
		TotalWakiAgencyFee float64   `form:"total-waki-agency-fee" json:"total_waki_agency_fee"`
		Spending           float64   `form:"spending" json:"spending"`
		Budget             float64   `form:"budget" json:"budget"`
		Mo                 float64   `form:"mo" json:"mo"`
		Cr                 float64   `form:"cr" json:"cr"`
	}

	RevenueMonitoringChart struct {
		SummaryDate string  `form:"summary-date" json:"summary_date"`
		Mo          float64 `form:"mo" json:"mo"`
		Revenue     float64 `form:"revenue" json:"revenue"`
		Spending    float64 `form:"spending" json:"spending"`
	}

	DataRevenue struct {
		RevenueMonitoringChart []RevenueMonitoringChart `json:"revenue_monitoring_chart"`
		TimeRevenue            float64                  `form:"time-revenue" json:"time_revenue"`
		TimeInternalRevenue    float64                  `form:"time-internal-revenue" json:"time_internal_revenue"`
		TimeExternalRevenue    float64                  `form:"time-external-revenue" json:"time_external_revenue"`
	}

	DefaultInput struct {
		AgencyFee         float64 `form:"agency-fee" json:"agency_fee"`
		TechnicalFee      float64 `form:"technical-fee" json:"technical_fee"`
		CostPerConversion float64 `form:"cost-per-conversion" json:"cost_per_conversion"`
		TargetDailyBudget float64 `form:"target-daily-budget" json:"target_daily_budget"`
	}

	EditTargetBudgetRequest struct {
		Level    string  `form:"level"     json:"level"`
		Country  string  `form:"country"   json:"country"`
		Operator string  `form:"operator"  json:"operator"`
		Partner  string  `form:"partner"   json:"partner"`
		Service  string  `form:"service"   json:"service"`
		Adnet    string  `form:"adnet"     json:"adnet"`
		Year     int     `form:"year"      json:"year"`
		Month    int     `form:"month"     json:"month"`
		Budget   float64 `form:"budget"    json:"budget"`
		DataType string  `form:"data_type" json:"data_type"`
	}

	BudgetDetailItem struct {
		Country     string  `json:"country"      gorm:"column:country"`
		Operator    string  `json:"operator"     gorm:"column:operator"`
		Partner     string  `json:"partner"      gorm:"column:partner"`
		Service     string  `json:"service"      gorm:"column:service"`
		Adnet       string  `json:"adnet"        gorm:"column:adnet"`
		Year        int     `json:"year"         gorm:"column:year"`
		Month       int     `json:"month"        gorm:"column:month"`
		Budget      float64 `json:"budget"       gorm:"column:budget"`
		Spending    float64 `json:"spending"     gorm:"column:spending"`
		BudgetUsage float64 `json:"budget_usage" gorm:"column:budget_usage"`
	}

	BudgetSummaryItem struct {
		Country  string  `json:"country"  gorm:"column:country"`
		Year     int     `json:"year"     gorm:"column:year"`
		Month    int     `json:"month"    gorm:"column:month"`
		Budget   float64 `json:"budget"   gorm:"column:budget"`
		Spending float64 `json:"spending" gorm:"column:spending"`
	}

	BudgetAggEntry struct {
		Country  string  `gorm:"column:country"`
		Operator string  `gorm:"column:operator"`
		Partner  string  `gorm:"column:partner"`
		Service  string  `gorm:"column:service"`
		Adnet    string  `gorm:"column:adnet"`
		Year     int     `gorm:"column:year"`
		Month    int     `gorm:"column:month"`
		Budget   float64 `gorm:"column:budget"`
	}
)
type Tabler interface {
	TableName() string
}

func (CampaignSummaryMonitoring) TableName() string {
	return "summary_campaigns"
}
