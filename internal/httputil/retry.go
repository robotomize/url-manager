package httputil

import (
	"context"
	"net/http"
	"time"
)

type BackoffFunc func(min, max time.Duration, num int) time.Duration

func NewRetryClient(
	client Client,
	count int,
	minWait time.Duration,
	maxWait time.Duration,
	backoffFunc BackoffFunc,
) Client {
	return &RetryClient{Client: client, count: count, minWait: minWait, maxWait: maxWait, backoffFunc: backoffFunc}
}

type RetryClient struct {
	Client
	count       int
	minWait     time.Duration
	maxWait     time.Duration
	backoffFunc BackoffFunc
}

func (r *RetryClient) Do(req *http.Request) (resp *http.Response, err error) {
	resp, err = r.Client.Do(req)
	if err == nil || r.count < 1 {
		return
	}

	for i := 0; i < r.count && err != nil; i++ {
		waitTime := r.backoffFunc(r.minWait, r.maxWait, i)
		wait(req.Context(), waitTime)
		resp, err = r.Client.Do(req)
	}

	return
}

func wait(ctx context.Context, t time.Duration) {
	ticker := time.NewTicker(t)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return
	case <-ticker.C:
	}
}
