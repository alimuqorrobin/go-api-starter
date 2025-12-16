package service

import (
	"context"
	"go-api-starter/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
	GetUser(ctx context.Context, id uint) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, id uint, req *domain.UpdateUserRequest) (*domain.User, error)
	DeleteUser(ctx context.Context, id uint) error
}

type ProductService interface {
	CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error)
	GetProduct(ctx context.Context, id uint) (*domain.Product, error)
	GetAllProducts(ctx context.Context, limit, offset int) ([]domain.Product, int64, error)
	UpdateProduct(ctx context.Context, id uint, req *domain.UpdateProductRequest) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id uint) error
	GetProductByName(ctx context.Context, name string) (*domain.Product, error)
}

type AuthService interface {
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenResponse, error)
	RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.TokenResponse, error)
	Logout(ctx context.Context, userID uint) error
	ValidateAccessToken(ctx context.Context, token string) (uint, error)
}
