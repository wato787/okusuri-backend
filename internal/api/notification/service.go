package notification

import (
	"context"
	"fmt"
	"okusuri-backend/internal/common/user"
	"okusuri-backend/pkg/config"

	"firebase.google.com/go/v4/messaging"
)

// Service は通知関連のビジネスロジックを提供する
type Service struct {
	repository *Repository
}

// NewService は新しいServiceインスタンスを作成する
func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

// GetSettingByUserID はユーザーIDに基づいて通知設定を取得する
func (s *Service) GetSettingByUserID(userID string) (*NotificationSetting, error) {
	// リポジトリから通知設定を取得
	setting, err := s.repository.GetSettingByUserID(userID)
	if err != nil {
		return nil, err
	}

	return setting, nil
}

// RegisterSetting は通知設定を登録する
func (s *Service) RegisterSetting(setting *NotificationSetting) error {
	err := s.repository.RegisterSetting(setting)
	if err != nil {
		return err
	}

	return nil
}

// GetAllSettings は全ての通知設定を取得する
func (s *Service) GetAllSettings() ([]NotificationSetting, error) {
	settings, err := s.repository.GetAllSettings()
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// SendNotification は通知を送信する
func (s *Service) SendNotification(user user.User, setting NotificationSetting, message string) error {
	// fcmTokenが空でない場合、通知を送信する
	if setting.FcmToken != "" {
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
			Token: setting.FcmToken,
		}

		// 通知の送信
		_, err = client.Send(ctx, msg)
		if err != nil {
			return fmt.Errorf("通知送信エラー: %v", err)
		}
	}

	return nil
}
