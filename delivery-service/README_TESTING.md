# Delivery Service - Тестирование

## Структура тестов

Delivery Service имеет полный набор тестов, покрывающий все слои архитектуры:

### 1. Unit тесты для Service (`internal/service/service_test.go`)

Тесты бизнес-логики с использованием моков для репозитория:

- `TestNew` - создание сервиса
- `TestCreateDelivery_Success` - успешное создание доставки
- `TestGetDelivery_Success` - успешное получение доставки
- `TestGetDelivery_NotFound` - получение несуществующей доставки
- `TestGetDeliveryByOrderID_Success` - получение доставки по order_id
- `TestUpdateDeliveryStatus_Success` - успешное обновление статуса
- `TestUpdateDeliveryStatus_NotFound` - обновление статуса несуществующей доставки
- `TestUpdateDeliveryStatus_UpdateError` - ошибка при обновлении

### 2. Unit тесты для Handler (`internal/handler/handler_test.go`)

Тесты gRPC обработчиков с использованием моков для сервиса:

- `TestCreateDelivery_Success` - успешное создание доставки
- `TestCreateDelivery_InvalidOrderID` - валидация order_id
- `TestCreateDelivery_InvalidAddress` - валидация address
- `TestGetDelivery_Success` - успешное получение доставки
- `TestGetDelivery_NotFound` - обработка NotFound
- `TestGetDelivery_InvalidID` - валидация delivery_id
- `TestUpdateDeliveryStatus_Success` - успешное обновление статуса
- `TestUpdateDeliveryStatus_InvalidID` - валидация delivery_id
- `TestUpdateDeliveryStatus_InvalidStatus` - валидация status
- `TestUpdateDeliveryStatus_NotFound` - обработка NotFound

### 3. Integration тесты для Repository (`internal/repository/repository_test.go`)

Тесты работы с реальной PostgreSQL базой данных:

- `TestCreateDelivery_Success` - успешное создание доставки
- `TestCreateDelivery_DuplicateOrderID` - проверка уникальности order_id
- `TestGetDelivery_Success` - успешное получение доставки
- `TestGetDelivery_NotFound` - получение несуществующей доставки
- `TestGetDeliveryByOrderID_Success` - получение доставки по order_id
- `TestGetDeliveryByOrderID_NotFound` - получение несуществующей доставки по order_id
- `TestUpdateDeliveryStatus_Success` - успешное обновление статуса
- `TestUpdateDeliveryStatus_NotFound` - обновление несуществующей доставки

## Запуск тестов

### Запуск всех тестов

```bash
cd delivery-service
go test ./... -v
```

### Запуск тестов конкретного пакета

```bash
# Тесты сервиса
go test ./internal/service -v

# Тесты обработчика
go test ./internal/handler -v

# Тесты репозитория (требует работающую БД)
go test ./internal/repository -v
```

### Запуск с покрытием

```bash
go test ./... -cover
```

### Интеграционные тесты

**Важно:** Интеграционные тесты требуют запущенный PostgreSQL контейнер:

```bash
# Убедитесь, что контейнер запущен
docker ps | grep postgres-delivery

# Если нет, запустите его
docker-compose up -d postgres-delivery

# Подождите несколько секунд для инициализации БД
sleep 3

# Запустите тесты
go test ./internal/repository -v
```

## Требования для тестов

1. **Go 1.25+**
2. **PostgreSQL** (для интеграционных тестов):
   - Host: `localhost`
   - Port: `5436`
   - Database: `deliveries_db`
   - User: `user`
   - Password: `password`

3. **Зависимости:**
   - `github.com/stretchr/testify` - для assertions и моков
   - `github.com/lib/pq` - драйвер PostgreSQL

## Архитектура тестов

Тесты следуют тем же принципам, что и в других сервисах:

- **Моки** используют `stretchr/testify/mock`
- **Интеграционные тесты** используют реальную БД с очисткой перед каждым тестом
- **Проверка ошибок** включает проверку gRPC статус кодов

## Примеры использования моков

### Service тесты

```go
mockRepo := new(MockRepository)
service := New(mockRepo)

mockRepo.On("CreateDelivery", ctx, mock.AnythingOfType("*model.Delivery")).Return(nil)
```

### Handler тесты

```go
mockService := new(MockDeliveryService)
handler := New(mockService)

mockService.On("GetDelivery", ctx, int64(1)).Return(expectedDelivery, nil)
```

## Примечания

- Все интеграционные тесты используют `setupTestDBWithCleanup`, который автоматически очищает таблицу перед каждым тестом
- Тесты проверяют правильность обработки gRPC статус кодов (`codes.InvalidArgument`, `codes.NotFound`, `codes.Internal`)
- Валидация входных данных проверяется на уровне handler

