package packing

import (
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var packs = []int{250, 500, 1000, 2000, 5000}
var packsNonDivisible = []int{11, 34, 59, 70}
var packsSmall = []int{3, 4, 5}

var largeExpected []int

func TestMain(m *testing.M) {
	for i := 0; i < 7; i++ {
		largeExpected = append(largeExpected, 5000)
	}

	largeExpected = append(largeExpected, []int{2000, 2000, 500, 250}...)

	os.Exit(m.Run())
}

func TestAlgorithm(t *testing.T) {
	tests := []struct {
		packs    []int
		order    int
		expected []int
	}{
		{packs: packs, order: 1, expected: []int{250}},
		{packs: packs, order: 250, expected: []int{250}},
		{packs: packs, order: 251, expected: []int{500}},
		{packs: packs, order: 501, expected: []int{500, 250}},
		{packs: packs, order: 749, expected: []int{500, 250}},
		{packs: packs, order: 751, expected: []int{1000}},
		{packs: packs, order: 12001, expected: []int{5000, 5000, 2000, 250}},
		{packs: packs, order: 9499, expected: []int{5000, 2000, 2000, 500}},
		{packs: packs, order: 39501, expected: largeExpected},
		{packs: packsNonDivisible, order: 131, expected: []int{11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11}},
		{packs: packsNonDivisible, order: 178, expected: []int{34, 34, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11}},
		{packs: packsSmall, order: 10, expected: []int{5, 5}},
		{packs: packs, order: 0, expected: []int{}},
		{packs: packs, order: -1, expected: []int{}},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			sort.Sort(sort.Reverse(sort.IntSlice(test.packs[:])))
			repo := &Packager{}
			out := repo.Calculate(test.order, test.packs)
			assert.Equal(t, test.expected, out)
		})
	}
}

func TestGetLeastShortest(t *testing.T) {
	tests := []struct {
		arr      [][]int
		expected []int
		target   int
	}{
		{arr: [][]int{{5000}, {2000}, {1000}, {500}, {250}}, expected: []int{250}, target: 250},
		{arr: [][]int{{5000}, {2000}, {1000}, {500}, {250, 250}}, expected: []int{500}, target: 251},
		{arr: [][]int{{5000}, {2000}, {1000}, {500, 500}, {250, 500}, {250, 250, 250}}, expected: []int{250, 500}, target: 501},
		{arr: [][]int{{5000}, {2000}, {1000}, {500, 500}, {250, 250, 500}, {250, 250, 250, 250}}, expected: []int{1000}, target: 751},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual := getLeastShortest(test.arr, test.target)

			assert.Equal(t, test.expected, actual)
		})
	}
}

func BenchmarkCalculate(b *testing.B) {
	// we have unit tested this order to make sure the response is valid,
	// and we're not benching invalid data.
	testOrder := 39501
	repo := &Packager{}
	for n := 0; n < b.N; n++ {
		repo.Calculate(testOrder, packs)
	}
}
