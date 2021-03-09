package ackermann

import (
	"fmt"
)

func ackermann(m, n uint) uint {
	switch uint(0) {
	case m:
		return n + uint(1)
	case n:
		return ackermann(m-uint(1), uint(1))
	}
	return ackermann(n-uint(1), ackermann(m, n-uint(1)))
}

func main() {
	fmt.Printf("Ackermann Function")
	fmt.Printf("Test with (3, 1)")
	fmt.Printf("Result: ")
	fmt.Println(ackermann(uint(3), uint(1)))
}
