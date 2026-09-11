package dto

import (
	"time"

	"entrevista/core/domain/entity"
)

type CreateTaskRequest struct {
	Titulo        string    `json:"titulo" binding:"required,max=255"`
	Description   string    `json:"description" binding:"required"`
	ResponsableID uint      `json:"responsable_id" binding:"required"`
	FechaLimite   time.Time `json:"fecha_limite" binding:"required"`
	Estado        string    `json:"estado" binding:"omitempty,max=100"`
}

type UpdateTaskRequest struct {
	Titulo      string     `json:"titulo" binding:"omitempty,max=255"`
	Description string     `json:"description" binding:"omitempty"`
	FechaLimite *time.Time `json:"fecha_limite" binding:"omitempty"`
	Estado      string     `json:"estado" binding:"omitempty,max=100"`
}

// AssignTaskRequest reasigna una tarea existente a otro empleado.
type AssignTaskRequest struct {
	ResponsableID uint `json:"responsable_id" binding:"required"`
}

type UpdateEstadoRequest struct {
	Estado string `json:"estado" binding:"required,max=100"`
}

type TaskResponse struct {
	ID            uint              `json:"id"`
	Titulo        string            `json:"titulo"`
	Description   string            `json:"description"`
	ResponsableID uint              `json:"responsable_id"`
	Responsable   *EmployeeResponse `json:"responsable,omitempty"`
	FechaLimite   time.Time         `json:"fecha_limite"`
	Estado        string            `json:"estado"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

func NewTaskResponse(t *entity.Task) TaskResponse {
	res := TaskResponse{
		ID:            t.ID,
		Titulo:        t.Titulo,
		Description:   t.Description,
		ResponsableID: t.ResponsableID,
		FechaLimite:   t.FechaLimite,
		Estado:        t.Estado,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
	// El responsable solo viene cargado cuando el repositorio hizo Preload.
	if t.Responsable.ID != 0 {
		responsable := NewEmployeeResponse(&t.Responsable)
		res.Responsable = &responsable
	}
	return res
}

func NewTaskListResponse(tasks []entity.Task) []TaskResponse {
	out := make([]TaskResponse, 0, len(tasks))
	for i := range tasks {
		out = append(out, NewTaskResponse(&tasks[i]))
	}
	return out
}
