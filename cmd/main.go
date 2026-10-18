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
	cfg, err := readConfig()
	if err != nil {
		panic(err)
	}

	app := application.New()
	googleDrive := googledrive.New(cfg.GoogleDriveApiKey)
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

type Config struct {
	GoogleDriveApiKey string
}

func readConfig() (Config, error) {
	var apikey string
	apikey, ok := os.LookupEnv("GOOGLE_DRIVE_API_KEY")

	if !ok {
		fmt.Println("enter google api key: ")
		_, err := fmt.Scan(&apikey)
		if err != nil {
			return Config{}, fmt.Errorf("google drive api key read: %w", err)
		}
		if len(apikey) == 0 {
			return Config{}, fmt.Errorf("google drive api key cannot be empty")
		}
	}

	return Config{
		GoogleDriveApiKey: apikey,
	}, nil
}
