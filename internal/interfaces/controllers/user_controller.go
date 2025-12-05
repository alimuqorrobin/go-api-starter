package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-api-starter/internal/usecase"
	"go-api-starter/internal/interfaces/middleware"
	"go-api-starter/internal/utils"
)

type UserController struct {
	usecase usecase.UserUsecase
}

func NewUserController(u usecase.UserUsecase) *UserController {
	return &UserController{
		usecase: u,
	}
}

func (ctl *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := ctl.usecase.GetUserByID(c, uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, "User found", user)
}

func (ctl *UserController) CreateUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	user, err := ctl.usecase.CreateUser(c, req.Name, req.Email)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Created(c, "User created", user)
}
