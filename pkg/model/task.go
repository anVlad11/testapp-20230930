package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Task struct {
	ID                 int64     `db:"id" gorm:"primarykey"`
	URL                string    `db:"url" gorm:"uniqueIndex;index"`
	Root               string    `db:"root" gorm:"index"`
	ContentType        string    `db:"content_type"`
	IsContentTypeValid bool      `db:"is_content_type_valid"`
	Content            string    `db:"-" gorm:"-"`
	Downloaded         bool      `db:"downloaded"`
	Links              Links     `db:"links" gorm:"type:json"`
	Extracted          bool      `db:"extracted"`
	Done               bool      `db:"done"`
	Error              string    `db:"error"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}

func (m *Task) TableName() string {
	return "tasks"
}

type Links []string

func (m *Links) Scan(value interface{}) error {
	var bytes []byte
	switch typedValue := value.(type) {
	case string:
		bytes = []byte(typedValue)
	case []byte:
		bytes = typedValue
	default:
		return fmt.Errorf("failed to unmarshal links, unsupported type")
	}

	var m2 Links
	err := json.Unmarshal(bytes, &m2)
	if err != nil {
		return err
	}

	*m = m2

	return nil
}

func (m Links) Value() (driver.Value, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return string(bytes), nil
}
