package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/yourname/go-backend-enterprise/internal/server"
    "github.com/yourname/go-backend-enterprise/internal/pkg/jwt"
    "github.com/yourname/go-backend-enterprise/internal/pkg/response"
)

type loginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func LoginHandler(d *server.Dependencies) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req loginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            response.Error(c, http.StatusBadRequest, "invalid payload")
            return
        }

        u, err := d.UserService.Authenticate(req.Username, req.Password)
        if err != nil {
            response.Error(c, http.StatusUnauthorized, "invalid credentials")
            return
        }

        token, err := jwt.GenerateToken(map[string]interface{}{
            "sub":  u.ID,
            "name": u.Username,
            "exp":  time.Now().Add(time.Hour * time.Duration(d.Config.JwtTtlHours)).Unix(),
        }, d.Config.JwtSecret)
        if err != nil {
            response.Error(c, http.StatusInternalServerError, "failed to generate token")
            return
        }

        response.Success(c, gin.H{"token": token})
    }
}

type registerRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func RegisterHandler(d *server.Dependencies) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req registerRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            response.Error(c, http.StatusBadRequest, "invalid payload")
            return
        }
        id, err := d.UserService.CreateUser(req.Username, req.Password)
        if err != nil {
            response.Error(c, http.StatusInternalServerError, err.Error())
            return
        }
        response.Success(c, gin.H{"id": id})
    }
}

func MeHandler(d *server.Dependencies) gin.HandlerFunc {
    return func(c *gin.Context) {
        sub, exists := c.Get("sub")
        if !exists {
            response.Error(c, http.StatusUnauthorized, "missing subject")
            return
        }
        var uid int64
        switch v := sub.(type) {
        case float64:
            uid = int64(v)
        case int64:
            uid = v
        case int:
            uid = int64(v)
        default:
            response.Error(c, http.StatusUnauthorized, "invalid subject")
            return
        }

        u, err := d.UserService.FindByID(uid)
        if err != nil {
            response.Error(c, http.StatusNotFound, "user not found")
            return
        }
        response.Success(c, u)
    }
}
