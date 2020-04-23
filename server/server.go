package server

import (
	"github.com/dalmarcogd/go-worker-pool/monitoring"
	"github.com/dalmarcogd/go-worker-pool/runtime"
	"github.com/dalmarcogd/go-worker-pool/worker"
	"log"
)

var (
	DefaultConfig = map[string]interface{}{
		"port":        8001,
		"host":        "localhost",
		"basePath":    "/workers",
		"stats":       false,
		"healthCheck": false,
		"debugPprof":  false,
	}
)

type server struct {
	port        int
	host        string
	basePath    string
	stats       bool
	healthCheck bool
	debugPprof  bool
	workers     []*worker.Worker
	handleError func(w *worker.Worker, err error)
}

//New
func New() *server {
	return NewWithConfig(DefaultConfig)
}

//NewWithConfig
func NewWithConfig(configs map[string]interface{}) *server {
	port := 8001
	if p, ok := configs["port"]; ok {
		port = p.(int)
	}
	host := "localhost"
	if h, ok := configs["host"]; ok {
		host = h.(string)
	}

	basePath := "/workers"
	if bp, ok := configs["basePath"]; ok {
		basePath = bp.(string)
	}

	stats := false
	if s, ok := configs["stats"]; ok {
		stats = s.(bool)
	}

	healthCheck := false
	if hc, ok := configs["healthCheck"]; ok {
		healthCheck = hc.(bool)
	}

	debugPprof := false
	if dpp, ok := configs["debugPprof"]; ok {
		debugPprof = dpp.(bool)
	}

	s := &server{
		port:        port,
		host:        host,
		basePath:    basePath,
		stats:       stats,
		healthCheck: healthCheck,
		debugPprof:  debugPprof,
		workers:     []*worker.Worker{},
	}
	runtime.SetServerRun(s)
	return s
}

//Stats
func (s *server) Stats() *server {
	s.stats = true
	return s
}

//HealthCheck
func (s *server) HealthCheck() *server {
	s.healthCheck = true
	return s
}

//DebugPprof
func (s *server) DebugPprof() *server {
	s.debugPprof = true
	return s
}

//HandleError
func (s *server) HandleError(handle func(w *worker.Worker, err error)) *server {
	s.handleError = handle
	return s
}

//Worker
func (s *server) Worker(name string, handle func() error, concurrency int, restartAlways bool) *server {
	s.workers = append(s.workers, worker.NewWorker(name, handle, concurrency, restartAlways))
	return s
}

//Run
func (s *server) Run() error {
	monitoring.SetupHttp(map[string]interface{}{
		"port":        s.port,
		"host":        s.host,
		"stats":       s.stats,
		"healthCheck": s.healthCheck,
		"debugPprof":  s.debugPprof,
		"basePath":    s.basePath,
	})
	defer func() {
		if err := monitoring.CloseHttp(); err != nil {
			log.Printf("Error when closed monitoring server at: %s", err)
		}
	}()
	return worker.RunWorkers(s.workers, s.handleError)
}

//Workers
func (s *server) Workers() []*worker.Worker {
	return s.workers
}
