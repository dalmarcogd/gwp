package worker

import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestRunWorkers(t *testing.T) {
	handleErrors := func(w *Worker, err error) {
		log.Print(err)
	}

	workers := []*Worker{
		NewWorker("w1",
			func() error {
				<-time.After(1 * time.Second)
				return nil
			},
			1,
			true),
	}

	go func() {
		if err := RunWorkers(workers, handleErrors); err != nil {
			t.Error(err)
		}
	}()

	<-time.After(2 * time.Second)
	for _, worker := range workers {
		if worker.Restarts <= 0 {
			t.Error("Worker setup to restart always but not restarted")
		}
	}
}

func Test_runWorkerHandleError(t *testing.T) {
	hasError := false
	handleErrors := func(w *Worker, err error) {
		log.Print(err)
		hasError = true
	}

	workers := []*Worker{
		NewWorker("w1",
			func() error {
				<-time.After(1 * time.Second)
				return errors.New("happened some error")
			},
			1,
			false),
	}

	go func() {
		if err := RunWorkers(workers, handleErrors); err != nil {
			t.Error(err)
		}
	}()

	<-time.After(3 * time.Second)
	for _, worker := range workers {
		if worker.Status()["w1-1"] != ERROR {
			t.Error("Worker setup to return error but not returned")
		}
	}
	if !hasError {
		t.Error("Expected error handled")
	}
}
