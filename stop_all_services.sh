#!/bin/bash

echo "Stopping all services..."

# Остановка всех процессов
for service in users-service goods-service order-service payment-service delivery-service notify-service api-gateway; do
    if [ -f "/tmp/${service}.pid" ]; then
        PID=$(cat "/tmp/${service}.pid")
        if ps -p $PID > /dev/null 2>&1; then
            echo "Stopping $service (PID: $PID)..."
            kill $PID
            rm "/tmp/${service}.pid"
        fi
    fi
done

# Остановка Docker контейнеров
echo "Stopping Docker containers..."
docker-compose down

echo "All services stopped."