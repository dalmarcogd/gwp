package runtime

var currentServer Server

func SetServerRun(s Server) {
	currentServer = s
}

func GetServerRun() Server {
	return currentServer
}
