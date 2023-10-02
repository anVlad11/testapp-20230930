package worker

import (
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"io"
	"net/http"
	"strings"
)

func (s *Service) downloadContent(task *model.Task) error {
	contentType, err := s.getContentType(task.URL)
	if err != nil {
		return err
	}
	task.ContentType = contentType
	if !s.checkContentType(contentType) {
		return nil
	}

	task.IsContentTypeValid = true

	request, err := http.NewRequest(http.MethodGet, task.URL, nil)
	if err != nil {
		return err
	}

	for k, v := range s.requestHeaders {
		request.Header.Set(k, v)
	}

	response, err := s.httpClient.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return err
	}

	var data []byte
	data, err = io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	task.Content = string(data)

	return nil
}

func (s *Service) getContentType(url string) (string, error) {
	request, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return "", err
	}

	for k, v := range s.requestHeaders {
		request.Header.Set(k, v)
	}

	response, err := s.httpClient.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return "", err
	}

	contentType := response.Header.Get("Content-Type")
	contentTypePart := strings.Split(contentType, ";")[0]

	return contentTypePart, nil
}

func (s *Service) checkContentType(contentType string) bool {
	_, exists := s.contentTypes[contentType]

	return exists
}
