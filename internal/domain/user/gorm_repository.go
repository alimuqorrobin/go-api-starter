package user

import (
	"gorm.io/gorm"
)

type gormRepo struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepo{db: db}
}

func (r *gormRepo) FindByID(id uint) (*User, error) {
	var u User
	err := r.db.First(&u, id).Error
	return &u, err
}

func (r *gormRepo) Create(user *User) error {
	return r.db.Create(user).Error
}
