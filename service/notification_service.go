package service

import (
	"okusuri-backend/model"
	"okusuri-backend/repository"
)

type NotificationService struct {
	NotificationRepository *repository.NotificationRepository
}

// NewNotificationService は新しいNotificationServiceのインスタンスを作成する
func NewNotificationService(notificationRepository *repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		NotificationRepository: notificationRepository,
	}
}

func (s *NotificationService) GetNotificationSettingByUserID(userID string) (*model.NotificationSetting, error) {
	// ユーザーIDに基づいて通知設定を取得
	notificationSetting, err := s.NotificationRepository.GetNotificationSettingByUserID(userID)
	if err != nil {
		return nil, err
	}

	return notificationSetting, nil
}
