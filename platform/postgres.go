package platform

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresSQL() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
