package repository

import (
	"context"
	"errors"
	"go-api-starter/internal/database"
	"go-api-starter/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db database.Database
}

func NewUserRepository(db database.Database) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	result := r.db.GetDB().WithContext(ctx).Create(user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return domain.ErrEmailAlreadyExists
		}
		return domain.ErrDatabaseConnection
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	user := &domain.User{}
	result := r.db.GetDB().WithContext(ctx).First(user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, domain.ErrDatabaseConnection
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	result := r.db.GetDB().WithContext(ctx).Where("email = ?", email).First(user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, domain.ErrDatabaseConnection
	}

	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	result := r.db.GetDB().WithContext(ctx).Find(&users)

	if result.Error != nil {
		return nil, domain.ErrDatabaseConnection
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	result := r.db.GetDB().WithContext(ctx).Save(user)

	if result.Error != nil {
		return domain.ErrDatabaseConnection
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.GetDB().WithContext(ctx).Delete(&domain.User{}, id)

	if result.Error != nil {
		return domain.ErrDatabaseConnection
	}

	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
