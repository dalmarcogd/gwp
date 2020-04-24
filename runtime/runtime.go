package runtime

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
