package entity

type (
	//"[{\"id\":1,\"urlservicekey\":\"idtelgempastelmbv\",\"campaign_id\":\"ID01\",\"name\":\"ID Pass Telesat Gazy\",\"objective\":\"cpa\",\"country\":\"ID\",\"advertiser\":\"PT\",\"operator\":\"telkomsel\",\"partner\":\"pass\",\"aggregator\":\"telesat\",\"adnet\":\"mbv\",\"service\":\"gazy\",\"keyword\":\"gazy\",\"subkeyword\":\"\",\"is_billable\":false,\"plan\":\"\",\"short_code\":\"1234\",\"device_type\":\"all\",\"os\":\"linux\",\"url_type\":\"click2sms\",\"click_type\":2,\"click_delay\":0,\"client_type\":\"internal\",\"traffic_source\":false,\"unique_click\":true,\"url_banner\":\"http://\",\"url_landing\":\"http://\",\"url_warp_landing\":\"http://\",\"url_service\":\"http://\",\"url_tfc_or_smartlink\":\"http://\",\"custom_integration\":\"\",\"ip_address\":\"\",\"is_active\":true,\"mo_capping\":0,\"counter_mo_capping\":0,\"status_capping\":false,\"kpi_upper_limit_capping\":500,\"is_machine_learning_capping\":false,\"ratio_send\":1,\"ratio_receive\":1,\"status_ratio\":true,\"kpi_upper_limit_ratio_send\":1,\"kpi_upper_limit_ratio_receive\":2,\"is_machine_learning_ratio\":false,\"api_url\":\"http://\",\"last_update\":\"\",\"last_update_capping\":\"\"}]"

	GlobalSetting struct {
		CostPerConversion string `json:"cost_per_conversion"`
		AgencyFee         string `json:"agency_fee"`
		TargetDailyBudget string `json:"target_daily_budget"`
		TechnicalFee      string `json:"technical_fee"`
	}

	DataSingleURL struct {
		SingleURLServiceKey string                `json:"single_urlservicekey"`
		Country             string                `json:"country"`
		Objective           string                `json:"objective"`
		MCC                 string                `json:"mcc"`
		Data                []URLServiceKeyChilds `json:"data"`
	}

	URLServiceKeyChilds struct {
		URLServiceKey string `json:"urlservicekey"`
		Operator      string `json:"operator"`
	}

	DataCampaignAction struct {
		Action       string       `json:"action"`
		Id           int          `json:"id"`
		CampaignId   string       `json:"campaign_id"`
		CampaignName string       `json:"name"`
		Objective    string       `json:"campaign_objective"`
		Country      string       `json:"country"`
		Advertiser   string       `json:"advertiser"`
		IsDelCamp    bool         `json:"is_del_camp"`
		CPCR         string       `json:"cost_per_conversion"`
		AgencyFee    string       `json:"agency_fee"`
		MCC          string       `json:"mcc"`
		DataConfig   []DataConfig `json:"data"`
	}

	DataConfig struct {
		Id                        int     `json:"id"`            // <-
		URLServiceKey             string  `json:"urlservicekey"` // lp-idxlgazpastelmbv
		CampaignId                string  `json:"campaign_id"`   // ID1
		CampaignName              string  `json:"name"`
		Objective                 string  `json:"objective"`
		Country                   string  `json:"country"`
		Advertiser                string  `json:"advertiser"`
		Operator                  string  `json:"operator"`
		Partner                   string  `json:"partner"`
		Aggregator                string  `json:"aggregator"`
		Adnet                     string  `json:"adnet"`
		Service                   string  `json:"service"`
		Keyword                   string  `json:"keyword"`
		SubKeyword                string  `json:"subkeyword"`
		IsBillable                bool    `json:"is_billable"`
		Plan                      string  `json:"plan"`
		PO                        string  `json:"po"`
		Cost                      string  `json:"cost"`
		PubId                     string  `json:"pubid"`
		ShortCode                 string  `json:"short_code"`
		DeviceType                string  `json:"device_type"`
		OS                        string  `json:"os"`
		URLType                   string  `json:"url_type"`
		ClickType                 int     `json:"click_type"`
		ClickDelay                int     `json:"click_delay"`
		ClientType                string  `json:"client_type"`
		TrafficSource             bool    `json:"traffic_source"`
		UniqueClick               bool    `json:"unique_click"`
		URLBanner                 string  `json:"url_banner"`
		URLLanding                string  `json:"url_landing"`
		URLWarpLanding            string  `json:"url_warp_landing"`
		URLService                string  `json:"url_service"`
		URLTFCSmartlink           string  `json:"url_tfc_or_smartlink"`
		GlobPost                  bool    `json:"glob_post"`
		URLGlobPost               string  `json:"url_glob_post"`
		CustomIntegration         string  `json:"custom_integration"`
		IPAddress                 []uint8 `json:"ip_address"`
		ISP                       string  `json:"isp"`
		IsActive                  bool    `json:"is_active"`
		MOCapping                 int     `json:"mo_capping"`
		MOCappingService          int     `json:"mo_capping_service"`
		CounterMOCapping          int     `json:"counter_mo_capping"`
		CounterMOCappingService   int     `json:"counter_mo_capping_service"`
		StatusCapping             bool    `json:"status_capping"`
		KPIUpperLimitCapping      int     `json:"kpi_upper_limit_capping"`
		IsMachineLearningCapping  bool    `json:"is_machine_learning_capping"`
		RatioSend                 int     `json:"ratio_send"`
		RatioReceive              int     `json:"ratio_receive"`
		CounterMORatio            int     `json:"counter_mo_ratio"`
		StatusRatio               bool    `json:"status_ratio"`
		KPIUpperLimitRatioSend    int     `json:"kpi_upper_limit_ratio_send"`
		KPIUpperLimitRatioReceive int     `json:"kpi_upper_limit_ratio_receive"`
		IsMachineLearningRatio    bool    `json:"is_machine_learning_ratio"`
		APIURL                    string  `json:"api_url"`
		LastUpdate                string  `json:"last_update"`
		LastUpdateCapping         string  `json:"last_update_capping"`
		CPCR                      string  `json:"cost_per_conversion"`
		AgencyFee                 string  `json:"agency_fee"`
		TargetDailyBudget         string  `json:"target_daily_budget"`
		BudgetUsage               string  `json:"budget_usage"`
		TechnicalFee              string  `json:"technical_fee"`
		URLPostback               string  `json:"url_postback"`
		PostbackMethod            string  `json:"postback_method"`
		MainstreamLpType          string  `json:"mainstream_lp_type"`
		Title                     string  `json:"title"`
		TitleOriginal             string  `json:"title_original"`
		TitleColor                string  `json:"title_color"`
		TitleStyle                string  `json:"title_style"`
		TitlePageType             string  `json:"title_page_type"`
		TitleFontSize             string  `json:"title_font_size"`
		SubTitle                  string  `json:"sub_title"`
		SubTitleOriginal          string  `json:"sub_title_original"`
		SubTitleColor             string  `json:"sub_title_color"`
		SubTitleStyle             string  `json:"sub_title_style"`
		SubTitlePageType          string  `json:"sub_title_page_type"`
		SubTitleFontSize          string  `json:"sub_title_font_size"`
		BackgroundURL             string  `json:"background_url"`
		BackgroundColor           string  `json:"background_color"`
		LogoURL                   string  `json:"logo_url"`
		URLBannerOriginal         string  `json:"url_banner_original"`
		Tnc                       string  `json:"tnc"`
		TncOriginal               string  `json:"tnc_original"`
		TncColor                  string  `json:"tnc_color"`
		TncStyle                  string  `json:"tnc_style"`
		TncPageType               string  `json:"tnc_page_type"`
		TncFontSize               string  `json:"tnc_font_size"`
		ButtonSubscribe           string  `json:"button_subscribe"`
		ButtonSubscribeOriginal   string  `json:"button_subscribe_original"`
		ButtonSubscribeColor      string  `json:"button_subscribe_color"`
		StatusSubmitKeyMainstream bool    `json:"status_submit_key_mainstream"`
		KeyMainstream             string  `json:"key_mainstream"`
		Channel                   string  `json:"channel"`
		GoogleSheet               string  `json:"google_sheet"`
		GoogleSheetBillable       string  `json:"google_sheet_billable"`
		Currency                  string  `json:"currency"`
		MCC                       string  `json:"mcc"`
		ClickableAnywhere         bool    `json:"clickable_anywhere"`
		NonTargetURL              string  `json:"non_target_url"`
		EnableIpRanges            bool    `json:"enable_ip_ranges"`
		ConversionName            string  `json:"conversion_name"`
		DomainService             string  `json:"domain_service"`
		CampaignDetailName        string  `json:"campaign_detail_name"`
		Prefix                    string  `json:"prefix"`
		CountryDialingCode        string  `json:"country_dialing_code"`
		UnusedTrafficRedirectType string  `json:"unused_traffic_redirect_type"`
		CompanyLegalName          string  `json:"company_legal_name"`
		CompanyAddress            string  `json:"company_address"`
		CompanyEmail              string  `json:"company_email"`
		CompanyPhone              string  `json:"company_phone"`
		ServicePrice              float64 `json:"service_price"`
		PortalURL                 string  `json:"portal_url"`
		IsEvina 				  bool    `json:"is_evina"`
		EvinaRedirectFraudURL     string  `json:"evina_redirect_fraud_url"`
	}

	//'{"id":1,"urlservicekey":"idtelgempastelmbv","campaign_id":"ID01","country":"ID","partner":"pass","operator":"telkomsel","aggregator":"telesat","service":"gazy","short_code":"1234","adnet":"mbv","keyword":"gazy","subkeyword":"","is_billable":false,"plan":"","traffic":0,"landing":0,"click":0,"redirect":0,"traffic_data":[],"landing_data":[],"click_data":[],"redirect_data":[]}'

	DataCounter struct {
		CampaignDetailId int                         `json:"campaign_detail_id"`
		URLServiceKey    string                      `json:"urlservicekey"`
		CampaignId       string                      `json:"campaign_id"`
		Country          string                      `json:"country"`
		Partner          string                      `json:"partner"`
		Operator         string                      `json:"operator"`
		Aggregator       string                      `json:"aggregator"`
		Service          string                      `json:"service"`
		ShortCode        string                      `json:"short_code"`
		Adnet            string                      `json:"adnet"`
		Keyword          string                      `json:"keyword"`
		SubKeyword       string                      `json:"subkeyword"`
		IsBillable       bool                        `json:"is_billable"`
		Plan             string                      `json:"plan"`
		Traffic          int                         `json:"traffic"`
		Landing          int                         `json:"landing"`
		Click            int                         `json:"click"`
		Redirect         int                         `json:"redirect"`
		MOReceived       int                         `json:"moreceived"`
		Postback         int                         `json:"postback"`
		TotalFP          int                         `json:"totalfp"`
		TrafficData      []DataCounterDetail         `json:"traffic_data"`
		LandingData      []DataCounterDetail         `json:"landing_data"`
		ClickData        []DataCounterDetail         `json:"click_data"`
		RedirectData     []DataCounterDetail         `json:"redirect_data"`
		MOData           []DataCounterDetailInternal `json:"mo_data"`
		PostbackData     []PixelStorage              `json:"postback_data"`
		FPData           []DataCounterDetailInternal `json:"fp_data"`
	}

	DataCounterDetail struct {
		Date       string `json:"date"`
		Time       string `json:"time"`
		HTTPStatus string `json:"http_status"`
		OS         string `json:"os"`
		IsOperator bool   `json:"is_operator"`
	}

	DataCounterDetailInternal struct {
		Date   string `json:"date"`
		Msisdn string `json:"msisdn"`
		TrxId  string `json:"trxid"`
		Pixel  string `json:"pixel"`
		Code   string `json:"code"`
		Status bool   `json:"status"`
	}

	AlertData struct {
		Platform string     `json:"platform"`
		DataUser []DataUser `json:"data"`
	}

	DataUser struct {
		Name   string `json:"name"`
		UserId int    `json:"user_id"`
	}

	Summary struct {
		SummaryDate        string `json:"summary_date"`
		CRMO               string `json:"cr_mo"`
		CRPostback         string `json:"cr_postback"`
		SuccessFP          string `json:"success_fp"`
		BillRate           string `json:"billrate"`
		SBAF               string `json:"sbaf"`
		SAAF               string `json:"saaf"`
		CPA                string `json:"cpa"`
		Revenue            string `json:"revenue"`
		URLWarpLanding     string `json:"url_warp_landing"`
		URLLanding         string `json:"url_landing"`
		MOCapping          int    `json:"mo_capping"`
		RatioSend          int    `json:"ratio_send"`
		RatioReceive       int    `json:"ratio_receive"`
		ClientType         string `json:"client_type"`
		CPCR               string `json:"cost_per_conversion"`
		AgencyFee          string `json:"agency_fee"`
		TechnicalFee       string `json:"technical_fee"`
		TotalWakiAgencyFee string `json:"total_waki_agency_fee"`
		TDB                string `json:"target_daily_budget"`
		BudgetUsage        string `json:"budget_usage"`
		PricePerMO         string `json:"price_per_mo"`
		IsActive           bool   `json:"is_active"`
		CampaignId         string `json:"campaign_id"`
		CampaignName       string `json:"campaign_name"`
		Country            string `json:"country"`
		Partner            string `json:"partner"`
		Operator           string `json:"operator"`
		URLServiceKey      string `json:"urlservicekey"`
		Aggregator         string `json:"aggregator"`
		Service            string `json:"service"`
		Adnet              string `json:"adnet"`
		ShortCode          string `json:"shortcode"`
		PO                 string `json:"po"`
		TotalTraffic       int    `json:"total_traffic"`
		TotalLanding       int    `json:"total_landing"`
		TotalClicked       int    `json:"total_clicked"`
		TotalRedirect      int    `json:"total_redirect"`
		TotalMOReceived    int    `json:"total_moreceived"`
		TotalPostback      int    `json:"total_postback"`
		TotalFP            int    `json:"total_fp"`
		TotalROI           int    `json:"total_roi"`
		DataReserved1      string `json:"data_reserved1"`
		DataReserved2      string `json:"data_reserved2"`
		DataReserved3      string `json:"data_reserved3"`
		DataReserved4      string `json:"data_reserved4"`
		DataReserved5      string `json:"data_reserved5"`
		DataReserved6      string `json:"data_reserved6"`
		DataReserved7      string `json:"data_reserved7"`
		DataReserved8      string `json:"data_reserved8"`
		DataReserved9      string `json:"data_reserved9"`
		DataReserved10     string `json:"data_reserved10"`
		TotalReserved1     int    `json:"total_reserved1"`
		TotalReserved2     int    `json:"total_reserved2"`
		TotalReserved3     int    `json:"total_reserved3"`
		TotalReserved4     int    `json:"total_reserved4"`
		TotalReserved5     int    `json:"total_reserved5"`
		TotalReserved6     int    `json:"total_reserved6"`
		TotalReserved7     int    `json:"total_reserved7"`
		TotalReserved8     int    `json:"total_reserved8"`
		TotalReserved9     int    `json:"total_reserved9"`
		TotalReserved10    int    `json:"total_reserved10"`
		CampaignObjective  string `json:"campaign_objective"`
		Channel            string `json:"channel"`
	}

	InquiryCampID struct {
		URLServiceKey string `json:"urlservicekey" query:"urlservicekey"`
		Country       string `json:"country" query:"country"`
		Operator      string `json:"operator" query:"operator"`
		Partner       string `json:"partner" query:"partner"`
		Aggregator    string `json:"aggregator" query:"aggregator"`
		Adnet         string `json:"adnet" query:"adnet"`
		Service       string `json:"service" query:"service"`
		URLLanding    string `json:"url_landing" query:"url_landing"`
		URLService    string `json:"url_service" query:"url_service"`
	}

	InquiryAPICampID struct {
		Country   string `json:"country" query:"country"`
		Operator  string `json:"operator" query:"operator"`
		Service   string `json:"service" query:"service"`
		Adnet     string `json:"adnet" query:"adnet"`
	}

	InquiryAPICampIDResult struct {
		URLServiceKey string `json:"urlservicekey"`
		Country       string `json:"country"`
		Operator      string `json:"operator"`
		Partner       string `json:"partner"`
		Aggregator    string `json:"aggregator"`
		Adnet         string `json:"adnet"`
		Service       string `json:"service"`
	}
)
