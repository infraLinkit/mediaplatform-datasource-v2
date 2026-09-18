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

// GetDisplaySummaryBudgetIO — matches datasource's daily-granularity summary_budget_ios
// schema: buckets each day into week1-5 by day-of-month, LEFT JOINed against budget_ios
// (v2's approval-workflow schema: monthly_spend_target/monthly_mo_target/target_roi in
// place of datasource's io_target/mo_target/roi; datasource's ltv/roas have no v2
// equivalent column, so they default to 0 same as a LEFT JOIN miss would).
func (r *BaseModel) GetDisplaySummaryBudgetIO(o entity.DisplaySummaryBudgetIO) ([]entity.SummaryBudgetIOAgg, int64, error) {
	var whereClause []string
	var args []interface{}

	if o.Action == "Search" {
		if o.Country != "" {
			cs := []string{}
			for _, c := range strings.Split(o.Country, ",") {
				if c = strings.TrimSpace(c); c != "" {
					cs = append(cs, c)
				}
			}
			if len(cs) == 1 {
				whereClause = append(whereClause, "s.country = ?")
				args = append(args, cs[0])
			} else if len(cs) > 1 {
				ph := make([]string, len(cs))
				for i := range ph {
					ph[i] = "?"
				}
				whereClause = append(whereClause, fmt.Sprintf("s.country IN (%s)", strings.Join(ph, ",")))
				for _, c := range cs {
					args = append(args, c)
				}
			}
		}
		if o.Partner != "" {
			whereClause = append(whereClause, "s.partner = ?")
			args = append(args, o.Partner)
		}
		if o.Continent != "" {
			whereClause = append(whereClause, "s.continent = ?")
			args = append(args, o.Continent)
		}
		if o.CampaignType != "" {
			whereClause = append(whereClause, "s.campaign_type LIKE ?")
			args = append(args, "%"+strings.ToUpper(o.CampaignType)+"%")
		}
		if o.Channel != "" {
			whereClause = append(whereClause, "s.channel = ?")
			args = append(args, o.Channel)
		}
		if o.Company != "" {
			whereClause = append(whereClause, "s.company = ?")
			args = append(args, o.Company)
		}
		if o.Operator != "" {
			whereClause = append(whereClause, "s.operator = ?")
			args = append(args, o.Operator)
		}
		if o.Service != "" {
			whereClause = append(whereClause, "s.service = ?")
			args = append(args, o.Service)
		}
		if o.ClientType != "" {
			whereClause = append(whereClause, "s.client_type = ?")
			args = append(args, o.ClientType)
		}
		if o.DateRange != "" {
			whereClause = append(whereClause, "s.month = ?")
			args = append(args, o.DateRange)
		} else if o.DateBefore != "" && o.DateAfter != "" {
			whereClause = append(whereClause, "s.summary_date BETWEEN ?::date AND ?::date")
			args = append(args, o.DateBefore, o.DateAfter)
		} else {
			whereClause = append(whereClause, "s.month = TO_CHAR(CURRENT_DATE, 'YYYY-MM')")
		}
	} else {
		whereClause = append(whereClause, "s.month = TO_CHAR(CURRENT_DATE, 'YYYY-MM')")
	}

	where := ""
	if len(whereClause) > 0 {
		where = "WHERE " + strings.Join(whereClause, " AND ")
	}

	orderBy := "ORDER BY s.month DESC, s.continent ASC, s.country ASC, s.operator ASC, s.channel ASC"
	if o.OrderColumn != "" {
		dir := "ASC"
		if strings.ToUpper(o.OrderDir) == "DESC" {
			dir = "DESC"
		}
		orderBy = fmt.Sprintf("ORDER BY %s %s", o.OrderColumn, dir)
	}

	rawSQL := fmt.Sprintf(`
		SELECT
			s.country, s.continent, s.company, s.partner, s.operator, s.channel,
			s.campaign_type, s.month,
			MAX(s.summary_date)::text AS last_date,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 1  AND 7  THEN s.actual_cost ELSE 0 END) AS actual_week_1,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 8  AND 14 THEN s.actual_cost ELSE 0 END) AS actual_week_2,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 15 AND 21 THEN s.actual_cost ELSE 0 END) AS actual_week_3,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 22 AND 28 THEN s.actual_cost ELSE 0 END) AS actual_week_4,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) >= 29             THEN s.actual_cost ELSE 0 END) AS actual_week_5,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 1  AND 7  THEN s.mo_count ELSE 0 END) AS mo_week1,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 8  AND 14 THEN s.mo_count ELSE 0 END) AS mo_week2,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 15 AND 21 THEN s.mo_count ELSE 0 END) AS mo_week3,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) BETWEEN 22 AND 28 THEN s.mo_count ELSE 0 END) AS mo_week4,
			SUM(CASE WHEN EXTRACT(DAY FROM s.summary_date) >= 29             THEN s.mo_count ELSE 0 END) AS mo_week5,
			COALESCE(b.id, 0)                   AS budget_io_id,
			COALESCE(b.monthly_spend_target, 0) AS io_target,
			COALESCE(b.monthly_mo_target, 0)    AS mo_target,
			COALESCE(b.target_cac, 0)           AS target_cac,
			0                                   AS ltv,
			0                                   AS roas,
			COALESCE(b.target_roi, 0)           AS roi
		FROM summary_budget_ios s
		LEFT JOIN budget_ios b ON
			b.country = s.country AND b.month = s.month
			AND b.deleted_at IS NULL
		%s
		GROUP BY
			s.country, s.continent, s.company, s.partner, s.operator, s.channel,
			s.campaign_type, s.month,
			b.id, b.monthly_spend_target, b.monthly_mo_target, b.target_cac, b.target_roi
		%s
	`, where, orderBy)

	var ss []entity.SummaryBudgetIOAgg
	if err := r.DB.Raw(rawSQL, args...).Scan(&ss).Error; err != nil {
		return nil, 0, err
	}
	return ss, int64(len(ss)), nil
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
