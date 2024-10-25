package server

import (
	"net"
)

// acceptor is the server entry point.
type acceptor struct {
	address string
}

// NewAcceptor creates a new acceptor instance.
func NewAcceptor(address string) *acceptor {
	logger.Infof("Creating acceptor with address %s", address)
	acceptor := &acceptor{
		address: address,
	}

	return acceptor
}

// Accept starts the acceptor, accepts client TCP connections, and starts the new connection read and write goroutines.
func (a *acceptor) Accept() {
	var l net.Listener
	var err error

	// TODO - use chan, WaitGroup for shutdown
	l, err = net.Listen("tcp", a.address)
	if err != nil {
		logger.Errorf("Acceptor listen error", err)
		defer l.Close()
	} else {
		for {
			conn, _ := l.Accept()
			c := newConn(conn)

			go c.read()
			go c.write()
		}
	}
}
