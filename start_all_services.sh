#!/usr/bin/env bash
set -euo pipefail

# Скрипт запускает все backend-сервисы в фоне и сохраняет PID/логи.

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RUN_DIR="$BASE_DIR/.run"
LOG_DIR="$RUN_DIR/logs"
PID_DIR="$RUN_DIR/pids"

GOCACHE_DIR="$BASE_DIR/.gocache"
GOMODCACHE_DIR="$BASE_DIR/.gomodcache"

mkdir -p "$LOG_DIR" "$PID_DIR" "$GOCACHE_DIR" "$GOMODCACHE_DIR"

SERVICES=(
  "users-service"
  "goods-service"
  "payment-service"
  "delivery-service"
  "order-service"
  "notify-service"
  "api-gateway"
)

if ! command -v go >/dev/null 2>&1; then
  echo "Ошибка: go не найден в PATH"
  exit 1
fi

start_service() {
  local service="$1"
  local service_dir="$BASE_DIR/$service"
  local pid_file="$PID_DIR/$service.pid"
  local log_file="$LOG_DIR/$service.log"

  if [ ! -d "$service_dir" ]; then
    echo "[WARN] Пропуск $service: директория не найдена"
    return 0
  fi

  if [ -f "$pid_file" ]; then
    local existing_pid
    existing_pid="$(cat "$pid_file")"
    if kill -0 "$existing_pid" >/dev/null 2>&1; then
      echo "[INFO] $service уже запущен (PID $existing_pid)"
      return 0
    fi
    rm -f "$pid_file"
  fi

  echo "[INFO] Запуск $service..."
  (
    cd "$service_dir"
    GOCACHE="$GOCACHE_DIR" GOMODCACHE="$GOMODCACHE_DIR" \
      nohup go run ./cmd/main.go >"$log_file" 2>&1 &
    echo $! >"$pid_file"
  )

  sleep 1

  local pid
  pid="$(cat "$pid_file")"
  if kill -0 "$pid" >/dev/null 2>&1; then
    echo "[OK] $service запущен (PID $pid), лог: $log_file"
  else
    echo "[ERROR] Не удалось запустить $service. Последние строки лога:"
    tail -n 20 "$log_file" || true
    return 1
  fi
}

for service in "${SERVICES[@]}"; do
  start_service "$service"
done

echo
echo "Все сервисы обработаны."
echo "API Gateway: http://localhost:8080"
echo "Swagger: http://localhost:8080/swagger/index.html"
echo "Остановка: ./stop_all_services.sh"
