#!/bin/bash

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Starting all ecommerce services...${NC}"

# Проверка Docker
if ! docker ps > /dev/null 2>&1; then
    echo "Docker is not running. Please start Docker first."
    exit 1
fi

# Запуск инфраструктуры
echo -e "${GREEN}Starting infrastructure (Docker)...${NC}"
docker-compose up -d

# Ожидание готовности БД
echo -e "${YELLOW}Waiting for databases to be ready...${NC}"
sleep 5

# Функция для запуска сервиса в фоне
start_service() {
    local service_name=$1
    local service_path=$2
    local port=$3
    local env_vars="${4:-}"
    
    echo -e "${GREEN}Starting $service_name on port $port...${NC}"
    cd "$service_path"
    if [ -n "$env_vars" ]; then
        eval "$env_vars go run ./cmd/main.go > /tmp/${service_name}.log 2>&1 &"
    else
        go run ./cmd/main.go > "/tmp/${service_name}.log" 2>&1 &
    fi
    echo $! > "/tmp/${service_name}.pid"
    cd - > /dev/null
    sleep 2
}

# Запуск всех сервисов
start_service "users-service" "users-service" "8001" "ADMIN_EMAIL=admin@example.com ADMIN_PASSWORD=admin123"
start_service "goods-service" "goods-service" "8002"
start_service "order-service" "order-service" "8003"
start_service "payment-service" "payment-service" "8004"
start_service "delivery-service" "delivery-service" "8005"
start_service "notify-service" "notify-service" "8006"
start_service "api-gateway" "api-gateway" "8080"

echo -e "${GREEN}All services started!${NC}"
echo ""
echo "Services running:"
echo "  - Users Service:    http://localhost:8001"
echo "  - Goods Service:    http://localhost:8002"
echo "  - Order Service:    http://localhost:8003"
echo "  - Payment Service:  http://localhost:8004"
echo "  - Delivery Service: http://localhost:8005"
echo "  - Notify Service:   http://localhost:8006 (Kafka consumer)"
echo "  - API Gateway:      http://localhost:8080"
echo ""
echo "To stop all services:"
echo "  ./stop_all_services.sh"
echo ""
echo "To view logs:"
echo "  tail -f /tmp/users-service.log"