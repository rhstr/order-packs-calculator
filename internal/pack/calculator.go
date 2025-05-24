package pack

import (
	"errors"
	"slices"
	"sort"
)

// Packing represents the packing of items into boxes of a specific size.
type Packing struct {
	BoxSize  int `json:"boxSize"`
	Quantity int `json:"quantity"`
}

// CalculatePacking calculates the optimal packing of items into boxes of given sizes.
func CalculatePacking(order int, packs ...int) ([]Packing, error) {
	if order <= 0 {
		return nil, errors.New("order size must be greater than zero")
	}

	if len(packs) == 0 {
		return nil, errors.New("at least one pack size must be provided")
	}

	for _, pack := range packs {
		if pack <= 0 {
			return nil, errors.New("pack sizes must be greater than zero")
		}
	}

	if arePacksDuplicated(packs) {
		return nil, errors.New("pack sizes must be unique")
	}

	calculations := calculatePackingDP(order, packs...)
	if calculations == nil {
		return nil, errors.New("no packing solution found")
	}

	result := make([]Packing, 0, len(calculations))
	for packSize, quantity := range calculations {
		result = append(result, Packing{
			BoxSize:  packSize,
			Quantity: quantity,
		})
	}

	slices.SortFunc(result, func(a, b Packing) int {
		return b.BoxSize - a.BoxSize
	})

	return result, nil
}

// arePacksDuplicated checks if there are any duplicate pack sizes.
func arePacksDuplicated(packs []int) bool {
	if len(packs) < 2 {
		return false
	}

	sort.Ints(packs)
	for i := 1; i < len(packs); i++ {
		if packs[i] == packs[i-1] {
			return true
		}
	}

	return false
}

// calculatePackingDP uses dynamic programming to find the optimal packing solution.
func calculatePackingDP(order int, packs ...int) map[int]int {
	slices.SortFunc(packs, func(a, b int) int { return b - a })

	limit := order + packs[0] // allow overpacking up to one largest pack

	dp := make([]*state, limit+1)
	dp[0] = &state{packs: 0, items: 0, prev: -1, pack: 0}

	for i := 0; i <= limit; i++ {
		if dp[i] == nil {
			continue
		}

		for _, p := range packs {
			next := i + p
			if next > limit {
				continue
			}

			newPacks := dp[i].packs + 1
			newItems := dp[i].items + p

			if dp[next] == nil ||
				newItems < dp[next].items ||
				(newItems == dp[next].items && newPacks < dp[next].packs) {
				dp[next] = &state{packs: newPacks, items: newItems, prev: i, pack: p}
			}
		}
	}

	// Find the best solution with minimal items, then minimal packs
	best := -1
	for i := order; i <= limit; i++ {
		if dp[i] == nil {
			continue
		}

		if best == -1 ||
			dp[i].items < dp[best].items ||
			(dp[i].items == dp[best].items && dp[i].packs < dp[best].packs) {
			best = i
		}
	}

	if best == -1 {
		return nil
	}

	// Reconstruct solution
	result := make(map[int]int)
	for s := best; s > 0; {
		p := dp[s].pack
		result[p]++
		s = dp[s].prev
	}

	return result
}

// state represents the intermediate state of the packing calculation based on dynamic programming.
type state struct {
	packs int
	items int
	prev  int
	pack  int
}
