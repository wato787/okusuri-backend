package medication

import (
	"fmt"
)

// Service は薬の服用記録に関するビジネスロジックを提供する
type Service struct {
	repository Repository
}

// NewService は新しいServiceインスタンスを作成する
func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

// RegisterLog は服用記録を登録する
func (s *Service) RegisterLog(userID string, log MedicationLog) error {
	// 服用履歴を登録
	err := s.repository.RegisterLog(userID, log)
	if err != nil {
		return fmt.Errorf("failed to register medication log: %w", err)
	}

	return nil
}

// GetLogsByUserID はユーザーIDに基づいて服用履歴を取得する
func (s *Service) GetLogsByUserID(userID string) ([]MedicationLog, error) {
	// ユーザーIDに基づいて服用履歴を取得
	logs, err := s.repository.GetLogsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get medication logs: %w", err)
	}

	return logs, nil
}
