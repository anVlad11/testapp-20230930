package worker

import (
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func (s *Service) saveToDisk(task *model.Task) error {
	var err error

	if task.Content == "" {
		return nil
	}

	uri, _ := url.Parse(task.URL)

	path := uri.Path
	if uri.RawQuery != "" {
		path += "?" + uri.RawQuery
	}
	partsRaw := append([]string{uri.Host}, strings.Split(path, "/")...)
	invalidChars := []string{"<", ">", ":", "\"", "\\", "|", "?", "*"}

	parts := []string{}
	for i := range partsRaw {
		part := partsRaw[i]
		for _, char := range invalidChars {
			part = strings.ReplaceAll(part, char, "_")
		}
		if part != "" {
			parts = append(parts, part)
		}
	}

	if len(parts) == 1 {
		parts = append(parts, "index")
	}

	parts = append([]string{s.dataDir}, parts...)

	path = strings.Join(parts, string(os.PathSeparator))

	if filepath.Ext(path) == "" {
		if extension, exists := s.contentTypes[task.ContentType]; exists {
			if !strings.HasSuffix(path, extension) {
				path = path + extension
			}
		}
	}

	_, err = os.Stat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	err = nil

	err = os.MkdirAll(filepath.Dir(path), 0755)

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(task.Content)
	if err != nil {
		return err
	}

	return nil
}
