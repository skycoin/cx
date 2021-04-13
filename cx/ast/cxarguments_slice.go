package ast

type CXArgumentSlice struct {
	IsSlice      bool
	IsArray      bool
	IsArrayFirst bool // and then dereference
}
