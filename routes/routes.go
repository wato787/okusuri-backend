package routes

import (
	"okusuri-backend/middleware"
	"okusuri-backend/repository"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	// userRepositoryを初期化
	userRepository := repository.NewUserRepository()

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

		// 認証が必要なルートグループ
		secured := api.Group("/")
		secured.Use(middleware.Auth(userRepository))
		{
			// ここに認証が必要なルートを追加
			// 例: RegisterUserRoutes(secured, userController)
		}
	}

	return router
}
