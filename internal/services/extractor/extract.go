package extractor

import (
	"github.com/anvlad11/testapp-20230930/internal/model"
	"golang.org/x/net/html"
	"net/url"
	"strings"
)

func (s *Service) extract(task *model.Task) error {
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
