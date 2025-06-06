package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"okusuri-backend/internal/dto"
	"okusuri-backend/internal/model"
	"okusuri-backend/internal/repository"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// テスト用のgin.Contextを作成する関数
func createTestContext(method, path string, body interface{}, userID string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// ユーザーIDをコンテキストに設定
	if userID != "" {
		c.Set("userID", userID)
	}

	return c, w
}

// TestMedicationHandler_RegisterLog は服薬ログ登録のテスト（構造体の基本テスト）
func TestMedicationHandler_RegisterLog(t *testing.T) {
	t.Run("リクエストボディの構造テスト", func(t *testing.T) {
		// リクエストボディの基本構造をテスト
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		req := dto.MedicationLogRequest{
			HasBleeding: true,
			Date:        &testDate,
		}

		assert.True(t, req.HasBleeding)
		assert.Equal(t, testDate, *req.Date)
	})

	t.Run("ユーザーIDなしのエラーハンドリング", func(t *testing.T) {
		// ハンドラーの作成
		handler := NewMedicationHandler(repository.NewMedicationRepository())

		// リクエストの作成（ユーザーIDなし）
		requestBody := dto.MedicationLogRequest{HasBleeding: true}
		c, w := createTestContext("POST", "/api/medication-log", requestBody, "")

		// ハンドラーの実行
		handler.RegisterLog(c)

		// レスポンスの検証
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid user ID", response["error"])
	})
}

// TestMedicationHandler_GetLogs は服薬ログ取得のテスト（構造体の基本テスト）
func TestMedicationHandler_GetLogs(t *testing.T) {
	t.Run("ユーザーIDなしのエラーハンドリング", func(t *testing.T) {
		// ハンドラーの作成
		handler := NewMedicationHandler(repository.NewMedicationRepository())

		// リクエストの作成（ユーザーIDなし）
		c, w := createTestContext("GET", "/api/medication-log", nil, "")

		// ハンドラーの実行
		handler.GetLogs(c)

		// レスポンスの検証
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid user ID", response["error"])
	})

	t.Run("MedicationLogの構造テスト", func(t *testing.T) {
		// MedicationLogの基本構造をテスト
		logs := []model.MedicationLog{
			{
				ID:          1,
				UserID:      "test-user-1",
				HasBleeding: true,
				CreatedAt:   time.Now(),
			},
			{
				ID:          2,
				UserID:      "test-user-1",
				HasBleeding: false,
				CreatedAt:   time.Now().Add(-24 * time.Hour),
			},
		}

		assert.Len(t, logs, 2)
		assert.Equal(t, "test-user-1", logs[0].UserID)
		assert.Equal(t, "test-user-1", logs[1].UserID)
		assert.True(t, logs[0].HasBleeding)
		assert.False(t, logs[1].HasBleeding)
	})
}

// TestMedicationHandler_GetLogByID は特定ログ取得のテスト（構造体の基本テスト）
func TestMedicationHandler_GetLogByID(t *testing.T) {
	t.Run("ユーザーIDなしのエラーハンドリング", func(t *testing.T) {
		// ハンドラーの作成
		handler := NewMedicationHandler(repository.NewMedicationRepository())

		// リクエストの作成（ユーザーIDなし）
		c, w := createTestContext("GET", "/api/medication-log/1", nil, "")
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		// ハンドラーの実行
		handler.GetLogByID(c)

		// レスポンスの検証
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid user ID", response["error"])
	})

	t.Run("パラメータの構造テスト", func(t *testing.T) {
		// パラメータの基本構造をテスト
		params := []gin.Param{{Key: "id", Value: "1"}}

		assert.Len(t, params, 1)
		assert.Equal(t, "id", params[0].Key)
		assert.Equal(t, "1", params[0].Value)
	})
}

// TestMedicationHandler_UpdateLog は服薬ログ更新のテスト（構造体の基本テスト）
func TestMedicationHandler_UpdateLog(t *testing.T) {
	t.Run("ユーザーIDなしのエラーハンドリング", func(t *testing.T) {
		// ハンドラーの作成
		handler := NewMedicationHandler(repository.NewMedicationRepository())

		// リクエストの作成（ユーザーIDなし）
		requestBody := dto.MedicationLogRequest{HasBleeding: true}
		c, w := createTestContext("PUT", "/api/medication-log/1", requestBody, "")
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		// ハンドラーの実行
		handler.UpdateLog(c)

		// レスポンスの検証
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid user ID", response["error"])
	})

	t.Run("更新リクエストの構造テスト", func(t *testing.T) {
		// 更新リクエストの基本構造をテスト
		requestBody := dto.MedicationLogRequest{HasBleeding: true}

		assert.True(t, requestBody.HasBleeding)
		assert.Nil(t, requestBody.Date)
	})
}

// TestMedicationHandler_New は MedicationHandler 作成のテスト
func TestMedicationHandler_New(t *testing.T) {
	t.Run("MedicationHandlerが正常に作成される", func(t *testing.T) {
		repo := repository.NewMedicationRepository()
		handler := NewMedicationHandler(repo)

		assert.NotNil(t, handler)
		assert.Equal(t, repo, handler.medicationRepo)
	})
}
