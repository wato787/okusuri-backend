package helper

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"google.golang.org/api/oauth2/v1"
	"google.golang.org/api/option"
)

// JWTトークンを検証してユーザーIDを取得するメソッド
func ValidateToken(tokenString string) (uint, error) {
	// 環境変数を読み込む
	err := godotenv.Load()
	if err != nil {
		return 0, fmt.Errorf("環境変数の読み込みに失敗しました: %w", err)
	}

	JWTSecret := os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		return 0, fmt.Errorf("JWT_SECRETが設定されていません")
	}

	// トークンを解析
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("トークンの検証に失敗しました: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, fmt.Errorf("無効なトークン")
}

// ユーザーIDを受け取り、JWTトークンを生成するメソッド
func GenerateToken(userID uint) (string, int64, error) {
	err := godotenv.Load()
	if err != nil {
		return "", 0, fmt.Errorf("環境変数の読み込みに失敗しました: %w", err)
	}
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

	JWTSecret := os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		return "", 0, fmt.Errorf("JWT_SECRETが設定されていません")
	}

	// シークレットキーでトークンに署名
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// GoogleのIDトークンを検証するメソッド
func VerifyGoogleIDToken(idToken string) *oauth2.Tokeninfo {
	httpClient := &http.Client{}
	oauth2Service, err := oauth2.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return nil
	}
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil
	}
	return tokenInfo
}
