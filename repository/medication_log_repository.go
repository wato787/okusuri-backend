package repository

import (
	"okusuri-backend/model"
	"okusuri-backend/pkg/config"
)

type MedicationLogRepository struct{}

func NewMedicationLogRepository() *MedicationLogRepository {
	return &MedicationLogRepository{}
}
func (r *MedicationLogRepository) RegisterMedicationLog(userID string, medicationLog model.MedicationLog) error {
	// DB接続
	db := config.DB

	// ユーザーIDに基づいて服用履歴を登録
	if err := db.Create(&medicationLog).Error; err != nil {
		return err
	}

	return nil
}

func (r *MedicationLogRepository) GetMedicationLogsByUserID(userID string) ([]model.MedicationLog, error) {
	// DB接続
	db := config.DB

	// ユーザーIDに基づいて服用履歴を取得
	var medicationLogs []model.MedicationLog
	if err := db.Where("user_id = ?", userID).Find(&medicationLogs).Error; err != nil {
		return nil, err
	}

	return medicationLogs, nil
}
