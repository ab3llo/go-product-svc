package services

import (
	"context"
	"net/http"

	"github.com/ab3llo/go-product-svc/pkg/db"
	"github.com/ab3llo/go-product-svc/pkg/models"
	"github.com/ab3llo/go-product-svc/pkg/product/pb"
	"github.com/google/uuid"
)

type Server struct {
	pb.UnimplementedProductServiceServer
	DbConnection db.DatabaseConnection
}

func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	var product models.Product

	product.Id = uuid.New().String()
	product.Name = req.Name
	product.Stock = req.Stock
	product.Price = req.Price

	if result := s.DbConnection.DB.Create(&product); result.Error != nil {
		return &pb.CreateProductResponse{
			Status: http.StatusConflict,
			Error:  result.Error.Error(),
		}, nil
	}

	return &pb.CreateProductResponse{
		Status: http.StatusCreated,
		Id:     product.Id,
	}, nil
}

func (s *Server) FindOne(ctx context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error) {
	var product models.Product

	if result := s.DbConnection.DB.First(&product, "id = ?", req.Id); result.Error != nil {
		return &pb.FindOneResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	data := &pb.FindOneData{
		Id:    product.Id,
		Name:  product.Name,
		Stock: product.Stock,
		Price: product.Price,
	}

	return &pb.FindOneResponse{
		Status: http.StatusOK,
		Data:   data,
	}, nil
}

func (s *Server) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
	var product models.Product

	if result := s.DbConnection.DB.First(&product, "id = ?", req.Id); result.Error != nil {
		return &pb.DecreaseStockResponse{
				Status: http.StatusNotFound,
				Error:  result.Error.Error(),
			},
			nil
	}

	if product.Stock <= 0 {
		return &pb.DecreaseStockResponse{
			Status: http.StatusBadRequest,
			Error:  "Stock too low",
		}, nil
	}

	var log models.StockDecreaseLog

	if result := s.DbConnection.DB.Where(&models.StockDecreaseLog{OrderId: req.OrderId}).First(&log); result.Error == nil {
		return &pb.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  "Stock already decreased",
		}, nil
	}

	product.Stock = product.Stock - 1
	s.DbConnection.DB.Save(&product)

	log.OrderId = req.OrderId
	log.ProductRefer = product.Id

	s.DbConnection.DB.Create(&log)

	return &pb.DecreaseStockResponse{
		Status: http.StatusOK,
	}, nil
}
