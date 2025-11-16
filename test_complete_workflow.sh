#!/bin/bash

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Полная проверка всех сервисов ecommerce системы ===${NC}"
echo ""

# Функция для проверки сервиса
check_service() {
    local name=$1
    local port=$2
    
    if grpcurl -plaintext -connect-timeout 2 localhost:$port list > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC} $name (порт $port) - OK"
        return 0
    else
        echo -e "${RED}✗${NC} $name (порт $port) - FAILED"
        return 1
    fi
}

# Проверка gRPC сервисов
echo -e "${YELLOW}1. Проверка gRPC сервисов:${NC}"
check_service "Users Service" "8001"
check_service "Goods Service" "8002"
check_service "Order Service" "8003"
check_service "Payment Service" "8004"
check_service "Delivery Service" "8005"
echo ""

# Проверка API Gateway
echo -e "${YELLOW}2. Проверка API Gateway:${NC}"
if curl -s http://localhost:8080/api/v1/auth/register > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} API Gateway (порт 8080) - OK"
else
    echo -e "${RED}✗${NC} API Gateway (порт 8080) - FAILED"
    exit 1
fi
echo ""

# Проверка Docker контейнеров
echo -e "${YELLOW}3. Проверка Docker контейнеров:${NC}"
if docker ps | grep -q "postgres-users"; then
    echo -e "${GREEN}✓${NC} PostgreSQL контейнеры - OK"
else
    echo -e "${RED}✗${NC} PostgreSQL контейнеры - FAILED"
fi

if docker ps | grep -q "kafka"; then
    echo -e "${GREEN}✓${NC} Kafka - OK"
else
    echo -e "${RED}✗${NC} Kafka - FAILED"
fi

if docker ps | grep -q "zookeeper"; then
    echo -e "${GREEN}✓${NC} Zookeeper - OK"
else
    echo -e "${RED}✗${NC} Zookeeper - FAILED"
fi
echo ""

# Создание тестовых данных
echo -e "${YELLOW}4. Создание тестовых товаров в БД:${NC}"
PGPASSWORD=password psql -h localhost -p 5433 -U user -d goods_db -c "
INSERT INTO goods (name, description, price, stock, created_at, updated_at)
VALUES 
    ('Laptop', 'High performance laptop', 1299.99, 50, NOW(), NOW()),
    ('Mouse', 'Wireless mouse', 29.99, 100, NOW(), NOW()),
    ('Keyboard', 'Mechanical keyboard', 89.99, 75, NOW(), NOW())
ON CONFLICT DO NOTHING;
" > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} Товары созданы"
else
    echo -e "${RED}✗${NC} Ошибка создания товаров"
fi
echo ""

# Генерируем уникальный email для теста
TIMESTAMP=$(date +%s)
TEST_EMAIL="testuser${TIMESTAMP}@example.com"
TEST_PASSWORD="testpass123"

echo -e "${YELLOW}5. Тестирование полного workflow:${NC}"
echo "Используемый email: $TEST_EMAIL"
echo ""

# 5.1 Регистрация
echo "  5.1 Регистрация пользователя..."
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"name\": \"Test User\",
    \"password\": \"$TEST_PASSWORD\"
  }")

USER_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.id // empty')
if [ "$USER_ID" != "null" ] && [ ! -z "$USER_ID" ]; then
    echo -e "    ${GREEN}✓${NC} Пользователь создан (ID: $USER_ID)"
else
    echo -e "    ${RED}✗${NC} Ошибка регистрации"
    echo "$REGISTER_RESPONSE" | jq '.'
    exit 1
fi

# 5.2 Login
echo "  5.2 Получение JWT токена..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\"
  }")

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token // empty')
if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
    echo -e "    ${RED}✗${NC} Ошибка получения токена"
    echo "$LOGIN_RESPONSE" | jq '.'
    exit 1
fi
echo -e "    ${GREEN}✓${NC} Токен получен"

# 5.3 Получение информации о пользоватеle
echo "  5.3 Получение информации о пользователе..."
USER_INFO=$(curl -s -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN")
if echo "$USER_INFO" | jq -e '.id' > /dev/null 2>&1; then
    echo -e "    ${GREEN}✓${NC} Информация о пользователе получена"
else
    echo -e "    ${RED}✗${NC} Ошибка получения информации"
    exit 1
fi

# 5.4 Получение списка товаров
echo "  5.4 Получение списка товаров..."
GOODS_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/v1/goods?limit=10&offset=0" \
  -H "Authorization: Bearer $TOKEN")
GOOD_ID=$(echo "$GOODS_RESPONSE" | jq -r '.goods[0].id // empty')

if [ -z "$GOOD_ID" ] || [ "$GOOD_ID" == "null" ]; then
    echo -e "    ${RED}✗${NC} Товары не найдены"
    exit 1
fi
echo -e "    ${GREEN}✓${NC} Товары получены (ID первого товара: $GOOD_ID)"

# 5.5 Создание заказа
echo "  5.5 Создание заказа..."
ORDER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"items\": [{
      \"good_id\": $GOOD_ID,
      \"quantity\": 1,
      \"price\": 1299.99
    }]
  }")

ORDER_ID=$(echo "$ORDER_RESPONSE" | jq -r '.id // empty')
if [ "$ORDER_ID" != "null" ] && [ ! -z "$ORDER_ID" ]; then
    echo -e "    ${GREEN}✓${NC} Заказ создан (ID: $ORDER_ID)"
else
    echo -e "    ${RED}✗${NC} Ошибка создания заказа"
    echo "$ORDER_RESPONSE" | jq '.'
    exit 1
fi

# 5.6 Проверка платежа
echo "  5.6 Проверка платежа..."
sleep 2  # Даем время на обработку платежа
PAYMENT_RESPONSE=$(grpcurl -plaintext -d "{
  \"order_id\": $ORDER_ID
}" localhost:8004 pb.PaymentsService/GetPaymentByOrderID 2>&1)

PAYMENT_ID=$(echo "$PAYMENT_RESPONSE" | jq -r '.id // empty' 2>/dev/null)
if [ ! -z "$PAYMENT_ID" ] && [ "$PAYMENT_ID" != "null" ]; then
    echo -e "    ${GREEN}✓${NC} Платеж создан"
    echo "      Payment ID: $PAYMENT_ID"
else
    echo -e "    ${YELLOW}⚠${NC} Платеж не найден (возможно, еще обрабатывается)"
    echo "      Response: $PAYMENT_RESPONSE"
fi

# 5.7 Создание доставки
echo "  5.7 Создание доставки..."
DELIVERY_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/deliveries \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"order_id\": $ORDER_ID,
    \"address\": \"123 Main St, City, Country\"
  }")

DELIVERY_ID=$(echo "$DELIVERY_RESPONSE" | jq -r '.id // empty')
if [ "$DELIVERY_ID" != "null" ] && [ ! -z "$DELIVERY_ID" ]; then
    echo -e "    ${GREEN}✓${NC} Доставка создана (ID: $DELIVERY_ID)"
else
    echo -e "    ${RED}✗${NC} Ошибка создания доставки"
    echo "$DELIVERY_RESPONSE" | jq '.'
fi

# 5.8 Проверка Kafka событий
echo "  5.8 Проверка Kafka событий..."
sleep 1  # Даем время на обработку события
if tail -30 /tmp/notify-service.log 2>/dev/null | grep -qE "(Received event|order)"; then
    echo -e "    ${GREEN}✓${NC} Kafka события обрабатываются"
    tail -5 /tmp/notify-service.log 2>/dev/null | grep -E "(Received event|order)" || true
else
    echo -e "    ${YELLOW}⚠${NC} Kafka события не обнаружены (возможно, notify-service не запущен или события еще не обработаны)"
    echo "      Последние строки лога:"
    tail -5 /tmp/notify-service.log 2>/dev/null || echo "      Лог файл не найден"
fi

echo ""
echo -e "${GREEN}=== Все проверки завершены ===${NC}"
echo ""
echo "Созданные данные:"
echo "  User ID: $USER_ID"
echo "  Order ID: $ORDER_ID"
if [ ! -z "$DELIVERY_ID" ] && [ "$DELIVERY_ID" != "null" ]; then
    echo "  Delivery ID: $DELIVERY_ID"
fi
echo ""

