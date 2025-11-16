# Настройка мониторинга ecommerce системы

## Архитектура мониторинга

Система использует полный стек observability:
- **Prometheus** - сбор метрик
- **Grafana** - визуализация метрик
- **Jaeger** - распределенный трейсинг
- **Logger** (slog) - структурированное логирование

## Портовая схема

### gRPC сервисы (основные порты)
- `users-service`: 8001
- `goods-service`: 8002
- `order-service`: 8003
- `payment-service`: 8004
- `delivery-service`: 8005
- `notify-service`: 8006 (Kafka consumer)
- `api-gateway`: 8080 (HTTP)

### HTTP метрики Prometheus
- `users-service`: 9001
- `goods-service`: 9002
- `order-service`: 9003
- `payment-service`: 9004
- `delivery-service`: 9005
- `notify-service`: 9006
- `api-gateway`: 9007

### Инфраструктура мониторинга
- **Prometheus**: 9090
- **Grafana**: 3000 (admin/admin)
- **Jaeger UI**: 16686

## Запуск мониторинга

### Важно: Сначала запустите все сервисы приложения!

```bash
# 1. Запустить все сервисы приложения
./start_all_services.sh

# 2. Проверить, что сервисы запущены
curl http://localhost:9001/metrics  # users-service
curl http://localhost:9002/metrics  # goods-service
curl http://localhost:9003/metrics  # order-service
curl http://localhost:9004/metrics  # payment-service
curl http://localhost:9005/metrics  # delivery-service
curl http://localhost:9006/metrics  # notify-service
curl http://localhost:9007/metrics  # api-gateway

# 3. Запустить сервисы мониторинга
docker-compose up -d prometheus grafana jaeger

# 4. Проверить статус
docker ps | grep -E "(prometheus|grafana|jaeger)"
```

### Решение проблемы "connection refused"

Если Prometheus показывает все targets как "down":
1. Убедитесь, что все сервисы запущены (см. выше)
2. Prometheus использует `host.docker.internal` для доступа к хосту из контейнера
3. На macOS это работает автоматически
4. На Linux может потребоваться добавить `--add-host=host.docker.internal:host-gateway` при запуске Docker

## Настройка Grafana

Grafana настроена автоматически! Подробная инструкция в файле `monitoring/GRAFANA_SETUP.md`.

### Быстрый старт:

1. Откройте http://localhost:3000
2. Войдите с `admin/admin`
3. Prometheus datasource уже настроен автоматически
4. Импортируйте дашборд из `monitoring/dashboard-services.json` (см. `GRAFANA_SETUP.md`)

## Проверка метрик в Prometheus

1. Откройте http://localhost:9090
2. Проверьте статус targets: Status → Targets
3. Все сервисы должны быть в статусе "UP"
4. Примеры запросов:
   - `up` - статус всех сервисов
   - `go_goroutines` - количество горутин
   - `go_memstats_alloc_bytes` - использование памяти

## Метрики по умолчанию

Все сервисы экспортируют стандартные метрики Go:
- `go_goroutines` - количество горутин
- `go_memstats_*` - статистика памяти
- `process_*` - метрики процесса

## Кастомные метрики

Для добавления кастомных метрик (request count, latency, errors):
1. Используйте `github.com/prometheus/client_golang/prometheus`
2. Регистрируйте метрики в handlers
3. Метрики автоматически будут доступны на `/metrics`

## Пример проверки метрик

```bash
# Проверить метрики users-service
curl http://localhost:9001/metrics | grep go_goroutines

# Проверить все targets в Prometheus
curl http://localhost:9090/api/v1/targets | jq '.data.activeTargets[].health'
```

