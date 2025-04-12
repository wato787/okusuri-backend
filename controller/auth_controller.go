package controller

import (
	"net/http"
	"okusuri-backend/config"
	"okusuri-backend/dto"
	"okusuri-backend/repository"
	"okusuri-backend/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController() *AuthController {
	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo, &config.Config{})

	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	// ユーザー登録処理
	user, err := ac.authService.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー登録に失敗しました"})
		return
	} else if user == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "ユーザーは既に存在します"})
		return
	}

	// トークン生成処理
	token, expiresAt, err := ac.authService.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, dto.AuthResponse{
		User: dto.UserResponse{
			UserID:   user.ID,
			Email:    user.Email,
			Name:     user.Name,
			ImageURL: user.ImageUrl,
		},
		Token:     token,
		ExpiresAt: expiresAt,
		Message:   "ユーザー登録が成功しました",
	})

}
