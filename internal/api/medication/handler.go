package medication

import (
	"net/http"

	"okusuri-backend/pkg/helper"

	"github.com/gin-gonic/gin"
)

// Handler は薬の服用記録に関するHTTPハンドラー
type Handler struct {
	service *Service
}

// NewHandler は新しいHandler インスタンスを作成する
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterLog は服用記録を登録するハンドラー
func (h *Handler) RegisterLog(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// リクエストボディを構造体にバインド
	var req MedicationLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	medicationLog := MedicationLog{
		UserID:      userID,
		HasBleeding: req.HasBleeding,
	}

	// 服用記録を作成
	err = h.service.RegisterLog(userID, medicationLog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register medication log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "medication log registered successfully"})
}

// GetLogs はユーザーの服用記録を取得するハンドラー
func (h *Handler) GetLogs(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// 服用記録を取得
	logs, err := h.service.GetLogsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get medication logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}
