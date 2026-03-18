#!/usr/bin/env bash
set -euo pipefail

# Генерирует нагрузку на Kafka и считает скорость публикации сообщений.

KAFKA_CONTAINER="${KAFKA_CONTAINER:-kafka}"
BROKER="${BROKER:-localhost:9092}"
TOPIC="${TOPIC:-order-events}"
MESSAGES="${MESSAGES:-10000}"

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "Требуется команда: $1" >&2
    exit 1
  }
}

require_cmd docker
require_cmd awk

if ! docker ps --format '{{.Names}}' | grep -qx "$KAFKA_CONTAINER"; then
  echo "Kafka контейнер '$KAFKA_CONTAINER' не запущен" >&2
  exit 1
fi

START_TS="$(date +%s)"

i=1
while [ "$i" -le "$MESSAGES" ]; do
  printf '{"event_type":"order.created","order_id":%d,"user_id":1,"status":"paid","total_price":199.99}\n' "$i"
  i=$((i + 1))
done | docker exec -i "$KAFKA_CONTAINER" kafka-console-producer \
  --broker-list "$BROKER" \
  --topic "$TOPIC" >/dev/null

END_TS="$(date +%s)"
ELAPSED=$((END_TS - START_TS))
if [ "$ELAPSED" -le 0 ]; then
  ELAPSED=1
fi

MSG_PER_SEC=$(awk -v m="$MESSAGES" -v s="$ELAPSED" 'BEGIN { printf "%.2f", m/s }')

echo "=== KAFKA PRODUCE RESULT ==="
echo "topic: $TOPIC"
echo "messages_sent: $MESSAGES"
echo "elapsed_sec: $ELAPSED"
echo "messages_per_sec_avg: $MSG_PER_SEC"
