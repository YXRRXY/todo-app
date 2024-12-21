package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/YXRRXY/todo-app/model"
	"github.com/YXRRXY/todo-app/repository"
	"github.com/golang-jwt/jwt"
)

type UserService struct {
	Repo      *repository.UserRepo
	JwtSecret string
}

func (service *UserService) Register(username, password string) (*model.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("用户或密码不能为空")
	}
	if existingUser, _ := service.Repo.FindUserByUsername(username); existingUser != nil {
		return nil, errors.New("用户已存在")
	}
	user := &model.User{
		Username: username,
		Password: password,
	}
	err := service.Repo.SaveUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) Login(username, password string) (*model.User, string, error) {
	user, err := service.Repo.FindUserByUsername(username)
	if err != nil {
		return nil, "", err
	}
	if user.Password != password {
		return nil, "", errors.New("密码错误")
	}
	token, err := service.generateToken(user.Id)
	if err != nil {
		return nil, "", err
	}
	user.Token = token
	return user, token, nil
}

func (service *UserService) generateToken(userID uint) (string, error) {
	claims := jwt.StandardClaims{
		Id:        strconv.FormatUint(uint64(userID), 10),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(service.JwtSecret))
	return signedToken, err
}
