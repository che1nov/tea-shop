# API Gateway

API Gateway для микросервисной платформы интернет-магазина чая.

## Описание

API Gateway является единой точкой входа для всех клиентских запросов. Он преобразует HTTP REST запросы в gRPC вызовы к соответствующим микросервисам.

## Функционал

- ✅ HTTP REST API
- ✅ JWT аутентификация
- ✅ CORS настройки
- ✅ Swagger документация
- ✅ Prometheus метрики
- ✅ Graceful shutdown

## Swagger документация

После запуска сервиса, Swagger UI доступен по адресу:
```
http://localhost:8080/swagger/index.html
```

### Генерация документации

После изменения API endpoints:

```bash
swag init -g cmd/main.go
```

## Конфигурация

Конфигурация находится в `config/config.go`. Переменные окружения:

- `API_GATEWAY_PORT` - порт API Gateway (по умолчанию 8080)
- `JWT_SECRET` - секретный ключ для JWT токенов
- `USERS_SERVICE` - адрес users-service (по умолчанию localhost:8001)
- `GOODS_SERVICE` - адрес goods-service (по умолчанию localhost:8002)
- `ORDERS_SERVICE` - адрес orders-service (по умолчанию localhost:8003)
- `PAYMENTS_SERVICE` - адрес payments-service (по умолчанию localhost:8004)
- `DELIVERY_SERVICE` - адрес delivery-service (по умолчанию localhost:8005)

## Запуск

```bash
go run ./cmd/main.go
```

## Структура

```
api-gateway/
├── cmd/
│   └── main.go          # Точка входа, настройка роутера
├── config/
│   └── config.go        # Конфигурация
├── internal/
│   ├── handler/
│   │   └── handler.go   # HTTP handlers
│   └── middleware/
│       └── auth.go       # JWT middleware
└── docs/
    ├── docs.go          # Сгенерированная Swagger документация
    ├── swagger.json     # JSON спецификация
    └── swagger.yaml     # YAML спецификация
```

