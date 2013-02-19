package reservoir

type Worker struct {
	location string
	port     uint16
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
