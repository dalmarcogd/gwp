package stats

import (
	"encoding/json"
	"github.com/dalmarcogd/go-worker-pool/runtime"
	"net/http"
)

//Handler
func Handler(writer http.ResponseWriter, request *http.Request) {
	if http.MethodGet != request.Method {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	response := map[string][]map[string]interface{}{
		"workers": {},
	}

	workers := runtime.GetServerRun().Workers()
	for _, worker := range workers {
		response["workers"] = append(response["workers"],  map[string]interface{}{
			"id": worker.Id,
			"name": worker.Name,
			"replicas": worker.Replicas,
			"status": worker.Status(),
		})
	}

	_ = json.NewEncoder(writer).Encode(response)
	writer.Header().Set("Content-Type", "application/json")
}