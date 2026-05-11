package persistence

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
	"gorm.io/gorm"
)

func (r *BaseModel) GetAlertReportAll(o entity.DisplayAlertReport, v string) ([]entity.SummaryAll, error) {
	var (
		rows *sql.Rows
		// err  error
		query *gorm.DB
	)

	switch v {
	case "all":
		var summaryAll []entity.SummaryAll
		for _, v := range []string{"nomo", "crdrop", "ratioalert", "capreached"} {
			s, err := r.GetAlertReportAll(o, v)
			if err != nil {
				return []entity.SummaryAll{}, err
			}
			summaryAll = append(summaryAll, s...)
		}
		return summaryAll, nil
	case "nomo":
		query = r.DB.Model(&entity.SummaryMo{}).Select("summary_mos.*, 'No MO' as error")
	case "crdrop":
		query = r.DB.Model(&entity.SummaryCr{}).Select("summary_crs.*, 'CR Drop' as error")
	case "ratioalert":
		query = r.DB.Model(&entity.SummaryRatio{}).Select("summary_ratios.*, 'Ratio Alert' as error")
	case "capreached":
		query = r.DB.Model(&entity.SummaryCapping{}).Select("summary_cappings.*, 'Cap Reached' as error")
	default:
		return []entity.SummaryAll{}, nil
	}

	rows, err := r.filter(o, query)
	if err != nil {
		r.Logs.Error(fmt.Sprintf("Error %s when reading table", err))
		return []entity.SummaryAll{}, err
	}

	if rows == nil {
		return []entity.SummaryAll{}, nil
	}
	defer rows.Close()

	var ss []entity.SummaryAll

	for rows.Next() {
		var s entity.SummaryAll
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", len(ss)))
	return ss, rows.Err()

}

func (r *BaseModel) filter(o entity.DisplayAlertReport, query *gorm.DB) (rows *sql.Rows, err error) {
	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("LOWER(country) = LOWER(?)", o.Country)
		}
		if o.Operator != "" {
			query = query.Where("LOWER(operator) = LOWER(?)", o.Operator)
		}
		if o.CampaignName != "" {
			query = query.Where("LOWER(campaign_name) = LOWER(?)", o.CampaignName)
		}
		if o.Service != "" {
			query = query.Where("LOWER(service) = LOWER(?)", o.Service)
		}
		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("date_send = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '1 DAY' AND CURRENT_DATE")
			case "LAST7DAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST30DAY":
				query = query.Where("date_send BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("date_send >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LAST MONTH":
				query = query.Where("summary_date BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOM RANGE":
				query = query.Where("summary_date BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			default:
				query = query.Where("summary_date = ?", o.DateRange)
			}
		}
		rows, err = query.Order("summary_date DESC").Order("ID DESC").Rows()

		if err != nil {
			return nil, err
		}
	} else {
		rows, err = query.Order("summary_date DESC").Order("ID DESC").Rows()
		if err != nil {
			return nil, err
		}
	}

	return rows, nil
}

func (r *BaseModel) UpdateStatusAlert(ID string, Status bool, time, v string) error {
	var result error

	switch v {
	case "nomo":
		result = r.DB.Model(&entity.SummaryMo{}).Where("id = ? ", ID).Update("status", Status).Error
	case "crdrop":
		result = r.DB.Model(&entity.SummaryCr{}).Where("id = ? ", ID).Update("status", Status).Error
	case "ratioalert":
		result = r.DB.Model(&entity.SummaryRatio{}).Where("id = ? ", ID).Update("status", Status).Error
	case "capreached":
		result = r.DB.Model(&entity.SummaryCapping{}).Where("id = ? ", ID).Update("status", Status).Error
	}
	r.Logs.Debug(fmt.Sprintf("Update Status Alert %s on %s with result : %v", v, time, result))
	return result

}
