package runtime

import "github.com/dalmarcogd/gwp/worker"

type (
	//Server interface that define the contract to be used between monior.http and workerServer
	Server interface {
		Workers() []*worker.Worker
	}

	//FakeServer interface that define the contract to be used between monior.http and workerServer
	//for tests only
	FakeServer struct{}
)
