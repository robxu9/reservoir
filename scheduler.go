package reservoir

/*
	SchedulerStatus shows the current status of Reservoir's scheduler:

	0 : Stopped
	1 : Starting
	2 : Waiting for available workers
	3 : Waiting for jobs
	4 : Dispatching job
*/
const (
	SCHEDULER_STOPPED uint8 = iota
	SCHEDULER_STARTING
	SCHEDULER_WAITWORKER
	SCHEDULER_WAITJOB
	SCHEDULER_DISPATCHJOB
)

var SchedulerStatus uint8

var runScheduler bool
var exittask string

var JobsChannel chan *Job      // Pending Jobs
var WorkerChannel chan *Worker // Idle Workers (minus 1, as it's pulled by job)

type ready struct {
	job    *Job
	worker *Worker
}

var readyChannel chan *ready

func init() {
	JobsChannel = make(chan *Job, 10000)
	WorkerChannel = make(chan *Worker, 1000)
	readyChannel = make(chan *ready)
	SchedulerStatus = SCHEDULER_STOPPED
	runScheduler = false
}

func Scheduler_QueueJob(j *Job) {
	JobsChannel <- j
}

func Scheduler_QueueWorker(w *Worker) {
	WorkerChannel <- w
}

func Scheduler_Stop() {
	runScheduler = false
	RmExitTask(exittask)
}

func Scheduler_Run() {
	runScheduler = true
	exittask = AddExitTask(Scheduler_Stop)
	SchedulerStatus = SCHEDULER_STARTING
	t := make(chan bool)
	go scheduler_terminate(t)
	go scheduler_dispatcher()

	go func() {
		for {
			select {
			case rjob := <-readyChannel:
				go rjob.worker.DispatchJob(rjob.job)
			case <-t:
				return
			}
		}
	}()
}

func scheduler_dispatcher() {
	SchedulerStatus = SCHEDULER_WAITWORKER
	worker := <-WorkerChannel
	SchedulerStatus = SCHEDULER_WAITJOB
	job := <-JobsChannel
	SchedulerStatus = SCHEDULER_DISPATCHJOB
	readyChannel <- &ready{job, worker}
}

func scheduler_terminate(t chan bool) {
	for {
		if !runScheduler {
			t <- true
			return
		}
	}
}
