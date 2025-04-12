package service

import (
	"fmt"
	"okusuri-backend/config"
	"okusuri-backend/dto"
	"okusuri-backend/model"
	"okusuri-backend/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) RegisterUser(req dto.SignupRequest) (*model.User, error) {
	// IDトークンの検証

	// ユーザーの存在確認
	existingUser, err := s.userRepo.FindByProviderId(req.ProviderID)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, nil // ユーザーは既に存在する
	}

	// 新しいユーザーの作成
	newUser := &model.User{
		Email:      req.Email,
		Name:       req.Name,
		ImageUrl:   req.ImageURL,
		Provider:   model.ProviderGoogle,
		ProviderId: req.ProviderID,
	}

	// ユーザーの保存
	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil

}

// ユーザーIDを受け取り、JWTトークンを生成するメソッド
func (s *AuthService) GenerateToken(userID uint) (string, int64, error) {
	// トークンの有効期限を設定
	expirationTime := time.Now().Add(2400 * time.Hour)
	expiresAt := expirationTime.Unix()

	// トークンの内容（クレーム）を作成
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt,
		"iat":     time.Now().Unix(), // 発行時刻（Issued At）
	}

	// トークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// シークレットキーでトークンに署名
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}
