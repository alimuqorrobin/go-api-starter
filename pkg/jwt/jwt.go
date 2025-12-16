package jwt

import (
	"go-api-starter/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secret             string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func NewManager() *Manager {
	return &Manager{
		secret:             "your-secret-key-change-in-production",
		accessTokenExpiry:  24 * time.Hour,
		refreshTokenExpiry: 7 * 24 * time.Hour,
	}
}

func (m *Manager) GenerateAccessToken(userID uint) (string, time.Time) {
	expiryTime := time.Now().Add(m.accessTokenExpiry)
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(m.secret))

	return tokenString, expiryTime
}

func (m *Manager) GenerateRefreshToken(userID uint) (string, time.Time) {
	expiryTime := time.Now().Add(m.refreshTokenExpiry)
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(m.secret))

	return tokenString, expiryTime
}

func (m *Manager) ValidateToken(tokenString string) (uint, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secret), nil
	})

	if err != nil {
		return 0, domain.ErrInvalidToken
	}

	if !token.Valid {
		return 0, domain.ErrTokenExpired
	}

	return claims.UserID, nil
}
