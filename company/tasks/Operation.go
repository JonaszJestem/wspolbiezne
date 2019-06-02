package tasks

type Operation string

const (
	PLUS  Operation = "+"
	TIMES Operation = "*"
)

var Operations = [2]Operation{PLUS, TIMES}
