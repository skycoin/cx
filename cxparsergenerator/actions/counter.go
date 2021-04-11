package actions

import "fmt"

var (
	GenSymCounter int
)

// MakeGenSym ...
func MakeGenSym(name string) string {
	gensym := fmt.Sprintf("%s_%d", name, GenSymCounter)
	GenSymCounter++

	return gensym
}
