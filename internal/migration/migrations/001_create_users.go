package migrations

import (
	"go-api-starter/internal/domain"
	"gorm.io/gorm"
)

type CreateUsersTable struct{}

func NewCreateUsersTable() *CreateUsersTable {
	return &CreateUsersTable{}
}

func (m *CreateUsersTable) Name() string {
	return "001_create_users_table"
}

func (m *CreateUsersTable) Up(db *gorm.DB) error {
	return db.AutoMigrate(&domain.User{})
}
