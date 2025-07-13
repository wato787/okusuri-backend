package repository

import (
	"okusuri-backend/internal/model"

	"gorm.io/gorm"
)

// AccountRepository はアカウント関連のリポジトリ
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository は新しいAccountRepositoryを作成
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create はアカウントを作成
func (r *AccountRepository) Create(account *model.Account) error {
	return r.db.Create(account).Error
}

// FindByID はIDでアカウントを検索
func (r *AccountRepository) FindByID(id string) (*model.Account, error) {
	var account model.Account
	err := r.db.Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// FindByUserID はユーザーIDでアカウントを検索
func (r *AccountRepository) FindByUserID(userID string) ([]*model.Account, error) {
	var accounts []*model.Account
	err := r.db.Where("user_id = ?", userID).Find(&accounts).Error
	return accounts, err
}

// FindByProviderAndAccountID はプロバイダIDとアカウントIDでアカウントを検索
func (r *AccountRepository) FindByProviderAndAccountID(providerID, accountID string) (*model.Account, error) {
	var account model.Account
	err := r.db.Where("provider_id = ? AND account_id = ?", providerID, accountID).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// Update はアカウントを更新
func (r *AccountRepository) Update(account *model.Account) error {
	return r.db.Save(account).Error
}

// Delete はアカウントを削除
func (r *AccountRepository) Delete(id string) error {
	return r.db.Delete(&model.Account{}, "id = ?", id).Error
}

// DeleteByUserID はユーザーIDでアカウントを削除
func (r *AccountRepository) DeleteByUserID(userID string) error {
	return r.db.Delete(&model.Account{}, "user_id = ?", userID).Error
}