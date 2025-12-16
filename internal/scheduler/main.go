package scheduler

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

type Task interface {
	Name() string
	Execute() error
}

type Scheduler struct {
	tasks   []Task
	logger  *zap.SugaredLogger
	mu      sync.RWMutex
	done    chan bool
	running bool
}

func NewScheduler(logger *zap.SugaredLogger) *Scheduler {
	return &Scheduler{
		tasks:  make([]Task, 0),
		logger: logger,
		done:   make(chan bool),
	}
}

// AddTask registers task dengan interval
func (s *Scheduler) AddTask(interval time.Duration, task Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks = append(s.tasks, task)

	// Launch goroutine untuk task ini
	go s.runTask(interval, task)
}

func (s *Scheduler) runTask(interval time.Duration, task Task) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.logger.Infof("Executing scheduled task: %s", task.Name())

			// Gunakan goroutine untuk prevent blocking
			go func() {
				if err := task.Execute(); err != nil {
					s.logger.Errorf("Task %s failed: %v", task.Name(), err)
				} else {
					s.logger.Infof("Task %s completed successfully", task.Name())
				}
			}()
		case <-s.done:
			s.logger.Infof("Stopping scheduled task: %s", task.Name())
			return
		}
	}
}

func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.running = true
	s.logger.Info("Scheduler started")
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	s.running = false
	close(s.done)
	s.logger.Info("Scheduler stopped")
}

func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}
