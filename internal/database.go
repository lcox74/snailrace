package internal

import (
	"github.com/lcox74/snailrace/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDatabase() (*gorm.DB, error) {
	// Open Database Connection, or Create Database if it doesn't exist
	db, err := gorm.Open(sqlite.Open("snailrace.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schemas
	err = db.AutoMigrate(&models.Snail{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
