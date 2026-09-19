package persistence

import (
	"database/sql"
	"strings"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func (m *BaseModel) CreateOperatorAlias(o *entity.OperatorAlias) error {
	return m.DB.Create(o).Error
}

func (m *BaseModel) UpdateOperatorAlias(o *entity.OperatorAlias) error {
	return m.DB.Updates(o).Error
}

func (m *BaseModel) DeleteOperatorAlias(id uint) error {
	return m.DB.Delete(&entity.OperatorAlias{}, id).Error
}

func (r *BaseModel) GetOperatorAliasList(o entity.GlobalRequestFromDataTable) ([]entity.OperatorAlias, int64, error) {

	var (
		rows       *sql.Rows
		total_rows int64
	)

	query := r.DB.Model(&entity.OperatorAlias{}).Where("type = ?", "API")
	if o.Search != "" {
		search_value := strings.Trim(o.Search, " ")
		query = query.Where("operator ILIKE ? OR alias ILIKE ? OR service ILIKE ? OR country ILIKE ?",
			"%"+search_value+"%", "%"+search_value+"%", "%"+search_value+"%", "%"+search_value+"%")
	}

	query.Unscoped().Count(&total_rows)

	query_limit := query.Limit(o.PageSize)
	if o.Page > 0 {
		query_limit = query_limit.Offset((o.Page - 1) * o.PageSize)
	}

	rows, _ = query_limit.Order("operator").Rows()
	defer rows.Close()

	var ss []entity.OperatorAlias
	for rows.Next() {
		var s entity.OperatorAlias
		r.DB.ScanRows(rows, &s)
		ss = append(ss, s)
	}

	return ss, total_rows, rows.Err()
}
