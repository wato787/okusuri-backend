# syntax=docker/dockerfile:1

FROM golang:1.23.2 AS builder
WORKDIR /app

# GOTOOLCHAIN を local に固定してコンテナ内で追加ダウンロードを防ぐ
ENV GOTOOLCHAIN=local
ENV CGO_ENABLED=0
ENV GOOS=linux

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags "-s -w" -o server ./cmd/server

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app

COPY --from=builder /app/server ./server

EXPOSE 8080
ENV PORT=8080
ENV GIN_MODE=release

ENTRYPOINT ["./server"]
