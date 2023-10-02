package task

import (
	"gorm.io/gorm"
	"sync"
)

type Repository struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}
