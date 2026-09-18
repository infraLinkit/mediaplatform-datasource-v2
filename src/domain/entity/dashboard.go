package entity

type (
	DisplayDashboard struct {
		Period     string `form:"period" json:"period"`
		Metrics    string `form:"metrics" json:"metrics"`
		Country    string `form:"country" json:"country"`
		Adnet      string `form:"adnet" json:"adnet"`
		Operator   string `form:"operator" json:"operator"`
		Service    string `form:"service" json:"service"`
		Page       int    `form:"page" json:"page"`
		PageSize   int    `form:"page_size" json:"page_size"`
		DateRange  string `form:"date_range" json:"date_range"`
		DateBefore string `form:"date_before" json:"date_before"`
		DateAfter  string `form:"date_after" json:"date_after"`
		Action     string `form:"action" json:"action"`
		ClientType string `form:"client_type" json:"client_type"`
	}

	SummaryDashboardDetail struct {
		Date                    string  `json:"date"`
		TotalMO                 int     `json:"total_mo"`
		TotalActiveAdnet        int     `json:"total_active_adnet"`
		TotalSpending           float64 `json:"total_spending"`
		TotalSaaf               float64 `json:"total_saaf"`
		SpendingToAdnets        float64 `json:"spending_to_adnets"`
		TotalTechnicalFee       float64 `json:"total_technical_fee"`
		TotalS2sSpending        float64 `json:"total_s2s_spending"`
		TotalApiSpending        float64 `json:"total_api_spending"`
		TotalMainstreamSpending float64 `json:"total_mainstream_spending"`
		TotalDSPSpending        float64 `json:"total_dsp_spending"`
		InternalRevenue         float64 `json:"internal_revenue"`
		ExternalRevenue         float64 `json:"external_revenue"`
		InternalSpend           float64 `json:"internal_spend"`
		ExternalSpend           float64 `json:"external_spend"`
		S2SRevenue              float64 `json:"s2s_revenue"`
		MainstreamRevenue       float64 `json:"mainstream_revenue"`
		DSPRevenue              float64 `json:"dsp_revenue"`
		TotalLanding            int     `json:"total_landing"`
		TotalClicked            int     `json:"total_clicked"`
		TotalPostback           int     `json:"total_postback"`
	}

	DetailChartData struct {
		Date                   string  `json:"date"`
		TotalMO                int     `json:"total_mo"`
		TotalSpending          float64 `json:"total_spending"`
		TotalRevenue           float64 `json:"total_revenue"`
		TotalROAS              float64 `json:"total_roas"`
		TotalTechnicalFee      float64 `json:"total_technical_fee"`
		LastMonthTotalMO       int     `json:"lastmonth_total_mo"`
		LastMonthTotalSpending float64 `json:"lastmonth_total_spending"`
		LastMonthTotalRevenue  float64 `json:"lastmonth_total_revenue"`
		LastMonthDate          string  `json:"last_month_date"`
	}

	SummaryDashboardData struct {
		Revenue                 float64 `json:"revenue"`
		TotalSpending           float64 `json:"total_spending"`
		SpendingToAdnets        float64 `json:"spending_to_adnets"`
		TotalTechnicalFee       float64 `json:"total_technical_fee"`
		TargetBudget            float64 `json:"target_budget"`
		Profit                  float64 `json:"profit"`
		MarginPct               float64 `json:"margin_pct"`
		ROAS                    float64 `json:"roas"`
		EstROAS                 float64 `json:"est_roas"`
		ROI                     float64 `json:"roi"`
		EstROI                  float64 `json:"est_roi"`
		InternalRevenue         float64 `json:"internal_revenue"`
		ExternalRevenue         float64 `json:"external_revenue"`
		InternalSpend           float64 `json:"internal_spend"`
		ExternalSpend           float64 `json:"external_spend"`
		InternalROAS            float64 `json:"internal_roas"`
		ExternalROAS            float64 `json:"external_roas"`
		TotalMO                 int     `json:"total_mo"`
		TotalActiveAdnet        int     `json:"total_active_adnet"`
		TotalAdnet              int     `json:"total_adnet"`
		ECPA                    float64 `json:"ecpa"`
		CAC                     float64 `json:"cac"`
		RecoveryDays            float64 `json:"recovery_days"`
		TotalS2SSpending        float64 `json:"total_s2s_spending"`
		TotalAPISpending        float64 `json:"total_api_spending"`
		TotalMainstreamSpending float64 `json:"total_mainstream_spending"`
		TotalDSPSpending        float64 `json:"total_dsp_spending"`
		S2SRevenue              float64 `json:"s2s_revenue"`
		APIRevenue              float64 `json:"api_revenue"`
		MainstreamRevenue       float64 `json:"mainstream_revenue"`
		DSPRevenue              float64 `json:"dsp_revenue"`
		ForecastMO              int     `json:"forecast_mo"`
		ForecastRevenue         float64 `json:"forecast_revenue"`
		ForecastSpending        float64 `json:"forecast_spending"`
		ForecastProfit          float64 `json:"forecast_profit"`
		RunningDays             int     `json:"running_days"`
		DaysInMonth             int     `json:"days_in_month"`
		TotalLanding            int     `json:"total_landing"`
		TotalClicked            int     `json:"total_clicked"`
		TotalPostback           int     `json:"total_postback"`
		DateList                []string          `json:"date_list"`
		DetailChartData         []DetailChartData `json:"detail_chart_data"`
	}

	SummaryDashboardReportDetail struct {
		Date             string  `json:"date"`
		MOReceived       int     `json:"mo_received"`
		MOSent           int     `json:"mo_sent"`
		SpendingToAdnets float64 `json:"spending_to_adnets"`
		Spending         float64 `json:"spending"`
		WAKIRevenue      float64 `json:"waki_revenue"`
	}

	SummaryDashboardReport struct {
		DateRange string                         `json:"date_range"`
		DateList  []string                       `json:"date_list"`
		Detail    []SummaryDashboardReportDetail `json:"detail"`
	}

	SummaryMODetail struct {
		DateNow string `json:"date_now"`
	}

	TopCampaign struct {
		CampaignID    string  `json:"campaign_id"`
		URLServiceKey string  `json:"url_service_key"`
		Country       string  `json:"country"`
		CountryName   string  `json:"country_name"`
		Operator      string  `json:"operator"`
		Service       string  `json:"service"`
		ClientType    string  `json:"client_type"`
		Adnet         string  `json:"adnet"`
		Landing       int     `json:"landing"`
		MO            int     `json:"mo_received"`
		Postback      int     `json:"postback"`
		CRMO          float64 `json:"cr_mo"`
		CRPostback    float64 `json:"cr_postback"`
		URL           string  `json:"url"`
		ECPA          string  `json:"e_cpa"`
		Revenue       float64 `json:"revenue"`
		Spend         float64 `json:"spend"`
		SpendToAdnets float64 `json:"spend_to_adnets"`
		TechnicalFee  float64 `json:"technical_fee"`
		Profit        float64 `json:"profit"`
		ROAS          float64 `json:"roas"`
		EstROAS       float64 `json:"est_roas"`
		ROIMonths     float64 `json:"roi_months"`
		HasROI        bool    `json:"has_roi"`
	}

	SummaryTopBestCampaign struct {
		Campaign []TopCampaign `json:"campaign"`
	}

	SummaryTopWorstCampaign struct {
		Campaign []TopCampaign `json:"campaign"`
	}

	TopPartnerSpend struct {
		Partner string  `json:"partner"`
		Spend   float64 `json:"spend"`
		Pct     float64 `json:"pct"`
	}

	CountryStat struct {
		Country string  `json:"country"`
		Spend   float64 `json:"spend"`
		Revenue float64 `json:"revenue"`
		MO      int     `json:"mo"`
		ROAS    float64 `json:"roas"`
		Share   float64 `json:"share"`
	}

	OpsStats struct {
		AvgLoadTime    float64 `json:"avg_load_time"`
		TotalCampaigns int     `json:"total_campaigns"`
		S2SCampaigns   int     `json:"s2s_campaigns"`
		APICampaigns   int     `json:"api_campaigns"`
		MSCampaigns    int     `json:"ms_campaigns"`
		DSPCampaigns   int     `json:"dsp_campaigns"`
	}

	AlertItem struct {
		Type   string `json:"type"`
		Head   string `json:"head"`
		Meta   string `json:"meta"`
		Action string `json:"action"`
	}

	RollupRow struct {
		Country      string  `json:"country"`
		Operator     string  `json:"operator"`
		Service      string  `json:"service"`
		ClientType   string  `json:"client_type"`
		MO           int     `json:"mo"`
		Spend        float64 `json:"spend"`
		Revenue      float64 `json:"revenue"`
		ROAS         float64 `json:"roas"`
		EstROAS      float64 `json:"est_roas"`
		ROIMonths    float64 `json:"roi_months"`
		HasROI       bool    `json:"has_roi"`
		MarginPct    float64 `json:"margin_pct"`
		RecoveryDays float64 `json:"recovery_days"`
		CAC          float64 `json:"cac"`
	}

	HierarchyCampaignRow struct {
		Country      string  `json:"country"`
		Operator     string  `json:"operator"`
		Service      string  `json:"service"`
		Adnet        string  `json:"adnet"`
		CampaignID   string  `json:"campaign_id"`
		ClientType   string  `json:"client_type"`
		MO           int     `json:"mo"`
		Spend        float64 `json:"spend"`
		Revenue      float64 `json:"revenue"`
		ROAS         float64 `json:"roas"`
		EstROAS      float64 `json:"est_roas"`
		HasEstROAS   bool    `json:"has_est_roas"`
		ROIMonths    float64 `json:"roi_months"`
		HasROI       bool    `json:"has_roi"`
		RecoveryDays float64 `json:"recovery_days"`
		CR           float64 `json:"cr"`
		CPA          float64 `json:"cpa"`
		Payout       float64 `json:"payout"`
		Status       string  `json:"status"`
		Source       string  `json:"source"`
	}

	AdnetStat struct {
		Adnet        string  `json:"adnet"`
		Spend        float64 `json:"spend"`
		Revenue      float64 `json:"revenue"`
		MO           int     `json:"mo"`
		Campaigns    int     `json:"campaigns"`
		ROAS         float64 `json:"roas"`
		EstROAS      float64 `json:"est_roas"`
		ROIMonths    float64 `json:"roi_months"`
		HasROI       bool    `json:"has_roi"`
		RecoveryDays float64 `json:"recovery_days"`
		CAC          float64 `json:"cac"`
	}

	HeatmapCell struct {
		Campaign string  `json:"campaign"`
		Adnet    string  `json:"adnet"`
		Service  string  `json:"service"`
		ROAS     float64 `json:"roas"`
		EstROAS  float64 `json:"est_roas"`
		Spend    float64 `json:"spend"`
		MO       int     `json:"mo"`
	}

	HeatmapData struct {
		Campaigns []string      `json:"campaigns"`
		Adnets    []string      `json:"adnets"`
		Cells     []HeatmapCell `json:"cells"`
	}

	CampaignDailyStat struct {
		Date       string  `json:"date"`
		MO         int     `json:"mo"`
		Spend      float64 `json:"spend"`
		Revenue    float64 `json:"revenue"`
		ROAS       float64 `json:"roas"`
		EstROAS    float64 `json:"est_roas"`
		HasEstROAS bool    `json:"has_est_roas"`
		ROIMonths  float64 `json:"roi_months"`
		HasROI     bool    `json:"has_roi"`
	}

	FilterOptions struct {
		Countries []string `json:"countries"`
		Services  []string `json:"services"`
	}
)
