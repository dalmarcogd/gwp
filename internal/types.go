package internal

import "github.com/dalmarcogd/gwp/worker"

type (
	//Server interface that define the contract to be used between monitor.http and workerServer
	Server interface {
		Workers() []*worker.Worker
		Healthy() bool
		Infos() map[string]interface{}
	}

	//FakeServer interface that define the contract to be used between monitor.http and workerServer
	//for tests only
	FakeServer struct{}
)
