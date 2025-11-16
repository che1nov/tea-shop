#!/bin/bash

echo "Checking all services..."

# Проверка gRPC сервисов
services=("8001:Users" "8002:Goods" "8003:Order" "8004:Payment" "8005:Delivery")

for service in "${services[@]}"; do
    IFS=':' read -r port name <<< "$service"
    if grpcurl -plaintext -connect-timeout 2 localhost:$port list > /dev/null 2>&1; then
        echo "✓ $name Service (port $port) - OK"
    else
        echo "✗ $name Service (port $port) - FAILED"
    fi
done

# Проверка API Gateway
if curl -s http://localhost:8080/api/v1/auth/register > /dev/null 2>&1; then
    echo "✓ API Gateway (port 8080) - OK"
else
    echo "✗ API Gateway (port 8080) - FAILED"
fi

# Проверка Docker контейнеров
echo ""
echo "Docker containers:"
docker ps --format "table {{.Names}}\t{{.Status}}" | grep -E "postgres|kafka|zookeeper"