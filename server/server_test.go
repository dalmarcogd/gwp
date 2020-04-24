package server

import (
	"fmt"
	"github.com/dalmarcogd/go-worker-pool/worker"
	"testing"
)

func TestNew(t *testing.T) {
	s := New()
	if s.port != DefaultConfig["port"] {
		t.Errorf("Port is different of default port %d != %d", s.port, DefaultConfig["port"])
	}
	if s.host != DefaultConfig["host"] {
		t.Errorf("Host is different of default host %s != %s", s.host, DefaultConfig["host"])
	}
	if s.basePath != DefaultConfig["basePath"] {
		t.Errorf("BasePath is different of default basePath %s != %s", s.basePath, DefaultConfig["basePath"])
	}
	if s.stats != DefaultConfig["stats"] {
		t.Errorf("Stats is different of default stats %t != %t", s.stats, DefaultConfig["stats"])
	}
	if s.healthCheck != DefaultConfig["healthCheck"] {
		t.Errorf("HealthCheck is different of default healthCheck %t != %t", s.healthCheck, DefaultConfig["healthCheck"])
	}
	if s.debugPprof != DefaultConfig["debugPprof"] {
		t.Errorf("DebugPprof is different of default debugPprof %t != %t", s.debugPprof, DefaultConfig["debugPprof"])
	}
}

func TestNewWithConfig(t *testing.T) {
	config := map[string]interface{}{
		"port":        8002,
		"host":        "google",
		"basePath":    "/workers/test",
		"stats":       true,
		"healthCheck": true,
		"debugPprof":  true,
	}
	s := NewWithConfig(config)
	if s.port != config["port"] {
		t.Errorf("Port is different of default port %d != %d", s.port, config["port"])
	}
	if s.host != config["host"] {
		t.Errorf("Host is different of default host %s != %s", s.host, config["host"])
	}
	if s.basePath != config["basePath"] {
		t.Errorf("BasePath is different of default basePath %s != %s", s.basePath, config["basePath"])
	}
	if s.stats != config["stats"] {
		t.Errorf("Stats is different of default stats %t != %t", s.stats, config["stats"])
	}
	if s.healthCheck != config["healthCheck"] {
		t.Errorf("HealthCheck is different of default healthCheck %t != %t", s.healthCheck, config["healthCheck"])
	}
	if s.debugPprof != config["debugPprof"] {
		t.Errorf("DebugPprof is different of default debugPprof %t != %t", s.debugPprof, config["debugPprof"])
	}
}

func Test_server_HandleError(t *testing.T) {
	f := func(w *worker.Worker, err error) {}
	s := New().HandleError(f)
	if s.handleError != nil && fmt.Sprintf("%p", f) != fmt.Sprintf("%p", s.handleError) {
		t.Errorf("HandleError is different of f %v != %v", fmt.Sprintf("%p", f), fmt.Sprintf("%p", s.handleError))
	}
}

func Test_server_HealthCheck(t *testing.T) {
	s := New().HealthCheck()
	if !s.healthCheck {
		t.Error("HealthCheck setup on workerServer and his not enable")
	}
}

func Test_server_Stats(t *testing.T) {
	s := New().Stats()
	if !s.stats {
		t.Error("Stats setup on workerServer and his not enable")
	}
}

func Test_server_DebugPprof(t *testing.T) {
	s := New().DebugPprof()
	if !s.debugPprof {
		t.Error("DebugPprof setup on workerServer and his not enable")
	}
}

func Test_server_Run(t *testing.T) {
	s := New().Worker("w1", func() error { return nil }, 1, false)
	if err := s.Run(); err != nil {
		t.Errorf("Error when run workerServer %v", err)
	}
	s = New().HealthCheck().DebugPprof().Stats().Worker("w2", func() error { return nil }, 1, false)
	if err := s.Run(); err != nil {
		t.Errorf("Error when run workerServer %v", err)
	}
}

func Test_server_Worker(t *testing.T) {
	nameWorker := "w1"
	handleWorker := func() error { return nil }
	concurrencyWorker := 1
	restartAlwaysWorker := false
	s := New().Worker(nameWorker, handleWorker, concurrencyWorker, restartAlwaysWorker)
	workers := s.Workers()
	if len(workers) != 1 {
		t.Error("Number of workers is different from setup")
	}
	w := workers[0]
	if w.Name != nameWorker {
		t.Errorf("Name of worker if different from setup %s != %s", w.Name, nameWorker)
	}
	if w.Handle != nil &&  fmt.Sprintf("%p", w.Handle) != fmt.Sprintf("%p", handleWorker) {
		t.Errorf("Name of worker if different from setup %p != %p", w.Handle, handleWorker)
	}
	if w.Concurrency != concurrencyWorker {
		t.Errorf("Concurrency of worker if different from setup %d != %d", w.Concurrency, concurrencyWorker)
	}
	if w.RestartAlways != restartAlwaysWorker {
		t.Errorf("RestartAlaways of worker if different from setup %t != %t", w.RestartAlways, restartAlwaysWorker)
	}
}
