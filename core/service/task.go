package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"entrevista/core/domain"
	"entrevista/core/domain/dto"
	"entrevista/core/domain/entity"
	"entrevista/core/port"
)

type taskService struct {
	tasks     port.TaskRepository
	employees port.EmployeeRepository
}

// NewTaskService inyecta los repositorios que necesita el caso de uso de tareas.
// El repositorio de empleados se usa para validar que el responsable exista
// antes de asignarle una tarea.
func NewTaskService(tasks port.TaskRepository, employees port.EmployeeRepository) port.TaskService {
	return &taskService{tasks: tasks, employees: employees}
}

func (s *taskService) Create(ctx context.Context, req dto.CreateTaskRequest) (dto.TaskResponse, error) {
	if err := s.validarResponsable(ctx, req.ResponsableID); err != nil {
		return dto.TaskResponse{}, err
	}

	estado := strings.TrimSpace(req.Estado)
	if estado == "" {
		estado = entity.EstadoPendiente
	}
	if !entity.EstadoValido(estado) {
		return dto.TaskResponse{}, fmt.Errorf("%w: estado %q no soportado", domain.ErrInvalidInput, estado)
	}

	task := entity.Task{
		Titulo:        strings.TrimSpace(req.Titulo),
		Description:   strings.TrimSpace(req.Description),
		ResponsableID: req.ResponsableID,
		FechaLimite:   req.FechaLimite,
		Estado:        estado,
	}
	if err := s.tasks.Create(ctx, &task); err != nil {
		return dto.TaskResponse{}, err
	}

	// Se relee para devolver el responsable ya precargado.
	return s.GetByID(ctx, task.ID)
}

func (s *taskService) GetByID(ctx context.Context, id uint) (dto.TaskResponse, error) {
	task, err := s.tasks.GetByID(ctx, id)
	if err != nil {
		return dto.TaskResponse{}, err
	}
	return dto.NewTaskResponse(task), nil
}

func (s *taskService) List(ctx context.Context) ([]dto.TaskResponse, error) {
	tasks, err := s.tasks.List(ctx)
	if err != nil {
		return nil, err
	}
	return dto.NewTaskListResponse(tasks), nil
}

func (s *taskService) ListByResponsable(ctx context.Context, responsableID uint) ([]dto.TaskResponse, error) {
	if err := s.validarResponsable(ctx, responsableID); err != nil {
		return nil, err
	}
	tasks, err := s.tasks.ListByResponsable(ctx, responsableID)
	if err != nil {
		return nil, err
	}
	return dto.NewTaskListResponse(tasks), nil
}

func (s *taskService) Update(ctx context.Context, id uint, req dto.UpdateTaskRequest) (dto.TaskResponse, error) {
	task, err := s.tasks.GetByID(ctx, id)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	if req.Titulo != "" {
		task.Titulo = strings.TrimSpace(req.Titulo)
	}
	if req.Description != "" {
		task.Description = strings.TrimSpace(req.Description)
	}
	if req.FechaLimite != nil {
		task.FechaLimite = *req.FechaLimite
	}
	if req.Estado != "" {
		estado := strings.TrimSpace(req.Estado)
		if !entity.EstadoValido(estado) {
			return dto.TaskResponse{}, fmt.Errorf("%w: estado %q no soportado", domain.ErrInvalidInput, estado)
		}
		task.Estado = estado
	}

	if err := s.tasks.Update(ctx, task); err != nil {
		return dto.TaskResponse{}, err
	}
	return dto.NewTaskResponse(task), nil
}

// Assign reasigna una tarea existente a otro empleado.
func (s *taskService) Assign(ctx context.Context, id uint, req dto.AssignTaskRequest) (dto.TaskResponse, error) {
	task, err := s.tasks.GetByID(ctx, id)
	if err != nil {
		return dto.TaskResponse{}, err
	}
	if err := s.validarResponsable(ctx, req.ResponsableID); err != nil {
		return dto.TaskResponse{}, err
	}

	task.ResponsableID = req.ResponsableID
	if err := s.tasks.Update(ctx, task); err != nil {
		return dto.TaskResponse{}, err
	}
	return s.GetByID(ctx, task.ID)
}

func (s *taskService) UpdateEstado(ctx context.Context, id uint, req dto.UpdateEstadoRequest) (dto.TaskResponse, error) {
	estado := strings.TrimSpace(req.Estado)
	if !entity.EstadoValido(estado) {
		return dto.TaskResponse{}, fmt.Errorf("%w: estado %q no soportado", domain.ErrInvalidInput, estado)
	}

	task, err := s.tasks.GetByID(ctx, id)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	task.Estado = estado
	if err := s.tasks.Update(ctx, task); err != nil {
		return dto.TaskResponse{}, err
	}
	return dto.NewTaskResponse(task), nil
}

func (s *taskService) Delete(ctx context.Context, id uint) error {
	return s.tasks.Delete(ctx, id)
}

// validarResponsable confirma que el empleado exista antes de asignarle trabajo.
func (s *taskService) validarResponsable(ctx context.Context, responsableID uint) error {
	if _, err := s.employees.GetByID(ctx, responsableID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fmt.Errorf("%w: el empleado %d no existe", domain.ErrInvalidInput, responsableID)
		}
		return err
	}
	return nil
}
