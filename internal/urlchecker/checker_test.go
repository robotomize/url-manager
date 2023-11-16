package urlchecker

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func newMockBodyCloser(r io.Reader) *mockBodyCloser {
	return &mockBodyCloser{Reader: r}
}

type mockBodyCloser struct {
	io.Reader
}

func (m *mockBodyCloser) Close() error {
	return nil
}

type mockHTTPClient struct {
	StatusCode    int
	Status        string
	ContentLength int64
	Content       string
	Err           error
}

func (m *mockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: m.StatusCode,
		Status:     m.Status,
		Body:       newMockBodyCloser(strings.NewReader(m.Content)),
	}

	return resp, m.Err
}

func TestChecker_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		clientMock  *mockHTTPClient
		url         string
		expected    Response
		expectedErr error
	}{
		{
			name: "test_ok_response",
			clientMock: &mockHTTPClient{
				StatusCode:    200,
				Status:        "OK",
				ContentLength: 11,
				Content:       "hello world",
			},
			url: "https://example.com",
			expected: Response{
				StatusCode:    200,
				Status:        "OK",
				ContentLength: 11,
			},
			expectedErr: nil,
		},
		{
			name: "test_http_error",
			clientMock: &mockHTTPClient{
				Err: errors.New("mock error"),
			},
			url:         "https://example.com",
			expected:    Response{},
			expectedErr: errors.New("http client Do: mock error"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(
			tc.name, func(t *testing.T) {
				t.Parallel()

				ch := New(tc.clientMock)

				result, err := ch.Check(context.Background(), tc.url)
				if err != nil {
					if tc.expectedErr == nil || err.Error() != tc.expectedErr.Error() {
						t.Errorf("expected: %v, got: %v", tc.expectedErr, err)
					}
				}

				if result.ContentLength != tc.expected.ContentLength {
					t.Errorf("expected: %v, got: %v", tc.expected.ContentLength, result.ContentLength)
				}

				if result.Status != tc.expected.Status {
					t.Errorf("expected: %v, got: %v", tc.expected.Status, result.Status)
				}

				if result.StatusCode != tc.expected.StatusCode {
					t.Errorf("expected: %v, got: %v", tc.expected.StatusCode, result.StatusCode)
				}
			},
		)
	}
}
