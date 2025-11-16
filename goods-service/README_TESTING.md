# Тестирование Goods Service

## Структура тестов

### Unit-тесты

#### 1. Service тесты (`internal/service/service_test.go`)
- ✅ `TestNew` - создание сервиса
- ✅ `TestCreateGood_Success` - создание товара
- ✅ `TestGetGood_Success` - получение товара
- ✅ `TestListGoods_Success` - список товаров
- ✅ `TestGetTotalGoods_Success` - общее количество товаров
- ✅ `TestCheckStock_Available` - проверка наличия (достаточно)
- ✅ `TestCheckStock_NotAvailable` - проверка наличия (недостаточно)
- ✅ `TestCheckStock_GoodNotFound` - товар не найден
- ✅ `TestReserveStock_Success` - успешное резервирование
- ✅ `TestReserveStock_InsufficientStock` - недостаточно товара
- ✅ `TestReserveStock_Error` - ошибка при резервировании

#### 2. Handler тесты (`internal/handler/handler_test.go`)
- ✅ `TestNew` - создание handler
- ✅ `TestGetGood_Success` - получение товара
- ✅ `TestGetGood_NotFound` - товар не найден
- ✅ `TestListGoods_Success` - список товаров
- ✅ `TestCheckStock_Available` - проверка наличия
- ✅ `TestCheckStock_NotAvailable` - товара нет
- ✅ `TestReserveStock_Success` - резервирование успешно
- ✅ `TestReserveStock_InsufficientStock` - недостаточно товара
- ✅ `TestReserveStock_Error` - ошибка резервирования

#### 3. Repository тесты (`internal/repository/repository_test.go`)
- ✅ `TestCreateGood_Success` - создание товара
- ✅ `TestGetGood_Success` - получение по ID
- ✅ `TestGetGood_NotFound` - товар не найден
- ✅ `TestListGoods_Success` - список товаров
- ✅ `TestGetTotalGoods_Success` - общее количество

## Запуск тестов

```bash
cd goods-service

# Все тесты
go test ./...

# С покрытием
go test -cover ./...

# Конкретный пакет
go test ./internal/service -v
go test ./internal/handler -v
```

## Покрытие

- **Service**: бизнес-логика, проверка наличия, резервирование
- **Handler**: gRPC endpoints, преобразование данных
- **Repository**: CRUD операции с БД

