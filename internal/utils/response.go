package utils

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Respond(c *gin.Context, httpCode int, status string, message string, data interface{}) {
	c.JSON(httpCode, APIResponse{
		Status:  status,
		Code:    httpCode,
		Message: message,
		Data:    data,
	})
}

func Success(c *gin.Context, message string, data interface{}) {
	Respond(c, 200, "success", message, data)
}

func Created(c *gin.Context, message string, data interface{}) {
	Respond(c, 201, "success", message, data)
}

func Error(c *gin.Context, httpCode int, message string) {
	Respond(c, httpCode, "error", message, nil)
}
