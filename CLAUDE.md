# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GopherAI — AI chat assistant with image recognition. Go backend (Gin) + Vue 3 frontend.

## Tech Stack

- **Backend**: Go 1.24, Gin, GORM (MySQL), CloudWeGo Eino (AI SDK), Redis, RabbitMQ
- **Frontend**: Vue 3, Element Plus, Vue Router, Axios, Vue CLI
- **Auth**: JWT (golang-jwt)

## Architecture

```
main.go → router → controller → service → dao (GORM) → MySQL
                  → middleware (JWT auth)
                  → common/aihelper (AI model abstraction)
                  → common/rabbitmq (async message persistence)
```

### AI Layer (common/aihelper)

- `AIModel` interface (model.go) with `GenerateResponse`/`StreamResponse`
- Implementations: OpenAIModel (env: OPENAI_API_KEY, OPENAI_MODEL_NAME, OPENAI_BASE_URL) and OllamaModel
- `AIFactory` creates models by type code ("1"=OpenAI, "2"=Ollama)
- `AIHelperManager` manages user→session→AIHelper mappings (singleton)
- `AIHelper` wraps model + message history + save callback (async via RabbitMQ)

### Data Flow

1. User sends message → controller → service → AIHelper.GenerateResponse/StreamResponse
2. Message history stored in memory, persisted async via RabbitMQ → MySQL
3. Stream responses use callback pattern for real-time frontend push

## Key Commands

```bash
# Backend
go run main.go                        # Start server (config from config/config.toml)
go build -o gopherai .                # Build binary
go mod tidy                           # Sync dependencies

# Frontend
cd vue-frontend && npm install        # Install deps
cd vue-frontend && npm run serve      # Dev server (port 8080)
cd vue-frontend && npm run build      # Production build
```

## Configuration

`config/config.toml` — MySQL, Redis, RabbitMQ, JWT, Email, server port/host.
AI model env vars: `OPENAI_API_KEY`, `OPENAI_MODEL_NAME`, `OPENAI_BASE_URL`.

## API Endpoints

All under `/api/v1`:
- `/user/*` — register/login (no auth)
- `/AI/*` — AI chat (JWT required)
- `/image/*` — image recognition (JWT required)
