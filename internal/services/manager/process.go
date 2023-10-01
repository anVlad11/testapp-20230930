package manager

import (
	"errors"
	"fmt"
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"net/url"
)

func (s *Service) Process(task *model.Task) error {
	if task.URL == "" {
		return nil
	}

	if task.Error != "" {
		return errors.New(task.Error)
	}

	if task.Root == "" {
		domain, err := url.Parse(task.URL)
		if err != nil {
			return err
		}
		task.Root = domain.String()
	}

	if _, isDone := s.done.Load(task.URL); isDone {
		return nil
	}

	if !task.Downloaded {
		s.downloadQueue <- task
		return nil
	}

	if !task.IsContentTypeValid {
		task.Done = true
		s.processing.Delete(task.URL)
		s.done.Store(task.URL, true)
		fmt.Printf("invalid: %s\n", task.URL)

		return nil
	}

	if !task.Extracted {
		s.extractQueue <- task
		return nil
	}

	for _, link := range task.Links {
		if _, exists := s.processing.Load(link); exists {
			continue
		}
		if _, exists := s.done.Load(link); exists {
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
		s.processing.Delete(task.URL)
		s.done.Store(task.URL, true)
		fmt.Printf("done: %s\n", task.URL)
	}

	return nil
}
