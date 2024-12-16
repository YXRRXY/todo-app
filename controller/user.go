package controller

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/yxrxy/todo-app/service"
)

type UserController struct {
	UserService *service.UserService
}

func (uc UserController) Register(ctx context.Context, c *app.RequestContext) {
	var userRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindAndValidate(&userRequest); err != nil {
		c.JSON(400, map[string]string{"error": "Invalid request"})
		return
	}

	user, err := uc.UserService.Register(userRequest.Username, userRequest.Password)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func (uc *UserController) Login(ctx context.Context, c *app.RequestContext) {
	var userRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindAndValidate(&userRequest); err != nil {
		c.JSON(400, map[string]string{"error": "Invalid request"})
		return
	}

	user, token, err := uc.UserService.Login(userRequest.Username, userRequest.Password)
	if err != nil {
		c.JSON(401, map[string]string{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	c.JSON(200, response)
}
