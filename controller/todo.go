package controller

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/yxrxy/todo-app/service"
)

type TodoController struct {
	TodoService *service.TodoService
}

func (tc TodoController) AddTodo(ctx context.Context, c *app.RequestContext) {
	var todoRequest struct {
		Title     string `json:"title"`
		Content   string `json:"content"`
		StartTime int64  `json:"start_time"`
		EndTime   int64  `json:"end_time"`
	}
	if err := c.BindAndValidate(&todoRequest); err != nil {
		c.JSON(400, map[string]string{"error": "Invalid request"})
		return
	}

	todo, err := tc.TodoService.AddTodo(todoRequest.Title, todoRequest.Content, todoRequest.StartTime, todoRequest.EndTime)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(200, todo)
}

func (tc *TodoController) GetTodos(ctx context.Context, c *app.RequestContext) {
	page := 1
	pageSize := 10

	if p, ok := c.GetQuery("page"); ok {
		if pInt, err := strconv.ParseInt(p, 10, 64); err == nil {
			page = int(pInt)
		}
	}
	if ps, ok := c.GetQuery("page_size"); ok {
		if psInt, err := strconv.ParseInt(ps, 10, 64); err == nil {
			pageSize = int(psInt)
		}
	}

	todos, total, err := tc.TodoService.GetTodos(page, pageSize)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"status": 200,
		"data": map[string]interface{}{
			"items": todos,
			"total": total,
		},
		"msg":   "ok",
		"error": "",
	}

	c.JSON(200, response)
}

func (tc *TodoController) SearchTodos(ctx context.Context, c *app.RequestContext) {
	keyword := c.DefaultQuery("keyword", "")
	page := 1
	pageSize := 10

	if p, ok := c.GetQuery("page"); ok {
		if pInt, err := strconv.ParseInt(p, 10, 64); err == nil {
			page = int(pInt)
		}
	}
	if ps, ok := c.GetQuery("page_size"); ok {
		if psInt, err := strconv.ParseInt(ps, 10, 64); err == nil {
			pageSize = int(psInt)
		}
	}

	todos, total, err := tc.TodoService.SearchTodos(keyword, page, pageSize)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"status": 200,
		"data": map[string]interface{}{
			"items": todos,
			"total": total,
		},
		"msg":   "ok",
		"error": "",
	}

	c.JSON(200, response)
}

func (tc *TodoController) UpdateTodoStatus(ctx context.Context, c *app.RequestContext) {
	todoIDStr := c.Param("id")
	statusStr := c.Param("status")

	todoID, err := strconv.ParseInt(todoIDStr, 10, 64)
	if err != nil {
		c.JSON(400, map[string]string{"error": "无效的代办事项ID"})
		return
	}
	status, err := strconv.ParseInt(statusStr, 10, 64)
	if err != nil {
		c.JSON(400, map[string]string{"error": "无效的状态值"})
		return
	}
	err = tc.TodoService.UpdateTodoStatus(uint(todoID), int(status))
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]string{"msg": "更新成功"})
}
