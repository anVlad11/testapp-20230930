package manager

import (
	"fmt"
	"github.com/anvlad11/testapp-20230930/internal/interfaces/repositories"
	"github.com/anvlad11/testapp-20230930/internal/interfaces/services"
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"time"
)

type Service struct {
	workQueue  chan *model.Task
	inputQueue chan *model.Task

	taskRepository repositories.TaskRepository
}

func NewService(taskRepository repositories.TaskRepository, inputQueue chan *model.Task) *Service {
	svc := &Service{
		workQueue:  make(chan *model.Task),
		inputQueue: inputQueue,

		taskRepository: taskRepository,
	}

	return svc
}

func (s *Service) Start() {
	go func() {
		for {
			taskCount := s.load()
			if taskCount == 0 {
				time.Sleep(1 * time.Second)
			}
		}
	}()

	go func() {
		for task := range s.inputQueue {
			err := s.process(task)
			if err != nil {
				fmt.Printf("could not process task: %v", err)
			}
		}
	}()
}

func (s *Service) load() int {
	processingTasks, err := s.taskRepository.GetProcessing()
	if err != nil {
		fmt.Printf("could not load tasks from db: %v", err)
	}

	for _, task := range processingTasks {
		s.inputQueue <- task
	}

	return len(processingTasks)
}

func (s *Service) AddWorker(worker services.WorkerService) {
	worker.SetInput(s.workQueue)
}
