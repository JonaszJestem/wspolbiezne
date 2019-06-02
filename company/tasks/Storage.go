package tasks

import "fmt"

type Storage struct {
	Tasks       []Task
	StorageSize int
}

func RunStorage(storage Storage, addTask <-chan Task, getTask <-chan TaskRequest, getTasks <-chan bool) {
	for {
		select {
		case newTask := <-newTaskGuard(hasSpace(storage), addTask):
			storage.Tasks = append(storage.Tasks, newTask)
		case taskRequest := <-taskRequestGuard(len(storage.Tasks) > 0, getTask):
			taskRequest.Task <- storage.Tasks[0]
			storage.Tasks = storage.Tasks[1:]
		case <-getTasks:
			for _, task := range storage.Tasks {
				fmt.Println(task)
			}
		}
	}
}

func newTaskGuard(condition bool, channel <-chan Task) <-chan Task {
	if condition {
		return channel
	}
	return nil
}

func taskRequestGuard(condition bool, channel <-chan TaskRequest) <-chan TaskRequest {
	if condition {
		return channel
	}
	return nil
}

func CreateStorage(storageSize int) Storage {
	return Storage{make([]Task, 0, storageSize), storageSize}
}

func hasSpace(storage Storage) bool {
	return len(storage.Tasks) < storage.StorageSize
}
