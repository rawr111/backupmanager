package googledrive

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type googleResponse struct {
	statusCode int
	body       string
}

type method string

const methodGet method = "GET"
const methodPost method = "POST"
const host string = "https://www.googleapis.com"

func (drive *GoogleDrive) doHttpRequest(ctx context.Context, method method, path string, body io.Reader) (googleResponse, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		string(method),
		fmt.Sprintf("%s/%s", host, path),
		body,
	)
	if err != nil {
		return googleResponse{}, fmt.Errorf("send google drive request: %w", err)
	}

	//need to get it from previous step
	req.Header.Set("Content-type", "")
	req.Header.Set("X-goog-api-key", drive.apikey)

	resp, err := drive.cli.Do(req)
	if err != nil {
		return googleResponse{}, fmt.Errorf("send google drive request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return googleResponse{}, err
	}

	return googleResponse{
		statusCode: resp.StatusCode,
		body:       string(respBody),
	}, nil
}
