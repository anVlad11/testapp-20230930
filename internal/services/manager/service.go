package manager

import (
	"fmt"
	"github.com/anvlad11/testapp-20230930/internal/interfaces/services"
	"github.com/anvlad11/testapp-20230930/internal/model"
	"net/url"
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
					err := s.ProcessTask(task)
					s.mu.Unlock()
					if err != nil {
						delete(s.processing, task.URL)
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

func (s *Service) ProcessTask(task *model.Task) error {
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
