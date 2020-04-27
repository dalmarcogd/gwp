package monitoring

import (
	"context"
	"errors"
	"fmt"
	"github.com/dalmarcogd/go-worker-pool/monitoring/healthcheck"
	"github.com/dalmarcogd/go-worker-pool/monitoring/stats"
	"log"
	"net/http"
	"net/http/pprof"
	"strconv"
)

var serverHTTP *http.Server

//SetupHTTP the http server to be used for monitoring the workers
func SetupHTTP(configs map[string]interface{}) {
	port := 0
	if p, ok := configs["port"]; ok {
		port = p.(int)
	}
	host := ""
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

	var statsFunc func(http.ResponseWriter, *http.Request)
	if stf, ok := configs["statsFunc"]; ok && stf != nil {
		statsFunc = stf.(func(http.ResponseWriter, *http.Request))
	} else {
		statsFunc = stats.Handler
	}

	hc := false
	if h, ok := configs["healthCheck"]; ok {
		hc = h.(bool)
	}

	var healthCheckFunc func(http.ResponseWriter, *http.Request)
	if hcf, ok := configs["healthCheckFunc"]; ok && hcf != nil {
		healthCheckFunc = hcf.(func(http.ResponseWriter, *http.Request))
	} else {
		healthCheckFunc = healthcheck.Handler
	}

	debugPprof := false
	if dp, ok := configs["debugPprof"]; ok {
		debugPprof = dp.(bool)
	}

	needHTTP := st || hc || debugPprof
	if needHTTP {
		serveMux := http.NewServeMux()
		if hc {
			serveMux.HandleFunc(fmt.Sprintf("%s/health-check", basePath), healthCheckFunc)
		}

		if st {
			serveMux.HandleFunc(fmt.Sprintf("%s/stats", basePath), statsFunc)
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
			serverHTTP = &http.Server{Addr: address, Handler: serveMux}
			if err := serverHTTP.ListenAndServe(); err != nil {
				log.Print(err)
			}
		}(serveMux)
	}
}



//CloseHTTP the http server to be used by monitoring
func CloseHTTP() error {
	if serverHTTP != nil {
		defer log.Printf("Shutdown monitoring server at %s", serverHTTP.Addr)
		defer func() {serverHTTP = nil}()
		return serverHTTP.Shutdown(context.Background())
	}

	return errors.New("the serverHTTP is not configured")
}
