package service

import (
	"fmt"
	"okusuri-backend/model"
	"okusuri-backend/repository"
)

type MedicationLogService struct {
	MedicationLogRepository *repository.MedicationLogRepository
}

// NewMedicationLogService は新しいMedicationLogServiceのインスタンスを作成する
func NewMedicationLogService(medicationLogRepository *repository.MedicationLogRepository) *MedicationLogService {
	return &MedicationLogService{
		MedicationLogRepository: medicationLogRepository,
	}
}
func (s *MedicationLogService) RegisterMedicationLog(userID string, medicationLog model.MedicationLog) error {
	// 服用履歴を登録
	err := s.MedicationLogRepository.RegisterMedicationLog(userID, medicationLog)
	if err != nil {
		return fmt.Errorf("failed to register medication log: %w", err)
	}

	return nil
}
