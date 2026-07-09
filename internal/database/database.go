package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"ticket-system/internal/models"
)

// Connect opens a SQLite database at the given path and auto-migrates
// the schema. It exits the process on failure since the service cannot
// run without a working database.
func Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Ticket{}); err != nil {
		log.Fatalf("failed to auto-migrate database: %v", err)
	}

	return db
}
