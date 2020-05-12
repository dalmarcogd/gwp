package worker

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"strconv"
	"sync"
	"time"
)

//NewWorker is a constructor for #Worker and give
//for user some default settings
func NewWorker(name string, handle func() error, configs ...Config) *Worker {
	id, _ := uuid.NewUUID()
	w := &Worker{
		ID:          id.String(),
		Name:        name,
		Handle:      handle,
		Concurrency: 1,
		subWorkers:  make(map[string]*SubWorker),
		ctx:         context.Background(),
	}
	for _, config := range configs {
		config.Apply(w)
	}
	return w
}

//Run is a executed inside goroutine by #RunWorkers
//He administrate the number of concurrency
func (w *Worker) Run(errors chan WrapperHandleError) {
	var wg sync.WaitGroup
	w.StartAt = time.Now().UTC()
	for i := 1; i <= w.Concurrency; i++ {
		ctx, cancel := getContext(w)

		s := newSubWorker(ctx, i, w)
		w.subWorkers[s.Name()] = s

		wg.Add(1)
		go func(subWorker *SubWorker, c context.Context, cancelFunc context.CancelFunc) {
			defer wg.Done()
			defer cancelFunc()

			select {
			case <-c.Done():
				log.Printf("Worker [%v] finished: %v", subWorker.Name(), ctx.Err())
			case <-s.Run(errors):
				log.Printf("Worker [%v] finished", subWorker.Name())
			}
		}(s, ctx, cancel)
	}
	wg.Wait()
}

func newSubWorker(ctx context.Context, id int, w *Worker) *SubWorker {
	return &SubWorker{ID: id, Status: Started, Worker: w, ctx: ctx}
}

func getContext(w *Worker) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(w.ctx)
	if w.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, w.Timeout)
	}
	if !w.Deadline.IsZero() {
		ctx, cancel = context.WithDeadline(ctx, w.Deadline)
	}
	return ctx, cancel
}

//Status return a map with status from the each #SubWorker
func (w *Worker) Status() map[string]string {
	status := map[string]string{}
	for _, subWorker := range w.subWorkers {
		status[subWorker.Name()] = subWorker.Status
	}
	return status
}

//Healthy check if anyone #SubWorker still #Started,
//this survey responds if it is running
func (w *Worker) Healthy() bool {
	for _, v := range w.Status() {
		if v == Started {
			return true
		}
	}
	return false
}

//Name return the name of #SubWorker
//the pattern is %s-%s <- #Worker.Name, #SubWorker.ID
func (s SubWorker) Name() string {
	return fmt.Sprintf("%s-%s", s.Worker.Name, strconv.Itoa(s.ID))
}

//Run method execute on instance fo worker. If you have config two concurrences
//this will execute two twice
func (s *SubWorker) Run(errors chan WrapperHandleError) chan bool {
	done := make(chan bool, 1)

	go func(handle func() error, sw *SubWorker) {
		if err := handle(); err != nil {
			errors <- WrapperHandleError{subWorker: sw, err: err}
			sw.Status = Error
			sw.Error = err
		} else {
			sw.Status = Finished
		}
		done <- true
	}(s.Worker.Handle, s)

	return done
}

//RunWorkers is a function that administrate the workers and yours errors
func RunWorkers(workers []*Worker, handleError func(w *Worker, err error)) error {
	var wg sync.WaitGroup

	for _, worker := range workers {
		errors := make(chan WrapperHandleError, worker.Concurrency)

		wg.Add(2)
		go func(w *Worker) {
			defer wg.Done()
			runWorker(w, errors)
		}(worker)
		go func(w *Worker) {
			defer wg.Done()
			runWorkerHandleError(handleError, w, errors)
		}(worker)
	}
	//Waiting all goroutines
	wg.Wait()

	return nil
}

func runWorker(w *Worker, errors chan WrapperHandleError) {
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

func runWorkerHandleError(handleError func(w *Worker, err error), worker *Worker, errors chan WrapperHandleError) {
	defer log.Printf("Worker [%s] handleError finished", worker.Name)
	log.Printf("Worker [%s] handleError started", worker.Name)
	for err := range errors {
		if handleError != nil {
			done := make(chan bool, 1)
			go func() {
				handleError(err.subWorker.Worker, err.err)
				done <- true
			}()

			select {
			case <-time.After(10 * time.Second):
				log.Printf("Worker [%s] handleError timeout for handling error: %v", worker.Name, err.err)
				break
			case <-done:
				log.Printf("Worker [%s] handleError handled: %v", worker.Name, err.err)
				break
			}
		} else {
			log.Printf("Worker [%s] error [%v] ignored", worker.Name, err.err)
		}
	}
}
