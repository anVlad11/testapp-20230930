package main

import (
	"flag"
	"github.com/anvlad11/testapp-20230930/internal/repositories/task"
	"github.com/anvlad11/testapp-20230930/internal/services/downloader"
	"github.com/anvlad11/testapp-20230930/internal/services/extractor"
	"github.com/anvlad11/testapp-20230930/internal/services/manager"
	"github.com/anvlad11/testapp-20230930/pkg/database"
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"log"
	"time"
)

var (
	rootVar    = flag.String("root", "", "Page to start crawling from")
	dataDirVar = flag.String("data-dir", "./data", "Folder to store crawled pages content")
)

func main() {
	flag.Parse()
	if rootVar == nil {
		log.Fatal("root is empty")
	}
	root := *rootVar

	initTask := &model.Task{URL: root}

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	taskRepository := task.NewRepository(db, *dataDirVar)

	managerService := manager.NewService(taskRepository)
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

	err = managerService.Process(initTask)
	if err != nil {
		log.Fatalf("could not process root task: %v", err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
