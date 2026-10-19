package googledrive

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type GoogleDrive struct {
	cli         http.Client
	accessToken string

	refreshToken string
	clientId     string
	clientSecret string
}

type Opts struct {
	RefreshToken string
	ClientId     string
	SecretId     string
}

var OptsValidationError = errors.New("google drive settings are not valid")

// use pointer receiver because we dont want copy http.client
func New(ctx context.Context, opts Opts) (*GoogleDrive, error) {
	if opts.ClientId == "" || opts.RefreshToken == "" || opts.SecretId == "" {
		return nil, fmt.Errorf("%w: clientid, refreshtoken or SecretId are empty", OptsValidationError)
	}

	drive := &GoogleDrive{
		clientId:     opts.ClientId,
		clientSecret: opts.SecretId,
		refreshToken: opts.RefreshToken,
		accessToken:  "",
		cli:          http.Client{},
	}

	accessToken, err := drive.refreshAccessToken(ctx)
	if err != nil {
		fmt.Printf("access token refresh error: %s\n", err)
	}
	drive.accessToken = accessToken

	go func() {
		ticker := time.NewTicker(time.Minute * 50)
		for {
			select {
			case <-ticker.C:
				accessToken, err := drive.refreshAccessToken(ctx)
				if err != nil {
					fmt.Printf("access token refresh error: %s\n", err)
				}
				drive.accessToken = accessToken
			case <-ctx.Done():
				return
			}
		}
	}()

	return drive, nil
}

const uploadFilePath string = "upload/drive/v3/files?uploadType=media"

func (drive *GoogleDrive) UploadFile(ctx context.Context, filePath string, alias string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file %s failed: %w", filePath, err)
	}
	defer file.Close()

	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	resp, err := drive.doHttpRequest(ctx,
		methodPost,
		uploadFilePath,
		file,
		headers{
			{
				key:   "Content-type",
				value: contentType,
			},
		},
	)
	if err != nil {
		return err
	}

	if resp.statusCode != 200 {
		return fmt.Errorf("upload google drive file: %s", resp.body.err.message)
	}

	return nil
}

func (drive *GoogleDrive) DownloadFile(fileName string) (error, []byte) {
	return nil, nil
}
