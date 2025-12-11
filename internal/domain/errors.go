package domain

import "errors"

var (
    // User errors
    ErrUserNotFound      = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrInvalidCredentials = errors.New("invalid credentials")
    
    // Auth errors
    ErrInvalidToken      = errors.New("invalid token")
    ErrExpiredToken      = errors.New("token has expired")
    ErrUnauthorized      = errors.New("unauthorized")
    
    // Validation errors
    ErrInvalidInput      = errors.New("invalid input")
    ErrValidationFailed  = errors.New("validation failed")
    
    // Database errors
    ErrDatabaseConnection = errors.New("database connection error")
    ErrDatabaseQuery      = errors.New("database query error")
)