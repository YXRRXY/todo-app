package controller

import (
	"context"
	"strconv"

	"github.com/YXRRXY/todo-app/service"
	"github.com/cloudwego/hertz/pkg/app"
)

type TodoController struct {
	TodoService *service.TodoService
}

func (tc TodoController) AddTodo(ctx context.Context, c *app.RequestContext) {
	userID := c.GetUint("user_id")

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

	todo, err := tc.TodoService.AddTodo(userID, todoRequest.Title, todoRequest.Content, todoRequest.StartTime, todoRequest.EndTime)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(200, todo)
}

func (tc *TodoController) GetTodos(ctx context.Context, c *app.RequestContext) {
	userID := c.GetUint("user_id")
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

	var status *int
	if statusStr, ok := c.GetQuery("status"); ok {
		if statusInt, err := strconv.Atoi(statusStr); err == nil {
			status = &statusInt
		}
	}

	todos, total, err := tc.TodoService.GetTodos(userID, page, pageSize, status)
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
	userID := c.GetUint("user_id")

	var status *int
	if statusStr, ok := c.GetQuery("status"); ok {
		if statusInt, err := strconv.Atoi(statusStr); err == nil {
			status = &statusInt
		}
	}

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

	todos, total, err := tc.TodoService.SearchTodos(userID, keyword, page, pageSize, status)
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
	userID := c.GetUint("user_id")
	todoIDStr := c.Param("id")
	statusStr := c.Param("status")

	todoID, err := strconv.ParseUint(todoIDStr, 10, 64)
	if err != nil {
		c.JSON(400, map[string]string{"error": "无效的代办事项ID"})
		return
	}
	status, err := strconv.ParseInt(statusStr, 10, 64)
	if err != nil {
		c.JSON(400, map[string]string{"error": "无效的状态值"})
		return
	}
	err = tc.TodoService.UpdateTodoStatus(userID, uint(todoID), int(status))
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]string{"msg": "更新成功"})
}

func (tc *TodoController) BatchUpdateStatus(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Status        int    `json:"status"`
		CurrentStatus *int   `json:"current_status,omitempty"`
		IDs           []uint `json:"ids,omitempty"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, map[string]string{"error": "无效的请求参数"})
		return
	}

	userID := c.GetUint("user_id")

	count, err := tc.TodoService.BatchUpdateStatus(userID, req.Status, req.CurrentStatus, req.IDs)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]interface{}{
		"status": 200,
		"msg":    "更新成功",
		"data": map[string]int64{
			"updated_count": count,
		},
	})
}

func (tc *TodoController) BatchDelete(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Status *int   `json:"status,omitempty"`
		IDs    []uint `json:"ids,omitempty"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, map[string]string{"error": "无效的请求参数"})
		return
	}

	userID := c.GetUint("user_id")

	count, err := tc.TodoService.BatchDelete(userID, req.Status, req.IDs)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]interface{}{
		"status": 200,
		"msg":    "删除成功",
		"data": map[string]int64{
			"deleted_count": count,
		},
	})
}

func (tc *TodoController) DeleteTodo(ctx context.Context, c *app.RequestContext) {
	userID := c.GetUint("user_id")
	todoIDStr := c.Param("id")

	todoID, err := strconv.ParseUint(todoIDStr, 10, 64)
	if err != nil {
		c.JSON(400, map[string]string{"error": "无效的待办事项ID"})
		return
	}

	err = tc.TodoService.DeleteTodo(userID, uint(todoID))
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]string{"msg": "删除成功"})
}
