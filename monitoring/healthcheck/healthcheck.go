package healthcheck

import (
	"encoding/json"
	"net/http"
)

//Handler
func Handler(writer http.ResponseWriter, request *http.Request) {
	if http.MethodGet != request.Method {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(map[string]interface{}{
		"status": true,
	})
}
