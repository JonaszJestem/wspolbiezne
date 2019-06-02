package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	. "wspolbiezne/company"
	. "wspolbiezne/company/shop"
	. "wspolbiezne/company/tasks"
)

func main() {
	rand.Seed(time.Now().Unix())

	addProduct, sellProduct, getProducts := createShop()
	addTask, getTask, getTasks := createStorage()
	brokenMachines := make(chan *Machine)
	machineRepaired := make(chan *Machine)
	repairRequests := make(chan RepairRequest)
	createRepairWorkers(repairRequests, machineRepaired)

	createBoss(addTask)
	go CreateRepairService(brokenMachines, repairRequests, machineRepaired)

	var addingMachines = createMachines("Adder_", NumberOfAdders, AdderWorkingTime)
	var multiplyingMachines = createMachines("Multiplier_", NumberOfMultipliers, MultiplierWorkingTime)
	var workers = createWorkers(brokenMachines, getTask, addProduct, addingMachines, multiplyingMachines, NumberOfWorkers)
	createCustomers(sellProduct)

	if Loud {
		wait()
	} else {
		getCommands(getProducts, getTasks, workers)
	}
}

func createRepairWorkers(repairRequests chan RepairRequest, machineRepaired chan *Machine) {
	for i := 0; i < NumberOfRepairWorkers; i++ {
		worker := CreateRepairWorker("mechanic_"+strconv.Itoa(i), repairRequests, RepairTime, Loud, machineRepaired)
		go StartRepairWorker(worker)
	}
}

func createMachines(namePrefix string, numberOfMachines int, sleepTime time.Duration) []*Machine {
	var machines []*Machine
	for i := 0; i < numberOfMachines; i++ {
		var tasks = make(chan Task)
		var machine = CreateMachine(namePrefix+strconv.Itoa(i), tasks, sleepTime, MachineReliability, Loud)
		machines = append(machines, &machine)
		go StartMachine(machine)
	}
	return machines
}

func getCommands(getProducts chan bool, getTasks chan bool, workers []*Worker) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter command: ")
		fmt.Println("t - list of tasks")
		fmt.Println("s - list of products")
		fmt.Println("w - list of workers")

		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("t", text) == 0 {
			getTasks <- true
		} else if strings.Compare("s", text) == 0 {
			getProducts <- true
		} else if strings.Compare("w", text) == 0 {
			for _, worker := range workers {
				fmt.Printf("%s has done %d tasks and is patient: %t\n", worker.Name, worker.TasksDone, worker.Patient)
			}
		}
	}
}

func createCustomers(sellProduct chan Offer) {
	for i := 0; i < NumberOfCustomers; i++ {
		var customer = CreateCustomer(fmt.Sprintf("customer_%d", i), sellProduct, Loud)
		go StartCustomer(customer)
	}
}

func createWorkers(repairRequests chan *Machine, getTask chan TaskRequest, addProduct chan Product, addingMachines []*Machine, multiplyingMachines []*Machine, howMany int) []*Worker {
	var workers []*Worker
	for i := 0; i < howMany; i++ {
		var shouldBePatient = rand.Float32() < 0.5
		var worker = CreateWorker(fmt.Sprintf("worker_%d", i), repairRequests, getTask, addProduct, addingMachines, multiplyingMachines, shouldBePatient, WorkerTimeout, Loud)
		workers = append(workers, &worker)
		go StartWorker(&worker)
	}
	return workers
}

func createBoss(addTask chan Task) {
	var boss = CreateBoss("boss", addTask, Loud)
	go StartBossing(boss)
}

func createStorage() (chan Task, chan TaskRequest, chan bool) {
	var storage = CreateStorage(StorageSize)
	var addTask = make(chan Task)
	var getTask = make(chan TaskRequest)
	var getTasks = make(chan bool)
	go RunStorage(storage, addTask, getTask, getTasks)
	return addTask, getTask, getTasks
}

func createShop() (chan Product, chan Offer, chan bool) {
	var shop = CreateShop(ShopSize)
	var addProduct = make(chan Product)
	var sellProduct = make(chan Offer)
	var getProducts = make(chan bool)
	go RunShop(shop, addProduct, sellProduct, getProducts)
	return addProduct, sellProduct, getProducts
}

func wait() {
	for {
		time.Sleep(time.Second)
	}
}
