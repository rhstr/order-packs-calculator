package pack_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rhstr/order-packs-calculator/internal/pack"
)

func TestCalculatePacking(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		order int
		packs []int
		want  map[int]int
	}{
		{
			name:  "Edge case",
			order: 500000,
			packs: []int{23, 31, 53},
			want:  map[int]int{23: 2, 31: 7, 53: 9429},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := pack.CalculatePacking(tc.order, tc.packs...)
			assert.Equal(t, tc.want, got)
		})
	}
}
