package util

const (
	MaxInt  = 1<<(IntBits-1) - 1
	MinInt  = -MaxInt - 1
	MaxUint = 1<<IntBits - 1
	IntBits = 1 << (^uint(0)>>32&1 + ^uint(0)>>16&1 + ^uint(0)>>8&1 + 3)
)

func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// Min returns the smaller of a and b.
func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
