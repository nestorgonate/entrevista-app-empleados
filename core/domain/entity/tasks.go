package entity

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Titulo        string `gorm:"type:varchar(255);not null"`
	Description   string `gorm:"type:text;not null"`
	ResponsableID uint   `gorm:"type:bigint;not null"`
	//Declara ResponsableID como foreign key a employee
	//Actualmente hace set null en delete, pero se podria restringir como buena practica para no perder registros
	Responsable Employee  `gorm:"foreignKey:ResponsableID;references:ID"`
	FechaLimite time.Time `gorm:"type:timestamptz;not null"`
	Estado      string    `gorm:"type:varchar(100);not null"`
}
