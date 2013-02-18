package reservoir

type Worker struct {
	location string
	port uint16
}

// Send a job to the worker.
func (w *Worker) sendJob(job *Job) {
	
}

// Return the status of a job, with 0 being unfinished
// and 255 being completed.
func (w *Worker) statusJob(job *Job) uint8 {
	
	return 0
}

func (w *Worker) ping() bool {
	return true
}
