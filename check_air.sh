#!/bin/bash

# GOPATHを表示
echo "GOPATH: $GOPATH"

# GOBINを表示
echo "GOBIN: $GOBIN"

# 現在のPATHを表示
echo "PATH: $PATH"

# airコマンドの場所を探す
which air || echo "airコマンドが見つかりません"

# GOPATHのbinディレクトリにairが存在するか確認
if [ -f "$GOPATH/bin/air" ]; then
  echo "airが$GOPATH/bin/に存在します"
else
  echo "airが$GOPATH/bin/に存在しません"
fi

# ユーザーのホームディレクトリ内のgoディレクトリを確認
if [ -f "$HOME/go/bin/air" ]; then
  echo "airが$HOME/go/bin/に存在します"
else
  echo "airが$HOME/go/bin/に存在しません"
fi
