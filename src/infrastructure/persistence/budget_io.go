package persistence

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func (m *BaseModel) CreateBudgetIO(budgetIO *entity.BudgetIO) error {
	return m.DB.Create(budgetIO).Error
}

func (r *BaseModel) GetDisplayBudgetIO(o entity.DisplayBudgetIO, allowedCompanies []string) ([]entity.BudgetIO, int64, error) {
	var rows *sql.Rows
	var err error
	var total_rows int64


	query := r.DB.Model(&entity.BudgetIO{}).Where("company IN ?", allowedCompanies).Where("status != ?", "approved")

	if o.CampaignType != "" {
		query.Where("campaign_type = ? ", o.CampaignType)
	} 

	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Partner != "" {
			query = query.Where("partner = ?", o.Partner)
		}
		if o.Continent != "" {
			query = query.Where("continent = ?", o.Continent)
		}

		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("DATE(created_at) = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("DATE(created_at) = CURRENT_DATE - INTERVAL '1 DAY'")
			case "LAST7DAY":
				query = query.Where("DATE(created_at) BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST30DAY":
				query = query.Where("DATE(created_at) BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("DATE(created_at) >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("DATE(created_at) BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				query = query.Where("DATE(created_at) BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			default:
				query = query.Where("DATE(created_at) = ?", o.DateRange)
			}
		} else {
			query = query.Where("DATE(created_at) >= DATE_TRUNC('month', CURRENT_DATE)")
		}
	}

	if o.OrderColumn != "" {
		dir := "ASC"
		if strings.ToUpper(o.OrderDir) == "DESC" {
			dir = "DESC"
		}

		switch o.OrderColumn {
		default:
			query = query.Order(fmt.Sprintf("%s %s", o.OrderColumn, dir))
		}
	} else {
		query = query.Order("DATE(created_at) DESC").Order("id DESC")
	}

	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, err = query_limit.Rows()
	if err != nil {
		return []entity.BudgetIO{}, 0, err
	}
	defer rows.Close()

	var ss []entity.BudgetIO
	for rows.Next() {
		var s entity.BudgetIO
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", len(ss)))

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetDisplayBudgetIOAll(o entity.DisplayBudgetIO) ([]entity.BudgetIO, int64, error) {
	var rows *sql.Rows
	var err error
	var total_rows int64


	query := r.DB.Model(&entity.BudgetIO{}).Where("status != ?", "approved")

	if o.CampaignType != "" {
		query.Where("campaign_type = ? ", o.CampaignType)
	} 

	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Partner != "" {
			query = query.Where("partner = ?", o.Partner)
		}
		if o.Continent != "" {
			query = query.Where("continent = ?", o.Continent)
		}

		if o.DateRange != "" {
			switch strings.ToUpper(o.DateRange) {
			case "TODAY":
				query = query.Where("DATE(created_at) = CURRENT_DATE")
			case "YESTERDAY":
				query = query.Where("DATE(created_at) = CURRENT_DATE - INTERVAL '1 DAY'")
			case "LAST7DAY":
				query = query.Where("DATE(created_at) BETWEEN CURRENT_DATE - INTERVAL '7 DAY' AND CURRENT_DATE")
			case "LAST30DAY":
				query = query.Where("DATE(created_at) BETWEEN CURRENT_DATE - INTERVAL '30 DAY' AND CURRENT_DATE")
			case "THISMONTH":
				query = query.Where("DATE(created_at) >= DATE_TRUNC('month', CURRENT_DATE)")
			case "LASTMONTH":
				query = query.Where("DATE(created_at) BETWEEN DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 MONTH') AND DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 DAY'")
			case "CUSTOMRANGE":
				query = query.Where("DATE(created_at) BETWEEN ? AND ?", o.DateBefore, o.DateAfter)
			default:
				query = query.Where("DATE(created_at) = ?", o.DateRange)
			}
		} else {
			query = query.Where("DATE(created_at) >= DATE_TRUNC('month', CURRENT_DATE)")
		}
	}

	if o.OrderColumn != "" {
		dir := "ASC"
		if strings.ToUpper(o.OrderDir) == "DESC" {
			dir = "DESC"
		}

		switch o.OrderColumn {
		default:
			query = query.Order(fmt.Sprintf("%s %s", o.OrderColumn, dir))
		}
	} else {
		query = query.Order("DATE(created_at) DESC").Order("id DESC")
	}

	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, err = query_limit.Rows()
	if err != nil {
		return []entity.BudgetIO{}, 0, err
	}
	defer rows.Close()

	var ss []entity.BudgetIO
	for rows.Next() {
		var s entity.BudgetIO
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", len(ss)))

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetDisplayBudgetIOApproved(o entity.DisplayBudgetIO, allowedCompanies []string) ([]entity.BudgetIO, int64, error) {
	var rows *sql.Rows
	var err error
	var total_rows int64


	query := r.DB.Model(&entity.BudgetIO{}).Where("company IN ?", allowedCompanies).Where("status = ?", "approved")

	if o.CampaignType != "" {
		query.Where("campaign_type = ? ", o.CampaignType)
	} 

	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Partner != "" {
			query = query.Where("partner = ?", o.Partner)
		}
		if o.Continent != "" {
			query = query.Where("continent = ?", o.Continent)
		}

		if o.Keyword != "" {
			keywordPattern := "%" + o.Keyword + "%"
			query = query.Where(
				r.DB.Where("io_id ILIKE ?", keywordPattern).
					Or("company_group_name ILIKE ?", keywordPattern).
					Or("submitted_by ILIKE ?", keywordPattern),
			)
		}

		if o.DateRange != "" {
			query = query.Where("month = ?", o.DateRange)
		} else {
			query = query.Where("month = TO_CHAR(CURRENT_DATE, 'YYYY-MM')")
		}
	}

	if o.OrderColumn != "" {
		dir := "ASC"
		if strings.ToUpper(o.OrderDir) == "DESC" {
			dir = "DESC"
		}

		switch o.OrderColumn {
		default:
			query = query.Order(fmt.Sprintf("%s %s", o.OrderColumn, dir))
		}
	} else {
		query = query.Order("DATE(created_at) DESC").Order("id DESC")
	}

	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, err = query_limit.Rows()
	if err != nil {
		return []entity.BudgetIO{}, 0, err
	}
	defer rows.Close()

	var ss []entity.BudgetIO
	for rows.Next() {
		var s entity.BudgetIO
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", len(ss)))

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetDisplayBudgetIOApprovedAll(o entity.DisplayBudgetIO) ([]entity.BudgetIO, int64, error) {
	var rows *sql.Rows
	var err error
	var total_rows int64


	query := r.DB.Model(&entity.BudgetIO{}).Where("status = ?", "approved")

	if o.CampaignType != "" {
		query.Where("campaign_type = ? ", o.CampaignType)
	} 

	if o.Action == "Search" {
		if o.Country != "" {
			query = query.Where("country = ?", o.Country)
		}
		if o.Partner != "" {
			query = query.Where("partner = ?", o.Partner)
		}
		if o.Continent != "" {
			query = query.Where("continent = ?", o.Continent)
		}

		if o.Keyword != "" {
			keywordPattern := "%" + o.Keyword + "%"
			query = query.Where(
				r.DB.Where("io_id ILIKE ?", keywordPattern).
					Or("company_group_name ILIKE ?", keywordPattern).
					Or("submitted_by ILIKE ?", keywordPattern),
			)
		}

		if o.DateRange != "" {
			query = query.Where("month = ?", o.DateRange)
		} else {
			query = query.Where("month = TO_CHAR(CURRENT_DATE, 'YYYY-MM')")
		}
	}

	if o.OrderColumn != "" {
		dir := "ASC"
		if strings.ToUpper(o.OrderDir) == "DESC" {
			dir = "DESC"
		}

		switch o.OrderColumn {
		default:
			query = query.Order(fmt.Sprintf("%s %s", o.OrderColumn, dir))
		}
	} else {
		query = query.Order("DATE(created_at) DESC").Order("id DESC")
	}

	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, err = query_limit.Rows()
	if err != nil {
		return []entity.BudgetIO{}, 0, err
	}
	defer rows.Close()

	var ss []entity.BudgetIO
	for rows.Next() {
		var s entity.BudgetIO
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", len(ss)))

	return ss, total_rows, rows.Err()
}

func (r *BaseModel) GetDisplaySummaryBudgetIO(o entity.DisplaySummaryBudgetIO) ([]entity.SummaryBudgetIO, int64, error) {

    query := r.DB.Model(&entity.SummaryBudgetIO{})

    if o.Action == "Search" {
        if o.Country != "" {
            query = query.Where("country = ?", o.Country)
        }
        if o.Partner != "" {
            query = query.Where("partner = ?", o.Partner)
        }
        if o.Continent != "" {
            query = query.Where("continent = ?", o.Continent)
        }

        if o.DateRange != "" {
            query = query.Where("month = ?", o.DateRange)
        } else {
            query = query.Where("month = TO_CHAR(CURRENT_DATE, 'YYYY-MM')")
        }
    }

    if o.OrderColumn != "" {
        dir := "ASC"
        if strings.ToUpper(o.OrderDir) == "DESC" {
            dir = "DESC"
        }
        query = query.Order(fmt.Sprintf("%s %s", o.OrderColumn, dir))
    } else {
        query = query.Order("month DESC, id DESC")
    }

    rows, err := query.Rows()
    if err != nil {
        return []entity.SummaryBudgetIO{}, 0, err
    }
    defer rows.Close()

    var ss []entity.SummaryBudgetIO
    for rows.Next() {
        var s entity.SummaryBudgetIO
        r.DB.ScanRows(rows, &s)
        ss = append(ss, s)
    }

    r.Logs.Debug(fmt.Sprintf("Total data : %d ... \n", len(ss)))

    return ss, int64(len(ss)), rows.Err()
}

func (r *BaseModel) UpdateSummaryBudgetIO(req entity.UpdateSummaryBudgetIORequest) error {
	updates := map[string]interface{}{"updated_at": time.Now()}
	if req.MOTarget != nil {
		updates["monthly_mo_target"] = *req.MOTarget
	}
	if req.IOTarget != nil {
		updates["monthly_spend_target"] = *req.IOTarget
	}
	if req.TargetCAC != nil {
		updates["target_cac"] = *req.TargetCAC
	}
	if req.ROI != nil {
		updates["target_roi"] = int(*req.ROI)
	}

	if req.ID > 0 {
		result := r.DB.Model(&entity.BudgetIO{}).Where("id = ?", req.ID).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			return nil
		}
	}

	if req.Country == "" || req.Month == "" {
		return fmt.Errorf("UpdateSummaryBudgetIO: id or (country, month) required")
	}
	return r.DB.Model(&entity.BudgetIO{}).
		Where("country = ? AND month = ?", req.Country, req.Month).
		Updates(updates).Error
}
