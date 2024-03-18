package main

import (
	"fmt"
	"retask/internal/packing"
)

// If tests are not enough, use this package for debugging and running the algorithm alone without the API.
func main() {
	packs := []int{3, 5, 8, 10}
	target := 51

	repo := packing.New()
	o := repo.Calculate(packs, target)
	fmt.Println(o)
}
