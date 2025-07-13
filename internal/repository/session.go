package repository

import (
	"okusuri-backend/internal/model"

	"gorm.io/gorm"
)

// SessionRepository はセッション関連のリポジトリ
type SessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository は新しいSessionRepositoryを作成
func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create はセッションを作成
func (r *SessionRepository) Create(session *model.Session) error {
	return r.db.Create(session).Error
}

// FindByToken はトークンでセッションを検索
func (r *SessionRepository) FindByToken(token string) (*model.Session, error) {
	var session model.Session
	err := r.db.Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// FindByUserID はユーザーIDでセッションを検索
func (r *SessionRepository) FindByUserID(userID string) ([]*model.Session, error) {
	var sessions []*model.Session
	err := r.db.Where("user_id = ?", userID).Find(&sessions).Error
	return sessions, err
}

// Update はセッションを更新
func (r *SessionRepository) Update(session *model.Session) error {
	return r.db.Save(session).Error
}

// Delete はセッションを削除
func (r *SessionRepository) Delete(id string) error {
	return r.db.Delete(&model.Session{}, "id = ?", id).Error
}

// DeleteByToken はトークンでセッションを削除
func (r *SessionRepository) DeleteByToken(token string) error {
	return r.db.Delete(&model.Session{}, "token = ?", token).Error
}

// DeleteByUserID はユーザーIDでセッションを削除
func (r *SessionRepository) DeleteByUserID(userID string) error {
	return r.db.Delete(&model.Session{}, "user_id = ?", userID).Error
}

// DeleteExpired は期限切れセッションを削除
func (r *SessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < NOW()").Delete(&model.Session{}).Error
}