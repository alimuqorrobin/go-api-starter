package service

import (
    "context"

    "golang-api-starter/internal/domain"
    "golang-api-starter/pkg/jwt"
)

type AuthService struct {
    userService  *UserService
    tokenService *jwt.TokenService
}

func NewAuthService(userService *UserService, tokenService *jwt.TokenService) *AuthService {
    return &AuthService{
        userService:  userService,
        tokenService: tokenService,
    }
}

// Login authenticates user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
    // Find user by email
    user, err := s.userService.GetUserByEmail(ctx, req.Email)
    if err != nil {
        return nil, domain.ErrInvalidCredentials
    }

    // Verify password
    if err := s.userService.VerifyPassword(user.PasswordHash, req.Password); err != nil {
        return nil, domain.ErrInvalidCredentials
    }

    // Check if user is active
    if !user.IsActive {
        return nil, domain.ErrUnauthorized
    }

    // Generate tokens
    accessToken, err := s.tokenService.GenerateToken(user.ID, user.Username, user.Email)
    if err != nil {
        return nil, err
    }

    refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, user.Username, user.Email)
    if err != nil {
        return nil, err
    }

    return &domain.LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
        ExpiresIn:    86400, // 24 hours
        User:         user.ToResponse(),
    }, nil
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.UserResponse, error) {
    createReq := &domain.CreateUserRequest{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
        FullName: req.FullName,
    }

    return s.userService.CreateUser(ctx, createReq)
}

// RefreshToken generates new access token from refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.LoginResponse, error) {
    // Validate refresh token
    claims, err := s.tokenService.ValidateToken(refreshToken)
    if err != nil {
        return nil, domain.ErrInvalidToken
    }

    // Get user
    user, err := s.userService.GetUserByEmail(ctx, claims.Email)
    if err != nil {
        return nil, domain.ErrUserNotFound
    }

    // Generate new tokens
    newAccessToken, err := s.tokenService.GenerateToken(user.ID, user.Username, user.Email)
    if err != nil {
        return nil, err
    }

    newRefreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, user.Username, user.Email)
    if err != nil {
        return nil, err
    }

    return &domain.LoginResponse{
        AccessToken:  newAccessToken,
        RefreshToken: newRefreshToken,
        TokenType:    "Bearer",
        ExpiresIn:    86400,
        User:         user.ToResponse(),
    }, nil
}