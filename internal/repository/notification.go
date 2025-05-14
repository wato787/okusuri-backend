package repository

import (
	"context"
	"fmt"
	"okusuri-backend/internal/model"
	"okusuri-backend/pkg/config"

	"firebase.google.com/go/v4/messaging"
)

type NotificationRepository struct{}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

// GetSettingByUserID はユーザーIDに基づいて通知設定を取得する
func (r *NotificationRepository) GetSettingByUserID(userID string) (*model.NotificationSetting, error) {
	// DB接続
	db := config.DB

	// ユーザーIDに基づいて通知設定を取得
	var setting model.NotificationSetting
	if err := db.Where("user_id = ?", userID).First(&setting).Error; err != nil {
		return nil, err
	}

	return &setting, nil
}

// RegisterSetting は通知設定を登録する
func (r *NotificationRepository) RegisterSetting(setting *model.NotificationSetting) error {
	// DB接続
	db := config.DB

	if err := db.Create(setting).Error; err != nil {
		return err
	}

	return nil
}

// GetAllSettings は全ての通知設定を取得する
func (r *NotificationRepository) GetAllSettings() ([]model.NotificationSetting, error) {
	// DB接続
	db := config.DB

	var settings []model.NotificationSetting
	if err := db.Find(&settings).Error; err != nil {
		return nil, err
	}

	return settings, nil
}

// SendNotification は通知を送信する
func (r *NotificationRepository) SendNotification(user model.User, setting model.NotificationSetting, message string) error {
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
