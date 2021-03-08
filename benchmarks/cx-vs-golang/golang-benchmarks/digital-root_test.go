package digitalRoot

import "testing"

func BenchmarkDigitalRoot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DigitalRoot(79563, 10)
	}
}
