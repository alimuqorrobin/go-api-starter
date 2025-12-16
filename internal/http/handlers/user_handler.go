package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-api-starter/internal/domain"
	"go-api-starter/internal/http/response"
	"go-api-starter/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	user, err := h.userService.CreateUser(r.Context(), &req)
	if err != nil {
		if err == domain.ErrEmailAlreadyExists {
			response.BadRequest(w, "Email already exists")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.Created(w, user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUser(r.Context(), uint(id))
	if err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(w, "User not found")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, user)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers(r.Context())
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, users)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}

	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), uint(id), &req)
	if err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(w, "User not found")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(r.Context(), uint(id)); err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(w, "User not found")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "User deleted successfully"})
}
