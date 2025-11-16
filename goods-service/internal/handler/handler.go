package handler

import (
	"context"

	"github.com/che1nov/tea-shop/goods-service/internal/model"
	"github.com/che1nov/tea-shop/goods-service/internal/service"
	pb "github.com/che1nov/tea-shop/shared/pb"
)

type GoodsHandler struct {
	service service.GoodsServiceInterface
	pb.UnimplementedGoodsServiceServer
}

func New(svc service.GoodsServiceInterface) *GoodsHandler {
	return &GoodsHandler{
		service: svc,
	}
}

func (h *GoodsHandler) CreateGood(ctx context.Context, req *pb.CreateGoodRequest) (*pb.Good, error) {
	createReq := &model.CreateGoodRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	good, err := h.service.CreateGood(ctx, createReq)
	if err != nil {
		return nil, err
	}

	return &pb.Good{
		Id:          good.ID,
		Sku:         good.SKU,
		Name:        good.Name,
		Description: good.Description,
		Price:       good.Price,
		Stock:       good.Stock,
		CreatedAt:   good.CreatedAt.Unix(),
	}, nil
}

func (h *GoodsHandler) GetGood(ctx context.Context, req *pb.GetGoodRequest) (*pb.Good, error) {
	good, err := h.service.GetGood(ctx, req.GoodId)
	if err != nil {
		return nil, err
	}

	if good == nil {
		return nil, nil
	}

	return &pb.Good{
		Id:          good.ID,
		Sku:         good.SKU,
		Name:        good.Name,
		Description: good.Description,
		Price:       good.Price,
		Stock:       good.Stock,
		CreatedAt:   good.CreatedAt.Unix(),
	}, nil
}

func (h *GoodsHandler) ListGoods(ctx context.Context, req *pb.ListGoodsRequest) (*pb.ListGoodsResponse, error) {
	goods, err := h.service.ListGoods(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	total, err := h.service.GetTotalGoods(ctx)
	if err != nil {
		return nil, err
	}

	pbGoods := make([]*pb.Good, len(goods))
	for i, good := range goods {
		pbGoods[i] = &pb.Good{
			Id:          good.ID,
			Sku:         good.SKU,
			Name:        good.Name,
			Description: good.Description,
			Price:       good.Price,
			Stock:       good.Stock,
			CreatedAt:   good.CreatedAt.Unix(),
		}
	}

	return &pb.ListGoodsResponse{
		Goods: pbGoods,
		Total: total,
	}, nil
}

func (h *GoodsHandler) CheckStock(ctx context.Context, req *pb.CheckStockRequest) (*pb.CheckStockResponse, error) {
	available, err := h.service.CheckStock(ctx, req.GoodId, req.Quantity)
	if err != nil {
		return nil, err
	}

	return &pb.CheckStockResponse{
		Available: available,
	}, nil
}

func (h *GoodsHandler) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
	success, err := h.service.ReserveStock(ctx, req.GoodId, req.Quantity, req.OrderId)
	if err != nil {
		return &pb.ReserveStockResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	if !success {
		return &pb.ReserveStockResponse{
			Success: false,
			Error:   "insufficient stock",
		}, nil
	}

	return &pb.ReserveStockResponse{
		Success: true,
	}, nil
}

func (h *GoodsHandler) UpdateGood(ctx context.Context, req *pb.UpdateGoodRequest) (*pb.Good, error) {
	updateReq := &model.UpdateGoodRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.Sku,
	}

	good, err := h.service.UpdateGood(ctx, req.Id, updateReq)
	if err != nil {
		return nil, err
	}

	if good == nil {
		return nil, nil
	}

	return &pb.Good{
		Id:          good.ID,
		Sku:         good.SKU,
		Name:        good.Name,
		Description: good.Description,
		Price:       good.Price,
		Stock:       good.Stock,
		CreatedAt:   good.CreatedAt.Unix(),
	}, nil
}

func (h *GoodsHandler) DeleteGood(ctx context.Context, req *pb.DeleteGoodRequest) (*pb.DeleteGoodResponse, error) {
	err := h.service.DeleteGood(ctx, req.GoodId)
	if err != nil {
		return &pb.DeleteGoodResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeleteGoodResponse{
		Success: true,
		Message: "Good deleted successfully",
	}, nil
}
