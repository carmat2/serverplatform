package main

import (
	"serverplatform/server"
)

func main() {
	server.CreateLoggerTarget()
	defer server.ShutdownLogger()

	a := server.NewAcceptor("localhost:8081")
	a.Accept()
}
