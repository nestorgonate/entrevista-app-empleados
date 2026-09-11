package dto

import (
	"time"

	"entrevista/core/domain/entity"
)

type CreateEmployeeRequest struct {
	Nombre string `json:"nombre" binding:"required,max=25"`
	Correo string `json:"correo" binding:"required,email,max=255"`
	Cargo  string `json:"cargo" binding:"required"`
}

type UpdateEmployeeRequest struct {
	Nombre string `json:"nombre" binding:"omitempty,max=25"`
	Correo string `json:"correo" binding:"omitempty,email,max=255"`
	Cargo  string `json:"cargo" binding:"omitempty"`
}

type EmployeeResponse struct {
	ID        uint      `json:"id"`
	Nombre    string    `json:"nombre"`
	Correo    string    `json:"correo"`
	Cargo     string    `json:"cargo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewEmployeeResponse(e *entity.Employee) EmployeeResponse {
	return EmployeeResponse{
		ID:        e.ID,
		Nombre:    e.Nombre,
		Correo:    e.Correo,
		Cargo:     e.Cargo,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func NewEmployeeListResponse(employees []entity.Employee) []EmployeeResponse {
	out := make([]EmployeeResponse, 0, len(employees))
	for i := range employees {
		out = append(out, NewEmployeeResponse(&employees[i]))
	}
	return out
}
