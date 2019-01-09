package base

// CXCall ...
type CXCall struct {
	Operator     *CXFunction
	Line         int
	FramePointer int
}

/*
  Packages
*/

// CXConstant ...
type CXConstant struct {
	// native constants. only used for pre-packaged constants (e.g. math package's PI)
	// these fields are used to feed WritePrimary
	Value []byte
	Type  int
}
