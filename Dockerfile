# syntax=docker/dockerfile:1

FROM golang:1.22.2 AS builder
WORKDIR /app

# Go 1.24 ツールチェーンが未提供のため、ローカルツールチェーンを強制
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
