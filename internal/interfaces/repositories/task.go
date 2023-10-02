package repositories

import "github.com/anvlad11/testapp-20230930/pkg/model"

type TaskRepository interface {
	Save(task *model.Task) (*model.Task, error)
	GetProcessing(limit int) ([]*model.Task, error)
	GetByURL(url string) (*model.Task, error)
}
