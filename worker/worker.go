package worker

import (
	"github.com/google/uuid"
	"sync"
)

//NewWorker
func NewWorker(name string, handle func() error, replicas int) *Worker {
	id, _ := uuid.NewUUID()
	return &Worker{Id: id.String(), Name: name, Handle: handle, Replicas: replicas, subWorkers: make(map[string]*SubWorker)}
}

//Run
func (w *Worker) Run(errors chan WrapperHandleError) {
	var wg sync.WaitGroup
	for i := 1; i <= w.Replicas; i++ {
		s := &SubWorker{Id: i, Status: STARTED, Worker: w}
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

//Status
func (w *Worker) Status() map[string]int {
	status := map[string]int{}
	for _, subWorker := range w.subWorkers {
		status[subWorker.Name()] = subWorker.Status
	}
	return status
}
