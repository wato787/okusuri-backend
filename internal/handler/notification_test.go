package handler

import (
	"encoding/json"
	"net/http"
	"okusuri-backend/internal/dto"
	"okusuri-backend/internal/model"
	"okusuri-backend/internal/repository"
	"okusuri-backend/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNotificationHandler_GetSetting は通知設定取得のテスト（構造体の基本テスト）
func TestNotificationHandler_GetSetting(t *testing.T) {
	t.Run("ユーザーIDなしのエラーハンドリング", func(t *testing.T) {
		// ハンドラーの作成
		handler := NewNotificationHandler(
			repository.NewNotificationRepository(),
			repository.NewUserRepository(),
			service.NewNotificationService(),
			repository.NewMedicationRepository(),
			service.NewMedicationService(repository.NewMedicationRepository()),
		)

		// リクエストの作成（ユーザーIDなし）
		c, w := createTestContext("GET", "/api/notification/setting", nil, "")

		// ハンドラーの実行
		handler.GetSetting(c)

		// レスポンスの検証
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid user ID", response["error"])
	})

	t.Run("NotificationSettingの構造テスト", func(t *testing.T) {
		// NotificationSettingの基本構造をテスト
		setting := model.NotificationSetting{
			ID:           1,
			UserID:       "test-user-1",
			Subscription: "test-subscription",
			IsEnabled:    true,
			UpdatedAt:    time.Now(),
		}

		assert.NotZero(t, setting.ID)
		assert.Equal(t, "test-user-1", setting.UserID)
		assert.Equal(t, "test-subscription", setting.Subscription)
		assert.True(t, setting.IsEnabled)
		assert.NotZero(t, setting.UpdatedAt)
	})
}

// TestNotificationHandler_RegisterSetting は通知設定登録のテスト（構造体の基本テスト）
func TestNotificationHandler_RegisterSetting(t *testing.T) {
	t.Run("リクエストボディの構造テスト", func(t *testing.T) {
		// リクエストボディの基本構造をテスト
		req := dto.RegisterNotificationSettingRequest{
			Subscription: "test-subscription-data",
			IsEnabled:    true,
			Platform:     "web",
		}

		assert.Equal(t, "test-subscription-data", req.Subscription)
		assert.True(t, req.IsEnabled)
	})

	t.Run("ユーザーIDなしのエラーハンドリング", func(t *testing.T) {
		// ハンドラーの作成
		handler := NewNotificationHandler(
			repository.NewNotificationRepository(),
			repository.NewUserRepository(),
			service.NewNotificationService(),
			repository.NewMedicationRepository(),
			service.NewMedicationService(repository.NewMedicationRepository()),
		)

		// リクエストの作成（ユーザーIDなし）
		requestBody := dto.RegisterNotificationSettingRequest{
			Subscription: "test",
			IsEnabled:    true,
			Platform:     "web",
		}
		c, w := createTestContext("POST", "/api/notification/setting", requestBody, "")

		// ハンドラーの実行
		handler.RegisterSetting(c)

		// レスポンスの検証
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid user ID", response["error"])
	})

	t.Run("リクエストバリデーション構造テスト", func(t *testing.T) {
		// リクエストバリデーションの基本構造をテスト
		requestBody := dto.RegisterNotificationSettingRequest{
			Subscription: "",
			IsEnabled:    true,
			Platform:     "web",
		}

		assert.Empty(t, requestBody.Subscription)
		assert.True(t, requestBody.IsEnabled)
		assert.Equal(t, "web", requestBody.Platform)
	})
}

// TestNotificationHandler_New は NotificationHandler 作成のテスト
func TestNotificationHandler_New(t *testing.T) {
	t.Run("NotificationHandlerが正常に作成される", func(t *testing.T) {
		notificationRepo := repository.NewNotificationRepository()
		userRepo := repository.NewUserRepository()
		notificationSvc := service.NewNotificationService()
		medicationRepo := repository.NewMedicationRepository()

		medicationSvc := service.NewMedicationService(medicationRepo)

		handler := NewNotificationHandler(
			notificationRepo,
			userRepo,
			notificationSvc,
			medicationRepo,
			medicationSvc,
		)

		assert.NotNil(t, handler)
		assert.Equal(t, notificationRepo, handler.notificationRepo)
		assert.Equal(t, userRepo, handler.userRepo)
		assert.Equal(t, notificationSvc, handler.notificationSvc)
		assert.Equal(t, medicationSvc, handler.medicationSvc)
	})
}
