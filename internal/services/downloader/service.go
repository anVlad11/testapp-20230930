package downloader

import (
	"github.com/anvlad11/testapp-20230930/internal/model"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Service struct {
	input  chan *model.Task
	output chan *model.Task

	isRunning bool
	mu        sync.Mutex
}

func NewService() *Service {
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
						task.Error = err
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

func (s *Service) download(task *model.Task) error {
	contentType, err := s.getContentType(task.URL)
	if err != nil {
		return err
	}
	task.ContentType = contentType
	if !s.checkContentType(contentType) {
		return nil
	}

	task.IsContentTypeValid = true

	httpClient := http.Client{}
	resp, err := httpClient.Get(task.URL)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	var data []byte
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	task.Content = string(data)

	return nil
}

func (s *Service) getContentType(url string) (string, error) {
	httpClient := http.Client{}

	resp, err := httpClient.Head(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}

	contentType := resp.Header.Get("Content-Type")
	contentTypePart := strings.Split(contentType, ";")[0]

	return contentTypePart, nil
}

func (s *Service) checkContentType(contentType string) bool {
	validContentTypes := map[string]bool{
		"text/html":              true,
		"text/css":               true,
		"application/javascript": true,
		"application/json":       true,
	}

	valid, exists := validContentTypes[contentType]

	return valid && exists
}
