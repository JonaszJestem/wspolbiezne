package tasks

import (
	"fmt"
	"time"
)

type RepairWorker struct {
	name          string
	RepairRequest chan RepairRequest
	repairTime    time.Duration
	loud          bool
}

func CreateRepairWorker(name string, repairRequest chan RepairRequest, repairTime time.Duration, loud bool) RepairWorker {
	return RepairWorker{name, repairRequest, repairTime, loud}
}

func StartRepairWorker(worker RepairWorker) {
	for {
		brokenMachines := make(chan *Machine)
		worker.RepairRequest <- RepairRequest{brokenMachines}
		repairRequest := <-worker.RepairRequest
		brokenMachine := <-repairRequest.BrokenMachine
		time.Sleep(worker.repairTime)
		brokenMachine.backdoor <- true

		if worker.loud {
			fmt.Printf("%s repaired \t%v\n", worker.name, brokenMachine.Name)
		}
	}
}
