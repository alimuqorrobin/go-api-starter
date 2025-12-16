package service

import (
	"context"
	"go-api-starter/internal/domain"
	"go-api-starter/internal/repository"
	"go-api-starter/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.JWTTokenRepository
	jwtMgr    *jwt.Manager
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.JWTTokenRepository,
	jwtMgr *jwt.Manager,
) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtMgr:    jwtMgr,
	}
}

func (s *authService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, domain.ErrInvalidPassword
	}

	accessToken, accessExpires := s.jwtMgr.GenerateAccessToken(user.ID)
	refreshToken, refreshExpires := s.jwtMgr.GenerateRefreshToken(user.ID)

	tokenRecord := &domain.JWTToken{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshExpires,
	}

	if err := s.tokenRepo.Save(ctx, tokenRecord); err != nil {
		return nil, err
	}

	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    accessExpires.Unix(),
		TokenType:    "Bearer",
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.TokenResponse, error) {
	tokenRecord, err := s.tokenRepo.GetByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, tokenRecord.UserID)
	if err != nil {
		return nil, err
	}

	newAccessToken, newAccessExpires := s.jwtMgr.GenerateAccessToken(user.ID)
	newRefreshToken, newRefreshExpires := s.jwtMgr.GenerateRefreshToken(user.ID)

	s.tokenRepo.DeleteByRefreshToken(ctx, req.RefreshToken)

	newTokenRecord := &domain.JWTToken{
		UserID:       user.ID,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    newRefreshExpires,
	}

	if err := s.tokenRepo.Save(ctx, newTokenRecord); err != nil {
		return nil, err
	}

	return &domain.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    newAccessExpires.Unix(),
		TokenType:    "Bearer",
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID uint) error {
	tokens, err := s.tokenRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, token := range tokens {
		s.tokenRepo.DeleteByRefreshToken(ctx, token.RefreshToken)
	}

	return nil
}

func (s *authService) ValidateAccessToken(ctx context.Context, token string) (uint, error) {
	userID, err := s.jwtMgr.ValidateToken(token)
	if err != nil {
		return 0, err
	}

	_, err = s.tokenRepo.GetByAccessToken(ctx, token)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
