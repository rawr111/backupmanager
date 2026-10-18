package usecases

import (
	"context"
	"fmt"
	"time"
)

type BackupMngr interface {
	BackupFile(ctx context.Context, filePath string, alias string) error
}

type Usecase struct {
	mngr BackupMngr
}

func New(mngr BackupMngr) *Usecase {
	return &Usecase{
		mngr: mngr,
	}
}

type file struct {
	path  string
	alias string
}

func (u *Usecase) BackupJob(ctx context.Context) error {
	// filesToCopy := []file{
	// 	{
	// 		"/Users/gladkikhmikhail/Downloads/250.zip",
	// 		"some_big_file",
	// 	},
	// }
	fmt.Println("sleeep...")
	time.Sleep(time.Hour * 10)
	return u.mngr.BackupFile(ctx, "", "")
}
