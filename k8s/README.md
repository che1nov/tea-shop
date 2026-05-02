# Kubernetes запуск

Это минимальный dev/staging-набор манифестов для локального кластера Kubernetes: приложения, frontend, Kafka, Zookeeper, пять PostgreSQL, Prometheus и Grafana.

## Сборка образов

Из корня проекта:

```sh
docker build -f Dockerfile --build-arg SERVICE=api-gateway -t tea-shop/api-gateway:latest .
docker build -f Dockerfile --build-arg SERVICE=users-service -t tea-shop/users-service:latest .
docker build -f Dockerfile --build-arg SERVICE=goods-service -t tea-shop/goods-service:latest .
docker build -f Dockerfile --build-arg SERVICE=order-service -t tea-shop/order-service:latest .
docker build -f Dockerfile --build-arg SERVICE=payment-service -t tea-shop/payment-service:latest .
docker build -f Dockerfile --build-arg SERVICE=delivery-service -t tea-shop/delivery-service:latest .
docker build -f Dockerfile --build-arg SERVICE=notify-service -t tea-shop/notify-service:latest .
docker build -f frontend/Dockerfile -t tea-shop/frontend:latest frontend
```

Для `kind` нужно загрузить локальные образы в кластер:

```sh
kind load docker-image tea-shop/api-gateway:latest tea-shop/users-service:latest tea-shop/goods-service:latest tea-shop/order-service:latest tea-shop/payment-service:latest tea-shop/delivery-service:latest tea-shop/notify-service:latest tea-shop/frontend:latest
```

Для `minikube` удобнее перед сборкой выполнить:

```sh
eval "$(minikube docker-env)"
```

## Применение

```sh
kubectl apply -k k8s
kubectl -n tea-shop get pods
```

## Локальный доступ

Frontend:

```sh
kubectl -n tea-shop port-forward svc/frontend 5173:80
```

API Gateway:

```sh
kubectl -n tea-shop port-forward svc/api-gateway 8080:8080
```

Prometheus и Grafana:

```sh
kubectl -n tea-shop port-forward svc/prometheus 9090:9090
kubectl -n tea-shop port-forward svc/grafana 3000:3000
```

Grafana: `admin/admin`.

## Секреты

Перед production-запуском поменяй значения в `k8s/secrets.yaml`: `jwt_secret`, `db_password`, `admin_password`, `email_password`.
