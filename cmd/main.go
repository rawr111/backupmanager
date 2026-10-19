package main

import (
	application "backupmanager/internal/app"
	"backupmanager/internal/repositories/googledrive"
	"backupmanager/internal/services/backupmanager"
	"backupmanager/internal/usecases"
	"context"
	"errors"
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

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := application.New()
	googleDrive, err := googledrive.New(ctx, cfg.GoogleDrive)
	//если что-то не так с параметрами при запуске приложения, то разумно выйти с паникой
	if errors.Is(err, googledrive.OptsValidationError) {
		panic(err)
	}
	mngr := backupmanager.New(googleDrive)
	usecase := usecases.New(mngr)

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
	GoogleDrive googledrive.Opts
}

// если что-то не так с параметрами при запуске приложения, то разумно выйти с паникой
func readConfig() (Config, error) {
	refreshToken, ok := os.LookupEnv("GOOGLE_DRIVE_REFRESH_TOKEN")
	if !ok {
		panic("no refreshToken for googledrive")
	}
	clientId, ok := os.LookupEnv("GOOGLE_DRIVE_CLIENT_ID")
	if !ok {
		panic("no clientId for googledrive")
	}
	secretId, ok := os.LookupEnv("GOOGLE_DRIVE_SECRET_ID")
	if !ok {
		panic("no secretId for googledrive")
	}

	return Config{
		GoogleDrive: googledrive.Opts{
			RefreshToken: refreshToken,
			ClientId:     clientId,
			SecretId:     secretId,
		},
	}, nil
}
