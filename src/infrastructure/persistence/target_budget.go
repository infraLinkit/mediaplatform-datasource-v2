package persistence

import (
	"fmt"
	"time"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func (r *BaseModel) GetTargetBudget(country string, startDate time.Time, endDate time.Time, operator string, partner string, service string, adnet string) ([]entity.TargetBudget, bool) {

	start := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	nextMonth := time.Date(endDate.Year(), endDate.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	end := nextMonth.AddDate(0, 0, -1).Format("2006-01-02")

	query := r.DB.Table("target_budgets as a").
		Select(`a.country, 
            a.year, 
            a.month, 
            a.budget, 
            SUM(b.saaf) as spending`).
		Joins(`LEFT JOIN summary_campaigns as b ON 
           a.country = b.country AND 
           a.year = EXTRACT(year from b.summary_date) AND 
           a.month = EXTRACT(month from b.summary_date)`).
		Where("(a.year||'-'||a.month||'-'||'2')::date BETWEEN ? AND ? AND b.deleted_at IS NULL",
			start, end).
		Group("a.country, a.year, a.month, a.budget")

	if country != "" {
		query.Where("a.country=?", country)
	}

	if operator != "" {
		query.Where("b.operator=?", operator)
	}

	if partner != "" {
		query.Where("b.partner=?", partner)
	}

	if service != "" {
		query.Where("b.service=?", service)
	}

	if adnet != "" {
		query.Where("b.adnet=?", adnet)
	}

	if rows, err := query.Rows(); err == nil {
		var ss []entity.TargetBudget
		defer rows.Close()
		for rows.Next() {
			var s entity.TargetBudget
			err := rows.Scan(&s.Country, &s.Year, &s.Month, &s.Budget, &s.Spending)
			if err != nil {
				fmt.Println("ERROR REPORTED::", err)
			}
			fmt.Println("COUTRY:: ", s.Country)
			ss = append(ss, s)
		}

		return ss, true
	} else {

	}

	return []entity.TargetBudget{}, false
}

func (r *BaseModel) GetTargetBudgetList(country string, startDate time.Time, endDate time.Time, operator string, partner string, service string, adnet string) ([]entity.TargetBudgetDetail, bool) {
	where := " AND TRUE "
	if country != "" {
		where += " AND a.country='" + country + "'"
	}
	if operator != "" {
		where += " AND a.operator='" + operator + "'"
	}

	if partner != "" {
		where += " AND a.partner='" + partner + "'"
	}

	if service != "" {
		where += " AND a.service='" + service + "'"
	}

	if adnet != "" {
		where += " AND a.adnet='" + adnet + "'"
	}

	start := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	nextMonth := time.Date(endDate.Year(), endDate.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	end := nextMonth.AddDate(0, 0, -1).Format("2006-01-02")

	SQL := `WITH MonthlySpending AS (
			SELECT 
				EXTRACT(year from summary_date) as year,
				EXTRACT(month from summary_date) as month,
				country, 
				operator, 
				partner, 
				service, 
				adnet,
				SUM(saaf) as total_spending
			FROM summary_campaigns
			WHERE summary_date BETWEEN '` + start + ` 00:00:00' AND '` + end + ` 00:00:00' ` +
		where +
		` AND deleted_at IS NULL
			GROUP BY 1, 2, 3, 4, 5, 6, 7
		)
		SELECT 
			m.year,
			m.month,
			m.country, 
			m.operator, 
			m.partner, 
			m.service, 
			m.adnet,
			CASE WHEN b.budget IS NULL then 0 ELSE b.budget END as budget, 
			CASE WHEN b.budget_per_day IS NULL then 0 ELSE b.budget_per_day END as budget_per_day,
			m.total_spending as spending,
			CASE 
				WHEN b.budget > 0 THEN m.total_spending / b.budget 
				ELSE 0 
			END as budget_usage
		FROM MonthlySpending m
		LEFT JOIN target_budget_details b ON
			m.country = b.country AND
			m.operator = b.operator AND
			m.partner = b.partner AND
			m.service = b.service AND
			m.adnet = b.adnet AND
			m.year = b.year::integer AND
			m.month = b.month::integer;`

	query := r.DB.Raw(SQL)

	ss := []entity.TargetBudgetDetail{}

	if rows, err := query.Rows(); err == nil {
		defer rows.Close()
		for rows.Next() {
			var s entity.TargetBudgetDetail
			rows.Scan(&s.Year, &s.Month, &s.Country, &s.Operator,
				&s.Partner, &s.Service, &s.Adnet, &s.Budget, &s.BudgetPerDay, &s.Spending, &s.BudgetUsage)
			ss = append(ss, s)
		}
		return ss, false
	}

	fmt.Println(ss)
	return ss, true

}

func (r *BaseModel) AddTargetBudget(s entity.TargetBudget, data []entity.TargetBudgetDetail) error {
	SQL := `
	INSERT INTO target_budgets (country,month,year,budget,updated_at) 
	VALUES
	(?,?,?,?,NOW())
	ON CONFLICT (country,month,year) 
	DO UPDATE SET 
		updated_at=NOW(),
		budget = ? `

	q := r.DB.Exec(SQL, s.Country, s.Month, s.Year, s.Budget, s.Budget)

	for _, s := range data {
		SQL := `
			INSERT INTO target_budget_details(
			created_at,updated_at, country, year, month, operator, partner, service, adnet, 
			budget, budget_per_day
			) VALUES
			(NOW(),NOW(),?,?,?,?,?,?,?,?,0.0)
			ON CONFLICT (country, year, month, operator, partner, service, adnet) 
			DO UPDATE SET 
				updated_at=NOW(),
				budget = ? `
		q = r.DB.Exec(SQL, s.Country, s.Year, s.Month, s.Operator, s.Partner, s.Service, s.Adnet, s.Budget, s.Budget)
	}
	return q.Error
}
