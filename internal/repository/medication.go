package repository

import (
	"okusuri-backend/internal/model"
	"okusuri-backend/pkg/config"
)

type MedicationRepository struct{}

func NewMedicationRepository() *MedicationRepository {
	return &MedicationRepository{}
}

// RegisterLog はユーザーの服用記録をデータベースに登録する
func (r *MedicationRepository) RegisterLog(userID string, log model.MedicationLog) error {
	// DB接続
	db := config.DB

	// ユーザーIDに基づいて服用履歴を登録
	if err := db.Create(&log).Error; err != nil {
		return err
	}

	return nil
}

// GetLogsByUserID はユーザーIDに基づいて服用履歴をデータベースから取得する
func (r *MedicationRepository) GetLogsByUserID(userID string) ([]model.MedicationLog, error) {
	// DB接続
	db := config.DB

	// ユーザーIDに基づいて服用履歴を取得
	var logs []model.MedicationLog
	if err := db.Where("user_id = ?", userID).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}
