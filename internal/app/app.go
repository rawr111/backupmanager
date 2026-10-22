package app

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type JobResult struct {
	Err error
	Ok  bool
}

type job struct {
	id string
}

type App struct {
	currentJobs map[string]job
	mu          sync.RWMutex
	wg          sync.WaitGroup
	done        chan struct{}
}

// use pointer receiver because we cannot copy sync.waitgroup struct
func New() *App {
	return &App{
		wg:          sync.WaitGroup{},
		done:        make(chan struct{}),
		currentJobs: make(map[string]job),
	}
}

// how many jobs are currently running
func (a *App) JobCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return len(a.currentJobs)
}

func (a *App) AddJob(ctx context.Context, f func(ctx context.Context) error) <-chan JobResult {
	a.wg.Add(1)
	result := make(chan JobResult)
	jobId := uuid.NewString()

	a.mu.Lock()
	defer a.mu.Unlock()
	a.currentJobs[jobId] = job{jobId}

	go func() {
		err := f(ctx)
		result <- JobResult{
			Err: err,
			Ok:  err == nil,
		}
		a.mu.Lock()
		delete(a.currentJobs, jobId)
		a.mu.Unlock()

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
