package gwp_test

import (
	"errors"
	"github.com/dalmarcogd/gwp"
	"github.com/dalmarcogd/gwp/worker"
	"log"
	"time"
)

func Example_Simple_Worker() {
	if err := gwp.
		New().
		Stats().
		HealthCheck().
		DebugPprof().
		HandleError(func(w *worker.Worker, err error) {
			log.Printf("Worker [%s] error: %s", w.Name, err)
		}).
		Worker(
			"w1",
			func() error {
				time.Sleep(10 * time.Second)
				return errors.New("test")
			},
			1,
			true).
		Worker(
			"w2",
			func() error {
				time.Sleep(30 * time.Second)
				return nil
			},
			1,
			false).
		Worker(
			"w3",
			func() error {
				time.Sleep(1 * time.Minute)
				return errors.New("test")
			},
			1,
			false).
		Run(); err != nil {
		panic(err)
	}
}
