package worker

import (
	"context"
	"time"
)

const (
	// STARTED is a value used for control which are running
	STARTED = "Started"
	// FINISHED is a value used for control which are finished
	FINISHED = "Finished"
	// ERROR is a value used for control if has error
	ERROR = "Error"
)

type (
	// SubWorker is a type that represents the concurrency, for the number of concurrency
	// has an #SubWorker
	SubWorker struct {
		Worker *Worker
		ID     int
		Status string
		Error  error
		ctx    context.Context
	}

	// Worker is a type that represents an group of concurrency and keep some settings
	Worker struct {
		ID            string
		Name          string
		StartAt       time.Time
		FinishedAt    time.Time
		Handle        func() error
		Concurrency   int
		RestartAlways bool
		Restarts      int
		Timeout       time.Duration
		Deadline      time.Time
		subWorkers    map[string]*SubWorker
		ctx           context.Context
	}

	// WrapperHandleError is a wrapper to transport worker and the error generate inside worker
	WrapperHandleError struct {
		worker *Worker
		err    error
	}

	Config struct {
		k func(w *Worker)
	}
)
