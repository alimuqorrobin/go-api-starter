package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "go-api-starter/internal/pkg/response"
    deps "go-api-starter/internal/deps"
)

type UserHandler struct { 
    deps *deps.Dependencies 
}

func NewUserHandler(d *deps.Dependencies) *UserHandler {
    return &UserHandler{deps: d}
}

func (h *UserHandler) Me(c *gin.Context) {
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

    u, err := h.deps.UserService.FindByID(uid)
    if err != nil {
        response.Error(c, http.StatusNotFound, "user not found")
        return
    }

    response.Success(c, u)
}
