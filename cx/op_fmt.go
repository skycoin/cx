package cxcore

import (
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/util2"
	"strconv"
)

func buildString(expr *ast.CXExpression, fp int) []byte {
	inp1 := expr.Inputs[0]

	fmtStr := ast.ReadStr(fp, inp1)

	var res []byte
	var specifiersCounter int
	var lenStr = int(len(fmtStr))

	for c := 0; c < len(fmtStr); c++ {
		var nextCh byte
		ch := fmtStr[c]
		if c < lenStr-1 {
			nextCh = fmtStr[c+1]
		}
		if ch == '\\' {
			switch nextCh {
			case '%':
				c++
				res = append(res, nextCh)
				continue
			case 'n':
				c++
				res = append(res, '\n')
				continue
			default:
				res = append(res, ch)
				continue
			}
		}
		if ch == '%' {
			if specifiersCounter+1 == len(expr.Inputs) {
				res = append(res, []byte(fmt.Sprintf("%%!%c(MISSING)", nextCh))...)
				c++
				continue
			}

			inp := expr.Inputs[specifiersCounter+1]
			switch nextCh {
			case 's':
				res = append(res, []byte(util2.CheckForEscapedChars(ast.ReadStr(fp, inp)))...)
			case 'd':
				switch inp.Type {
				case constants.TYPE_I8:
					res = append(res, []byte(strconv.FormatInt(int64(ast.ReadI8(fp, inp)), 10))...)
				case constants.TYPE_I16:
					res = append(res, []byte(strconv.FormatInt(int64(ast.ReadI16(fp, inp)), 10))...)
				case constants.TYPE_I32:
					res = append(res, []byte(strconv.FormatInt(int64(ast.ReadI32(fp, inp)), 10))...)
				case constants.TYPE_I64:
					res = append(res, []byte(strconv.FormatInt(ast.ReadI64(fp, inp), 10))...)
				case constants.TYPE_UI8:
					res = append(res, []byte(strconv.FormatUint(uint64(ast.ReadUI8(fp, inp)), 10))...)
				case constants.TYPE_UI16:
					res = append(res, []byte(strconv.FormatUint(uint64(ast.ReadUI16(fp, inp)), 10))...)
				case constants.TYPE_UI32:
					res = append(res, []byte(strconv.FormatUint(uint64(ast.ReadUI32(fp, inp)), 10))...)
				case constants.TYPE_UI64:
					res = append(res, []byte(strconv.FormatUint(ast.ReadUI64(fp, inp), 10))...)
				}
			case 'f':
				switch inp.Type {
				case constants.TYPE_F32:
					res = append(res, []byte(strconv.FormatFloat(float64(ast.ReadF32(fp, inp)), 'f', 7, 32))...)
				case constants.TYPE_F64:
					res = append(res, []byte(strconv.FormatFloat(ast.ReadF64(fp, inp), 'f', 16, 64))...)
				}
			case 'v':
				res = append(res, []byte(ast.GetPrintableValue(fp, inp))...)
			}
			c++
			specifiersCounter++
		} else {
			res = append(res, ch)
		}
	}

	if specifiersCounter != len(expr.Inputs)-1 {
		extra := "%!(EXTRA "
		// for _, inp := range expr.ProgramInput[:specifiersCounter] {
		lInps := len(expr.Inputs[specifiersCounter+1:])
		for c := 0; c < lInps; c++ {
			inp := expr.Inputs[specifiersCounter+1+c]
			elt := ast.GetAssignmentElement(inp)
			typ := ""
			_ = typ
			if elt.CustomType != nil {
				// then it's custom type
				typ = elt.CustomType.Name
			} else {
				// then it's native type
				typ = constants.TypeNames[elt.Type]
			}

			if c == lInps-1 {
				extra += fmt.Sprintf("%s=%s", typ, ast.GetPrintableValue(fp, elt))
			} else {
				extra += fmt.Sprintf("%s=%s, ", typ, ast.GetPrintableValue(fp, elt))
			}

		}

		extra += ")"

		res = append(res, []byte(extra)...)
	}

	return res
}

func opSprintf(expr *ast.CXExpression, fp int) {
	ast.WriteString(fp, string(buildString(expr, fp)), expr.Outputs[0])
}

func opPrintf(expr *ast.CXExpression, fp int) {
	fmt.Print(string(buildString(expr, fp)))
}
