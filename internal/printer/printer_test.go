package printer

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type mockWriter struct {
	Written []byte
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	m.Written = append(m.Written, p...)
	return len(p), nil
}

type mockErrorWriter struct{}

func (m *mockErrorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("mocked error")
}

func TestPrinterOutputEntry(t *testing.T) {
	t.Parallel()

	mw := &mockWriter{}
	p := New(mw)
	url := "https://example.com"
	contentLength := int64(1024)
	statusCode := 200
	status := "OK"
	ts := time.Second

	n, err := p.OutputEntry(url, contentLength, statusCode, status, ts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := fmt.Sprintf(
		"GET: %s\tSTATUS: %s\tCODE: %d\tSIZE: %d\tDURATION: %s\n",
		url, status, statusCode, contentLength, ts.String(),
	)

	if string(mw.Written) != expected {
		t.Errorf("expected:\n%s\nActual output:\n%s\n", expected, mw.Written)
	}

	if n != len(expected) {
		t.Errorf("expected: %d, got %d", len(expected), n)
	}
}

func TestPrinterOutputValidationError(t *testing.T) {
	t.Parallel()

	mw := &mockWriter{}
	p := New(mw)
	url := "https://example.com"
	n, err := p.OutputValidationError(url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := fmt.Sprintf("INVALID URL: %s\n", url)

	if string(mw.Written) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s\n", expected, mw.Written)
	}

	if n != len(expected) {
		t.Errorf("expected: %d, got: %d", len(expected), n)
	}
}

func TestPrinterError(t *testing.T) {
	t.Parallel()

	p := printer{w: &mockErrorWriter{}}
	url := "https://example.com"
	contentLength := int64(1024)
	statusCode := 200
	status := "OK"
	ts := time.Second

	n, err := p.OutputEntry(url, contentLength, statusCode, status, ts)
	if err == nil {
		t.Fatal("expected error, but got nil")
	}

	if n != 0 {
		t.Errorf("expected: 0, got %d", n)
	}
}
