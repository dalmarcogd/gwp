package server

import (
	"github.com/dalmarcogd/go-worker-pool/monitoring"
	"github.com/dalmarcogd/go-worker-pool/runtime"
	"github.com/dalmarcogd/go-worker-pool/worker"
	"log"
	"net/http"
)

var (
	//DefaultConfig is a default config for start the #workerServer
	DefaultConfig = map[string]interface{}{
		"port":        8001,
		"host":        "localhost",
		"basePath":    "/workers",
		"stats":       false,
		"healthCheck": false,
		"debugPprof":  false,
	}
)

// New build an #workerServer with #DefaultConfig
func New() *workerServer {
	return NewWithConfig(DefaultConfig)
}

// NewWithConfig build an #workerServer by the settings
func NewWithConfig(configs map[string]interface{}) *workerServer {
	s := &workerServer{
		config:  configs,
		workers: map[string]*worker.Worker{},
	}
	runtime.SetServerRun(s)
	return s
}

// Stats setup for the server to start with /stats
func (s *workerServer) Stats() *workerServer {
	s.config["stats"] = true
	return s
}

// StatsFunc setup the handler for /stats
func (s *workerServer) StatsFunc(f func(writer http.ResponseWriter, request *http.Request)) *workerServer {
	s.Stats().config["statsFunc"] = f
	return s
}

// HealthCheck setup for the server to start with /health-check
func (s *workerServer) HealthCheck() *workerServer {
	s.config["healthCheck"] = true
	return s
}

// HealthCheckFunc setup the handler for /health-check
func (s *workerServer) HealthCheckFunc(f func(writer http.ResponseWriter, request *http.Request)) *workerServer {
	s.HealthCheck().config["healthCheckFunc"] = f
	return s
}

// DebugPprof setup for the server to start with /debug/pprof*
func (s *workerServer) DebugPprof() *workerServer {
	s.config["debugPprof"] = true
	return s
}

// HandleError setup the a function that will called when to occur and error
func (s *workerServer) HandleError(handle func(w *worker.Worker, err error)) *workerServer {
	s.handleError = handle
	return s
}

// Worker build an #Worker and add to execution with #workerServer
func (s *workerServer) Worker(name string, handle func() error, concurrency int, restartAlways bool) *workerServer {
	w := worker.NewWorker(name, handle, concurrency, restartAlways)
	s.workers[w.ID] = w
	return s
}

// Workers return the slice of #Worker configured
func (s *workerServer) Workers() []*worker.Worker {
	v := make([]*worker.Worker, 0, len(s.workers))

	for _, value := range s.workers {
		v = append(v, value)
	}
	return v
}

// Configs return the configs from #workerServer
func (s *workerServer) Configs() map[string]interface{} {
	return s.config
}

// Run user to start the #workerServer
func (s *workerServer) Run() error {
	monitoring.SetupHTTP(s.config)
	defer func() {
		if err := monitoring.CloseHTTP(); err != nil {
			log.Printf("Error when closed monitoring workerServer at: %s", err)
		}
	}()
	return worker.RunWorkers(s.Workers(), s.handleError)
}
