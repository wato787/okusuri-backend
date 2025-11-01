# Makefile

.PHONY: dev build run test clean install-deps install-air fly-launch fly-deploy fly-status fly-secrets

# Get Go bin directory path
GO_BIN := $(shell go env GOPATH)/bin

# デフォルトターゲット
all: dev

# 依存関係をインストール
install-deps:
	go mod download

# airをインストール
install-air:
	go install github.com/air-verse/air@latest

# 開発モードで実行（ホットリロード）
dev: install-air
	$(GO_BIN)/air

# ビルド
build:
	go build -o ./bin/server ./cmd/server

# 実行
run: build
	./bin/server

# テスト実行
test:
	go test -v ./...

# クリーンアップ
clean:
	rm -rf ./bin ./tmp

fly-launch:
	fly launch --copy-config --no-deploy --yes

fly-deploy:
	fly deploy --strategy rolling --wait-timeout 300

fly-status:
	fly status

fly-secrets:
	@test -n "$(SECRET)" || (echo "SECRET変数に設定値を指定してください (例: make fly-secrets SECRET=\"DATABASE_URL=...\")" && exit 1)
	fly secrets set $(SECRET)
