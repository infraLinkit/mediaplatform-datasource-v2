package entity

import (
	"time"

	"gorm.io/gorm"
)

type (
	IncSummaryCampaign struct {
		gorm.Model
		ID                int       `gorm:"primaryKey;autoIncrement" json:"id"`
		SummaryDate       time.Time `gorm:"type:date;uniqueIndex:idx_incsumunique" json:"summary_date"`
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
		Landing           int       `gorm:"length:20;default:0" json:"landing"`
		Clicked 		  int       `gorm:"length:20;default:0" json:"clicked"`
		MoReceived        int       `gorm:"length:20;default:0" json:"mo_received"`
		Postback          int       `gorm:"length:20;default:0" json:"postback"`
		POAF              float64   `gorm:"type:double precision;default:0" json:"po_af"`
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}
)

// HOOK or Trigger
// Important Create Name Func per hook entity!, Name returns the name of the plugin
/* func (sum *IncSummaryCampaign) Name() string {
	return "inc_summary_campaign_trigger_update_landing_mo_postback"
}

// Initialize is called when you register the plugin via db.Use()
func (sum *IncSummaryCampaign) Initialize(db *gorm.DB) error {
	// Register a callback before the standard Create operation
	db.Callback().Update().After("gorm:update").
		Register(sum.Name(), sum.AfterUpdate)
	return nil
}

func (sum *IncSummaryCampaign) AfterUpdate(db *gorm.DB) {

	if db.Error == nil {
		fmt.Println("--> Global Hook: A record is about to be updated!")
		// You can access the data being created:
		// model := db.Statement.Dest
		if _, ok := db.Statement.Dest.(*IncSummaryCampaign); ok {

			haystack := []string{
				"ID-147",
				"ID-150",
				"ID-151",
			}

			if sum.MoReceived > 0 && external.InArray(sum.URLServiceKey, haystack) {

				var (
					c  Campaign
					cd CampaignDetail
				)

				db.Model(&CampaignDetail{}).
					Where("url_service_key = ?", sum.URLServiceKey).First(&cd)

				db.Model(&Campaign{}).
					Where("id = ?", cd.CampaignId).First(&c)

				sc := sum.FormulaCPA(db, cd, c)

				// Update the record in the database
				result := db.Model(&SummaryCampaign{}).Where("summary_date = ? AND url_service_key = ?", sum.SummaryDate.Format("2006-01-02"), sum.URLServiceKey).Updates(&sc)

				if result.RowsAffected == 0 {

					// Create new record if not exist
					db.Create(sc)
				}
			}

		}
	}
}

func (sum *IncSummaryCampaign) FormulaCPA(tx *gorm.DB, cd CampaignDetail, c Campaign) SummaryCampaign {

	var sc SummaryCampaign
	sc.URLServiceKey = cd.URLServiceKey
	sc.Status = cd.IsActive
	sc.SummaryDate = sum.SummaryDate
	sc.CampaignId = cd.CampaignId
	sc.CampaignName = c.Name
	sc.Company = c.Advertiser
	sc.Country = cd.Country
	sc.Partner = cd.Partner
	sc.Operator = cd.Operator
	sc.Aggregator = cd.Aggregator
	sc.Service = cd.Service
	sc.Adnet = cd.Adnet
	sc.ShortCode = cd.ShortCode
	sc.Traffic = sum.Landing
	sc.Landing = sum.Landing
	sc.MoReceived = sum.MoReceived
	sc.CR = sc.CrMO
	sc.Postback = sum.Postback
	sc.TotalFP = 0
	sc.SuccessFP = 0
	sc.Billrate = 0

	po, _ := strconv.ParseFloat(cd.PO, 64)
	sc.PO = po
	sc.Cost = 0
	sc.URLAfter = cd.URLWarpLanding
	sc.URLBefore = cd.URLLanding
	sc.MOLimit = cd.MOCapping
	sc.ROI = 0
	sc.RatioSend = cd.RatioSend
	sc.RatioReceive = cd.RatioReceive
	sc.ClientType = cd.ClientType
	sc.CostPerConversion = cd.CostPerConversion
	sc.AgencyFee = cd.AgencyFee
	sc.URLServiceKey = cd.URLServiceKey
	sc.CampaignObjective = c.CampaignObjective
	sc.Channel = cd.Channel

	// CR conversion to mo received
	//var cr_mo float64
	if sum.MoReceived > 0 && sum.Landing > 0 {
		sc.CrMO, _ = strconv.ParseFloat(fmt.Sprintf("%f", float64(sum.MoReceived)/float64(sum.Landing)), 64)
	} else {
		sc.CrMO, _ = strconv.ParseFloat("0", 64)
	}

	//var cr_postback float64
	if sum.Postback > 0 && sum.Landing > 0 {
		sc.CrPostback, _ = strconv.ParseFloat(fmt.Sprintf("%f", float64(sum.Postback)/float64(sum.Landing)), 64)
	} else {
		sc.CrPostback, _ = strconv.ParseFloat("0", 64)
	}

	// GET SBAF (spending before agency fee)
	//payout, _ := strconv.ParseFloat(strings.TrimSpace(cd.PO), 64)
	//payout := sum.PO
	//h.Logs.Info(fmt.Sprintf("PO : %f", payout))

	mo_sent := float64(sum.Postback)
	//sbaf := payout * mo_sent
	sc.SBAF = po * mo_sent //sbaf

	// GET total waki agency fee
	//sc.CostPerConversion, _ = strconv.ParseFloat(strings.TrimSpace(gs.CPCR), 64)
	//cost_per_conversion := sum.CostPerConversion
	//sum.AgencyFee, _ = strconv.ParseFloat(strings.TrimSpace(cd.AgencyFee), 64)
	//cd.AgencyFee = cd.AgencyFee / 100
	mo_received := float64(sum.MoReceived)
	sc.TotalWakiAgencyFee = (cd.CostPerConversion * mo_received) + (cd.AgencyFee * (cd.CostPerConversion + (cd.CostPerConversion * mo_received)))

	// GET SAAF (spending after agency fee)
	//saaf := total_waki_agency_fee + sbaf
	//tech_fee, _ := strconv.Atoi(gs.TechnicalFee)
	//sum.TechnicalFee, _ = strconv.ParseFloat(strings.TrimSpace(gs.TechnicalFee), 64)
	//sum.TechnicalFee = sum.TechnicalFee / 100

	sc.TechnicalFee = cd.TechnicalFee * (sc.SBAF + sc.TotalWakiAgencyFee)

	sc.SAAF = sc.TotalWakiAgencyFee + sc.SBAF + sc.TechnicalFee //saaf
	if strings.ToLower(cd.ClientType) == "external" {
		sc.SAAF = mo_received * sum.POAF
	}

	// GET eCPA
	//cpa := float64(0)
	if sc.SAAF > 0 && mo_received > 0 {
		sc.CPA = sc.SAAF / mo_received
	}

	// Revenue
	//revenue := float64(0)
	if sc.SAAF > 0 && sc.SBAF > 0 {
		sc.Revenue = sc.SAAF - sc.SBAF
	}

	sc.PricePerMO = sc.SAAF / mo_received

	if sum.Landing == 0 || sum.MoReceived == 0 || sum.Postback == 0 {
		//price_per_mo_string = "0"
		sc.PricePerMO = 0
	}

	return sc
} */
