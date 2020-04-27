package runtime

import "github.com/dalmarcogd/go-worker-pool/worker"

var currentServer Server

func init() {
	SetServerRun(FakeServer{})
}

//SetServerRun set the instance of server that will running
func SetServerRun(s Server) {
	currentServer = s
}

//GetServerRun return the instance of server that still running
func GetServerRun() Server {
	return currentServer
}

//Workers return the worker from FakeServer
func (f FakeServer) Workers() []*worker.Worker {
	return []*worker.Worker{}
}