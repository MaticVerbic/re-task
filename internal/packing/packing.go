package packing

import (
	"maps"
	"math"
	"sort"
)

// This is a repo responsible for packing.
// In this simple case it is not necessary but for the sake of completeness we make a struct,
// so that we can introduce a consumer interface architecture, which in a real world scenario would
// allow us to mock this package much easier in testing and inject dependencies.

// Packager is a packaging repo.
type Packager struct {
	// dependency injections in here
}

func New() *Packager {
	return &Packager{}
}

// sumSlice returns the sum of the slice.
func sumSlice(in []int) int {
	sum := 0
	for _, i := range in {
		sum += i
	}

	return sum
}

// Calculate calculates the best possible pack distribution for a given set of packages and an order. It uses a recursive
// approach to generate all combinations that sum to equal or above the order, then calculates the sum of each, and finds the
// one with the smallest positive difference to order and smallest box count.
func (p *Packager) Calculate(items int, packs []int) []int {
	if items <= 0 {
		return []int{}
	}

	if items <= packs[len(packs)-1] {
		return []int{packs[len(packs)-1]}
	}
	var c []map[int]int

	combinations(items, packs, map[int]int{}, &c)

	var outs [][]int
	for j, item := range c {
		outs = append(outs, []int{})
		for k, v := range item {
			for _ = v; v > 0; v-- {
				outs[j] = append(outs[j], k)
			}
		}
	}

	optimalSolution := getLeastShortest(outs, items)
	sort.Sort(sort.Reverse(sort.IntSlice(optimalSolution)))

	return optimalSolution
}

// combinations calculates the maximum for each given pack until the order is reduced to no more options
// if the order is above the maximum then it recursively follows the same pattern until no more matches are found.
func combinations(items int, packs []int, previousOrder map[int]int, out *[]map[int]int) {
	for i, pack := range packs {
		maximum := int(math.Ceil(float64(items) / float64(pack)))
		order := map[int]int{}
		for k, v := range previousOrder {
			order[k] = v
		}
		order[pack] += maximum
		*out = append(*out, maps.Clone(order))

		for j := maximum - 1; j > 0; j-- {
			order[pack] = j
			combinations(items-pack*j, packs[i+1:], order, out)
		}
	}
}

// getLeastShortest creates a map of sums of all arrays, then finds the closest match to the target, and returns
// the shortest array for the given match.
func getLeastShortest(arr [][]int, target int) []int {
	sums := map[int][][]int{}
	for _, item := range arr {
		sums[sumSlice(item)] = append(sums[sumSlice(item)], item)
	}

	closestHigher := math.MaxInt32
	idx := math.MaxInt32
	for k := range sums {
		diff := k - target
		if diff >= 0 && diff < closestHigher {
			closestHigher = diff
			idx = k
		}
	}

	shortest := sums[idx][0]
	smallestLen := len(sums[idx][0])
	for i := 0; i < len(sums[idx]); i++ {
		if len(sums[idx][i]) <= smallestLen {
			shortest = sums[idx][i]
			smallestLen = len(sums[idx][i])
		}
	}

	return shortest
}
