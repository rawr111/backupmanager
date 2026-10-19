package googledrive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type payload struct {
	err struct {
		code    int
		message string
	}
}

type googleResponse struct {
	statusCode int
	body       payload
}

type method string

const methodGet method = "GET"
const methodPost method = "POST"
const host string = "https://www.googleapis.com"

type headers []struct {
	key   string
	value string
}

func (h headers) Set(request *http.Request) {
	for _, header := range h {
		request.Header.Set(header.key, header.value)
	}
}

func (drive *GoogleDrive) doHttpRequest(ctx context.Context, method method, path string, body io.Reader, headers headers) (googleResponse, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		string(method),
		fmt.Sprintf("%s/%s", host, path),
		body,
	)
	if err != nil {
		return googleResponse{}, fmt.Errorf("send google drive request: %w", err)
	}

	headers.Set(req)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", drive.accessToken))

	resp, err := drive.cli.Do(req)
	if err != nil {
		return googleResponse{}, fmt.Errorf("send google drive request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return googleResponse{}, err
	}

	var payload payload
	err = json.Unmarshal(respBody, &payload)
	if err != nil {
		//если гугл нам ответил не стандартной ошибкой вдруг
		payload.err = struct {
			code    int
			message string
		}{
			code:    resp.StatusCode,
			message: string(respBody),
		}
		return googleResponse{
			statusCode: resp.StatusCode,
			body:       payload,
		}, nil
	}

	return googleResponse{
		statusCode: resp.StatusCode,
		body:       payload,
	}, nil
}
