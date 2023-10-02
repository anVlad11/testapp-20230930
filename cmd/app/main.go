package main

import (
	"flag"
	"github.com/anvlad11/testapp-20230930/internal/config"
	"github.com/anvlad11/testapp-20230930/internal/repositories/task"
	"github.com/anvlad11/testapp-20230930/internal/services/manager"
	"github.com/anvlad11/testapp-20230930/internal/services/worker"
	"github.com/anvlad11/testapp-20230930/pkg/database"
	"github.com/anvlad11/testapp-20230930/pkg/model"
	"log"
	"time"
)

var configPath = flag.String(
	"config-path",
	"./config.yaml",
	"Path to the application config",
)

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

	taskRepository := task.NewRepository(db)

	inputQueue := make(chan *model.Task, 1)
	managerService := manager.NewService(taskRepository, inputQueue)
	managerService.Start()

	for i := 1; i <= cfg.WorkerCount; i++ {
		downloaderService := worker.NewService(
			cfg.ContentTypes,
			cfg.RequestHeaders,
			cfg.DataDirectory,
			taskRepository,
		)
		downloaderService.Start()
		managerService.AddWorker(downloaderService)
	}

	inputQueue <- &model.Task{URL: cfg.RootURL}

	for {
		time.Sleep(1 * time.Second)
	}
}
