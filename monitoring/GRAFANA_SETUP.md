# Настройка Grafana для ecommerce системы

## Быстрый старт

Grafana уже настроена и доступна по адресу: http://localhost:3000

### Учетные данные по умолчанию:
- **Username**: `admin`
- **Password**: `admin`

## Автоматическая настройка

После перезапуска Grafana, Prometheus datasource будет настроен автоматически.

### Перезапуск Grafana:

```bash
docker-compose restart grafana
```

## Ручная настройка Prometheus Datasource

Если автоматическая настройка не сработала:

1. Откройте http://localhost:3000
2. Войдите с учетными данными `admin/admin`
3. Перейдите в **Configuration** → **Data Sources**
4. Нажмите **Add data source**
5. Выберите **Prometheus**
6. В поле **URL** введите: `http://prometheus:9090`
7. Нажмите **Save & Test**

## Создание дашборда

### Вариант 1: Импорт готового дашборда

1. Откройте Grafana: http://localhost:3000
2. Перейдите в **Dashboards** → **Import**
3. Скопируйте содержимое файла `monitoring/dashboard-services.json`
4. Вставьте JSON в поле **Import via panel json**
5. Нажмите **Load**
6. Выберите **Prometheus** datasource
7. Нажмите **Import**

### Вариант 2: Создание дашборда вручную

#### 1. Создание нового дашборда

1. Нажмите **+** → **Create Dashboard**
2. Нажмите **Add visualization**

#### 2. Добавление панелей

##### Панель "Service Health"

- **Query**: `up{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}`
- **Visualization**: **Stat**
- **Format**: **Table**

##### Панель "Go Routines"

- **Query**: `go_goroutines{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}`
- **Visualization**: **Time series**
- **Legend**: `{{job}}`

##### Панель "Memory Usage"

- **Query**: `go_memstats_alloc_bytes{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"} / 1024 / 1024`
- **Visualization**: **Time series**
- **Unit**: **Mbytes**
- **Legend**: `{{job}}`

##### Панель "CPU Usage"

- **Query**: `rate(process_cpu_seconds_total{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}[5m]) * 100`
- **Visualization**: **Time series**
- **Unit**: **Percent (0-100)**
- **Legend**: `{{job}}`

##### Панель "HTTP Requests Rate"

- **Query**: `rate(promhttp_metric_handler_requests_total{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}[5m])`
- **Visualization**: **Time series**
- **Unit**: **reqps (requests/sec)**
- **Legend**: `{{job}} - {{code}}`

## Полезные запросы Prometheus

### Проверка доступности сервисов

```promql
up{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}
```

### Количество горутин по сервисам

```promql
go_goroutines{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}
```

### Использование памяти

```promql
go_memstats_alloc_bytes{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"} / 1024 / 1024
```

### Использование CPU

```promql
rate(process_cpu_seconds_total{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}[5m]) * 100
```

### Количество HTTP запросов к метрикам

```promql
rate(promhttp_metric_handler_requests_total{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}[5m])
```

## Проверка работы

1. Убедитесь, что Prometheus собирает метрики:
   ```bash
   curl http://localhost:9090/api/v1/targets | jq '.data.activeTargets[].health'
   ```

2. Убедитесь, что Grafana доступна:
   ```bash
   curl http://localhost:3000/api/health
   ```

3. Проверьте, что datasource настроен:
   - Откройте http://localhost:3000
   - Configuration → Data Sources
   - Должен быть Prometheus datasource

## Проверка данных через Explore

Если дашборд не показывает данные, проверьте через Explore:

1. Откройте Grafana: http://localhost:3000
2. Нажмите на иконку **Explore** (компас) в левом меню
3. Выберите datasource: **Prometheus**
4. Введите запрос: `up{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}`
5. Нажмите **Run query**
6. Должны увидеть 7 метрик со значением `1`

Если данные есть в Explore, но нет в дашборде:
- Проверьте, что в дашборде выбран правильный datasource (Prometheus)
- Проверьте временной диапазон (Last 15 minutes или Last 1 hour)
- Попробуйте обновить дашборд (F5 или кнопка Refresh)

## Troubleshooting

### Grafana не видит Prometheus

1. Проверьте, что оба контейнера запущены:
   ```bash
   docker ps | grep -E "(prometheus|grafana)"
   ```

2. Проверьте логи Grafana:
   ```bash
   docker logs grafana
   ```

3. Проверьте, что Prometheus доступен из контейнера Grafana:
   ```bash
   docker exec grafana wget -qO- http://prometheus:9090/api/v1/status/config
   ```

### Дашборд не показывает данные

1. Проверьте, что Prometheus datasource правильно настроен
2. Проверьте, что запросы PromQL корректны (используйте Prometheus UI для тестирования)
3. Убедитесь, что метрики собираются:
   ```bash
   curl http://localhost:9090/api/v1/query?query=up
   ```

## Дополнительные ресурсы

- [Grafana Documentation](https://grafana.com/docs/)
- [Prometheus Query Language](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Grafana Dashboard Examples](https://grafana.com/grafana/dashboards/)

