package usecases

import (
	"context"
	"fmt"
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

type Job struct {
	Func func(ctx context.Context) error
	Name string
}

func (u *Usecase) BackupFiles(ctx context.Context) []Job {
	filesToCopy := []file{
		{
			"/Users/gladkikhmikhail/Downloads/250.zip",
			"some_big_file",
		},
	}

	result := []Job{}
	for _, file := range filesToCopy {
		result = append(result, Job{
			Func: func(ctx context.Context) error {
				return u.mngr.BackupFile(ctx, file.path, file.alias)
			},
			Name: fmt.Sprintf("backup %s file as %s", file.path, file.alias),
		})
	}

	return result
}
