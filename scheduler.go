package reservoir

/*
	SchedulerStatus shows the current status of Reservoir's scheduler:

	0 : Stopped
	1 : Waiting for available workers
	2 : Waiting for jobs
	3 : Dispatching job
*/
var SchedulerStatus uint8

var runScheduler bool

var jobChannel chan *Job
var workerChannel chan *Worker

func init() {
	jobChannel = make(chan *Job, 100)
	workerChannel = make(chan *Worker, 100)
	SchedulerStatus = 0
	runScheduler = false
}

func Scheduler_QueueJob(j *Job) {
	jobChannel <- j
}

func Scheduler_QueueWorker(w *Worker) {
	workerChannel <- w
}

func Scheduler_Stop() {
	runScheduler = false
}

func Scheduler_Run() {
	runScheduler = true
	for runScheduler {
		SchedulerStatus = 1
		worker := <-workerChannel
		SchedulerStatus = 2
		job := <-jobChannel
		SchedulerStatus = 3
		go worker.DispatchJob(job)
	}
}
