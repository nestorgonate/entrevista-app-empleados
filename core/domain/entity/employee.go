package entity

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	Nombre string `gorm:"type:varchar(25);not null"`
	Correo string `gorm:"type:varchar(255);not null"`
	Cargo  string `gorm:"type:text;not null"`
}
