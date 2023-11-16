package urlchecker

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

type mockHTTPClient struct {
	StatusCode    int
	Status        string
	ContentLength int64
	Err           error
}

func (m *mockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode:    m.StatusCode,
		Status:        m.Status,
		ContentLength: m.ContentLength,
		Body:          nil,
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
			name: "test_ok_reponse",
			clientMock: &mockHTTPClient{
				StatusCode:    200,
				Status:        "OK",
				ContentLength: 100,
			},
			url: "https://example.com",
			expected: Response{
				StatusCode:    200,
				Status:        "OK",
				ContentLength: 100,
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
		t.Run(
			tc.name, func(t *testing.T) {
				t.Parallel()

				ch := New(tc.clientMock)

				result, err := ch.Check(context.Background(), tc.url)

				if err != nil {
					if tc.expectedErr == nil || err.Error() != tc.expectedErr.Error() {
						t.Errorf("expected: %v, got: %v", tc.expectedErr, err)
					}
				} else if tc.expectedErr != nil {
					t.Errorf("expected: %v, got: nil", tc.expectedErr)
				}

				if result != tc.expected {
					t.Errorf("expected: %v, got: %v", tc.expected, result)
				}
			},
		)
	}
}
