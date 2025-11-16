# Поток обработки заказов

Документация описывает полный цикл обработки заказа от создания до отправки.

## 📋 Обзор процесса

```
Пользователь → Frontend → API Gateway → Order Service → 
    ├─ Goods Service (резервирование товаров)
    ├─ Payment Service (обработка платежа)
    └─ Kafka (публикация события)
        └─ Notify Service (уведомление пользователя)
```

## 🔄 Детальный поток

### 1. Создание заказа (Frontend)

**Файл:** `frontend/src/pages/Cart.tsx`

Пользователь нажимает кнопку "Оформить заказ" в корзине:

```typescript
// 1. Получаем товары из корзины
const cartItems = useCartStore((state) => state.items)

// 2. Формируем заказ
const orderItems = cartItems.map(item => ({
  good_id: item.id,
  quantity: item.quantity,
  price: item.price
}))

// 3. Отправляем запрос через API Gateway
POST /api/v1/orders
{
  "items": orderItems
}
```

**Важно:** 
- `user_id` автоматически извлекается из JWT токена (через middleware)
- Товары берутся из корзины (Zustand store)

---

### 2. API Gateway - Обработка запроса

**Файл:** `api-gateway/internal/handler/handler.go`

**Эндпоинт:** `POST /api/v1/orders`

**Метод:** `CreateOrder`

```go
func (h *APIHandler) CreateOrder(c *gin.Context) {
    // 1. Извлекаем user_id из JWT токена (AuthMiddleware)
    userID, _ := c.Get("user_id").(int64)
    
    // 2. Парсим товары из запроса
    var req struct {
        Items []struct {
            GoodID   int64   `json:"good_id"`
            Quantity int32   `json:"quantity"`
            Price    float64 `json:"price"`
        } `json:"items"`
    }
    
    // 3. Вызываем order-service через gRPC
    order, err := h.ordersClient.CreateOrder(ctx, &pb.CreateOrderRequest{
        UserId: userID,
        Items:  items,
    })
    
    // 4. Возвращаем созданный заказ
    c.JSON(http.StatusCreated, order)
}
```

**Требования:**
- ✅ JWT токен (через `AuthMiddleware`)
- ✅ Товары в корзине

---

### 3. Order Service - Создание заказа

**Файл:** `order-service/internal/service/service.go`

**Метод:** `CreateOrder`

#### Шаг 3.1: Проверка товаров

```go
// Для каждого товара в заказе:
for _, item := range req.Items {
    // 1. Получаем информацию о товаре
    good, err := s.goodsServiceConn.GetGood(ctx, &pb.GetGoodRequest{
        GoodId: item.GoodID,
    })
    
    // 2. Проверяем наличие товара на складе
    checkResp, err := s.goodsServiceConn.CheckStock(ctx, &pb.CheckStockRequest{
        GoodId:   item.GoodID,
        Quantity: item.Quantity,
    })
    
    if !checkResp.Available {
        return nil, errors.New("товара недостаточно")
    }
    
    // 3. Рассчитываем цену
    item.Price = good.Price
    totalPrice += good.Price * float64(item.Quantity)
}
```

#### Шаг 3.2: Создание заказа в БД

```go
order := &model.Order{
    UserID:     req.UserID,
    Items:      req.Items,
    Status:     "pending",
    TotalPrice: totalPrice,
}

// Сохраняем заказ в PostgreSQL
err := s.repo.CreateOrder(ctx, order)
```

**Статус:** `pending` (заказ создан, ожидает обработки)

#### Шаг 3.3: Резервирование товаров

```go
// Для каждого товара резервируем количество
for _, item := range req.Items {
    _, err := s.goodsServiceConn.ReserveStock(ctx, &pb.ReserveStockRequest{
        GoodId:   item.GoodID,
        Quantity: item.Quantity,
        OrderId:  order.ID,
    })
    
    if err != nil {
        // Откатываем заказ, если резервирование не удалось
        return nil, err
    }
}
```

**Что происходит в Goods Service:**
- Уменьшается `stock` (остаток товара)
- Создается запись в `stock_reservations` (связь с заказом)
- Все выполняется в транзакции

#### Шаг 3.4: Обработка платежа

```go
paymentResp, err := s.paymentServiceConn.ProcessPayment(ctx, &pb.ProcessPaymentRequest{
    OrderId: order.ID,
    Amount:  totalPrice,
    Method:  "card",
})

if paymentResp.Status == "completed" {
    order.Status = "paid"
    s.repo.UpdateOrderStatus(ctx, order.ID, "paid")
} else {
    order.Status = "payment_failed"
    s.repo.UpdateOrderStatus(ctx, order.ID, "payment_failed")
}
```

**Возможные статусы:**
- `paid` - платеж успешен
- `payment_failed` - платеж не прошел

#### Шаг 3.5: Публикация события в Kafka

```go
s.producer.PublishOrderCreated(ctx, &kafka.OrderEvent{
    OrderID:    order.ID,
    UserID:     order.UserID,
    Status:     order.Status,
    TotalPrice: order.TotalPrice,
})
```

**Топик:** `order_created`

**Формат события:**
```json
{
  "order_id": 1,
  "user_id": 1,
  "status": "paid",
  "total_price": 599.98
}
```

---

### 4. Notify Service - Обработка уведомлений

**Файл:** `notify-service/internal/kafka/consumer.go`

#### Шаг 4.1: Получение события из Kafka

```go
// Consumer подписывается на топик "order_created"
message, err := consumer.ReadMessage(ctx)

var event kafka.OrderEvent
json.Unmarshal(message.Value, &event)
```

#### Шаг 4.2: Отправка уведомления

```go
// Отправляем email пользователю
err := s.sendOrderConfirmationEmail(ctx, event.UserID, event.OrderID)
```

**Что отправляется:**
- Подтверждение создания заказа
- Информация о заказе (товары, сумма)
- Статус заказа

---

### 5. Delivery Service - Создание доставки

**После успешного создания заказа можно создать доставку:**

**Эндпоинт:** `POST /api/v1/delivery`

**Метод:** `CreateDelivery`

```go
func (h *APIHandler) CreateDelivery(c *gin.Context) {
    var req struct {
        OrderID int64  `json:"order_id" binding:"required"`
        Address string `json:"address" binding:"required"`
    }
    
    delivery, err := h.deliveryClient.CreateDelivery(ctx, &pb.CreateDeliveryRequest{
        OrderId: req.OrderID,
        Address: req.Address,
    })
    
    c.JSON(http.StatusCreated, delivery)
}
```

**Статусы доставки:**
- `pending` - доставка создана
- `in_transit` - в пути
- `delivered` - доставлено
- `cancelled` - отменено

---

## 📊 Статусы заказа

| Статус | Описание |
|--------|----------|
| `pending` | Заказ создан, ожидает обработки |
| `paid` | Платеж успешен |
| `payment_failed` | Платеж не прошел |
| `processing` | Заказ обрабатывается |
| `completed` | Заказ выполнен |
| `cancelled` | Заказ отменен |

---

## 🔍 Где обрабатываются заказы

### Создание заказа
- **Frontend:** `frontend/src/pages/Cart.tsx` - кнопка "Оформить заказ"
- **API Gateway:** `api-gateway/internal/handler/handler.go` - `CreateOrder`
- **Order Service:** `order-service/internal/service/service.go` - `CreateOrder`

### Резервирование товаров
- **Order Service:** вызывает `goods-service` через gRPC
- **Goods Service:** `goods-service/internal/repository/repository.go` - `ReserveStock`

### Обработка платежа
- **Order Service:** вызывает `payment-service` через gRPC
- **Payment Service:** `payment-service/internal/service/service.go` - `ProcessPayment`

### Уведомления
- **Order Service:** публикует событие в Kafka (топик `order_created`)
- **Notify Service:** `notify-service/internal/kafka/consumer.go` - получает событие и отправляет уведомление

### Доставка
- **API Gateway:** `api-gateway/internal/handler/handler.go` - `CreateDelivery`
- **Delivery Service:** `delivery-service/internal/service/service.go` - `CreateDelivery`

---

## 🔄 Схема потока данных

```
┌─────────────┐
│   Frontend  │
│   (Cart)    │
└──────┬──────┘
       │ POST /api/v1/orders
       ▼
┌─────────────┐
│ API Gateway │
│  (JWT Auth) │
└──────┬──────┘
       │ gRPC CreateOrder
       ▼
┌─────────────┐
│Order Service│
└──────┬──────┘
       │
       ├─► gRPC CheckStock ───► Goods Service
       │
       ├─► gRPC ReserveStock ──► Goods Service
       │                          (уменьшает stock)
       │
       ├─► gRPC ProcessPayment ─► Payment Service
       │                          (обрабатывает платеж)
       │
       └─► Kafka Producer ──────► Kafka (order_created)
                                     │
                                     ▼
                              ┌─────────────┐
                              │Notify Service│
                              │ (Consumer)   │
                              └─────────────┘
```

---

## 🚀 Как протестировать

### 1. Создать заказ через Frontend

1. Откройте `http://localhost:5173`
2. Войдите в систему
3. Добавьте товары в корзину
4. Перейдите в корзину (`/cart`)
5. Нажмите "Оформить заказ"

### 2. Проверить создание заказа

```bash
# Проверить заказ через API Gateway
curl -X GET http://localhost:8080/api/v1/orders/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. Проверить резервирование товаров

```bash
# Проверить остаток товара
curl http://localhost:8080/api/v1/goods/1
```

### 4. Проверить уведомления

```bash
# Проверить логи notify-service
tail -f /tmp/notify-service.log
```

---

## 📝 Важные замечания

1. **Транзакции:** Все операции с заказом выполняются в транзакциях
2. **Резервирование:** Товары резервируются перед созданием платежа
3. **Откат:** При ошибке резервирования или платежа заказ откатывается
4. **Асинхронность:** Уведомления отправляются асинхронно через Kafka
5. **Масштабируемость:** Каждый сервис может масштабироваться независимо

---

## 🔗 Связанные документы

- [Order Service README](./order-service/README.md)
- [Goods Service README](./goods-service/README.md)
- [Payment Service README](./payment-service/README.md)
- [Notify Service README](./notify-service/README.md)
- [Delivery Service README](./delivery-service/README.md)

