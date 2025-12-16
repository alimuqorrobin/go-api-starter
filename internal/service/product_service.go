package service

import (
	"context"
	"go-api-starter/internal/domain"
	"go-api-starter/internal/repository"
)

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error) {
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) GetProduct(ctx context.Context, id uint) (*domain.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *productService) GetAllProducts(ctx context.Context, limit, offset int) ([]domain.Product, int64, error) {
	return s.productRepo.GetAll(ctx, limit, offset)
}

func (s *productService) UpdateProduct(ctx context.Context, id uint, req *domain.UpdateProductRequest) (*domain.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	return s.productRepo.Delete(ctx, id)
}

func (s *productService) GetProductByName(ctx context.Context, name string) (*domain.Product, error) {
	return s.productRepo.GetByName(ctx, name)
}
