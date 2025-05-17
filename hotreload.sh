#!/bin/bash

# ホットリロードスクリプト
# 監視対象ディレクトリ
WATCH_DIRS="cmd internal pkg migrations"

# ビルドコマンド
BUILD_CMD="go build -o ./tmp/main ./cmd/server"

# 実行コマンド
RUN_CMD="./tmp/main"

# 一時ディレクトリの作成
mkdir -p tmp

echo "ホットリロードを開始します..."
echo "監視対象ディレクトリ: $WATCH_DIRS"

# サーバープロセスID
SERVER_PID=""

# サーバーの起動
function start_server() {
  if [ -n "$SERVER_PID" ]; then
    echo "サーバーを再起動します..."
    kill -TERM $SERVER_PID 2>/dev/null
    wait $SERVER_PID 2>/dev/null
  else
    echo "サーバーを起動します..."
  fi
  
  # ビルド
  echo "ビルド中..."
  $BUILD_CMD
  
  if [ $? -eq 0 ]; then
    # 成功したらサーバーを開始
    $RUN_CMD &
    SERVER_PID=$!
    echo "サーバーが起動しました (PID: $SERVER_PID)"
  else
    echo "ビルドに失敗しました"
  fi
}

# 初回起動
start_server

# ファイル監視
while true; do
  # 変更を検知する
  echo "ファイルの変更を監視しています..."
  
  CHANGED_FILES=$(find $WATCH_DIRS -type f -name "*.go" -newer ./tmp/main 2>/dev/null)
  
  if [ -n "$CHANGED_FILES" ]; then
    echo "変更が検出されました:"
    echo "$CHANGED_FILES"
    start_server
  fi
  
  # 1秒待機
  sleep 1
done
