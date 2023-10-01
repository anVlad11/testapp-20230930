package services

import "github.com/anvlad11/testapp-20230930/internal/model"

type WorkerService interface {
	SetPipe(input chan *model.Task, output chan *model.Task)
	Start()
	Stop()
}
