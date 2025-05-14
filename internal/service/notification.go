// internal/service/notification_service.go
package service

import (
	"context"
	"fmt"
	"okusuri-backend/internal/model"
	"okusuri-backend/pkg/config"

	"firebase.google.com/go/v4/messaging"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// SendNotification は通知を送信する
func (s *NotificationService) SendNotification(user model.User, setting model.NotificationSetting, message string) error {
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
