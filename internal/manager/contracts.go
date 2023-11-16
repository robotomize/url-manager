package manager

import (
	"context"
	"time"

	"github.com/robotomize/url-manager/internal/urlchecker"
)

type reader interface {
	ReadURL() (string, error)
}

type logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
}

type urlChecker interface {
	Check(ctx context.Context, u string) (urlchecker.Response, error)
}

type printer interface {
	OutputEntry(url string, contentLength int64, statusCode int, status string, ts time.Duration) (int, error)
	OutputValidationError(url string) (int, error)
}
