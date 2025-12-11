package domain

import (
    "time"
)

// User entity
type User struct {
    ID           int64      `json:"id"`
    Username     string     `json:"username"`
    Email        string     `json:"email"`
    PasswordHash string     `json:"-"` // Never expose password
    FullName     string     `json:"full_name"`
    IsActive     bool       `json:"is_active"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}

// CreateUserRequest DTO
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    FullName string `json:"full_name" binding:"required,min=3"`
}

// UpdateUserRequest DTO
type UpdateUserRequest struct {
    Username string `json:"username" binding:"omitempty,min=3,max=50"`
    Email    string `json:"email" binding:"omitempty,email"`
    FullName string `json:"full_name" binding:"omitempty,min=3"`
    IsActive *bool  `json:"is_active"`
}

// UserResponse DTO
type UserResponse struct {
    ID        int64     `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    FullName  string    `json:"full_name"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
    return &UserResponse{
        ID:        u.ID,
        Username:  u.Username,
        Email:     u.Email,
        FullName:  u.FullName,
        IsActive:  u.IsActive,
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }
}