package app

import (
	"context"
	"sync"
	"testing"
	"time"
)

type testOperation func(ctx context.Context) error

type jobMaker struct {
	mu           sync.RWMutex
	countJobDone int
}

func (j *jobMaker) makeLongTestOpeartion(delaySeconds int) testOperation {
	f := func(ctx context.Context) error {
		t := time.NewTimer(time.Duration(delaySeconds) * time.Second)
		select {
		case <-ctx.Done():
		case <-t.C:
		}

		j.mu.RLock()
		j.countJobDone++
		j.mu.RUnlock()

		return nil
	}
	return f
}

func TestCuncurrentJobsProcessing(t *testing.T) {
	app := New()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	jm := &jobMaker{}

	jobs := []func(ctx context.Context) error{
		jm.makeLongTestOpeartion(1),
		jm.makeLongTestOpeartion(2),
		jm.makeLongTestOpeartion(1),
		jm.makeLongTestOpeartion(2),
		jm.makeLongTestOpeartion(3),
		jm.makeLongTestOpeartion(1),
		jm.makeLongTestOpeartion(2),
		jm.makeLongTestOpeartion(1),
		jm.makeLongTestOpeartion(1),
		jm.makeLongTestOpeartion(1),
		jm.makeLongTestOpeartion(1),
		jm.makeLongTestOpeartion(1),
	}

	go func() {
		time.Sleep(5 * time.Second)
	}()

	for _, j := range jobs {
		c := app.AddJob(ctx, j)
		go func() {
			<-c
		}()
	}

	c := make(chan struct{})
	go func() {
		app.Wait(ctx)
		c <- struct{}{}
	}()

	select {
	case <-c:
		if jm.countJobDone != len(jobs) {
			t.Fatalf("not all jobs were procceded")
		}
		return
	case <-ctx.Done():
		t.Fatalf("jobs are not ended")
	}
}
