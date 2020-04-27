package worker

import (
	"time"
)

const (
	//STARTED
	STARTED  = "Started"
	//FINISHED
	FINISHED = "Finished"
	//ERROR
	ERROR    = "Error"
)

type (
	//SubWorker
	SubWorker struct {
		Worker *Worker
		ID     int
		Status string
		Error  error
	}

	//Worker
	Worker struct {
		ID            string
		Name          string
		StartAt       time.Time
		FinishedAt    time.Time
		Handle        func() error
		Concurrency   int
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
