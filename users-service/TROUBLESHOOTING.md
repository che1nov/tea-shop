# Решение проблем с зависимостями

## Проблема: "no required module provides package"

Если возникает ошибка:
```
no required module provides package github.com/che1nov/ecommerce/shared/pb
```

### Решение:

1. **Обновить зависимости:**
```bash
cd users-service
go mod tidy
```

2. **Очистить кэш модулей:**
```bash
go clean -modcache
go mod download
```

3. **Проверить replace директиву в go.mod:**
Убедитесь, что в `users-service/go.mod` есть:
```go
replace github.com/che1nov/ecommerce/shared => ../shared
```

4. **Проверить, что shared модуль корректен:**
```bash
cd ../shared
go build ./pb
```

5. **Перезапустить IDE** (VS Code, GoLand и т.д.) для обновления индексации

6. **Проверить версию Go:**
```bash
go version
# Должна быть 1.25.1 или выше
```

### Проверка работоспособности:

```bash
# Проверить импорт
go list github.com/che1nov/ecommerce/shared/pb

# Проверить сборку
go build ./...

# Проверить тесты
go test ./...
```

## Структура модулей

```
ecommerce/
├── shared/
│   ├── go.mod (module: github.com/che1nov/ecommerce/shared)
│   └── pb/
│       └── *.pb.go (package: pb)
└── users-service/
    ├── go.mod (module: github.com/che1nov/ecommerce/users-service)
    └── internal/
        └── handler/
            └── handler.go (import: github.com/che1nov/ecommerce/shared/pb)
```

## Важно

- `shared` должен быть отдельным модулем с собственным `go.mod`
- В `users-service/go.mod` должна быть `replace` директива для локального модуля
- После изменения `shared` модуля нужно выполнить `go mod tidy` в `users-service`

