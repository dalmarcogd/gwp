package worker

import (
	"fmt"
	"strconv"
	"time"
)

const (
	STARTED  = 1
	FINISHED = 2
	ERROR    = 3
)

type SubWorker struct {
	Worker *Worker
	Id     int
	Status int
	Error  error
}

func (s SubWorker) Name() string {
	return fmt.Sprintf("%s-%s", s.Worker.Name, strconv.Itoa(s.Id))
}

type Worker struct {
	Id         string
	Name       string
	StartAt    time.Time
	FinishedAt time.Time
	Handle     func() error
	Replicas   int
	Errors     chan error
	subWorkers map[string]*SubWorker
}

type WrapperHandleError struct {
	worker *Worker
	err    error
}
