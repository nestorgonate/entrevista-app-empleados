package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"entrevista/core/domain"
	"entrevista/core/domain/entity"
	"entrevista/core/port"
)

type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository construye el adaptador GORM del puerto TaskRepository.
func NewTaskRepository(db *gorm.DB) port.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *entity.Task) error {
	// Omit evita que GORM haga upsert del empleado asociado al crear la tarea.
	return r.db.WithContext(ctx).Omit("Responsable").Create(task).Error
}

func (r *taskRepository) GetByID(ctx context.Context, id uint) (*entity.Task, error) {
	var task entity.Task
	if err := r.db.WithContext(ctx).Preload("Responsable").First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) List(ctx context.Context) ([]entity.Task, error) {
	var tasks []entity.Task
	if err := r.db.WithContext(ctx).Preload("Responsable").Order("id asc").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) ListByResponsable(ctx context.Context, responsableID uint) ([]entity.Task, error) {
	var tasks []entity.Task
	err := r.db.WithContext(ctx).
		Preload("Responsable").
		Where("responsable_id = ?", responsableID).
		Order("id asc").
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *entity.Task) error {
	return r.db.WithContext(ctx).Omit("Responsable").Save(task).Error
}

func (r *taskRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&entity.Task{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
