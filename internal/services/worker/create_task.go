package worker

import "github.com/anvlad11/testapp-20230930/pkg/model"

func (s *Service) createTask(link string, root string) error {
	var err error

	var existingTask *model.Task
	existingTask, err = s.taskRepository.GetByURL(link)
	if err != nil {
		return err
	}

	if existingTask != nil {
		return nil
	}

	newTask := &model.Task{
		URL:  link,
		Root: root,
	}

	newTask, err = s.taskRepository.Save(newTask)
	if err != nil {
		return err
	}

	return nil
}
