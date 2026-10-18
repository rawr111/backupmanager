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

func (drive *GoogleDrive) doHttpRequest(ctx context.Context, method method, path string, body io.Reader) (googleResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "", "", body)
	if err != nil {
		return googleResponse{}, fmt.Errorf("send google drive request: %w", err)
	}

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
