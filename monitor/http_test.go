package monitor

import (
	"github.com/dalmarcogd/gwp/monitor/healthcheck"
	"github.com/dalmarcogd/gwp/monitor/stats"
	"net/http"
	"testing"
	"time"
)

func TestSetupHTTP(t *testing.T) {
	SetupHTTP(map[string]interface{}{
		"port":            8002,
		"host":            "localhost",
		"stats":           true,
		"statsFunc":       stats.Handler,
		"healthCheck":     true,
		"healthCheckFunc": healthcheck.Handler,
		"debugPprof":      true,
		"basePath":        "",
	})

	<-time.After(1 * time.Second)

	response := make(chan *http.Response, 3)

	go func() {
		r, err := http.Get("http://localhost:8002/stats")
		if err != nil || r == nil || r.StatusCode != http.StatusOK {
			statusCode := 0
			if r != nil {
				statusCode = r.StatusCode
			}
			t.Errorf("Request for /stats returned error: %s, %d", err, statusCode)
		}
		response <- r
	}()

	go func() {
		r, err := http.Get("http://localhost:8002/health-check")
		if err != nil || r == nil || r.StatusCode != http.StatusOK {
			statusCode := 0
			if r != nil {
				statusCode = r.StatusCode
			}
			t.Errorf("Request for /health-check returned error: %s, %d", err, statusCode)
		}
		response <- r
	}()

	go func() {
		r, err := http.Get("http://localhost:8002/debug/pprof")
		if err != nil || r == nil || r.StatusCode != http.StatusOK {
			statusCode := 0
			if r != nil {
				statusCode = r.StatusCode
			}
			t.Errorf("Request for /debug/pprof returned error: %s, %d", err, statusCode)
		}
		response <- r
	}()

	<-response

	<-time.After(1 * time.Second)

	err := CloseHTTP()
	if err != nil {
		t.Errorf("Error when close http: %v", err)
	}
}

func TestCloseHTTP(t *testing.T) {
	SetupHTTP(map[string]interface{}{
		"port":        8002,
		"host":        "localhost",
		"stats":       true,
		"healthCheck": true,
		"debugPprof":  true,
		"basePath":    "",
	})

	<-time.After(1 * time.Second)
	err := CloseHTTP()
	if err != nil {
		t.Errorf("Error when close http: %v", err)
	}

	serverHTTP = nil
	err = CloseHTTP()
	if err == nil {
		t.Errorf("Error when close http: %v", err)
	}
}
