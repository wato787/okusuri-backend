package dto

// 新規登録のリクエストDTO
type SignupRequest struct {
	IDToken    string `json:"id_token" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Name       string `json:"name" binding:"required"`
	ImageURL   string `json:"image_url"` // 任意フィールド
	ProviderID string `json:"provider_id" binding:"required"`
}

type SignupResponse struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	ImageURL  string `json:"image_url"`
	Token     string `json:"token"`      // JWTアクセストークン
	ExpiresAt int64  `json:"expires_at"` // トークンの有効期限（UNIXタイムスタンプ）
}
