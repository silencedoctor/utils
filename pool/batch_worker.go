package pool

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

// Error ...
type Error struct {
	Index int64
	Err   error
}

// Error ...
func (e *Error) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

type BatchWorker struct {
	pool        *Pool
	workerIndex int64
	workerWg    sync.WaitGroup

	handleErrWg    sync.WaitGroup
	handleErrStart sync.Once

	errChan chan error
	errs    []error
}

// NewBatchWorker ...
func NewBatchWorker(pool *Pool) *BatchWorker {
	errChan := make(chan error, 64)
	errs := make([]error, 0, 64)

	return &BatchWorker{
		pool:    pool,
		errChan: errChan,
		errs:    errs,
	}
}

// Do ...
func (bw *BatchWorker) Do(t Task) {
	bw.handleErrStart.Do(func() {
		bw.handleErrWg.Add(1)
		go bw.handleError()
	})

	idx := atomic.AddInt64(&bw.workerIndex, 1) - 1
	bw.workerWg.Add(1)

	bw.pool.Submit(func() (err error) {
		defer func() {
			if pn := recover(); pn != nil {
				err = errors.New(fmt.Sprintf("worker panic recovered: %v", pn))
			}
			bw.workerWg.Done()
		}()

		err = t()
		if err != nil {
			err := &Error{
				Index: idx,
				Err:   err,
			}
			bw.errChan <- err
			return err
		}
		return nil
	})
}

func (bw *BatchWorker) handleError() {
	for {
		select {
		case err, ok := <-bw.errChan:
			if !ok {
				bw.handleErrWg.Done()
				return
			}
			bw.errs = append(bw.errs, err)
		}
	}
}

// Wait ...
func (bw *BatchWorker) Wait() []error {
	bw.workerWg.Wait()

	close(bw.errChan)
	bw.handleErrWg.Wait()

	return bw.errs
}
