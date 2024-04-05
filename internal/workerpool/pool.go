package workerpool

import (
	"log"
	"sync"
	"time"
)

type Job struct {
	URL string
}

type Pool struct {
	worker       *worker
	workersCount int
	jobs         chan Job
	results      chan Result
	wg           *sync.WaitGroup
	stopped      bool
}

func New(workersCount int, timeout time.Duration, results chan Result) *Pool {
	return &Pool{
		worker:       newWorker(timeout),
		workersCount: workersCount,
		jobs:         make(chan Job),
		results:      results,
		wg:           new(sync.WaitGroup),
	}
}

func (p *Pool) initWorker(id int) {
	for job := range p.jobs {
		time.Sleep(time.Second)
		p.results <- p.worker.process(job)
		p.wg.Done()
	}

	log.Printf("[worker ID %d] finished proccesing", id)
}

func (p *Pool) Init() {
	for i := 0; i < p.workersCount; i++ {
		go p.initWorker(i)
	}
}

func (p *Pool) Push(j Job) {
	if p.stopped {
		return
	}

	p.jobs <- j
	p.wg.Add(1)
}

func (p *Pool) Stop() {
	p.stopped = true
	close(p.jobs)
	p.wg.Wait()
}
