package repositories

import "github.com/anvlad11/testapp-20230930/pkg/model"

type TaskRepository interface {
	Save(task *model.Task) error
	GetAll() ([]*model.Task, error)
}
