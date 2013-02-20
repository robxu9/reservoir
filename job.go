package reservoir

type Job struct {
	Id     string
	Script string
	Worker *Worker
}
