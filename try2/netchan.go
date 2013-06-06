package reservoir

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	WORKER_VERSION = 104

	JOB_PING byte = iota // send PING, return HandlerProperties in a nutshell
	JOB_ACK
	JOB_SEND   // send SEND, OK means that it will accept a job
	JOB_SENT   // sent the whole job
	JOB_NOSEND // no longer sending jobs, finish remaining, send us NOSEND, and close
	JOB_SUCCESS
	JOB_FAILURE
	JOB_WARNING
)

type JobChan struct {
	Receive  chan *Job
	Send     chan *Job
	mutex    *sync.RWMutex // For handlers
	handlers []*Handler
}

type Handler struct {
	properties HandlerProperties
	connection *net.TCPConn
}

type HandlerProperties struct {
	Arch        string
	Hostname    string
	CPU         int
	MaxProcs    int
	GoVersion   string
	Environment []string
	GID         int
	UID         int

	CurrentWorkers int
	MaxWorkers     int

	PingTime int64
}

func (h *Handler) Send(code byte, details string) {
	_, err := fmt.Fprintf(h.connection, "%d:%s", code, details)
	if err != nil {
		log.Printf("Error occurred sending to %s: %s\n", h.connection.RemoteAddr().String(), err)
	}
}

// not reusable
func NewJobChan() *JobChan {
	newChannel := &JobChan{
		make(chan *Job),
		make(chan *Job),
		&sync.RWMutex{},
		make(map[HandlerProperties]*net.TCPConn),
	}
}

func (j *JobChan) AddNewConn(tcpAddr *net.TCPAddr) bool {
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("Error adding new handler %s: %s\n", tcpAddr.String(), err)
		return false
	}
	fmt.Fprintln(tcpConn, JOB_PING)
	result, err := bufio.NewReader(tcpConn).ReadString("\n")
	if err != nil {
		log.Printf("Error pinging new handler %s: %s\n", tcpAddr.String(), err)
		tcpConn.Close()
		return false
	}
	var handler HandlerProperties
	err = json.Unmarshal(result, &handler)
	if err != nil {
		log.Printf("Error parsing new handler %s: %s\n", tcpAddr.String(), err)
		tcpConn.Close()
		return false
	}

	j.mutex.Lock()
	j.handlers[handler] = tcpConn
	j.mutex.Unlock()

	go func() { // Receive monitor
		reader := bufio.NewReader(tcpConn)
		for {
			code, err := reader.ReadString("\n")
		}

	}()

	return true
}

func (j *JobChan) PingMonitor() {
	for {
		j.mutex.RLock()
		for handler, conn := range j.handlers {
			bftime := time.Now().UnixNano()
			resultChan :=
				fmt.Fprintln(conn, JOB_PING)
			result, err := bufio.NewReader(conn).ReadString("\n")
			if err != nil {
				log.Printf("Error pinging handler %s: %s\n", tcpAddr.String(), err)
				tcpConn.Close()
				return false
			}
			var handler HandlerProperties
			err = json.Unmarshal(result, &handler)
			if err != nil {
				log.Printf("Error parsing handler %s: %s\n", tcpAddr.String(), err)
				tcpConn.Close()
				return false
			}
			j.mutex.Lock()
			j.handlers[handler] = tcpConn
			j.mutex.Unlock()
		}
	}
}

func (j *JobChan) SendMonitor() {
	for {
		next, ok := <-j.Send
		if !ok {
			break
		}
	}
	j.mutex.Lock()
	// notify all writing connections of breakoff
	for handle, conn := range j.connections {
		err := conn.CloseWrite()
		log.Printf("Closed writing connection %s with error %s.\n", conn.RemoteAddr().String(), err)
		delete(j.connections[handle])
	}
	j.mutex.Unlock()
}
