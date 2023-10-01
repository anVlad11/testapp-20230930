package manager

import (
	"fmt"
	"github.com/anvlad11/testapp-20230930/internal/interfaces/services"
	"github.com/anvlad11/testapp-20230930/internal/model"
	"sync"
	"time"
)

type Service struct {
	downloadQueue chan *model.Task
	extractQueue  chan *model.Task
	inputQueue    chan *model.Task

	isRunning bool
	mu        sync.Mutex

	processing map[string]bool
	done       map[string]bool
}

func NewService() *Service {
	svc := &Service{
		downloadQueue: make(chan *model.Task, 1),
		extractQueue:  make(chan *model.Task, 1),
		inputQueue:    make(chan *model.Task, 1),

		processing: map[string]bool{},
		done:       map[string]bool{},
	}

	return svc
}

func (s *Service) Start() {
	s.isRunning = true
	go func() {
		for s.isRunning {
			s.mu.Lock()

			select {
			case task, ok := <-s.inputQueue:
				if !ok {
					return
				}
				go func() {
					err := s.Process(task)
					s.mu.Unlock()
					if err != nil {
						delete(s.processing, task.URL)
						s.done[task.URL] = false
						fmt.Println(err)
					}
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

func (s *Service) AddDownloader(downloader services.WorkerService) {
	downloader.SetPipe(s.downloadQueue, s.inputQueue)
}

func (s *Service) AddExtractor(downloader services.WorkerService) {
	downloader.SetPipe(s.extractQueue, s.inputQueue)
}
