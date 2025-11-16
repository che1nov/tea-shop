# Users Service

Сервис управления пользователями и аутентификации.

## Описание

Сервис отвечает за:
- Регистрацию новых пользователей
- Аутентификацию пользователей
- Генерацию и валидацию JWT токенов
- Управление данными пользователей

## Технологии

- **Go 1.25+**
- **gRPC** - для межсервисного взаимодействия
- **PostgreSQL** - база данных (users_db)
- **JWT** - для токенов аутентификации
- **bcrypt** - для хеширования паролей
- **Prometheus** - метрики

## Архитектура

Сервис следует принципам Clean Architecture:

```
Handler (gRPC) → Service → Repository → Database
```

### Слои

- **Handler** (`internal/handler/`) - обработка gRPC запросов
- **Service** (`internal/service/`) - бизнес-логика (хеширование паролей, генерация JWT)
- **Repository** (`internal/repository/`) - работа с БД
- **Model** (`internal/model/`) - доменные модели

## API Документация

API описан в proto файле: `shared/pb/users.proto`

Для тестирования gRPC API используйте:
- **grpcui** (веб-интерфейс): `grpcui -plaintext localhost:8001`
- `grpcurl` (CLI, см. `README_API_TESTING.md`)
- BloomRPC (GUI)
- Postman (gRPC support)

## gRPC API

### Методы

#### CreateUser
Создает нового пользователя в системе.

**Request:**
```protobuf
message CreateUserRequest {
  string email = 1;
  string name = 2;
  string password = 3;
}
```

**Response:**
```protobuf
message User {
  int64 id = 1;
  string email = 2;
  string name = 3;
  string password_hash = 4;
  int64 created_at = 5;
}
```

#### GetUser
Получает информацию о пользователе по ID.

**Request:**
```protobuf
message GetUserRequest {
  int64 user_id = 1;
}
```

#### Login
Аутентифицирует пользователя и возвращает JWT токен.

**Request:**
```protobuf
message LoginRequest {
  string email = 1;
  string password = 2;
}
```

**Response:**
```protobuf
message LoginResponse {
  string token = 1;
  User user = 2;
}
```

#### ValidateToken
Валидирует JWT токен и возвращает информацию о пользователе.

**Request:**
```protobuf
message ValidateTokenRequest {
  string token = 1;
}
```

**Response:**
```protobuf
message ValidateTokenResponse {
  bool valid = 1;
  int64 user_id = 2;
  string email = 3;
}
```

## Структура базы данных

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_users_email ON users(email);
```

## Конфигурация

Переменные окружения:

- `DB_HOST` - хост БД (по умолчанию: localhost)
- `DB_PORT` - порт БД (по умолчанию: 5432)
- `DB_USER` - пользователь БД (по умолчанию: user)
- `DB_PASSWORD` - пароль БД (по умолчанию: password)
- `DB_NAME` - имя БД (по умолчанию: users_db)
- `JWT_SECRET` - секретный ключ для JWT (по умолчанию: your-secret-key-change-in-production)

## Запуск

```bash
go run ./cmd/main.go
```

Сервис будет доступен на порту **8001** (gRPC).

Метрики Prometheus доступны на порту **9001**.

## Особенности реализации

1. **Хеширование паролей**: Используется bcrypt с дефолтной стоимостью
2. **JWT токены**: Срок действия 24 часа
3. **Уникальность email**: Проверяется на уровне БД и приложения
4. **Health checks**: Реализованы gRPC health checks
5. **Транзакции**: Все операции с БД безопасны

## Тестирование

```bash
go test ./...
```

## Зависимости

- **goods-service** - нет
- **order-service** - нет (опционально для валидации пользователей)
- **payment-service** - нет
- **delivery-service** - нет

## Логирование

Используется структурированное логирование через `shared/pkg/logger`.

## Мониторинг

Метрики Prometheus доступны по адресу: `http://localhost:9001/metrics`

