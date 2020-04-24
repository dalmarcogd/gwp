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

//New
func New() *workerServer {
	return NewWithConfig(DefaultConfig)
}

//NewWithConfig
func NewWithConfig(configs map[string]interface{}) *workerServer {
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

	s := &workerServer{
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
func (s *workerServer) Stats() *workerServer {
	s.stats = true
	return s
}

//HealthCheck
func (s *workerServer) HealthCheck() *workerServer {
	s.healthCheck = true
	return s
}

//DebugPprof
func (s *workerServer) DebugPprof() *workerServer {
	s.debugPprof = true
	return s
}

//HandleError
func (s *workerServer) HandleError(handle func(w *worker.Worker, err error)) *workerServer {
	s.handleError = handle
	return s
}

//Worker
func (s *workerServer) Worker(name string, handle func() error, concurrency int, restartAlways bool) *workerServer {
	s.workers = append(s.workers, worker.NewWorker(name, handle, concurrency, restartAlways))
	return s
}

//Run
func (s *workerServer) Run() error {
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
			log.Printf("Error when closed monitoring workerServer at: %s", err)
		}
	}()
	return worker.RunWorkers(s.workers, s.handleError)
}
