package handler

import (
    "net/http"

    "go-api-starter/internal/domain"
    "go-api-starter/internal/service"
    "go-api-starter/pkg/response"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{
        authService: authService,
    }
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Register request"
// @Success 201 {object} response.Response{data=domain.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    var req domain.RegisterRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err.Error())
        return
    }

    user, err := h.authService.Register(c.Request.Context(), &req)
    if err != nil {
        if err == domain.ErrUserAlreadyExists {
            response.Error(c, http.StatusConflict, "User already exists", err.Error())
            return
        }
        response.InternalServerError(c, "Failed to register user", err.Error())
        return
    }

    response.Success(c, http.StatusCreated, "User registered successfully", user)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login request"
// @Success 200 {object} response.Response{data=domain.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    var req domain.LoginRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err.Error())
        return
    }

    loginResponse, err := h.authService.Login(c.Request.Context(), &req)
    if err != nil {
        if err == domain.ErrInvalidCredentials {
            response.Unauthorized(c, "Invalid email or password")
            return
        }
        response.InternalServerError(c, "Login failed", err.Error())
        return
    }

    response.Success(c, http.StatusOK, "Login successful", loginResponse)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} response.Response{data=domain.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
    var req domain.RefreshTokenRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err.Error())
        return
    }

    loginResponse, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
    if err != nil {
        response.Unauthorized(c, "Invalid or expired refresh token")
        return
    }

    response.Success(c, http.StatusOK, "Token refreshed successfully", loginResponse)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user (client should delete tokens)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
    // In a stateless JWT system, logout is handled client-side by deleting tokens
    // Here we just return success
    response.Success(c, http.StatusOK, "Logged out successfully", nil)
}
```

**Penjelasan:**
- HTTP handlers untuk authentication
- Swagger annotations (untuk auto-generate docs)
- Input validation dengan `ShouldBindJSON`
- Error handling yang proper
- Generic response format

**Handler Pattern:**
```
1. Bind & validate input
2. Call service layer
3. Handle errors
4. Return response