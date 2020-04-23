package runtime

var currentServer Server

//SetServerRun
func SetServerRun(s Server) {
	currentServer = s
}

//GetServerRun
func GetServerRun() Server {
	return currentServer
}
