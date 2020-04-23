package runtime

import "github.com/dalmarcogd/go-worker-pool/worker"

//Server
type Server interface {
	Workers() []*worker.Worker
}
