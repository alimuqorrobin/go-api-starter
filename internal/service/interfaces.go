package service

import (
	"context"
	"go-api-starter/internal/domain"
)

// UserService interface
type UserService interface {
	CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
	GetUser(ctx context.Context, id uint) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, id uint, req *domain.UpdateUserRequest) (*domain.User, error)
	DeleteUser(ctx context.Context, id uint) error
}

// ProductService interface - UPDATED with locking methods
type ProductService interface {
	CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error)
	GetProduct(ctx context.Context, id uint) (*domain.Product, error)
	GetAllProducts(ctx context.Context, limit, offset int) ([]domain.Product, int64, error)
	UpdateProduct(ctx context.Context, id uint, req *domain.UpdateProductRequest) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id uint) error
	GetProductByName(ctx context.Context, name string) (*domain.Product, error)
	
	// NEW: Safe stock operations dengan locking (race condition protected)
	DeductStock(ctx context.Context, productID uint, quantity int) error
	AddStock(ctx context.Context, productID uint, quantity int) error
	TransferStock(ctx context.Context, fromProductID, toProductID uint, quantity int) error
}

// AuthService interface
type AuthService interface {
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenResponse, error)
	RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.TokenResponse, error)
	Logout(ctx context.Context, userID uint) error
	ValidateAccessToken(ctx context.Context, token string) (uint, error)
}