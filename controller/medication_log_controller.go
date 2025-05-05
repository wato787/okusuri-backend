package controller

import (
	"net/http"
	"okusuri-backend/dto"
	"okusuri-backend/helper"
	"okusuri-backend/model"
	"okusuri-backend/repository"
	"okusuri-backend/service"

	"github.com/gin-gonic/gin"
)

type MedicationLogController struct {
	MedicationLogRepository *repository.MedicationLogRepository
	MedicationLogService    *service.MedicationLogService
}

// NewMedicationLogController は新しいMedicationLogControllerのインスタンスを作成する
func NewMedicationLogController() *MedicationLogController {
	medicationLogRepository := repository.NewMedicationLogRepository()
	medicationLogService := service.NewMedicationLogService(medicationLogRepository)
	return &MedicationLogController{
		MedicationLogRepository: medicationLogRepository,
		MedicationLogService:    medicationLogService,
	}
}

func (mc *MedicationLogController) RegisterMedicationLog(c *gin.Context) {
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

	// MedicationLogを作成
	err = mc.MedicationLogService.RegisterMedicationLog(userID, medicationLog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register medication log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "medication log registered successfully"})
}

func (mc *MedicationLogController) GetMedicationLogs(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// MedicationLogを取得
	medicationLogs, err := mc.MedicationLogService.GetMedicationLogsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get medication logs"})
		return
	}

	c.JSON(http.StatusOK, medicationLogs)
}
