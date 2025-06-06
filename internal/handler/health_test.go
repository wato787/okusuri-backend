package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Ginのテストモードに設定
	gin.SetMode(gin.TestMode)

	// テスト用のルーターを作成
	router := gin.New()

	// ヘルスチェックエンドポイントを追加
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// テストケース
	t.Run("正常なヘルスチェック", func(t *testing.T) {
		// リクエストを作成
		req, err := http.NewRequest("GET", "/api/health", nil)
		assert.NoError(t, err)

		// レスポンスレコーダーを作成
		w := httptest.NewRecorder()

		// リクエストを実行
		router.ServeHTTP(w, req)

		// ステータスコードを確認
		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスボディを確認
		expected := `{"status":"ok"}`
		assert.JSONEq(t, expected, w.Body.String())
	})
}
