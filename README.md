![Go](https://github.com/dalmarcogd/go-worker-pool/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/dalmarcogd/go-worker-pool/branch/master/graph/badge.svg)](https://codecov.io/gh/dalmarcogd/go-worker-pool)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/beee52f22195471abea544a19ee6304a)](https://www.codacy.com/manual/dalmarco.gd/go-worker-pool?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=dalmarcogd/go-worker-pool&amp;utm_campaign=Badge_Grade)

# go-worker-pool

This package wants to offer the community and implement workers with the pure Go code for Golangers, without any other dependency just Uuid. It allows you to expose an http server to answer the response of health checks, stats, debug pprof and the main "workers". Workers for consumer queues, channel processes and other things that you think worker needs.

# Prerequisites
Golang version >= [1.14](https://golang.org/doc/devel/release.html#go1.14)

# Features
- Setup http server to monitoring yours;
  - /stats with workers, showing statuses her, number of goroutines, number of cpus and more;
  - /health-check that look for status of workers;
  - /debug/pprof expose all endpoints of investivate golang runtime [http](https://golang.org/pkg/net/http/pprof/);
- Allow multiple concurrencies of work, handle errors and restart always worker;

# Documentation
TODO

# Get started
#### [Simple Usage](examples/simpleWorker.go) ###

```go
package main

import (
	"errors"
	"github.com/dalmarcogd/go-worker-pool/server"
	"github.com/dalmarcogd/go-worker-pool/worker"
	"log"
	"time"
)

func main() {
	if err := server.
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
				return errors.New("teste")
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
				return errors.New("teste")
			},
			1,
			false).
		//Worker(
		//	"w4",
		//	func() error {
		//		time.Sleep(1000)
		//		return nil
		//	},
		//	1).
		Run(); err != nil {
		panic(err)
	}
}
```
