package user

import "go-api-starter/internal/domain/common"

type User struct {
	common.BaseEntity

	Name  string `json:"name"`
	Email string `json:"email"`
}
