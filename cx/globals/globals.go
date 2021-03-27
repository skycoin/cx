package globals

import (
	"fmt"
	"os"
)

// Var
var (
	HeapOffset    int
)

//Path is only used by os module and only to get working directory
//Path and working directory should not be hard coded into program struct (etc, when serialized)
//Working directory is property of executable and can be retrieved with golang library
var CxProgramPath string = ""

var CXPATH = os.Getenv("CXPATH") + "/"
var BINPATH = CXPATH + "bin/" // TODO @evanlinjin: Not used.
var PKGPATH = CXPATH + "pkg/" // TODO @evanlinjin: Not used.
var SRCPATH = CXPATH + "src/"

var (
	GenSymCounter int
)

// MakeGenSym ...
func MakeGenSym(name string) string {
	gensym := fmt.Sprintf("%s_%d", name, GenSymCounter)
	GenSymCounter++

	return gensym
}

