#!/usr/bin/env bash
set -euo pipefail

# Запускает k6 и измеряет:
# 1) RPS по количеству метрик http_reqs
# 2) Прирост количества сообщений в Kafka по offset

TOPIC="${TOPIC:-order-events}"
BROKER="${BROKER:-localhost:9092}"
KAFKA_CONTAINER="${KAFKA_CONTAINER:-kafka}"
K6_SCRIPT="${K6_SCRIPT:-tests/load/api_gateway_goods_read.js}"
K6_ARGS="${K6_ARGS:-}"
OUT_DIR="${OUT_DIR:-tests/load/results}"

mkdir -p "$OUT_DIR"
TS="$(date +%Y%m%d_%H%M%S)"
K6_OUT_JSON="$OUT_DIR/k6_${TS}.json"
SUMMARY_OUT="$OUT_DIR/summary_${TS}.txt"

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "Требуется команда: $1" >&2
    exit 1
  }
}

require_cmd docker
require_cmd k6
require_cmd awk
require_cmd grep

get_topic_offset_sum() {
  if ! docker exec "$KAFKA_CONTAINER" sh -lc "timeout 5 kafka-topics --bootstrap-server $BROKER --list" \
    | grep -qx "$TOPIC"; then
    echo 0
    return 0
  fi

  local offsets
  offsets=$(docker exec "$KAFKA_CONTAINER" sh -lc \
    "timeout 5 kafka-run-class kafka.tools.GetOffsetShell --broker-list $BROKER --topic $TOPIC --time -1" || true)

  if [ -z "$offsets" ]; then
    echo 0
    return 0
  fi

  printf '%s\n' "$offsets" | awk -F: '{sum += $3} END {print sum+0}'
}

if ! docker ps --format '{{.Names}}' | grep -qx "$KAFKA_CONTAINER"; then
  echo "Kafka контейнер '$KAFKA_CONTAINER' не запущен" >&2
  exit 1
fi

START_TS="$(date +%s)"
START_OFFSET="$(get_topic_offset_sum)"

echo "=== START ==="
echo "topic: $TOPIC"
echo "start_offset: $START_OFFSET"
echo "k6_script: $K6_SCRIPT"

# shellcheck disable=SC2086
k6 run $K6_ARGS --out "json=$K6_OUT_JSON" "$K6_SCRIPT"

END_TS="$(date +%s)"
END_OFFSET="$(get_topic_offset_sum)"

DELTA_MSG=$((END_OFFSET - START_OFFSET))
ELAPSED=$((END_TS - START_TS))
if [ "$ELAPSED" -le 0 ]; then
  ELAPSED=1
fi
MSG_PER_SEC=$(awk -v m="$DELTA_MSG" -v s="$ELAPSED" 'BEGIN { printf "%.2f", m/s }')

HTTP_REQS=$(grep -c '"metric":"http_reqs"' "$K6_OUT_JSON" || true)
RPS=$(awk -v c="$HTTP_REQS" -v s="$ELAPSED" 'BEGIN { printf "%.2f", c/s }')

{
  echo "=== RESULT ==="
  echo "topic: $TOPIC"
  echo "elapsed_sec: $ELAPSED"
  echo "http_reqs_total: $HTTP_REQS"
  echo "rps_avg: $RPS"
  echo "start_offset: $START_OFFSET"
  echo "end_offset: $END_OFFSET"
  echo "kafka_messages_delta: $DELTA_MSG"
  echo "kafka_msgs_per_sec_avg: $MSG_PER_SEC"
  echo "k6_json: $K6_OUT_JSON"
} | tee "$SUMMARY_OUT"

echo "summary_file: $SUMMARY_OUT"
