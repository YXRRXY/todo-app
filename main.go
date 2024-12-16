package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yxrxy/todo-app/config"
	"github.com/yxrxy/todo-app/controller"
	"github.com/yxrxy/todo-app/repository"
	"github.com/yxrxy/todo-app/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:zth20041017@tcp(localhost:3306)/todo-app?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败")
	}
	userRepo := &repository.UserRepo{DB: db}
	todoRepo := &repository.TodoRepo{DB: db}

	userService := &service.UserService{Repo: userRepo, JwtSecret: config.GlobalConfig.JwtSecret}
	todoService := &service.TodoService{Repo: todoRepo}

	userController := &controller.UserController{UserService: userService}
	todoController := &controller.TodoController{TodoService: todoService}

	h := server.Default()
	h.POST("/user/register", userController.Register)
	h.POST("/user/login", userController.Login)

	h.POST("/todo/add", todoController.AddTodo)
	h.GET("/todo/list", todoController.GetTodos)
	h.GET("/todo/search", todoController.SearchTodos)
	h.PUT("/todo/:id/status/:status", todoController.UpdateTodoStatus)

	h.Spin()

}
