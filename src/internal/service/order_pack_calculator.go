package service

import (
	"math"
	"slices"
)

type OrderPackCalculator interface {
	// SolvePacks calculates minimize the number of packs/quantity
	// packs represent available packs
	// order is given order amount
	// return type is map of packSize with count for each ordered packs
	SolvePacks(packs []int, order int) map[int]int
}

type orderPackCalculator struct {
}

// SolvePacks implements OrderPackCalculator.
func (o *orderPackCalculator) SolvePacks(packs []int, order int) map[int]int {

	MAX_INT := math.MaxInt

	// we need an upper bound because the order might not be achievable exactly.
	// example: packs=[23,31], order=1 → smallest achievable is 23.
	// worst case ceiling is order + largest pack size.
	UPPER_BOUND := order + slices.Max(packs)

	// dp[i] = minimum number of packs needed to reach exactly i items.
	// we size it to UPPER_BOUND to cover cases where order cannot be met exactly.
	dp := make([]int, UPPER_BOUND+1)

	// choice[i] = which pack size was used to arrive at position i.
	// we use this later to backtrack and reconstruct the solution
	choice := make([]int, UPPER_BOUND+1)

	dp[0] = 0

	// we fill slots with MAX number of Integer
	// it will help to get minimum value while we're finding minimize the number of packs/quantity
	for i := 1; i <= UPPER_BOUND; i++ {
		dp[i] = MAX_INT
	}

	// process possibilities
	// think of it as: for each pack size, we try to extend every reachable position
	// by adding one more pack of the size.
	for _, pack := range packs {
		for left := pack; left <= UPPER_BOUND; left++ {

			// only proceed if (left-pack) is reachable
			// also check if this path uses less packs than what we found before.
			// because we're trying to minimize packs/quantities
			if dp[left-pack] != MAX_INT && dp[left-pack]+1 < dp[left] {
				dp[left] = dp[left-pack] + 1

				// remember which pack we used to get here
				// so we can trace back the full solution later
				choice[left] = pack
			}
		}
	}

	// find best achievable amount starting from the order
	// this ensures we never send more items than necessary
	best := -1
	for left := order; left <= UPPER_BOUND; left++ {
		if dp[left] != MAX_INT {
			best = left
			break
		}
	}

	if best == -1 {
		// if no solution found, return nil
		// this shouldn't possible if pack sizes are valid
		return nil
	}

	// backtrack from best to 0 using the choice
	// each step tells us which pack did i use to arrive here? -> subtract it -> repeat
	// example: best=500, choice[500]=250 -> pos=250, choice[250]=250 -> post=0 -> Done
	result := make(map[int]int)
	left := best
	for left > 0 {
		pack := choice[left]
		result[pack]++
		left -= pack
	}

	return result
}

func NewOrderPackCalculator() OrderPackCalculator {
	calculator := &orderPackCalculator{}

	return calculator
}
