package main

import (
	application "backupmanager/internal/app"
	"backupmanager/internal/repositories/googledrive"
	"backupmanager/internal/services/backupmanager"
	"backupmanager/internal/usecases"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := application.New()
	googleDrive := googledrive.New()
	mngr := backupmanager.New(googleDrive)
	usecase := usecases.New(mngr)

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	app.AddJob(ctx, usecase.BackupJob)

	go func() {
		for err := range app.Result(ctx) {
			fmt.Printf("job run error: %s\n", err)
		}
	}()

	app.Wait(ctx)
}
