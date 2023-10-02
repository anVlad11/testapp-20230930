package task

import "github.com/anvlad11/testapp-20230930/pkg/model"

func (r *Repository) GetProcessing() ([]*model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var tasks []*model.Task

	err := r.db.
		Model(&model.Task{}).
		Where("done = 0").
		Find(&tasks).
		Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}
