package gwp

import (
	"github.com/dalmarcogd/gwp/monitor"
	"github.com/dalmarcogd/gwp/runtime"
	"github.com/dalmarcogd/gwp/worker"
	"log"
	"net/http"
)

//WorkerServer is a server that administrate the workers and the monitor
type WorkerServer struct {
	config      map[string]interface{}
	workers     map[string]*worker.Worker
	handleError func(w *worker.Worker, err error)
	healthy     []func() bool
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
		healthy: []func() bool {},
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

// CheckHealth includes to server checker the health
func (s *WorkerServer) CheckHealth(check func() bool) *WorkerServer {
	s.healthy = append(s.healthy, check)
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
	return s.CheckHealth(func() bool {
		return w.Healthy()
	})
}

// Workers return the slice of #Worker configured
func (s *WorkerServer) Workers() []*worker.Worker {
	v := make([]*worker.Worker, 0, len(s.workers))
	for _, value := range s.workers {
		v = append(v, value)
	}
	return v
}

// Workers return the slice of #Worker configured
func (s *WorkerServer) Healthy() bool {
	status := true
	for _, healthy := range s.healthy {
		if !healthy() {
			status = false
			break
		}
	}
	return status
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
