package ast

import "github.com/skycoin/cx/cx/constants"

/*
 * CXEXPR_TYPE enum contains CX expressions types for CXExpression struct
 */
type CXEXPR_TYPE int

const (
	CXEXPR_UNUSED CXEXPR_TYPE = iota
	CXEXPR_METHOD_CALL
	CXEXPR_STRUCT_LITERAL
	CXEXPR_ARRAY_LITERAL
	CXEXPR_SCOPE_NEW
	CXEXPR_SCOPE_DEL
)

// String returns alias for constants defined for cx edpression type
func (cxet CXEXPR_TYPE) String() string {
	return [...]string{"Unused", "MethodCall", "StructLiteral", "ArrayLiteral", "ScopeNew", "ScopeDel"}[int(cxet)]
}

// CXExpression is used represent a CX expression.
//
// All statements in CX are expressions, including for loops and other control
// flow.
//
type CXExpression struct {
	// Contents
	Inputs   []*CXArgument
	Outputs  []*CXArgument
	Label    string
	Operator *CXFunction
	Function *CXFunction
	Package  *CXPackage

	// debugging
	FileName string
	FileLine int

	// used for jmp statements
	ThenLines int
	ElseLines int

	ExpressionType CXEXPR_TYPE
}

// IsMethodCall checks if expression type is method call
func (cxe CXExpression) IsMethodCall() bool {
	return cxe.ExpressionType == CXEXPR_METHOD_CALL
}

// IsStructLiteral checks if expression type is struct literal
func (cxe CXExpression) IsStructLiteral() bool {
	return cxe.ExpressionType == CXEXPR_STRUCT_LITERAL
}

// IsArrayLiteral checks if expression type is array literal
func (cxe CXExpression) IsArrayLiteral() bool {
	return cxe.ExpressionType == CXEXPR_ARRAY_LITERAL
}

// IsBreak checks if expression type is break
func (cxe CXExpression) IsBreak() bool {
	return cxe.Operator != nil && cxe.Operator.AtomicOPCode == constants.OP_BREAK
}

// IsContinue checks if expression type is continue
func (cxe CXExpression) IsContinue() bool {
	return cxe.Operator != nil && cxe.Operator.AtomicOPCode == constants.OP_CONTINUE
}

// IsUndType checks if expression type is und type
func (cxe CXExpression) IsUndType() bool {
	return cxe.Operator != nil && IsOperator(cxe.Operator.AtomicOPCode)
}

// IsScopeNew checks if expression type is scope new
func (cxe CXExpression) IsScopeNew() bool {
	return cxe.ExpressionType == CXEXPR_SCOPE_NEW
}

// IsScopeDel checks if expression type is scope del
func (cxe CXExpression) IsScopeDel() bool {
	return cxe.ExpressionType == CXEXPR_SCOPE_DEL
}

// ----------------------------------------------------------------
//                             `CXExpression` Getters

/*
func (expr *CXExpression) GetInputs() ([]*CXArgument, error) {
	if expr.Inputs != nil {
		return expr.Inputs, nil
	}
	return nil, errors.New("expression has no arguments")

}
*/

// ----------------------------------------------------------------
//                     `CXExpression` Member handling

// AddInput ...
func (expr *CXExpression) AddInput(param *CXArgument) *CXExpression {
	// param.Package = expr.Package
	expr.Inputs = append(expr.Inputs, param)
	if param.ArgDetails.Package == nil {
		param.ArgDetails.Package = expr.Package
	}
	return expr
}

// RemoveInput ...
func (expr *CXExpression) RemoveInput() {
	if len(expr.Inputs) > 0 {
		expr.Inputs = expr.Inputs[:len(expr.Inputs)-1]
	}
}

// AddOutput ...
func (expr *CXExpression) AddOutput(param *CXArgument) *CXExpression {
	// param.Package = expr.Package
	expr.Outputs = append(expr.Outputs, param)
	if param.ArgDetails.Package == nil {
		param.ArgDetails.Package = expr.Package
	}
	return expr
}

// RemoveOutput ...
func (expr *CXExpression) RemoveOutput() {
	if len(expr.Outputs) > 0 {
		expr.Outputs = expr.Outputs[:len(expr.Outputs)-1]
	}
}

// AddLabel ...
func (expr *CXExpression) AddLabel(lbl string) *CXExpression {
	expr.Label = lbl
	return expr
}
