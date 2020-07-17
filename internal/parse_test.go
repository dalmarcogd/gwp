package internal

import (
	"context"
	"github.com/dalmarcogd/gwp/pkg/worker"
	"testing"
	"time"
)

func TestParseServerInfos(t *testing.T) {
	infos := ParseServerInfos(FakeServer{})
	if _, ok := infos["cpus"]; !ok {
		t.Error("Expected key cpus on result from ParseServerInfos")
	}
	if _, ok := infos["goroutines"]; !ok {
		t.Error("Expected key goroutines on result from ParseServerInfos")
	}
	if _, ok := infos["workers"]; !ok {
		t.Error("Expected key workers on result from ParseServerInfos")
	}

	infos = ParseServerInfos(ParseFakeServer{})

	if _, ok := infos["cpus"]; !ok {
		t.Error("Expected key cpus on result from ParseServerInfos")
	}
	if _, ok := infos["goroutines"]; !ok {
		t.Error("Expected key goroutines on result from ParseServerInfos")
	}
	if ws, ok := infos["workers"]; !ok || len(ws.([]map[string]interface{})) != 1 {
		t.Error("Expected key workers on result from ParseServerInfos")
	}
}

type ParseFakeServer struct{}

func (s ParseFakeServer) Infos() map[string]interface{} {
	return ParseServerInfos(s)
}

func (ParseFakeServer) Healthy() bool {
	return true
}

func (ParseFakeServer) Workers() []*worker.Worker {
	w := worker.NewWorker("w1", func(ctx context.Context) error {
		return nil
	})
	deadline := worker.WithDeadline(time.Now())
	deadline.Apply(w)
	w.FinishedAt = time.Now().UTC()
	return []*worker.Worker{
		w,
	}
}
