package main

import (
	"serverplatform/server"
)

func main() {
	a := server.NewAcceptor("localhost:8081")
	a.Accept()
}
