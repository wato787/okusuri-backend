package handler

import (
	"net/http"
	"okusuri-backend/internal/dto"
	"okusuri-backend/internal/model"
	"okusuri-backend/internal/repository"
	"okusuri-backend/pkg/helper"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MedicationHandler struct {
	medicationRepo *repository.MedicationRepository
}

func NewMedicationHandler(medicationRepo *repository.MedicationRepository) *MedicationHandler {
	return &MedicationHandler{
		medicationRepo: medicationRepo,
	}
}

func (h *MedicationHandler) RegisterLog(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// リクエストボディを構造体にバインド
	var req dto.MedicationLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	medicationLog := model.MedicationLog{
		UserID:      userID,
		HasBleeding: req.HasBleeding,
	}

	// リポジトリを直接呼び出す
	err = h.medicationRepo.RegisterLog(userID, medicationLog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register medication log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "medication log registered successfully"})
}

// GetLogs はユーザーの服用記録を取得するハンドラー
func (h *MedicationHandler) GetLogs(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// 服用記録を取得
	logs, err := h.medicationRepo.GetLogsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get medication logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetLogByID は特定のIDの服薬ログを取得するハンドラー
func (h *MedicationHandler) GetLogByID(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// URLからIDパラメータを取得
	logIDStr := c.Param("id")
	logID, err := strconv.ParseUint(logIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log ID"})
		return
	}

	// 服薬ログを取得
	log, err := h.medicationRepo.GetLogByID(userID, uint(logID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "medication log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}
