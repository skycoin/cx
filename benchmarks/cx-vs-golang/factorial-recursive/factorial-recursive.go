package factorialRec

func factorial(in int) int {
	if in == 0 {
		return 1
	} else {
		return in * factorial(in-1)
	}
}
