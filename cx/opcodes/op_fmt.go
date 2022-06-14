package opcodes

import (
	"fmt"
	"strconv"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/types"
)

func buildString(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) []byte {
	fmtStr := inputs[0].Get_str(prgrm)

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
			if specifiersCounter+1 == len(inputs) {
				res = append(res, []byte(fmt.Sprintf("%%!%c(MISSING)", nextCh))...)
				c++
				continue
			}

			inp := &inputs[specifiersCounter+1]
			var inpArg *ast.CXArgument
			if inp.TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				inpArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inp.TypeSignature.Meta))
			}
			switch nextCh {
			case 's':
				res = append(res, []byte(CheckForEscapedChars(inp.Get_str(prgrm)))...)
			case 'd':
				switch inp.Type {
				case types.I8:
					res = append(res, []byte(strconv.FormatInt(int64(inp.Get_i8(prgrm)), 10))...)
				case types.I16:
					res = append(res, []byte(strconv.FormatInt(int64(inp.Get_i16(prgrm)), 10))...)
				case types.I32:
					res = append(res, []byte(strconv.FormatInt(int64(inp.Get_i32(prgrm)), 10))...)
				case types.I64:
					res = append(res, []byte(strconv.FormatInt(inp.Get_i64(prgrm), 10))...)
				case types.UI8:
					res = append(res, []byte(strconv.FormatUint(uint64(inp.Get_ui8(prgrm)), 10))...)
				case types.UI16:
					res = append(res, []byte(strconv.FormatUint(uint64(inp.Get_ui16(prgrm)), 10))...)
				case types.UI32:
					res = append(res, []byte(strconv.FormatUint(uint64(inp.Get_ui32(prgrm)), 10))...)
				case types.UI64:
					res = append(res, []byte(strconv.FormatUint(inp.Get_ui64(prgrm), 10))...)
				}
			case 'f':
				switch inp.Type {
				case types.F32:
					res = append(res, []byte(strconv.FormatFloat(float64(inp.Get_f32(prgrm)), 'f', 7, 32))...)
				case types.F64:
					res = append(res, []byte(strconv.FormatFloat(inp.Get_f64(prgrm), 'f', 16, 64))...)
				}
			case 'v':
				res = append(res, []byte(ast.GetPrintableValue(prgrm, inp.FramePointer, inpArg))...)
			case 'b':
				res = append(res, []byte(strconv.FormatBool(inp.Get_bool(prgrm)))...)
			}
			c++
			specifiersCounter++
		} else {
			res = append(res, ch)
		}
	}

	if specifiersCounter != len(inputs)-1 {
		extra := "%!(EXTRA "
		// for _, inp := range expr.ProgramInput[:specifiersCounter] {
		lInps := len(inputs[specifiersCounter+1:])
		for c := 0; c < lInps; c++ {
			inp := &inputs[specifiersCounter+1+c]
			var inpArg *ast.CXArgument
			if inp.TypeSignature.Type == ast.TYPE_CXARGUMENT_DEPRECATE {
				inpArg = prgrm.GetCXArgFromArray(ast.CXArgumentIndex(inp.TypeSignature.Meta))
			}

			elt := inpArg.GetAssignmentElement(prgrm)
			typ := ""
			_ = typ
			if elt.StructType != nil {
				// then it's struct type
				typ = elt.StructType.Name
			} else {
				// then it's native type
				typ = elt.Type.Name()
			}

			if c == lInps-1 {
				extra += fmt.Sprintf("%s=%s", typ, ast.GetPrintableValue(prgrm, inp.FramePointer, elt))
			} else {
				extra += fmt.Sprintf("%s=%s, ", typ, ast.GetPrintableValue(prgrm, inp.FramePointer, elt))
			}

		}

		extra += ")"

		res = append(res, []byte(extra)...)
	}

	return res
}

func opSprintf(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_str(prgrm, string(buildString(prgrm, inputs, outputs)))
}

func opPrintf(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	fmt.Print(string(buildString(prgrm, inputs, outputs)))
}

//Only used in op_fmt.go, once
func CheckForEscapedChars(str string) []byte {
	var res []byte
	var lenStr = int(len(str))
	for c := 0; c < len(str); c++ {
		var nextCh byte
		ch := str[c]
		if c < lenStr-1 {
			nextCh = str[c+1]
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

		} else {
			res = append(res, ch)
		}
	}

	return res
}
