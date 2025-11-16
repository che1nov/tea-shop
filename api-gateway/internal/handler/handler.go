package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/che1nov/tea-shop/shared/pb"
)

type APIHandler struct {
	usersClient    pb.UsersServiceClient
	goodsClient    pb.GoodsServiceClient
	ordersClient   pb.OrdersServiceClient
	paymentsClient pb.PaymentsServiceClient
	deliveryClient pb.DeliveryServiceClient

	// Храним соединения для graceful shutdown
	usersConn    *grpc.ClientConn
	goodsConn    *grpc.ClientConn
	ordersConn   *grpc.ClientConn
	paymentsConn *grpc.ClientConn
	deliveryConn *grpc.ClientConn
}

func New(
	usersService string,
	goodsService string,
	ordersService string,
	paymentsService string,
	deliveryService string,
) (*APIHandler, error) {
	usersConn, err := grpc.Dial(usersService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	goodsConn, err := grpc.Dial(goodsService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		usersConn.Close()
		return nil, err
	}

	ordersConn, err := grpc.Dial(ordersService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		usersConn.Close()
		goodsConn.Close()
		return nil, err
	}

	paymentsConn, err := grpc.Dial(paymentsService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		usersConn.Close()
		goodsConn.Close()
		ordersConn.Close()
		return nil, err
	}

	deliveryConn, err := grpc.Dial(deliveryService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		usersConn.Close()
		goodsConn.Close()
		ordersConn.Close()
		paymentsConn.Close()
		return nil, err
	}

	return &APIHandler{
		usersClient:    pb.NewUsersServiceClient(usersConn),
		goodsClient:    pb.NewGoodsServiceClient(goodsConn),
		ordersClient:   pb.NewOrdersServiceClient(ordersConn),
		paymentsClient: pb.NewPaymentsServiceClient(paymentsConn),
		deliveryClient: pb.NewDeliveryServiceClient(deliveryConn),
		usersConn:      usersConn,
		goodsConn:      goodsConn,
		ordersConn:     ordersConn,
		paymentsConn:   paymentsConn,
		deliveryConn:   deliveryConn,
	}, nil
}

// Close закрывает все gRPC соединения
func (h *APIHandler) Close() error {
	var errs []error

	if h.usersConn != nil {
		if err := h.usersConn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if h.goodsConn != nil {
		if err := h.goodsConn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if h.ordersConn != nil {
		if err := h.ordersConn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if h.paymentsConn != nil {
		if err := h.paymentsConn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if h.deliveryConn != nil {
		if err := h.deliveryConn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}

	return nil
}

// RegisterUser регистрирует нового пользователя
// @Summary      Регистрация нового пользователя
// @Description  Создает нового пользователя в системе
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "Данные регистрации"  example({"email":"user@example.com","name":"Иван Иванов","password":"password123"})
// @Success      201      {object}  object  "Пользователь создан"
// @Failure      400      {object}  object  "Ошибка валидации"
// @Failure      500      {object}  object  "Внутренняя ошибка сервера"
// @Router       /auth/register [post]
func (h *APIHandler) RegisterUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.usersClient.CreateUser(context.Background(), &pb.CreateUserRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login выполняет вход в систему
// @Summary      Вход в систему
// @Description  Аутентификация пользователя и получение JWT токена
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "Данные для входа"  example({"email":"user@example.com","password":"password123"})
// @Success      200      {object}  object  "Успешный вход"
// @Failure      401      {object}  object  "Неверный email или пароль"
// @Failure      400      {object}  object  "Ошибка валидации"
// @Router       /auth/login [post]
func (h *APIHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.usersClient.Login(context.Background(), &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUser возвращает информацию о текущем пользователе
// @Summary      Получить информацию о текущем пользователе
// @Description  Возвращает информацию о авторизованном пользователе
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  object  "Информация о пользователе"
// @Failure      401  {object}  object  "Не авторизован"
// @Failure      500  {object}  object  "Внутренняя ошибка сервера"
// @Router       /users/me [get]
func (h *APIHandler) GetUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.usersClient.GetUser(context.Background(), &pb.GetUserRequest{
		UserId: userID.(int64),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateGood создает новый товар
// @Summary      Создать товар
// @Description  Создает новый товар. Требует роль администратора (role: "admin") в JWT токене.
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "Данные товара"  example({"name":"Зеленый чай","description":"Высококачественный зеленый чай","price":299.99,"stock":50})
// @Success      201      {object}  object  "Товар создан"
// @Failure      400      {object}  object  "Ошибка валидации"
// @Failure      401      {object}  object  "Не авторизован"
// @Failure      403      {object}  object  "Доступ запрещен: требуется роль администратора"
// @Failure      500      {object}  object  "Внутренняя ошибка сервера"
// @Router       /admin/goods [post]
func (h *APIHandler) CreateGood(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Price       float64 `json:"price" binding:"required,min=0"`
		Stock       int32   `json:"stock" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	good, err := h.goodsClient.CreateGood(context.Background(), &pb.CreateGoodRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, good)
}

// UpdateGood обновляет информацию о товаре
// @Summary      Обновить товар
// @Description  Обновляет информацию о товаре. Требует роль администратора (role: "admin") в JWT токене.
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      int     true  "ID товара"
// @Param        request  body      object  true  "Данные для обновления"  example({"name":"Зеленый чай","description":"Обновленное описание","price":349.99,"stock":60})
// @Success      200      {object}  object  "Товар обновлен"
// @Failure      400      {object}  object  "Ошибка валидации"
// @Failure      404      {object}  object  "Товар не найден"
// @Failure      401      {object}  object  "Не авторизован"
// @Failure      403      {object}  object  "Доступ запрещен: требуется роль администратора"
// @Failure      500      {object}  object  "Внутренняя ошибка сервера"
// @Router       /admin/goods/{id} [put]
func (h *APIHandler) UpdateGood(c *gin.Context) {
	goodID := c.Param("id")
	goodIDInt, _ := strconv.ParseInt(goodID, 10, 64)

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int32   `json:"stock"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	good, err := h.goodsClient.UpdateGood(context.Background(), &pb.UpdateGoodRequest{
		Id:          goodIDInt,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if good == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "good not found"})
		return
	}

	c.JSON(http.StatusOK, good)
}

// DeleteGood удаляет товар
// @Summary      Удалить товар
// @Description  Удаляет товар из каталога. Требует роль администратора (role: "admin") в JWT токене.
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "ID товара"
// @Success      200  {object}  object  "Товар удален"
// @Failure      400  {object}  object  "Ошибка удаления"
// @Failure      404  {object}  object  "Товар не найден"
// @Failure      401  {object}  object  "Не авторизован"
// @Failure      403  {object}  object  "Доступ запрещен: требуется роль администратора"
// @Failure      500  {object}  object  "Внутренняя ошибка сервера"
// @Router       /admin/goods/{id} [delete]
func (h *APIHandler) DeleteGood(c *gin.Context) {
	goodID := c.Param("id")
	goodIDInt, err := strconv.ParseInt(goodID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid good id"})
		return
	}

	response, err := h.goodsClient.DeleteGood(context.Background(), &pb.DeleteGoodRequest{
		GoodId: goodIDInt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !response.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": response.Message})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListGoods возвращает список товаров
// @Summary      Получить список товаров
// @Description  Возвращает список товаров с пагинацией
// @Tags         Goods
// @Produce      json
// @Param        limit   query     int     false  "Количество товаров"  default(10)
// @Param        offset  query     int     false  "Смещение"  default(0)
// @Success      200     {object}  object  "Список товаров"
// @Failure      500     {object}  object  "Внутренняя ошибка сервера"
// @Router       /goods [get]
func (h *APIHandler) ListGoods(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	limitInt, _ := strconv.ParseInt(limit, 10, 32)
	offsetInt, _ := strconv.ParseInt(offset, 10, 32)

	goods, err := h.goodsClient.ListGoods(context.Background(), &pb.ListGoodsRequest{
		Limit:  int32(limitInt),
		Offset: int32(offsetInt),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goods)
}

// GetGood возвращает товар по ID
// @Summary      Получить товар по ID
// @Description  Возвращает детальную информацию о товаре
// @Tags         Goods
// @Produce      json
// @Param        id   path      int     true  "ID товара"
// @Success      200  {object}  object  "Информация о товаре"
// @Failure      404  {object}  object  "Товар не найден"
// @Failure      500  {object}  object  "Внутренняя ошибка сервера"
// @Router       /goods/{id} [get]
func (h *APIHandler) GetGood(c *gin.Context) {
	goodID := c.Param("id")
	goodIDInt, _ := strconv.ParseInt(goodID, 10, 64)

	good, err := h.goodsClient.GetGood(context.Background(), &pb.GetGoodRequest{
		GoodId: goodIDInt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if good == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "good not found"})
		return
	}

	c.JSON(http.StatusOK, good)
}

// CreateOrder создает новый заказ
// @Summary      Создать заказ
// @Description  Создает новый заказ для текущего пользователя
// @Tags         Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "Данные заказа"  example({"items":[{"good_id":1,"quantity":2,"price":299.99}],"address":"г. Москва, ул. Примерная, д. 1, кв. 10"})
// @Success      201      {object}  object  "Заказ создан"
// @Failure      400      {object}  object  "Ошибка валидации"
// @Failure      401      {object}  object  "Не авторизован"
// @Failure      500      {object}  object  "Внутренняя ошибка сервера"
// @Router       /orders [post]
func (h *APIHandler) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Items []struct {
			GoodID   int64   `json:"good_id" binding:"required"`
			Quantity int32   `json:"quantity" binding:"required"`
			Price    float64 `json:"price"`
		} `json:"items" binding:"required"`
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items := make([]*pb.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = &pb.OrderItem{
			GoodId:   item.GoodID,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	order, err := h.ordersClient.CreateOrder(context.Background(), &pb.CreateOrderRequest{
		UserId:  userID.(int64),
		Items:   items,
		Address: req.Address,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrder возвращает заказ по ID
// @Summary      Получить заказ по ID
// @Description  Возвращает информацию о заказе
// @Tags         Orders
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int     true  "ID заказа"
// @Success      200  {object}  object  "Информация о заказе"
// @Failure      404  {object}  object  "Заказ не найден"
// @Failure      401  {object}  object  "Не авторизован"
// @Failure      500  {object}  object  "Внутренняя ошибка сервера"
// @Router       /orders/{id} [get]
func (h *APIHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	orderIDInt, _ := strconv.ParseInt(orderID, 10, 64)

	order, err := h.ordersClient.GetOrder(context.Background(), &pb.GetOrderRequest{
		OrderId: orderIDInt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetPayment возвращает информацию о платеже
// @Summary      Получить информацию о платеже
// @Description  Возвращает информацию о платеже по ID
// @Tags         Payments
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int     true  "ID платежа"
// @Success      200  {object}  object  "Информация о платеже"
// @Failure      404  {object}  object  "Платеж не найден"
// @Failure      401  {object}  object  "Не авторизован"
// @Failure      500  {object}  object  "Внутренняя ошибка сервера"
// @Router       /payments/{id} [get]
func (h *APIHandler) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")
	paymentIDInt, _ := strconv.ParseInt(paymentID, 10, 64)

	payment, err := h.paymentsClient.GetPayment(context.Background(), &pb.GetPaymentRequest{
		PaymentId: paymentIDInt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if payment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// CreateDelivery создает доставку
// @Summary      Создать доставку
// @Description  Создает новую доставку для заказа
// @Tags         Deliveries
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "Данные доставки"  example({"order_id":1,"address":"г. Москва, ул. Примерная, д. 1, кв. 10"})
// @Success      201      {object}  object  "Доставка создана"
// @Failure      400      {object}  object  "Ошибка валидации"
// @Failure      401      {object}  object  "Не авторизован"
// @Failure      500      {object}  object  "Внутренняя ошибка сервера"
// @Router       /deliveries [post]
func (h *APIHandler) CreateDelivery(c *gin.Context) {
	var req struct {
		OrderID int64  `json:"order_id" binding:"required"`
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delivery, err := h.deliveryClient.CreateDelivery(context.Background(), &pb.CreateDeliveryRequest{
		OrderId: req.OrderID,
		Address: req.Address,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, delivery)
}

// GetDelivery возвращает информацию о доставке
// @Summary      Получить информацию о доставке
// @Description  Возвращает информацию о доставке по ID
// @Tags         Deliveries
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int     true  "ID доставки"
// @Success      200  {object}  object  "Информация о доставке"
// @Failure      404  {object}  object  "Доставка не найдена"
// @Failure      401  {object}  object  "Не авторизован"
// @Failure      500  {object}  object  "Внутренняя ошибка сервера"
// @Router       /deliveries/{id} [get]
func (h *APIHandler) GetDelivery(c *gin.Context) {
	deliveryID := c.Param("id")
	deliveryIDInt, _ := strconv.ParseInt(deliveryID, 10, 64)

	delivery, err := h.deliveryClient.GetDelivery(context.Background(), &pb.GetDeliveryRequest{
		DeliveryId: deliveryIDInt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if delivery == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "delivery not found"})
		return
	}

	c.JSON(http.StatusOK, delivery)
}

// ListDeliveries возвращает список доставок (только для админа)
// @Summary      Получить список доставок
// @Description  Возвращает список всех доставок с возможностью фильтрации по статусу. Требует роль администратора.
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Param        status  query     string  false  "Фильтр по статусу (pending, in_transit, delivered, cancelled)"
// @Param        limit   query     int     false  "Лимит результатов (по умолчанию: 100)"
// @Param        offset  query     int     false  "Смещение для пагинации (по умолчанию: 0)"
// @Success      200     {object}  object  "Список доставок"
// @Failure      401     {object}  object  "Не авторизован"
// @Failure      403     {object}  object  "Доступ запрещен: требуется роль администратора"
// @Failure      500     {object}  object  "Внутренняя ошибка сервера"
// @Router       /admin/deliveries [get]
func (h *APIHandler) ListDeliveries(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")
	status := c.Query("status")

	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	offset, _ := strconv.ParseInt(offsetStr, 10, 32)

	response, err := h.deliveryClient.ListDeliveries(context.Background(), &pb.ListDeliveriesRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
		Status: status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateDeliveryStatus обновляет статус доставки (только для админа)
// @Summary      Обновить статус доставки
// @Description  Обновляет статус доставки. Требует роль администратора.
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      int     true  "ID доставки"
// @Param        request  body      object  true  "Новый статус"  example({"status":"in_transit"})
// @Success      200      {object}  object  "Статус доставки обновлен"
// @Failure      400      {object}  object  "Ошибка валидации"
// @Failure      401      {object}  object  "Не авторизован"
// @Failure      403      {object}  object  "Доступ запрещен: требуется роль администратора"
// @Failure      500      {object}  object  "Внутренняя ошибка сервера"
// @Router       /admin/deliveries/{id}/status [put]
func (h *APIHandler) UpdateDeliveryStatus(c *gin.Context) {
	deliveryID := c.Param("id")
	deliveryIDInt, err := strconv.ParseInt(deliveryID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid delivery id"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delivery, err := h.deliveryClient.UpdateDeliveryStatus(context.Background(), &pb.UpdateDeliveryStatusRequest{
		DeliveryId: deliveryIDInt,
		Status:     req.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, delivery)
}
