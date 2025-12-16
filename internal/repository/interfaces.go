package repository

import (
	"context"
	"go-api-starter/internal/domain"
)

// UserRepository interface
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
}

// ProductRepository interface
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uint) (*domain.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]domain.Product, int64, error)
	GetByName(ctx context.Context, name string) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uint) error
	UpdateStock(ctx context.Context, id uint, quantity int) error
}

// JWTTokenRepository interface
type JWTTokenRepository interface {
	Save(ctx context.Context, token *domain.JWTToken) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.JWTToken, error)
	GetByAccessToken(ctx context.Context, accessToken string) (*domain.JWTToken, error)
	DeleteByRefreshToken(ctx context.Context, refreshToken string) error
	DeleteExpiredTokens(ctx context.Context) error
	GetByUserID(ctx context.Context, userID uint) ([]domain.JWTToken, error)
}
