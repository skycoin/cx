package cxbenchmark_test

import (
	"os/exec"
	"testing"
)

func BenchmarkCX(b *testing.B) {
	for i := 0; i < b.N; i++ {

		cmd := exec.Command("./bin/cx", "./test_files/test.cx")
		_, err := cmd.CombinedOutput()
		if err != nil {
			b.Fatal(err)
		}

	}
}
