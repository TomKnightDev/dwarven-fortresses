package components

type Worker struct {
	HasJob     bool
	JobID      int
	GettingJob bool
}

func NewWorker() Worker {
	return Worker{}
}
