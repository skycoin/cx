// +build !base,!extra,!full

package cxcore

// For the parser. These shouldn't be used in the runtime for performance reasons
var (
	ConstNames = map[int]string{}
	ConstCodes = map[string]int{}
	Constants  = map[int]CXConstant{}
)
