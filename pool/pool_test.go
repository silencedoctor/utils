package pool_test

import (
	"context"
	"fmt"
	"my/utils/pool"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	p, err := pool.NewPool(pool.WithCapacity(1000), pool.WithExpiryDuration(time.Second*30))
	if err != nil {
		fmt.Printf("new pool error ")
		return
	}

	ctx := context.Background()
	for i := 0; i < 10000; i++ {
		p.Submit(taskFmt(ctx, i))
	}

	time.Sleep(30 * time.Second)
}

func TestPoolClear(t *testing.T) {
	p, err := pool.NewPool(pool.WithCapacity(100), pool.WithExpiryDuration(time.Second*20))
	if err != nil {
		fmt.Printf("new pool error ")
		return
	}

	ctx := context.Background()

	for i := 0; i < 1000; i++ {
		p.Submit(taskFmt(ctx, i))
	}

	time.Sleep(10 * time.Second)
	for i := 0; i < 50; i++ {
		p.Submit(taskFmt(ctx, i))
	}
	fmt.Printf("pool is ok running %v \n", p.Running())

	time.Sleep(10 * time.Second)
	for i := 0; i < 25; i++ {
		p.Submit(taskFmt(ctx, i))
	}
	fmt.Printf("pool is ok running %v \n", p.Running())

	time.Sleep(10 * time.Second)
	for i := 0; i < 10; i++ {
		p.Submit(taskFmt(ctx, i))
	}
	fmt.Printf("pool is ok running %v \n", p.Running())

	time.Sleep(10 * time.Second)
	fmt.Printf("pool is ok running %v \n", p.Running())
	time.Sleep(10 * time.Second)
	fmt.Printf("pool is ok running %v \n", p.Running())
	time.Sleep(10 * time.Second)
	fmt.Printf("pool is ok running %v \n", p.Running())
	time.Sleep(10 * time.Second)
	fmt.Printf("pool is ok running %v \n", p.Running())
}

func doTaskFmt(ctx context.Context, i int) error {
	fmt.Printf("%v \n", i)
	time.Sleep(time.Second)
	return nil
}

// taskFmt ...
func taskFmt(ctx context.Context, idx int) pool.Task {
	return func() error {
		return doTaskFmt(ctx, idx)
	}
}

func TestBatchWorker(t *testing.T) {
	p, err := pool.NewPool(pool.WithCapacity(100), pool.WithExpiryDuration(time.Second*30))
	if err != nil {
		fmt.Printf("new pool error ")
		return
	}

	ctx := context.Background()

	bw := pool.NewBatchWorker(p)
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	for i := 0; i < len(data); i++ {
		idx := data[i]
		bw.Do(taskFmtErr(ctx, idx))
	}
	errs := bw.Wait()

	for _, e := range errs {
		ee := e.(*pool.Error)
		fmt.Printf("%#v %v %v\n", ee.Index, ee.Err.Error(), data[ee.Index])
	}
	fmt.Println("hhhhhhh -------->", errs)
	time.Sleep(1 * time.Second)
}

func doTaskFmtErr(ctx context.Context, i int) error {
	fmt.Printf("%v \n", i)
	if i%5 == 0 {
		return fmt.Errorf("i %v", i)
	}
	return nil
}

// taskFmt ...
func taskFmtErr(ctx context.Context, idx int) pool.Task {
	return func() error {
		return doTaskFmtErr(ctx, idx)
	}
}

func TestPoolClose(t *testing.T) {
	p, err := pool.NewPool(pool.WithCapacity(100), pool.WithExpiryDuration(time.Second*30))
	if err != nil {
		fmt.Printf("new pool error \n")
		return
	}

	ctx := context.Background()

	for i := 0; i < 10000; i++ {
		p.Submit(taskFmt(ctx, i))
	}

	fmt.Printf("pool woker num %v \n", p.Running())
	time.Sleep(10 * time.Second)
	fmt.Printf("pool woker num %v \n", p.Running())

	p.Close()
	fmt.Printf("pool woker num %v \n", p.Running())
	time.Sleep(10 * time.Second)
	fmt.Printf("pool woker num %v \n", p.Running())
}