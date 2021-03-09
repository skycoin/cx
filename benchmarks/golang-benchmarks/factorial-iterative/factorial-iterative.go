package factorialIter

func factorial(in int) (out int) {
	out = 1
	for i := 1; i < in; i += 1 {
		out *= (i + 1)
	}
	return
}
