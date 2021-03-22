package util

import (
	"testing"
)

var (
	MaxIntP1       = MaxInt
	MaxUintP1 uint = MaxUint
	MinIntM1       = MinInt
)

func TestMax(t *testing.T) {
	tests := []struct{ a, b, e int }{
		{MinInt, MinIntM1, MaxInt},
		{MinIntM1, MinInt, MaxInt},
		{MinIntM1, MinIntM1, MaxInt},

		{MinInt, MinInt, MinInt},
		{MinInt + 1, MinInt, MinInt + 1},
		{MinInt, MinInt + 1, MinInt + 1},

		{-1, -1, -1},
		{-1, 0, 0},
		{-1, 1, 1},

		{0, -1, 0},
		{0, 0, 0},
		{0, 1, 1},

		{1, -1, 1},
		{1, 0, 1},
		{1, 1, 1},

		{MaxInt, MaxInt, MaxInt},
		{MaxInt - 1, MaxInt, MaxInt},
		{MaxInt, MaxInt - 1, MaxInt},

		{MaxIntP1, MaxInt, MaxInt},
		{MaxInt, MaxIntP1, MaxInt},
		{MaxIntP1, MaxIntP1, MinInt},
	}

	for _, test := range tests {
		if g, e := Max(test.a, test.b), test.e; g != e {
			t.Fatal(test.a, test.b, g, e)
		}
	}
}

func TestMin(t *testing.T) {
	tests := []struct{ a, b, e int }{
		{MinIntM1, MinInt, MinInt},
		{MinInt, MinIntM1, MinInt},
		{MinIntM1, MinIntM1, MaxInt},

		{MinInt, MinInt, MinInt},
		{MinInt + 1, MinInt, MinInt},
		{MinInt, MinInt + 1, MinInt},

		{-1, -1, -1},
		{-1, 0, -1},
		{-1, 1, -1},

		{0, -1, -1},
		{0, 0, 0},
		{0, 1, 0},

		{1, -1, -1},
		{1, 0, 0},
		{1, 1, 1},

		{MaxInt, MaxInt, MaxInt},
		{MaxInt - 1, MaxInt, MaxInt - 1},
		{MaxInt, MaxInt - 1, MaxInt - 1},

		{MaxIntP1, MaxInt, MinInt},
		{MaxInt, MaxIntP1, MinInt},
		{MaxIntP1, MaxIntP1, MinInt},
	}

	for _, test := range tests {
		if g, e := Min(test.a, test.b), test.e; g != e {
			t.Fatal(test.a, test.b, g, e)
		}
	}
}
