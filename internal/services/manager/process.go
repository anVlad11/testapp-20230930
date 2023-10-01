package manager

import (
	"fmt"
	"github.com/anvlad11/testapp-20230930/internal/model"
	"net/url"
)

func (s *Service) Process(task *model.Task) error {
	if task.URL == "" {
		return nil
	}

	if task.Error != nil {
		return task.Error
	}

	if task.Root == "" {
		domain, err := url.Parse(task.URL)
		if err != nil {
			return err
		}
		task.Root = domain.String()
	}

	if _, isDone := s.done[task.URL]; isDone {
		return nil
	}

	if !task.Downloaded {
		s.downloadQueue <- task
		return nil
	}

	if !task.IsContentTypeValid {
		task.Done = true
		delete(s.processing, task.URL)
		s.done[task.URL] = true
		fmt.Printf("invalid: %s\n", task.URL)

		return nil
	}

	if !task.Extracted {
		s.extractQueue <- task
		return nil
	}

	for _, link := range task.Links {
		if _, exists := s.processing[link]; exists {
			continue
		}
		if _, exists := s.done[link]; exists {
			continue
		}
		newTask := &model.Task{
			URL:   link,
			Links: []string{},
			Root:  task.Root,
		}
		go func(newTask *model.Task) { s.downloadQueue <- newTask }(newTask)
	}

	if !task.Done {
		task.Done = true
		delete(s.processing, task.URL)
		s.done[task.URL] = true
		fmt.Printf("done: %s\n", task.URL)
	}

	return nil
}
