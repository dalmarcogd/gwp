package internal

import (
	"github.com/dalmarcogd/gwp/pkg/worker"
)

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

//Healthy return the health of server
func (f FakeServer) Healthy() bool {
	return true
}

//Infos return the empty infos
func (f FakeServer) Infos() map[string]interface{} {
	return ParseServerInfos(f)
}
