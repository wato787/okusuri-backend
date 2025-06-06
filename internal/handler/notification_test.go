package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// NotificationHandlerの基本テスト
func TestNotificationHandler_New(t *testing.T) {
	t.Run("NotificationHandlerが正常に作成される", func(t *testing.T) {
		// NotificationHandlerの作成をテスト
		// 実際の依存関係は使わずに、nilで作成してもパニックしないことを確認
		assert.NotPanics(t, func() {
			NewNotificationHandler(nil, nil, nil, nil, nil)
		})
	})
}
