package runtime

import "github.com/dalmarcogd/go-worker-pool/worker"

type (
	//Server
	Server interface {
		Workers() []*worker.Worker
	}

	//FakeServer
	FakeServer struct{}
)
