package port

import (
	"context"

	"entrevista/core/domain/dto"
)

// EmployeeService es el puerto de entrada para los casos de uso de empleados.
type EmployeeService interface {
	Create(ctx context.Context, req dto.CreateEmployeeRequest) (dto.EmployeeResponse, error)
	GetByID(ctx context.Context, id uint) (dto.EmployeeResponse, error)
	List(ctx context.Context) ([]dto.EmployeeResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdateEmployeeRequest) (dto.EmployeeResponse, error)
	Delete(ctx context.Context, id uint) error
}

// TaskService es el puerto de entrada para los casos de uso de tareas.
type TaskService interface {
	Create(ctx context.Context, req dto.CreateTaskRequest) (dto.TaskResponse, error)
	GetByID(ctx context.Context, id uint) (dto.TaskResponse, error)
	List(ctx context.Context) ([]dto.TaskResponse, error)
	ListByResponsable(ctx context.Context, responsableID uint) ([]dto.TaskResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdateTaskRequest) (dto.TaskResponse, error)
	Assign(ctx context.Context, id uint, req dto.AssignTaskRequest) (dto.TaskResponse, error)
	UpdateEstado(ctx context.Context, id uint, req dto.UpdateEstadoRequest) (dto.TaskResponse, error)
	Delete(ctx context.Context, id uint) error
}
