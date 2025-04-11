package routes

import (
	"okusuri-backend/controller"
	"okusuri-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	// コントローラを初期化
	authController := controller.NewAuthController()

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

		// 認証ルートの登録
		RegisterAuthRoutes(api, authController)

		// 認証が必要なルートグループ
		secured := api.Group("/")
		secured.Use(middleware.JWTAuth())
		{
			// ここに認証が必要なルートを追加
			// 例: RegisterUserRoutes(secured, userController)
		}
	}

	return router
}
