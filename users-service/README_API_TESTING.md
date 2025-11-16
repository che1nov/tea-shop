# Тестирование Users Service API

## Способы тестирования gRPC API

### 1. Использование grpcurl (рекомендуется)

`grpcurl` уже установлен. Используйте его для тестирования API.

#### Базовые команды:

```bash
# Список всех сервисов
grpcurl -plaintext localhost:8001 list

# Список методов UsersService
grpcurl -plaintext localhost:8001 list pb.UsersService

# Описание метода
grpcurl -plaintext localhost:8001 describe pb.UsersService.CreateUser
```

#### Тестирование CreateUser:

```bash
grpcurl -plaintext -d '{
  "email": "user@example.com",
  "name": "John Doe",
  "password": "securepass123"
}' localhost:8001 pb.UsersService/CreateUser
```

#### Тестирование GetUser:

```bash
grpcurl -plaintext -d '{
  "user_id": 1
}' localhost:8001 pb.UsersService/GetUser
```

#### Тестирование ValidateToken:

```bash
grpcurl -plaintext -d '{
  "token": "your-jwt-token-here"
}' localhost:8001 pb.UsersService/ValidateToken
```

### 2. Использование готового скрипта

Запустите автоматизированный тест:

```bash
cd users-service
./test_api.sh
```

### 3. Использование BloomRPC (GUI инструмент)

1. Установите BloomRPC: https://github.com/bloomrpc/bloomrpc
2. Импортируйте proto файл: `shared/pb/users.proto`
3. Подключитесь к `localhost:8001`
4. Выберите метод и отправьте запрос

### 4. Health Check

Проверка здоровья сервиса:

```bash
grpcurl -plaintext localhost:8001 list grpc.health.v1.Health

grpcurl -plaintext -d '{
  "service": "pb.UsersService"
}' localhost:8001 grpc.health.v1.Health/Check
```

## Примеры ответов

### Успешное создание пользователя:
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "password_hash": "$2a$10$...",
  "created_at": 1234567890
}
```

### Ошибка валидации:
```
ERROR:
  Code: InvalidArgument
  Message: email is required
```

### Пользователь не найден:
```
ERROR:
  Code: NotFound
  Message: user with id 999 not found
```

### Дубликат email:
```
ERROR:
  Code: AlreadyExists
  Message: user with email test@example.com already exists
```

## Запуск сервиса

```bash
cd users-service
go run ./cmd/main.go
```

Или с переменными окружения:

```bash
DB_HOST=localhost DB_PORT=5432 DB_USER=user DB_PASSWORD=password DB_NAME=users_db \
go run ./cmd/main.go
```

## Порт

По умолчанию сервис слушает на порту **8001**.

Можно изменить через переменную окружения или в файле `config/config.go`.

