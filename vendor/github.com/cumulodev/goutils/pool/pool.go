// package pool is a simple and typesafe implementation of an executioner pool.
package pool

import (
	"sync"
)

// A Job represents the client job that will be executed by a worker of the pool.
type Job interface {
	// Work is called to concurrently / parallelly calculate the result
	// of a job. To store results, add the appropriate fields to your
	// Job backing data structure and save them there.
	Work()

	// Save is called only from single goroutine and only one Save() is active
	// at any time for a single pool. This can be used to get the results out
	// of the job data structure and store them somewhere else.
	Save()
}

// A Pool represents an executioner pool of workers.
type Pool struct {
	size int
	wg   *sync.WaitGroup
	in   chan Job
	out  chan Job
}

// New creates a new Pool with the given number of worker goroutines. The
// pool will not be started until Start() is called.
func New(size int) *Pool {
	return &Pool{
		size: size,
		wg:   new(sync.WaitGroup),
	}
}

// Start starts the workers and colelctor goroutines of this pool.
func (p *Pool) Start() {
	p.in = make(chan Job)
	p.out = make(chan Job)

	go p.collector()
	for i := 0; i < p.size; i++ {
		go p.worker()
	}
}

// Stop stops the pool.
func (p *Pool) Stop() {
	close(p.in)
	close(p.out)
}

// Wait waits for the pool to complete all jobs.
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Add adds a job to the pool. It blocks until a worker can take care of the
// job. This functions should only be called after the pool is started!
func (p *Pool) Add(j Job) {
	p.wg.Add(1)
	p.in <- j
}

func (p *Pool) worker() {
	for j := range p.in {
		j.Work()
		p.out <- j
	}
}

func (p *Pool) collector() {
	for j := range p.out {
		j.Save()
		p.wg.Done()
	}
}
