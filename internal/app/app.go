package app

import (
	"context"
	"sync"
)

type App struct {
	wg     sync.WaitGroup
	errors chan error
}

func (a *App) Result(ctx context.Context) <-chan error {
	done := ctx.Done()
	go func() {
		<-done
		close(a.errors)
	}()

	return a.errors
}

func (a *App) AddJob(ctx context.Context, f func(ctx context.Context) error) {
	a.wg.Add(1)

	go func() {
		err := f(ctx)
		a.wg.Done()
		a.errors <- err
	}()
}

// wait until all jobs complete
func (a *App) Wait(ctx context.Context) {
	<-ctx.Done()
	a.wg.Wait()
}

// use pointer receiver because we cannot copy sync.waitgroup struct
func New() *App {
	return &App{
		errors: make(chan error, 1000),
		wg:     sync.WaitGroup{},
	}
}
