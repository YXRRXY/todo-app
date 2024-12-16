package repository

import (
	"github.com/yxrxy/todo-app/model"
	"gorm.io/gorm"
)

type TodoRepo struct {
	DB *gorm.DB
}

func (repo *TodoRepo) CreateTodo(todo *model.Todo) error {
	return repo.DB.Create(todo).Error
}

func (repo *TodoRepo) GetTodos(page, pageSize int) ([]model.Todo, int64, error) {
	var todos []model.Todo
	var total int64

	//获取总条数
	if err := repo.DB.Model(&model.Todo{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	//分页
	if err := repo.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&todos).Error; err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

func (repo *TodoRepo) UpdateTodoStatus(todoID uint, status int) error {
	return repo.DB.Model(&model.Todo{}).Where("id = ?", todoID).Update("status", status).Error
}

func (repo *TodoRepo) SearchTodos(keyword string, page, pageSize int) ([]model.Todo, int64, error) {
	var todos []model.Todo
	var total int64
	if err := repo.DB.Model(&model.Todo{}).
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := repo.DB.Model(&model.Todo{}).
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&todos).Error; err != nil {
		return nil, 0, err
	}
	return todos, total, nil
}
