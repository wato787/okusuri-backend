// internal/service/notification.go
package service

import (
	"context"
	"fmt"
	"okusuri-backend/internal/model"
	"okusuri-backend/pkg/config"
	"sync"
	"time"

	"firebase.google.com/go/v4/messaging"
)

type NotificationService struct {
	// 直近に送信したトークンとタイムスタンプを保持するマップ
	// キー: FCMトークン, 値: 最後の送信時刻
	recentSends     map[string]time.Time
	recentSendMutex sync.Mutex
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		recentSends: make(map[string]time.Time),
	}
}

// 最近送信した通知かどうかをチェック（5分以内に同じトークンに送信したか）
func (s *NotificationService) isRecentlySent(token string) bool {
	s.recentSendMutex.Lock()
	defer s.recentSendMutex.Unlock()

	lastSent, exists := s.recentSends[token]
	if !exists {
		return false
	}

	// 5分以内の送信なら重複とみなす
	timeSinceLast := time.Since(lastSent)
	fmt.Printf(">> 前回の送信からの経過時間: %v (トークン: %s...)\n", 
		timeSinceLast.Round(time.Second), token[:10])
	return timeSinceLast < 5*time.Minute
}

// 送信記録を更新
func (s *NotificationService) markAsSent(token string) {
	s.recentSendMutex.Lock()
	defer s.recentSendMutex.Unlock()

	s.recentSends[token] = time.Now()
	fmt.Printf(">> トークン %s... を送信済みとしてマークしました\n", token[:10])

	// 古い記録をクリーンアップ（1時間以上前のものを削除）
	for t, lastSent := range s.recentSends {
		if time.Since(lastSent) > time.Hour {
			delete(s.recentSends, t)
			fmt.Printf(">> 古い送信記録を削除: %s...\n", t[:10])
		}
	}
}

// SendNotification は通知を送信する
func (s *NotificationService) SendNotification(user model.User, setting model.NotificationSetting, message string) error {
	// fcmTokenが空でない場合、通知を送信する
	if setting.FcmToken != "" {
		tokenPreview := setting.FcmToken
		if len(tokenPreview) > 10 {
			tokenPreview = tokenPreview[:10] + "..."
		}
		
		fmt.Printf("\n>> 通知送信サービス: ユーザーID: %s の処理を開始します\n", user.ID)
		fmt.Printf(">> FCMトークン: %s\n", tokenPreview)
		
		// 最近送信済みなら重複送信をスキップ
		if s.isRecentlySent(setting.FcmToken) {
			fmt.Printf(">> 通知送信サービス: トークン %s は最近送信済みのためスキップします\n", 
				tokenPreview)
			return nil // エラーにせず成功扱いでスキップ
		}

		ctx := context.Background()

		// 初期化済みのFCMクライアントを取得
		client, err := config.GetMessagingClient(ctx)
		if err != nil {
			fmt.Printf(">> 通知送信サービス: FCMクライアント取得エラー: %v\n", err)
			return fmt.Errorf("FCMクライアント取得エラー: %v", err)
		}

		// 通知メッセージの作成
		messageID := fmt.Sprintf("medication-%d", time.Now().UnixNano())
		fmt.Printf(">> 通知送信サービス: メッセージID: %s を作成\n", messageID)
		
		msg := &messaging.Message{
			Notification: &messaging.Notification{
				Title: "通知",
				Body:  message,
			},
			// 重複排除のためにメッセージIDを明示的に設定
			Data: map[string]string{
				"messageId": messageID,
				"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
				"userId":    user.ID,
			},
			Token: setting.FcmToken,
		}

		// 送信前にログ
		fmt.Printf(">> 通知送信サービス: FCMに送信リクエスト実行...\n")
		startTime := time.Now()
		
		// 通知の送信
		sendResult, err := client.Send(ctx, msg)
		elapsedTime := time.Since(startTime)
		
		if err != nil {
			fmt.Printf(">> 通知送信サービス: 通知送信エラー: %v (所要時間: %v)\n", 
				err, elapsedTime.Round(time.Millisecond))
			return fmt.Errorf("通知送信エラー: %v", err)
		}

		// 送信済みとしてマーク
		s.markAsSent(setting.FcmToken)
		
		fmt.Printf(">> 通知送信サービス: 通知送信成功 - FCM MessageID: %s (所要時間: %v)\n", 
			sendResult, elapsedTime.Round(time.Millisecond))
		fmt.Printf(">> 通知送信サービス: ユーザーID %s の処理完了\n", user.ID)
	} else {
		fmt.Printf(">> 通知送信サービス: ユーザーID: %s のFCMトークンが空です\n", user.ID)
	}

	return nil
}
