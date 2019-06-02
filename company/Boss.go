package company

import (
	. "wspolbiezne/company/tasks"
	"fmt"
	"math/rand"
	"time"
)

type Boss struct {
	name    string
	addTask chan<- Task
	loud    bool
}

func CreateBoss(name string, addTask chan<- Task, loud bool) Boss {
	return Boss{name, addTask, loud}
}

func StartBossing(boss Boss) {
	for {
		task := createRandomTask()
		boss.addTask <- task

		if boss.loud {
			fmt.Printf("%s produces new task:\t\t %v\n", boss.name, task)
		}
		time.Sleep(NewTaskInterval)
	}
}

func createRandomTask() Task {
	leftOperand := rand.Int()%10 + 1
	rightOperand := rand.Int()%10 + 1
	randomOperation := Operations[rand.Intn(len(Operations))]

	return CreateTask(leftOperand, rightOperand, randomOperation)
}
