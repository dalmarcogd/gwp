package server

import "github.com/dalmarcogd/go-worker-pool/worker"

type workerServer struct {
	port        int
	host        string
	basePath    string
	stats       bool
	healthCheck bool
	debugPprof  bool
	workers     []*worker.Worker
	handleError func(w *worker.Worker, err error)
}

//Workers
func (s *workerServer) Workers() []*worker.Worker {
	return s.workers
}
