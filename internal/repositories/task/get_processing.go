package task

import "github.com/anvlad11/testapp-20230930/pkg/model"

func (r *Repository) GetProcessing(limit int) ([]*model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var tasks []*model.Task

	q := r.db.
		Model(&model.Task{}).
		Where("done = 0")

	if limit > 0 {
		q = q.Limit(limit)
	}

	err := q.Find(&tasks).
		Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}
