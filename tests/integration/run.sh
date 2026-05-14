#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

cd "$ROOT_DIR"

DOCKER_BIN="${DOCKER_BIN:-docker}"
GO_BIN="${GO_BIN:-go}"

if ! command -v "$DOCKER_BIN" >/dev/null 2>&1 && [ -x /usr/local/bin/docker ]; then
  DOCKER_BIN="/usr/local/bin/docker"
fi

if ! command -v "$GO_BIN" >/dev/null 2>&1 && [ -x /opt/homebrew/bin/go ]; then
  GO_BIN="/opt/homebrew/bin/go"
fi

"$DOCKER_BIN" compose up -d \
  postgres-users postgres-goods postgres-orders postgres-payments postgres-delivery

cleanup() {
  if [ "${INTEGRATION_CLEANUP:-0}" = "1" ]; then
    "$DOCKER_BIN" compose down --remove-orphans
  fi
}
trap cleanup EXIT

wait_for_postgres() {
  local container="$1"
  local db="$2"

  echo "[INFO] Waiting for $container..."
  for _ in $(seq 1 30); do
    if "$DOCKER_BIN" compose exec -T "$container" pg_isready -U user -d "$db" >/dev/null 2>&1; then
      echo "[OK] $container is ready"
      return 0
    fi
    sleep 1
  done

  echo "[ERROR] $container is not ready"
  "$DOCKER_BIN" compose logs "$container" || true
  return 1
}

wait_for_postgres postgres-users users_db
wait_for_postgres postgres-goods goods_db
wait_for_postgres postgres-orders orders_db
wait_for_postgres postgres-payments payments_db
wait_for_postgres postgres-delivery deliveries_db

export RUN_INTEGRATION_TESTS=1

run_repo_tests() {
  local service="$1"

  echo "[INFO] Running integration tests for $service"
  (
    cd "$ROOT_DIR/$service"
    "$GO_BIN" test -tags=integration ./internal/repository
  )
}

run_repo_tests users-service
run_repo_tests goods-service
run_repo_tests order-service
run_repo_tests payment-service
run_repo_tests delivery-service

echo "[OK] Integration tests passed"
