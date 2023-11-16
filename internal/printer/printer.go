package printer

import (
	"fmt"
	"io"
	"time"
)

type Printer interface {
	OutputEntry(url string, contentLength int64, statusCode int, status string, ts time.Duration) (int, error)
	OutputValidationError(url string) (int, error)
}

var _ Printer = (*printer)(nil)

func New(w io.Writer) Printer {
	return &printer{w: w}
}

type printer struct {
	w io.Writer
}

func (s printer) OutputEntry(url string, contentLength int64, statusCode int, status string, ts time.Duration) (
	int, error,
) {
	n, err := fmt.Fprintf(
		s.w, "GET: %s	STATUS: %s	CODE: %d	SIZE: %d	DURATION: %s\n", url, status, statusCode,
		contentLength,
		ts.String(),
	)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (s printer) OutputValidationError(url string) (int, error) {
	n, err := fmt.Fprintf(s.w, "INVALID URL: %s\n", url)
	if err != nil {
		return 0, err
	}
	return n, nil
}
