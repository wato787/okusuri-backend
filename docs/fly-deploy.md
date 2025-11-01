# Fly.io デプロイガイド

本ドキュメントは Render から Fly.io へ移行する際の作業手順と運用ポイントをまとめています。`fly.toml` および `Dockerfile` はリポジトリ直下に配置済みです。

## 事前準備
- Fly.io アカウントの作成
- mise CLI のインストール（`curl https://mise.run | sh` を実行し、シェルを再読み込み）
- ツールのセットアップ
  ```bash
  mise install
  ```
- Fly.io で利用するリージョンの決定（本プロジェクトでは `nrt` を想定）

## 初期セットアップ
1. `mise run fly-login` で Fly.io にログインします。
2. `fly.toml` の `app` 名が一意であることを確認し、必要に応じて変更します。
3. 既存構成を用いてアプリを作成します。
   ```bash
   mise run fly-launch
   ```
   `--copy-config` オプションを利用しているため、`fly.toml` の内容がそのまま反映されます。

## シークレット設定
アプリ起動時に必要となる環境変数は Fly.io シークレットとして登録します。最低限以下を設定してください。

- `DATABASE_URL`: PostgreSQL 接続情報（Fly Postgres または外部サービス）
- `GOOGLE_CLIENT_ID`: Google OAuth 用クライアントID
- `GOOGLE_CLIENT_SECRET`: Google OAuth 用クライアントシークレット
- `APP_URL`: バックエンドの公開URL（例: `https://<app>.fly.dev`）
- `FRONTEND_URL`: フロントエンドの公開URL
- `VAPID_PUBLIC_KEY`: Web Push VAPID 公開鍵
- `VAPID_PRIVATE_KEY`: Web Push VAPID 秘密鍵

```bash
SECRET="DATABASE_URL=postgres://... APP_URL=https://..." mise run fly-secrets
# もしくは直接 fly secrets set KEY=VALUE ... を実行
```

## デプロイ
```bash
mise run fly-deploy
```

`fly.toml` に設定されている `min_machines_running = 1` により、最低1台の Machine が常時稼働しコールドスタートを防ぎます。ローリングデプロイ戦略を採用しているため、リリース時のダウンタイムを最小化できます。

## 運用コマンド例
- ステータス確認: `mise run fly-status`
- ログ監視: `fly logs`
- スケール調整: `fly scale vm shared-cpu-1x`、`fly scale memory 512`
- 強制再起動: `fly machine restart <machine_id>`

## 監視とヘルスチェック
`fly.toml` の `[[checks]]` により `GET /api/health` の HTTP チェックが 30 秒間隔で実行されます。失敗が続く場合は Machine が自動的に再起動します。アプリ側でも `/api/health` が DB 接続状態を返すことを前提にしています。

## データベースについて
- Fly Postgres を利用する場合は `fly pg create` でクラスタを作成し、`fly pg attach` でアプリに接続します。
- 外部の PostgreSQL（例: Neon, Supabase）を利用する場合は、VPC Peering もしくはパブリックエンドポイント経由で接続し、`DATABASE_URL` をそのエンドポイントに設定してください。
- マイグレーションはアプリ起動時に自動実行されますが、大規模な変更を行う際はメンテナンスウィンドウを確保してください。

## トラブルシューティング
- **`fly deploy` 時にビルドが失敗する**: ローカルで Docker が利用できる場合は `docker build .` で検証し、`GOTOOLCHAIN=local` が有効になっているか確認してください。
- **環境変数の設定漏れで起動しない**: `fly logs` を確認し、足りないシークレットを `fly secrets set` で追加してください。
- **コールドスタートが発生する**: `fly scale machine` で `min_machines_running` が 1 以上になっているかを確認します。`fly toml` を変更した際は `fly deploy` で再適用が必要です。

---

これらの手順により、Render でのコールドスタート問題を回避しつつ Fly.io 上で安定した運用が可能になります。
