package response

import (
    "github.com/gin-gonic/gin"
)

// Response adalah struktur response generic untuk semua endpoint
type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   interface{} `json:"error,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

// Meta untuk pagination
type Meta struct {
    Page       int   `json:"page,omitempty"`
    Limit      int   `json:"limit,omitempty"`
    TotalRows  int64 `json:"total_rows,omitempty"`
    TotalPages int   `json:"total_pages,omitempty"`
}

// Success mengembalikan response sukses
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
    c.JSON(statusCode, Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

// SuccessWithMeta mengembalikan response sukses dengan metadata pagination
func SuccessWithMeta(c *gin.Context, statusCode int, message string, data interface{}, meta *Meta) {
    c.JSON(statusCode, Response{
        Success: true,
        Message: message,
        Data:    data,
        Meta:    meta,
    })
}

// Error mengembalikan response error
func Error(c *gin.Context, statusCode int, message string, err interface{}) {
    c.JSON(statusCode, Response{
        Success: false,
        Message: message,
        Error:   formatError(err),
    })
}

// ValidationError mengembalikan response error validasi
func ValidationError(c *gin.Context, errors interface{}) {
    c.JSON(400, Response{
        Success: false,
        Message: "Validation failed",
        Error:   errors,
    })
}

// Unauthorized mengembalikan response unauthorized
func Unauthorized(c *gin.Context, message string) {
    c.JSON(401, Response{
        Success: false,
        Message: message,
    })
}

// Forbidden mengembalikan response forbidden
func Forbidden(c *gin.Context, message string) {
    c.JSON(403, Response{
        Success: false,
        Message: message,
    })
}

// NotFound mengembalikan response not found
func NotFound(c *gin.Context, message string) {
    c.JSON(404, Response{
        Success: false,
        Message: message,
    })
}

// InternalServerError mengembalikan response internal server error
func InternalServerError(c *gin.Context, message string, err interface{}) {
    c.JSON(500, Response{
        Success: false,
        Message: message,
        Error:   formatError(err),
    })
}

// formatError formats error untuk response
func formatError(err interface{}) interface{} {
    if err == nil {
        return nil
    }
    
    switch e := err.(type) {
    case error:
        return e.Error()
    case string:
        return e
    default:
        return err
    }
}

// Success
// response.Success(c, 200, "User created", user)

// // Error
// response.Error(c, 400, "Invalid input", err)

// // Pagination
// response.SuccessWithMeta(c, 200, "Users retrieved", users, meta)