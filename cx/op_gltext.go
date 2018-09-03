package base

import (
	"github.com/go-gl/gltext"
)

var fonts map[string]*gltext.Font = make(map[string]*gltext.Font, 0)

func op_gltext_LoadTrueType(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4, inp5, inp6 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4], expr.Inputs[5]

	if theFont, err := gltext.LoadTruetype(openFiles[ReadStr(fp, inp2)], ReadI32(fp, inp3), rune(ReadI32(fp, inp4)), rune(ReadI32(fp, inp5)), gltext.Direction(ReadI32(fp, inp6))); err == nil {
		fonts[ReadStr(fp, inp1)] = theFont
	} else {
		panic(err)
	}
}

func op_gltext_Printf(expr *CXExpression, fp int) {
	inp1, inp2, inp3, inp4 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3]

	if err := fonts[ReadStr(fp, inp1)].Printf(ReadF32(fp, inp2), ReadF32(fp, inp3), ReadStr(fp, inp4)); err != nil {
		panic(err)
	}
}
