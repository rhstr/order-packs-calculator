package pack

import (
	"slices"
)

func CalculatePacking(order int, packs ...int) map[int]int {
	if order <= 0 || len(packs) == 0 {
		return nil
	}
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

type state struct {
	packs int
	items int
	prev  int
	pack  int
}
