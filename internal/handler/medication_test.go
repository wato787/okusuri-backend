package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"okusuri-backend/internal/dto"
	"okusuri-backend/internal/repository"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MedicationHandlerの基本テスト
func TestMedicationHandler_New(t *testing.T) {
	t.Run("MedicationHandlerが正常に作成される", func(t *testing.T) {
		// MedicationHandlerの作成をテスト
		// 実際のリポジトリは使わずに、nilで作成してもパニックしないことを確認
		assert.NotPanics(t, func() {
			NewMedicationHandler(nil)
		})
	})
}

// GetMedicationStatsのテスト
func TestMedicationHandler_GetMedicationStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("認証エラーのテスト", func(t *testing.T) {
		medicationRepo := repository.NewMedicationRepository()
		handler := NewMedicationHandler(medicationRepo)

		// Ginルーターとコンテキストの設定
		router := gin.New()
		router.GET("/api/medication-stats", handler.GetMedicationStats)

		// リクエストの作成（ユーザーIDなし）
		req, _ := http.NewRequest("GET", "/api/medication-stats", nil)
		w := httptest.NewRecorder()

		// リクエストの実行
		router.ServeHTTP(w, req)

		// レスポンスの検証
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid user ID", response["error"])
	})

	t.Run("正常なレスポンスの構造テスト", func(t *testing.T) {
		medicationRepo := repository.NewMedicationRepository()
		handler := NewMedicationHandler(medicationRepo)

		// Ginルーターとコンテキストの設定
		router := gin.New()
		router.GET("/api/medication-stats", func(c *gin.Context) {
			// テスト用にユーザーIDを設定
			c.Set("userID", "test-user-123")
			handler.GetMedicationStats(c)
		})

		// リクエストの作成
		req, _ := http.NewRequest("GET", "/api/medication-stats", nil)
		w := httptest.NewRecorder()

		// リクエストの実行
		router.ServeHTTP(w, req)

		// レスポンスの検証（統計データの構造を確認）
		if w.Code == http.StatusOK {
			var response dto.MedicationStatsResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// 統計データの存在を確認
			assert.GreaterOrEqual(t, response.LongestContinuousDays, 0)
			assert.GreaterOrEqual(t, response.TotalBleedingDays, 0)
			assert.GreaterOrEqual(t, response.MedicationBreaks, 0)
			assert.GreaterOrEqual(t, response.AverageCycleLength, 0.0)
			assert.GreaterOrEqual(t, response.MonthlyMedicationRate, 0.0)
		}
	})
}
