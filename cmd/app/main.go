package main

import (
	"flag"
	"github.com/anvlad11/testapp-20230930/internal/config"
	"github.com/anvlad11/testapp-20230930/internal/repositories/task"
	"github.com/anvlad11/testapp-20230930/internal/services/downloader"
	"github.com/anvlad11/testapp-20230930/internal/services/extractor"
	"github.com/anvlad11/testapp-20230930/internal/services/manager"
	"github.com/anvlad11/testapp-20230930/pkg/database"
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"log"
	"time"
)

var configPath = flag.String("config-path", "./config.yaml", "Path to the application config")

func main() {
	flag.Parse()

	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	db, err := database.NewDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	taskRepository := task.NewRepository(db, cfg.DataDirectory, cfg.ContentTypes)

	managerService := manager.NewService(taskRepository)
	managerService.Start()

	for i := 1; i <= cfg.DownloaderCount; i++ {
		downloaderService := downloader.NewService(cfg.ContentTypes)
		downloaderService.Start()
		managerService.AddDownloader(downloaderService)
	}

	for i := 1; i <= cfg.ExtractorCount; i++ {
		extractorService := extractor.NewService()
		extractorService.Start()
		managerService.AddExtractor(extractorService)
	}

	initTask := &model.Task{URL: cfg.RootURL}

	err = managerService.Process(initTask)
	if err != nil {
		log.Fatalf("could not process root task: %v", err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
