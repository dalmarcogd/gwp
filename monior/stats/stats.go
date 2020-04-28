package stats

import (
	"encoding/json"
	gwpr "github.com/dalmarcogd/go-worker-pool/runtime"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

//Handler
func Handler(writer http.ResponseWriter, request *http.Request) {
	if http.MethodGet != request.Method {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	response := map[string]interface{}{
		"cpus":       strconv.Itoa(runtime.NumCPU()),
		"goroutines": strconv.Itoa(runtime.NumGoroutine()),
		"workers":    []map[string]interface{}{},
	}

	for _, worker := range gwpr.GetServerRun().Workers() {
		finishedAt := ""
		if !worker.FinishedAt.IsZero() {
			finishedAt = worker.FinishedAt.Format(time.RFC3339)
		}
		response["workers"] = append(response["workers"].([]map[string]interface{}), map[string]interface{}{
			"id":             worker.ID,
			"name":           worker.Name,
			"concurrency":    worker.Concurrency,
			"restart_always": worker.RestartAlways,
			"restarts":       worker.Restarts,
			"started_at":     worker.StartAt.Format(time.RFC3339),
			"finished_at":    finishedAt,
			"status":         worker.Status(),
		})
	}

	_ = json.NewEncoder(writer).Encode(response)
	writer.Header().Set("Content-Type", "application/json")
}
