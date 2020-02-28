package grpool

import (
	"sync"
	"time"
)

// Gorouting instance which can accept client jobs
type worker struct {
	workerPool chan *worker
	jobChannel chan Job
	Jobresult  chan Jobresult
	jobtimeout time.Duration
	stop       chan bool
}

func (w *worker) start(pool *Pool) {
	go func() {
		var job Job
		for {
			// worker free, add it to pool
			w.workerPool <- w
			//fmt.Println("w.free", w.workerPool, "\n")
			select {
			case job = <-w.jobChannel:
				res, timeout, err := Do(w.jobtimeout, func() (interface{}, error) {
					res, err := job.Jobfunc()
					//fmt.Printf("job.Jobfunc", res, err, "\n")
					return res, err
				})
			//fmt.Printf("withtimeout.Do", res, err, "\n")
				w.Jobresult <- Jobresult{
					Jobid:    job.Jobid,
					Timedout: timeout,
					Result:   res,
					Err:      err,
				}
				pool.wg.Done()
			case stop := <-w.stop:
				//fmt.Printf("w.stop", stop, "\n")
				if stop {
					w.stop <- true
					return
				}
			}
		}
	}()
}

func newWorker(pool chan *worker, jobresult chan Jobresult) *worker {
	return &worker{
		workerPool: pool,
		jobChannel: make(chan Job),
		Jobresult:  jobresult,
		stop:       make(chan bool),
	}
}

// Accepts jobs from clients, and waits for first free worker to deliver job
type dispatcher struct {
	workerPool chan *worker
	jobQueue   chan Job
	Jobresult  chan Jobresult
	jobtimeout time.Duration
	stop       chan bool
}

func (d *dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			worker := <-d.workerPool
			worker.jobtimeout = d.jobtimeout
			worker.jobChannel <- job
		case stop := <-d.stop:
			//fmt.Printf("dispatcher.stop", stop, "\n")
			if stop {
				for i := 0; i < cap(d.workerPool); i++ {
					worker := <-d.workerPool
					worker.stop <- true
					<-worker.stop
				}
				d.stop <- true
				return
			}
		}
	}
}

func newDispatcher(workerPool chan *worker, jobQueue chan Job, jobresult chan Jobresult, timeout time.Duration ,pool *Pool) *dispatcher {
	d := &dispatcher{
		workerPool: workerPool,
		jobQueue:   jobQueue,
		Jobresult:  jobresult,
		jobtimeout: timeout,
		stop:       make(chan bool),
	}

	for i := 0; i < cap(d.workerPool); i++ {
		worker := newWorker(d.workerPool, d.Jobresult)
		worker.start(pool)
	}

	go d.dispatch()
	return d
}

type Jobresult struct {
	Jobid    interface{}
	Timedout bool
	Result   interface{}
	Err      error
}

// Represents user request, function which should be executed in some worker.
type Job struct {
	Jobid   interface{}
	Jobfunc func() (interface{}, error)
}

type Pool struct {
	JobQueue   chan Job
	dispatcher *dispatcher
	Jobresult  chan Jobresult
	wg         sync.WaitGroup
}


// Will make pool of gorouting workers.
// numWorkers - how many workers will be created for this pool
// queueLen - how many jobs can we accept until we block
//
// Returned object contains JobQueue reference, which you can use to send job to pool.
func NewPool(numWorkers int, jobQueueLen int, timeout time.Duration) *Pool {
	jobQueue := make(chan Job, jobQueueLen)
	workerPool := make(chan *worker, numWorkers)
	PoolJobresult := make(chan Jobresult, jobQueueLen)
	pool := &Pool{
		JobQueue:   jobQueue,
		Jobresult:  PoolJobresult,
	}
	pool.dispatcher=newDispatcher(workerPool, jobQueue, PoolJobresult, timeout,pool)
	return pool
}

// In case you are using WaitAll fn, you should call this method
// every time your job is done.
//
// If you are not using WaitAll then we assume you have your own way of synchronizing.
func (p *Pool) JobDone() {
	p.wg.Done()
}

// How many jobs we should wait when calling WaitAll.
// It is using WaitGroup Add/Done/Wait
func (p *Pool) WaitCount(count int) {
	if count != cap(p.JobQueue) {
		panic("WaitCount != jobQueueLen")
	}
	p.wg.Add(count)
}

// Will wait for all jobs to finish.
func (p *Pool) WaitAll() {
	p.wg.Wait()
	close(p.Jobresult)
}

// Will release resources used by pool
func (p *Pool) Release() {
	p.dispatcher.stop <- true
	<-p.dispatcher.stop
}
