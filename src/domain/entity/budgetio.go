package entity

import "time"

type (
	BudgetIORequest struct {
		CPName            string               `json:"cp_name"`
		PICName           string               `json:"pic_name"`
		ContactEmail      string               `json:"contact_email"`
		CPBusinessPICName string               `json:"cp_business_pic_name"`
		Signature         string               `json:"signature"`
		SubmittedBy       string               `json:"submitted_by"`
		Data              []BudgetIORowRequest `json:"data"`
	}

	BudgetIORowRequest struct {
		IOID               string  `json:"io_id"`
		CampaignType       string  `json:"campaign_type"`
		Month              string  `json:"month"`
		Country            string  `json:"country"`
		CountryName        string  `json:"country_name"`
		Continent          string  `json:"continent"`
		CompanyGroupName   string  `json:"company_group_name"`
		Company            string  `json:"company"`
		Partner            string  `json:"partner"`
		Service            string  `json:"service"`
		TargetCAC          float64 `json:"target_cac"`
		TargetROI          int     `json:"target_roi"`
		MonthlyMOTarget    float64 `json:"monthly_mo_target"`
		MonthlySpendTarget float64 `json:"monthly_spend_target"`
	}

	DisplayBudgetIO struct {
		ID         int       `json:"id"`
		Status     string    `json:"status"`
		VerifiedAt time.Time `json:"verified_at"`

		SubmittedBy string `json:"submitted_by"`
		VerifiedBy  string `json:"verified_by"`

		IOID             string `json:"io_id"`
		CampaignType     string `json:"campaign_type"`
		Month            string `json:"month"`
		Country          string `json:"country"`
		CountryName      string `json:"country_name"`
		Continent        string `json:"continent"`
		CompanyGroupName string `json:"company_group_name"`
		Company          string `json:"company"`
		Partner          string `json:"partner"`
		Service          string `json:"service"`

		TargetCAC          float64 `json:"target_cac"`
		TargetROI          int     `json:"target_roi"`
		MonthlyMOTarget    float64 `json:"monthly_mo_target"`
		MonthlySpendTarget float64 `json:"monthly_spend_target"`

		CPName            string `json:"cp_name"`
		PICName           string `json:"pic_name"`
		ContactEmail      string `json:"contact_email"`
		CPBusinessPICName string `json:"business_pic_name"`
		Signature         string `json:"signature"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`

		// ===== Extra fields for datatables / filtering =====
		Keyword     string `json:"keyword"`
		PageSize    int    `json:"page_size"`
		Page        int    `json:"page"`
		Action      string `json:"action"`
		DateRange   string `json:"date_range"`
		DateBefore  string `json:"date_before"`
		DateAfter   string `json:"date_after"`
		Draw        int    `json:"draw"`
		Reload      string `json:"reload"`
		OrderColumn string `json:"order_column"`
		OrderDir    string `json:"order_dir"`
	}

	TotalBudgetIO struct {
		MonthlySpendTarget float64 `json:"monthly_spend_target"`
	}

	ContinentReport struct {
		Continent string `json:"continent"`
		Month     string `json:"month"`

		ActualCostWeek1Continent float64 `json:"actual_cost_week1_continent"`
		KPIWeek1Continent        float64 `json:"kpi_week1_continent"`
		ActualCostWeek2Continent float64 `json:"actual_cost_week2_continent"`
		KPIWeek2Continent        float64 `json:"kpi_week2_continent"`
		ActualCostWeek3Continent float64 `json:"actual_cost_week3_continent"`
		KPIWeek3Continent        float64 `json:"kpi_week3_continent"`
		ActualCostWeek4Continent float64 `json:"actual_cost_week4_continent"`
		KPIWeek4Continent        float64 `json:"kpi_week4_continent"`

		TotalActualCostContinent float64 `json:"total_actual_cost_continent"`
		TotalIOTargetContinent   float64 `json:"total_io_target_continent"`
		BudgetUsageContinent     float64 `json:"budget_usage_continent"`

		GMVContinent  float64 `json:"gmv_continent"`
		LTVContinent  float64 `json:"ltv_continent"`
		ROASContinent float64 `json:"roas_continent"`
		ROIContinent  float64 `json:"roi_continent"`

		Countries []CountryReport `json:"countries"`
	}

	CountryReport struct {
		Country string `json:"country"`

		ActualCostWeek1Country float64 `json:"actual_cost_week1_country"`
		KPIWeek1Country        float64 `json:"kpi_week1_country"`
		ActualCostWeek2Country float64 `json:"actual_cost_week2_country"`
		KPIWeek2Country        float64 `json:"kpi_week2_country"`
		ActualCostWeek3Country float64 `json:"actual_cost_week3_country"`
		KPIWeek3Country        float64 `json:"kpi_week3_country"`
		ActualCostWeek4Country float64 `json:"actual_cost_week4_country"`
		KPIWeek4Country        float64 `json:"kpi_week4_country"`

		TotalActualCostCountry float64 `json:"total_actual_cost_country"`
		TotalIOTargetCountry   float64 `json:"total_io_target_country"`
		BudgetUsageCountry     float64 `json:"budget_usage_country"`

		GMVCountry  float64 `json:"gmv_country"`
		LTVCountry  float64 `json:"ltv_country"`
		ROASCountry float64 `json:"roas_country"`
		ROICountry  float64 `json:"roi_country"`
	}

	DisplaySummaryBudgetIO struct {
		ID        int    `json:"id"`
		Month     string `json:"month"`
		Continent string `json:"continent"`
		Country   string `json:"country"`
		Company   string `json:"company"`
		Partner   string `json:"partner"`
		Service   string `json:"service"`

		TotalMonthlySpendTarget float64 `json:"total_monthly_spend_target"`
		ActualWeek1             float64 `json:"actual_week_1"`
		ActualWeek2             float64 `json:"actual_week_2"`
		ActualWeek3             float64 `json:"actual_week_3"`
		ActualWeek4             float64 `json:"actual_week_4"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`

		Keyword     string `json:"keyword"`
		PageSize    int    `json:"page_size"`
		Page        int    `json:"page"`
		Action      string `json:"action"`
		DateRange   string `json:"date_range"`
		DateBefore  string `json:"date_before"`
		DateAfter   string `json:"date_after"`
		Draw        int    `json:"draw"`
		Reload      string `json:"reload"`
		OrderColumn string `json:"order_column"`
		OrderDir    string `json:"order_dir"`
	}

	SummaryBudgetIOAgg struct {
		Country      string  `gorm:"column:country"`
		Continent    string  `gorm:"column:continent"`
		Company      string  `gorm:"column:company"`
		Partner      string  `gorm:"column:partner"`
		Operator     string  `gorm:"column:operator"`
		Channel      string  `gorm:"column:channel"`
		Service      string  `gorm:"column:service"`
		ClientType   string  `gorm:"column:client_type"`
		CampaignType string  `gorm:"column:campaign_type"`
		Month        string  `gorm:"column:month"`
		LastDate     string  `gorm:"column:last_date"`

		ActualWeek1 float64 `gorm:"column:actual_week_1"`
		ActualWeek2 float64 `gorm:"column:actual_week_2"`
		ActualWeek3 float64 `gorm:"column:actual_week_3"`
		ActualWeek4 float64 `gorm:"column:actual_week_4"`
		MOWeek1     float64 `gorm:"column:mo_week1"`
		MOWeek2     float64 `gorm:"column:mo_week2"`
		MOWeek3     float64 `gorm:"column:mo_week3"`
		MOWeek4     float64 `gorm:"column:mo_week4"`

		BudgetIOID int     `gorm:"column:budget_io_id"`
		IOTarget   float64 `gorm:"column:io_target"`
		MOTarget   float64 `gorm:"column:mo_target"`
		TargetCAC  float64 `gorm:"column:target_cac"`
		LTV        float64 `gorm:"column:ltv"`
		ROAS       float64 `gorm:"column:roas"`
		ROI        float64 `gorm:"column:roi"`
	}

	IOReportRow struct {
		BudgetIOID int    `json:"budget_io_id"`
		Region     string `json:"region"`
		Country    string `json:"country"`
		Company    string `json:"company"`
		Partner    string `json:"partner"`
		Operator   string `json:"operator"`
		Channel    string `json:"channel"`
		Service    string `json:"service"`
		ClientType string `json:"client_type"`
		Month      string `json:"month"`

		MOWeek1 float64 `json:"mo_week1"`
		MOWeek2 float64 `json:"mo_week2"`
		MOWeek3 float64 `json:"mo_week3"`
		MOWeek4 float64 `json:"mo_week4"`

		CostWeek1 float64 `json:"cost_week1"`
		CostWeek2 float64 `json:"cost_week2"`
		CostWeek3 float64 `json:"cost_week3"`
		CostWeek4 float64 `json:"cost_week4"`

		IOTarget    float64 `json:"io_target"`
		MOTarget    float64 `json:"mo_target"`
		TargetCAC   float64 `json:"target_cac"`
		EstLTV      float64 `json:"est_ltv"`
		EstROAS     float64 `json:"est_roas"`
		ROI         float64 `json:"roi"`
		RecordedDay int     `json:"recorded_day"`
	}

	UpdateSummaryBudgetIORequest struct {
		ID        int      `json:"id"`
		Country   string   `json:"country"`
		Month     string   `json:"month"`
		MOTarget  *float64 `json:"mo_target,omitempty"`
		IOTarget  *float64 `json:"io_target,omitempty"`
		TargetCAC *float64 `json:"target_cac,omitempty"`
		LTV       *float64 `json:"ltv,omitempty"`
		ROAS      *float64 `json:"roas,omitempty"`
		ROI       *float64 `json:"roi,omitempty"`
	}
)
