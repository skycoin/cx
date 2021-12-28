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
	ExpressionType CXEXPR_TYPE

	// For new CX AST Format
	Index int
	Type  CXOPERATION_TYPE
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
func (cxe CXExpression) IsBreak(prgrm *CXProgram) bool {
	cxAtomicOp, _, _, err := prgrm.GetOperation(&cxe)
	if err != nil {
		panic(err)
	}

	return cxAtomicOp.Operator != nil && cxAtomicOp.Operator.AtomicOPCode == constants.OP_BREAK
}

// IsContinue checks if expression type is continue
func (cxe CXExpression) IsContinue(prgrm *CXProgram) bool {
	cxAtomicOp, _, _, err := prgrm.GetOperation(&cxe)
	if err != nil {
		panic(err)
	}

	return cxAtomicOp.Operator != nil && cxAtomicOp.Operator.AtomicOPCode == constants.OP_CONTINUE
}

// IsUndType checks if expression type is und type
func (cxe CXExpression) IsUndType(prgrm *CXProgram) bool {
	cxAtomicOp, _, _, err := prgrm.GetOperation(&cxe)
	if err != nil {
		panic(err)
	}

	return cxAtomicOp.Operator != nil && IsOperator(cxAtomicOp.Operator.AtomicOPCode)

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

func (cxe CXExpression) GetLabel(prgrm *CXProgram) string {
	cxAtomicOp, _, _, err := prgrm.GetOperation(&cxe)
	if err != nil {
		panic(err)
	}

	return cxAtomicOp.Label
}
