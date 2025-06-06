# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

- **Development with hot reload**: `make dev` (uses air for hot reloading)
- **Build**: `make build` (outputs to ./bin/server)
- **Run**: `make run` (builds then runs)
- **Test**: `make test` (runs all tests with -v flag)
- **Clean**: `make clean` (removes ./bin and ./tmp directories)

## Architecture Overview

This is a Go REST API backend for a medication tracking application (`okusuri-backend`).

### Project Structure
- **cmd/server/main.go**: Entry point - sets up DB, runs migrations, configures routes, starts server on :8080
- **internal/**: Core business logic organized by layer
  - **handler/**: HTTP request handlers (medication, notification)
  - **service/**: Business logic services
  - **repository/**: Data access layer
  - **model/**: Database models (MedicationLog, NotificationSetting, User)
  - **dto/**: Data transfer objects
  - **middleware/**: Auth, CORS, logging middleware
  - **routes.go**: Router setup with dependency injection pattern
- **pkg/config/**: Database configuration and connection management
- **migrations/**: Database migrations

### Key Technologies
- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM ORM
- **Environment**: godotenv for .env loading
- **Hot Reload**: Air (installed via make dev)

### Database Setup
Requires `DATABASE_URL` environment variable. The application automatically runs migrations on startup.

### API Structure
All endpoints are under `/api` prefix:
- `/api/health`: Health check
- `/api/notification`: Notification sending
- `/api/notification/setting`: Notification settings (auth required)
- `/api/medication-log`: Medication logging endpoints (auth required)
- `/api/medication-status`: Medication status (auth required)

### Authentication
Uses JWT-based authentication middleware for protected routes. User ID is extracted from JWT tokens.

### Core Models
- **MedicationLog**: Tracks medication intake with bleeding status
- **NotificationSetting**: Manages user notification preferences per platform
- **User**: User authentication and identification