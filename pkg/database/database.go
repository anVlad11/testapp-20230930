package database

import (
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func NewDatabase(databasePath string) (*gorm.DB, error) {
	dbLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             500 * time.Millisecond,
		LogLevel:                  logger.Error,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	})

	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.Task{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
