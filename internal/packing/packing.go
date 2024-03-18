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

func (p *Packager) Calculate(packs []int, target int) []int {
	// in API, we don't need this, since we ensure this only happens once when new package sizes are created, but
	// for completeness’s sake, if the tester runs this algorithm on its own with the API we do it here too.
	// This would not have needed to be done in real world scenarios.
	sort.Ints(packs)

	// eliminate negative or zero. We do this in the API already, but as above for completeness's sake we do it here to.
	if target <= 0 {
		return []int{}
	}

	// if the order is smaller than the smallest package, then just return that.
	if target < packs[0] {
		return []int{packs[0]}
	}

	// initialize the difference checking vars
	//  we will use these to first calculate the smallest difference in items (2. point)
	//  then the smallest length (3. point) in case some of the results come up to the same sum
	smallestDifference := math.MaxInt32
	smallestLength := math.MaxInt32

	// make an array to count boxes
	boxes := make([]int, target+packs[len(packs)-1]+1)
	// here we will create combinations
	solutions := make([][]int, len(boxes))
	for i := 0; i < len(boxes); i++ {
		// make sure the size is 32bit, otherwise we will cause overflow and get negative value
		boxes[i] = math.MaxInt32
	}

	// first box is always the 0th case, making the least viable solution to any order to send 0 boxes.
	boxes[0] = 0
	// here we will store all solutions, we could omit this, but then we would have to rebuild each solution dynamically
	// aka. create a cartesian product for each row.
	solutions[0] = make([]int, len(packs))

	// we will store our output here
	var closestMatch []int
	for i := 1; i < len(boxes); i++ {
		for j := 0; j < len(packs); j++ {
			// dynamically look packs[j] behind if the box has already been accounted for.
			// here we could also look ahead as per tabulation standard, but since we're building a list of combinations,
			// looking behind is better since we can also take the existing array and just add instead of having to rebuild it.
			if packs[j] <= i && boxes[i-packs[j]]+1 < boxes[i] {
				// increment the number of boxes
				boxes[i] = boxes[i-packs[j]] + 1

				// reuse the existing array - memoization
				if solutions[i] == nil {
					solutions[i] = make([]int, len(packs))
				}
				for k := range packs {
					solutions[i][k] = solutions[i-packs[j]][k]
				}

				// add one current box
				solutions[i][j]++

				// sum the array in order to check whether the sum is smaller than the previous smaller sum.
				s := sum(solutions[i], packs)
				if s >= target && s <= smallestDifference {
					// if diff is definitely smaller, we don't check for number of boxes
					if s < smallestDifference {
						smallestDifference = s
						smallestLength = boxes[i]
						closestMatch = remap(solutions[i], packs)
						continue
					}
					// if diff is not definitely smaller (e.g. is equal), we check for number of boxes
					if boxes[i] <= smallestLength {
						smallestDifference = s
						smallestLength = boxes[i]
						closestMatch = remap(solutions[i], packs)
					}
				}
			}
		}
	}

	// we return the closest match
	return closestMatch
}

// remap function changes the working array (len(packs), with number of repetitions for each package), to array with
// each package repeated as many times as required. We sort it in reverse order for clarity.
// example: packs = [5, 10, 20, 30], arr = [1, 2, 3, 4], output = [30, 30, 30, 30, 20, 20, 20, 10, 10, 5]
func remap(arr, packs []int) []int {
	out := []int{}

	for i := 0; i < len(arr); i++ {
		for j := arr[i]; j > 0; j-- {
			out = append(out, packs[i])
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(out[:])))

	return out
}

// sum the array of repetitions, so we can calculate the smallest positive difference to the target.
// example: packs = [2, 3], arr = [1, 4], algo = 1*2 + 3*4, output = 14
func sum(arr []int, packs []int) int {
	out := 0
	for i, item := range arr {
		out += item * packs[i]
	}

	return out
}

// The code below this line is not used, it will simply be used by my for the review after the task is completed. Please
// feel free to ignore it.

// Calculate calculates the best possible pack distribution for a given set of packages and an order. It uses a recursive
// approach to generate all combinations that sum to equal or above the order, then calculates the sum of each, and finds the
// one with the smallest positive difference to order and smallest box count.
func (p *Packager) calculateOld(items int, packs []int) []int {
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
