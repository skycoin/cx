package base

import (
	"fmt"
	// "github.com/skycoin/skycoin/src/cipher/encoder"
)

func op_assert_value(expr *CXExpression, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	

	var byts1, byts2 []byte
	
	if inp1.Type == TYPE_STR {
		byts1 = []byte(ReadStr(fp, inp1))
		byts2 = []byte(ReadStr(fp, inp2))
	} else {
		byts1 = ReadMemory(GetFinalOffset(fp, inp1), inp1)
		byts2 = ReadMemory(GetFinalOffset(fp, inp2), inp2)
	}

	var same bool
	same = true

	if len(byts1) != len(byts2) {
		same = false
		fmt.Println("byts1", inp1.Type, byts1, inp1.Name, inp1.Size, inp1.TotalSize)
		fmt.Println("byts2", inp2.Type, byts2, inp2.Name)
	} else {
		for i, byt := range byts1 {
			if byt != byts2[i] {
				same = false
				fmt.Println("byts1", inp1.Type, byts1, inp1.Name, inp1.Size, inp1.TotalSize)
				fmt.Println("byts2", inp2.Type, byts2, inp2.Name)
			}
		}
	}
	

	var message string
	message = ReadStr(fp, inp3)

	if !same {
		if message != "" {
			fmt.Printf("%s: %d: result was not equal to the expected value; %s\n", expr.FileName, expr.FileLine, message)
		} else {
			fmt.Printf("%s: %d: result was not equal to the expected value\n", expr.FileName, expr.FileLine)
		}
	} else {

	}
}
