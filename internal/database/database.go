package database

import (
	"suai-queue/internal/config"
	"suai-queue/internal/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DB.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Exec("PRAGMA foreign_keys = ON")

	err = db.AutoMigrate(&domain.Group{}, &domain.Student{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
