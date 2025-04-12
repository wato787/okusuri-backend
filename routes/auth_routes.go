package routes

import (
	"okusuri-backend/controller"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, authController *controller.AuthController) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", authController.Signup)
		auth.POST("/login", authController.Login)
		// 認証が必要なルート
		// authorized := auth.Group("/")
		// authorized.Use(middleware.JWTAuth())
		{
			// authorized.POST("/logout", controllers.Logout)
		}
	}
}
