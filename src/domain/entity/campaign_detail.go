package entity

import (
	"time"

	"gorm.io/gorm"
)

type (
	CampaignDetail struct {
		gorm.Model
		ID                        int       `gorm:"primaryKey;autoIncrement" json:"id"`
		URLServiceKey             string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId                string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country                   string    `gorm:"not null;size:50" json:"country"`
		Operator                  string    `gorm:"not null;size:50" json:"operator"`
		Partner                   string    `gorm:"not null;size:50" json:"partner"`
		Aggregator                string    `gorm:"not null;size:50" json:"aggregator"`
		Adnet                     string    `gorm:"not null;size:50" json:"adnet"`
		Service                   string    `gorm:"not null;size:50" json:"service"`
		Keyword                   string    `gorm:"not null;size:50" json:"keyword"`
		Subkeyword                string    `gorm:"not null;size:50" json:"subkeyword"`
		IsBillable                bool      `gorm:"not null;size:50" json:"is_billable"`
		Plan                      string    `gorm:"not null;size:50" json:"plan"`
		PO                        string    `gorm:"size:50" json:"po"`
		Cost                      string    `gorm:"not null;size:50" json:"cost"`
		PubId                     string    `gorm:"not null;size:50" json:"pubid"`
		ShortCode                 string    `gorm:"not null;size:50" json:"short_code"`
		DeviceType                string    `gorm:"not null;size:50" json:"device_type"`
		OS                        string    `gorm:"not null;size:50" json:"os"`
		URLType                   string    `gorm:"not null;size:50;default:wap" json:"url_type"`
		ClickType                 int       `gorm:"not null;length:1;default:1" json:"click_type"`
		ClickDelay                int       `gorm:"not null;length:1;default:1" json:"click_delay"`
		ClientType                string    `gorm:"not null;size:50" json:"client_type"`
		TrafficSource             bool      `gorm:"not null;default:false" json:"traffic_source"`
		UniqueClick               bool      `gorm:"not null;default:false" json:"unique_click"`
		URLBanner                 string    `gorm:"type:text;default:NA" json:"url_banner"`
		URLLanding                string    `gorm:"size:255;default:NA" json:"url_landing"`
		URLWarpLanding            string    `gorm:"size:255;default:NA" json:"url_warp_landing"`
		URLService                string    `gorm:"size:255;default:NA" json:"url_service"`
		URLTFCORSmartlink         string    `gorm:"size:255;default:NA" json:"url_tfc_or_smartlink"`
		GlobPost                  bool      `gorm:"not null;default:false" json:"glob_post"`
		URLGlobPost               string    `gorm:"size:255;default:NA" json:"url_glob_post"`
		CustomIntegration         string    `gorm:"size:30;default:NA" json:"custom_integration"`
		IpAddress                 []string  `gorm:"type:text[]"`
		IsActive                  bool      `gorm:"not null;default:true" json:"is_active"`
		MOCapping                 int       `gorm:"not null;length:10;default:0" json:"mo_capping"`
		MOCappingService          int       `gorm:"not null;length:10;default:0" json:"mo_capping_service"`
		CounterMOCapping          int       `gorm:"not null;length:10;default:0" json:"counter_mo_capping"`
		CounterMOCappingService   int       `gorm:"not null;length:10;default:0" json:"counter_mo_capping_service"`
		StatusCapping             bool      `gorm:"not null;default:false" json:"status_capping"`
		KPIUpperLimitCapping      int       `gorm:"not null;length:20;default:1" json:"kpi_upper_limit_capping"`
		IsMachineLearningCapping  bool      `gorm:"not null;default:false" json:"is_machine_learning_capping"`
		RatioSend                 int       `gorm:"not null;length:10;default:1" json:"ratio_send"`
		RatioReceive              int       `gorm:"not null;length:10;default:4" json:"ratio_receive"`
		CounterMORatio            int       `gorm:"not null;length:10;default:0" json:"counter_mo_ratio"`
		StatusRatio               bool      `gorm:"not null;default:false" json:"status_ratio"`
		KPIUpperLimitRatioSend    int       `gorm:"not null;length:10;default:1" json:"kpi_upper_limit_ratio_send"`
		KPIUpperLimitRatioReceive int       `gorm:"not null;length:10;default:4" json:"kpi_upper_limit_ratio_receive"`
		IsMachineLearningRatio    bool      `gorm:"not null;default:false" json:"is_machine_learning_ratio"`
		APIURL                    string    `gorm:"type:text" json:"api_url"`
		LastUpdate                time.Time `json:"last_update"`
		LastUpdateCapping         time.Time `json:"last_update_capping"`
		CostPerConversion         float64   `gorm:"type:double precision" json:"cost_per_conversion"`
		AgencyFee                 float64   `gorm:"type:double precision" json:"agency_fee"`
		TargetDailyBudget         float64   `gorm:"type:double precision" json:"target_daily_budget"`
		TechnicalFee              float64   `gorm:"type:double precision" json:"technical_fee"`
		URLPostback               string    `gorm:"size:255;default:NA" json:"url_postback"`
		MainstreamLpType          string    `gorm:"size:50;default:NA" json:"mainstream_lp_type"`
		Title                     string    `gorm:"type:text;default:NA" json:"title"`
		TitleOriginal             string    `gorm:"type:text;default:NA" json:"title_original"`
		TitleColor                string    `gorm:"size:50;default:NA" json:"title_color"`
		TitleStyle                string    `gorm:"size:50;default:NA" json:"title_style"`
		TitlePageType             string    `gorm:"size:50;default:NA" json:"title_page_type"`
		TitleFontSize             string    `gorm:"size:50;default:NA" json:"title_font_size"`
		SubTitle                  string    `gorm:"type:text;default:NA" json:"sub_title"`
		SubTitleOriginal          string    `gorm:"type:text;default:NA" json:"sub_title_original"`
		SubTitleColor             string    `gorm:"size:50;default:NA" json:"sub_title_color"`
		SubTitleStyle             string    `gorm:"size:50;default:NA" json:"sub_title_style"`
		SubTitlePageType          string    `gorm:"size:50;default:NA" json:"sub_title_page_type"`
		SubTitleFontSize          string    `gorm:"size:50;default:NA" json:"sub_title_font_size"`
		BackgroundURL             string    `gorm:"type:text;default:NA" json:"background_url"`
		BackgroundColor           string    `gorm:"size:50;default:NA" json:"background_color"`
		LogoURL                   string    `gorm:"type:text;default:NA" json:"logo_url"`
		URLBannerOriginal         string    `gorm:"type:text;default:NA" json:"url_banner_original"`
		Tnc                       string    `gorm:"type:text;default:NA" json:"tnc"`
		TncOriginal               string    `gorm:"type:text;default:NA" json:"tnc_original"`
		TncColor                  string    `gorm:"size:50;default:NA" json:"tnc_color"`
		TncStyle                  string    `gorm:"size:50;default:NA" json:"tnc_style"`
		TncPageType               string    `gorm:"size:50;default:NA" json:"tnc_page_type"`
		TncFontSize               string    `gorm:"size:50;default:NA" json:"tnc_font_size"`
		ButtonSubscribe           string    `gorm:"type:text;default:NA" json:"button_subscribe"`
		ButtonSubscribeOriginal   string    `gorm:"type:text;default:NA" json:"button_subscribe_original"`
		ButtonSubscribeColor      string    `gorm:"size:100;default:NA" json:"button_subscribe_color"`
		StatusSubmitKeyMainstream bool      `gorm:"not null;default:false" json:"status_submit_key_mainstream"`
		KeyMainstream             string    `gorm:"size:50;default:NA" json:"key_mainstream"`
		Channel                   string    `gorm:"size:50;default:NA" json:"channel"`
		GoogleSheet               string    `gorm:"type:text;default:NA" json:"google_sheet"`
		GoogleSheetBillable       string    `gorm:"type:text;default:NA" json:"google_sheet_billable"`
		Currency                  string    `gorm:"size:10;default:NA" json:"currency"`
		MCC                       string    `gorm:"size:10;default:NA" json:"mcc"`
		ClickableAnywhere         bool      `gorm:"not null;default:false" json:"clickable_anywhere"`
		NonTargetURL              string    `gorm:"type:text;default:NA" json:"non_target_url"`
		EnableIpRanges            bool      `gorm:"not null;default:false" json:"enable_ip_ranges"`
		ConversionName            string    `gorm:"size:50;default:NA" json:"conversion_name"`
		DomainService             string    `gorm:"type:varchar(80)" json:"domain_service"`
		CampaignDetailName        string    `gorm:"type:varchar(80)" json:"campaign_detail_name"`
		Prefix                    string    `gorm:"type:varchar(80)" json:"prefix"`
		CountryDialingCode        string    `gorm:"type:varchar(80)" json:"country_dialing_code"`
		UnusedTrafficRedirectType string    `gorm:"type:varchar(80);default:NA" json:"unused_traffic_redirect_type"`
		CompanyLegalName          string    `gorm:"type:varchar(100);default:NA" json:"company_legal_name"`
		CompanyAddress            string    `gorm:"type:text;default:NA" json:"company_address"`
		CompanyEmail              string    `gorm:"type:varchar(100);default:NA" json:"company_email"`
		CompanyPhone              string    `gorm:"type:varchar(20);default:NA" json:"company_phone"`
		ServicePrice              float64   `gorm:"type:double precision;default:0" json:"service_price"`
		PortalURL                 string    `gorm:"type:text;default:NA" json:"portal_url"`
		IsEvina                   bool      `gorm:"not null;default:false" json:"is_evina"`
		EvinaRedirectFraudURL     string    `gorm:"type:text;default:NA" json:"evina_redirect_fraud_url"`
		CRThreshold               float64   `gorm:"type:double precision;default:0" json:"cr_threshold"`
		ApiKey                    string    `gorm:"column:api_key;size:50" json:"api_key"`
		CreatedAt                 time.Time
		UpdatedAt                 time.Time
	}

	HistoryCappingKey struct {
		gorm.Model
		ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
		URLServiceKey string    `gorm:"index:idx_urlservicekey;not null;size:50;uniqueIndex:idx_hck_conflict_key" json:"urlservicekey"`
		CreatedAt     time.Time `gorm:"uniqueIndex:idx_hck_conflict_key"`
		UpdatedAt     time.Time
	}
)

// HOOK or Trigger
// func (cd *CampaignDetail) AfterUpdate(db *gorm.DB) (err error) {
// 	if cd.CounterMOCapping >= cd.MOCapping {
// 		db.Model(&CampaignDetail{}).Where("id = ?", cd.ID).Update("last_update_capping", cd.LastUpdate)
// 	}

// 	return nil
// }

// Important Create Name Func per hook entity!, Name returns the name of the plugin
/* func (cd *CampaignDetail) Name() string {
	return "campaign_detail_trigger"
}

// Initialize is called when you register the plugin via db.Use()
func (cd *CampaignDetail) Initialize(db *gorm.DB) error {
	// Register a callback before the standard Create operation

	//db.Callback().Create().After("gorm:insert").Register(cd.Name(), cd.AfterUpdate)
	db.Callback().Update().After("gorm:update").Register("campaign_detail_trigger", cd.AfterSave)

	return nil
}

func (cd *CampaignDetail) AfterUpdate(db *gorm.DB) {
	if cd.CounterMOCapping >= cd.MOCapping {
		db.Model(&CampaignDetail{}).Where("id = ?", cd.ID).Update("last_update_capping", cd.LastUpdate)
	}

	cd.SummaryContainer(db)
}

func (cd *CampaignDetail) AfterSave(db *gorm.DB) {

	if db.Error == nil {
		// You can access the data being created:
		// model := db.Statement.Dest
		//if _, ok := db.Statement.Dest.(*CampaignDetail); ok {

		cd.SummaryContainer(db)
		//}
	}
}

func (cd *CampaignDetail) SummaryContainer(tx *gorm.DB) {

	var (
		c  Campaign
		o  IncSummaryCampaign
		sc SummaryCampaign
	)

	fmt.Printf("--> Campaign Detail AfterSave Hook, %#v\n", cd)

	tx.Model(&Campaign{}).
		Where("id = ?", cd.CampaignId).First(&c)

	curdate_time, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

	pos, _ := strconv.ParseFloat(strings.TrimSpace(cd.PO), 64)

	o = IncSummaryCampaign{
		SummaryDate:       curdate_time,
		URLServiceKey:     cd.URLServiceKey,
		CampaignId:        cd.CampaignId,
		CampaignObjective: c.CampaignObjective,
		Country:           cd.Country,
		Operator:          cd.Operator,
		Partner:           cd.Partner,
		Aggregator:        cd.Aggregator,
		Adnet:             cd.Adnet,
		Service:           cd.Service,
		ShortCode:         cd.ShortCode,
		Landing:           0,
		MoReceived:        0,
		Postback:          0,
		POAF:              0,
		CreatedAt:         curdate_time,
	}

	sc = SummaryCampaign{
		Status:             cd.IsActive,
		SummaryDate:        curdate_time,
		CampaignId:         cd.CampaignId,
		CampaignName:       c.Name,
		Company:            c.Advertiser,
		Country:            cd.Country,
		Partner:            cd.Partner,
		Operator:           cd.Operator,
		Aggregator:         cd.Aggregator,
		Service:            cd.Service,
		Adnet:              cd.Adnet,
		ShortCode:          cd.ShortCode,
		Traffic:            0,
		Landing:            0,
		MoReceived:         0,
		CR:                 0,
		Postback:           0,
		TotalFP:            0,
		SuccessFP:          0,
		Billrate:           0,
		PO:                 pos,
		Cost:               0,
		SBAF:               0,
		SAAF:               0,
		CPA:                0,
		Revenue:            0,
		URLAfter:           cd.URLWarpLanding,
		URLBefore:          cd.URLLanding,
		MOLimit:            cd.MOCapping,
		ROI:                0,
		RatioSend:          cd.RatioSend,
		RatioReceive:       cd.RatioReceive,
		ClientType:         cd.ClientType,
		CostPerConversion:  0,
		AgencyFee:          0,
		TotalWakiAgencyFee: 0,
		TargetDailyBudget:  0,
		BudgetUsage:        0,
		TechnicalFee:       0,
		URLServiceKey:      cd.URLServiceKey,
		CampaignObjective:  c.CampaignObjective,
		Channel:            cd.Channel,
		PricePerMO:         0,
		CrMO:               0,
		CrPostback:         0,
		POAF:               pos,
	}

	result := tx.Model(&o).
		Where("summary_date = '"+o.SummaryDate.Format("2006-01-02")+"' AND url_service_key = ?", o.URLServiceKey).
		First(&o)

	b := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if b { // Data Not Found

		result := tx.Create(&o)

		fmt.Printf("NewIncSummaryCampaign :%s-%s, affected: %d, is error : %#v\n", o.URLServiceKey, o.SummaryDate, result.RowsAffected, result.Error)

		result = tx.Create(&sc)

		fmt.Printf("NewSummaryCampaign :%s-%s, affected: %d, is error : %#v\n", o.URLServiceKey, o.SummaryDate, result.RowsAffected, result.Error)
	}
} */
