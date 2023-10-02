package worker

import (
	"fmt"
	"github.com/anvlad11/testapp-20230930/internal/interfaces/repositories"
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"log"
	"net/http"
	"time"
)

type Service struct {
	input chan *model.Task

	dataDir        string
	contentTypes   map[string]string
	httpClient     *http.Client
	requestHeaders map[string]string

	taskRepository repositories.TaskRepository
}

func NewService(
	contentTypes map[string]string,
	requestHeaders map[string]string,
	dataDir string,
	taskRepository repositories.TaskRepository,
) *Service {
	return &Service{
		contentTypes:   contentTypes,
		dataDir:        dataDir,
		httpClient:     http.DefaultClient,
		requestHeaders: requestHeaders,
		taskRepository: taskRepository,
	}
}

func (s *Service) SetInput(input chan *model.Task) {
	s.input = input
}

func (s *Service) Start() {
	go func() {
		for task := range s.input {
			s.process(task)
		}
	}()
}

func (s *Service) process(task *model.Task) {
	start := time.Now()
	err := s.downloadContent(task)

	if err != nil {
		task.Error = err.Error()
	}

	if task.IsContentTypeValid {
		err = s.saveToDisk(task)
		if err != nil {
			task.Error = err.Error()
		}

		err = s.extractLinks(task)
		if err != nil {
			task.Error = err.Error()
		}
	}

	task.Done = true

	fmt.Printf(
		"done, valid: %v, took %.2dms: %s\n",
		task.IsContentTypeValid,
		time.Since(start).Milliseconds(),
		task.URL,
	)

	task, err = s.taskRepository.Save(task)
	if err != nil {
		fmt.Printf("could not save task after processing: %v\n", err.Error())
		return
	}

	for _, link := range task.Links {
		err = s.createTask(link, task.Root)
		if err != nil {
			log.Printf("could not save new task for %s: %v", link, err)
		}
	}

}
