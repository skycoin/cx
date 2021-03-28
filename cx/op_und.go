package cxcore

import (
	"bufio"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"os"
	"strings"
)

//Reads input from Standard Inpuit
func OpReadStdin(inputs []ast.CXValue, outputs []ast.CXValue) {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
    if err != nil {
		panic(constants.CX_INTERNAL_ERROR)
	}

	// text = strings.Trim(text, " \n")
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)
    outputs[0].Set_str(text)
}