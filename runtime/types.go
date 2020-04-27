package runtime

import "github.com/dalmarcogd/go-worker-pool/worker"

type (
	//Server interface that define the contract to be used between monitoring.http and workerServer
	Server interface {
		Workers() []*worker.Worker
	}

	//FakeServer interface that define the contract to be used between monitoring.http and workerServer
	//for tests only
	FakeServer struct{}
)
