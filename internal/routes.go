package internal

import (
	"okusuri-backend/internal/handler"
	"okusuri-backend/internal/middleware"
	"okusuri-backend/internal/repository"
	"okusuri-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	// リポジトリの初期化
	userRepo := repository.NewUserRepository()
	sessionRepo := repository.NewSessionRepository(userRepo.GetDB())
	accountRepo := repository.NewAccountRepository(userRepo.GetDB())
	medicationRepo := repository.NewMedicationRepository()
	notificationRepo := repository.NewNotificationRepository()

	// サービスの初期化
	notificationService := service.NewNotificationService()
	medicationService := service.NewMedicationService(medicationRepo)

	// ハンドラーの初期化
	authHandler := handler.NewAuthHandler(userRepo, sessionRepo, accountRepo)
	medicationHandler := handler.NewMedicationHandler(medicationRepo)
	notificationHandler := handler.NewNotificationHandler(
		notificationRepo,
		userRepo,
		notificationService,
		medicationRepo,
		medicationService,
	)

	// Ginのルーターを作成
	router := gin.Default()

	// グローバルミドルウェアの設定
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})

		// 認証関連のエンドポイント
		auth := api.Group("/auth")
		{
			auth.GET("/google", authHandler.SignInWithGoogle)
			auth.GET("/callback/google", authHandler.GoogleCallback)
			auth.GET("/session", authHandler.GetSession)
			auth.POST("/signout", authHandler.SignOut)
		}

		api.POST(("/notification"), notificationHandler.SendNotification)

		// 新しいエンドポイントを追加
		api.GET("/medication-status", middleware.Auth(userRepo), medicationHandler.GetMedicationStatus)

		notificationSetting := api.Group("/notification/setting")
		notificationSetting.Use(middleware.Auth(userRepo))
		{
			notificationSetting.GET("", notificationHandler.GetSetting)
			notificationSetting.POST("", notificationHandler.RegisterSetting)
		}

		medicationLog := api.Group("/medication-log")
		medicationLog.Use(middleware.Auth(userRepo))
		{
			medicationLog.POST("", medicationHandler.RegisterLog)
			medicationLog.GET("", medicationHandler.GetLogs)
			medicationLog.GET("/:id", medicationHandler.GetLogByID)
			medicationLog.PATCH("/:id", medicationHandler.UpdateLog)
		}
	}

	return router
}
