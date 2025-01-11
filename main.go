package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

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

func initDatabase() (*gorm.DB, error) {
	dbRootDSN := fmt.Sprintf("%s:%s@tcp(%s:3306)/",
		config.GlobalConfig.DBUser,
		config.GlobalConfig.DBPassword,
		config.GlobalConfig.DBHost)

	db, err := sql.Open("mysql", dbRootDSN)
	if err != nil {
		return nil, fmt.Errorf("连接MySQL失败: %v", err)
	}
	defer db.Close()

	// 创建数据库
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", config.GlobalConfig.DBName))
	if err != nil {
		return nil, fmt.Errorf("创建数据库失败: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GlobalConfig.DBUser,
		config.GlobalConfig.DBPassword,
		config.GlobalConfig.DBHost,
		config.GlobalConfig.DBName)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 自动迁移表结构
	err = gormDB.AutoMigrate(
		&model.User{},
		&model.Todo{},
	)
	if err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %v", err)
	}

	return gormDB, nil
}

func main() {
	db, err := initDatabase()
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

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

	h.Spin()

}
