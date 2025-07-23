package handler

import (
	"net/http"
	"okusuri-backend/internal/dto"
	"okusuri-backend/internal/model"
	"okusuri-backend/internal/repository"
	"okusuri-backend/internal/service"
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
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	medicationLog := model.MedicationLog{
		UserID:      userID,
		HasBleeding: req.HasBleeding,
	}

	// 日付が指定されている場合は、その日付を使用
	if req.Date != nil {
		medicationLog.CreatedAt = *req.Date
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

// UpdateLog は指定されたIDの服薬ログを更新するハンドラー
func (h *MedicationHandler) UpdateLog(c *gin.Context) {
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

	// リクエストボディを構造体にバインド
	var req dto.MedicationLogRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.medicationRepo.UpdateLog(userID, uint(logID), req.HasBleeding)
	if err != nil {
		if err.Error() == "log not found or user not authorized" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update medication log"})
		return

	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "medication log updated successfully",
	})
}

// GetMedicationStatus は現在の服薬ステータスを取得するハンドラー
func (h *MedicationHandler) GetMedicationStatus(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// サービスから服薬ステータスを取得
	medicationService := service.NewMedicationService(h.medicationRepo)
	status, err := medicationService.GetMedicationStatus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get medication status"})
		return
	}

	c.JSON(http.StatusOK, status)
}


