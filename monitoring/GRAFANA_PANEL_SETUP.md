# Настройка панели в Grafana - Пошаговая инструкция

## Режим редактирования панели

### 1. Query Editor (Query A - Prometheus)

#### Metric / Query:
В поле запроса введите один из вариантов:

**Для CPU:**
```
rate(process_cpu_seconds_total{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}[5m]) * 100
```

**Для памяти:**
```
go_memstats_heap_inuse_bytes{job=~"users-service|goods-service|order-service|order-service|payment-service|delivery-service|notify-service|api-gateway"} / 1024 / 1024
```

**Для горутин:**
```
go_goroutines{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}
```

**Для статуса сервисов:**
```
up{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}
```

#### Label filters:
- **Select label:** `job`
- **Select value:** Можно выбрать конкретный сервис или оставить пустым (тогда будут все)

#### Operations:
- Для суммирования: `+ Operations` → `Aggregations` → `Sum`
- Для среднего: `+ Operations` → `Aggregations` → `Mean`

### 2. Options (под запросом)

- **Legend:** `{{job}}` - показывает имя сервиса
- **Format:** `Time series` (для графиков) или `Table` (для таблиц)
- **Step:** `auto` (автоматически) или `15s` (каждые 15 секунд)
- **Type:** `Range` (для временных рядов)
- **Exemplars:** `false` (обычно не нужны)

### 3. Visualization (правая панель)

- **Type:** 
  - `Time series` - для графиков временных рядов
  - `Stat` - для одной метрики
  - `Gauge` - для gauge индикаторов
  - `Table` - для таблиц

### 4. Tooltip (в Visualization)

- **Tooltip mode:** 
  - `Single` - показывает одну точку при наведении
  - `All` - показывает все серии
  - `Hidden` - скрывает tooltip

### 5. Legend (в Visualization)

- **Visibility:** `On` (включить легенду)
- **Mode:** 
  - `List` - список внизу
  - `Table` - таблица с значениями
- **Placement:** 
  - `Bottom` - снизу
  - `Right` - справа
- **Values:** Можно выбрать `Last`, `Max`, `Min` для отображения в легенде

### 6. Axis (в Visualization)

- **Time zone:** `Default` или конкретная зона
- **Placement:** Настройка осей (обычно `auto`)

### 7. Сохранение

- **Save dashboard** (синяя кнопка справа вверху)
- **Discard panel changes** (красная кнопка) - отменить изменения

## Примеры готовых запросов

### CPU Usage по всем сервисам:
```
rate(process_cpu_seconds_total{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}[5m]) * 100
```
Legend: `{{job}}`
Visualization: `Time series`

### Memory Usage (MB):
```
go_memstats_heap_inuse_bytes{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"} / 1024 / 1024
```
Legend: `{{job}}`
Visualization: `Time series`
Unit: `mbytes`

### Go Routines:
```
go_goroutines{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}
```
Legend: `{{job}}`
Visualization: `Time series`

### Service Status:
```
up{job=~"users-service|goods-service|order-service|payment-service|delivery-service|notify-service|api-gateway"}
```
Legend: `{{job}}`
Visualization: `Stat` или `Table`
Format: `Table`

## Быстрые советы

1. **Если нет данных:**
   - Проверьте временной диапазон (Last 5 minutes может быть мало)
   - Выберите "Last 15 minutes" или "Last 1 hour"
   - Нажмите "Run queries"

2. **Если слишком много данных:**
   - Добавьте фильтр по label: `job="users-service"`
   - Или используйте операции для агрегации

3. **Для лучшей визуализации:**
   - Используйте `Time series` для трендов
   - Используйте `Stat` для быстрых метрик
   - Используйте `Gauge` для процентов (CPU, Memory)







