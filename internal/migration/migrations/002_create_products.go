package migrations

import (
	"go-api-starter/internal/domain"
	"gorm.io/gorm"
)

type CreateProductsTable struct{}

func NewCreateProductsTable() *CreateProductsTable {
	return &CreateProductsTable{}
}

func (m *CreateProductsTable) Name() string {
	return "002_create_products_table"
}

func (m *CreateProductsTable) Up(db *gorm.DB) error {
	return db.AutoMigrate(&domain.Product{})
}
