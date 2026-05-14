#!/usr/bin/env bash
set -euo pipefail

APP_DIR="${APP_DIR:-$HOME/tea-shop}"
REPO_URL="${REPO_URL:-https://github.com/che1nov/tea-shop.git}"
BRANCH="${BRANCH:-main}"

if ! command -v git >/dev/null 2>&1; then
  echo "[ERROR] git is not installed"
  exit 1
fi

if ! command -v docker >/dev/null 2>&1; then
  echo "[ERROR] docker is not installed"
  exit 1
fi

if [ ! -d "$APP_DIR/.git" ]; then
  git clone --branch "$BRANCH" "$REPO_URL" "$APP_DIR"
fi

cd "$APP_DIR"

git fetch origin "$BRANCH"
git checkout "$BRANCH"
git pull --ff-only origin "$BRANCH"

if [ ! -f .env ]; then
  echo "[ERROR] $APP_DIR/.env is missing"
  echo "Create it from .env.example and set production secrets before deploy."
  exit 1
fi

docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build --remove-orphans
docker compose -f docker-compose.yml -f docker-compose.prod.yml ps
