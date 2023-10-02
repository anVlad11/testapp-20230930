package task

import (
	"errors"
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"gorm.io/gorm"
)

func (r *Repository) GetByURL(url string) (*model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var task *model.Task

	err := r.db.
		Where("url = ?", url).
		First(&task).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return task, nil
}
