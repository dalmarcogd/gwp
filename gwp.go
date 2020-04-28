package gwp

import (
	"github.com/dalmarcogd/gwp/monitor"
	"github.com/dalmarcogd/gwp/runtime"
	"github.com/dalmarcogd/gwp/worker"
	"log"
	"net/http"
)

type WorkerServer struct {
	config      map[string]interface{}
	workers     map[string]*worker.Worker
	handleError func(w *worker.Worker, err error)
}

var (
	//defaultConfig is a default config for start the #WorkerServer
	defaultConfig = map[string]interface{}{
		"port":        8001,
		"host":        "localhost",
		"basePath":    "/workers",
		"stats":       false,
		"healthCheck": false,
		"debugPprof":  false,
	}
)

// New build an #WorkerServer with #defaultConfig
func New() *WorkerServer {
	return NewWithConfig(defaultConfig)
}

// NewWithConfig build an #WorkerServer by the settings
func NewWithConfig(configs map[string]interface{}) *WorkerServer {
	for k, v := range defaultConfig {
		if _, ok := configs[k]; !ok {
			configs[k] = v
		}
	}

	s := &WorkerServer{
		config:  configs,
		workers: map[string]*worker.Worker{},
	}
	runtime.SetServerRun(s)
	return s
}

// Stats setup for the server to start with /stats
func (s *WorkerServer) Stats() *WorkerServer {
	s.config["stats"] = true
	return s
}

// StatsFunc setup the handler for /stats
func (s *WorkerServer) StatsFunc(f func(writer http.ResponseWriter, request *http.Request)) *WorkerServer {
	s.Stats().config["statsFunc"] = f
	return s
}

// HealthCheck setup for the server to start with /health-check
func (s *WorkerServer) HealthCheck() *WorkerServer {
	s.config["healthCheck"] = true
	return s
}

// HealthCheckFunc setup the handler for /health-check
func (s *WorkerServer) HealthCheckFunc(f func(writer http.ResponseWriter, request *http.Request)) *WorkerServer {
	s.HealthCheck().config["healthCheckFunc"] = f
	return s
}

// DebugPprof setup for the server to start with /debug/pprof*
func (s *WorkerServer) DebugPprof() *WorkerServer {
	s.config["debugPprof"] = true
	return s
}

// HandleError setup the a function that will called when to occur and error
func (s *WorkerServer) HandleError(handle func(w *worker.Worker, err error)) *WorkerServer {
	s.handleError = handle
	return s
}

// Worker build an #Worker and add to execution with #WorkerServer
func (s *WorkerServer) Worker(name string, handle func() error, concurrency int, restartAlways bool) *WorkerServer {
	w := worker.NewWorker(name, handle, concurrency, restartAlways)
	s.workers[w.ID] = w
	return s
}

// Workers return the slice of #Worker configured
func (s *WorkerServer) Workers() []*worker.Worker {
	v := make([]*worker.Worker, 0, len(s.workers))

	for _, value := range s.workers {
		v = append(v, value)
	}
	return v
}

// Configs return the configs from #WorkerServer
func (s *WorkerServer) Configs() map[string]interface{} {
	return s.config
}

// Run user to start the #WorkerServer
func (s *WorkerServer) Run() error {
	monitor.SetupHTTP(s.config)
	defer func() {
		if err := monitor.CloseHTTP(); err != nil {
			log.Printf("Error when closed monitor WorkerServer at: %s", err)
		}
	}()
	return worker.RunWorkers(s.Workers(), s.handleError)
}
