package platform

import (
	"gorm.io/gorm"

	"entrevista/core/domain/entity"
)

// Migrate crea o actualiza el esquema de las entidades del dominio.
// El orden importa: Employee primero, porque Task tiene la FK hacia el.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Employee{},
		&entity.Task{},
	)
}
