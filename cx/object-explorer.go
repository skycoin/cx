// +build base extra full

package base

import (
	"fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func GetAllObjects(prgrm *CXProgram) {
	fp := 0

	for c := 0; c <= prgrm.CallCounter; c++ {
		op := prgrm.CallStack[c].Operator

		for _, ptr := range op.ListOfPointers {
			var heapOffset int32
			encoder.DeserializeAtomic(prgrm.Memory[fp+ptr.Offset:fp+ptr.Offset+TYPE_POINTER_SIZE], &heapOffset)

			var byts []byte

			if ptr.CustomType != nil {
				// then it's a pointer to a struct
				// use CustomStruct to match the fields against the bytes
				// for _, fld := range ptr.Fields {

				// }

				byts = prgrm.Memory[int(heapOffset)+OBJECT_HEADER_SIZE : int(heapOffset)+OBJECT_HEADER_SIZE+ptr.CustomType.Size]
			}

			if len(ptr.Lengths) > 0 {

			}

			// var currLengths []int
			// var currCustom *CXStruct

			// for c := len(ptr.DeclarationSpecifiers) - 1; c >= 0; c-- {
			// 	// we need to go backwards in here

			// 	switch ptr.DeclarationSpecifiers[c] {
			// 	case DECL_POINTER:
			// 		// we might not need to do anything
			// 	case DECL_ARRAY:
			// 		currLengths = ptr.Lengths
			// 	case DECL_SLICE:
			// 	case DECL_STRUCT:
			// 		currCustom = ptr.CustomType
			// 	case DECL_BASIC:
			// 	}
			// }

			// if len(ptr.Lengths) > 0 {
			// 	fmt.Println("ARRAY")
			// }

			// if ptr.CustomType != nil {
			// 	fmt.Println("STRUCT")
			// }

			fmt.Println("declarat", ptr.DeclarationSpecifiers)

			fmt.Println("obj", ptr.Name, ptr.CustomType, prgrm.Memory[heapOffset:int(heapOffset)+op.Size], byts)
		}

		fp += op.Size
	}
}
