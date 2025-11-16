# Goods Service

Сервис управления каталогом товаров и остатками.

## Описание

Сервис отвечает за:
- Управление каталогом товаров (CRUD операции)
- Управление остатками товаров
- Резервирование товаров при создании заказов
- Проверку наличия товаров на складе

## Технологии

- **Go 1.25+**
- **gRPC** - для межсервисного взаимодействия
- **PostgreSQL** - база данных (goods_db)
- **Prometheus** - метрики

## Архитектура

Сервис следует принципам Clean Architecture:

```
Handler (gRPC) → Service → Repository → Database
```

### Слои

- **Handler** (`internal/handler/`) - обработка gRPC запросов
- **Service** (`internal/service/`) - бизнес-логика
- **Repository** (`internal/repository/`) - работа с БД
- **Model** (`internal/model/`) - доменные модели

## API Документация

API описан в proto файле: `shared/pb/goods.proto`

Для тестирования gRPC API используйте:
- **grpcui** (веб-интерфейс): `grpcui -plaintext localhost:8002`
- `grpcurl` (CLI)
- BloomRPC (GUI)

## gRPC API

### Методы

#### CreateGood
Создает новый товар в каталоге.

**Request:**
```protobuf
message CreateGoodRequest {
  string name = 1;
  string description = 2;
  double price = 3;
  int32 stock = 4;
}
```

**Response:**
```protobuf
message Good {
  int64 id = 1;
  string name = 2;
  string description = 3;
  double price = 4;
  int32 stock = 5;
  int64 created_at = 6;
}
```

#### GetGood
Получает информацию о товаре по ID.

**Request:**
```protobuf
message GetGoodRequest {
  int64 good_id = 1;
}
```

#### ListGoods
Получает список товаров с пагинацией.

**Request:**
```protobuf
message ListGoodsRequest {
  int32 limit = 1;
  int32 offset = 2;
}
```

**Response:**
```protobuf
message ListGoodsResponse {
  repeated Good goods = 1;
  int32 total = 2;
}
```

#### UpdateGood
Обновляет информацию о товаре.

**Request:**
```protobuf
message UpdateGoodRequest {
  int64 id = 1;
  string name = 2;
  string description = 3;
  double price = 4;
  int32 stock = 5;
}
```

#### DeleteGood
Удаляет товар из каталога. Автоматически удаляет связанные резервации.

**Request:**
```protobuf
message DeleteGoodRequest {
  int64 good_id = 1;
}
```

**Response:**
```protobuf
message DeleteGoodResponse {
  bool success = 1;
  string message = 2;
}
```

#### CheckStock
Проверяет наличие товара на складе.

**Request:**
```protobuf
message CheckStockRequest {
  int64 good_id = 1;
  int32 quantity = 2;
}
```

**Response:**
```protobuf
message CheckStockResponse {
  bool available = 1;
}
```

#### ReserveStock
Резервирует товар для заказа.

**Request:**
```protobuf
message ReserveStockRequest {
  int64 good_id = 1;
  int32 quantity = 2;
  int64 order_id = 3;
}
```

**Response:**
```protobuf
message ReserveStockResponse {
  bool success = 1;
  string error = 2;
}
```

## Структура базы данных

```sql
CREATE TABLE goods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE stock_reservations (
    id SERIAL PRIMARY KEY,
    good_id INTEGER NOT NULL REFERENCES goods(id) ON DELETE CASCADE,
    order_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Конфигурация

Переменные окружения:

- `DB_HOST` - хост БД (по умолчанию: localhost)
- `DB_PORT` - порт БД (по умолчанию: 5433)
- `DB_USER` - пользователь БД (по умолчанию: user)
- `DB_PASSWORD` - пароль БД (по умолчанию: password)
- `DB_NAME` - имя БД (по умолчанию: goods_db)

## Запуск

```bash
go run ./cmd/main.go
```

Сервис будет доступен на порту **8002** (gRPC).

Метрики Prometheus доступны на порту **9002**.

## Особенности реализации

1. **Резервирование товаров**: При создании заказа товары резервируются в таблице `stock_reservations`
2. **Удаление товаров**: Сначала удаляются все резервации, затем товар (транзакция)
3. **Проверка остатков**: Учитываются зарезервированные товары
4. **Транзакции**: Все операции с остатками выполняются в транзакциях

## Тестирование

```bash
go test ./...
```

## Зависимости

- **users-service** - нет
- **order-service** - используется для резервирования товаров
- **payment-service** - нет
- **delivery-service** - нет

## Логирование

Используется структурированное логирование через `shared/pkg/logger`.

## Мониторинг

Метрики Prometheus доступны по адресу: `http://localhost:9002/metrics`

