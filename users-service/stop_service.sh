#!/bin/bash

# Скрипт для остановки Users Service
# Использование: ./stop_service.sh

PORT=8001

echo "Остановка сервиса на порту $PORT..."

# Находим процесс на порту
PID=$(lsof -ti :$PORT)

if [ -z "$PID" ]; then
    echo "Сервис на порту $PORT не запущен"
    exit 0
fi

echo "Найден процесс PID: $PID"
kill $PID

sleep 1

# Проверяем, остановлен ли процесс
if lsof -ti :$PORT > /dev/null 2>&1; then
    echo "Принудительная остановка..."
    kill -9 $PID
    sleep 1
fi

if lsof -ti :$PORT > /dev/null 2>&1; then
    echo "Ошибка: не удалось остановить процесс"
    exit 1
else
    echo "✓ Сервис остановлен"
fi

