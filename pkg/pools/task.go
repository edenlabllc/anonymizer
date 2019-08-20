package pools

import (
	"sync"
)

type Message struct {
	ID       string
	Data     string
	Callback string
	Object   interface{}
	Table    interface{}
}

// Task encapsulates a work item that should go in a work
// pool.
type Task struct {
	Err              error
	App              interface{}
	Message          *Message
	f                func(task *Task) error
	SuccessIteration int
}

// NewTask initializes a new task based on a given work
// function.
func NewTask(f func(task *Task) error, message *Message, app interface{}) *Task {
	return &Task{Message: message, f: f, App: app, SuccessIteration: 0}
}

// Run runs a Task and does appropriate accounting via a
// given sync.WorkGroup.
func (t *Task) Run(wg *sync.WaitGroup) {
	t.Err = t.f(t)
	defer wg.Done()
}
