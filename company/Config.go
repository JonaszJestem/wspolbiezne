package company

import "time"

const (
	//General
	Loud = false

	//Customer
	NumberOfCustomers int = 1
	newOfferInterval      = time.Second

	//Worker
	NumberOfWorkers       int = 3
	NumberOfAdders        int = 1
	AdderWorkingTime          = time.Second
	NumberOfMultipliers   int = 1
	MultiplierWorkingTime     = time.Second
	WorkerSleep               = time.Second
	WorkerTimeout             = time.Second * 3

	//Boss
	NewTaskInterval = time.Second / 3

	//Storages
	StorageSize int = 10
	ShopSize    int = 10

	//Repairs
	NumberOfRepairWorkers int     = 2
	RepairTime                    = time.Second * 1
	MachineReliability    float32 = 0.8
)
