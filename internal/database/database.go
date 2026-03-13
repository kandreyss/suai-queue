package database

import (
	"suai-queue/internal/config"
	"suai-queue/pkg/student"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DB.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&student.Student{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
