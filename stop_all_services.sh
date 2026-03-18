#!/usr/bin/env bash
set -euo pipefail

# Скрипт останавливает backend-сервисы, запущенные через start_all_services.sh.

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PID_DIR="$BASE_DIR/.run/pids"

if [ ! -d "$PID_DIR" ]; then
  echo "PID-директория не найдена: $PID_DIR"
  exit 0
fi

stop_pid_file() {
  local pid_file="$1"
  local service
  service="$(basename "$pid_file" .pid)"
  local pid
  pid="$(cat "$pid_file")"

  if ! kill -0 "$pid" >/dev/null 2>&1; then
    echo "[INFO] $service уже остановлен"
    rm -f "$pid_file"
    return 0
  fi

  echo "[INFO] Останавливаю $service (PID $pid)..."
  kill "$pid" >/dev/null 2>&1 || true

  for _ in {1..10}; do
    if ! kill -0 "$pid" >/dev/null 2>&1; then
      echo "[OK] $service остановлен"
      rm -f "$pid_file"
      return 0
    fi
    sleep 1
  done

  echo "[WARN] Принудительная остановка $service (PID $pid)"
  kill -9 "$pid" >/dev/null 2>&1 || true
  rm -f "$pid_file"
}

shopt -s nullglob
pid_files=("$PID_DIR"/*.pid)
shopt -u nullglob

if [ ${#pid_files[@]} -eq 0 ]; then
  echo "Нет запущенных сервисов (pid-файлы не найдены)."
  exit 0
fi

for pid_file in "${pid_files[@]}"; do
  stop_pid_file "$pid_file"
done

echo "Готово: сервисы остановлены."
