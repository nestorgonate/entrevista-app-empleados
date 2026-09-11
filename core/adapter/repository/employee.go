package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"entrevista/core/domain"
	"entrevista/core/domain/entity"
	"entrevista/core/port"
)

type employeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository construye el adaptador GORM del puerto EmployeeRepository.
func NewEmployeeRepository(db *gorm.DB) port.EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(ctx context.Context, employee *entity.Employee) error {
	return r.db.WithContext(ctx).Create(employee).Error
}

func (r *employeeRepository) GetByID(ctx context.Context, id uint) (*entity.Employee, error) {
	var employee entity.Employee
	if err := r.db.WithContext(ctx).First(&employee, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) GetByCorreo(ctx context.Context, correo string) (*entity.Employee, error) {
	var employee entity.Employee
	if err := r.db.WithContext(ctx).Where("correo = ?", correo).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) List(ctx context.Context) ([]entity.Employee, error) {
	var employees []entity.Employee
	if err := r.db.WithContext(ctx).Order("id asc").Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *employeeRepository) Update(ctx context.Context, employee *entity.Employee) error {
	return r.db.WithContext(ctx).Save(employee).Error
}

func (r *employeeRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&entity.Employee{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
