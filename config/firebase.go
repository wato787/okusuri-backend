package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App

// InitFirebase はFirebaseを初期化する関数
func InitFirebase() error {
	ctx := context.Background()

	// 個別の環境変数から認証情報JSONを構築
	credentialJSON, err := buildCredentialJSON()
	if err != nil {
		return fmt.Errorf("認証情報の構築エラー: %v", err)
	}

	// JSONデータから認証情報を設定
	opt := option.WithCredentialsJSON([]byte(credentialJSON))

	firebaseApp, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("firebase初期化エラー: %v", err)
	}

	return nil
}

// buildCredentialJSON は環境変数から認証情報JSONを構築する関数
func buildCredentialJSON() (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not loaded")
	}
	// サービスアカウントJSONの構造に対応した構造体
	type ServiceAccountCredential struct {
		Type                    string `json:"type"`
		ProjectID               string `json:"project_id"`
		PrivateKeyID            string `json:"private_key_id"`
		PrivateKey              string `json:"private_key"`
		ClientEmail             string `json:"client_email"`
		ClientID                string `json:"client_id"`
		AuthURI                 string `json:"auth_uri"`
		TokenURI                string `json:"token_uri"`
		AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
		ClientX509CertURL       string `json:"client_x509_cert_url"`
	}

	// 環境変数から値を取得
	credential := ServiceAccountCredential{
		Type:                    os.Getenv("FIREBASE_TYPE"),
		ProjectID:               os.Getenv("FIREBASE_PROJECT_ID"),
		PrivateKeyID:            os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
		PrivateKey:              os.Getenv("FIREBASE_PRIVATE_KEY"),
		ClientEmail:             os.Getenv("FIREBASE_CLIENT_EMAIL"),
		ClientID:                os.Getenv("FIREBASE_CLIENT_ID"),
		AuthURI:                 os.Getenv("FIREBASE_AUTH_URI"),
		TokenURI:                os.Getenv("FIREBASE_TOKEN_URI"),
		AuthProviderX509CertURL: os.Getenv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL"),
		ClientX509CertURL:       os.Getenv("FIREBASE_CLIENT_X509_CERT_URL"),
	}

	// 必須フィールドの検証
	if credential.Type == "" || credential.ProjectID == "" || credential.PrivateKey == "" || credential.ClientEmail == "" {
		return "", fmt.Errorf("必須の環境変数が設定されていません")
	}

	// 構造体をJSON文字列に変換
	jsonData, err := json.Marshal(credential)
	if err != nil {
		return "", fmt.Errorf("JSON変換エラー: %v", err)
	}

	return string(jsonData), nil
}

// GetMessagingClient はFCMクライアントを取得する関数
func GetMessagingClient(ctx context.Context) (*messaging.Client, error) {
	if firebaseApp == nil {
		return nil, fmt.Errorf("firebaseが初期化されていません")
	}

	client, err := firebaseApp.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("FCMクライアント作成エラー: %v", err)
	}

	return client, nil
}
