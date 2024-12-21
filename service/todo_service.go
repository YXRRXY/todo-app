package service

import (
	"errors"
	"time"

	"github.com/YXRRXY/todo-app/model"
	"github.com/YXRRXY/todo-app/repository"
)

type TodoService struct {
	Repo *repository.TodoRepo
}

func (s *TodoService) AddTodo(content, title string, startTime, endTime int64) (*model.Todo, error) {
	todo := &model.Todo{
		Title:     title,
		Content:   content,
		Status:    0,
		CreatedAt: time.Now(),
		StartTime: time.Unix(startTime, 0),
		EndTime:   time.Unix(endTime, 0),
	}

	err := s.Repo.CreateTodo(todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *TodoService) GetTodos(page, pageSize int) ([]model.Todo, int64, error) {
	return s.Repo.GetTodos(page, pageSize)
}

func (s *TodoService) UpdateTodoStatus(todoID uint, status int) error {
	if status != 0 && status != 1 {
		return errors.New("无效的状态")
	}
	return s.Repo.UpdateTodoStatus(todoID, status)
}

func (s *TodoService) SearchTodos(keyword string, page, pageSize int) ([]model.Todo, int64, error) {
	return s.Repo.SearchTodos(keyword, page, pageSize)
}
