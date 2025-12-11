package handler

import (
    "net/http"
    "strconv"

    "go-api-starter/internal/domain"
    "go-api-starter/internal/middleware"
    "go-api-starter/internal/service"
    "go-api-starter/pkg/response"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateUserRequest true "Create user request"
// @Success 201 {object} response.Response{data=domain.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req domain.CreateUserRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err.Error())
        return
    }

    user, err := h.userService.CreateUser(c.Request.Context(), &req)
    if err != nil {
        if err == domain.ErrUserAlreadyExists {
            response.Error(c, http.StatusConflict, "User already exists", err.Error())
            return
        }
        response.InternalServerError(c, "Failed to create user", err.Error())
        return
    }

    response.Success(c, http.StatusCreated, "User created successfully", user)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get user details by user ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.Response{data=domain.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid user ID", err.Error())
        return
    }

    user, err := h.userService.GetUserByID(c.Request.Context(), id)
    if err != nil {
        if err == domain.ErrUserNotFound {
            response.NotFound(c, "User not found")
            return
        }
        response.InternalServerError(c, "Failed to get user", err.Error())
        return
    }

    response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

// ListUsers godoc
// @Summary List all users
// @Description Get list of all users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} response.Response{data=[]domain.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }

    users, total, err := h.userService.ListUsers(c.Request.Context(), page, limit)
    if err != nil {
        response.InternalServerError(c, "Failed to retrieve users", err.Error())
        return
    }

    totalPages := int(total) / limit
    if int(total)%limit != 0 {
        totalPages++
    }

    meta := &response.Meta{
        Page:       page,
        Limit:      limit,
        TotalRows:  total,
        TotalPages: totalPages,
    }

    response.SuccessWithMeta(c, http.StatusOK, "Users retrieved successfully", users, meta)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body domain.UpdateUserRequest true "Update user request"
// @Success 200 {object} response.Response{data=domain.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid user ID", err.Error())
        return
    }

    var req domain.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err.Error())
        return
    }

    user, err := h.userService.UpdateUser(c.Request.Context(), id, &req)
    if err != nil {
        if err == domain.ErrUserNotFound {
            response.NotFound(c, "User not found")
            return
        }
        response.InternalServerError(c, "Failed to update user", err.Error())
        return
    }

    response.Success(c, http.StatusOK, "User updated successfully", user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid user ID", err.Error())
        return
    }

    if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
        if err == domain.ErrUserNotFound {
            response.NotFound(c, "User not found")
            return
        }
        response.InternalServerError(c, "Failed to delete user", err.Error())
        return
    }

    response.Success(c, http.StatusOK, "User deleted successfully", nil)
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Get authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=domain.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
    userID, exists := middleware.GetUserID(c)
    if !exists {
        response.Unauthorized(c, "User not authenticated")
        return
    }

    user, err := h.userService.GetUserByID(c.Request.Context(), userID)
    if err != nil {
        response.InternalServerError(c, "Failed to get profile", err.Error())
        return
    }

    response.Success(c, http.StatusOK, "Profile retrieved successfully", user)
}

// BulkCreateUsers godoc
// @Summary Bulk create users
// @Description Create multiple users concurrently
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body []domain.CreateUserRequest true "Bulk create users request"
// @Success 201 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /users/bulk [post]
func (h *UserHandler) BulkCreateUsers(c *gin.Context) {
    var requests []domain.CreateUserRequest

    if err := c.ShouldBindJSON(&requests); err != nil {
        response.ValidationError(c, err.Error())
        return
    }

    if len(requests) == 0 {
        response.Error(c, http.StatusBadRequest, "No users to create", nil)
        return
    }

    result := h.userService.BulkCreateUsers(c.Request.Context(), requests)

    response.Success(c, http.StatusCreated, "Bulk user creation completed", result)
}

// URL: /users?page=2&limit=20
// page := c.DefaultQuery("page", "1")   // default 1
// limit := c.DefaultQuery("limit", "10") // default 10
// URL: /users/123
// id := c.Param("id") // "123"
// Dari JWT auth middleware
// userID, _ := middleware.GetUserID(c)