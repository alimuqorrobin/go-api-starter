package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Page    int         `json:"page,omitempty"`
	Limit   int         `json:"limit,omitempty"`
	Total   int64       `json:"total,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: true,
		Data:    data,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

func Error(w http.ResponseWriter, statusCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: false,
		Error:   errMsg,
	}

	json.NewEncoder(w).Encode(response)
}

func Created(w http.ResponseWriter, data interface{}) {
	Success(w, http.StatusCreated, data, "Resource created successfully")
}

func OK(w http.ResponseWriter, data interface{}) {
	Success(w, http.StatusOK, data, "Request successful")
}

func BadRequest(w http.ResponseWriter, errMsg string) {
	Error(w, http.StatusBadRequest, errMsg)
}

func Unauthorized(w http.ResponseWriter, errMsg string) {
	Error(w, http.StatusUnauthorized, errMsg)
}

func Forbidden(w http.ResponseWriter, errMsg string) {
	Error(w, http.StatusForbidden, errMsg)
}

func NotFound(w http.ResponseWriter, errMsg string) {
	Error(w, http.StatusNotFound, errMsg)
}

func InternalError(w http.ResponseWriter, errMsg string) {
	Error(w, http.StatusInternalServerError, errMsg)
}

func TooManyRequests(w http.ResponseWriter) {
	Error(w, http.StatusTooManyRequests, "Too many requests, please try again later")
}

func Paginated(w http.ResponseWriter, data interface{}, page, limit int, total int64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := PaginatedResponse{
		Success: true,
		Data:    data,
		Page:    page,
		Limit:   limit,
		Total:   total,
	}

	json.NewEncoder(w).Encode(response)
}
