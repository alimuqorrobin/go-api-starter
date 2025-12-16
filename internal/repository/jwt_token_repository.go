package repository

import (
	"context"
	"errors"
	"go-api-starter/internal/database"
	"go-api-starter/internal/domain"
	"gorm.io/gorm"
	"time"
)

type jwtTokenRepository struct {
	db database.Database
}

func NewJWTTokenRepository(db database.Database) JWTTokenRepository {
	return &jwtTokenRepository{db: db}
}

func (r *jwtTokenRepository) Save(ctx context.Context, token *domain.JWTToken) error {
	result := r.db.GetDB().WithContext(ctx).Create(token)

	if result.Error != nil {
		return domain.ErrDatabaseConnection
	}

	return nil
}

func (r *jwtTokenRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.JWTToken, error) {
	token := &domain.JWTToken{}
	result := r.db.GetDB().WithContext(ctx).
		Where("refresh_token = ? AND expires_at > ?", refreshToken, time.Now()).
		First(token)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTokenExpired
		}
		return nil, domain.ErrDatabaseConnection
	}

	return token, nil
}

func (r *jwtTokenRepository) GetByAccessToken(ctx context.Context, accessToken string) (*domain.JWTToken, error) {
	token := &domain.JWTToken{}
	result := r.db.GetDB().WithContext(ctx).
		Where("access_token = ? AND expires_at > ?", accessToken, time.Now()).
		First(token)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrInvalidToken
		}
		return nil, domain.ErrDatabaseConnection
	}

	return token, nil
}

func (r *jwtTokenRepository) DeleteByRefreshToken(ctx context.Context, refreshToken string) error {
	result := r.db.GetDB().WithContext(ctx).
		Where("refresh_token = ?", refreshToken).
		Delete(&domain.JWTToken{})

	if result.Error != nil {
		return domain.ErrDatabaseConnection
	}

	return nil
}

func (r *jwtTokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	result := r.db.GetDB().WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&domain.JWTToken{})

	if result.Error != nil {
		return domain.ErrDatabaseConnection
	}

	return nil
}

func (r *jwtTokenRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.JWTToken, error) {
	var tokens []domain.JWTToken
	result := r.db.GetDB().WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&tokens)

	if result.Error != nil {
		return nil, domain.ErrDatabaseConnection
	}

	return tokens, nil
}
