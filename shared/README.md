# Shared Module

Общий модуль с переиспользуемыми компонентами.

## Описание

Модуль содержит общие компоненты, используемые всеми микросервисами:
- Protocol Buffers определения
- Общие утилиты (logger, errors)
- Переиспользуемые типы

## Структура

```
shared/
├── pb/              # Protocol Buffers определения
│   ├── users.proto
│   ├── goods.proto
│   ├── orders.proto
│   ├── payments.proto
│   └── delivery.proto
└── pkg/             # Переиспользуемые пакеты
    ├── logger/      # Структурированное логирование
    └── errors/      # Общие ошибки
```

## Protocol Buffers

### Файлы

- **users.proto** - определения для Users Service
- **goods.proto** - определения для Goods Service
- **orders.proto** - определения для Orders Service
- **payments.proto** - определения для Payments Service
- **delivery.proto** - определения для Delivery Service

### Генерация кода

После изменения `.proto` файлов, нужно сгенерировать Go код:

```bash
cd shared/pb
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       *.proto
```

### Зависимости

Для генерации нужны:
- `protoc` - компилятор Protocol Buffers
- `protoc-gen-go` - плагин для генерации Go кода
- `protoc-gen-go-grpc` - плагин для генерации gRPC кода

Установка:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Пакеты

### logger

Структурированное логирование с использованием `zap`.

**Использование:**
```go
import "github.com/che1nov/tea-shop/shared/pkg/logger"

func main() {
    logger.Init()
    logger.Info("Service started", "port", 8080)
    logger.Error("Error occurred", "error", err)
}
```

### errors

Общие определения ошибок (если нужны).

## Использование в сервисах

Все сервисы импортируют shared модуль:

```go
import (
    pb "github.com/che1nov/tea-shop/shared/pb"
    "github.com/che1nov/tea-shop/shared/pkg/logger"
)
```

## Версионирование

Модуль использует `replace` директиву в go.mod каждого сервиса:

```go
replace github.com/che1nov/tea-shop/shared => ../shared
```

Это позволяет работать с локальной версией модуля без необходимости публикации в репозиторий.

## Важно

- При изменении proto файлов нужно пересобрать код во всех сервисах
- Изменения в shared модуле влияют на все сервисы
- Рекомендуется использовать semantic versioning при публикации

