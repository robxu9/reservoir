package reservoir

import (
	"net"
)

type Worker struct {
	workername string
	workersub  uint64
	host       string
	connection *net.TCPConn
}

func (w *Worker) DispatchJob(job *Job) {

}

// Send a message to the worker
func (w *Worker) SendMessage(msg *Message) bool {

	return false
}

func (w *Worker) Ping() bool {
	return true
}

func (w *Worker) Shutdown() {
	// Finish remaining jobs and shutdown
	w.connection.Close()
}

func (w *Worker) IsShutdown() bool {
	return false
}
