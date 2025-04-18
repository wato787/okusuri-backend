package service

import (
	"context"
	"fmt"
	"okusuri-backend/config"
	"okusuri-backend/model"
	"okusuri-backend/repository"

	"firebase.google.com/go/v4/messaging"
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

func (s *NotificationService) RegisterNotificationSetting(notificationSetting *model.NotificationSetting) error {
	err := s.NotificationRepository.RegisterNotificationSetting(notificationSetting)
	if err != nil {
		return err
	}

	return nil
}

func (s *NotificationService) GetAllNotificationSettings() ([]model.NotificationSetting, error) {
	notificationSettings, err := s.NotificationRepository.GetAllNotificationSettings()
	if err != nil {
		return nil, err
	}

	return notificationSettings, nil
}

func (s *NotificationService) SendNotification(user model.User, notificationSetting model.NotificationSetting, message string) error {

	// fcmTokenが空でない場合、通知を送信する
	if notificationSetting.FcmToken != "" {
		ctx := context.Background()

		// 初期化済みのFCMクライアントを取得
		client, err := config.GetMessagingClient(ctx)
		if err != nil {
			return fmt.Errorf("FCMクライアント取得エラー: %v", err)
		}

		// 通知メッセージの作成
		msg := &messaging.Message{
			Notification: &messaging.Notification{
				Title: "通知",
				Body:  message,
			},
			Token: notificationSetting.FcmToken,
		}

		// 通知の送信
		_, err = client.Send(ctx, msg)
		if err != nil {
			return fmt.Errorf("通知送信エラー: %v", err)
		}
	}

	return nil
}
