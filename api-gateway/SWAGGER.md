# Swagger документация

## Установка swag

Для генерации Swagger документации необходимо установить `swag`:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Генерация документации

После добавления или изменения аннотаций Swagger, выполните:

```bash
cd api-gateway
swag init -g cmd/main.go
```

Это создаст/обновит файлы в папке `docs/`:
- `docs.go` - сгенерированный код
- `swagger.json` - JSON спецификация
- `swagger.yaml` - YAML спецификация

## Доступ к Swagger UI

После запуска API Gateway, Swagger UI будет доступен по адресу:

```
http://localhost:8080/swagger/index.html
```

## Использование

1. Откройте Swagger UI в браузере
2. Для тестирования защищенных endpoints:
   - Нажмите кнопку "Authorize"
   - Введите ваш JWT токен (полученный после `/api/v1/auth/login`)
   - Формат: `Bearer {ваш_токен}` или просто `{ваш_токен}`
3. Тестируйте endpoints прямо из интерфейса

## Обновление документации

При изменении API endpoints:
1. Обновите Swagger аннотации в `internal/handler/handler.go`
2. Выполните `swag init -g cmd/main.go`
3. Перезапустите сервис

