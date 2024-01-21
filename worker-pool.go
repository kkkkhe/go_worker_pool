package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	fn   func(interface{}) error
	args interface{}
}

type Worker struct {
	wg   *sync.WaitGroup
	tq   chan *Task
	quit chan bool
}

type Pool struct {
	size    int
	tq      chan *Task
	runInBg chan bool
}

func createTask(fn func(interface{}) error, args interface{}) *Task {
	return &Task{fn, args}
}
func createWorker(wg *sync.WaitGroup, taskQueue chan *Task) *Worker {
	return &Worker{wg, taskQueue, make(chan bool)}
}
func createPool(size int) *Pool {
	return &Pool{size, make(chan *Task), make(chan bool)}
}

func (w Worker) Run() {

	go func() {
		for {
			select {
			case task, ok := <-w.tq:
				if ok {
					task.fn(task.args)
				}
			case <-w.quit:
				return
			}
		}
	}()
}

func (p Pool) NewPool() {
	wg := &sync.WaitGroup{}
	for s := 0; s < p.size; s++ {
		worker := createWorker(wg, p.tq)
		worker.Run()
	}
	p.runInBg = make(chan bool)
	<-p.runInBg
}
func testFunc(id interface{}) error {
	fmt.Printf("Execute function %d\n", id)
	time.Sleep(300 * time.Millisecond)
	if id == 6 {
		return errors.New("Error occured")
	}
	return nil
}

func (p Pool) AddFunc(task *Task) {
	p.tq <- task
}

func main() {
	pool := createPool(5)
	go func() {
		for {
			taskID := rand.Intn(100)
			task := createTask(testFunc, taskID)
			pool.AddFunc(task)
		}
	}()
	pool.NewPool()
}
