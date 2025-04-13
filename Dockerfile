# ベースイメージとして公式のGoイメージを使用
FROM golang:1.24

# コレがミソ
COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.8.4 /lambda-adapter /opt/extensions/lambda-adapter

# 作業ディレクトリを設定
WORKDIR /app

# Goモジュールを使用するためのファイルをコピー
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN go build -o /app

# コンテナが実行するコマンドを指定
CMD [ "/app" ]