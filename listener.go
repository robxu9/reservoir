package reservoir

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

var n1 byte = 10
var messageProcessors map[string]Processable = make(map[string]Processable)
var run bool

func processMessage(msg string) {
	msgobj := new(Message)
	err := json.Unmarshal(msg, msgobj)
	if err != nil {
		Errorf("Could not process message: %s\n", err)
	}
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
			Errorf("Error %s from %s.\n", err.Error(), conn.RemoteAddr().String())
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
	l, err := net.ListenTCP("tcp", net.ResolveTCPAddr("tcp", ":24096"))

	if err != nil {
		Panicf(err)
	}

	run = true

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			Errorf("Couldn't accept connection: %s\n", err.Error())
			continue
		}
		go handle(conn)
	}
}
