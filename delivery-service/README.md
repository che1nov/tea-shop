# Delivery Service

Сервис управления доставкой заказов.

## Описание

Сервис отвечает за:
- Создание доставок для заказов
- Управление статусами доставок
- Отслеживание адресов доставки

## Технологии

- **Go 1.25+**
- **gRPC** - для межсервисного взаимодействия
- **PostgreSQL** - база данных (deliveries_db)
- **Prometheus** - метрики

## Архитектура

Сервис следует принципам Clean Architecture:

```
Handler (gRPC) → Service → Repository → Database
```

### Слои

- **Handler** (`internal/handler/`) - обработка gRPC запросов
- **Service** (`internal/service/`) - бизнес-логика доставок
- **Repository** (`internal/repository/`) - работа с БД
- **Model** (`internal/model/`) - доменные модели

## API Документация

API описан в proto файле: `shared/pb/delivery.proto`

Для тестирования gRPC API используйте:
- **grpcui** (веб-интерфейс): `grpcui -plaintext localhost:8005`
- `grpcurl` (CLI)
- BloomRPC (GUI)

## gRPC API

### Методы

#### CreateDelivery
Создает доставку для заказа.

**Request:**
```protobuf
message CreateDeliveryRequest {
  int64 order_id = 1;
  string address = 2;
}
```

**Response:**
```protobuf
message Delivery {
  int64 id = 1;
  int64 order_id = 2;
  string address = 3;
  string status = 4;
  int64 created_at = 5;
  int64 updated_at = 6;
}
```

#### GetDelivery
Получает информацию о доставке по ID.

**Request:**
```protobuf
message GetDeliveryRequest {
  int64 delivery_id = 1;
}
```

#### UpdateDeliveryStatus
Обновляет статус доставки.

**Request:**
```protobuf
message UpdateDeliveryStatusRequest {
  int64 delivery_id = 1;
  string status = 2;
}
```

## Структура базы данных

```sql
CREATE TABLE deliveries (
    id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    address TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_deliveries_order_id ON deliveries(order_id);
```

## Статусы доставки

- `pending` - доставка ожидает обработки
- `preparing` - заказ готовится к отправке
- `in_transit` - доставка в пути
- `delivered` - доставлено
- `cancelled` - доставка отменена

## Конфигурация

Переменные окружения:

- `DB_HOST` - хост БД (по умолчанию: localhost)
- `DB_PORT` - порт БД (по умолчанию: 5436)
- `DB_USER` - пользователь БД (по умолчанию: user)
- `DB_PASSWORD` - пароль БД (по умолчанию: password)
- `DB_NAME` - имя БД (по умолчанию: deliveries_db)

## Запуск

```bash
go run ./cmd/main.go
```

Сервис будет доступен на порту **8005** (gRPC).

Метрики Prometheus доступны на порту **9005**.

## Особенности реализации

1. **Транзакции**: Все операции выполняются в транзакциях
2. **Статусы**: Управление жизненным циклом доставки
3. **Адреса**: Хранение полных адресов доставки
4. **Интеграция**: Может вызываться из API Gateway или других сервисов

## Тестирование

```bash
go test ./...
```

## Зависимости

- **users-service** - нет
- **goods-service** - нет
- **order-service** - связан через order_id
- **payment-service** - нет

## Логирование

Используется структурированное логирование через `shared/pkg/logger`.

## Мониторинг

Метрики Prometheus доступны по адресу: `http://localhost:9005/metrics`

