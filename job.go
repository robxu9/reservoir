package reservoir

import (
	"sync"
)

const (
	JOBID_LENGTH = 15
)

var random *RS = NewAlphaNumericRS()

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

func NewReservoirJob(script string) *Job {
	return &reservoirJob{
		"",
		script,
		nil,
		sync.Once{},
	}
}

func (r *reservoirJob) init() {
	r.once.Do(func() {
		r.Id = random.NewRandomString(15)
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
