package pool

import (
	"fmt"
	"time"
)

type Worker struct {
	pool *Pool // 归属协程池

	task chan Task //任务

	recycleTime time.Time // 超时时间
}

func (w *Worker) work() {
	go w.doWork()
}

func (w *Worker) doWork() {
	for {
		select {
		case t := <-w.task:
			func() {
				defer func() {
					if pn := recover(); pn != nil {
						fmt.Printf("worker panic recovered: %v", pn)
					}
				}()

				if t == nil {
					w.pool.decRunning()
					return
				}
				t()
				w.pool.recycleWorker(w)
			}()
		case <-w.pool.ctx.Done():
			w.pool.decRunning()
			return
		}
	}
}
