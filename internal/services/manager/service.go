package manager

import (
	"fmt"
	"github.com/anvlad11/testapp-20230930/internal/interfaces/repositories"
	"github.com/anvlad11/testapp-20230930/internal/interfaces/services"
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"github.com/anvlad11/testapp-20230930/pkg/xsync"
	"sync"
	"time"
)

type Service struct {
	downloadQueue chan *model.Task
	extractQueue  chan *model.Task
	inputQueue    chan *model.Task

	isRunning bool
	mu        sync.Mutex

	processing xsync.Map[string, bool]
	done       xsync.Map[string, bool]

	taskRepository repositories.TaskRepository
}

func NewService(taskRepository repositories.TaskRepository) *Service {
	svc := &Service{
		downloadQueue: make(chan *model.Task, 1),
		extractQueue:  make(chan *model.Task, 1),
		inputQueue:    make(chan *model.Task, 1),

		processing: xsync.Map[string, bool]{},
		done:       xsync.Map[string, bool]{},

		taskRepository: taskRepository,
	}

	return svc
}

func (s *Service) Start() {
	s.isRunning = true

	s.processing = xsync.Map[string, bool]{}
	s.done = xsync.Map[string, bool]{}

	s.load()

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
						s.processing.Delete(task.URL)
						s.done.Store(task.URL, false)
						fmt.Println(err)
					}
					if task != nil {
						err = s.taskRepository.Save(task)
						if err != nil {
							fmt.Printf("could not save task: %v", err)
						}
					}
				}()
			default:
				s.mu.Unlock()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()
}

func (s *Service) load() {
	existingTasks, err := s.taskRepository.GetAll()
	if err != nil {
		fmt.Printf("could not load tasks from db: %v", err)
	}

	for _, task := range existingTasks {
		if task.Done {
			s.done.Store(task.URL, true)
			continue
		}
		s.processing.Store(task.URL, true)
		go func(task *model.Task) { _ = s.Process(task) }(task)
	}
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
