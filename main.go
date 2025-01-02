package main

import (
	"context"
	"strconv"
	"strings"

	"github.com/YXRRXY/todo-app/config"
	"github.com/YXRRXY/todo-app/controller"
	"github.com/YXRRXY/todo-app/model"
	"github.com/YXRRXY/todo-app/repository"
	"github.com/YXRRXY/todo-app/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func corsMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if string(ctx.Method()) == "OPTIONS" {
			ctx.AbortWithStatus(consts.StatusNoContent)
			return
		}

		ctx.Next(c)
	}
}

func authMiddleware(jwtSecret string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		auth := string(c.GetHeader("Authorization"))
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(401, map[string]string{"error": "未授权访问"})
			return
		}

		tokenString := auth[7:]
		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, map[string]string{"error": "无效的token"})
			return
		}

		claims := token.Claims.(*jwt.StandardClaims)
		userID, _ := strconv.ParseUint(claims.Id, 10, 64)
		c.Set("user_id", uint(userID))
		c.Next(ctx)
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

	h.Use(corsMiddleware())

	auth := h.Group("/", authMiddleware(config.GlobalConfig.JwtSecret))
	{
		auth.POST("/todo/add", todoController.AddTodo)
		auth.GET("/todo/list", todoController.GetTodos)
		auth.GET("/todo/search", todoController.SearchTodos)
		auth.PUT("/todo/:id/status/:status", todoController.UpdateTodoStatus)
		auth.PUT("/todo/status/batch", todoController.BatchUpdateStatus)
		auth.DELETE("/todo/batch", todoController.BatchDelete)
		auth.DELETE("/todo/:id", todoController.DeleteTodo)
	}

	h.POST("/user/register", userController.Register)
	h.POST("/user/login", userController.Login)

	err = db.AutoMigrate(&model.Todo{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	h.Spin()

}
