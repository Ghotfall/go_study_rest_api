package db

import (
	"go_study_rest_api/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func init() {
	sqliteDB, dbError := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if dbError != nil {
		log.Fatalf("Failed to open DB: %s\n", dbError.Error())
	}

	DB = sqliteDB

	migrateErr := DB.AutoMigrate(&models.Student{})
	if migrateErr != nil {
		log.Fatalf("Failed to migrate schema: %s\n", migrateErr.Error())
	}
}
