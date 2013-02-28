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

type workerMap struct {
	lock     *sync.Mutex
	requests chan workerMapReq
	workMap  map[WorkerID]*WorkerConnection
}

func (w *workerMap) Keys() []WorkerID {
	slice := make([]WorkerID, 0)
	w.lock.Lock()
	for key := range w.workMap {
		slice = append(slice, key)
	}
	w.lock.Unlock()
	return slice
}

func (w *workerMap) Get(key WorkerID) *WorkerConnection {
	retchan := make(chan *WorkerConnection)
	Workers.requests <- workerMapReq{true, key, nil, retchan}
	return <-retchan
}

func (w *workerMap) Set(key WorkerID, value *WorkerConnection) *WorkerConnection {
	retchan := make(chan *WorkerConnection)
	Workers.requests <- workerMapReq{false, key, value, retchan}
	return <-retchan
}

type workerMapReq struct {
	get     bool
	key     WorkerID
	value   *WorkerConnection
	retchan chan *WorkerConnection
}

func init() {
	Workers = workerMap{
		&sync.Mutex{},
		make(chan workerMapReq),
		make(map[WorkerID]*WorkerConnection),
	}

	term := func() {
		workerExitChan <- true
	}

	workerExitChan = make(chan bool)
	if AddDefinedExitTask("Worker", term) {
		workerexittask = "Worker"
	} else {
		workerexittask = AddExitTask(term)
	}

	go func() {
		for {
			select {
			case req := <-Workers.requests:
				go workers_modify(req)
			case <-workerExitChan:
				return
			}
		}
	}()
}

func workers_modify(req workerMapReq) {
	Workers.lock.Lock()
	oldval := Workers.workMap[req.key]
	if req.value == nil && !req.get {
		delete(Workers.workMap, req.key)
	} else if req.value != nil {
		Workers.workMap[req.key] = req.value
	}
	Workers.lock.Unlock()
	req.retchan <- oldval
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
