package model

import "gorm.io/gorm"

// 服用履歴の構造体
type MedicationLog struct {
	gorm.Model
	UserID      string `json:"user_id" gorm:"not null;uniqueIndex:idx_user_id"`
	HasBleeding bool   `json:"has_bleeding" gorm:"default:false"`
}
