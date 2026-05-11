package persistence

import (
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

// CreateMenu inserts a new menu into the database
func (r *BaseModel) CreateMenu(menu *entity.Menu) error {
	return r.DB.Create(menu).Error
}

// GetAllMenus retrieves all menus
func (r *BaseModel) GetAllMenus() ([]entity.Menu, error) {
	var menus []entity.Menu
	err := r.DB.Find(&menus).Error
	return menus, err
}

// GetMenuByID retrieves a menu by ID
func (r *BaseModel) GetMenuByID(id uint) (*entity.Menu, error) {
	var menu entity.Menu
	err := r.DB.First(&menu, id).Error
	return &menu, err
}

// UpdateMenu updates an existing menu
func (r *BaseModel) UpdateMenu(menu *entity.Menu) error {
	return r.DB.Save(menu).Error
}

// DeleteMenu deletes a menu by ID
func (r *BaseModel) DeleteMenu(id uint) error {
	return r.DB.Delete(&entity.Menu{}, id).Error
}
