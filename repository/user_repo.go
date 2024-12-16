package repository

import (
	"github.com/yxrxy/todo-app/model"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func (repo *UserRepo) SaveUser(user *model.User) error {
	return repo.DB.Create(user).Error
}

func (repo *UserRepo) FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := repo.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
