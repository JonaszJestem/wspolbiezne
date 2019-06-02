package tasks

import "fmt"

type RepairService struct {
	BrokenMachines    []*Machine
	RepairingMachines []*Machine
}

type RepairRequest struct {
	BrokenMachine chan *Machine
}

func RunRepairService(service *RepairService, brokenMachines <-chan *Machine, repairRequest chan RepairRequest, machineRepaired <-chan *Machine) {
	for {
		select {
		case brokenMachine := <-brokenMachines:
			alreadyRequested := false
			for _, broken := range service.BrokenMachines {
				if brokenMachine.Name == broken.Name {
					alreadyRequested = true
					fmt.Println("Already requested!")
					break
				}
			}
			for _, broken := range service.RepairingMachines {
				if brokenMachine.Name == broken.Name {
					alreadyRequested = true
					fmt.Println("Already being repaired")
					break
				}
			}
			if !alreadyRequested {
				service.BrokenMachines = append(service.BrokenMachines, brokenMachine)
			}
		case repairRequest := <-repairGuard(len(service.BrokenMachines) > 0, repairRequest):
			machineToBeRepaired := service.BrokenMachines[0]
			repairRequest.BrokenMachine <- machineToBeRepaired
			service.RepairingMachines = append(service.RepairingMachines, machineToBeRepaired)
			service.BrokenMachines = service.BrokenMachines[1:]
		case repairedMachine := <-machineRepaired:
			for i, v := range service.RepairingMachines {
				if v == repairedMachine {
					service.RepairingMachines = append(service.RepairingMachines[:i], service.RepairingMachines[i+1:]...)
					break
				}
			}
			fmt.Println("Repaired confirm")
		}
	}
}

func repairGuard(condition bool, channel chan RepairRequest) chan RepairRequest {
	if condition {
		return channel
	}
	return nil
}

func CreateRepairService(brokenMachines <-chan *Machine, repairRequests chan RepairRequest) {
	service := RepairService{make([]*Machine, 0), make([]*Machine, 0)}
	RunRepairService(&service, brokenMachines, repairRequests)
}
