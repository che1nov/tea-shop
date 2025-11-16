# Тестирование Users Service

## Структура тестов

### Unit-тесты

#### 1. Service тесты (`internal/service/service_test.go`)
- ✅ `TestNew` - создание сервиса
- ✅ `TestCreateUser_Success` - успешное создание пользователя
- ✅ `TestCreateUser_RepositoryError` - ошибка при создании
- ✅ `TestGetUser_Success` - получение пользователя
- ✅ `TestGenerateToken_Success` - генерация JWT токена
- ✅ `TestValidateToken_ValidToken` - валидация валидного токена
- ✅ `TestValidateToken_InvalidToken` - валидация невалидного токена
- ✅ `TestValidateToken_WrongSecret` - токен с неправильным секретом
- ✅ `TestLogin_Success` - успешный логин
- ✅ `TestLogin_UserNotFound` - логин несуществующего пользователя

#### 2. Handler тесты (`internal/handler/handler_test.go`)
- ✅ `TestNew` - создание handler
- ✅ `TestCreateUser_Success` - успешное создание через API
- ✅ `TestCreateUser_MissingEmail` - валидация: отсутствует email
- ✅ `TestCreateUser_MissingName` - валидация: отсутствует name
- ✅ `TestCreateUser_MissingPassword` - валидация: отсутствует password
- ✅ `TestCreateUser_EmailAlreadyExists` - ошибка дубликата email
- ✅ `TestGetUser_Success` - успешное получение пользователя
- ✅ `TestGetUser_InvalidUserId` - валидация: невалидный user_id
- ✅ `TestGetUser_NotFound` - пользователь не найден
- ✅ `TestValidateToken_ValidToken` - валидация валидного токена
- ✅ `TestValidateToken_InvalidToken` - валидация невалидного токена
- ✅ `TestValidateToken_EmptyToken` - пустой токен

#### 3. Repository тесты (`internal/repository/repository_test.go`)
- ✅ `TestCreateUser_Success` - создание пользователя
- ✅ `TestCreateUser_DuplicateEmail` - проверка уникальности email
- ✅ `TestGetUserByID_Success` - получение по ID
- ✅ `TestGetUserByID_NotFound` - пользователь не найден
- ✅ `TestGetUserByEmail_Success` - получение по email
- ✅ `TestGetUserByEmail_NotFound` - email не найден

## Запуск тестов

### Запуск всех тестов
```bash
cd users-service
go test ./...
```

### Запуск с подробным выводом
```bash
go test -v ./...
```

### Запуск конкретного пакета
```bash
go test -v ./internal/service
go test -v ./internal/handler
go test -v ./internal/repository
```

### Запуск конкретного теста
```bash
go test -v ./internal/service -run TestCreateUser_Success
```

### Запуск с покрытием кода
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Покрытие тестами

Текущее покрытие:
- **Service**: бизнес-логика, JWT токены, валидация
- **Handler**: gRPC endpoints, валидация входных данных, обработка ошибок
- **Repository**: CRUD операции, уникальные ограничения

## Что тестируется

### Service Layer
- ✅ Создание пользователей с хешированием паролей
- ✅ Получение пользователей
- ✅ Генерация и валидация JWT токенов
- ✅ Логин пользователей
- ✅ Обработка ошибок репозитория

### Handler Layer
- ✅ Валидация входных данных (email, name, password, user_id)
- ✅ Преобразование между protobuf и domain моделями
- ✅ Правильные gRPC статус-коды (InvalidArgument, NotFound, AlreadyExists, Internal)
- ✅ Обработка всех типов ошибок

### Repository Layer
- ✅ CRUD операции с БД
- ✅ Уникальное ограничение на email
- ✅ Поиск по ID и email
- ✅ Обработка случая "не найдено"

## Зависимости для тестов

```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
go get github.com/stretchr/testify/require
```

## Интеграционные тесты

Repository тесты используют реальную БД. Убедитесь, что:
1. PostgreSQL запущен (Docker контейнер)
2. БД `users_db` существует
3. Пользователь `user` с паролем `password` имеет доступ

Для изоляции тестов можно использовать:
- Testcontainers (Docker контейнеры для тестов)
- In-memory база данных
- Отдельная тестовая БД

## Следующие шаги

1. Добавить интеграционные тесты end-to-end
2. Добавить тесты производительности (benchmarks)
3. Настроить CI/CD для автоматического запуска тестов
4. Добавить тесты для конфигурации

