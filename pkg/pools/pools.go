package pools

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Pool is a worker group that runs a number of tasks at a
// configured concurrency.
type Pool struct {
	Tasks       []*Task
	timeWork    time.Duration
	concurrency int
	tasksChan   chan *Task
	wg          sync.WaitGroup
	quitChannel chan os.Signal
	//shutdownChannel chan struct{}
}

// NewPool initializes a new pool with the given tasks and
// at the given concurrency.
func NewPool(tasks []*Task, concurrency int, durationTime int64) (*Pool, error) {
	return &Pool{
		Tasks:       tasks,
		timeWork:    time.Duration(durationTime) * time.Microsecond,
		concurrency: concurrency,
		tasksChan:   make(chan *Task, 1),
	}, nil
}

func (p *Pool) Stop() {
	_, ok := <-p.tasksChan
	if ok {
		close(p.tasksChan)
	}
}

func (p *Pool) initSignal() {
	quitChannel := make(chan os.Signal, 3)
	signal.Notify(quitChannel, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)
	p.quitChannel = quitChannel
}

func (p *Pool) Signal() {
	for {
		s := <-p.quitChannel
		switch s {
		case os.Interrupt, syscall.SIGTERM, syscall.SIGKILL:
			fmt.Println("Stop called, closing quit channel")
			p.Stop()
			fmt.Println("Done")
			os.Exit(1)
			return
		}
	}
}

// Run runs all work within the pool and blocks until it's
// finished.
func (p *Pool) Run() {
	p.initSignal()

	for i := 0; i < p.concurrency; i++ {
		go p.work()
	}

	p.wg.Add(len(p.Tasks))
	for _, task := range p.Tasks {
		p.tasksChan <- task
	}

	go p.Signal()
	close(p.tasksChan)
	p.wg.Wait()
}

// The work loop for any single goroutine.
func (p *Pool) work() {
	for task := range p.tasksChan {
		task.Run(&p.wg)
		time.Sleep(p.timeWork)
	}
}
