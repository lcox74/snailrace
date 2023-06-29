package internal

import (
	"github.com/lcox74/snailrace/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

func SetupDatabase() (*gorm.DB, error) {
	// Open Database Connection, or Create Database if it doesn't exist
	db, err := gorm.Open(sqlite.Open("snailrace.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schemas
	return db, MigrateSchemas(db)
}

func MigrateSchemas(db *gorm.DB) error {

	// Schema List
	schemas := []interface{}{
		&models.User{},
		&models.Snail{},
	}

	// Migrate the schemas
	for _, schema := range schemas {
		err := db.AutoMigrate(schema)
		if err != nil {
			log.Warnf("Failed migrating schema: %s", err)
			return err
		}
	}

	return nil
}
