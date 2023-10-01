package downloader

import (
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"sync"
	"time"
)

type Service struct {
	input  chan *model.Task
	output chan *model.Task

	isRunning bool
	mu        sync.Mutex

	validContentTypes map[string]bool
}

func NewService(contentTypes map[string]string) *Service {
	validContentTypes := map[string]bool{}
	for contentType := range contentTypes {
		validContentTypes[contentType] = true
	}

	return &Service{}
}

func (s *Service) SetPipe(input chan *model.Task, output chan *model.Task) {
	s.output = output
	s.input = input
}

func (s *Service) Start() {
	s.isRunning = true
	go func() {
		for s.isRunning {
			s.mu.Lock()
			select {
			case task, ok := <-s.input:
				if !ok {
					return
				}
				go func() {
					err := s.download(task)
					s.mu.Unlock()

					if err != nil {
						task.Error = err.Error()
					}
					task.Downloaded = true
					s.output <- task
				}()
			default:
				s.mu.Unlock()

				time.Sleep(10 * time.Millisecond)
			}
		}
	}()
}

func (s *Service) Stop() {
	s.isRunning = false
}
