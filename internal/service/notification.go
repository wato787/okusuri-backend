package service

import (
	"encoding/json"
	"fmt"
	"okusuri-backend/internal/model"
	"os"
	"sync"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"
)

// NotificationService は通知を送信するサービス
type NotificationService struct {
	// 直近に送信したサブスクリプションとタイムスタンプを保持するマップ
	recentSends     map[string]time.Time
	recentSendMutex sync.Mutex
}

// サブスクリプションデータの構造体
type PushSubscription struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

// 通知データの構造体
type NotificationData struct {
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Data  map[string]string `json:"data,omitempty"`
}

// 新しいNotificationServiceのインスタンスを作成
func NewNotificationService() *NotificationService {
	return &NotificationService{
		recentSends: make(map[string]time.Time),
	}
}

// 最近送信した通知かどうかをチェック（5分以内に同じサブスクリプションに送信したか）
func (s *NotificationService) isRecentlySent(subKey string) bool {
	s.recentSendMutex.Lock()
	defer s.recentSendMutex.Unlock()

	lastSent, exists := s.recentSends[subKey]
	if !exists {
		return false
	}

	// 5分以内の送信なら重複とみなす
	timeSinceLast := time.Since(lastSent)
	fmt.Printf(">> 前回の送信からの経過時間: %v (サブスクリプション: %s...)\n",
		timeSinceLast.Round(time.Second), subKey[:10])
	return timeSinceLast < 5*time.Minute
}

// 送信記録を更新
func (s *NotificationService) markAsSent(subKey string) {
	s.recentSendMutex.Lock()
	defer s.recentSendMutex.Unlock()

	s.recentSends[subKey] = time.Now()
	fmt.Printf(">> サブスクリプション %s... を送信済みとしてマークしました\n", subKey[:10])

	// 古い記録をクリーンアップ（1時間以上前のものを削除）
	for key, lastSent := range s.recentSends {
		if time.Since(lastSent) > time.Hour {
			delete(s.recentSends, key)
			fmt.Printf(">> 古い送信記録を削除: %s...\n", key[:10])
		}
	}
}

// SendNotificationWithDays は連続服薬日数を含めて通知を送信する
func (s *NotificationService) SendNotificationWithDays(
	user model.User, setting model.NotificationSetting, message string, consecutiveDays int,
) error {
	// subscriptionが空の場合
	if setting.Subscription == "" {
		fmt.Printf(">> 通知サービス: ユーザーID: %s のサブスクリプションが空です\n", user.ID)
		return fmt.Errorf("サブスクリプションが見つかりません")
	}

	subscriptionPreview := setting.Subscription
	if len(subscriptionPreview) > 10 {
		subscriptionPreview = subscriptionPreview[:10] + "..."
	}

	fmt.Printf("\n>> 通知サービス: ユーザーID: %s の処理を開始します\n", user.ID)
	fmt.Printf(">> サブスクリプション: %s\n", subscriptionPreview)

	// JSON文字列をパース
	var subscription PushSubscription
	err := json.Unmarshal([]byte(setting.Subscription), &subscription)
	if err != nil {
		fmt.Printf(">> 通知サービス: サブスクリプションのパースに失敗: %v\n", err)
		return fmt.Errorf("サブスクリプションのパースに失敗: %v", err)
	}

	// 最近送信済みなら重複送信をスキップ
	subKey := subscription.Endpoint
	if s.isRecentlySent(subKey) {
		fmt.Printf(">> 通知サービス: サブスクリプション %s は最近送信済みのためスキップします\n",
			subscriptionPreview)
		return nil // エラーにせず成功扱いでスキップ
	}

	// VAPID鍵の取得
	vapidPublicKey := os.Getenv("VAPID_PUBLIC_KEY")
	vapidPrivateKey := os.Getenv("VAPID_PRIVATE_KEY")

	if vapidPublicKey == "" || vapidPrivateKey == "" {
		fmt.Printf(">> 通知サービス: VAPID鍵が設定されていません\n")
		return fmt.Errorf("VAPID鍵が設定されていません")
	}

	// 通知内容の作成（連続服薬日数を含める）
	notificationData := NotificationData{
		Title: "お薬通知",
		Body:  message,
		Data: map[string]string{
			"messageId":       fmt.Sprintf("medication-%d", time.Now().UnixNano()),
			"timestamp":       fmt.Sprintf("%d", time.Now().Unix()),
			"userId":          user.ID,
			"consecutiveDays": fmt.Sprintf("%d", consecutiveDays),
		},
	}

	// 通知内容をJSONに変換
	payload, err := json.Marshal(notificationData)
	if err != nil {
		fmt.Printf(">> 通知サービス: 通知内容のJSON変換に失敗: %v\n", err)
		return fmt.Errorf("通知内容のJSON変換に失敗: %v", err)
	}

	// Web Push通知の送信
	_, err = webpush.SendNotification(
		payload,
		&webpush.Subscription{
			Endpoint: subscription.Endpoint,
			Keys: webpush.Keys{
				P256dh: subscription.Keys.P256dh,
				Auth:   subscription.Keys.Auth,
			},
		},
		&webpush.Options{
			VAPIDPublicKey:  vapidPublicKey,
			VAPIDPrivateKey: vapidPrivateKey,
			TTL:             30,
			Subscriber:      "example@example.com", // 開発者のメールアドレス
		},
	)

	if err != nil {
		fmt.Printf(">> 通知サービス: 通知送信エラー: %v\n", err)
		return fmt.Errorf("通知送信エラー: %v", err)
	}

	// 送信済みとしてマーク
	s.markAsSent(subKey)

	fmt.Printf(">> 通知サービス: 通知送信成功\n")
	fmt.Printf(">> 通知サービス: ユーザーID %s の処理完了\n", user.ID)

	return nil
}

// SendNotification は通知を送信する（後方互換性のため）
func (s *NotificationService) SendNotification(
	user model.User, setting model.NotificationSetting, message string,
) error {
	return s.SendNotificationWithDays(user, setting, message, 0)
}
