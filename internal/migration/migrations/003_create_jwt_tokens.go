package migrations

import (
	"go-api-starter/internal/domain"
	"gorm.io/gorm"
)

type CreateJWTTokensTable struct{}

func NewCreateJWTTokensTable() *CreateJWTTokensTable {
	return &CreateJWTTokensTable{}
}

func (m *CreateJWTTokensTable) Name() string {
	return "003_create_jwt_tokens_table"
}

func (m *CreateJWTTokensTable) Up(db *gorm.DB) error {
	return db.AutoMigrate(&domain.JWTToken{})
}
