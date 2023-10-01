package database

import (
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(databasePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.Task{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
