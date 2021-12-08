package ast

/*
 * CXOPERATION_TYPE enum contains CX operations types for CXOperation struct
 */
type CXOPERATION_TYPE int

const (
	UNUSED CXOPERATION_TYPE = iota
	CX_ATOMIC_OPERATOR
	CX_ARGUMENT
	CX_LINE
)

// type CXOperation struct {
// 	Index int
// 	Type  CXOPERATION_TYPE
// }
