# okusuri-backend

## 概要

`okusuri-backend`は、服薬管理アプリケーションのバックエンドAPIサーバーです。Go言語で開発され、服薬記録の管理、通知機能、ユーザー認証などの機能を提供します。

## プロジェクトの特徴

### 🏗️ アーキテクチャ
- **レイヤードアーキテクチャ**を採用し、責務を明確に分離
- **依存性注入パターン**を使用して、コンポーネント間の結合度を低減
- **RESTful API**設計に準拠

### 🔧 技術スタック
- **言語**: Go 1.24
- **Webフレームワーク**: Gin
- **データベース**: PostgreSQL + GORM ORM
- **認証**: Google OAuth 2.0 + JWT
- **通知**: Web Push API (webpush-go)
- **環境設定**: godotenv
- **開発ツール**: Air (ホットリロード)

### 📁 ディレクトリ構造

```
okusuri-backend/
├── cmd/server/          # エントリーポイント
├── internal/            # 内部パッケージ
│   ├── dto/            # データ転送オブジェクト
│   ├── handler/        # HTTPハンドラー
│   ├── middleware/     # ミドルウェア
│   ├── model/          # データベースモデル
│   ├── repository/     # データアクセス層
│   ├── service/        # ビジネスロジック
│   └── routes.go       # ルーティング設定
├── pkg/                # 外部公開パッケージ
│   ├── config/         # 設定管理
│   └── helper/         # ユーティリティ関数
├── migrations/         # データベースマイグレーション
└── scripts/            # 開発・運用スクリプト
```

## 主要機能

### 1. ユーザー認証・管理
- **Google OAuth 2.0**による認証
- **セッション管理**（トークンベース）
- **ユーザー情報管理**（名前、メール、画像など）

### 2. 服薬管理
- **服薬記録の登録・更新・取得**
- **出血状態の記録**（`hasBleeding`フラグ）
- **服薬ステータス計算**
  - 現在の連続服用日数
  - 休薬期間の判定（4日間）
  - 連続出血日数の計算

### 3. 通知システム
- **Web Push通知**による服薬リマインダー
- **通知設定の管理**（プラットフォーム別）
- **重複送信防止**（5分間の制限）
- **サブスクリプション管理**

### 4. API エンドポイント

#### 認証関連
- `GET /api/auth/google` - Google OAuth認証開始
- `GET /api/auth/callback/google` - OAuthコールバック処理
- `GET /api/auth/session` - セッション情報取得
- `POST /api/auth/signout` - サインアウト

#### 服薬管理
- `GET /api/medication-status` - 服薬ステータス取得（認証必須）
- `POST /api/medication-log` - 服薬記録登録（認証必須）
- `GET /api/medication-log` - 服薬記録一覧取得（認証必須）
- `GET /api/medication-log/:id` - 特定の服薬記録取得（認証必須）
- `PATCH /api/medication-log/:id` - 服薬記録更新（認証必須）

#### 通知管理
- `POST /api/notification` - 通知送信
- `GET /api/notification/setting` - 通知設定取得（認証必須）
- `POST /api/notification/setting` - 通知設定登録（認証必須）

#### ヘルスチェック
- `GET /api/health` - ヘルスチェック

## データモデル

### User
```go
type User struct {
    ID            string    `json:"id" gorm:"primary_key"`
    Name          string    `json:"name"`
    Email         string    `json:"email" gorm:"unique"`
    EmailVerified bool      `json:"emailVerified"`
    Image         *string   `json:"image"`
    CreatedAt     time.Time `json:"createdAt"`
    UpdatedAt     time.Time `json:"updatedAt"`
}
```

### MedicationLog
```go
type MedicationLog struct {
    ID          uint       `json:"id" gorm:"primarykey"`
    UserID      string     `json:"userId" gorm:"not null;index:idx_user_id"`
    HasBleeding bool       `json:"hasBleeding" gorm:"default:false"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   time.Time  `json:"updatedAt"`
    DeletedAt   *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}
```

### NotificationSetting
```go
type NotificationSetting struct {
    ID           uint           `json:"id" gorm:"primarykey"`
    UserID       string         `json:"userId" gorm:"not null;index:idx_user_platform,unique:true,part:1"`
    Platform     string         `json:"platform" gorm:"not null;index:idx_user_platform,unique:true,part:2"`
    IsEnabled    bool           `json:"isEnabled" gorm:"default:true"`
    Subscription string         `json:"subscription" gorm:"type:text"`
    CreatedAt    time.Time      `json:"createdAt"`
    UpdatedAt    time.Time      `json:"updatedAt"`
    DeletedAt    gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}
```

## 開発環境

### 前提条件
- Go 1.24以上
- PostgreSQL
- Air（ホットリロード用）

### セットアップ

1. **依存関係のインストール**
   ```bash
   make install-deps
   ```

2. **Airのインストール**
   ```bash
   make install-air
   ```

3. **環境変数の設定**
   ```bash
   # .envファイルを作成
   DATABASE_URL="postgres://username:password@localhost:5432/dbname"
   GOOGLE_CLIENT_ID="your_google_client_id"
   APP_URL="http://localhost:3000"
   ```

4. **開発サーバーの起動**
   ```bash
   make dev  # ホットリロード有効
   ```

### 開発コマンド

```bash
make dev      # 開発モード（ホットリロード）
make build    # ビルド
make run      # ビルドして実行
make test     # テスト実行
make clean    # クリーンアップ
```

## 設計パターン

### 1. レイヤードアーキテクチャ
- **Handler層**: HTTPリクエストの処理
- **Service層**: ビジネスロジック
- **Repository層**: データアクセス
- **Model層**: データ構造定義

### 2. 依存性注入
- 各レイヤー間の依存関係を明示的に注入
- テスト時のモック化が容易
- コンポーネントの再利用性向上

### 3. ミドルウェアパターン
- 認証、CORS、ログ出力などの横断的関心事を分離
- チェーン形式で処理を組み合わせ

### 4. DTOパターン
- APIリクエスト・レスポンスの構造を明確化
- バリデーションルールの定義

## セキュリティ

### 認証・認可
- **JWTトークン**によるセッション管理
- **Google OAuth 2.0**による安全な認証
- **ミドルウェア**による認証必須エンドポイントの保護

### データ保護
- **環境変数**による機密情報の管理
- **SQLインジェクション対策**（GORMによるパラメータ化クエリ）
- **CORS設定**による適切なアクセス制御

## パフォーマンス・スケーラビリティ

### データベース最適化
- **インデックス**の適切な設定
- **GORM**による効率的なクエリ生成
- **自動マイグレーション**によるスキーマ管理

### 通知システム
- **重複送信防止**による無駄な処理の削減
- **非同期処理**によるレスポンス時間の短縮
- **サブスクリプション管理**による効率的な通知配信

## テスト

### テストフレームワーク
- **testify**を使用したテスト
- **単体テスト**と**統合テスト**の両方をサポート

### テスト実行
```bash
make test  # 全テスト実行
go test -v ./internal/handler  # 特定パッケージのテスト
```

## デプロイメント

### 環境変数
- `DATABASE_URL`: PostgreSQL接続文字列
- `GOOGLE_CLIENT_ID`: Google OAuthクライアントID
- `APP_URL`: アプリケーションのベースURL

### ビルド
```bash
make build  # バイナリファイル生成
./bin/server  # サーバー起動
```

## 監視・ログ

### ログ出力
- **構造化ログ**による情報の整理
- **ミドルウェア**によるリクエスト・レスポンスの記録
- **エラーハンドリング**による適切なログ出力

### ヘルスチェック
- `/api/health`エンドポイントによる死活監視
- データベース接続状態の確認

## 今後の拡張性

### 機能拡張
- **複数プラットフォーム**対応（iOS、Android）
- **プッシュ通知**の高度化
- **服薬スケジュール**管理
- **副作用記録**機能

### 技術的改善
- **GraphQL**APIの導入
- **gRPC**によるマイクロサービス化
- **キャッシュ層**の追加（Redis等）
- **メトリクス収集**（Prometheus等）

## ライセンス

このプロジェクトは独自のライセンスの下で提供されています。

## コントリビューション

開発への参加を歓迎します。プルリクエストやイシューの報告をお待ちしています。

---

**注意**: このREADMEは、モノレポ移行のためのリポジトリ説明として作成されています。
