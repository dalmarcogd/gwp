package runtime

import "github.com/dalmarcogd/go-worker-pool/worker"

type Server interface {
	Workers() []*worker.Worker
}
