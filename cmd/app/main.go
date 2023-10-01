package main

import (
	"errors"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	err := do()
	if err != nil {
		log.Fatal(err)
	}
}

var rootVar = flag.String("root", "", "Page to start crawling from")

func do() error {
	var err error

	flag.Parse()
	if rootVar == nil {
		return errors.New("root is empty")
	}

	root := *rootVar

	domain, err := url.Parse(root)
	if err != nil {
		return err
	}

	todo := map[string]string{}
	todo[root] = ""
	done := map[string]string{}
	for {
		if len(todo) == 0 {
			break
		}
		for path := range todo {
			delete(todo, path)
			if _, exists := done[path]; exists {
				continue
			}
			fmt.Println(path)
			var content string
			content, err = download(path)
			if err != nil {
				return err
			}
			if content == "" {
				continue
			}
			done[path] = content
			var paths []string
			paths, err = crawl(content, domain)
			if err != nil {
				return err
			}
			for _, newPath := range paths {
				if _, exists := done[newPath]; exists {
					continue
				}
				todo[newPath] = ""
			}
		}
	}

	return nil
}

func crawl(content string, root *url.URL) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	var links []string

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
							log.Printf("error parsing %s: %s", a.Val, err.Error())
							continue
						}

						resolved := root.ResolveReference(link)
						// We're not browser
						resolved.Fragment = ""
						if resolved.Host == root.Host {
							links = append(links, resolved.String())
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			htmlNodeQueue = append(htmlNodeQueue, c)
		}
	}

	return links, nil
}

func download(url string) (string, error) {
	contentType, err := getContentType(url)
	if err != nil {
		return "", err
	}

	if !checkContentType(contentType) {
		return "", nil
	}

	fmt.Println(contentType)

	httpClient := http.Client{}
	resp, err := httpClient.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}

	var data []byte
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getContentType(url string) (string, error) {
	httpClient := http.Client{}

	resp, err := httpClient.Head(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}

	contentType := resp.Header.Get("Content-Type")
	return contentType, nil
}

func checkContentType(contentType string) bool {
	contentTypePart := strings.Split(contentType, ";")[0]
	validContentTypes := map[string]bool{"text/html": true, "text/css": true, "application/javascript": true, "application/json": true}

	valid, exists := validContentTypes[contentTypePart]

	return valid && exists
}
