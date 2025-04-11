package service

import (
	"errors"
	"log"
	"okusuri-backend/repository"
)

// カスタムエラー
var (
	ErrEmailAlreadyExists  = errors.New("このメールアドレスは既に使用されています")
	ErrInvalidCredentials  = errors.New("メールアドレスまたはパスワードが無効です")
	ErrUserNotFound        = errors.New("ユーザーが見つかりません")
	ErrTokenGenerationFail = errors.New("トークン生成に失敗しました")
	ErrInvalidToken        = errors.New("無効なトークンです")
)

// AuthService は認証関連の機能を提供する
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService は新しいAuthServiceのインスタンスを作成する
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// RegisterUser は新規ユーザーを登録する
func (s *AuthService) RegisterUser(email, password string) {
	log.Println("RegisterUser called")
}
