package healthcheck

import (
	"context"
	"encoding/json"
	"github.com/dalmarcogd/gwp/internal"
	"github.com/dalmarcogd/gwp/worker"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	//We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)

	//Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	//directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body := map[string]interface{}{}
	err = json.NewDecoder(rr.Body).Decode(&body)
	if err != nil {
		t.Errorf("Error when decode body responde: %v", err)
	}

	if !body["status"].(bool) {
		t.Errorf("Was expected the status true but returned %t", body["status"].(bool))
	}

	internal.SetServerRun(HCFakeServer{})

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body = map[string]interface{}{}
	err = json.NewDecoder(rr.Body).Decode(&body)
	if err != nil {
		t.Errorf("Error when decode body responde: %v", err)
	}

	if !body["status"].(bool) {
		t.Errorf("Was expected the status true but returned %t", body["status"].(bool))
	}

	req, err = http.NewRequest(http.MethodPost, "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Error returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

type HCFakeServer struct{}

func (s HCFakeServer) Infos() map[string]interface{} {
	return internal.ParseServerInfos(s)
}

func (s HCFakeServer) Healthy() bool {
	return true
}

func (HCFakeServer) Workers() []*worker.Worker {
	w := worker.NewWorker("w1", func(ctx context.Context) error {
		return nil
	})
	w.FinishedAt = time.Now().UTC()
	return []*worker.Worker{
		w,
	}
}
