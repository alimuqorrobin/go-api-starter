package repository

import (
	"context"
	"errors"
	"go-api-starter/internal/database"
	"go-api-starter/internal/domain"
	"gorm.io/gorm"
)

type productRepository struct {
	db database.Database
}

func NewProductRepository(db database.Database) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	result := r.db.GetDB().WithContext(ctx).Create(product)

	if result.Error != nil {
		return domain.ErrDatabaseConnection
	}

	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	product := &domain.Product{}
	result := r.db.GetDB().WithContext(ctx).First(product, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, domain.ErrDatabaseConnection
	}

	return product, nil
}

func (r *productRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	// Get total count
	if err := r.db.GetDB().WithContext(ctx).Model(&domain.Product{}).Count(&total).Error; err != nil {
		return nil, 0, domain.ErrDatabaseConnection
	}

	// Get paginated data
	result := r.db.GetDB().WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&products)

	if result.Error != nil {
		return nil, 0, domain.ErrDatabaseConnection
	}

	return products, total, nil
}

func (r *productRepository) GetByName(ctx context.Context, name string) (*domain.Product, error) {
	product := &domain.Product{}
	result := r.db.GetDB().WithContext(ctx).Where("name = ?", name).First(product)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, domain.ErrDatabaseConnection
	}

	return product, nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	result := r.db.GetDB().WithContext(ctx).Save(product)

	if result.Error != nil {
		return domain.ErrDatabaseConnection
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.GetDB().WithContext(ctx).Delete(&domain.Product{}, id)

	if result.Error != nil {
		return domain.ErrDatabaseConnection
	}

	if result.RowsAffected == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func (r *productRepository) UpdateStock(ctx context.Context, id uint, quantity int) error {
	result := r.db.GetDB().WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", quantity))

	if result.Error != nil {
		return domain.ErrDatabaseConnection
	}

	if result.RowsAffected == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}
