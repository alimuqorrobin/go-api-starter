package middleware

import (
    "strings"

    "golang-api-starter/pkg/jwt"
    "golang-api-starter/pkg/response"
    "github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token
func AuthMiddleware(tokenService *jwt.TokenService) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            response.Unauthorized(c, "Authorization header is required")
            c.Abort()
            return
        }

        // Check Bearer token format
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            response.Unauthorized(c, "Invalid authorization header format")
            c.Abort()
            return
        }

        tokenString := parts[1]

        // Validate token
        claims, err := tokenService.ValidateToken(tokenString)
        if err != nil {
            response.Unauthorized(c, "Invalid or expired token")
            c.Abort()
            return
        }

        // Set user info in context
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("email", claims.Email)

        c.Next()
    }
}

// GetUserID gets user ID from context
func GetUserID(c *gin.Context) (int64, bool) {
    userID, exists := c.Get("user_id")
    if !exists {
        return 0, false
    }
    return userID.(int64), true
}

// GetUsername gets username from context
func GetUsername(c *gin.Context) (string, bool) {
    username, exists := c.Get("username")
    if !exists {
        return "", false
    }
    return username.(string), true
}

// Di handler
// userID, _ := middleware.GetUserID(c)
// username, _ := middleware.GetUsername(c)