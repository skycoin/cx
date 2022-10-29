package ast

import (
	"fmt"
	"os"

	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

type CXStructIndex int

// CXStruct is used to represent a CX struct.
type CXStruct struct {
	// Metadata
	Index   int
	Name    string         // Name of the struct
	Package CXPackageIndex // The package this struct belongs to

	// Contents
	Fields []CXTypeSignatureIndex // The fields of the struct
}

// ----------------------------------------------------------------
//                             `CXStruct` Getters

// GetField ...
func (strct *CXStruct) GetField(prgrm *CXProgram, name string) (*CXArgument, error) {
	for _, typeSignatureIdx := range strct.Fields {
		typeSignature := prgrm.GetCXTypeSignatureFromArray(typeSignatureIdx)
		if typeSignature.Name == name {
			// If type is struct
			if typeSignature.Type == TYPE_STRUCT {
				structDetails := prgrm.GetCXTypeSignatureStructFromArray(typeSignature.Meta)
				return &CXArgument{
					Name:                  typeSignature.Name,
					Type:                  types.STRUCT,
					DeclarationSpecifiers: []int{constants.DECL_STRUCT},
					StructType:            structDetails.StructType,
					Offset:                typeSignature.Offset,
				}, nil
				// If type is not cxargument deprecate
			} else if typeSignature.Type != TYPE_CXARGUMENT_DEPRECATE {
				return &CXArgument{
					Name:                  typeSignature.Name,
					Type:                  types.Code(typeSignature.Meta),
					DeclarationSpecifiers: []int{constants.DECL_BASIC},
					StructType:            nil,
					Offset:                typeSignature.Offset,
				}, nil
			}

			// If type is cxargument deprecate
			fldIdx := typeSignature.Meta
			return &prgrm.CXArgs[fldIdx], nil
		}
	}

	return nil, fmt.Errorf("field '%s' not found in struct '%s'", name, strct.Name)
}

// ----------------------------------------------------------------
//                     `CXStruct` Member handling

// MakeStruct ...
func MakeStruct(name string) *CXStruct {
	return &CXStruct{
		Name: name,
	}
}

// AddField ...
func (strct *CXStruct) AddField(prgrm *CXProgram, cxArgument *CXArgument) *CXStruct {
	// Check if field already exist
	for _, typeSignatureIdx := range strct.Fields {
		typeSignature := prgrm.GetCXTypeSignatureFromArray(typeSignatureIdx)
		if typeSignature.Name == cxArgument.Name {
			// fldIdx := typeSignature.Meta
			// fmt.Printf("%s : duplicate field", CompilationError(prgrm.CXArgs[fldIdx].ArgDetails.FileName, prgrm.CXArgs[fldIdx].ArgDetails.FileLine))
			fmt.Println("duplicate field")
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
	}

	newCXTypeSignature := GetCXTypeSignatureRepresentationOfCXArg_ForStructs(prgrm, cxArgument)
	typeSignatureIdx := prgrm.AddCXTypeSignatureInArray(newCXTypeSignature)
	strct.Fields = append(strct.Fields, typeSignatureIdx)

	// TODO: Found out the effect and completely remove this.
	// Dont remove this yet as removing this gives another error on cxfx
	// even though all cx tests passed already.
	// numFlds := len(strct.Fields)
	// if numFlds > 2 {
	// 	// Pre-compiling the offset of the field.
	// 	lastTypeSignature := strct.Fields[numFlds-2]
	// 	currentTypeSignature := strct.Fields[numFlds-1]
	// 	if lastTypeSignature.Type == TYPE_CXARGUMENT_DEPRECATE && currentTypeSignature.Type == TYPE_CXARGUMENT_DEPRECATE {
	// 		lastFldIdx := lastTypeSignature.Meta
	// 		currFldIdx := currentTypeSignature.Meta
	// 		prgrm.CXArgs[currFldIdx].Offset = prgrm.CXArgs[lastFldIdx].Offset + prgrm.CXArgs[lastFldIdx].TotalSize
	// 	} else if lastTypeSignature.Type == TYPE_CXARGUMENT_DEPRECATE && currentTypeSignature.Type != TYPE_CXARGUMENT_DEPRECATE {
	// 		lastFldIdx := lastTypeSignature.Meta
	// 		currentTypeSignature.Offset = prgrm.CXArgs[lastFldIdx].Offset + prgrm.CXArgs[lastFldIdx].TotalSize
	// 	} else if lastTypeSignature.Type != TYPE_CXARGUMENT_DEPRECATE && currentTypeSignature.Type == TYPE_CXARGUMENT_DEPRECATE {
	// 		currFldIdx := currentTypeSignature.Meta
	// 		prgrm.CXArgs[currFldIdx].Offset = lastTypeSignature.Offset + lastTypeSignature.GetSize(prgrm, false)
	// 	} else if lastTypeSignature.Type != TYPE_CXARGUMENT_DEPRECATE && currentTypeSignature.Type != TYPE_CXARGUMENT_DEPRECATE {
	// 		currentTypeSignature.Offset = lastTypeSignature.Offset + lastTypeSignature.GetSize(prgrm, false)
	// 	}
	// }

	return strct
}

// AddField ...
func (strct *CXStruct) AddField_TypeSignature(prgrm *CXProgram, fieldIdx CXTypeSignatureIndex) *CXStruct {
	// Check if field already exist
	for _, typeSignatureIdx := range strct.Fields {
		typeSignature := prgrm.GetCXTypeSignatureFromArray(typeSignatureIdx)
		field := prgrm.GetCXTypeSignatureFromArray(fieldIdx)
		if typeSignature.Name == field.Name {
			// fldIdx := typeSignature.Meta
			// fmt.Printf("%s : duplicate field", CompilationError(prgrm.CXArgs[fldIdx].ArgDetails.FileName, prgrm.CXArgs[fldIdx].ArgDetails.FileLine))
			fmt.Println("duplicate field")
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
	}

	strct.Fields = append(strct.Fields, fieldIdx)

	return strct
}

func (strct *CXStruct) AddField_CXAtomicOps(prgrm *CXProgram, fieldIdx CXTypeSignatureIndex) *CXStruct {
	strct.Fields = append(strct.Fields, fieldIdx)

	return strct
}

func (strct *CXStruct) AddField_Globals_CXAtomicOps(prgrm *CXProgram, cxArgIdx CXArgumentIndex) *CXStruct {
	cxArgument := prgrm.GetCXArgFromArray(cxArgIdx)

	// Check if field already exist
	for _, typeSignatureIdx := range strct.Fields {
		typeSignature := prgrm.GetCXTypeSignatureFromArray(typeSignatureIdx)
		if typeSignature.Name == cxArgument.Name {
			// fldIdx := typeSignature.Meta
			// fmt.Printf("%s : duplicate field", CompilationError(prgrm.CXArgs[fldIdx].ArgDetails.FileName, prgrm.CXArgs[fldIdx].ArgDetails.FileLine))
			fmt.Println("duplicate field")
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
	}

	newCXTypeSignature := GetCXTypeSignatureRepresentationOfCXArg(prgrm, cxArgument)
	newCXTypeSignatureIdx := prgrm.AddCXTypeSignatureInArray(newCXTypeSignature)
	strct.Fields = append(strct.Fields, newCXTypeSignatureIdx)

	return strct
}

// RemoveField ...
// func (strct *CXStruct) RemoveField(prgrm *CXProgram, fldName string) {
// 	if len(strct.Fields) > 0 {
// 		lenFlds := len(strct.Fields)
// 		for i, fldIdx := range strct.Fields {
// 			if prgrm.CXArgs[fldIdx].Name == fldName {
// 				if i == lenFlds-1 {
// 					strct.Fields = strct.Fields[:len(strct.Fields)-1]
// 				} else {
// 					strct.Fields = append(strct.Fields[:i], strct.Fields[i+1:]...)
// 				}
// 				break
// 			}
// 		}
// 	}
// }

func (strct *CXStruct) GetStructSize(prgrm *CXProgram) types.Pointer {
	var structSize types.Pointer
	for _, typeSignatureIdx := range strct.Fields {
		typeSignature := prgrm.GetCXTypeSignatureFromArray(typeSignatureIdx)
		structSize += typeSignature.GetSize(prgrm, false)
	}

	return structSize
}
