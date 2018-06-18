package digitalRoot

func Sum(value, base int) (sum int) {
	for ; value > 0; value /= base {
		sum += int(value % base)
	}
	return
}

func DigitalRoot(in, base int) (pers, root int) {
	root = int(in)
	for x := in; x >= base; x = root {
		root = Sum(x, base)
		pers += 1
	}
	return
}
