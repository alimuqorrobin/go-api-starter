package service

import (
	"context"
	"fmt"
	"go-api-starter/internal/database"
	"go-api-starter/internal/domain"
	"go-api-starter/internal/repository"
	"gorm.io/gorm"
)

// JANGAN DECLARE INTERFACE DI SINI
// Sudah ada di interfaces.go
// type ProductService interface { ... }

type productService struct {
	db          *gorm.DB
	productRepo repository.ProductRepository
}

func NewProductService(db *gorm.DB, productRepo repository.ProductRepository) ProductService {
	return &productService{
		db:          db,
		productRepo: productRepo,
	}
}

// CreateProduct creates new product
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

// GetProduct gets product by ID
func (s *productService) GetProduct(ctx context.Context, id uint) (*domain.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

// GetAllProducts gets all products with pagination
func (s *productService) GetAllProducts(ctx context.Context, limit, offset int) ([]domain.Product, int64, error) {
	return s.productRepo.GetAll(ctx, limit, offset)
}

// UpdateProduct updates product
func (s *productService) UpdateProduct(ctx context.Context, id uint, req *domain.UpdateProductRequest) (*domain.Product, error) {
	var product *domain.Product

	err := database.WithTx(s.db, ctx, func(tx *gorm.DB) error {
		product = &domain.Product{}
		if err := database.LockForUpdate(tx, product, id); err != nil {
			return fmt.Errorf("product not found: %w", err)
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

		if err := tx.Save(product).Error; err != nil {
			return fmt.Errorf("failed to update product: %w", err)
		}

		return nil
	})

	return product, err
}

// DeleteProduct deletes product
func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	return database.WithTx(s.db, ctx, func(tx *gorm.DB) error {
		product := &domain.Product{}
		if err := database.LockForUpdate(tx, product, id); err != nil {
			return fmt.Errorf("product not found: %w", err)
		}

		if err := tx.Delete(product).Error; err != nil {
			return fmt.Errorf("failed to delete product: %w", err)
		}

		return nil
	})
}

// GetProductByName gets product by name
func (s *productService) GetProductByName(ctx context.Context, name string) (*domain.Product, error) {
	return s.productRepo.GetByName(ctx, name)
}

// ============================================
// SAFE STOCK OPERATIONS WITH LOCKING
// ============================================

// DeductStock safely deducts product stock
func (s *productService) DeductStock(ctx context.Context, productID uint, quantity int) error {
	return database.WithTx(s.db, ctx, func(tx *gorm.DB) error {
		product := &domain.Product{}
		if err := database.LockForUpdate(tx, product, productID); err != nil {
			return fmt.Errorf("product not found")
		}

		if product.Stock < quantity {
			return fmt.Errorf("insufficient stock: have %d, need %d", product.Stock, quantity)
		}

		if err := tx.Model(product).
			Update("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
			return fmt.Errorf("failed to deduct stock: %w", err)
		}

		return nil
	})
}

// AddStock safely adds product stock
func (s *productService) AddStock(ctx context.Context, productID uint, quantity int) error {
	return database.WithTx(s.db, ctx, func(tx *gorm.DB) error {
		product := &domain.Product{}
		if err := database.LockForUpdate(tx, product, productID); err != nil {
			return fmt.Errorf("product not found")
		}

		if err := tx.Model(product).
			Update("stock", gorm.Expr("stock + ?", quantity)).Error; err != nil {
			return fmt.Errorf("failed to add stock: %w", err)
		}

		return nil
	})
}

// TransferStock safely transfers stock between products
func (s *productService) TransferStock(ctx context.Context, fromProductID, toProductID uint, quantity int) error {
	return database.WithTx(s.db, ctx, func(tx *gorm.DB) error {
		var from, to *domain.Product

		// Lock in order by ID to prevent deadlock
		if fromProductID < toProductID {
			from = &domain.Product{}
			if err := database.LockForUpdate(tx, from, fromProductID); err != nil {
				return fmt.Errorf("from product not found")
			}

			to = &domain.Product{}
			if err := database.LockForUpdate(tx, to, toProductID); err != nil {
				return fmt.Errorf("to product not found")
			}
		} else {
			to = &domain.Product{}
			if err := database.LockForUpdate(tx, to, toProductID); err != nil {
				return fmt.Errorf("to product not found")
			}

			from = &domain.Product{}
			if err := database.LockForUpdate(tx, from, fromProductID); err != nil {
				return fmt.Errorf("from product not found")
			}
		}

		if from.Stock < quantity {
			return fmt.Errorf("insufficient stock in source: have %d, need %d", from.Stock, quantity)
		}

		if err := tx.Model(from).
			Update("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
			return fmt.Errorf("failed to deduct from source: %w", err)
		}

		if err := tx.Model(to).
			Update("stock", gorm.Expr("stock + ?", quantity)).Error; err != nil {
			return fmt.Errorf("failed to add to destination: %w", err)
		}

		return nil
	})
}