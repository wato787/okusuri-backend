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

// 服薬状態を管理するモデル
type MedicationStatus struct {
	ID                      uint       `json:"id" gorm:"primarykey"`
	UserID                  string     `json:"userId" gorm:"uniqueIndex;not null"`
	IsRestPeriod            bool       `json:"isRestPeriod" gorm:"default:false"`
	CurrentStreak           int        `json:"currentStreak" gorm:"default:0"`           // 現在の服用日数
	ConsecutiveBleedingDays int        `json:"consecutiveBleedingDays" gorm:"default:0"` // 連続出血日数
	RestPeriodStartedAt     *time.Time `json:"restPeriodStartedAt"`                      // 休薬期間の開始日
	RestPeriodDays          int        `json:"restPeriodDays" gorm:"default:4"`          // 休薬期間の日数（デフォルト4日）
	CreatedAt               time.Time  `json:"createdAt"`
	UpdatedAt               time.Time  `json:"updatedAt"`
}
