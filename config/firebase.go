package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App

// InitFirebase はFirebaseを初期化する関数
func InitFirebase() error {
	ctx := context.Background()

	// 環境変数から秘密鍵のパスを取得するか、デフォルトパスを使用
	credentialPath := os.Getenv("FIREBASE_CREDENTIAL_PATH")
	if credentialPath == "" {
		// 環境変数が設定されていない場合、デフォルトパスを使用
		credentialPath = filepath.Join("config", "firebase-service-account.json")
	}

	// サービスアカウントのJSONファイルパスを指定
	opt := option.WithCredentialsFile(credentialPath)

	var err error
	firebaseApp, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("firebase初期化エラー: %v", err)
	}

	return nil
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
