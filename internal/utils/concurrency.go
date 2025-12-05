package utils

import (
	"sync"
)

type Task func()

type WorkerPool struct {
	TaskChan chan Task
	wg       sync.WaitGroup
}

func NewWorkerPool(workerCount int) *WorkerPool {
	pool := &WorkerPool{
		TaskChan: make(chan Task),
	}

	for i := 0; i < workerCount; i++ {
		go func() {
			for job := range pool.TaskChan {
				job()
				pool.wg.Done()
			}
		}()
	}

	return pool
}

func (p *WorkerPool) Submit(task Task) {
	p.wg.Add(1)
	p.TaskChan <- task
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
}
