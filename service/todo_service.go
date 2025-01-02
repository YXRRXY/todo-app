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

func (s *TodoService) AddTodo(userID uint, title, content string, startTime, endTime int64) (*model.Todo, error) {
	todo := &model.Todo{
		UserID:    userID,
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

func (s *TodoService) GetTodos(userID uint, page, pageSize int, status *int) ([]model.Todo, int64, error) {
	return s.Repo.GetTodos(userID, page, pageSize, status)
}

func (s *TodoService) UpdateTodoStatus(userID uint, todoID uint, status int) error {
	if status != 0 && status != 1 {
		return errors.New("无效的状态值")
	}
	return s.Repo.UpdateTodoStatus(userID, todoID, status)
}

func (s *TodoService) SearchTodos(userID uint, keyword string, page, pageSize int, status *int) ([]model.Todo, int64, error) {
	return s.Repo.SearchTodos(userID, keyword, page, pageSize, status)
}

func (s *TodoService) BatchUpdateStatus(userID uint, status int, currentStatus *int, ids []uint) (int64, error) {
	if status != 0 && status != 1 {
		return 0, errors.New("无效的状态值")
	}
	return s.Repo.BatchUpdateStatus(userID, status, currentStatus, ids)
}

func (s *TodoService) DeleteTodo(userID uint, todoID uint) error {
	return s.Repo.DeleteTodo(userID, todoID)
}

func (s *TodoService) BatchDelete(userID uint, status *int, ids []uint) (int64, error) {
	return s.Repo.BatchDelete(userID, status, ids)
}
