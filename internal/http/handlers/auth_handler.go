package handlers

import (
	"encoding/json"
	"net/http"

	"go-api-starter/internal/domain"
	"go-api-starter/internal/http/response"
	"go-api-starter/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	tokenResp, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		if err == domain.ErrUserNotFound || err == domain.ErrInvalidPassword {
			response.Unauthorized(w, "Invalid email or password")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, tokenResp)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req domain.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	tokenResp, err := h.authService.RefreshToken(r.Context(), &req)
	if err != nil {
		response.Unauthorized(w, "Invalid refresh token")
		return
	}

	response.OK(w, tokenResp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		response.Unauthorized(w, "User not found in context")
		return
	}

	if err := h.authService.Logout(r.Context(), userID); err != nil {
		response.InternalError(w, "Logout failed")
		return
	}

	response.OK(w, map[string]string{"message": "Logged out successfully"})
}
