package healthcheck

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

	status := true
	for _, worker := range runtime.GetServerRun().Workers() {
		if !worker.IsUp() {
			status = false
			break
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(map[string]interface{}{
		"status": status,
	})
}
