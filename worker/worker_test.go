package worker

import (
	"errors"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNewWorker(t *testing.T) {
	nameWorker := "w1"
	handleWorker := func() error { return nil }
	concurrencyWorker := 1
	restartAlwaysWorker := false
	w := NewWorker(nameWorker, handleWorker, concurrencyWorker, restartAlwaysWorker)
	if w.Name != nameWorker {
		t.Errorf("Name of worker if different from setup %s != %s", w.Name, nameWorker)
	}
	if w.Handle != nil && fmt.Sprintf("%p", w.Handle) != fmt.Sprintf("%p", handleWorker) {
		t.Errorf("Name of worker if different from setup %p != %p", w.Handle, handleWorker)
	}
	if w.Concurrency != concurrencyWorker {
		t.Errorf("Concurrency of worker if different from setup %d != %d", w.Concurrency, concurrencyWorker)
	}
	if w.RestartAlways != restartAlwaysWorker {
		t.Errorf("RestartAlaways of worker if different from setup %t != %t", w.RestartAlways, restartAlwaysWorker)
	}
}

func TestWorker_Run(t *testing.T) {
	nameWorker := "w1"
	handleWorker := func() error {
		<-time.After(1 * time.Second)
		return errors.New("happened error")
	}
	concurrencyWorker := 1
	restartAlwaysWorker := false
	w := NewWorker(nameWorker, handleWorker, concurrencyWorker, restartAlwaysWorker)
	go func() {
		errorsCh := make(chan WrapperHandleError, 1)
		w.Run(errorsCh)
		close(errorsCh)
	}()

	for _, v := range w.Status() {
		if v != STARTED {
			t.Errorf("Was expect that worker is with started status, but returned: %s", v)
		}
	}

	<-time.After(2 * time.Second)
	for _, v := range w.Status() {
		if v != ERROR {
			t.Errorf("Was expect that worker is with error status, but returned: %s", v)
		}
	}

}

func TestWorker_Status(t *testing.T) {
	nameWorker := "w1"
	handleWorker := func() error {
		<-time.After(3 * time.Second)
		return nil
	}
	concurrencyWorker := 1
	restartAlwaysWorker := false
	w := NewWorker(nameWorker, handleWorker, concurrencyWorker, restartAlwaysWorker)
	go func() {
		errorsCh := make(chan WrapperHandleError, 1)
		w.Run(errorsCh)
		close(errorsCh)
	}()

	<-time.After(1 * time.Second)
	for _, v := range w.Status() {
		if v != STARTED {
			t.Errorf("Was expect that worker is with started status, but returned: %s", v)
		}
	}

	nameWorker = "w1"
	handleWorker = func() error {
		return errors.New("happened error")
	}
	concurrencyWorker = 1
	restartAlwaysWorker = false
	w = NewWorker(nameWorker, handleWorker, concurrencyWorker, restartAlwaysWorker)
	go func() {
		errorsCh := make(chan WrapperHandleError, 1)
		w.Run(errorsCh)
		close(errorsCh)
	}()

	<-time.After(1 * time.Second)
	for _, v := range w.Status() {
		if v != ERROR {
			t.Errorf("Was expect that worker is with error status, but returned: %s", v)
		}
	}

	nameWorker = "w1"
	handleWorker = func() error {
		return nil
	}
	concurrencyWorker = 1
	restartAlwaysWorker = false
	w = NewWorker(nameWorker, handleWorker, concurrencyWorker, restartAlwaysWorker)
	go func() {
		errorsCh := make(chan WrapperHandleError, 1)
		w.Run(errorsCh)
		close(errorsCh)
	}()

	<-time.After(1 * time.Second)
	for _, v := range w.Status() {
		if v != FINISHED {
			t.Errorf("Was expect that worker is with finished status, but returned: %s", v)
		}
	}
}

func TestWorker_IsUp(t *testing.T) {
	nameWorker := "w1"
	handleWorker := func() error {
		<-time.After(3 * time.Second)
		return nil
	}
	concurrencyWorker := 1
	restartAlwaysWorker := false
	w := NewWorker(nameWorker, handleWorker, concurrencyWorker, restartAlwaysWorker)
	go func() {
		errors := make(chan WrapperHandleError)
		w.Run(errors)
		close(errors)
	}()
	<-time.After(1 * time.Second)
	if !w.IsUp() {
		t.Errorf("Was expect that worker is Up, but returned: Down")
	}
	<-time.After(3 * time.Second)
	if w.IsUp() {
		t.Errorf("Was expect that worker is Down, but returned: Up")
	}
}

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
