package user

import "time"

// User モデル
type User struct {
	ID            string    `json:"id" gorm:"primary_key"`
	Name          string    `json:"name"`
	Email         string    `json:"email" gorm:"unique"`
	EmailVerified bool      `json:"emailVerified"`
	Image         *string   `json:"image"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (User) TableName() string {
	return "user"
}

// Session モデル
type Session struct {
	ID        string    `json:"id" gorm:"primary_key;column:id"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"column:expiresAt"`
	Token     string    `json:"token" gorm:"unique;column:token"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt"`
	IPAddress *string   `json:"ipAddress" gorm:"column:ipAddress"`
	UserAgent *string   `json:"userAgent" gorm:"column:userAgent"`
	UserID    string    `json:"userId" gorm:"column:userId"`
}

func (Session) TableName() string {
	return "session"
}

// Account モデル
type Account struct {
	ID                    string     `json:"id" gorm:"primary_key"`
	AccountID             string     `json:"accountId"`
	ProviderID            string     `json:"providerId"`
	UserID                string     `json:"userId"`
	AccessToken           *string    `json:"accessToken"`
	RefreshToken          *string    `json:"refreshToken"`
	IDToken               *string    `json:"idToken"`
	AccessTokenExpiresAt  *time.Time `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt *time.Time `json:"refreshTokenExpiresAt"`
	Scope                 *string    `json:"scope"`
	Password              *string    `json:"password"`
	CreatedAt             time.Time  `json:"createdAt"`
	UpdatedAt             time.Time  `json:"updatedAt"`
}

func (Account) TableName() string {
	return "account"
}

// Verification モデル
type Verification struct {
	ID         string     `json:"id" gorm:"primary_key"`
	Identifier string     `json:"identifier"`
	Value      string     `json:"value"`
	ExpiresAt  time.Time  `json:"expiresAt"`
	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}

func (Verification) TableName() string {
	return "verification"
}
