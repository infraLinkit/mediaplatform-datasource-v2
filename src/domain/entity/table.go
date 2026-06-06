package entity

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type (
	Campaign struct {
		gorm.Model
		ID                int    `gorm:"primaryKey;autoIncrement" json:"id"`
		CampaignId        string `gorm:"index:idx_campid_unique;not null;size:50" json:"campaign_id"`
		Name              string `gorm:"not null;size:50" json:"name"`
		CampaignObjective string `gorm:"not null;size:50" json:"campaign_objective"`
		Country           string `gorm:"not null;size:50" json:"country"`
		Advertiser        string `gorm:"not null;size:50" json:"advertiser"`
		/* SingleURLServiceKey string `gorm:"not null;size:50" json:"singleurlservicekey"`
		URLLanding          string `gorm:"type:text;default:NA" json:"url_landing"`
		URLWarpLanding      string `gorm:"type:text;default:NA" json:"url_warp_landing"`
		MCC                 string `gorm:"not null;size:50" json:"mcc"` */
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	LpDesignType struct {
		gorm.Model
		ID                int    `gorm:"primaryKey;autoIncrement" json:"id"`
		MainstreamLPType  string `gorm:"size:50" json:"mainstream_lp_type"`
		Description       string `gorm:"not null;size:255;default:NA" json:"description"`
		SampleImageBanner string `gorm:"not null;size:255" json:"sample_image_banner"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	ResultCampaign struct {
		Name                      string
		CampaignObjective         string
		Advertiser                string
		ID                        int
		URLServiceKey             string
		CampaignId                string
		Country                   string
		Operator                  string
		Partner                   string
		Aggregator                string
		Adnet                     string
		Service                   string
		Keyword                   string
		Subkeyword                string
		IsBillable                bool
		Plan                      string
		PO                        string
		Cost                      string
		PubId                     string
		ShortCode                 string
		DeviceType                string
		OS                        string
		URLType                   string
		ClickType                 int
		ClickDelay                int
		ClientType                string
		TrafficSource             bool
		UniqueClick               bool
		URLBanner                 string
		URLLanding                string
		URLWarpLanding            string
		URLService                string
		URLTFCORSmartlink         string
		GlobPost                  bool
		URLGlobPost               string
		CustomIntegration         string
		IpAddress                 []string
		IsActive                  bool
		MOCapping                 int
		CounterMOCapping          int
		StatusCapping             bool
		KPIUpperLimitCapping      int
		IsMachineLearningCapping  bool
		RatioSend                 int
		RatioReceive              int
		CounterMORatio            int
		StatusRatio               bool
		KPIUpperLimitRatioSend    int
		KPIUpperLimitRatioReceive int
		IsMachineLearningRatio    bool
		APIURL                    string
		LastUpdate                time.Time
		LastUpdateCapping         time.Time
		CostPerConversion         float64
		AgencyFee                 float64
		TargetDailyBudget         float64
		TechnicalFee              float64
		URLPostback               string
		Channel                   string
	}

	MO struct {
		gorm.Model
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		CampaignDetailId  int       `gorm:"index:idx_mocampdetailid_unique;not null;length:20" json:"campaign_detail_id"`
		Pxdate            time.Time `gorm:"not null" json:"pxdate"`
		URLServiceKey     string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId        string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country           string    `gorm:"not null;size:50" json:"country"`
		Operator          string    `gorm:"not null;size:50" json:"operator"`
		Partner           string    `gorm:"not null;size:50" json:"partner"`
		Aggregator        string    `gorm:"size:50" json:"aggregator"`
		Adnet             string    `gorm:"not null;size:50" json:"adnet"`
		Service           string    `gorm:"not null;size:50" json:"service"`
		Keyword           string    `gorm:"size:50" json:"keyword"`
		Subkeyword        string    `gorm:"size:50" json:"subkeyword"`
		IsBillable        bool      `gorm:"default:false" json:"is_billable"`
		Plan              string    `gorm:"size:50;default:NA" json:"plan"`
		PO                string    `gorm:"size:50;default:NA" json:"po"`
		Cost              string    `gorm:"size:50;default:NA" json:"cost"`
		PubId             string    `gorm:"size:50;default:NA" json:"pubid"`
		ShortCode         string    `gorm:"size:50;default:NA" json:"short_code"`
		URL               string    `gorm:"size:255;default:NA" json:"url"`
		URLType           string    `gorm:"size:50;default:NA" json:"url_type"`
		Pixel             string    `gorm:"size:255;default:NA" json:"pixel"`
		Token             string    `gorm:"size:255;default:NA" json:"token"`
		TrxId             string    `gorm:"size:255;default:NA" json:"trx_id"`
		Msisdn            string    `gorm:"size:255;default:NA" json:"msisdn"`
		IsUsed            bool      `gorm:"not null;default:false" json:"is_used"`
		Browser           string    `gorm:"size:150;default:NA" json:"browser"`
		OS                string    `gorm:"size:150;default:NA" json:"os"`
		Ip                string    `gorm:"type:text" json:"ip"`
		ISP               string    `gorm:"size:150;default:NA" json:"isp"`
		ReferralURL       string    `gorm:"size:255;default:NA" json:"referral_url"`
		UserAgent         string    `gorm:"type:text" json:"user_agent"`
		TrafficSource     bool      `gorm:"not null;default:false" json:"traffic_source"`
		TrafficSourceData string    `gorm:"size:255;default:NA" json:"traffic_source_data"`
		UserRejected      bool      `gorm:"not null;default:false" json:"user_rejected"`
		UniqueClick       bool      `gorm:"not null;default:false" json:"unique_click"`
		UserDuplicated    bool      `gorm:"not null;default:false" json:"user_duplicated"`
		Handset           string    `gorm:"type:text" json:"handset"`
		HandsetCode       string    `gorm:"size:150;default:NA" json:"handset_code"`
		HandsetType       string    `gorm:"size:150;default:NA" json:"handset_type"`
		URLLanding        string    `gorm:"size:255;default:NA" json:"url_landing"`
		URLWarpLanding    string    `gorm:"size:255;default:NA" json:"url_warp_landing"`
		URLService        string    `gorm:"size:255;default:NA" json:"url_service"`
		URLTFCORSmartlink string    `gorm:"size:255;default:NA" json:"url_tfc_or_smartlink"`
		PixelUsedDate     time.Time `gorm:"not null" json:"pixel_used_date"`
		StatusPostback    bool      `gorm:"not null;default:false" json:"status_postback"`
		IsUnique          bool      `gorm:"not null;default:false" json:"is_unique"`
		URLPostback       string    `gorm:"size:255;default:NA" json:"url_postback"`
		StatusURLPostback string    `gorm:"size:150" json:"status_url_postback"`
		ReasonURLPostback string    `gorm:"size:255" json:"reason_url_postback"`
		IsActive          bool      `gorm:"not null;default:false" json:"is_active"`
		CounterMOCapping  int       `gorm:"not null;length:10" json:"counter_mo_capping"`
		MOCapping         int       `gorm:"not null;length:10" json:"mo_capping"`
		StatusCapping     bool      `gorm:"not null;default:false" json:"status_capping"`
		CounterMORatio    int       `gorm:"not null;length:10" json:"counter_mo_ratio"`
		RatioSend         int       `gorm:"not null;length:10;default:1" json:"ratio_send"`
		RatioReceive      int       `gorm:"not null;length:10;default:4" json:"ratio_receive"`
		StatusRatio       bool      `gorm:"not null;default:false" json:"status_ratio"`
		APIURL            string    `gorm:"type:text" json:"api_url"`
		CampaignObjective string    `gorm:"size:50;default:NA" json:"campaign_objective"`
		Channel           string    `gorm:"size:50;default:NA" json:"channel"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	PixelStorage struct {
		gorm.Model
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		CampaignDetailId  int       `gorm:"index:idx_mocampdetailid_unique;not null;length:20" json:"campaign_detail_id"`
		Pxdate            time.Time `gorm:"not null" json:"pxdate"`
		URLServiceKey     string    `gorm:"index:idx_urlservicekey;index:idx_token;index:idx_msisdn;index:idx_pixel,not null;size:50" json:"urlservicekey"`
		CampaignId        string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		CampaignName      string    `gorm:"size:100;deafult:NA" json:"campaign_name"`
		Country           string    `gorm:"not null;size:50" json:"country"`
		Operator          string    `gorm:"not null;size:50" json:"operator"`
		Partner           string    `gorm:"not null;size:50" json:"partner"`
		Aggregator        string    `gorm:"size:50" json:"aggregator"`
		Adnet             string    `gorm:"not null;size:50" json:"adnet"`
		Service           string    `gorm:"not null;size:50" json:"service"`
		Keyword           string    `gorm:"size:50" json:"keyword"`
		Subkeyword        string    `gorm:"size:50" json:"subkeyword"`
		IsBillable        bool      `gorm:"default:false" json:"is_billable"`
		Plan              string    `gorm:"size:50;default:NA" json:"plan"`
		PO                string    `gorm:"size:50;default:NA" json:"po"`
		Cost              string    `gorm:"size:50;default:NA" json:"cost"`
		PubId             string    `gorm:"size:50;default:NA" json:"pubid"`
		ShortCode         string    `gorm:"size:50;default:NA" json:"short_code"`
		URL               string    `gorm:"size:255;default:NA" json:"url"`
		URLType           string    `gorm:"size:50;default:NA" json:"url_type"`
		Pixel             string    `gorm:"index:idx_pixel,size:255;default:NA" json:"pixel"`
		Token             string    `gorm:"index:idx_token,size:255;default:NA" json:"token"`
		TrxId             string    `gorm:"size:255;default:NA" json:"trx_id"`
		TrxDate           time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"trx_date"`
		Msisdn            string    `gorm:"index:idx_msisdn,size:255;default:NA" json:"msisdn"`
		IsUsed            bool      `gorm:"not null;default:false" json:"is_used"`
		Browser           string    `gorm:"size:150;default:NA" json:"browser"`
		OS                string    `gorm:"size:150;default:NA" json:"os"`
		Ip                string    `gorm:"type:text" json:"ip"`
		ISP               string    `gorm:"size:150;default:NA" json:"isp"`
		ReferralURL       string    `gorm:"size:255;default:NA" json:"referral_url"`
		UserAgent         string    `gorm:"type:text" json:"user_agent"`
		TrafficSource     bool      `gorm:"not null;default:false" json:"traffic_source"`
		TrafficSourceData string    `gorm:"size:255;default:NA" json:"traffic_source_data"`
		UserRejected      bool      `gorm:"not null;default:false" json:"user_rejected"`
		UniqueClick       bool      `gorm:"not null;default:false" json:"unique_click"`
		UserDuplicated    bool      `gorm:"not null;default:false" json:"user_duplicated"`
		Handset           string    `gorm:"type:text" json:"handset"`
		HandsetCode       string    `gorm:"size:150;default:NA" json:"handset_code"`
		HandsetType       string    `gorm:"size:150;default:NA" json:"handset_type"`
		URLLanding        string    `gorm:"size:255;default:NA" json:"url_landing"`
		URLWarpLanding    string    `gorm:"size:255;default:NA" json:"url_warp_landing"`
		URLService        string    `gorm:"size:255;default:NA" json:"url_service"`
		URLTFCORSmartlink string    `gorm:"size:255;default:NA" json:"url_tfc_or_smartlink"`
		PixelUsedDate     time.Time `gorm:"not null" json:"pixel_used_date"`
		StatusPostback    bool      `gorm:"not null;default:false" json:"status_postback"`
		IsUnique          bool      `gorm:"not null;default:false" json:"is_unique"`
		URLPostback       string    `gorm:"size:255;default:NA" json:"url_postback"`
		StatusURLPostback string    `gorm:"size:150" json:"status_url_postback"`
		ReasonURLPostback string    `gorm:"size:255" json:"reason_url_postback"`
		IsActive          bool      `gorm:"not null;default:false" json:"is_active"`
		CounterMOCapping  int       `gorm:"length:10;default:0" json:"counter_mo_capping"`
		MOCapping         int       `gorm:"length:10;default:0" json:"mo_capping"`
		StatusCapping     bool      `gorm:"not null;default:false" json:"status_capping"`
		CounterMORatio    int       `gorm:"not null;length:10;default:0" json:"counter_mo_ratio"`
		RatioSend         int       `gorm:"not null;length:10;default:1" json:"ratio_send"`
		RatioReceive      int       `gorm:"not null;length:10;default:4" json:"ratio_receive"`
		StatusRatio       bool      `gorm:"not null;default:false" json:"status_ratio"`
		APIURL            string    `gorm:"type:text" json:"api_url"`
		CampaignObjective string    `gorm:"size:50;default:NA" json:"campaign_objective"`
		Channel           string    `gorm:"size:50;default:NA" json:"channel"`
		GoogleSheet       string    `gorm:"type:text;default:NA" json:"google_sheet"`
		Currency          string    `gorm:"size:10;default:NA" json:"currency"`
		PostbackMethod    string    `gorm:"size:50" json:"postback_method"`
		LandingTime       time.Time `json:"landing_time"`
		LandedTime        float64   `gorm:"type:double precision;default:0" json:"landed_time"`
		HttpStatus        int       `gorm:"not null;length:10;default:0" json:"http_status"`
		IsOperator        bool      `gorm:"not null;default:false" json:"is_operator"`
		MCampaignId       string    `gorm:"size:255;default:NA" json:"m_campaign_id"`
		MAdgroupId        string    `gorm:"size:255;default:NA" json:"m_adgroup_id"`
		MPlacement        string    `gorm:"size:255;default:NA" json:"m_placement"`
		MCreative         string    `gorm:"size:255;default:NA" json:"m_creative"`
		MGadCampaignId    string    `gorm:"size:255;default:NA" json:"m_gad_campaign_id"`
		MStatusCharge     bool      `gorm:"not null;default:false" json:"m_status_charge"`
		MStatusTimeCharge time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"m_status_time_charge"`
		MSnackCampaignId  string    `gorm:"size:255;default:NA" json:"m_snack_campaign_id"`
		IsEvina           bool      `gorm:"default:false" json:"is_evina"`
		StatusEvinaFraud  bool      `gorm:"default:false" json:"status_evina_fraud"`
		EvinaTI           string    `gorm:"type:text" json:"evina_ti"`
		EvinaTRX             string    `gorm:"type:text" json:"evina_trx"`
		EvinaTS              int64     `gorm:"default:0" json:"evina_ts"`
		EvinaASCode          int       `gorm:"default:0" json:"evina_as_code"`
		EvinaFTCode          int       `gorm:"default:0" json:"evina_ft_code"`
		StatusBillable       string    `gorm:"size:150" json:"status_billable"`
		StatusCodeBillable   string    `gorm:"size:150" json:"status_code_billable"`
		ReasonStatusBillable string    `gorm:"size:255" json:"reason_status_billable"`
		CreatedAt            time.Time
		UpdatedAt            time.Time
	}

	ClickStorage struct {
		gorm.Model
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		CampaignDetailId  int       `gorm:"index:idx_mocampdetailid_unique;not null;length:20" json:"campaign_detail_id"`
		Clkdate           time.Time `gorm:"not null" json:"clkdate"`
		URLServiceKey     string    `gorm:"index:idx_urlservicekey;index:idx_token;index:idx_msisdn;index:idx_pixel,not null;size:50" json:"urlservicekey"`
		CampaignId        string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		CampaignName      string    `gorm:"size:100;deafult:NA" json:"campaign_name"`
		Country           string    `gorm:"not null;size:50" json:"country"`
		Operator          string    `gorm:"not null;size:50" json:"operator"`
		Partner           string    `gorm:"not null;size:50" json:"partner"`
		Aggregator        string    `gorm:"size:50" json:"aggregator"`
		Adnet             string    `gorm:"not null;size:50" json:"adnet"`
		Service           string    `gorm:"not null;size:50" json:"service"`
		Keyword           string    `gorm:"size:50" json:"keyword"`
		Subkeyword        string    `gorm:"size:50" json:"subkeyword"`
		IsBillable        bool      `gorm:"default:false" json:"is_billable"`
		Plan              string    `gorm:"size:50;default:NA" json:"plan"`
		PO                string    `gorm:"size:50;default:NA" json:"po"`
		Cost              string    `gorm:"size:50;default:NA" json:"cost"`
		PubId             string    `gorm:"size:50;default:NA" json:"pubid"`
		ShortCode         string    `gorm:"size:50;default:NA" json:"short_code"`
		URL               string    `gorm:"size:255;default:NA" json:"url"`
		URLType           string    `gorm:"size:50;default:NA" json:"url_type"`
		Pixel             string    `gorm:"index:idx_pixel,size:255;default:NA" json:"pixel"`
		Token             string    `gorm:"index:idx_token,size:255;default:NA" json:"token"`
		TrxId             string    `gorm:"size:255;default:NA" json:"trx_id"`
		TrxDate           time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"trx_date"`
		Msisdn            string    `gorm:"index:idx_msisdn,size:255;default:NA" json:"msisdn"`
		IsUsed            bool      `gorm:"not null;default:false" json:"is_used"`
		Browser           string    `gorm:"size:150;default:NA" json:"browser"`
		OS                string    `gorm:"size:150;default:NA" json:"os"`
		Ip                string    `gorm:"type:text" json:"ip"`
		ISP               string    `gorm:"size:150;default:NA" json:"isp"`
		ReferralURL       string    `gorm:"size:255;default:NA" json:"referral_url"`
		UserAgent         string    `gorm:"type:text" json:"user_agent"`
		TrafficSource     bool      `gorm:"not null;default:false" json:"traffic_source"`
		TrafficSourceData string    `gorm:"size:255;default:NA" json:"traffic_source_data"`
		UserRejected      bool      `gorm:"not null;default:false" json:"user_rejected"`
		UniqueClick       bool      `gorm:"not null;default:false" json:"unique_click"`
		UserDuplicated    bool      `gorm:"not null;default:false" json:"user_duplicated"`
		Handset           string    `gorm:"type:text" json:"handset"`
		HandsetCode       string    `gorm:"size:150;default:NA" json:"handset_code"`
		HandsetType       string    `gorm:"size:150;default:NA" json:"handset_type"`
		URLLanding        string    `gorm:"size:255;default:NA" json:"url_landing"`
		URLWarpLanding    string    `gorm:"size:255;default:NA" json:"url_warp_landing"`
		URLService        string    `gorm:"size:255;default:NA" json:"url_service"`
		URLTFCORSmartlink string    `gorm:"size:255;default:NA" json:"url_tfc_or_smartlink"`
		PixelUsedDate     time.Time `gorm:"not null" json:"pixel_used_date"`
		StatusPostback    bool      `gorm:"not null;default:false" json:"status_postback"`
		IsUnique          bool      `gorm:"not null;default:false" json:"is_unique"`
		URLPostback       string    `gorm:"size:255;default:NA" json:"url_postback"`
		StatusURLPostback string    `gorm:"size:150" json:"status_url_postback"`
		ReasonURLPostback string    `gorm:"size:255" json:"reason_url_postback"`
		IsActive          bool      `gorm:"not null;default:false" json:"is_active"`
		CounterMOCapping  int       `gorm:"length:10;default:0" json:"counter_mo_capping"`
		MOCapping         int       `gorm:"length:10;default:0" json:"mo_capping"`
		StatusCapping     bool      `gorm:"not null;default:false" json:"status_capping"`
		CounterMORatio    int       `gorm:"not null;length:10;default:0" json:"counter_mo_ratio"`
		RatioSend         int       `gorm:"not null;length:10;default:1" json:"ratio_send"`
		RatioReceive      int       `gorm:"not null;length:10;default:4" json:"ratio_receive"`
		StatusRatio       bool      `gorm:"not null;default:false" json:"status_ratio"`
		APIURL            string    `gorm:"size:255;default:NA" json:"api_url"`
		CampaignObjective string    `gorm:"size:50;default:NA" json:"campaign_objective"`
		Channel           string    `gorm:"size:50;default:NA" json:"channel"`
		GoogleSheet       string    `gorm:"type:text;default:NA" json:"google_sheet"`
		Currency          string    `gorm:"size:10;default:NA" json:"currency"`
		PostbackMethod    string    `gorm:"size:50" json:"postback_method"`
		LandingTime       time.Time `json:"landing_time"`
		LandedTime        float64   `gorm:"type:double precision;default:0" json:"landed_time"`
		HttpStatus        int       `gorm:"not null;length:10;default:0" json:"http_status"`
		IsOperator        bool      `gorm:"not null;default:false" json:"is_operator"`
		MCampaignId       string    `gorm:"size:255;default:NA" json:"m_campaign_id"`
		MAdgroupId        string    `gorm:"size:255;default:NA" json:"m_adgroup_id"`
		MPlacement        string    `gorm:"size:255;default:NA" json:"m_placement"`
		MCreative         string    `gorm:"size:255;default:NA" json:"m_creative"`
		MGadCampaignId    string    `gorm:"size:255;default:NA" json:"m_gad_campaign_id"`
		MStatusCharge     bool      `gorm:"not null;default:false" json:"m_status_charge"`
		MStatusTimeCharge time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"m_status_time_charge"`
		MSnackCampaignId  string    `gorm:"size:255;default:NA" json:"m_snack_campaign_id"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	Postback struct {
		gorm.Model
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		CampaignDetailId  int       `gorm:"index:idx_mocampdetailid_unique;not null;length:20" json:"campaign_detail_id"`
		Pxdate            time.Time `gorm:"not null" json:"pxdate"`
		URLServiceKey     string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId        string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country           string    `gorm:"not null;size:50" json:"country"`
		Operator          string    `gorm:"not null;size:50" json:"operator"`
		Partner           string    `gorm:"not null;size:50" json:"partner"`
		Aggregator        string    `gorm:"size:50" json:"aggregator"`
		Adnet             string    `gorm:"not null;size:50" json:"adnet"`
		Service           string    `gorm:"not null;size:50" json:"service"`
		Keyword           string    `gorm:"size:50" json:"keyword"`
		Subkeyword        string    `gorm:"size:50" json:"subkeyword"`
		IsBillable        bool      `gorm:"default:false" json:"is_billable"`
		Plan              string    `gorm:"size:50;default:NA" json:"plan"`
		PO                string    `gorm:"size:50;default:NA" json:"po"`
		Cost              string    `gorm:"size:50;default:NA" json:"cost"`
		PubId             string    `gorm:"size:50;default:NA" json:"pubid"`
		ShortCode         string    `gorm:"size:50;default:NA" json:"short_code"`
		URL               string    `gorm:"size:255;default:NA" json:"url"`
		URLType           string    `gorm:"size:50;default:NA" json:"url_type"`
		Pixel             string    `gorm:"size:255;default:NA" json:"pixel"`
		Token             string    `gorm:"size:255;default:NA" json:"token"`
		TrxId             string    `gorm:"size:255;default:NA" json:"trx_id"`
		Msisdn            string    `gorm:"size:255;default:NA" json:"msisdn"`
		IsUsed            bool      `gorm:"not null;default:false" json:"is_used"`
		Browser           string    `gorm:"size:150;default:NA" json:"browser"`
		OS                string    `gorm:"size:150;default:NA" json:"os"`
		Ip                string    `gorm:"size:150;default:NA" json:"ip"`
		ISP               string    `gorm:"size:150;default:NA" json:"isp"`
		ReferralURL       string    `gorm:"size:255;default:NA" json:"referral_url"`
		UserAgent         string    `gorm:"type:text" json:"user_agent"`
		TrafficSource     bool      `gorm:"not null;default:false" json:"traffic_source"`
		TrafficSourceData string    `gorm:"size:255;default:NA" json:"traffic_source_data"`
		UserRejected      bool      `gorm:"not null;default:false" json:"user_rejected"`
		UniqueClick       bool      `gorm:"not null;default:false" json:"unique_click"`
		UserDuplicated    bool      `gorm:"not null;default:false" json:"user_duplicated"`
		Handset           string    `gorm:"type:text" json:"handset"`
		HandsetCode       string    `gorm:"size:150;default:NA" json:"handset_code"`
		HandsetType       string    `gorm:"size:150;default:NA" json:"handset_type"`
		URLLanding        string    `gorm:"size:255;default:NA" json:"url_landing"`
		URLWarpLanding    string    `gorm:"size:255;default:NA" json:"url_warp_landing"`
		URLService        string    `gorm:"size:255;default:NA" json:"url_service"`
		URLTFCORSmartlink string    `gorm:"size:255;default:NA" json:"url_tfc_or_smartlink"`
		PixelUsedDate     time.Time `gorm:"not null" json:"pixel_used_date"`
		StatusPostback    bool      `gorm:"not null;default:false" json:"status_postback"`
		IsUnique          bool      `gorm:"not null;default:false" json:"is_unique"`
		URLPostback       string    `gorm:"size:255;default:NA" json:"url_postback"`
		StatusURLPostback string    `gorm:"size:150" json:"status_url_postback"`
		ReasonURLPostback string    `gorm:"size:255" json:"reason_url_postback"`
		IsActive          bool      `gorm:"not null;default:false" json:"is_active"`
		CounterMOCapping  int       `gorm:"not null;length:10" json:"counter_mo_capping"`
		MOCapping         int       `gorm:"not null;length:10" json:"mo_capping"`
		StatusCapping     bool      `gorm:"not null;default:false" json:"status_capping"`
		CounterMORatio    int       `gorm:"not null;length:10" json:"counter_mo_ratio"`
		RatioSend         int       `gorm:"not null;length:10;default:1" json:"ratio_send"`
		RatioReceive      int       `gorm:"not null;length:10;default:4" json:"ratio_receive"`
		StatusRatio       bool      `gorm:"not null;default:false" json:"status_ratio"`
		APIURL            string    `gorm:"type:text" json:"api_url"`
		CampaignObjective string    `gorm:"not null;size:50" json:"campaign_objective"`
		Channel           string    `gorm:"not null;size:50" json:"channel"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	SummaryDashboard struct {
		gorm.Model
		ID                int       `gorm:"autoIncrement:false" json:"id"`
		SummaryDate       time.Time `gorm:"type:date;primaryKey" json:"summary_date"`
		Adnet             string    `gorm:"length:100;primaryKey" json:"adnet"`
		Company           string    `gorm:"length:200;primaryKey" json:"company"`
		TotalMO           int       `gorm:"length:20;default:0" json:"total_mo"`
		TotalCpaMO        int       `gorm:"length:20;default:0" json:"total_cpa_mo"`
		TotalSmsMO        int       `gorm:"length:20;default:0" json:"total_sms_mo"`
		TotalMainstreamMO int       `gorm:"length:20;default:0" json:"total_mainstream_mo"`

		TotalPostback           int `gorm:"length:20;default:0" json:"total_postback"`
		TotalCPAPostback        int `gorm:"length:20;default:0" json:"total_cpa_postback"`
		TotalSMSPostback        int `gorm:"length:20;default:0" json:"total_sms_postback"`
		TotalMainstreamPostback int `gorm:"length:20;default:0" json:"total_mainstream_postback"`

		TotalSpending           float64 `gorm:"type:double precision;default:0" json:"total_spending"`
		TotalCPASpending        float64 `gorm:"type:double precision;default:0" json:"total_cpa_spending"`
		TotalSMSSpending        float64 `gorm:"type:double precision;default:0" json:"total_sms_spending"`
		TotalMainstreamSpending float64 `gorm:"type:double precision;default:0" json:"total_mainstream_spending"`

		TotalSaaf           float64   `gorm:"type:double precision;default:0" json:"total_saaf"`
		TotalCPASaaf        float64   `gorm:"type:double precision;default:0" json:"total_cpa_saaf"`
		TotalSMSSaaf        float64   `gorm:"type:double precision;default:0" json:"total_sms_saaf"`
		TotalMainstreamSaaf float64   `gorm:"type:double precision;default:0" json:"total_mainstream_saaf"`
		UpdatedAt           time.Time `gorm:"type:timestamp" json:"updated_at"`
	}

	SummaryAPIDashboard struct {
		gorm.Model
		Date    time.Time `gorm:"type:date;primaryKey" json:"date"`
		TotalMO int       `gorm:"length:20;default:0" json:"total_mo"`
		Sbaf    float64   `gorm:"type:double precision;default:0" json:"sbaf"`
	}

	SummaryCampaign struct {
		gorm.Model
		ID                       int       `gorm:"primaryKey;autoIncrement" json:"id"`
		Status                   bool      `gorm:"not null;size:50" json:"status"`
		SummaryDate              time.Time `gorm:"type:date;uniqueIndex:idx_sumunique" json:"summary_date"`
		URLServiceKey            string    `gorm:"uniqueIndex:idx_sumunique;not null;size:50" json:"urlservicekey"`
		CampaignId               string    `gorm:"uniqueIndex:idx_sumunique;not null;size:50" json:"campaign_id"`
		CampaignName             string    `gorm:"not null;size:100" json:"campaign_name"`
		Country                  string    `gorm:"uniqueIndex:idx_sumunique;not null;size:50" json:"country"`
		Operator                 string    `gorm:"uniqueIndex:idx_sumunique;not null;size:50" json:"operator"`
		Partner                  string    `gorm:"uniqueIndex:idx_sumunique;not null;size:50" json:"partner"`
		Aggregator               string    `gorm:"not null;size:50" json:"aggregator"`
		Adnet                    string    `gorm:"uniqueIndex:idx_sumunique;not null;size:50" json:"adnet"`
		Service                  string    `gorm:"uniqueIndex:idx_sumunique;not null;size:50" json:"service"`
		ShortCode                string    `gorm:"not null;size:50" json:"short_code"`
		Traffic                  int       `gorm:"length:20;default:0" json:"traffic"`
		Landing                  int       `gorm:"length:20;default:0" json:"landing"`
		Clicked                  int       `gorm:"length:20;default:0" json:"clicked"`
		MoReceived               int       `gorm:"length:20;default:0" json:"mo_received"`
		CR                       float64   `gorm:"type:double precision;default:0" json:"cr"`
		Postback                 int       `gorm:"length:20;default:0" json:"postback"`
		TotalFP                  int       `gorm:"length:20;default:0" json:"total_fp"`
		SuccessFP                int       `gorm:"length:20;default:0" json:"success_fp"`
		Billrate                 float64   `gorm:"type:double precision;default:0" json:"billrate"`
		ROI                      float64   `gorm:"type:double precision;default:0" json:"roi"`
		PO                       float64   `gorm:"type:double precision;default:0" json:"po"`
		FirstPush                float64   `gorm:"type:double precision;default:0" json:"first_push"`
		Cost                     float64   `gorm:"type:double precision;length:20;default:0" json:"cost"`
		SBAF                     float64   `gorm:"type:double precision;length:20;default:0" json:"sbaf"`
		SAAF                     float64   `gorm:"type:double precision;length:20;default:0" json:"saaf"`
		CPA                      float64   `gorm:"type:double precision;default:0" json:"cpa"`
		Revenue                  float64   `gorm:"type:double precision;length:20;default:0" json:"revenue"`
		URLAfter                 string    `gorm:"size:255;default:NA" json:"url_after"`
		URLBefore                string    `gorm:"size:255;default:NA" json:"url_before"`
		MOLimit                  int       `gorm:"length:10;default:0" json:"mo_limit"`
		RatioSend                int       `gorm:"length:10;default:1" json:"ratio_send"`
		RatioReceive             int       `gorm:"length:10;default:4" json:"ratio_receive"`
		Company                  string    `gorm:"size:255;default:NA" json:"company"`
		ClientType               string    `gorm:"size:30;default:NA" json:"client_type"`
		CostPerConversion        float64   `gorm:"type:double precision;default:0" json:"cost_per_conversion"`
		AgencyFee                float64   `gorm:"type:double precision;default:0" json:"agency_fee"`
		TargetDailyBudget        float64   `gorm:"type:double precision;default:0" json:"target_daily_budget"`
		TargetMonthlyBudget      float64   `gorm:"type:double precision;default:0" json:"target_monthly_budget"`
		CrMO                     float64   `gorm:"type:double precision;default:0" json:"cr_mo"`
		CrPostback               float64   `gorm:"type:double precision;default:0" json:"cr_postback"`
		TotalWakiAgencyFee       float64   `gorm:"type:double precision;default:0" json:"total_waki_agency_fee"`
		BudgetUsage              float64   `gorm:"type:double precision;default:0" json:"budget_usage"`
		TargetDailyBudgetChanges int       `gorm:"length:12;default:0" json:"target_daily_budget_changes"`
		TechnicalFee             float64   `gorm:"type:double precision;default:0" json:"technical_fee"`
		CampaignObjective        string    `gorm:"uniqueIndex:idx_sumunique;size:50;default:NA" json:"campaign_objective"`
		Channel                  string    `gorm:"size:50;default:NA" json:"channel"`
		PricePerMO               float64   `gorm:"type:double precision;default:0" json:"price_per_mo"`
		POAF                     float64   `gorm:"type:double precision;default:0" json:"po_af"`
		CreatedAt                time.Time
		UpdatedAt                time.Time
	}

	SummaryCampaignBilling struct {
		gorm.Model
		ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
		SummaryDate   time.Time `gorm:"type:date;uniqueIndex:idx_sumbilunique" json:"summary_date"`
		URLServiceKey string    `gorm:"uniqueIndex:idx_sumbilunique;not null;size:50" json:"urlservicekey"`
		CampaignName  string    `gorm:"not null;size:100" json:"campaign_name"`
		Country       string    `gorm:"uniqueIndex:idx_sumbilunique;not null;size:50" json:"country"`
		Operator      string    `gorm:"uniqueIndex:idx_sumbilunique;not null;size:50" json:"operator"`
		Partner       string    `gorm:"uniqueIndex:idx_sumbilunique;not null;size:50" json:"partner"`
		Adnet         string    `gorm:"uniqueIndex:idx_sumbilunique;not null;size:50" json:"adnet"`
		Service       string    `gorm:"uniqueIndex:idx_sumbilunique;not null;size:50" json:"service"`
		Company       string    `gorm:"size:255;default:NA" json:"company"`
		CampaignId    string    `gorm:"uniqueIndex:idx_sumbilunique;not null;size:50" json:"campaign_id"`
		AdgroupID     string    `gorm:"uniqueIndex:idx_sumbilunique;not null;size:50" json:"adgroup_id"`
		Placement     string    `gorm:"type:text" json:"placement"`
		StatusSuccess int       `gorm:"length:20;default:0" json:"status_success"`
		StatusFailed  int       `gorm:"length:20;default:0" json:"status_failed"`
		TotalBill     int       `gorm:"length:20;default:0" json:"total_bill"`
		BillRate      float64   `gorm:"type:double precision;default:0" json:"bill_rate"`
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	IncSummaryCampaignHour struct {
		gorm.Model
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		SummaryDateHour   time.Time `gorm:"uniqueIndex:idx_incsumunique" json:"summary_date_hour"`
		URLServiceKey     string    `gorm:"uniqueIndex:idx_incsumunique;not null;size:50" json:"urlservicekey"`
		CampaignId        string    `gorm:"uniqueIndex:idx_incsumunique;not null;size:50" json:"campaign_id"`
		CampaignObjective string    `gorm:"uniqueIndex:idx_incsumunique;size:50;default:NA" json:"campaign_objective"`
		Country           string    `gorm:"uniqueIndex:idx_incsumunique;not null;size:50" json:"country"`
		Operator          string    `gorm:"uniqueIndex:idx_incsumunique;not null;size:50" json:"operator"`
		Partner           string    `gorm:"uniqueIndex:idx_incsumunique;not null;size:50" json:"partner"`
		Aggregator        string    `gorm:"not null;size:50" json:"aggregator"`
		Adnet             string    `gorm:"uniqueIndex:idx_incsumunique;not null;size:50" json:"adnet"`
		Service           string    `gorm:"uniqueIndex:idx_incsumunique;not null;size:50" json:"service"`
		ShortCode         string    `gorm:"not null;size:50" json:"short_code"`
		Clicked           int       `gorm:"length:20;default:0" json:"clicked"`
		Landing           int       `gorm:"length:20;default:0" json:"landing"`
		MoReceived        int       `gorm:"length:20;default:0" json:"mo_received"`
		Postback          int       `gorm:"length:20;default:0" json:"postback"`
		POAF              float64   `gorm:"type:double precision;default:0" json:"po_af"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	DataClicked struct {
		gorm.Model
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		ClickedTime       time.Time `json:"clicked_time"`
		ClickedButtonTime int       `gorm:"not null;length:20" json:"clicked_button_time"`
		HttpStatus        int       `gorm:"not null;length:10" json:"http_status"`
		URLServiceKey     string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId        string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country           string    `gorm:"not null;size:50" json:"country"`
		Operator          string    `gorm:"not null;size:50" json:"operator"`
		Partner           string    `gorm:"not null;size:50" json:"partner"`
		Aggregator        string    `gorm:"size:50" json:"aggregator"`
		Adnet             string    `gorm:"not null;size:50" json:"adnet"`
		Service           string    `gorm:"not null;size:50" json:"service"`
		ShortCode         string    `gorm:"size:50;default:NA" json:"short_code"`
		Keyword           string    `gorm:"size:50" json:"keyword"`
		Subkeyword        string    `gorm:"size:50" json:"subkeyword"`
		IsBillable        bool      `gorm:"not null;default:false" json:"is_billable"`
		Plan              string    `gorm:"size:50;default:NA" json:"plan"`
		OS                string    `gorm:"size:35;default:NA" json:"os"`
		IsOperator        bool      `gorm:"not null;default:false" json:"is_operator"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	DataLanding struct {
		gorm.Model
		ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
		LandingTime   time.Time `json:"landing_time"`
		LandedTime    int       `gorm:"not null;length:20" json:"landed_time"`
		HttpStatus    int       `gorm:"not null;length:10" json:"http_status"`
		URLServiceKey string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId    string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country       string    `gorm:"not null;size:50" json:"country"`
		Operator      string    `gorm:"not null;size:50" json:"operator"`
		Partner       string    `gorm:"not null;size:50" json:"partner"`
		Aggregator    string    `gorm:"size:50" json:"aggregator"`
		Adnet         string    `gorm:"not null;size:50" json:"adnet"`
		Service       string    `gorm:"not null;size:50" json:"service"`
		ShortCode     string    `gorm:"size:50;default:NA" json:"short_code"`
		Keyword       string    `gorm:"size:50" json:"keyword"`
		Subkeyword    string    `gorm:"size:50" json:"subkeyword"`
		IsBillable    bool      `gorm:"not null;default:false" json:"is_billable"`
		Plan          string    `gorm:"size:50;default:NA" json:"plan"`
		OS            string    `gorm:"size:35;default:NA" json:"os"`
		IsOperator    bool      `gorm:"not null;default:false" json:"is_operator"`
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	DataRedirect struct {
		gorm.Model
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		RedirectTime      time.Time `json:"redirect_time"`
		RedirectAddedTime int       `gorm:"not null;length:20" json:"redirect_added_time"`
		HttpStatus        int       `gorm:"not null;length:10" json:"http_status"`
		URLServiceKey     string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId        string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country           string    `gorm:"not null;size:50" json:"country"`
		Operator          string    `gorm:"not null;size:50" json:"operator"`
		Partner           string    `gorm:"not null;size:50" json:"partner"`
		Aggregator        string    `gorm:"size:50" json:"aggregator"`
		Adnet             string    `gorm:"not null;size:50" json:"adnet"`
		Service           string    `gorm:"not null;size:50" json:"service"`
		ShortCode         string    `gorm:"size:50;default:NA" json:"short_code"`
		Keyword           string    `gorm:"size:50" json:"keyword"`
		Subkeyword        string    `gorm:"size:50" json:"subkeyword"`
		IsBillable        bool      `gorm:"not null;default:false" json:"is_billable"`
		Plan              string    `gorm:"size:50;default:NA" json:"plan"`
		OS                string    `gorm:"size:35;default:NA" json:"os"`
		IsOperator        bool      `gorm:"not null;default:false" json:"is_operator"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	DataTraffic struct {
		gorm.Model
		ID               int       `gorm:"primaryKey;autoIncrement" json:"id"`
		TrafficTime      time.Time `json:"traffic_time"`
		TrafficAddedTime int       `gorm:"not null;length:20" json:"traffic_added_time"`
		HttpStatus       int       `gorm:"not null;length:10" json:"http_status"`
		URLServiceKey    string    `gorm:"index:idx_urlservicekey;not null;size:50" json:"urlservicekey"`
		CampaignId       string    `gorm:"index:idx_campdetailid_unique;not null;size:50" json:"campaign_id"`
		Country          string    `gorm:"not null;size:50" json:"country"`
		Operator         string    `gorm:"not null;size:50" json:"operator"`
		Partner          string    `gorm:"not null;size:50" json:"partner"`
		Aggregator       string    `gorm:"size:50" json:"aggregator"`
		Adnet            string    `gorm:"not null;size:50" json:"adnet"`
		Service          string    `gorm:"not null;size:50" json:"service"`
		ShortCode        string    `gorm:"size:50;default:NA" json:"short_code"`
		Keyword          string    `gorm:"size:50" json:"keyword"`
		Subkeyword       string    `gorm:"size:50" json:"subkeyword"`
		IsBillable       bool      `gorm:"not null;default:false" json:"is_billable"`
		Plan             string    `gorm:"size:50;default:NA" json:"plan"`
		OS               string    `gorm:"size:35;default:NA" json:"os"`
		IsOperator       bool      `gorm:"not null;default:false" json:"is_operator"`
		CreatedAt        time.Time
		UpdatedAt        time.Time
	}

	ApiPinReport struct {
		gorm.Model
		ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
		DateSend      time.Time `gorm:"type:date;uniqueIndex:idx_pin_unique" json:"date_send"`
		Country       string    `gorm:"not null;size:50;uniqueIndex:idx_pin_unique" json:"country"`
		Company       string    `gorm:"not null;size:50" json:"company"`
		Adnet         string    `gorm:"not null;size:50;uniqueIndex:idx_pin_unique" json:"adnet"`
		Operator      string    `gorm:"not null;size:50;uniqueIndex:idx_pin_unique" json:"operator"`
		Service       string    `gorm:"not null;size:50;uniqueIndex:idx_pin_unique" json:"service"`
		PayoutAdn     float64   `gorm:"type:double precision;default:0" json:"payout_adn"`
		PayoutAF      float64   `gorm:"type:double precision;default:0" json:"payout_af"`
		TotalMO       int       `gorm:"default:0" json:"total_mo"`
		TotalPostback int       `gorm:"default:0" json:"total_postback"`
		SBAF          float64   `gorm:"type:double precision;default:0" json:"sbaf"`
		SAAF          float64   `gorm:"type:double precision;default:0" json:"saaf"`
		PricePerMO    float64   `gorm:"type:double precision;default:0" json:"price_per_mo"`
		WakiRevenue   float64   `gorm:"type:double precision;default:0" json:"waki_revenue"`
		CampaignId    string    `gorm:"size:50;uniqueIndex:idx_pin_unique" json:"campaign_id"`
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	ApiPinPerformance struct {
		gorm.Model
		ID                     int       `gorm:"primaryKey;autoIncrement" json:"id"`
		DateSend               time.Time `gorm:"type:date;uniqueIndex:idx_pin_performance_unique" json:"date_send"`
		Country                string    `gorm:"not null;size:50;uniqueIndex:idx_pin_performance_unique" json:"country"`
		Company                string    `gorm:"not null;size:50" json:"company"`
		Adnet                  string    `gorm:"not null;size:50;uniqueIndex:idx_pin_performance_unique" json:"adnet"`
		Operator               string    `gorm:"not null;size:50;uniqueIndex:idx_pin_performance_unique" json:"operator"`
		Service                string    `gorm:"not null;size:50;uniqueIndex:idx_pin_performance_unique" json:"service"`
		PinRequest             int       `gorm:"length:20;default:0" json:"pin_request"`
		UniquePinRequest       int       `gorm:"length:20;default:0" json:"unique_pin_request"`
		PinSuccess             int       `gorm:"length:20;default:0" json:"pin_success"`
		PinFailed              int       `gorm:"length:20;default:0" json:"pin_failed"`
		PinVerifyRequest       int       `gorm:"length:20;default:0" json:"pin_verify_request"`
		PinVerifyRequestUnique int       `gorm:"length:20;default:0" json:"pin_verify_request_unique"`
		PinOK                  int       `gorm:"length:20;default:0" json:"pin_ok"`
		PinNotOK               int       `gorm:"length:20;default:0" json:"pin_not_ok"`
		PinOkRatio             int       `gorm:"length:20;default:0" json:"pin_ok_ratio"`
		CPA                    float64   `gorm:"type:double precision;not null;length:20;default:0" json:"cpa"`
		CPAWaki                float64   `gorm:"type:double precision;not null;length:20;default:0" json:"cpa_waki"`
		EstimatedARPU          float64   `gorm:"type:double precision;not null;length:20;default:0" json:"estimated_arpu"`
		SBAF                   float64   `gorm:"type:double precision;not null;length:20;default:0" json:"sbaf"`
		SAAF                   float64   `gorm:"type:double precision;not null;length:20;default:0" json:"saaf"`
		ChargedMO              float64   `gorm:"type:double precision;not null;length:20;default:0" json:"charged_mo"`
		SubsCR                 float64   `gorm:"type:double precision;not null;length:20;default:0" json:"subs_cr"`
		AdnetCR                float64   `gorm:"type:double precision;not null;length:20;default:0" json:"adnet_cr"`
		CAC                    float64   `gorm:"type:double precision;not null;length:20;default:0" json:"cac"`
		PaidCAC                float64   `gorm:"type:double precision;not null;length:20;default:0" json:"paid_cac"`
		TotalSpending          float64   `gorm:"type:double precision;not null;length:20;default:0" json:"total_spending"`
		TotalWakiAgencyFee     float64   `gorm:"type:double precision;not null;length:20;default:0" json:"total_waki_agency_fee"`
		CpaPerPO               float64   `gorm:"type:double precision;not null;length:20;default:0" json:"cpa_per_po"`
		TotalSpendingAfterWaki float64   `gorm:"type:double precision;not null;length:20;default:0" json:"total_spending_after_waki"`
		CreatedAt              time.Time
		UpdatedAt              time.Time
	}

	Menu struct {
		ID            int    `gorm:"primaryKey;autoIncrement" json:"id"`
		Name          string `gorm:"type:varchar(255);not null" json:"name"`
		Code          string `gorm:"type:varchar(255);not null" json:"code"`
		Url           string `gorm:"type:varchar(255);not null" json:"url"`
		Icon          string `gorm:"type:varchar(255)" json:"icon"`
		Parent        int    `gorm:"type:int" json:"parent"`
		Sort          string `gorm:"type:varchar(255)" json:"sort"`
		ShowOnSidebar bool   `gorm:"type:bool;default:false" json:"show_on_sidebar"`
		Permission    string `gorm:"type:varchar(255)" json:"permission"`
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	Role struct {
		ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`
		Name      string `gorm:"type:varchar(255);not null" json:"name"`
		Code      string `gorm:"type:varchar(255);not null" json:"code"`
		CreatedAt time.Time
		UpdatedAt time.Time
		Users     []User `gorm:"foreignKey:RoleID;references:ID"`
	}

	Permission struct {
		ID        int  `gorm:"primaryKey;autoIncrement" json:"id"`
		RoleID    int  `gorm:"type:int" json:"role_id"`
		MenuID    int  `gorm:"type:int" json:"menu_id"`
		Status    bool `gorm:"type:bool;" json:"status"`
		CreatedAt time.Time
		UpdatedAt time.Time

		// Relations
		Role Role `gorm:"foreignKey:RoleID;references:ID"`
		Menu Menu `gorm:"foreignKey:MenuID;references:ID"`
	}

	User struct {
		ID              int            `gorm:"primaryKey;autoIncrement" json:"id"`
		Name            string         `gorm:"type:varchar(255);not null" json:"name"`
		Username        string         `gorm:"type:varchar(255);unique;not null" json:"username"`
		Email           string         `gorm:"type:varchar(255);unique;not null" json:"email"`
		EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty"`
		Password        string         `gorm:"type:varchar(255);not null" json:"-"`
		RoleID          int            `json:"role_id"`
		LastLogin       *time.Time     `json:"last_login,omitempty"`
		IPAddress       string         `gorm:"type:varchar(255)" json:"ip_address,omitempty"`
		Handset         string         `gorm:"type:varchar(255)" json:"handset,omitempty"`
		Status          bool           `gorm:"default:true" json:"status"`
		TypeUser        string         `gorm:"type:varchar(255)" json:"type_user,omitempty"`
		IsVerify        bool           `gorm:"type:bool" json:"is_verify,omitempty"`
		VerifyBy        string         `gorm:"type:varchar(255)" json:"verify_by,omitempty"`
		VerifyAt        *time.Time     `json:"verify_at,omitempty"`
		CreatedBy       string         `gorm:"type:varchar(255)" json:"created_by,omitempty"`
		UpdatedBy       string         `gorm:"type:varchar(255)" json:"updated_by,omitempty"`
		DeletedAt       gorm.DeletedAt `gorm:"index"`
		ProfilePic      string         `gorm:"type:varchar(255);default:null" json:"profile_pic,omitempty"`
		CreatedAt       time.Time
		UpdatedAt       time.Time

		// Relations
		Role        Role          `gorm:"foreignKey:RoleID;references:ID"`
		DetailUser  []DetailUser  `gorm:"foreignKey:UserID;references:ID"`
		UserAdnet   []UserAdnet   `gorm:"foreignKey:UserID;references:ID"`
		UserCompany []UserCompany `gorm:"foreignKey:UserID;references:ID"`
	}

	CcEmail struct {
		ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`
		UserID    int    `gorm:"not null" json:"user_id"`
		Email     string `gorm:"type:varchar(255);not null" json:"email"`
		CreatedAt time.Time
		UpdatedAt time.Time

		// Relations
		// User User `gorm:"foreignKey:UserID;references:ID"`
	}

	Email struct {
		ID              int    `gorm:"primaryKey;autoIncrement" json:"id"`
		EmailName       string `gorm:"type:varchar(255)" json:"email_name"`
		EmailPurpose    string `gorm:"type:varchar(255)" json:"email_purpose"`
		EmailContentSub string `gorm:"type:text" json:"email_content_sub"`
		EmailContent    string `gorm:"type:text" json:"email_content"`
		EmailSubject    string `gorm:"type:text" json:"email_subject"`
		EmailBody       string `gorm:"type:text" json:"email_body"`
		EmailFooter     string `gorm:"type:text" json:"email_footer"`
		CreatedAt       time.Time
		UpdatedAt       time.Time
	}

	DetailUser struct {
		ID         int  `gorm:"primaryKey" json:"id"`
		UserID     int  `json:"user_id"`
		CountryID  int  `json:"country_id"`
		OperatorID int  `json:"operator_id"`
		ServiceID  int  `json:"service_id"`
		Status     bool `json:"status"`
		CreatedAt  time.Time
		UpdatedAt  time.Time

		// Relations
		User User `gorm:"foreignKey:UserID;references:ID"`
	}

	UserAdnet struct {
		ID        int  `gorm:"primaryKey" json:"id"`
		UserID    int  `json:"user_id"`
		AdnetID   int  `json:"adnet_id"`
		Status    bool `json:"status"`
		CreatedAt time.Time
		UpdatedAt time.Time

		// Relations
		User  User      `gorm:"foreignKey:UserID;references:ID"`
		Adnet AdnetList `gorm:"foreignKey:AdnetID;references:ID"`
	}

	Country struct {
		ID                uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		Code              string `gorm:"type:varchar(10)" json:"code" validate:"unique,required,max=10"`
		Code2             string `gorm:"type:varchar(10)" json:"code2" validate:"unique,required,max=10"`
		Name              string `gorm:"type:varchar(80)" json:"name"  validate:"unique,required,max=80"`
		NumericCode       string `gorm:"type:varchar(10)" json:"numeric_code" form:"numeric_code"  validate:"required,max=10"`
		MobileCountryCode string `gorm:"type:varchar(10)" json:"mobile_country_code" form:"mobile_country_code"  validate:"required,max=10"`
		Continent         string `gorm:"type:varchar(80)" json:"continent"  validate:"max=80"`
		IsActive          string `gorm:"type:bool;default:true" json:"is_active" form:"is_active" `
	}

	Company struct {
		ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		Name      string `gorm:"type:text" json:"name"`
		LegalName string `gorm:"type:text" json:"legal_name"`
		Address   string `gorm:"type:text" json:"address"`
		Email     string `gorm:"type:varchar(255);default:NA" json:"email"`
		Phone     string `gorm:"type:varchar(20);default:NA" json:"phone"`
	}

	CompanyGroup struct {
		ID      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		Name    string `gorm:"type:text" json:"name"`
		Company string `gorm:"size:255;default:NA" json:"company"`
		Country string `gorm:"size:255;default:NA" json:"country"`
		Partner string `gorm:"size:255;default:NA" json:"partner"`
	}

	Domain struct {
		ID  uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		Url string `gorm:"type:varchar(80)" json:"url"`
	}

	Operator struct {
		ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
		Country    string         `gorm:"type:varchar(10)" json:"country"`
		Name       string         `gorm:"type:varchar(30)" json:"name"`
		IsActive   string         `gorm:"type:bool;default:true" json:"is_active" form:"is_active" `
		Lastupdate time.Time      `gorm:"type:timestamp" json:"lastupdate"`
		Prefix     pq.StringArray `gorm:"type:text[]" json:"prefix"`
	}

	Partner struct {
		ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
		Country        string    `gorm:"type:varchar(10)" json:"country"`
		Name           string    `gorm:"type:varchar(75)" json:"name"`
		Operator       string    `gorm:"type:varchar(50)" json:"operator"`
		Service        string    `gorm:"type:varchar(50)" json:"service"`
		Client         string    `gorm:"type:varchar(50)" json:"client"`
		ClientType     string    `gorm:"type:varchar(50)" json:"client_type"`
		Company        string    `gorm:"type:varchar(75)" json:"company"`
		UrlPostback    string    `gorm:"type:text" json:"url_postback"`
		Postback       string    `gorm:"type:varchar(50)" json:"postback"`
		PostbackMethod string    `gorm:"type:varchar(50)" json:"postback_method"`
		Aggregator     string    `gorm:"type:varchar(50)" json:"aggregator"`
		IsActive       bool      `gorm:"type:bool;default:true" json:"is_active" form:"is_active"`
		IsBillable     bool      `gorm:"type:bool;default:false" json:"is_billable" form:"is_billable"`
		Lastupdate     time.Time `gorm:"type:timestamp" json:"lastupdate"`
	}

	Service struct {
		ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
		Service   string  `gorm:"type:varchar(55)" json:"service"`
		Adn       string  `gorm:"type:varchar(20)" json:"adn"`
		Partner   string  `gorm:"type:varchar(50)" json:"partner"`
		Country   string  `gorm:"type:varchar(50)" json:"country"`
		Operator  string  `gorm:"type:varchar(50)" json:"operator"`
		Currency  string  `gorm:"type:varchar(10)" json:"currency"`
		Price     float64 `gorm:"type:double precision;default:0" json:"price"`
		PortalURL string  `gorm:"type:text" json:"portal_url"`
	}

	AdnetList struct {
		ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
		Code         string         `gorm:"type:varchar(30)" json:"code"`
		Name         string         `gorm:"type:varchar(30)" json:"name"`
		ApiUrl       string         `gorm:"type:text" json:"api_url"`
		ApiUrlBefore string         `gorm:"column:api_url" json:"api_url_before"`
		ApiUrlDr     string         `gorm:"type:text" json:"api_url_dr"`
		IsActive     string         `gorm:"type:bool;default:true" json:"is_active" form:"is_active"`
		IsRemnant    string         `gorm:"type:bool;default:false" json:"is_remnant" form:"is_remnant"`
		Lastupdate   time.Time      `gorm:"type:timestamp" json:"lastupdate"`
		IsDummyPixel string         `gorm:"type:bool;default:false" json:"is_dummy_pixel"`
		IsDsp        string         `gorm:"type:bool;default:false" json:"is_dsp"`
		CcEmail      pq.StringArray `gorm:"type:text[]" json:"cc_email"`
	}

	Audit struct {
		ID            uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		UserType      string `gorm:"type:varchar(255)" json:"user_type"`
		UserID        int    `gorm:"type:int8" json:"user_id"`
		Event         string `gorm:"type:varchar(255)" json:"event"`
		AuditableType string `gorm:"type:varchar(255)" json:"auditable_type"`
		AuditableID   string `gorm:"type:varchar(255)" json:"auditable_id"`
		OldValues     string `gorm:"type:text" json:"old_values"`
		NewValues     string `gorm:"type:text" json:"new_values"`
		URL           string `gorm:"type:text" json:"url"`
		IPAddress     string `gorm:"type:inet" json:"ip_address"`
		UserAgent     string `gorm:"type:varchar(1023)" json:"user_agent"`
		Tags          string `gorm:"type:varchar(255)" json:"tags"`
		CreatedAt     time.Time
		UpdatedAt     time.Time
		ActionName    string `gorm:"type:varchar(255)" json:"action_name"`
	}

	Agency struct {
		ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		Code   string `gorm:"type:varchar(20)" json:"code"`
		Name   string `gorm:"type:varchar(80)" json:"name"`
		APIURL string `gorm:"size:255;default:NA" json:"api_url"`
	}

	Channel struct {
		ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		Name   string `gorm:"type:varchar(80)" json:"name" `
		ApiKey string `gorm:"type:text" json:"api_key"`
		Type   string `gorm:"type:varchar(80)" json:"type" `
	}

	MainstreamGroup struct {
		ID               uint    `gorm:"primaryKey;autoIncrement" json:"id"`
		Name             string  `gorm:"type:varchar(254)" json:"name" `
		Channel          string  `gorm:"type:varchar(80)" json:"channel"  `
		Agency           string  `gorm:"type:varchar(80)" json:"agency" `
		Country          string  `gorm:"type:varchar(10)" json:"country" `
		Operator         string  `gorm:"type:varchar(50)" json:"operator" `
		Service          string  `gorm:"type:varchar(80)" json:"service" `
		Company          string  `gorm:"type:varchar(80)" json:"company" `
		CompanyLegalName string  `gorm:"type:varchar(80)" json:"company_legal_name" `
		CompanyAddress   string  `gorm:"type:varchar(255)" json:"company_address" `
		CompanyEmail     string  `gorm:"type:varchar(255)" json:"company_email"`
		CompanyPhone     string  `gorm:"type:varchar(20)" json:"company_phone"`
		ServiceCurrency  string  `gorm:"type:varchar(10)" json:"service_currency"`
		ServicePrice     float64 `gorm:"type:double precision;default:0" json:"service_price"`
		UniqueDomain     string  `gorm:"type:varchar(80)" json:"unique_domain"`
		DomainService    string  `gorm:"type:varchar(80)" json:"domain_service"`
		PortalURL        string  `gorm:"type:text" json:"portal_url"`
		APIURL           string  `gorm:"type:text" json:"api_url"`
	}

	SummaryLanding struct {
		gorm.Model
		ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
		SummaryDateHour time.Time `gorm:"not null;size:50;uniqueIndex:idx_summary_unique;index:idx_totloadtime_check" json:"summary_date_hour"`
		URLServiceKey   string    `gorm:"not null;size:50;uniqueIndex:idx_summary_unique;index:idx_totloadtime_check" json:"urlservicekey"`
		CampaignId      string    `gorm:"not null;size:50;uniqueIndex:idx_summary_unique" json:"campaign_id"`
		CampaignName    string    `gorm:"not null;size:100" json:"campaign_name"`
		Country         string    `gorm:"not null;size:50;uniqueIndex:idx_summary_unique" json:"country"`
		Operator        string    `gorm:"not null;size:50;uniqueIndex:idx_summary_unique" json:"operator"`
		Partner         string    `gorm:"not null;size:50;uniqueIndex:idx_summary_unique" json:"partner"`
		Aggregator      string    `gorm:"size:50;uniqueIndex:idx_summary_unique" json:"aggregator"`
		Adnet           string    `gorm:"not null;size:50;uniqueIndex:idx_summary_unique" json:"adnet"`
		Service         string    `gorm:"not null;size:50;uniqueIndex:idx_summary_unique" json:"service"`

		URLCampaign      string  `gorm:"not null;size:255" json:"url_campaign"`
		ResponseTime     float64 `gorm:"type:double precision;default:0" json:"response_time"`
		TotalLoadTime    float64 `gorm:"type:double precision;default:0" json:"total_load_time"`
		Landing          int     `gorm:"default:0" json:"landing"`
		SuccessRate      float64 `gorm:"default:0" json:"success_rate"`
		ClickIOS         int     `gorm:"default:0" json:"click_ios"`
		ClickAndroid     int     `gorm:"default:0" json:"click_android"`
		ClickOperator    int     `gorm:"default:0" json:"click_operator"`
		ClickNonOperator int     `gorm:"default:0" json:"click_non_operator"`

		URLService             string  `gorm:"not null;size:255;default:NA" json:"url_service"`
		ResponseUrlServiceTime float64 `gorm:"type:double precision;default:0" json:"response_url_service_time"`

		CreatedAt time.Time
		UpdatedAt time.Time
	}

	IPRangeCsvRow struct {
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		IPType            string    `gorm:"type:varchar(50);uniqueIndex:idx_ip_range_csv_row_unique" json:"ip_type"`
		UploadDate        time.Time `gorm:"not null;uniqueIndex:idx_ip_range_csv_row_unique" json:"upload_date"`
		Network           string    `gorm:"type:varchar(50);uniqueIndex:idx_ip_range_csv_row_unique" json:"network"`
		ISP               string    `gorm:"type:varchar(50);uniqueIndex:idx_ip_range_csv_row_unique" json:"isp"`
		MobileCountryCode string    `gorm:"type:varchar(10)" json:"mobile_country_code"`
	}

	IPRange struct {
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		IPType            string    `gorm:"type:varchar(50);uniqueIndex:idx_ip_range_unique" json:"ip_type"`
		UploadDate        time.Time `gorm:"not null;size:50;uniqueIndex:idx_ip_range_unique" json:"upload_date"`
		Network           string    `gorm:"type:varchar(50);uniqueIndex:idx_ip_range_unique" json:"network"`
		ISP               string    `gorm:"type:varchar(50);uniqueIndex:idx_ip_range_unique" json:"isp"`
		MobileCountryCode string    `gorm:"type:varchar(10)" json:"mobile_country_code"`
		Country           string    `gorm:"type:varchar(10)" json:"country"`
	}

	SummaryTraffic struct {
		gorm.Model
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		Status            bool      `gorm:"not null;size:50" json:"status"`
		SummaryDateHour   time.Time `gorm:"uniqueIndex:idx_sumunique" json:"summary_date_hour"`
		URLServiceKey     string    `gorm:"uniqueIndex:idx_sumunique;not null;size:50" json:"urlservicekey"`
		CampaignId        string    `gorm:"uniqueIndex:idx_sumunique;not null;size:50" json:"campaign_id"`
		CampaignName      string    `gorm:"not null;size:100" json:"campaign_name"`
		Country           string    `gorm:"uniqueIndex:idx_sumunique;not null;size:50" json:"country"`
		Operator          string    `gorm:"uniqueIndex:idx_sumunique;not null;size:50" json:"operator"`
		Partner           string    `gorm:"uniqueIndex:idx_sumunique;not null;size:50" json:"partner"`
		Aggregator        string    `gorm:"not null;size:50" json:"aggregator"`
		Adnet             string    `gorm:"uniqueIndex:idx_sumunique;not null;size:50" json:"adnet"`
		Service           string    `gorm:"uniqueIndex:idx_sumunique;not null;size:50" json:"service"`
		ShortCode         string    `gorm:"not null;size:50" json:"short_code"`
		Landing           int       `gorm:"length:20;default:0" json:"landing"`
		MoReceived        int       `gorm:"length:20;default:0" json:"mo_received"`
		Postback          int       `gorm:"length:20;default:0" json:"postback"`
		FirstPush         float64   `gorm:"type:double precision;default:0" json:"first_push"`
		URLWarpLanding    string    `gorm:"size:255;default:NA" json:"url_landing"`
		URLLanding        string    `gorm:"size:255;default:NA" json:"url_warp_landing"`
		Company           string    `gorm:"size:255;default:NA" json:"company"`
		ClientType        string    `gorm:"size:30;default:NA" json:"client_type"`
		CrMO              float64   `gorm:"type:double precision;default:0" json:"cr_mo"`
		CrPostback        float64   `gorm:"type:double precision;default:0" json:"cr_postback"`
		CampaignObjective string    `gorm:"uniqueIndex:idx_sumunique;size:50;default:NA" json:"campaign_objective"`
		Channel           string    `gorm:"size:50;default:NA" json:"channel"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}

	Continent struct {
		ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		Name string `gorm:"type:varchar(80)" json:"name" `
	}

	BudgetIO struct {
		gorm.Model
		ID     int    `gorm:"primaryKey;autoIncrement" json:"id"`
		Status string `gorm:"default:'submitted'" json:"status"`

		VerifiedAt time.Time `gorm:"default:NULL" json:"verified_at"`

		SubmittedBy string `gorm:"size:255;default:NA" json:"submitted_by"`
		ApprovedBy  string `gorm:"size:255;default:NA" json:"approved_by"`

		IOID             string `gorm:"size:255;default:NA" json:"io_id"`
		CampaignType     string `gorm:"size:100;default:NA" json:"campaign_type"`
		Month            string `gorm:"size:20;default:NA" json:"month"`
		Country          string `gorm:"size:10;default:NA" json:"country"`
		CountryName      string `gorm:"size:255;default:NA" json:"country_name"`
		Continent        string `gorm:"size:50;default:NA" json:"continent"`
		CompanyGroupName string `gorm:"size:255;default:NA" json:"company_group_name"`
		Company          string `gorm:"size:255;default:NA" json:"company"`
		Partner          string `gorm:"size:255;default:NA" json:"partner"`
		Service          string `gorm:"size:255;default:NA" json:"service"`

		TargetCAC          float64 `gorm:"type:double precision" json:"target_cac"`
		TargetROI          int     `gorm:"type:int" json:"target_roi"`
		MonthlyMOTarget    float64 `gorm:"type:double precision" json:"monthly_mo_target"`
		MonthlySpendTarget float64 `gorm:"type:double precision" json:"monthly_spend_target"`

		CPName            string `gorm:"size:255;default:NA" json:"cp_name"`
		PICName           string `gorm:"size:255;default:NA" json:"pic_name"`
		ContactEmail      string `gorm:"size:255;default:NA" json:"contact_email"`
		CPBusinessPICName string `gorm:"size:255;default:NA" json:"business_pic_name"`
		Signature         string `gorm:"type:text" json:"signature"`

		CreatedAt time.Time
		UpdatedAt time.Time
	}

	SummaryBudgetIO struct {
		gorm.Model
		ID                      int     `gorm:"primaryKey;autoIncrement" json:"id"`
		CampaignType            string  `gorm:"size:100;default:NA" json:"campaign_type"`
		Month                   string  `gorm:"size:20;default:NA" json:"month"` // format YYYY-MM
		Country                 string  `gorm:"size:50;default:NA" json:"country"`
		Continent               string  `gorm:"size:50;default:NA" json:"continent"`
		Company                 string  `gorm:"size:50;default:NA" json:"company"`
		Partner                 string  `gorm:"size:50;default:NA" json:"partner"`
		Service                 string  `gorm:"size:50;default:NA" json:"service"`
		TotalMonthlySpendTarget float64 `gorm:"type:double precision;default:0" json:"total_monthly_spend_target"`
		ActualWeek1             float64 `gorm:"type:double precision;default:0" json:"actual_week_1"`
		ActualWeek2             float64 `gorm:"type:double precision;default:0" json:"actual_week_2"`
		ActualWeek3             float64 `gorm:"type:double precision;default:0" json:"actual_week_3"`
		ActualWeek4             float64 `gorm:"type:double precision;default:0" json:"actual_week_4"`

		GMV  float64 `gorm:"type:double precision;default:0" json:"gmv"`
		LTV  float64 `gorm:"type:double precision;default:0" json:"ltv"`
		ROAS float64 `gorm:"type:double precision;default:0" json:"roas"`
		ROI  float64 `gorm:"type:double precision;default:0" json:"roi"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	UserCompany struct {
		gorm.Model
		ID        int  `gorm:"primaryKey;autoIncrement" json:"id"`
		UserID    int  `gorm:"length:20" json:"user_id"`
		CompanyID int  `gorm:"length:20" json:"company_id"`
		Status    bool `gorm:"not null" json:"status"`
		CreatedAt time.Time
		UpdatedAt time.Time

		// Relations
		User    User    `gorm:"foreignKey:UserID;references:ID"`
		Company Company `gorm:"foreignKey:CompanyID;references:ID"`
	}

	DomainService struct {
		ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		Name string `gorm:"type:varchar(80)" json:"name"`
	}

	OperatorAlias struct {
		ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
		Country  string `gorm:"type:varchar(10)" json:"country"`
		Operator string `gorm:"type:varchar(100)" json:"operator"`
		Service  string `gorm:"type:varchar(100)" json:"service"`
		Alias    string `gorm:"type:varchar(100)" json:"alias"`
		Type     string `gorm:"type:varchar(50)" json:"type"`
	}
)

// Hook table campaign_details
/* func (o *CampaignDetail) AfterUpdate(tx *gorm.DB) (err error) {
	if o.CounterMOCapping >= o.MOCapping {
		tx.Model(&CampaignDetail{}).Where("id = ?", o.ID).Update("last_update_capping", o.LastUpdate)
	}
	return
} */
