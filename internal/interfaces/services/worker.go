package services

import "github.com/anvlad11/testapp-20230930/pkg/model"

type WorkerService interface {
	SetInput(input chan *model.Task)
}
