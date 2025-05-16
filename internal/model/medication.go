package model

import "time"

// 服用履歴の構造体
type MedicationLog struct {
	ID          uint       `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty" gorm:"index"`
	UserID      string     `json:"userId" gorm:"not null;index:idx_user_id"` // uniqueIndexからindexに変更
	HasBleeding bool       `json:"hasBleeding" gorm:"default:false"`
}
