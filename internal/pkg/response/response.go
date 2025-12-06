package response

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type ApiResponse struct {
    Status  string      `json:"status"`
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    TraceID string      `json:"trace_id"`
}

func Success(c *gin.Context, data interface{}) {
    traceID := uuid.New().String()
    c.JSON(http.StatusOK, ApiResponse{
        Status:  "success",
        Code:    http.StatusOK,
        Message: "OK",
        Data:    data,
        TraceID: traceID,
    })
}

func SuccessWithCode(c *gin.Context, code int, data interface{}) {
    traceID := uuid.New().String()
    c.JSON(code, ApiResponse{
        Status:  "success",
        Code:    code,
        Message: http.StatusText(code),
        Data:    data,
        TraceID: traceID,
    })
}

func Error(c *gin.Context, code int, message string) {
    traceID := uuid.New().String()
    c.JSON(code, ApiResponse{
        Status:  "error",
        Code:    code,
        Message: message,
        Data:    nil,
        TraceID: traceID,
    })
}
