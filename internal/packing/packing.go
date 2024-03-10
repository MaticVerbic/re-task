package packing

// This is a repo responsible for packing.
// In this simple case it is not necessary but for the sake of completeness we make a struct,
// so that we can introduce a consumer interface architecture, which in a real world scenario would
// allow us to mock this package much easier in testing and inject dependencies.

// Packager is a packaging repo.
type Packager struct {
	// dependency injections in here
}

// New returns a new packager.
func New() *Packager {
	return &Packager{}
}

// Calculate takes an order of items, and the sizes of all possible boxes,
// it then calculates the best possible scenario for shipping the amount of items,
// with the least overhead and the least boxes.
func (p *Packager) Calculate(items int, packs []int) []int {
	// calculate the next closest multiple of the smallest pack
	// Could also use math.Ceil(order/packs[0]) * packs[0]
	// but would have to convert to float and then back to int
	closest := (items + packs[0] - 1) / packs[0] * packs[0]
	var out []int

	// greedy algorithm to reduce to packs
	for i := len(packs) - 1; i >= 0; i-- {
		pack := packs[i]
		if no := closest / pack; no > 0 {
			// append the bucket no. of times to output slice
			// could also use map, and then further calculate
			// if we could reduce the number of packs ->
			// map = [500: 2] could become map[1000: 1]
			for i := 0; i < no; i++ {
				out = append(out, pack)
			}
			closest -= no * pack
		}
	}

	// in some cases where packs are not divisible with the smallest pack
	// we could have a miss in the last smallest pack,
	// so we sum the slice of outputs and compare with order,
	// if it's not big enough we add one smallest box.
	sum := sumSlice(out)
	if sum < items && sum > 0 {
		out = append(out, packs[0])
	}

	return out
}

// sumSlice returns the sum of the slice.
func sumSlice(in []int) int {
	sum := 0
	for _, i := range in {
		sum += i
	}

	return sum
}
