# Integration Tests

Интеграционные тесты проверяют repository layer на реальном PostgreSQL из `docker-compose.yml`.

Запуск:

```bash
bash tests/integration/run.sh
```

Скрипт поднимает PostgreSQL-контейнеры, ждет готовности и запускает:

```bash
RUN_INTEGRATION_TESTS=1 go test -tags=integration ./internal/repository
```

Тесты чистят свои таблицы через `TRUNCATE ... RESTART IDENTITY CASCADE`, поэтому запускайте их только против локальной test/dev базы.

Можно переопределить DSN:

```bash
USERS_TEST_DATABASE_URL="user=user password=password dbname=users_db host=localhost port=5432 sslmode=disable"
GOODS_TEST_DATABASE_URL="user=user password=password dbname=goods_db host=localhost port=5433 sslmode=disable"
ORDERS_TEST_DATABASE_URL="user=user password=password dbname=orders_db host=localhost port=5434 sslmode=disable"
PAYMENTS_TEST_DATABASE_URL="user=user password=password dbname=payments_db host=localhost port=5435 sslmode=disable"
DELIVERY_TEST_DATABASE_URL="user=user password=password dbname=deliveries_db host=localhost port=5436 sslmode=disable"
```
