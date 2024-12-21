package main

import (
	"context"

	"github.com/YXRRXY/todo-app/config"
	"github.com/YXRRXY/todo-app/controller"
	"github.com/YXRRXY/todo-app/repository"
	"github.com/YXRRXY/todo-app/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// CORS 中间件
func corsMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// 处理预检请求
		if string(ctx.Method()) == "OPTIONS" {
			ctx.AbortWithStatus(consts.StatusNoContent)
			return
		}

		ctx.Next(c)
	}
}

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

	// 添加 CORS 中间件
	h.Use(corsMiddleware())

	h.POST("/user/register", userController.Register)
	h.POST("/user/login", userController.Login)

	h.POST("/todo/add", todoController.AddTodo)
	h.GET("/todo/list", todoController.GetTodos)
	h.GET("/todo/search", todoController.SearchTodos)
	h.PUT("/todo/:id/status/:status", todoController.UpdateTodoStatus)

	h.Spin()

}
