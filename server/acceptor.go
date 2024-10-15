package server

import (
	"fmt"
	"net"
	"os"
)

// acceptor is the server entry point
type acceptor struct {
	address string
}

// NewAcceptor creates a new acceptor instance
func NewAcceptor(address string) *acceptor {
	fmt.Printf("creating acceptor with address %s\r\n", address)
	acceptor := &acceptor{
		address: address,
	}

	return acceptor
}

// Accept starts the acceptor, accepts client TCP connections, and starts the new connection read and write goroutines
func (a *acceptor) Accept() {
	var l net.Listener
	var err error

	// TODO - use chan, WaitGroup for shutdown
	l, err = net.Listen("tcp", a.address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, _ := l.Accept()
		c := newConn(conn)

		go c.read()
		go c.write()
	}
}
