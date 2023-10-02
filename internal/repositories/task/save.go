package task

import (
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"gorm.io/gorm/clause"
)

func (r *Repository) Save(task *model.Task) (*model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var err error
	if task.ID == 0 {
		err = r.create(task)
	} else {
		err = r.update(task)
	}

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *Repository) create(task *model.Task) error {
	var err error

	if err != nil {
		return err
	}

	err = r.db.
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}, {Name: "url"}},
				DoNothing: true,
			},
		).
		Create(&task).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) update(task *model.Task) error {
	var err error

	if err != nil {
		return err
	}

	err = r.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}, {Name: "url"}},
			UpdateAll: true,
		}).
		Model(&model.Task{}).
		Create(task).
		Error

	if err != nil {
		return err
	}

	return nil
}
