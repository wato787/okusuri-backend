package controller

import (
	"net/http"
	"okusuri-backend/dto"
	"okusuri-backend/helper"
	"okusuri-backend/repository"
	"okusuri-backend/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService *service.UserService
}

func NewAuthController() *AuthController {
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)

	return &AuthController{
		userService: userService,
	}
}

func (ac *AuthController) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	// IDトークンの検証
	switch req.Provider {
	case "google":
		tokenInfo := helper.VerifyGoogleIDToken(req.IDToken)
		if tokenInfo == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDトークン"})
			return
		}
	}

	// ユーザー登録処理
	user, err := ac.userService.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー登録に失敗しました"})
		return
	} else if user == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "ユーザーは既に存在します"})
		return
	}

	// トークン生成処理
	token, expiresAt, err := helper.GenerateToken(user.ID)
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

func (ac *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	// IDトークンの検証
	switch req.Provider {
	case "google":
		tokenInfo := helper.VerifyGoogleIDToken(req.IDToken)
		if tokenInfo == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なIDトークン"})
			return
		}
	}

	// ユーザー取得
	user, err := ac.userService.GetUserByProviderId(req.ProviderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー取得に失敗しました"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}
	// トークン生成処理
	token, expiresAt, err := helper.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, dto.AuthResponse{
		User: dto.UserResponse{
			UserID:   user.ID,
			Email:    user.Email,
			Name:     user.Name,
			ImageURL: user.ImageUrl,
		},
		Token:     token,
		ExpiresAt: expiresAt,
		Message:   "ログインが成功しました",
	})
}
