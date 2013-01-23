package job

type Job interface {
	GetId() int64
	GetType() string
	SetId(id int64)
}

type JobHandle interface {
	JobEnter(job *Job) error
	JobStatus(id int) (bool, error)
	JobFinish(id int) (*Job, error)
}

// Types and Queuer
var JobTypes map[string]JobHandle
var JobQueuer chan *Job

// ID assigner
var idcounter int64 = 0

func RegisterJobHandle(name string, handle JobHandle) bool {
	k, ok := JobTypes[name]
	if !ok {
		JobTypes[name] = handle
	}

	return !ok
}

func DeRegisterJobHandle(name string) {
	delete(JobTypes, name)
}

func init() {
	JobQueuer = make(chan *Job)
	go queueJob()
}

func queueJob() {
	for {
		job := <-JobQueuer
		if idcounter == 9223372036854775807 {
			idcounter = 0
		}
		job.SetId(idcounter)
		idcounter++

		handle, ok := JobTypes[job.GetType()]
		if !ok {
			job.SetId(-1) // signify that something went wrong
		} else {
			handle.JobEnter(job)
		}
	}
}
