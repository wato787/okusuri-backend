package service

import (
	"log"
	"okusuri-backend/repository"
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
