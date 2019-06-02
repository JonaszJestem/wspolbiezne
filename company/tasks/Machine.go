package tasks

import (
	"fmt"
	"math/rand"
	"time"
)

type Machine struct {
	Name        string
	Tasks       chan Task
	sleepTime   time.Duration
	loud        bool
	broken      bool
	reliability float32
	backdoor    chan bool
}

func CreateMachine(name string, tasks chan Task, sleepTime time.Duration, reliability float32, loud bool) Machine {
	return Machine{name, tasks, sleepTime, loud, false, reliability, make(chan bool)}
}

func StartMachine(machine Machine) {
	rand.Seed(time.Now().Unix())

	for {
		select {
		case task := <-machineGuard(!machine.broken, machine.Tasks):
			time.Sleep(machine.sleepTime)
			result, err := getResult(task, machine.reliability)
			if !machine.broken && err == nil {
				task.Successful = true
				task.Result = result
				if machine.loud {
					fmt.Printf("%s produces %d\n", machine.Name, task.Result)
				}
				machine.Tasks <- task
			} else {
				if machine.loud {
					fmt.Printf("%s is broken\n", machine.Name)
				}
				machine.broken = true
				task.BreakingTask = true
				machine.Tasks <- task
			}
		case <-machine.backdoor:
			machine.broken = false
		}
	}
}

func machineGuard(condition bool, channel chan Task) chan Task {
	if condition {
		return channel
	}
	return nil
}

func getResult(task Task, reliability float32) (int, error) {
	if rand.Float32() > reliability {
		return 0, fmt.Errorf("machine is broken")
	}

	switch operation := task.Operation; operation {
	case "+":
		return task.Left + task.Right, nil
	case "*":
		return task.Left * task.Right, nil
	default:
		return 0, fmt.Errorf("unrecognized operation %s", task.Operation)
	}
}
