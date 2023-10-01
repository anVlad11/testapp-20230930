package database

import (
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.Task{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
