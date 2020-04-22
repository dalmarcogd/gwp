package worker

import (
	"log"
	"sync"
)

func RunWorkers(workers []*Worker, handleError func(w *Worker, err error)) error {
	var wg sync.WaitGroup

	for _, worker := range workers {
		log.Printf("Worker [%s] started", worker.Name)

		wg.Add(1)
		errors := make(chan WrapperHandleError)
		go func(w *Worker, errors chan WrapperHandleError) {
			defer wg.Done()
			defer log.Printf("Worker [%s] finished", w.Name)
			defer close(errors)
			w.Run(errors)
		}(worker, errors)

		wg.Add(1)
		go func(worker *Worker, handle func(w *Worker, err error), errors chan WrapperHandleError, wg *sync.WaitGroup) {
			defer wg.Done()
			defer log.Printf("Worker [%s] handleError finished", worker.Name)
			log.Printf("Worker [%s] handleError started", worker.Name)
			for err := range errors {
				if handleError != nil {
					handleError(err.worker, err.err)
				}
			}
		}(worker, handleError, errors, &wg)
	}

	wg.Wait()

	return nil
}
