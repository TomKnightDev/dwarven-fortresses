package components

import "github.com/tomknightdev/dwarven-fortresses/enums"

type Task struct {
	enums.TaskTypeEnum
	RequiredActions int
	ActionsComplete int
	Position        Position
	Completed       bool
	Blocked         bool
}

func NewTask(pos Position, taskType enums.TaskTypeEnum, requiredActions int) *Task {
	return &Task{
		TaskTypeEnum:    taskType,
		RequiredActions: requiredActions,
		Position:        pos,
	}
}

func (t *Task) CompleteTask() {
	t.Completed = true
}
