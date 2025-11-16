# Тестовые запросы для проверки данных в Grafana

Проверьте эти запросы в Prometheus UI (http://localhost:9090) перед использованием в Grafana:

## 1. Проверка доступности сервисов

```promql
up{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}
```

**Ожидаемый результат**: 7 метрик со значением `1`

## 2. Количество горутин

```promql
go_goroutines{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}
```

**Ожидаемый результат**: Метрики для каждого сервиса (обычно 10-50 горутин)

## 3. Использование памяти

```promql
go_memstats_alloc_bytes{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"} / 1024 / 1024
```

**Ожидаемый результат**: Метрики в мегабайтах для каждого сервиса

## 4. Использование CPU

```promql
rate(process_cpu_seconds_total{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}[5m]) * 100
```

**Ожидаемый результат**: Процент использования CPU для каждого сервиса

## 5. Все метрики для одного сервиса

```promql
{job="users-service"}
```

**Ожидаемый результат**: Все метрики для users-service

## Проверка в Grafana

1. Откройте Grafana: http://localhost:3000
2. Перейдите в **Explore** (иконка компаса)
3. Выберите datasource **Prometheus**
4. Вставьте любой из запросов выше
5. Нажмите **Run query**
6. Должны появиться данные

## Если данных нет

1. Проверьте, что сервисы запущены:
   ```bash
   ./check_monitoring.sh
   ```

2. Проверьте, что Prometheus собирает метрики:
   ```bash
   curl http://localhost:9090/api/v1/targets | jq '.data.activeTargets[] | {job: .labels.job, health: .health}'
   ```

3. Проверьте доступность метрик напрямую:
   ```bash
   curl http://localhost:9001/metrics | head -20
   ```

4. Проверьте, что datasource правильно настроен:
   - Configuration → Data Sources → Prometheus
   - URL должен быть: `http://prometheus:9090`
   - Нажмите "Save & Test"







