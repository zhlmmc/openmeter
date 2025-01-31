package slicesx_test

import (
	"testing"

	"github.com/openmeterio/openmeter/pkg/slicesx"
	"github.com/stretchr/testify/assert"
)

func TestLast(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		pred     func(int) bool
		wantVal  *int
		wantIdx  int
		wantBool bool
	}{
		{
			name:     "nil slice",
			slice:    nil,
			pred:     func(i int) bool { return true },
			wantVal:  nil,
			wantIdx:  -1,
			wantBool: false,
		},
		{
			name:     "empty slice",
			slice:    []int{},
			pred:     func(i int) bool { return true },
			wantVal:  nil,
			wantIdx:  -1,
			wantBool: false,
		},
		{
			name:     "no match",
			slice:    []int{1, 2, 3},
			pred:     func(i int) bool { return i > 10 },
			wantVal:  nil,
			wantIdx:  -1,
			wantBool: false,
		},
		{
			name:     "match last element",
			slice:    []int{1, 2, 3},
			pred:     func(i int) bool { return i == 3 },
			wantVal:  intPtr(3),
			wantIdx:  2,
			wantBool: true,
		},
		{
			name:     "match first element",
			slice:    []int{1, 2, 3},
			pred:     func(i int) bool { return i == 1 },
			wantVal:  intPtr(1),
			wantIdx:  0,
			wantBool: true,
		},
		{
			name:     "match middle element",
			slice:    []int{1, 2, 3},
			pred:     func(i int) bool { return i == 2 },
			wantVal:  intPtr(2),
			wantIdx:  1,
			wantBool: true,
		},
		{
			name:     "multiple matches returns last",
			slice:    []int{1, 2, 2, 3},
			pred:     func(i int) bool { return i == 2 },
			wantVal:  intPtr(2),
			wantIdx:  2,
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotIdx, gotBool := slicesx.Last(tt.slice, tt.pred)

			if tt.wantVal == nil {
				assert.Nil(t, gotVal)
			} else {
				assert.Equal(t, *tt.wantVal, *gotVal)
			}
			assert.Equal(t, tt.wantIdx, gotIdx)
			assert.Equal(t, tt.wantBool, gotBool)
		})
	}
}

func intPtr(i int) *int {
	return &i
}
