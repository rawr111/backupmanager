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

	for _, job := range usecase.BackupFiles(ctx) {
		jobResult := app.AddJob(ctx, job.Func)
		go func() {
			res := <-jobResult
			if !res.Ok {
				fmt.Printf("job '%s' has been done with error: %v\n", job.Name, res.Err)
			} else {
				fmt.Printf("job '%s' has been done\n", job.Name)
			}
		}()
	}

	fmt.Println("wait...")
	app.Wait(ctx)
	fmt.Println("done")
}
