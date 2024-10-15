package server

import (
	"fmt"
	"net"
	"io"
)

// connection wraps a client TCP connection
// 
// chRead is the read request channel to ensure only one processor has a read in progress
// readBuff is the buffer used to read data from the TCP connection, created with len/cap set to 4096 bytes
// availableRead is 0 if the readBuff is empty, and set to the number of available bytes in the read buffer still to be processed
type connection struct {
	conn     net.Conn
	readBuff []byte
	availableRead	int
	chRead   chan (processor)
}

// Creates a new connection instance and starts the read loop with the signature validator process
func newConn(conn net.Conn) *connection {
	readBuff := make([]byte, 4096)
	chRead := make(chan processor, 1)

	c := &connection{
		conn:     conn,
		readBuff: readBuff,
		availableRead: 0,
		chRead:   chRead,
	}

	v := NewValidator()
	v.setDataLen(len(signature))
	c.chRead <- v

	return c
}

// read is the connection read loop
// on each iteration the current connection processor is passed available bytes for processing
func (c *connection) read() {
	var p processor
	for {
		// one processor can issue a read request to this connection
		p = <-c.chRead

		var count int
		if c.availableRead > 0 {
			// if bytes are already read but not processed, call the processor to process the pending data bytes 
			count = c.availableRead
		} else {
			var err error
			// blocking connection read
			count, err = c.conn.Read(c.readBuff)
			if err == io.EOF {
				fmt.Println("connection read() - err io.EOF, closing the connection")
				defer c.conn.Close()
				defer close(c.chRead)
				return;
			}
			if err != nil {
				fmt.Printf("connection read() - err %s, closing the connection", err)
				defer c.conn.Close()
				defer close(c.chRead)
				return;
			}
		}

		// call the current processor with the new data slice
		processed, errProtocol := p.processReadData(c.readBuff[:count])
		
		if errProtocol == nil || errProtocol.Status() == needMoreConnRead {
			if processed < count {
				// the current processor did not process all available bytes, so the extra data is moved to the beginning of the read buffer
				// and the availableRead var is set to the byte count
				fmt.Printf("connection read() - consolidating read buffer after processor handling - processed bytes %d, available bytes %d\r\n", processed, count)
				c.availableRead = count-processed
				for i := 0; i < c.availableRead; i++ {
					c.readBuff[i] = c.readBuff[processed+i]
				}
			} else if processed == count {
				// the current processor processed all available bytes from the read buffer
				// there will be a next read request issued with either the currect processor, or the next processor in the logic
				c.availableRead = 0
				fmt.Printf("connection read() - processor processed all available bytes %d\r\n", processed)
			} else {
				fmt.Printf("connection read() - processor return invalid processed value %d, count %d, closing connection\r\n", processed, count)
				defer c.conn.Close()
				defer close(c.chRead)
				return;
			}
		}

		// errProtocol nil means the current processor completed with success
		if errProtocol != nil {
			if errProtocol.Status() == needMoreConnRead {
				// this is not an actual error case, the current processor needs more data to complete its processing
				// so issue a new read request for the currrent processor	
				c.chRead <- p
			} else {
				fmt.Printf("connection read() - errProtocol %s, closing the connection\r\n", errProtocol)
				defer c.conn.Close()
				defer close(c.chRead)
				return;
			}
		} else {
			// the current processor completed successfully, start a goroutine to issue a read request on the next processor
			p = p.getNextProcessor()
			if p != nil {
				go func() {
					c.chRead <- p
				}()
			} else {
				// there is no next processor defined, force a connection close, and close the read/write channels
				// TODO handle
			}
		}
	}
}

// write is the connection write method
// it blocks until all bytes are written in the connection
func (c *connection) write() {
	//TODO
}
