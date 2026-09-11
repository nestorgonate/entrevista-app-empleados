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

type employeeService struct {
	repo port.EmployeeRepository
}

// NewEmployeeService inyecta el repositorio de empleados en el caso de uso.
func NewEmployeeService(repo port.EmployeeRepository) port.EmployeeService {
	return &employeeService{repo: repo}
}

func (s *employeeService) Create(ctx context.Context, req dto.CreateEmployeeRequest) (dto.EmployeeResponse, error) {
	correo := strings.ToLower(strings.TrimSpace(req.Correo))

	if _, err := s.repo.GetByCorreo(ctx, correo); err == nil {
		return dto.EmployeeResponse{}, fmt.Errorf("%w: ya existe un empleado con el correo %s", domain.ErrConflict, correo)
	} else if !errors.Is(err, domain.ErrNotFound) {
		return dto.EmployeeResponse{}, err
	}

	employee := entity.Employee{
		Nombre: strings.TrimSpace(req.Nombre),
		Correo: correo,
		Cargo:  strings.TrimSpace(req.Cargo),
	}
	if err := s.repo.Create(ctx, &employee); err != nil {
		return dto.EmployeeResponse{}, err
	}
	return dto.NewEmployeeResponse(&employee), nil
}

func (s *employeeService) GetByID(ctx context.Context, id uint) (dto.EmployeeResponse, error) {
	employee, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.EmployeeResponse{}, err
	}
	return dto.NewEmployeeResponse(employee), nil
}

func (s *employeeService) List(ctx context.Context) ([]dto.EmployeeResponse, error) {
	employees, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return dto.NewEmployeeListResponse(employees), nil
}

func (s *employeeService) Update(ctx context.Context, id uint, req dto.UpdateEmployeeRequest) (dto.EmployeeResponse, error) {
	employee, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.EmployeeResponse{}, err
	}

	if req.Nombre != "" {
		employee.Nombre = strings.TrimSpace(req.Nombre)
	}
	if req.Cargo != "" {
		employee.Cargo = strings.TrimSpace(req.Cargo)
	}
	if req.Correo != "" {
		correo := strings.ToLower(strings.TrimSpace(req.Correo))
		if correo != employee.Correo {
			existente, err := s.repo.GetByCorreo(ctx, correo)
			switch {
			case err == nil && existente.ID != employee.ID:
				return dto.EmployeeResponse{}, fmt.Errorf("%w: ya existe un empleado con el correo %s", domain.ErrConflict, correo)
			case err != nil && !errors.Is(err, domain.ErrNotFound):
				return dto.EmployeeResponse{}, err
			}
			employee.Correo = correo
		}
	}

	if err := s.repo.Update(ctx, employee); err != nil {
		return dto.EmployeeResponse{}, err
	}
	return dto.NewEmployeeResponse(employee), nil
}

func (s *employeeService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
