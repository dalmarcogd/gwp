package stats

import (
	"context"
	"encoding/json"
	"github.com/dalmarcogd/gwp/internal"
	"github.com/dalmarcogd/gwp/pkg/worker"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {

	req, err := http.NewRequest(http.MethodGet, "/stats", nil)
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

	if len(body["workers"].([]interface{})) != 0 {
		t.Errorf("Was expected any one worker but returned %d", len(body["workers"].([]interface{})))
	}

	internal.SetServerRun(STFakeServer{})

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

	if len(body["workers"].([]interface{})) != 1 {
		t.Errorf("Was expected one worker but returned %d", len(body["workers"].([]interface{})))
	}

	req, err = http.NewRequest(http.MethodPost, "/stats", nil)
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

type STFakeServer struct{}

func (s STFakeServer) Infos() map[string]interface{} {
	return internal.ParseServerInfos(s)
}

func (s STFakeServer) Healthy() bool {
	return true
}

func (STFakeServer) Workers() []*worker.Worker {
	w := worker.NewWorker("w1", func(ctx context.Context) error {
		return nil
	})
	w.FinishedAt = time.Now().UTC()
	return []*worker.Worker{
		w,
	}
}
