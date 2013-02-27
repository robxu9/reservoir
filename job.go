package reservoir

import (
	"sync"
)

const (
	JOBID_LENGTH = 15
)

var random *RS = NewAlphaNumericRS()
var randCall chan bool = make(chan bool, 10)
var randGet chan string = make(chan string, 10)

func init() {
	go func() {
		for {
			<-randCall
			randGet <- random.NewRandomString(15)
		}
	}()
}

type Job interface {
	GetId() string
	GetScript() string
	GetWorker() *Worker
	SetWorker(w *Worker)
}

type reservoirJob struct {
	id     string
	script string
	worker *Worker
	once   sync.Once
}

func NewReservoirJob(script string) Job {
	return &reservoirJob{
		"",
		script,
		nil,
		sync.Once{},
	}
}

func (r *reservoirJob) init() {
	r.once.Do(func() {
		randCall <- true
		r.id = <-randGet
	})
}

func (r *reservoirJob) GetId() string {
	r.init()
	return r.id
}

func (r *reservoirJob) GetScript() string {
	return r.script
}

func (r *reservoirJob) GetWorker() *Worker {
	return r.worker
}

func (r *reservoirJob) SetWorker(w *Worker) {
	r.worker = w
}
