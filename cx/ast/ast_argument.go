package ast

// GetAssignmentElement ...
func GetAssignmentElement(arg *CXArgument) *CXArgument {
	if len(arg.Fields) > 0 {
		return arg.Fields[len(arg.Fields)-1]
	}
	return arg

}

// GetType ...
func GetType(arg *CXArgument) int {
	fieldCount := len(arg.Fields)
	if fieldCount > 0 {
		return GetType(arg.Fields[fieldCount-1])
	}

	return arg.Type
}

