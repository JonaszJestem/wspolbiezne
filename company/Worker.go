package company

import (
	"fmt"
	"math/rand"
	"time"
	. "wspolbiezne/company/shop"
	. "wspolbiezne/company/tasks"
)

type Worker struct {
	Name          string
	repairRequest chan *Machine
	getTask       chan TaskRequest
	addProduct    chan Product
	multiplier    []*Machine
	adders        []*Machine
	Patient       bool
	loud          bool
	timeout       time.Duration
	TasksDone     int
}

func CreateWorker(name string, repairRequest chan *Machine, getTask chan TaskRequest, addProduct chan Product, multipliers []*Machine, adders []*Machine, loud bool, timeout time.Duration, patient bool) Worker {
	return Worker{name, repairRequest, getTask, addProduct, multipliers, adders, loud, patient, timeout, 0}
}

func StartWorker(worker *Worker) {
	for {
		var requestedTask = make(chan Task)
		taskRequest := TaskRequest{Task: requestedTask}
		worker.getTask <- taskRequest

		task := <-taskRequest.Task
		for {
			machines := SelectMachine(task, *worker)
			resolvedTask := UseMachines(task, machines, worker)

			if resolvedTask.Successful {
				var product = CreateProduct(resolvedTask.Result)
				worker.TasksDone++

				if worker.loud {
					fmt.Printf("%s consumes new task \t%v and produces %v\n", worker.Name, task, product)
				}

				worker.addProduct <- product
				break
			}

			time.Sleep(WorkerSleep)
		}

	}
}

func SelectMachine(task Task, worker Worker) []*Machine {
	if task.Operation == PLUS {
		return worker.adders
	} else {
		return worker.multiplier
	}
}

func UseMachines(task Task, machines []*Machine, worker *Worker) Task {
	var resolvedTask Task
	machine := machines[rand.Intn(len(machines))]
	if !worker.Patient {
		select {
		case machine.Tasks <- task:
			resolvedTask = <-machine.Tasks
		case <-time.After(worker.timeout):
			if worker.loud {
				fmt.Printf("%s times out\n", worker.Name)
			}
		}
	} else {
		machine.Tasks <- task
		resolvedTask = <-machine.Tasks
	}

	if resolvedTask.BreakingTask {
		worker.repairRequest <- machine
	}

	return resolvedTask
}
