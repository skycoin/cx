// +build base extra full

package base

import (
	// "fmt"
	"io/ioutil"
	"os"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

var openFiles map[string]*os.File = make(map[string]*os.File, 0)

func op_os_ReadFile(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	
	_ = out1

	if byts, err := ioutil.ReadFile(ReadStr(fp, inp1)); err == nil {
		_ = byts
		// sByts := encoder.Serialize(byts)
		// assignOutput(0, sByts, "[]byte", expr, call)
	} else {
		panic(err)
	}
}

func op_os_Open(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	name := ReadStr(fp, inp1)
	if file, err := os.Open(name); err == nil {
		openFiles[name] = file
	} else {
		panic(err)
	}
}

func op_os_Close(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	name := ReadStr(fp, inp1)
	if file, ok := openFiles[name]; ok {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}
}

func op_os_GetWorkingDirectory(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)
	
	byts := encoder.Serialize(PROGRAM.Path)
	WriteObject(out1Offset, byts)
}
