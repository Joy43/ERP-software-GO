#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

export PATH="$(go env GOPATH)/bin:$PATH"

export APP_ENV="${APP_ENV:-development}"
export APP_PORT="${APP_PORT:-8080}"
export APP_NAME="${APP_NAME:-assmi-super-shop-erp-backend}"

echo "[1/5] Downloading Go modules..."
go mod tidy

echo "[2/5] Ensuring swag CLI exists..."
if ! command -v swag >/dev/null 2>&1; then
  go install github.com/swaggo/swag/cmd/swag@latest
fi

echo "[3/5] Generating Swagger docs..."
swag init -g cmd/server/main.go --output docs/swagger

echo "[4/5] Ensuring air CLI exists for hot reload..."
if ! command -v air >/dev/null 2>&1; then
  go install github.com/air-verse/air@latest
fi

echo "[5/5] Starting hot reload server on :${APP_PORT}..."
air -c .air.toml
