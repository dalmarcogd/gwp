package monitoring

import (
	"context"
	"fmt"
	"github.com/dalmarcogd/go-worker-pool/monitoring/healthcheck"
	"github.com/dalmarcogd/go-worker-pool/monitoring/stats"
	"log"
	"net/http"
	"net/http/pprof"
	"strconv"
)

var serverHttp *http.Server

//SetupHttp
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

	debugPprof := false
	if dp, ok := configs["debugPprof"]; ok {
		debugPprof = dp.(bool)
	}

	needHttp := st || hc || debugPprof
	if needHttp {
		serveMux := http.NewServeMux()
		if hc {
			serveMux.HandleFunc(fmt.Sprintf("%s/health-check", basePath), healthcheck.Handler)
		}

		if st {
			serveMux.HandleFunc(fmt.Sprintf("%s/stats", basePath), stats.Handler)
		}

		if debugPprof {
			serveMux.HandleFunc(fmt.Sprintf("%s/debug/pprof/", basePath), pprof.Index)
			serveMux.HandleFunc(fmt.Sprintf("%s/debug/pprof/cmdline", basePath), pprof.Cmdline)
			serveMux.HandleFunc(fmt.Sprintf("%s/debug/pprof/trace", basePath), pprof.Trace)
			serveMux.HandleFunc(fmt.Sprintf("%s/debug/pprof/profile", basePath), pprof.Profile)
			serveMux.HandleFunc(fmt.Sprintf("%s/debug/pprof/symbol", basePath), pprof.Symbol)

			serveMux.Handle(fmt.Sprintf("%s/debug/pprof/goroutine", basePath), pprof.Handler("goroutine"))
			serveMux.Handle(fmt.Sprintf("%s/debug/pprof/heap", basePath), pprof.Handler("heap"))
			serveMux.Handle(fmt.Sprintf("%s/debug/pprof/threadcreate", basePath), pprof.Handler("threadcreate"))
			serveMux.Handle(fmt.Sprintf("%s/debug/pprof/block", basePath), pprof.Handler("block"))
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

//CloseHttp
func CloseHttp() error {
	if serverHttp != nil {
		defer log.Printf("Shutdown monitoring server at %s", serverHttp.Addr)
		return serverHttp.Shutdown(context.Background())
	}
	return nil
}
