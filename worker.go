package reservoir

import (
	"net"
)

var WorkerHosts map[string]map[uint64]*Worker = make(map[string]map[uint64]*Worker)

func AddWorker(worker *Worker) {
	if WorkerHosts[worker.WorkerName] == nil {
		WorkerHosts[worker.WorkerName] = make(map[uint64]*Worker)
	}
	WorkerHosts[worker.WorkerName][worker.WorkerSub] = worker
}

func RmWorker(worker *Worker) {
	delete(WorkerHosts[worker.WorkerName], worker.WorkerSub)
}

type Worker struct {
	WorkerName string
	WorkerSub  uint64
	Host       string
	Connection *net.TCPConn
	ExitTask   string
	Closed     bool
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

func (w *Worker) Dial() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", w.Host)
	if err != nil {
		return err
	}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	w.Connection = tcpConn
	w.ExitTask = AddExitTask(func() {
		w.Connection.Close()
	})
	w.Closed = false
	return nil
}

func (w *Worker) Shutdown() {
	// Finish remaining jobs and shutdown
	w.Connection.Close()
	RmExitTask(w.ExitTask)
	w.Closed = true
}

func (w *Worker) IsShutdown() bool {
	return w.Closed
}
