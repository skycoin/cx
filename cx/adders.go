package base

// AddType ...
func (arg *CXArgument) AddType(typ string) *CXArgument {
	if typCode, found := TypeCodes[typ]; found {
		arg.Type = typCode
		size := GetArgSize(typCode)
		arg.Size = size
		arg.TotalSize = size
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_BASIC)
	} else {
		arg.Type = TYPE_UNDEFINED
	}

	return arg
}
