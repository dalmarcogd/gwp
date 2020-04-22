package monitoring

import (
	"context"
	"fmt"
	"github.com/dalmarcogd/go-worker-pool/monitoring/healthcheck"
	"github.com/dalmarcogd/go-worker-pool/monitoring/stats"
	"log"
	"net/http"
	"strconv"
)

var serverHttp *http.Server

func SetupHttp(configs map[string]interface{}) {
	port := 80
	if p, ok := configs["port"]; ok {
		port = p.(int)
	}
	host := "localhost"
	if h, ok := configs["host"]; ok {
		host = h.(string)
	}

	basePath := ""
	if bp, ok := configs["basePath"]; ok {
		basePath = bp.(string)
	}

	st := false
	if s, ok := configs["stats"]; ok {
		st = s.(bool)
	}

	hc := false
	if h, ok := configs["healthCheck"]; ok {
		hc = h.(bool)
	}

	needHttp := st || hc
	if needHttp {
		serveMux := http.NewServeMux()
		if hc {
			serveMux.HandleFunc(fmt.Sprintf("%s/health-check", basePath), healthcheck.Handler)
		}

		if st {
			serveMux.HandleFunc(fmt.Sprintf("%s/stats", basePath), stats.Handler)
		}

		go func(serveMux *http.ServeMux) {
			address := fmt.Sprintf("%s:%s", host, strconv.Itoa(port))
			log.Printf("Started monitoring server at %s", address)
			serverHttp = &http.Server{Addr: address, Handler: serveMux}
			if err := serverHttp.ListenAndServe(); err != nil {
				log.Print(err)
			}
		}(serveMux)
	}
}

func CloseHttp() error {
	defer log.Printf("Shutdown monitoring server at %s", serverHttp.Addr)
	return serverHttp.Shutdown(context.Background())
}
