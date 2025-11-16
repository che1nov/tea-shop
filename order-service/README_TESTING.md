# Тестирование Order Service

## Структура тестов

### Unit-тесты

#### 1. Service тесты (`internal/service/service_test.go`)
- ✅ `TestNew` - создание сервиса
- ✅ `TestGetOrder_Success` - получение заказа
- ✅ `TestUpdateOrderStatus_Success` - обновление статуса
- ✅ `TestListUserOrders_Success` - список заказов пользователя

**Примечание:** `CreateOrder` требует сложных моков для gRPC клиентов и Kafka, поэтому тестируется через интеграционные тесты.

#### 2. Handler тесты (`internal/handler/handler_test.go`)
- ✅ `TestNew` - создание handler
- ✅ `TestGetOrder_Success` - получение заказа
- ✅ `TestGetOrder_NotFound` - заказ не найден
- ✅ `TestUpdateOrderStatus_Success` - обновление статуса
- ✅ `TestOrderToProto` - преобразование в protobuf

#### 3. Repository тесты (`internal/repository/repository_test.go`)
- ✅ `TestCreateOrder_Success` - создание заказа
- ✅ `TestGetOrder_Success` - получение по ID
- ✅ `TestGetOrder_NotFound` - заказ не найден
- ✅ `TestUpdateOrderStatus_Success` - обновление статуса
- ✅ `TestListUserOrders_Success` - список заказов пользователя

## Запуск тестов

```bash
cd order-service

# Все тесты
go test ./...

# С покрытием
go test -cover ./...
```

## Особенности

- Order service использует gRPC клиенты для goods и payment сервисов
- Использует Kafka для публикации событий
- Требует моки для внешних зависимостей

