package server

import(
	"fmt"
	"testing"
	"net"
	"time"
)

func TestMain(m *testing.M) {
	a := NewAcceptor("localhost:8081")
	go a.Accept()

	fmt.Println("testing connection - acceptor started in listen mode")
	time.Sleep(5*time.Second)

	m.Run()
}

func Test_Connection_Read(t *testing.T) {   
	fmt.Println("\r\n ****** testing connection read")

	addr, err := net.ResolveTCPAddr("tcp", "localhost:8081")
    if err != nil {
        t.Error("Failed to resolve TCP address localhost:8081")
    }

    conn, err := net.DialTCP("tcp", nil, addr)
    if err != nil {
        t.Error("TCP connection to address localhost:8081 failed")
    }

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
	fmt.Println("\r\n ****** testing connection read with invalid signature")

	addr, err := net.ResolveTCPAddr("tcp", "localhost:8081")
    if err != nil {
        t.Error("Failed to resolve TCP address localhost:8081")
    }

    conn, err := net.DialTCP("tcp", nil, addr)
    if err != nil {
        t.Error("TCP connection to address localhost:8081 failed")
    }

	write("$SIGNATUR#", conn, t)

	time.Sleep(5*time.Second)
}

func Test_Connection_Close(t *testing.T) {
	fmt.Println("\r\n ****** testing connection close")

	addr, err := net.ResolveTCPAddr("tcp", "localhost:8081")
    if err != nil {
        t.Error("Failed to resolve TCP address localhost:8081")
    }

    conn, err := net.DialTCP("tcp", nil, addr)
    if err != nil {
        t.Error("TCP connection to address localhost:8081 failed")
    }

	write("$SIGNA", conn, t)

	time.Sleep(5*time.Second)
	conn.Close()
}

func write(data string, conn *net.TCPConn, t *testing.T){
	count, err := conn.Write([]byte(data))
    if err != nil {
		t.Error("TCP write to server failed", err)
    }
	fmt.Printf("testing connection - write data completed %s, byte count %d\r\n", data, count)
}