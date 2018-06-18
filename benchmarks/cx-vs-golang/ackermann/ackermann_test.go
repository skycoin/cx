package ackermann

import "testing"

func BenchmarkAckermann(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ackermann(uint(3), uint(1))
	}
}
