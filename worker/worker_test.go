package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNewWorker(t *testing.T) {
	nameWorker := "w1"
	handleWorker := func(ctx context.Context) error { return nil }
	concurrencyWorker := 1
	w := NewWorker(nameWorker, handleWorker, WithConcurrency(concurrencyWorker))
	if w.Name != nameWorker {
		t.Errorf("Name of worker if different from setup %s != %s", w.Name, nameWorker)
	}
	if w.Handle != nil && fmt.Sprintf("%p", w.Handle) != fmt.Sprintf("%p", handleWorker) {
		t.Errorf("Name of worker if different from setup %p != %p", w.Handle, handleWorker)
	}
	if w.Concurrency != concurrencyWorker {
		t.Errorf("Concurrency of worker if different from setup %d != %d", w.Concurrency, concurrencyWorker)
	}
	if w.RestartAlways != false {
		t.Errorf("RestartAlaways of worker if different from setup %t != %t", w.RestartAlways, false)
	}
}

func TestWorker_Run(t *testing.T) {
	nameWorker := "w1"
	handleWorker := func(ctx context.Context) error {
		<-time.After(1 * time.Second)
		return errors.New("happened error")
	}
	concurrencyWorker := 1
	w := NewWorker(nameWorker, handleWorker, WithConcurrency(concurrencyWorker))
	go func() {
		errorsCh := make(chan WrapperHandleError, 1)
		w.Run(errorsCh)
		close(errorsCh)
	}()

	for _, v := range w.Status() {
		if v != Started {
			t.Errorf("Was expect that worker is with started status, but returned: %s", v)
		}
	}

	<-time.After(2 * time.Second)
	for _, v := range w.Status() {
		if v != Error {
			t.Errorf("Was expect that worker is with error status, but returned: %s", v)
		}
	}

}

func TestWorker_Status(t *testing.T) {
	nameWorker := "w1"
	handleWorker := func(ctx context.Context) error {
		<-time.After(3 * time.Second)
		return nil
	}
	concurrencyWorker := 1
	w := NewWorker(nameWorker, handleWorker, WithConcurrency(concurrencyWorker))
	go func() {
		errorsCh := make(chan WrapperHandleError, 1)
		w.Run(errorsCh)
		close(errorsCh)
	}()

	<-time.After(1 * time.Second)
	for _, v := range w.Status() {
		if v != Started {
			t.Errorf("Was expect that worker is with started status, but returned: %s", v)
		}
	}

	nameWorker = "w1"
	handleWorker = func(ctx context.Context) error {
		return errors.New("happened error")
	}
	concurrencyWorker = 1
	w = NewWorker(nameWorker, handleWorker, WithConcurrency(concurrencyWorker))
	go func() {
		errorsCh := make(chan WrapperHandleError, 1)
		w.Run(errorsCh)
		close(errorsCh)
	}()

	<-time.After(1 * time.Second)
	for _, v := range w.Status() {
		if v != Error {
			t.Errorf("Was expect that worker is with error status, but returned: %s", v)
		}
	}

	nameWorker = "w1"
	handleWorker = func(ctx context.Context) error {
		return nil
	}
	concurrencyWorker = 1
	w = NewWorker(nameWorker, handleWorker, WithConcurrency(concurrencyWorker))
	go func() {
		errorsCh := make(chan WrapperHandleError, 1)
		w.Run(errorsCh)
		close(errorsCh)
	}()

	<-time.After(1 * time.Second)
	for _, v := range w.Status() {
		if v != Finished {
			t.Errorf("Was expect that worker is with finished status, but returned: %s", v)
		}
	}
}

func TestWorker_Healthy(t *testing.T) {
	nameWorker := "w1"
	handleWorker := func(ctx context.Context) error {
		<-time.After(3 * time.Second)
		return nil
	}
	concurrencyWorker := 1
	w := NewWorker(nameWorker, handleWorker, WithConcurrency(concurrencyWorker))
	go func() {
		errs := make(chan WrapperHandleError, w.Concurrency)
		w.Run(errs)
		close(errs)
	}()
	<-time.After(1 * time.Second)
	if !w.Healthy() {
		t.Errorf("Was expect that worker is Up, but returned: Down")
	}
	<-time.After(3 * time.Second)
	if w.Healthy() {
		t.Errorf("Was expect that worker is Down, but returned: Up")
	}
}

func TestRunWorkers(t *testing.T) {
	handleErrors := func(w *Worker, err error) {
		log.Print(err)
	}

	workers := []*Worker{
		NewWorker("w1",
			func(ctx context.Context) error {
				select {
				case <-ctx.Done():
				case <-time.After(1 * time.Second):
				}

				return nil
			},
			WithRestartAlways(),
			WithTimeout(3*time.Second),
			WithDeadline(time.Now().Add(5*time.Second)),
		),
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

func TestRunWorkersCron(t *testing.T) {
	handleErrors := func(w *Worker, err error) {
		log.Print(err)
	}

	counter := 0
	workers := []*Worker{
		NewWorker("w1",
			func(ctx context.Context) error {
				select {
				case <-ctx.Done():
				default:
					counter++
					break
				}

				return nil
			},
			WithCron(time.Second),
			WithTimeout(2*time.Second),
		),
	}

	go func() {
		if err := RunWorkers(workers, handleErrors); err != nil {
			t.Error(err)
		}
	}()

	<-time.After(4 * time.Second)
	if counter < 3 {
		t.Error("Worker setup to execute every second by cron but not executed")
	}
}

func Test_runWorkerHandleError(t *testing.T) {
	hasErrorHandled := false
	handleErrors := func(w *Worker, err error) {
		log.Print(err)
		hasErrorHandled = true
	}

	workers := []*Worker{
		NewWorker("w1",
			func(ctx context.Context) error {
				<-time.After(2 * time.Second)
				return errors.New("happened some error")
			}),
	}

	go func() {
		if err := RunWorkers(workers, handleErrors); err != nil {
			t.Error(err)
		}
	}()

	<-time.After(3 * time.Second)
	for _, worker := range workers {
		if worker.Status()["w1-1"] != Error {
			t.Error("Worker setup to return error but not returned")
		}
	}
	if !hasErrorHandled {
		t.Error("Expected error handled")
	}

	hasErrorUnhandled := false
	handleErrors = func(w *Worker, err error) {
		<-time.After(11 * time.Second)
		hasErrorUnhandled = true
	}

	workers = []*Worker{
		NewWorker("w2",
			func(ctx context.Context) error {
				return errors.New("happened some error")
			}),
	}

	go func() {
		if err := RunWorkers(workers, handleErrors); err != nil {
			t.Error(err)
		}
	}()

	<-time.After(10 * time.Second)
	for _, worker := range workers {
		if worker.Status()["w2-1"] != Error {
			t.Error("Worker setup to return error but not returned")
		}
	}
	if hasErrorUnhandled {
		t.Error("Expected no error handled")
	}

	workers = []*Worker{
		NewWorker("w3",
			func(ctx context.Context) error {
				return errors.New("happened some error")
			}),
	}

	go func() {
		if err := RunWorkers(workers, nil); err != nil {
			t.Error(err)
		}
	}()

	<-time.After(1 * time.Second)
	for _, worker := range workers {
		if worker.Status()["w3-1"] != Error {
			t.Error("Worker setup to return error but not returned")
		}
	}
}
