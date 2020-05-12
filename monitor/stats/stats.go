package stats

import (
	"encoding/json"
	"github.com/dalmarcogd/gwp/internal"
	"net/http"
)

//Handler that return the stats from workerServer
func Handler(writer http.ResponseWriter, request *http.Request) {
	if http.MethodGet != request.Method {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_ = json.NewEncoder(writer).Encode(internal.GetServerRun().Infos())
	writer.Header().Set("Content-Type", "application/json")
}
