package googledrive

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

type GoogleDrive struct {
	cli    http.Client
	apikey string
}

// use pointer receiver because we dont want copy http.client
func New(apikey string) *GoogleDrive {
	return &GoogleDrive{
		apikey: apikey,
		cli:    http.Client{},
	}
}

const uploadFilePath string = "/upload/drive/v3/files?uploadType=media"

func (drive *GoogleDrive) UploadFile(ctx context.Context, filePath string, alias string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file %s failed: %w", filePath, err)
	}
	defer file.Close()

	resp, err := drive.doHttpRequest(ctx, methodPost, uploadFilePath, file)
	if err != nil {
		return err
	}

	if resp.statusCode != 200 {
		return fmt.Errorf("upload google drive file: %s", resp.body)
	}

	return nil
}

func (drive *GoogleDrive) DownloadFile(fileName string) (error, []byte) {
	return nil, nil
}
