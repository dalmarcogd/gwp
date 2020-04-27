package worker

import (
	"log"
	"sync"
	"time"
)

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
