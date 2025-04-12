package service

import (
	"fmt"
	"okusuri-backend/dto"
	"okusuri-backend/helper"
	"okusuri-backend/model"
	"okusuri-backend/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) RegisterUser(req dto.SignupRequest) (*model.User, error) {
	// IDトークンの検証
	switch req.Provider {
	case "google":
		tokenInfo := helper.VerifyGoogleIDToken(req.IDToken)
		if tokenInfo == nil {
			return nil, fmt.Errorf("無効なIDトークン")
		}
		if tokenInfo.Email != req.Email {
			return nil, fmt.Errorf("IDトークンのメールアドレスとリクエストのメールアドレスが一致しません")
		}
	}

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
