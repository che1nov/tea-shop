#!/bin/bash

# Скрипт для тестирования gRPC API Users Service
# Использование: ./test_api.sh

SERVER="localhost:8001"

echo "═══════════════════════════════════════════════════════════════"
echo "  Тестирование gRPC API Users Service"
echo "═══════════════════════════════════════════════════════════════"
echo ""

echo "1. Проверка доступности сервиса..."
if grpcurl -plaintext ${SERVER} list > /dev/null 2>&1; then
    echo "   ✓ Сервис доступен"
else
    echo "   ✗ Сервис недоступен! Убедитесь, что сервис запущен."
    exit 1
fi

echo ""
echo "2. Список доступных сервисов:"
grpcurl -plaintext ${SERVER} list

echo ""
echo "3. Методы UsersService:"
grpcurl -plaintext ${SERVER} list pb.UsersService

echo ""
echo "4. Описание метода CreateUser:"
grpcurl -plaintext ${SERVER} describe pb.UsersService.CreateUser

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  Тестирование методов"
echo "═══════════════════════════════════════════════════════════════"
echo ""

echo "5. Тест: CreateUser (создание нового пользователя)"
RESULT=$(grpcurl -plaintext -d '{
  "email": "test@example.com",
  "name": "Test User",
  "password": "testpass123"
}' ${SERVER} pb.UsersService/CreateUser 2>&1)
echo "$RESULT"
USER_ID=$(echo "$RESULT" | grep -o '"id": [0-9]*' | grep -o '[0-9]*' | head -1)

if [ -n "$USER_ID" ]; then
    echo ""
    echo "6. Тест: GetUser (получение созданного пользователя id=$USER_ID)"
    grpcurl -plaintext -d "{\"user_id\": $USER_ID}" ${SERVER} pb.UsersService/GetUser
    
    echo ""
    echo "7. Тест: CreateUser (попытка создать дубликат email)"
    grpcurl -plaintext -d '{
      "email": "test@example.com",
      "name": "Duplicate User",
      "password": "pass123"
    }' ${SERVER} pb.UsersService/CreateUser
    
    echo ""
    echo "8. Тест: GetUser (несуществующий пользователь)"
    grpcurl -plaintext -d '{
      "user_id": 99999
    }' ${SERVER} pb.UsersService/GetUser
else
    echo "   Не удалось создать пользователя для дальнейших тестов"
fi

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  Тестирование завершено"
echo "═══════════════════════════════════════════════════════════════"

