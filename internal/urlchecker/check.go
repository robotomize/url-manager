package urlchecker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Response struct {
	StatusCode    int
	Status        string
	ContentLength int64
	TS            time.Duration
}

type Checker interface {
	Check(ctx context.Context, u string) (Response, error)
}

func New(httpClient client) Checker {
	return &checker{httpClient: httpClient}
}

type checker struct {
	httpClient client
}

func (m *checker) Check(ctx context.Context, u string) (Response, error) {
	var result Response
	ts := time.Now().UTC()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return result, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("http client Do: %w", err)
	}

	defer resp.Body.Close()

	var n int64
	buf := make([]byte, 4096)

	for {
		read, err := resp.Body.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return result, fmt.Errorf("http client response Body Read: %w", err)
		}
		n += int64(read)
	}

	result.ContentLength = n
	result.StatusCode = resp.StatusCode
	result.Status = resp.Status
	result.TS = time.Since(ts)

	return result, nil
}
