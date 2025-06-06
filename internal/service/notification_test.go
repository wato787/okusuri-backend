package service

import (
	"testing"

	"okusuri-backend/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestNotificationService_New(t *testing.T) {
	t.Run("NotificationServiceが正常に作成される", func(t *testing.T) {
		// NotificationServiceの作成をテスト
		service := NewNotificationService()

		assert.NotNil(t, service)
		assert.NotNil(t, service.recentSends)
	})
}

func TestNotificationService_RecentSendCheck(t *testing.T) {
	service := NewNotificationService()

	t.Run("初回送信は重複ではない", func(t *testing.T) {
		// 初回送信は重複とみなされない
		isRecent := service.isRecentlySent("test-subscription-key")
		assert.False(t, isRecent)
	})

	t.Run("送信記録をマークできる", func(t *testing.T) {
		// 送信記録をマーク
		service.markAsSent("test-subscription-key")

		// 直後に確認すると重複とみなされる
		isRecent := service.isRecentlySent("test-subscription-key")
		assert.True(t, isRecent)
	})
}

func TestNotificationService_SendNotification(t *testing.T) {
	service := NewNotificationService()

	t.Run("空のサブスクリプションでエラー", func(t *testing.T) {
		user := model.User{ID: "test-user"}
		setting := model.NotificationSetting{
			UserID:       "test-user",
			Platform:     "web",
			IsEnabled:    true,
			Subscription: "", // 空のサブスクリプション
		}

		err := service.SendNotification(user, setting, "テストメッセージ")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "サブスクリプションが見つかりません")
	})
}

func TestNotificationService_SendNotificationWithDays(t *testing.T) {
	service := NewNotificationService()

	t.Run("空のサブスクリプションでエラー", func(t *testing.T) {
		user := model.User{ID: "test-user"}
		setting := model.NotificationSetting{
			UserID:       "test-user",
			Platform:     "web",
			IsEnabled:    true,
			Subscription: "", // 空のサブスクリプション
		}

		err := service.SendNotificationWithDays(user, setting, "テストメッセージ", 5)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "サブスクリプションが見つかりません")
	})
}
