package app

import (
	"context"
	"sync"
)

type JobResult struct {
	Err error
	Ok  bool
}

type App struct {
	wg   sync.WaitGroup
	done chan struct{}
}

// use pointer receiver because we cannot copy sync.waitgroup struct
func New() *App {
	return &App{
		wg:   sync.WaitGroup{},
		done: make(chan struct{}),
	}
}

func (a *App) AddJob(ctx context.Context, f func(ctx context.Context) error) <-chan JobResult {
	a.wg.Add(1)
	result := make(chan JobResult)

	go func() {
		err := f(ctx)
		result <- JobResult{
			Err: err,
			Ok:  err == nil,
		}
		a.wg.Done()
	}()

	return result
}

// wait until all jobs complete
func (a *App) Wait(ctx context.Context) {
	go func() {
		a.wg.Wait()
		a.done <- struct{}{}
	}()

	select {
	case <-a.done:
		return
	case <-ctx.Done():
		return
	}
}
