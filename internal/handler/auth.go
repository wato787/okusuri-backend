package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"okusuri-backend/internal/model"
	"okusuri-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler は認証関連のハンドラー
type AuthHandler struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	accountRepo *repository.AccountRepository
}

// NewAuthHandler は新しいAuthHandlerを作成
func NewAuthHandler(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository, accountRepo *repository.AccountRepository) *AuthHandler {
	return &AuthHandler{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		accountRepo: accountRepo,
	}
}

// GoogleOAuthInfo はGoogle OAuthユーザー情報
type GoogleOAuthInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// SignInWithGoogle はGoogle OAuth認証を開始
func (h *AuthHandler) SignInWithGoogle(c *gin.Context) {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Google OAuth設定が不足しています"})
		return
	}

	redirectURL := fmt.Sprintf("%s/api/auth/callback/google", os.Getenv("APP_URL"))
	oauthURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&scope=%s&response_type=code&state=%s",
		googleClientID,
		url.QueryEscape(redirectURL),
		url.QueryEscape("openid profile email"),
		"random_state", // 本来はCSRF対策のためランダム文字列を生成
	)

	c.JSON(http.StatusOK, gin.H{
		"url": oauthURL,
	})
}

// GoogleCallback はGoogle OAuthコールバック処理
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "認証コードがありません"})
		return
	}

	// アクセストークンを取得
	token, err := h.exchangeCodeForToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークン取得に失敗しました"})
		return
	}

	// ユーザー情報を取得
	userInfo, err := h.getGoogleUserInfo(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報取得に失敗しました"})
		return
	}

	// ユーザーを作成または取得
	user, err := h.findOrCreateUser(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー作成に失敗しました"})
		return
	}

	// セッションを作成
	session, err := h.createSession(user.ID, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "セッション作成に失敗しました"})
		return
	}

	// Accountレコードを作成/更新
	err = h.createOrUpdateAccount(user.ID, "google", userInfo.ID, token)
	if err != nil {
		fmt.Printf("アカウント情報更新エラー: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "アカウント情報更新に失敗しました"})
		return
	}

	// フロントエンドにリダイレクト（セッショントークンをクエリパラメータで渡す）
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5174"
	}

	redirectURL := fmt.Sprintf("%s?token=%s", frontendURL, session.Token)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// GetSession は現在のセッション情報を取得
func (h *AuthHandler) GetSession(c *gin.Context) {
	token := h.extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	session, err := h.sessionRepo.FindByToken(token)
	if err != nil || session == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なセッションです"})
		return
	}

	// セッションの有効期限確認
	if session.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "セッションが期限切れです"})
		return
	}

	user, err := h.userRepo.FindByID(session.UserID)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"session": gin.H{
			"id":        session.ID,
			"token":     session.Token,
			"expiresAt": session.ExpiresAt,
		},
	})
}

// SignOut はサインアウト処理
func (h *AuthHandler) SignOut(c *gin.Context) {
	token := h.extractToken(c)
	if token == "" {
		c.JSON(http.StatusOK, gin.H{"message": "サインアウトしました"})
		return
	}

	// セッションを削除
	err := h.sessionRepo.DeleteByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サインアウトに失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "サインアウトしました"})
}

// exchangeCodeForToken は認証コードをアクセストークンに交換
func (h *AuthHandler) exchangeCodeForToken(code string) (string, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := fmt.Sprintf("%s/api/auth/callback/google", os.Getenv("APP_URL"))

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURL)

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}

	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

// getGoogleUserInfo はGoogleからユーザー情報を取得
func (h *AuthHandler) getGoogleUserInfo(token string) (*GoogleOAuthInfo, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleOAuthInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// findOrCreateUser はユーザーを検索または作成
func (h *AuthHandler) findOrCreateUser(userInfo *GoogleOAuthInfo) (*model.User, error) {
	// 既存ユーザーを検索
	user, err := h.userRepo.FindByEmail(userInfo.Email)
	if err == nil && user != nil {
		// ユーザー情報を更新
		user.Name = userInfo.Name
		user.EmailVerified = userInfo.VerifiedEmail
		if userInfo.Picture != "" {
			user.Image = &userInfo.Picture
		}
		err = h.userRepo.Update(user)
		return user, err
	}

	// 新規ユーザーを作成
	newUser := &model.User{
		ID:            uuid.New().String(),
		Name:          userInfo.Name,
		Email:         userInfo.Email,
		EmailVerified: userInfo.VerifiedEmail,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if userInfo.Picture != "" {
		newUser.Image = &userInfo.Picture
	}

	err = h.userRepo.Create(newUser)
	return newUser, err
}

// createSession はセッションを作成
func (h *AuthHandler) createSession(userID, ipAddress, userAgent string) (*model.Session, error) {
	session := &model.Session{
		ID:        uuid.New().String(),
		Token:     uuid.New().String(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7日間有効
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if ipAddress != "" {
		session.IPAddress = &ipAddress
	}
	if userAgent != "" {
		session.UserAgent = &userAgent
	}

	err := h.sessionRepo.Create(session)
	return session, err
}

// createOrUpdateAccount はアカウント情報を作成または更新
func (h *AuthHandler) createOrUpdateAccount(userID, providerID, accountID, accessToken string) error {
	fmt.Printf("createOrUpdateAccount開始: userID=%s, providerID=%s, accountID=%s\n", userID, providerID, accountID)

	account, err := h.accountRepo.FindByProviderAndAccountID(providerID, accountID)
	if err != nil {
		fmt.Printf("既存アカウント検索エラー: %v\n", err)
	}

	if err == nil && account != nil {
		// 既存アカウントを更新
		fmt.Printf("既存アカウントを更新: %s\n", account.ID)
		account.AccessToken = &accessToken
		account.UpdatedAt = time.Now()
		err = h.accountRepo.Update(account)
		if err != nil {
			fmt.Printf("アカウント更新エラー: %v\n", err)
		}
		return err
	}

	// 新規アカウントを作成
	fmt.Printf("新規アカウントを作成\n")
	newAccount := &model.Account{
		ID:          uuid.New().String(),
		AccountID:   accountID,
		ProviderID:  providerID,
		UserID:      userID,
		AccessToken: &accessToken,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = h.accountRepo.Create(newAccount)
	if err != nil {
		fmt.Printf("アカウント作成エラー: %v\n", err)
	}
	return err
}

// extractToken はリクエストからトークンを抽出
func (h *AuthHandler) extractToken(c *gin.Context) string {
	// Bearerトークンから抽出
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	// クエリパラメータから抽出
	return c.Query("token")
}
