package controller

import (
	"log"
	"okusuri-backend/repository"
	"okusuri-backend/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController() *AuthController {
	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)

	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Signup(c *gin.Context) {
	log.Println("Signup called")
}
