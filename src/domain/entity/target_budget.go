package entity

import (
	"time"

	"gorm.io/gorm"
)

type (
	TargetBudget struct {
		gorm.Model
		ID       int     `gorm:"primaryKey;autoIncrement" json:"id"`
		Country  string  `gorm:"uniqueIndex:idx_target_budget;not null;size:50" json:"country"`
		Year     int     `gorm:"uniqueIndex:idx_target_budget;not null;size:50" json:"year"`
		Month    int     `gorm:"uniqueIndex:idx_target_budget;not null;size:50" json:"month"`
		Budget   float64 `gorm:"type:double precision;default:0" json:"budget"`
		Spending float64 `gorm:"type:double precision;default:0" json:"spending"`

		CreatedAt time.Time
		UpdatedAt time.Time
	}

	TargetBudgetDetail struct {
		gorm.Model
		ID           int     `gorm:"primaryKey;autoIncrement" json:"id"`
		Country      string  `gorm:"uniqueIndex:idx_target_budget_detail;not null;size:50" json:"country"`
		Year         int     `gorm:"uniqueIndex:idx_target_budget_detail;not null;size:50" json:"year"`
		Month        int     `gorm:"uniqueIndex:idx_target_budget_detail;not null;size:50" json:"month"`
		Operator     string  `gorm:"uniqueIndex:idx_target_budget_detail;not null;size:50" json:"operator"`
		Partner      string  `gorm:"uniqueIndex:idx_target_budget_detail;not null;size:50" json:"partner"`
		Service      string  `gorm:"uniqueIndex:idx_target_budget_detail;not null;size:50" json:"service"`
		Adnet        string  `gorm:"uniqueIndex:idx_target_budget_detail;not null;size:50" json:"adnet"`
		Budget       float64 `gorm:"type:double precision;default:0" json:"budget"`
		BudgetPerDay float64 `gorm:"type:double precision;default:0" json:"budget_per_day"`
		Spending     float64 `gorm:"type:double precision;default:0" json:"spending"`
		BudgetUsage  float64 `gorm:"type:double precision;default:0" json:"budget_usage"`

		CreatedAt time.Time
		UpdatedAt time.Time
	}
)
