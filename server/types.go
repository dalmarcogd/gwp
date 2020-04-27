package server

import (
	"github.com/dalmarcogd/go-worker-pool/worker"
)

type workerServer struct {
	config      map[string]interface{}
	workers     map[string]*worker.Worker
	handleError func(w *worker.Worker, err error)
}
