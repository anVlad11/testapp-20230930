package manager

import (
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"net/url"
)

func (s *Service) process(task *model.Task) error {
	if task.URL == "" {
		return nil
	}

	if task.Root == "" {
		domain, err := url.Parse(task.URL)
		if err != nil {
			return err
		}
		task.Root = domain.String()
	}

	if task.Links == nil {
		task.Links = []string{}
	}

	if !task.Done {
		s.workQueue <- task
	}

	return nil
}
