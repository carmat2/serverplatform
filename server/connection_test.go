package server

import(
	"fmt"
	"testing"
	"net"
	"time"
)

func TestMain(m *testing.M) {
	CreateLoggerTargetTesting()

	a := NewAcceptor("localhost:8081")
	go a.Accept()

	logger.Trace("testing connection - acceptor started in listen mode")
	time.Sleep(5*time.Second)

	m.Run()

	ShutdownLogger()
}

func Test_Connection_Read(t *testing.T) {   
	logger.Trace("Testing connection read")

	conn := connect(t)

	write("$SIGNA", conn, t)

	time.Sleep(5*time.Second)
 
	jsonMsgOpenSession := `{"name":"opensession","dest":"sp","payload":[{"plugin":"helloworld","token":"123"}]}`
	msgSize := fmt.Sprintf("[%04d]",len(jsonMsgOpenSession))
	
	write("TURE"+msgSize+jsonMsgOpenSession, conn, t)

	jsonMsgOpenSession = `{"name":"opensession","dest":"sp","payload":[{"plugin":"helloworld_hellogo","token":"123"}]}`
	msgSize = fmt.Sprintf("[%04d]",len(jsonMsgOpenSession))

	write(msgSize+jsonMsgOpenSession, conn, t)

	time.Sleep(5*time.Second)
}

func Test_Connection_Read_invalidData(t *testing.T) {   
	logger.Trace("Testing connection read with invalid signature")

	conn := connect(t)

	write("$SIGNATUR#", conn, t)

	time.Sleep(5*time.Second)
}

func Test_Connection_Close(t *testing.T) {
	logger.Trace("Testing connection close")

	conn := connect(t)

	write("$SIGNA", conn, t)

	time.Sleep(5*time.Second)
	conn.Close()
}

func connect(t *testing.T) (*net.TCPConn) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:8081")
    if err != nil {
        t.Error("Failed to resolve TCP address localhost:8081")
    }

    conn, err := net.DialTCP("tcp", nil, addr)
    if err != nil {
        t.Error("TCP connection to address localhost:8081 failed")
    }

	return conn
}

func write(data string, conn *net.TCPConn, t *testing.T){
	count, err := conn.Write([]byte(data))
    if err != nil {
		t.Error("TCP write to server failed", err)
    }
	logger.Tracef("testing connection - write data completed %s, byte count %d", data, count)
}