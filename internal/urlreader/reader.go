package urlreader

import (
	"bufio"
	"fmt"
	"io"
)

type Reader interface {
	ReadURL() (string, error)
}

var _ Reader = (*reader)(nil)

func New(r io.Reader) Reader {
	return &reader{r: bufio.NewScanner(r)}
}

type reader struct {
	r *bufio.Scanner
}

func (f *reader) ReadURL() (string, error) {
	if f.r.Scan() {
		return f.r.Text(), nil
	}

	if err := f.r.Err(); err != nil {
		return "", fmt.Errorf("bufio.Scanner Err: %w", err)
	}

	return "", io.EOF
}
