package task

import (
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Repository struct {
	dataDir string
	db      *gorm.DB
}

func NewRepository(db *gorm.DB, dataDir string) *Repository {
	return &Repository{db: db, dataDir: dataDir}
}

func (r *Repository) Save(task *model.Task) error {
	err := r.saveToFS(task)
	if err != nil {
		return err
	}

	err = r.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}, {Name: "url"}},
			UpdateAll: true,
		}).
		Model(&model.Task{}).
		Create(task).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAll() ([]*model.Task, error) {
	var tasks []*model.Task

	err := r.db.
		Model(&model.Task{}).
		Find(&tasks).
		Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) saveToFS(task *model.Task) error {
	var err error

	if !task.Done {
		return nil
	}

	if !task.IsContentTypeValid {
		return nil
	}

	if task.Content == "" {
		return nil
	}

	uri, _ := url.Parse(task.URL)

	path := uri.Path
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

	parts = append([]string{r.dataDir}, parts...)

	path = strings.Join(parts, string(os.PathSeparator))

	if filepath.Ext(path) == "" {
		extensions := map[string]string{
			"text/html":              ".html",
			"text/css":               ".css",
			"application/json":       ".json",
			"application/javascript": ".js",
		}

		if extension, exists := extensions[task.ContentType]; exists {
			if !strings.HasSuffix(path, extension) {
				path = path + extension
			}
		}
	}

	_, err = os.Stat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil {
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
