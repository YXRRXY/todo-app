package repository

import (
	"github.com/YXRRXY/todo-app/model"
	"gorm.io/gorm"
)

type TodoRepo struct {
	DB *gorm.DB
}

func (repo *TodoRepo) CreateTodo(todo *model.Todo) error {
	return repo.DB.Create(todo).Error
}

func (repo *TodoRepo) GetTodos(userID uint, page, pageSize int, status *int) ([]model.Todo, int64, error) {
	var todos []model.Todo
	var total int64
	query := repo.DB.Model(&model.Todo{}).Where("user_id = ?", userID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&todos).Error; err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

func (repo *TodoRepo) UpdateTodoStatus(userID uint, todoID uint, status int) error {
	return repo.DB.Model(&model.Todo{}).
		Where("id = ? AND user_id = ?", todoID, userID).
		Update("status", status).Error
}

func (repo *TodoRepo) SearchTodos(userID uint, keyword string, page, pageSize int, status *int) ([]model.Todo, int64, error) {
	var todos []model.Todo
	var total int64

	query := repo.DB.Model(&model.Todo{}).
		Where("user_id = ?", userID).
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&todos).Error; err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

func (repo *TodoRepo) BatchUpdateStatus(userID uint, status int, currentStatus *int, ids []uint) (int64, error) {
	query := repo.DB.Model(&model.Todo{}).Where("user_id = ?", userID)

	if len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}

	if currentStatus != nil {
		query = query.Where("status = ?", *currentStatus)
	}

	result := query.Update("status", status)
	return result.RowsAffected, result.Error
}

func (repo *TodoRepo) DeleteTodo(userID uint, todoID uint) error {
	return repo.DB.Where("id = ? AND user_id = ?", todoID, userID).Delete(&model.Todo{}).Error
}

func (repo *TodoRepo) BatchDelete(userID uint, status *int, ids []uint) (int64, error) {
	query := repo.DB.Where("user_id = ?", userID)

	if len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	result := query.Delete(&model.Todo{})
	return result.RowsAffected, result.Error
}
