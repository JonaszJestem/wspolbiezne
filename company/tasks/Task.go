package tasks

type Task struct {
	Left         int
	Operation    Operation
	Right        int
	Result       int
	Successful   bool
	BreakingTask bool
}

func CreateTask(left int, right int, operation Operation) Task {
	return Task{
		Left:         left,
		Right:        right,
		Operation:    operation,
		Successful:   false,
		BreakingTask: false,
	}
}

type TaskRequest struct {
	Task chan Task
}
