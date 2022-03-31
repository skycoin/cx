package actions

import "fmt"

var (
	GenSymCounter int
)

// MakeGenSym generates generated tmp name used for temporary variables.
func MakeGenSym(name string) string {
	gensym := fmt.Sprintf("%s_%d", name, GenSymCounter)
	GenSymCounter++

	return gensym
}
