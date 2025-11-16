# Исправление ошибки Swagger

## Проблема

Если при открытии `http://localhost:8080/swagger/index.html` получаете 404, это может быть из-за ошибки компиляции в `docs/docs.go`.

## Решение

Поля `LeftDelim` и `RightDelim` были удалены в новых версиях `swaggo/swag`.

### Автоматическое исправление

```bash
cd api-gateway
swag init -g cmd/main.go
# Если ошибка осталась, удалите строки LeftDelim и RightDelim вручную
```

### Ручное исправление

Отредактируйте `api-gateway/docs/docs.go` и удалите строки:
```go
LeftDelim:        "{{",
RightDelim:       "}}",
```

### Перезапуск

После исправления перезапустите API Gateway:

```bash
# Остановите текущий процесс
./stop_all_services.sh

# Запустите заново
./start_all_services.sh
```

## Проверка

После перезапуска откройте: `http://localhost:8080/swagger/index.html`

