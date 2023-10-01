package extractor

import (
	"github.com/anvlad11/testapp-20230930/internal/model"
	"golang.org/x/net/html"
	"net/url"
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
					err := s.crawl(task)
					s.mu.Unlock()

					if err != nil {
						task.Error = err
					}
					task.Extracted = true
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

func (s *Service) crawl(task *model.Task) error {
	doc, err := html.Parse(strings.NewReader(task.Content))
	if err != nil {
		return err
	}

	root, err := url.Parse(task.Root)
	if err != nil {
		return err
	}
	task.Links = []string{}

	htmlNodeQueue := []*html.Node{doc}
	for {
		if len(htmlNodeQueue) == 0 {
			break
		}
		n := htmlNodeQueue[0]

		htmlNodeQueue = htmlNodeQueue[1:]
		if n.Type == html.ElementNode {
			attrs := []string{"src", "href", "content"}
			for _, a := range n.Attr {
				for _, attr := range attrs {
					if a.Key == attr {
						var link *url.URL
						link, err = url.Parse(a.Val)
						if err != nil {
							//log.Printf("error parsing %s: %s", a.Val, err.Error())
							continue
						}

						resolved := root.ResolveReference(link)
						// We're not a browser
						resolved.Fragment = ""
						if resolved.Host == root.Host {
							task.Links = append(task.Links, resolved.String())
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			htmlNodeQueue = append(htmlNodeQueue, c)
		}
	}

	return nil
}
