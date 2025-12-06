package concurrency

import "sync"

type Task func()

type WorkerPool struct {
    tasks chan Task
    wg    sync.WaitGroup
}

func NewWorkerPool(workerCount int, queueSize int) *WorkerPool {
    p := &WorkerPool{
        tasks: make(chan Task, queueSize),
    }
    for i := 0; i < workerCount; i++ {
        go func() {
            for t := range p.tasks {
                func() {
                    defer func() { if r := recover(); r != nil { /* noop */ } }()
                    t()
                }()
                p.wg.Done()
            }
        }()
    }
    return p
}

func (p *WorkerPool) Submit(t Task) {
    p.wg.Add(1)
    p.tasks <- t
}

func (p *WorkerPool) Wait() {
    p.wg.Wait()
}

func (p *WorkerPool) Close() {
    close(p.tasks)
}
