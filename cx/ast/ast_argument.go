package ast

// GetAssignmentElement ...
func GetAssignmentElement(arg *CXArgument) *CXArgument {
	if len(arg.Fields) > 0 {
		return arg.Fields[len(arg.Fields)-1]
	}
	return arg

}
