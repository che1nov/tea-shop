# Нагрузочное тестирование (RPS + Kafka)

## Что есть в папке

- `api_gateway_smoke.js` — smoke-проверка доступности `api-gateway`.
- `api_gateway_goods_read.js` — read-heavy нагрузка на `GET /api/v1/goods`.
- `scripts/run_rps_and_kafka.sh` — запускает `k6` и считает:
  - `rps_avg`
  - прирост сообщений в Kafka: `kafka_messages_delta`
- `scripts/kafka_produce_load.sh` — чистая нагрузка на Kafka producer.

## Требования

- Поднят `api-gateway` на `http://localhost:8080`.
- Поднят Kafka контейнер `kafka` (или укажи `KAFKA_CONTAINER`).
- Установлены `k6` и `docker`.

## 1) Нагрузка + RPS + количество сообщений в Kafka

```bash
TOPIC=order-events \
VUS=50 DURATION=1m \
K6_SCRIPT=tests/load/api_gateway_goods_read.js \
bash tests/load/scripts/run_rps_and_kafka.sh
```

Скрипт сохранит результат в `tests/load/results/summary_*.txt`.

## 2) Чистая нагрузка на Kafka

```bash
TOPIC=order-events MESSAGES=20000 \
bash tests/load/scripts/kafka_produce_load.sh
```

Вывод:
- `messages_sent`
- `messages_per_sec_avg`

## 3) Прямой запуск k6 (если нужен)

```bash
k6 run tests/load/api_gateway_smoke.js
k6 run tests/load/api_gateway_goods_read.js
```

## Интерпретация

- `rps_avg` — среднее количество HTTP-запросов в секунду.
- `kafka_messages_delta` — сколько новых сообщений появилось в топике за время теста.
- `kafka_msgs_per_sec_avg` — средняя скорость прироста сообщений в Kafka.
- `http_req_failed` и `checks` — качество ответов под нагрузкой.
