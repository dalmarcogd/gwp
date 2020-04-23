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

type (
	//SubWorker
	SubWorker struct {
		Worker *Worker
		Id     int
		Status int
		Error  error
	}

	//Worker
	Worker struct {
		Id            string
		Name          string
		StartAt       time.Time
		FinishedAt    time.Time
		Handle        func() error
		Replicas      int
		RestartAlways bool
		Restarts      int
		subWorkers    map[string]*SubWorker
	}

	//WrapperHandleError
	WrapperHandleError struct {
		worker *Worker
		err    error
	}
)

func (s SubWorker) Name() string {
	return fmt.Sprintf("%s-%s", s.Worker.Name, strconv.Itoa(s.Id))
}
