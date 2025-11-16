# E-commerce Tea Shop Platform

Микросервисная платформа для интернет-магазина чая, построенная на Go и React.

## Архитектура

Проект реализован с использованием принципов **Clean Architecture** и **микросервисной архитектуры**:

- **Backend**: Go (gRPC, Gin)
- **Frontend**: React + TypeScript + Vite
- **База данных**: PostgreSQL (отдельная БД для каждого сервиса)
- **Message Queue**: Kafka
- **Мониторинг**: Prometheus + Grafana
- **Трейсинг**: Jaeger

## Структура проекта

```
ecommerce/
├── api-gateway/          # API Gateway (HTTP -> gRPC)
│   └── README.md         # Документация API Gateway
├── users-service/         # Сервис пользователей
│   └── README.md         # Документация Users Service
├── goods-service/         # Сервис товаров (каталог)
│   └── README.md         # Документация Goods Service
├── order-service/         # Сервис заказов
│   └── README.md         # Документация Order Service
├── payment-service/       # Сервис платежей
│   └── README.md         # Документация Payment Service
├── delivery-service/      # Сервис доставки
│   └── README.md         # Документация Delivery Service
├── notify-service/        # Сервис уведомлений (Kafka consumer)
│   └── README.md         # Документация Notify Service
├── shared/                # Общие типы и утилиты
│   ├── pb/               # Protocol Buffers
│   ├── pkg/              # Переиспользуемые пакеты
│   └── README.md         # Документация Shared Module
└── frontend/              # React фронтенд
    └── README.md          # Документация Frontend
```

Подробная документация по каждому сервису доступна в соответствующих README.md файлах.

## API Документация

### API Gateway (HTTP REST)
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- Интерактивная документация с возможностью тестирования API

### gRPC Сервисы
- **Proto файлы**: `shared/pb/*.proto` - описание API
- **Тестирование**: Используйте **grpcui** (веб-интерфейс) или другие инструменты

**Быстрый старт с grpcui:**
```bash
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
grpcui -plaintext localhost:8001  # Откроет веб-интерфейс на http://127.0.0.1:8080
```

Подробнее см.:
- [SWAGGER_SERVICES.md](./SWAGGER_SERVICES.md) - Swagger для API Gateway
- [GRPC_DOCUMENTATION.md](./GRPC_DOCUMENTATION.md) - работа с gRPC API
- [REAL_WORLD_PRACTICES.md](./REAL_WORLD_PRACTICES.md) - практики в индустрии

## Быстрый старт

### Требования

- Go 1.25+
- Docker & Docker Compose
- Node.js 18+ (для фронтенда)

### 1. Запуск инфраструктуры

```bash
docker-compose up -d
```

Это запустит:
- PostgreSQL (5 баз данных)
- Kafka + Zookeeper
- Prometheus
- Grafana
- Jaeger

### 2. Запуск микросервисов

```bash
./start_all_services.sh
```

Сервисы будут доступны на:
- API Gateway: `http://localhost:8080`
- Users Service: `localhost:8001`
- Goods Service: `localhost:8002`
- Order Service: `localhost:8003`
- Payment Service: `localhost:8004`
- Delivery Service: `localhost:8005`

### 3. Запуск фронтенда

```bash
cd frontend
npm install
npm run dev
```

Фронтенд будет доступен на `http://localhost:5173`

### 4. Остановка всех сервисов

```bash
./stop_all_services.sh
docker-compose down
```

## Функционал

### Для пользователей:
- Регистрация и авторизация (JWT)
- Просмотр каталога товаров
- Добавление товаров в корзину
- Оформление заказов
- Просмотр истории заказов

### Для администраторов:
- Админ-панель для управления товарами
- Создание, редактирование, удаление товаров
- Управление остатками

## Архитектурные принципы

### Clean Architecture

Проект следует принципам Clean Architecture:

```
Controllers → Use Cases ← Interfaces ← Adapters
     ↓           ↓
    DTO        Domain
```

- **Use Cases** - бизнес-логика (помещается на один экран)
- **Controllers** - HTTP/gRPC handlers (только декодирование/кодирование)
- **Adapters** - реализация интерфейсов (PostgreSQL, Kafka)
- **Domain** - доменные сущности с валидацией

### Dependency Injection

Все зависимости инжектируются через конструкторы:

```go
type ProfileUseCase struct {
    postgres PostgreSQLAdapter
    redis    RedisAdapter
    kafka    KafkaProducer
}

func NewProfileUseCase(
    postgres PostgreSQLAdapter,
    redis RedisAdapter,
    kafka KafkaProducer,
) *ProfileUseCase {
    return &ProfileUseCase{
        postgres: postgres,
        redis:    redis,
        kafka:    kafka,
    }
}
```

## API Endpoints

Полная документация API доступна в Swagger UI: `http://localhost:8080/swagger/index.html`

### Публичные:
- `POST /api/v1/auth/register` - Регистрация
- `POST /api/v1/auth/login` - Вход
- `GET /api/v1/goods` - Список товаров
- `GET /api/v1/goods/:id` - Детали товара

### Защищенные (требуют JWT):
- `GET /api/v1/users/me` - Информация о пользователе
- `POST /api/v1/orders` - Создание заказа
- `GET /api/v1/orders/:id` - Детали заказа
- `GET /api/v1/payments/:id` - Информация о платеже
- `POST /api/v1/deliveries` - Создание доставки
- `GET /api/v1/deliveries/:id` - Информация о доставке

### Админ (требуют JWT + роль "admin"):
- `POST /api/v1/admin/goods` - Создание товара
- `PUT /api/v1/admin/goods/:id` - Обновление товара
- `DELETE /api/v1/admin/goods/:id` - Удаление товара

**Важно**: Админ-эндпоинты требуют роль `"admin"` в JWT токене. Обычные пользователи получат ошибку 403 Forbidden.

**Админские credentials**: Настраиваются через переменные окружения `ADMIN_EMAIL` и `ADMIN_PASSWORD`. По умолчанию: `admin@example.com` / `admin123`. **Обязательно измените в production!** Подробнее см. [ADMIN_ROLES.md](./ADMIN_ROLES.md)

### Swagger UI

После запуска API Gateway, интерактивная документация доступна по адресу:
```
http://localhost:8080/swagger/index.html
```

В Swagger UI можно:
- Просматривать все доступные endpoints
- Видеть схемы запросов и ответов
- Тестировать API прямо из браузера
- Авторизоваться с JWT токеном

## Тестирование

Каждый сервис содержит unit-тесты:

```bash
cd users-service
go test ./...
```

## Мониторинг

- **Prometheus**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000` (admin/admin)
- **Jaeger**: `http://localhost:16686`

## Технологии

### Backend:
- Go 1.25
- gRPC
- Gin (HTTP router)
- PostgreSQL
- Kafka
- JWT для аутентификации

### Frontend:
- React 18
- TypeScript
- Vite
- Tailwind CSS
- React Query
- Zustand
- React Router

## Особенности реализации

1. **Транзакции** - все операции с БД выполняются в транзакциях
2. **Резервирование товаров** - при создании заказа товары резервируются
3. **Обработка ошибок** - доменные ошибки с понятными сообщениями
4. **Логирование** - структурированное логирование через zap
5. **Метрики** - Prometheus метрики для всех сервисов
6. **Health checks** - gRPC health checks для всех сервисов

## Безопасность

- JWT токены для аутентификации
- Валидация входных данных
- Защита от SQL инъекций (prepared statements)
- CORS настройки для фронтенда

## Разработка

### Добавление нового сервиса

1. Создайте структуру сервиса по аналогии с существующими
2. Добавьте proto определения в `shared/pb/`
3. Пересоберите proto файлы
4. Добавьте сервис в `start_all_services.sh`
5. Обновите API Gateway для проксирования запросов
6. Добавьте Swagger аннотации к новым endpoints

### Пересборка Proto файлов

```bash
cd shared/pb
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       *.proto
```

### Обновление Swagger документации

После изменения API endpoints:

```bash
cd api-gateway
swag init -g cmd/main.go
```

Swagger UI автоматически обновится после перезапуска API Gateway.

## Лицензия

MIT

## Автор

Проект разработан как демонстрация навыков разработки микросервисных приложений на Go.

