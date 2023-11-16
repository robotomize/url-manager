package urlreader

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestReader_ReadURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input       string
		expected    string
		expectedErr error
	}{
		{
			name:        "test_valid",
			input:       "https://example.com\n",
			expected:    "https://example.com",
			expectedErr: nil,
		},
		{
			name:        "test_empty_input",
			input:       "",
			expected:    "",
			expectedErr: io.EOF,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(
			tc.name, func(t *testing.T) {
				t.Parallel()

				r := New(strings.NewReader(tc.input))

				line, err := r.ReadURL()
				if err != nil {
					if tc.expectedErr == nil || err.Error() != tc.expectedErr.Error() {
						t.Errorf("expected: %v, got: %v", tc.expectedErr, err)
					}
				}

				if line != tc.expected {
					t.Errorf("expected: %s, got: %s", tc.expected, line)
				}

				{
					if _, err := r.ReadURL(); err != nil {
						if !errors.Is(err, io.EOF) {
							t.Errorf("unexpected error")
						}
					}
				}
			},
		)
	}
}
