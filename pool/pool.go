package pool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Task func() error

type Pool struct {
	ctx context.Context

	cancel context.CancelFunc

	running int64 // 当前正在运行的协程数量

	options *options

	workers []*Worker // worker 集合

	lock sync.Mutex // 锁

	isClose bool
}

func NewPool(options ...Option) (*Pool, error) {
	return NewPoolWithCtx(context.Background(), options...)
}

func NewPoolWithCtx(ctx context.Context, options ...Option) (*Pool, error) {
	opts := newDefaultOptions()
	for _, opt := range options {
		opt.apply(opts)
	}

	if opts.capacity <= 0 {
		return nil, errors.New(fmt.Sprintf("capacity is too small capacity:%v", opts.capacity))
	}

	ctx, cancel := context.WithCancel(ctx)
	p := &Pool{
		ctx:     ctx,
		cancel:  cancel,
		running: 0,
		workers: make([]*Worker, 0),
		lock:    sync.Mutex{},
		options: opts,
	}

	go p.periodicallyClear()

	return p, nil
}

func (p *Pool) Submit(task Task) error {
	if p.isClose {
		return fmt.Errorf("pool is close")
	}

	w := p.getWorker()
	for w == nil { // w 为nil 说明在task被发送之前刚好被超时清理掉了，重新获取worker
		w = p.getWorker()
	}
	w.task <- task

	return nil
}

func (p *Pool) getWorker() *Worker {
	w, waiting := p.getIdleWorkers()

	if waiting { // 当前运行的协程已满
		for {
			w, waiting = p.getIdleWorkers()
			if waiting {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}
	} else if w == nil {
		w = &Worker{
			pool: p,
			task: make(chan Task, 1),
		}
		w.work()
		p.incRunning()
	}

	return w
}

// Running 获取协程池中正在运行的worker数目
func (p *Pool) Running() int64 {
	return p.running
}

// Capacity 获取协程池设定的大小
func (p *Pool) Capacity() int64 {
	return p.options.capacity
}

// incRunning 正在运行的worker数目加1
func (p *Pool) incRunning() {
	atomic.AddInt64(&p.running, 1)
}

// decRunning 正在运行的worker数目减1
func (p *Pool) decRunning() {
	atomic.AddInt64(&p.running, -1)
}

// getIdleWorkers 获取空闲worker
func (p *Pool) getIdleWorkers() (*Worker, bool) {
	p.lock.Lock()
	idleWorkers := p.workers
	n := len(idleWorkers) - 1 // 当前队列中可用worker的数量
	if n <= 0 {
		waiting := p.Capacity() <= p.Running()
		p.lock.Unlock()

		return nil, waiting
	}

	work := idleWorkers[n] // 获取当前队列最后一位， 并将其从队列中剔除
	idleWorkers[n] = nil
	p.workers = idleWorkers[:n]
	p.lock.Unlock()

	return work, false
}

// recycleWorker 回收执行完成的work到空闲列表
func (p *Pool) recycleWorker(worker *Worker) {
	p.lock.Lock()
	worker.recycleTime = time.Now()
	p.workers = append(p.workers, worker)
	p.lock.Unlock()
}

func (p *Pool) periodicallyClear() {
	timeTicker := time.NewTicker(p.options.expiryDuration)
	defer timeTicker.Stop()

	for {
		select {
		case <-timeTicker.C:
			currentTime := time.Now()
			p.lock.Lock()
			idleWorkers := p.workers
			i := 0
			for ; i < len(idleWorkers); i++ {
				if currentTime.Sub(idleWorkers[i].recycleTime) <= p.options.expiryDuration {
					break
				}

				idleWorkers[i].task <- nil
				idleWorkers[i] = nil
			}

			p.workers = idleWorkers[i:]
			p.lock.Unlock()
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Pool) Close() {
	p.lock.Lock()
	p.isClose = true
	p.cancel()
	p.lock.Unlock()
}
