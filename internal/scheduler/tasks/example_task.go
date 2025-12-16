package tasks

import "go-api-starter/internal/scheduler"

// ExampleTask untuk demo
type ExampleTask struct {
	name string
}

func NewExampleTask() scheduler.Task {
	return &ExampleTask{
		name: "example-task",
	}
}

func (t *ExampleTask) Name() string {
	return t.name
}

func (t *ExampleTask) Execute() error {
	// Contoh: log something, cleanup, atau generate report
	// logger.Infof("Running %s at %v", t.Name(), time.Now())
	return nil
}
