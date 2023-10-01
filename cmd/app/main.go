package main

import (
	"flag"
	"github.com/anvlad11/testapp-20230930/internal/model"
	"github.com/anvlad11/testapp-20230930/internal/services/downloader"
	"github.com/anvlad11/testapp-20230930/internal/services/extractor"
	"github.com/anvlad11/testapp-20230930/internal/services/manager"
	"log"
	"time"
)

var (
	rootVar = flag.String("root", "", "Page to start crawling from")
)

func main() {
	flag.Parse()
	if rootVar == nil {
		log.Fatal("root is empty")
	}
	root := *rootVar

	managerService := manager.NewService()
	managerService.Start()

	for i := 0; i <= 3; i++ {
		downloaderService := downloader.NewService()
		downloaderService.Start()
		managerService.AddDownloader(downloaderService)
	}

	for i := 0; i <= 3; i++ {
		extractorService := extractor.NewService()
		extractorService.Start()
		managerService.AddExtractor(extractorService)
	}

	err := managerService.Process(&model.Task{URL: root})
	if err != nil {
		log.Fatalf("could not process root task: %v", err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
