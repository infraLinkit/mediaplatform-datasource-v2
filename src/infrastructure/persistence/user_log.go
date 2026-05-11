package persistence

import (
	"strings"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

func (m *BaseModel) CreateUserLog(userlog *entity.Audit) error {
	return m.DB.Create(userlog).Error
}

func (r *BaseModel) GetUserLogList(o entity.GlobalRequestFromDataTable) ([]entity.DisplayUserLogList, int64, error) {
	var (
		logs      []entity.DisplayUserLogList
		totalRows int64
	)

	query := r.DB.Table("audits").
		Select("audits.*, users.email, users.name, users.username, roles.name as role, audits.updated_at as action_date").
		Joins("JOIN users ON audits.user_id = users.id").
		Joins("JOIN roles ON users.role_id = roles.id").
		Joins("JOIN (SELECT user_id, MAX(created_at) as latest_audit FROM audits GROUP BY user_id) as latest ON audits.user_id = latest.user_id AND audits.created_at = latest.latest_audit")

	// Apply search filter
	if o.Search != "" {
		searchValue := strings.TrimSpace(o.Search)
		query = query.Where("users.email ILIKE ? OR users.name ILIKE ? OR users.username ILIKE ?", "%"+searchValue+"%", "%"+searchValue+"%", "%"+searchValue+"%")
	}

	query.Unscoped().Count(&totalRows)

	query = query.Limit(o.PageSize)
	if o.Page > 0 {
		query = query.Offset((o.Page - 1) * o.PageSize)
	}

	if err := query.Order("audits.updated_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, totalRows, nil
}

func (r *BaseModel) GetUserLogHistory(o entity.GlobalRequestFromDataTable, ID int) ([]entity.DisplayUserLogList, int64, error) {
	var (
		logs      []entity.DisplayUserLogList
		totalRows int64
	)

	query := r.DB.Table("audits").
		Select("audits.*, audits.updated_at as action_date").
		Where("user_id = ?", ID).
		Group("audits.id, audits.updated_at").
		Order("audits.updated_at DESC")

	// Apply search filter
	if o.Search != "" {
		searchValue := strings.TrimSpace(o.Search)
		query = query.Where("users.email ILIKE ? OR users.name ILIKE ? OR users.username ILIKE ?", "%"+searchValue+"%", "%"+searchValue+"%", "%"+searchValue+"%")
	}

	query.Unscoped().Count(&totalRows)

	query = query.Limit(o.PageSize)
	if o.Page > 0 {
		query = query.Offset((o.Page - 1) * o.PageSize)
	}

	if err := query.Order("audits.updated_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, totalRows, nil
}
