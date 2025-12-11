package service

import (
    "context"
    "sync"
)

// Job represents a unit of work
type Job struct {
    ID      int
    Payload interface{}
}

// Result represents the result of a job
type Result struct {
    Job   Job
    Value interface{}
    Error error
}

// WorkerPool manages a pool of workers for concurrent processing
type WorkerPool struct {
    workerCount int
    jobs        chan Job
    results     chan Result
    done        chan struct{}
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workerCount int) *WorkerPool {
    return &WorkerPool{
        workerCount: workerCount,
        jobs:        make(chan Job, 100),
        results:     make(chan Result, 100),
        done:        make(chan struct{}),
    }
}

// Start initializes and starts the worker pool
func (wp *WorkerPool) Start(ctx context.Context) {
    var wg sync.WaitGroup

    // Start workers
    for i := 0; i < wp.workerCount; i++ {
        wg.Add(1)
        go wp.worker(ctx, i, &wg)
    }

    // Wait for all workers to finish
    go func() {
        wg.Wait()
        close(wp.results)
        close(wp.done)
    }()
}

// worker processes jobs from the jobs channel
func (wp *WorkerPool) worker(ctx context.Context, id int, wg *sync.WaitGroup) {
    defer wg.Done()

    for {
        select {
        case job, ok := <-wp.jobs:
            if !ok {
                return
            }
            // Process job (implementation depends on use case)
            result := Result{
                Job:   job,
                Value: job.Payload,
                Error: nil,
            }
            wp.results <- result

        case <-ctx.Done():
            return
        }
    }
}

// Submit submits a job to the worker pool
func (wp *WorkerPool) Submit(job Job) {
    wp.jobs <- job
}

// Results returns the results channel
func (wp *WorkerPool) Results() <-chan Result {
    return wp.results
}

// Close closes the worker pool
func (wp *WorkerPool) Close() {
    close(wp.jobs)
}

// Done returns a channel that's closed when all workers are done
func (wp *WorkerPool) Done() <-chan struct{} {
    return wp.done
}
```

**Penjelasan:**
- Worker Pool pattern untuk concurrency
- `workerCount` = Jumlah worker concurrent (dari config)
- Jobs channel = Queue untuk tasks
- Results channel = Output dari workers
- Context-aware untuk cancellation

**Worker Pool Architecture:**
```
    //              Jobs Queue
    //                  │
    //     ┌────────────┼────────────┐
    //     ▼            ▼            ▼
    // Worker 1     Worker 2    Worker 3  ... Worker N
    //     │            │            │
    //     └────────────┼────────────┘
    //                  ▼
    //           Results Channel

// Bulk operations
// pool := NewWorkerPool(10)
// pool.Start(ctx)

// for _, user := range users {
//     pool.Submit(Job{ID: i, Payload: user})
// }

// for i := 0; i < len(users); i++ {
//     result := <-pool.Results()
//     // Handle result
// }

// pool.Close()