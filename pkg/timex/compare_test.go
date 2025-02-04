package timex_test

import (
	"testing"
	"time"

	"github.com/openmeterio/openmeter/pkg/timex"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		time1    time.Time
		time2    time.Time
		expected int
	}{
		{
			name:     "equal times",
			time1:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			time2:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "time1 after time2",
			time1:    time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			time2:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 86400000000000,
		},
		{
			name:     "time1 before time2",
			time1:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			time2:    time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			expected: -86400000000000,
		},
		{
			name:     "time1 milliseconds after time2",
			time1:    time.Date(2025, 1, 1, 0, 0, 0, 1000000, time.UTC),
			time2:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timex.Compare(tt.time1, tt.time2)
			if result != tt.expected {
				t.Errorf("Compare() = %v, want %v", result, tt.expected)
			}
		})
	}
}
