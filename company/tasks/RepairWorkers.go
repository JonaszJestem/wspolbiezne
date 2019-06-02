package tasks

import (
	"fmt"
	"time"
)

type RepairWorker struct {
	name            string
	RepairRequest   chan RepairRequest
	repairTime      time.Duration
	loud            bool
	MachineRepaired chan *Machine
}

func CreateRepairWorker(name string, repairRequest chan RepairRequest, repairTime time.Duration, loud bool, machineRepaired chan *Machine) RepairWorker {
	return RepairWorker{name, repairRequest, repairTime, loud, machineRepaired}
}

func StartRepairWorker(worker RepairWorker) {
	for {
		brokenMachines := make(chan *Machine)
		worker.RepairRequest <- RepairRequest{brokenMachines}
		brokenMachine := <-brokenMachines
		time.Sleep(worker.repairTime)
		brokenMachine.backdoor <- true
		worker.MachineRepaired <- brokenMachine
		if worker.loud {
			fmt.Printf("%s repaired \t%v\n", worker.name, brokenMachine.Name)
		}
	}
}
