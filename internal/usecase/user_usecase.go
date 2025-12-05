package usecase

import (
	"context"
	"errors"

	"go-api-starter/internal/domain/user"
	"go-api-starter/internal/infrastructure/database"
	"go-api-starter/internal/infrastructure/logger"

	"gorm.io/gorm"
)

type UserUsecase interface {
	GetUserByID(ctx context.Context, id uint) (*user.User, error)
	CreateUser(ctx context.Context, name, email string) (*user.User, error)
}

type userUsecase struct {
	repo user.Repository
	db   *gorm.DB
}

func NewUserUsecase(repo user.Repository, db *gorm.DB) UserUsecase {
	return &userUsecase{
		repo: repo,
		db:   db,
	}
}

func (uc *userUsecase) GetUserByID(ctx context.Context, id uint) (*user.User, error) {
	u, err := uc.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return u, nil
}

func (uc *userUsecase) CreateUser(ctx context.Context, name, email string) (*user.User, error) {
	userEntity := &user.User{
		Name:  name,
		Email: email,
	}

	err := database.WithTransaction(ctx, uc.db, func(tx *gorm.DB) error {
		// Inject tx DB ke repository jika repository Anda mendukung dynamic DB (opsional)
		return uc.repo.Create(userEntity)
	})

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return userEntity, nil
}
