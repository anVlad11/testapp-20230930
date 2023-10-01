package model

type Task struct {
	URL                string
	Root               string
	ContentType        string
	IsContentTypeValid bool
	Content            string
	Downloaded         bool
	Links              []string
	Extracted          bool
	Done               bool
	Error              error
}
