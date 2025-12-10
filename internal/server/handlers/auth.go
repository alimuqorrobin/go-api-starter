package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go-api-starter/internal/pkg/response"
    "go-api-starter/internal/pkg/jwt"
    deps "go-api-starter/internal/deps"
)

type AuthHandler struct {
    deps *deps.Dependencies
}

func NewAuthHandler(d *deps.Dependencies) *AuthHandler {
    return &AuthHandler{deps: d}
}

type loginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req loginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "invalid payload")
        return
    }

    u, err := h.deps.UserService.Authenticate(req.Username, req.Password)
    if err != nil {
        response.Error(c, http.StatusUnauthorized, "invalid credentials")
        return
    }

    token, err := jwt.GenerateToken(map[string]interface{}{
        "sub":  u.ID,
        "name": u.Username,
        "exp":  time.Now().Add(time.Hour * time.Duration(h.deps.Config.JwtTtlHours)).Unix(),
    }, h.deps.Config.JwtSecret)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "failed to generate token")
        return
    }

    response.Success(c, gin.H{"token": token})
}

func NewRegisterHandler(d *deps.Dependencies) *AuthHandler {
    return &AuthHandler{deps: d}
}

type registerRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req registerRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "invalid payload")
        return
    }

    id, err := h.deps.UserService.CreateUser(req.Username, req.Password)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(c, gin.H{"id": id})
}
