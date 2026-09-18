package entity

import (
	"time"

	"gorm.io/gorm"
)

type CampaignROASCohort struct {
	gorm.Model
	ID                        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	SummaryDate               time.Time `gorm:"type:date;uniqueIndex:idx_camproascohort_unique;not null" json:"summary_date"`
	URLServiceKey             string    `gorm:"uniqueIndex:idx_camproascohort_unique;not null;size:80" json:"url_service_key"`
	Country                   string    `gorm:"uniqueIndex:idx_camproascohort_unique;size:10" json:"country"`
	CountryName               string    `gorm:"size:80" json:"country_name"`
	Operator                  string    `gorm:"uniqueIndex:idx_camproascohort_unique;size:255" json:"operator"`
	Service                   string    `gorm:"uniqueIndex:idx_camproascohort_unique;size:255" json:"service"`
	Adnet                     string    `gorm:"uniqueIndex:idx_camproascohort_unique;size:255" json:"adnet"`
	MO                        int       `gorm:"column:mo;default:0" json:"mo"`
	EstimatedROAS             float64   `gorm:"type:double precision;default:0" json:"estimated_roas"`
	NetROASActual             float64   `gorm:"type:double precision;default:0" json:"net_roas_actual"`
	RoiMonthsPayback          float64   `gorm:"type:double precision;default:0" json:"roi_months_payback"`
	EstimatedLTV              float64   `gorm:"type:double precision;default:0" json:"estimated_ltv"`
	EstimatedGrossRevenueFull float64   `gorm:"type:double precision;default:0" json:"estimated_gross_revenue_full"`
	HealthStatus              string    `gorm:"size:30" json:"health_status"`
	ComputedAt                time.Time `gorm:"type:timestamptz" json:"computed_at"`
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

func (CampaignROASCohort) TableName() string {
	return "campaign_roas_cohorts"
}
