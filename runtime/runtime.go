package runtime

import "github.com/dalmarcogd/go-worker-pool/worker"

var currentServer Server

func init() {
	SetServerRun(FakeServer{})
}

//SetServerRun
func SetServerRun(s Server) {
	currentServer = s
}

//GetServerRun
func GetServerRun() Server {
	return currentServer
}

func (f FakeServer) Workers() []*worker.Worker {
	return []*worker.Worker{}
}