package reservoir

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
)

var n1 byte = 10
var run bool

func processMessage(msg string) {

}

func handle(conn *net.TCPConn) {
	addr := conn.RemoteAddr()
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	// MESSAGE
	var jsonMSG bytes.Buffer

	for {
		s, err := rw.ReadString(n1)
		if len(s) > 0 {
			jsonMSG.WriteString(s)
		} else if err == io.EOF {
			// End of file
			conn.Close()
			go processMessage(jsonMSG.String())
			return
		} else {
			fmt.Printf("Error %s from %s.\n", err.Error(), conn.RemoteAddr().String())
			conn.Close()
			return
		}
	}
}

func Listener_Status() bool {
	return run
}

func Listener_Stop() {
	run = false
}

func Listener_Run() {
	l, err := net.ListenTCP("tcp4", &net.TCPAddr{net.IPv4zero, 24098})

	if err != nil {
		panic(err)
	}

	run = true

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			fmt.Printf("Couldn't accept connection: %s\n", err.Error())
			continue
		}
		go handle(conn)
	}
}