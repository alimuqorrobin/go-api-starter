package domain

// LoginRequest DTO
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// LoginResponse DTO
type LoginResponse struct {
    AccessToken  string        `json:"access_token"`
    RefreshToken string        `json:"refresh_token"`
    TokenType    string        `json:"token_type"`
    ExpiresIn    int           `json:"expires_in"`
    User         *UserResponse `json:"user"`
}

// RefreshTokenRequest DTO
type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

// RegisterRequest DTO
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    FullName string `json:"full_name" binding:"required,min=3"`
}