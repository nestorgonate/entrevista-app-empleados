package port

import (
	"context"

	"entrevista/core/domain/entity"
)

// EmployeeRepository es el puerto de salida hacia la persistencia de empleados.
type EmployeeRepository interface {
	Create(ctx context.Context, employee *entity.Employee) error
	GetByID(ctx context.Context, id uint) (*entity.Employee, error)
	GetByCorreo(ctx context.Context, correo string) (*entity.Employee, error)
	List(ctx context.Context) ([]entity.Employee, error)
	Update(ctx context.Context, employee *entity.Employee) error
	Delete(ctx context.Context, id uint) error
}

// TaskRepository es el puerto de salida hacia la persistencia de tareas.
type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetByID(ctx context.Context, id uint) (*entity.Task, error)
	List(ctx context.Context) ([]entity.Task, error)
	ListByResponsable(ctx context.Context, responsableID uint) ([]entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id uint) error
}
