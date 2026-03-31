package service_test

import (
	"testing"

	"github.com/hayrullahcansu/order-packs-calculator/src/internal/service"
	"github.com/stretchr/testify/assert"
)

type testModel struct {
	packs     []int
	order     int
	correct   map[int]int
	incorrect []map[int]int
}

var defaultTestPacks = []int{250, 500, 1000, 2000, 5000}

func TestSolvePacks(t *testing.T) {
	calculator := service.NewOrderPackCalculator()
	tcs := map[string]testModel{
		"order_1": testModel{
			packs:   defaultTestPacks,
			order:   1,
			correct: map[int]int{250: 1},
			incorrect: []map[int]int{
				map[int]int{
					500: 1, // more items than necessary
				},
			},
		},
		"order_250": testModel{
			packs:   defaultTestPacks,
			order:   250,
			correct: map[int]int{250: 1},
			incorrect: []map[int]int{
				map[int]int{
					500: 1, // more items than necessary
				},
			},
		},
		"order_251": testModel{
			packs:   defaultTestPacks,
			order:   251,
			correct: map[int]int{500: 1},
			incorrect: []map[int]int{
				map[int]int{
					250: 2, // more packs than necessary
				},
			},
		},
		"order_501": testModel{
			packs:   defaultTestPacks,
			order:   501,
			correct: map[int]int{500: 1, 250: 1},
			incorrect: []map[int]int{
				map[int]int{
					1000: 1, // more items than necessary
				},
				map[int]int{
					250: 3, // more packs than necessary
				},
			},
		},
		"order_12001": testModel{
			packs:   defaultTestPacks,
			order:   12001,
			correct: map[int]int{5000: 2, 2000: 1, 250: 1},
			incorrect: []map[int]int{
				map[int]int{
					5000: 3, // more items than necessary
				},
			},
		},
	}

	for testName, testCase := range tcs {
		solvedPacks := calculator.SolvePacks(testCase.packs, testCase.order)
		assert.Equalf(t, testCase.correct, solvedPacks, "%s solved pack is not equal", testName)
		for _, incorrect := range testCase.incorrect {
			assert.NotEqualf(t, testCase.correct, incorrect, "%s solved pack is equal to wrong value", testName)
		}
	}
}

func TestSolvePacksEdgeCase1(t *testing.T) {
	calculator := service.NewOrderPackCalculator()
	tcs := map[string]testModel{
		"order_1": testModel{
			packs:   []int{23, 31, 53},
			order:   500000,
			correct: map[int]int{23: 2, 31: 7, 53: 9429},
		},
	}

	for testName, testCase := range tcs {
		solvedPacks := calculator.SolvePacks(testCase.packs, testCase.order)
		assert.Equalf(t, testCase.correct, solvedPacks, "%s solved pack is not equal", testName)
		for _, incorrect := range testCase.incorrect {
			assert.NotEqualf(t, testCase.correct, incorrect, "%s solved pack is equal to wrong value", testName)
		}
	}
}
