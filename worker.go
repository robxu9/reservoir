package reservoir

import (
	"net"
	"sync"
)

type WorkerID struct {
	WorkerName  string
	WorkerSubID uint64
}

type WorkerConnection struct {
	Host       string
	Connection *net.TCPConn
	ExitTask   string
	Closed     bool
}

var Workers workerMap
var workerExitChan chan bool
var workerexittask string
var workerwaitchan chan bool

type workerMap struct {
	lock     *sync.RWMutex
	workMap  map[WorkerID]*WorkerConnection
}

func (w *workerMap) Keys() []WorkerID {
	slice := make([]WorkerID, 0)
	w.lock.RLock()
	for key := range w.workMap {
		slice = append(slice, key)
	}
	w.lock.RUnlock()
	return slice
}

func (w *workerMap) Has(key WorkerID) bool {
	w.lock.RLock()
	_, ok := w.workMap[key]
	w.lock.RUnlock()
	return ok
}

func (w *workerMap) Get(key WorkerID) *WorkerConnection {
	w.lock.RLock()
	val := w.workMap[key]
	w.lock.RUnlock()
	return val
}

func (w *workerMap) Set(key WorkerID, value *WorkerConnection) *WorkerConnection {
	w.lock.Lock()
	oldVal := w.workMap[key]
	if value == nil {
		delete(w.workMap, key)
	} else {
		w.workMap[key] = value
	}
	w.lock.Unlock()
	return oldVal
}

func init() {
	Workers = workerMap{
		&sync.RWMutex{},
		make(map[WorkerID]*WorkerConnection),
	}
	
	workerwaitchan = make(chan bool)

	term := func() {
		workerExitChan <- true
		<-workerwaitchan
	}

	workerExitChan = make(chan bool)
	if AddDefinedExitTask("Worker", term) {
		workerexittask = "Worker"
	} else {
		workerexittask = AddExitTask(term)
	}
	
	go func() {
		<-workerExitChan
		// TODO
		// disconnect from all workers, pull jobs that haven't finished
		// and wait until completed.
	}()
}

type Worker struct {
	WorkerName string
	WorkerSub  uint64
	Host       string
	Connection *net.TCPConn
	ExitTask   string
	Closed     bool
}

func (w *Worker) DispatchJob(job *ReservoirJob) {

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
