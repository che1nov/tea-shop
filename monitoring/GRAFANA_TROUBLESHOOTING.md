# Решение проблем с данными в Grafana

## Проблема: Дашборд не показывает данные

### Шаг 1: Проверка данных в Prometheus

```bash
# Проверка метрик
curl 'http://localhost:9090/api/v1/query?query=up{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}'
```

Должны увидеть 7 метрик со значением `1`.

### Шаг 2: Проверка через Grafana Explore

1. Откройте http://localhost:3000
2. Нажмите **Explore** (иконка компаса)
3. Выберите datasource: **Prometheus**
4. Введите запрос: `up{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}`
5. Нажмите **Run query**

**Если данные есть в Explore** → проблема в дашборде
**Если данных нет в Explore** → проблема в datasource или Prometheus

### Шаг 3: Проверка datasource

1. Configuration → Data Sources → Prometheus
2. Нажмите **Save & Test**
3. Должно быть: "Data source is working"
4. Если ошибка: проверьте URL (`http://prometheus:9090`)

### Шаг 4: Проверка дашборда

1. Откройте дашборд
2. Нажмите на любую панель → Edit
3. Проверьте:
   - **Data source**: должен быть "Prometheus"
   - **Query**: должен быть правильный PromQL запрос
   - **Time range**: должен быть "Last 15 minutes" или больше
4. Нажмите **Run query** в панели

### Шаг 5: Пересоздание дашборда

Если ничего не помогает:

1. Удалите текущий дашборд
2. Импортируйте упрощенный дашборд:
   - Dashboards → Import
   - Upload JSON file → `monitoring/dashboard-simple.json`
   - Выберите Prometheus datasource
   - Import

## Проблема: "No data" в панелях

### Возможные причины:

1. **Неправильный временной диапазон**
   - Проверьте time picker (правый верхний угол)
   - Выберите "Last 15 minutes" или "Last 1 hour"

2. **Неправильный запрос**
   - Проверьте PromQL запрос в панели
   - Скопируйте запрос в Explore и проверьте там

3. **Нет данных в Prometheus**
   - Проверьте: `curl http://localhost:9090/api/v1/targets`
   - Все targets должны быть "up"

4. **Сервисы не запущены**
   - Проверьте: `./check_monitoring.sh`
   - Все сервисы должны быть OK

## Проблема: "Query failed" или ошибка подключения

1. Проверьте подключение Grafana к Prometheus:
   ```bash
   docker exec grafana wget -qO- http://prometheus:9090/api/v1/status/config
   ```

2. Проверьте, что контейнеры в одной сети:
   ```bash
   docker inspect prometheus grafana | grep -A 5 "Networks"
   ```

3. Перезапустите контейнеры:
   ```bash
   docker-compose restart prometheus grafana
   ```

## Быстрая проверка

```bash
# Проверка всех компонентов
./check_monitoring.sh
./check_grafana.sh

# Проверка данных в Prometheus
curl 'http://localhost:9090/api/v1/query?query=up' | jq '.data.result | length'
# Должно быть 7

# Проверка метрик горутин
curl 'http://localhost:9090/api/v1/query?query=go_goroutines' | jq '.data.result | length'
# Должно быть 7
```

## Полезные запросы для тестирования

### Простой запрос (всегда работает):
```promql
up
```

### По конкретному сервису:
```promql
{job="users-service"}
```

### Горутины:
```promql
go_goroutines{job="users-service"}
```

### Память:
```promql
go_memstats_alloc_bytes{job="users-service"} / 1024 / 1024
```

## Контакты для помощи

Если проблема не решена:
1. Проверьте логи Grafana: `docker logs grafana`
2. Проверьте логи Prometheus: `docker logs prometheus`
3. Проверьте статус сервисов: `./check_monitoring.sh`







