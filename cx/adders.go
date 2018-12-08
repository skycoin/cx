package base

import (
	// "fmt"
	// "github.com/skycoin/skycoin/src/cipher/encoder"
)

func (arg *CXArgument) AddType(typ string) *CXArgument {
	// arg.Typ = typ
	if typCode, found := TypeCodes[typ]; found {
		arg.Type = typCode
		size := GetArgSize(typCode)
		arg.Size = size
		arg.TotalSize = size
	} else {
        arg.Type = TYPE_UNDEFINED
    }

	return arg
}

func (expr *CXExpression) AddInput(param *CXArgument) *CXExpression {
	// param.Package = expr.Package
	expr.Inputs = append(expr.Inputs, param)
	if param.Package == nil {
		param.Package = expr.Package
	}
	return expr
}

func (expr *CXExpression) AddOutput(param *CXArgument) *CXExpression {
	// param.Package = expr.Package
	expr.Outputs = append(expr.Outputs, param)
	param.Package = expr.Package
	return expr
}

func (expr *CXExpression) AddLabel(lbl string) *CXExpression {
	expr.Label = lbl
	return expr
}

// func (expr *CXExpression) AddOutputName(outName string) *CXExpression {
// 	if len(expr.Operator.Outputs) > 0 {
// 		nextOutIdx := len(expr.Outputs)

// 		var typ string
// 		if expr.Operator.Name == ID_FN || expr.Operator.Name == INIT_FN {
// 			var tmp string
// 			encoder.DeserializeRaw(*expr.Inputs[0].Value, &tmp)

// 			if expr.Operator.Name == INIT_FN {
// 				// then tmp is the type (e.g. initDef("i32") to initialize an i32)
// 				typ = tmp
// 			} else {
// 				var err error
// 				// then tmp is an identifier
// 				if typ, err = GetIdentType(tmp, expr.FileLine, expr.FileName, expr.Program); err == nil {
// 				} else {
// 					panic(err)
// 				}
// 			}
// 		} else {
// 			typ = expr.Operator.Outputs[nextOutIdx].Typ
// 		}

// 		outDef := MakeArgument(outName, "", -1).AddValue(MakeDefaultValue(expr.Operator.Outputs[nextOutIdx].Typ)).AddType(typ)

// 		outDef.Package = expr.Package
// 		outDef.Program = expr.Program

// 		expr.Outputs = append(expr.Outputs, outDef)
// 	}

// 	return expr
// }
