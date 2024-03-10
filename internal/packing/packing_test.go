package packing

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var packs = []int{250, 500, 1000, 2000, 5000}
var packsNonDivisible = []int{11, 34, 59, 70}

var largeExpected []int

func TestMain(m *testing.M) {
	for i := 0; i < 27; i++ {
		largeExpected = append(largeExpected, 5000)
	}

	largeExpected = append(largeExpected, []int{2000, 1000, 500, 250}...)

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
		{packs: packs, order: 12001, expected: []int{5000, 5000, 2000, 250}},
		{packs: packs, order: 9499, expected: []int{5000, 2000, 2000, 500}},
		{packs: packs, order: 138501, expected: largeExpected},
		{packs: packsNonDivisible, order: 131, expected: []int{70, 59, 11}},
		{packs: packsNonDivisible, order: 178, expected: []int{70, 70, 34, 11}},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			repo := &Packager{}
			out := repo.Calculate(test.order, test.packs)
			t.Log(test.order, out)
			assert.Equal(t, test.expected, out)
		})
	}
}

func BenchmarkCalculate(b *testing.B) {
	// we have unit tested this order to make sure the response is valid,
	// and we're not benching invalid data.
	testOrder := 138501
	repo := &Packager{}
	for n := 0; n < b.N; n++ {
		repo.Calculate(testOrder, packs)
	}
}
