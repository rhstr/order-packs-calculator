package pack_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rhstr/order-packs-calculator/internal/pack"
)

func TestCalculatePacking(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		orderSize int
		packs     []int
		want      []pack.Packing
		wantErr   bool
	}{
		{
			name:      "Edge case handled correctly",
			orderSize: 500000,
			packs:     []int{23, 31, 53},
			want: []pack.Packing{
				{BoxSize: 53, Quantity: 9429},
				{BoxSize: 31, Quantity: 7},
				{BoxSize: 23, Quantity: 2},
			},
		},
		{
			name:      "error: order size is zero",
			orderSize: 0,
			packs:     []int{23, 31, 53},
			wantErr:   true,
		},
		{
			name:      "error: order size is negative",
			orderSize: -1,
			packs:     []int{23, 31, 53},
			wantErr:   true,
		},
		{
			name:      "error: no packs provided",
			orderSize: 10,
			packs:     nil,
			wantErr:   true,
		},
		{
			name:      "error: one of the packs is zero",
			orderSize: 10,
			packs:     []int{23, 0, 53},
			wantErr:   true,
		},
		{
			name:      "error: packs are duplicated",
			orderSize: 10,
			packs:     []int{2, 23, 31, 53, 31},
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			calc := pack.NewCalculator(1000000)

			got, err := calc.CalculatePacking(tc.orderSize, tc.packs...)
			if tc.wantErr {
				assert.Error(t, err)

				return
			}

			assert.Equal(t, tc.want, got)
			assert.NoError(t, err)
		})
	}
}
