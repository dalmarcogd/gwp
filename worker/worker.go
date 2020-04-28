package worker

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"strconv"
	"sync"
	"time"
)

// NewWorker is a constructor for #Worker and give
// for user some default settings
func NewWorker(name string, handle func() error, concurrency int, restartAlways bool) *Worker {
	id, _ := uuid.NewUUID()
	return &Worker{
		ID:            id.String(),
		Name:          name,
		Handle:        handle,
		Concurrency:   concurrency,
		RestartAlways: restartAlways,
		subWorkers:    make(map[string]*SubWorker),
		Restarts:      0,
	}
}

// Run is a executed inside goroutine by #RunWorkers
// He administrate the number of concurrency
func (w *Worker) Run(errors chan WrapperHandleError) {
	var wg sync.WaitGroup
	w.StartAt = time.Now().UTC()
	for i := 1; i <= w.Concurrency; i++ {
		s := &SubWorker{ID: i, Status: STARTED, Worker: w}
		w.subWorkers[s.Name()] = s

		wg.Add(1)
		go func(handle func() error, sw *SubWorker) {
			defer wg.Done()
			if err := handle(); err != nil {
				errors <- WrapperHandleError{worker: w, err: err}
				sw.Status = ERROR
				sw.Error = err
			} else {
				sw.Status = FINISHED
			}
		}(w.Handle, s)
	}
	wg.Wait()
}

// Status return a map with status from the each #SubWorker
func (w *Worker) Status() map[string]string {
	status := map[string]string{}
	for _, subWorker := range w.subWorkers {
		status[subWorker.Name()] = subWorker.Status
	}
	return status
}

// IsUp check if anyone #SubWorker still #STARTED,
// this survey responds if it is running
func (w *Worker) IsUp() bool {
	for _, v := range w.Status() {
		if v == STARTED {
			return true
		}
	}
	return false
}

// Name return the name of #SubWorker
// the pattern is %s-%s <- #Worker.Name, #SubWorker.ID
func (s SubWorker) Name() string {
	return fmt.Sprintf("%s-%s", s.Worker.Name, strconv.Itoa(s.ID))
}

// RunWorkers is a function that administrate the workers and yours errors
func RunWorkers(workers []*Worker, handleError func(w *Worker, err error)) error {
	var wg sync.WaitGroup

	for _, worker := range workers {
		errors := make(chan WrapperHandleError)

		wg.Add(1)
		go runWorker(&wg)(worker, errors)

		wg.Add(1)
		go runWorkerHandleError(handleError)(worker, handleError, errors, &wg)
	}
	// Waiting all goroutines
	wg.Wait()

	return nil
}

func runWorker(wg *sync.WaitGroup) func(w *Worker, errors chan WrapperHandleError) {
	return func(w *Worker, errors chan WrapperHandleError) {
		defer wg.Done()
		log.Printf("Worker [%s] started", w.Name)
		defer log.Printf("Worker [%s] finished", w.Name)
		defer close(errors)
		for {
			w.Run(errors)
			if !w.RestartAlways {
				w.FinishedAt = time.Now().UTC()
				break
			}
			w.Restarts++
			log.Printf("Worker [%s] restarted", w.Name)
		}
	}
}

func runWorkerHandleError(handleError func(w *Worker, err error)) func(worker *Worker, handle func(w *Worker, err error), errors chan WrapperHandleError, wg *sync.WaitGroup) {
	return func(worker *Worker, handle func(w *Worker, err error), errors chan WrapperHandleError, wg *sync.WaitGroup) {
		defer wg.Done()
		defer log.Printf("Worker [%s] handleError finished", worker.Name)
		log.Printf("Worker [%s] handleError started", worker.Name)
		for err := range errors {
			if handleError != nil {
				handleError(err.worker, err.err)
			}
		}
	}
}
