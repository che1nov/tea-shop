# Contributing

## Development Setup

1. Установите зависимости:
```bash
go mod download
cd frontend && npm install
```

2. Запустите инфраструктуру:
```bash
docker-compose up -d
```

3. Запустите сервисы:
```bash
./start_all_services.sh
```

## Code Style

- Следуем принципам Clean Architecture
- Используем dependency injection
- Все Use Cases должны помещаться на один экран
- Комментарии на русском языке
- Тесты для всех публичных методов

## Commit Messages

Используйте понятные сообщения коммитов:
- `feat: добавлена функция X`
- `fix: исправлена ошибка Y`
- `refactor: рефакторинг модуля Z`

