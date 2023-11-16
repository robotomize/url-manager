package httputil

import (
	"math"
	"testing"
	"time"
)

func TestExponentialBackoff(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		min      time.Duration
		max      time.Duration
		num      int
		expected time.Duration
	}{
		{
			name:     "test_min_value_with_zero_attempt",
			min:      time.Millisecond,
			max:      time.Second,
			num:      0,
			expected: time.Millisecond,
		},
		{
			name:     "test_with_exp",
			min:      time.Millisecond,
			max:      time.Second,
			num:      3,
			expected: time.Duration(math.Pow(2, 3)) * time.Millisecond,
		},
		{
			name:     "test_check_max",
			min:      time.Millisecond,
			max:      time.Second,
			num:      10,
			expected: time.Second,
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(
			tc.name, func(t *testing.T) {
				t.Parallel()

				result := ExponentialBackoff(tc.min, tc.max, tc.num)

				if result != tc.expected {
					t.Errorf("Unexpected result. Expected: %v, Got: %v", tc.expected, result)
				}
			},
		)
	}
}
