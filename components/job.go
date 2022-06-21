package components

type Job struct {
	Tasks       []*Task
	ClaimedByID int
	EntityId    int
}

func NewJob(tasks ...*Task) Job {
	return Job{
		Tasks: tasks,
	}
}
