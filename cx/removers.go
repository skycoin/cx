package base

import ()

func (strct *CXStruct) RemoveField(fldName string) {
	if len(strct.Fields) > 0 {
		lenFlds := len(strct.Fields)
		for i, fld := range strct.Fields {
			if fld.Name == fldName {
				if i == lenFlds-1 {
					strct.Fields = strct.Fields[:len(strct.Fields)-1]
				} else {
					strct.Fields = append(strct.Fields[:i], strct.Fields[i+1:]...)
				}
				break
			}
		}
	}
}

func (expr *CXExpression) RemoveInput() {
	if len(expr.Inputs) > 0 {
		expr.Inputs = expr.Inputs[:len(expr.Inputs)-1]
	}
}

func (expr *CXExpression) RemoveOutput() {
	if len(expr.Outputs) > 0 {
		expr.Outputs = expr.Outputs[:len(expr.Outputs)-1]
	}
}
