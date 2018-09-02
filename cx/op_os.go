package base

import (
	// "fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func op_os_GetWorkingDirectory(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)
	
	byts := encoder.Serialize(expr.Program.Path)
	WriteObject(out1Offset, byts)
}
