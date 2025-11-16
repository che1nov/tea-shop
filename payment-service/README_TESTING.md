# Тестирование Payment Service

## Структура тестов

### Unit-тесты

#### 1. Service тесты (`internal/service/service_test.go`)
- ✅ `TestNew` - создание сервиса
- ✅ `TestProcessPayment_Success` - обработка платежа
- ✅ `TestProcessPayment_CreateError` - ошибка при создании
- ✅ `TestGetPayment_Success` - получение платежа
- ✅ `TestGetPayment_NotFound` - платеж не найден
- ✅ `TestGetPaymentByOrderID_Success` - получение по order_id

#### 2. Handler тесты (`internal/handler/handler_test.go`)
- ✅ `TestNew` - создание handler
- ✅ `TestProcessPayment_Success` - обработка платежа
- ✅ `TestGetPayment_Success` - получение платежа
- ✅ `TestGetPayment_NotFound` - платеж не найден

#### 3. Repository тесты (`internal/repository/repository_test.go`)
- ✅ `TestCreatePayment_Success` - создание платежа
- ✅ `TestGetPayment_Success` - получение по ID
- ✅ `TestGetPayment_NotFound` - платеж не найден
- ✅ `TestUpdatePaymentStatus_Success` - обновление статуса
- ✅ `TestGetPaymentByOrderID_Success` - получение по order_id

## Запуск тестов

```bash
cd payment-service

# Все тесты
go test ./...

# С покрытием
go test -cover ./...
```

