# Payment Service

Сервис обработки платежей.

## Описание

Сервис отвечает за:
- Обработку платежей для заказов
- Управление статусами платежей
- Хранение истории платежей

## Технологии

- **Go 1.25+**
- **gRPC** - для межсервисного взаимодействия
- **PostgreSQL** - база данных (payments_db)
- **Prometheus** - метрики

## Архитектура

Сервис следует принципам Clean Architecture:

```
Handler (gRPC) → Service → Repository → Database
```

### Слои

- **Handler** (`internal/handler/`) - обработка gRPC запросов
- **Service** (`internal/service/`) - бизнес-логика обработки платежей
- **Repository** (`internal/repository/`) - работа с БД
- **Model** (`internal/model/`) - доменные модели

## API Документация

API описан в proto файле: `shared/pb/payments.proto`

Для тестирования gRPC API используйте:
- **grpcui** (веб-интерфейс): `grpcui -plaintext localhost:8004`
- `grpcurl` (CLI)
- BloomRPC (GUI)

## gRPC API

### Методы

#### ProcessPayment
Обрабатывает платеж для заказа.

**Request:**
```protobuf
message ProcessPaymentRequest {
  int64 order_id = 1;
  double amount = 2;
  string method = 3;
}
```

**Response:**
```protobuf
message Payment {
  int64 id = 1;
  int64 order_id = 2;
  double amount = 3;
  string status = 4;
  int64 created_at = 5;
  int64 updated_at = 6;
}
```

#### GetPayment
Получает информацию о платеже по ID.

**Request:**
```protobuf
message GetPaymentRequest {
  int64 payment_id = 1;
}
```

#### GetPaymentByOrderID
Получает информацию о платеже по ID заказа.

**Request:**
```protobuf
message GetPaymentByOrderIDRequest {
  int64 order_id = 1;
}
```

## Структура базы данных

```sql
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    method VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_order_id ON payments(order_id);
```

## Статусы платежей

- `pending` - платеж ожидает обработки
- `processing` - платеж обрабатывается
- `completed` - платеж успешно обработан
- `failed` - платеж не удался
- `refunded` - платеж возвращен

## Конфигурация

Переменные окружения:

- `DB_HOST` - хост БД (по умолчанию: localhost)
- `DB_PORT` - порт БД (по умолчанию: 5435)
- `DB_USER` - пользователь БД (по умолчанию: user)
- `DB_PASSWORD` - пароль БД (по умолчанию: password)
- `DB_NAME` - имя БД (по умолчанию: payments_db)

## Запуск

```bash
go run ./cmd/main.go
```

Сервис будет доступен на порту **8004** (gRPC).

Метрики Prometheus доступны на порту **9004**.

## Особенности реализации

1. **Транзакции**: Все операции выполняются в транзакциях
2. **Статусы**: Управление жизненным циклом платежа
3. **История**: Все платежи сохраняются для аудита
4. **Интеграция**: Вызывается из order-service при создании заказа

## Тестирование

```bash
go test ./...
```

## Зависимости

- **users-service** - нет
- **goods-service** - нет
- **order-service** - вызывается из order-service
- **delivery-service** - нет

## Логирование

Используется структурированное логирование через `shared/pkg/logger`.

## Мониторинг

Метрики Prometheus доступны по адресу: `http://localhost:9004/metrics`

