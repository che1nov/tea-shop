# Order Service

Сервис управления заказами.

## Описание

Сервис отвечает за:
- Создание заказов
- Управление статусами заказов
- Резервирование товаров через goods-service
- Создание платежей через payment-service
- Публикацию событий в Kafka для уведомлений

## Технологии

- **Go 1.25+**
- **gRPC** - для межсервисного взаимодействия
- **PostgreSQL** - база данных (orders_db)
- **Kafka** - для публикации событий
- **Prometheus** - метрики

## Архитектура

Сервис следует принципам Clean Architecture:

```
Handler (gRPC) → Service → Repository → Database
                    ↓
              Kafka Producer
```

### Слои

- **Handler** (`internal/handler/`) - обработка gRPC запросов
- **Service** (`internal/service/`) - бизнес-логика, интеграция с другими сервисами
- **Repository** (`internal/repository/`) - работа с БД
- **Model** (`internal/model/`) - доменные модели
- **Kafka** (`internal/kafka/`) - producer для публикации событий

## API Документация

API описан в proto файле: `shared/pb/orders.proto`

Для тестирования gRPC API используйте:
- **grpcui** (веб-интерфейс): `grpcui -plaintext localhost:8003`
- `grpcurl` (CLI)
- BloomRPC (GUI)

## gRPC API

### Методы

#### CreateOrder
Создает новый заказ. При создании:
1. Проверяет наличие товаров через goods-service
2. Резервирует товары
3. Создает платеж через payment-service
4. Публикует событие в Kafka

**Request:**
```protobuf
message CreateOrderRequest {
  int64 user_id = 1;
  repeated OrderItem items = 2;
}

message OrderItem {
  int64 good_id = 1;
  int32 quantity = 2;
  double price = 3;
}
```

**Response:**
```protobuf
message Order {
  int64 id = 1;
  int64 user_id = 2;
  repeated OrderItem items = 3;
  string status = 4;
  double total_price = 5;
  int64 created_at = 6;
  int64 updated_at = 7;
}
```

#### GetOrder
Получает информацию о заказе по ID.

**Request:**
```protobuf
message GetOrderRequest {
  int64 order_id = 1;
}
```

#### UpdateOrderStatus
Обновляет статус заказа.

**Request:**
```protobuf
message UpdateOrderStatusRequest {
  int64 order_id = 1;
  string status = 2;
}
```

Возможные статусы:
- `pending` - заказ создан, ожидает обработки
- `processing` - заказ обрабатывается
- `completed` - заказ выполнен
- `cancelled` - заказ отменен

## Структура базы данных

```sql
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    total_price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    good_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Конфигурация

Переменные окружения:

- `DB_HOST` - хост БД (по умолчанию: localhost)
- `DB_PORT` - порт БД (по умолчанию: 5434)
- `DB_USER` - пользователь БД (по умолчанию: user)
- `DB_PASSWORD` - пароль БД (по умолчанию: password)
- `DB_NAME` - имя БД (по умолчанию: orders_db)
- `KAFKA_BROKERS` - адреса брокеров Kafka (по умолчанию: localhost:9092)

## Запуск

```bash
go run ./cmd/main.go
```

Сервис будет доступен на порту **8003** (gRPC).

Метрики Prometheus доступны на порту **9003**.

## Интеграции

### Goods Service
- `CheckStock` - проверка наличия товаров
- `ReserveStock` - резервирование товаров

### Payment Service
- `ProcessPayment` - создание платежа для заказа

### Kafka
Публикует события в топик `order_created`:
```json
{
  "order_id": 1,
  "user_id": 1,
  "total_price": 599.98,
  "status": "pending"
}
```

## Особенности реализации

1. **Транзакции**: Все операции с заказом выполняются в транзакциях
2. **Резервирование**: Товары резервируются перед созданием платежа
3. **События**: После создания заказа публикуется событие в Kafka
4. **Интеграция**: Синхронные вызовы к goods-service и payment-service

## Тестирование

```bash
go test ./...
```

## Зависимости

- **users-service** - для валидации пользователей (опционально)
- **goods-service** - для проверки и резервирования товаров
- **payment-service** - для создания платежей
- **Kafka** - для публикации событий

## Логирование

Используется структурированное логирование через `shared/pkg/logger`.

## Мониторинг

Метрики Prometheus доступны по адресу: `http://localhost:9003/metrics`

