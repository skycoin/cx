package base

import (
// "fmt"
// "github.com/skycoin/skycoin/src/cipher/encoder"
)

func (arg *CXArgument) AddType(typ string) *CXArgument {
	// arg.Typ = typ
	if typCode, found := TypeCodes[typ]; found {
		arg.Type = typCode
		size := GetArgSize(typCode)
		arg.Size = size
		arg.TotalSize = size
	} else {
		arg.Type = TYPE_UNDEFINED
	}

	return arg
}
