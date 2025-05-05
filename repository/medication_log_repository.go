package repository

import (
	"okusuri-backend/config"
	"okusuri-backend/model"
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
