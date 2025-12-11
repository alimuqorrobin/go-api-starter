package jwt

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var (
    ErrInvalidToken = errors.New("invalid token")
    ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}

type TokenService struct {
    secretKey         string
    expirationHours   int
    refreshExpiration int
}

func NewTokenService(secretKey string, expirationHours, refreshExpiration int) *TokenService {
    return &TokenService{
        secretKey:         secretKey,
        expirationHours:   expirationHours,
        refreshExpiration: refreshExpiration,
    }
}

// GenerateToken generates a new JWT token
func (s *TokenService) GenerateToken(userID int64, username, email string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        Email:    email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.expirationHours))),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.secretKey))
}

// GenerateRefreshToken generates a new refresh token
func (s *TokenService) GenerateRefreshToken(userID int64, username, email string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        Email:    email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.refreshExpiration))),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.secretKey))
}

// ValidateToken validates and parses JWT token
func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidToken
        }
        return []byte(s.secretKey), nil
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, ErrInvalidToken
    }

    // Check expiration
    if claims.ExpiresAt.Time.Before(time.Now()) {
        return nil, ErrExpiredToken
    }

    return claims, nil
}
```

**Penjelasan:**
- JWT (JSON Web Token) untuk authentication
- **Access Token** = Short-lived (24 hours)
- **Refresh Token** = Long-lived (7 days)
- HMAC SHA256 signing algorithm

**JWT Structure:**
```
Header.Payload.Signature

Header:
{
  "alg": "HS256",
  "typ": "JWT"
}

Payload (Claims):
{
  "user_id": 1,
  "username": "john",
  "email": "john@example.com",
  "exp": 1735689600,  // Expiration
  "iat": 1735603200,  // Issued At
  "nbf": 1735603200   // Not Before
}

Signature:
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  secret_key
)
```

**Token Flow:**
```
// Login
//   ↓
// Generate Access Token (24h)
// Generate Refresh Token (7d)
//   ↓
// Return both tokens to client

// Client stores tokens
//   ↓
// Every request: Send Access Token in header
// Authorization: Bearer <access_token>
//   ↓
// Access Token expired?
//   ↓
// Use Refresh Token to get new Access Token