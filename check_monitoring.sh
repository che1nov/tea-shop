#!/bin/bash

echo "=== Проверка сервисов и метрик ==="
echo ""

all_ok=true

check_service() {
  local port=$1
  local service_name=$2
  
  if curl -s --max-time 2 "http://localhost:$port/metrics" > /dev/null 2>&1; then
    echo "✓ $service_name (порт $port) - OK"
  else
    echo "✗ $service_name (порт $port) - НЕ ДОСТУПЕН"
    all_ok=false
  fi
}

check_service "9001" "users-service"
check_service "9002" "goods-service"
check_service "9003" "order-service"
check_service "9004" "payment-service"
check_service "9005" "delivery-service"
check_service "9006" "notify-service"
check_service "9007" "api-gateway"

echo ""
echo "=== Статус Prometheus targets ==="
if curl -s --max-time 2 "http://localhost:9090/api/v1/targets" > /dev/null 2>&1; then
  curl -s http://localhost:9090/api/v1/targets | jq -r '.data.activeTargets[] | "\(.labels.job): \(.health)"' 2>/dev/null | sort || echo "Ошибка получения данных из Prometheus"
else
  echo "Prometheus не доступен (проверьте: docker ps | grep prometheus)"
fi

echo ""
if [ "$all_ok" = true ]; then
  echo "✓ Все сервисы запущены и доступны для мониторинга"
else
  echo "⚠ Некоторые сервисы не запущены. Запустите: ./start_all_services.sh"
fi

