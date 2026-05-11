package persistence

import (
	"errors"

	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

// CreateRole inserts a new role into the database
func (r *BaseModel) CreateRole(role *entity.Role) error {
    return r.DB.Create(role).Error
}

// CreateOrUpdatePermission creates or updates permission data
func (r *BaseModel) CreateOrUpdatePermission(permission *entity.Permission) error {
    var existing entity.Permission
    result := r.DB.Where("role_id = ? AND menu_id = ?", permission.RoleID, permission.MenuID).First(&existing)

    if result.RowsAffected > 0 {
        return r.DB.Model(&existing).Save(permission).Error
    }
    
    return r.DB.Create(permission).Error
}

// GetAllRoles retrieves all roles
func (r *BaseModel) GetAllRoles() ([]entity.Role, error) {
	var roles []entity.Role
	err := r.DB.Find(&roles).Error
	return roles, err
}

// GetRoleByID retrieves a role by its ID
func (r *BaseModel) GetRoleByID(id int) (*entity.Role, error) {
	var role entity.Role
	if err := r.DB.Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByCode retrieves a role by its Code
func (r *BaseModel) GetRoleByCode(code string) (*entity.Role, error) {
	var role entity.Role
	if err := r.DB.Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// UpdateRole updates an existing role
func (r *BaseModel) UpdateRole(role *entity.Role) error {
	if role.ID == 0 {
		return errors.New("role ID is required")
	}
	return r.DB.Model(&entity.Role{}).Where("id = ?", role.ID).Save(role).Error
}

// DeletePermissionsByRoleID deletes all permissions linked to a role
func (r *BaseModel) DeletePermissionsByRoleID(roleID int) error {
	return r.DB.Where("role_id = ?", roleID).Delete(&entity.Permission{}).Error
}

// DeleteRole deletes a role by ID
func (r *BaseModel) DeleteRole(roleID int) error {
	result := r.DB.Where("id = ?", roleID).Delete(&entity.Role{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("role not found")
	}
	return nil
}

func (r *BaseModel) GetAllRolesWithPermission() ([]entity.RoleManagementData, error) {
	var roles []entity.RoleManagementData

	rows, err := r.DB.Table("roles").
		Select("roles.id, roles.code, roles.name, menus.name as permission_name").
		Joins("INNER JOIN permissions ON roles.id = permissions.role_id").
		Joins("INNER JOIN menus ON permissions.menu_id = menus.id").
		Where("permissions.status = ?", true).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roleMap := make(map[int]*entity.RoleManagementData)

	for rows.Next() {
		var roleID int
		var roleCode, roleName, permissionName string

		if err := rows.Scan(&roleID, &roleCode, &roleName, &permissionName); err != nil {
			return nil, err
		}

		if _, exists := roleMap[roleID]; !exists {
			roleMap[roleID] = &entity.RoleManagementData{
				ID:          roleID,
				Code:        roleCode,
				Name:        roleName,
				Permissions: []string{},
			}
		}

		if permissionName != "" {
			roleMap[roleID].Permissions = append(roleMap[roleID].Permissions, permissionName)
		}
	}

	for _, role := range roleMap {
		roles = append(roles, *role)
	}

	return roles, nil
}
