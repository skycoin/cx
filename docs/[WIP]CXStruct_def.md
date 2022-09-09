# CXStruct Definition

## CXStruct Usage
These have the same layout format as a “struct” definition
1. A struct
2. the variable layout of a function
3. the inputs to a function
4. the outputs of a function

## CXStruct
- StructId int
- NameStringId  int
- PackageNameStringId int
- Fields []CXTypeSignature

## CXTypeSignature
- NameStringId  int
- Offset int
- Type enum
    - Atomic
    - PointerAtomic
    - ArrayAtomic
    - ArrayPointerAtomic
    - SliceAtomic
    - SlicePointerAtomic
    - Struct
    - PointerStruct
    - ArrayStruct
    - ArrayPointerStruct
    - SliceStruct
    - SlicePointerStruct
    - Complex
    - PointerComplex
    - ArrayComplex
    - ArrayPointerComplex
    - SliceComplex
    - SlicePointerComplex
- Meta int
    - if Type is Atomic, the atomic type
    - if Type is Struct, the struct id
    - if Type is Complex, the complex id
    - if Type is Array, the type signature array id

## CXTypeSignature_Array
- Type int
- Length int
    
Note:
- Implement SizeOf() method for CXTypeSignature, panics if asked for type for a complex/non fixed size, only sizeof simple types
- Make one struct for complex type.
- glob all complex types which will just be cxargs
- From complex types, move types one by one to FieldTypes

## Progress
- [ ] CXStruct Implementation
    - [x] CXStruct for function Inputs
    - [x] CXStruct for function Outputs
    - [ ] CXStruct for function variable layout
    - [ ] CXStruct for struct definitions
- [ ] CXTypeSignature Implementation
    - [x] atomic
    - [x] PointerAtomic
    - [ ] ArrayAtomic
    - [ ] ArrayPointerAtomic
    - [ ] SliceAtomic
    - [ ] SlicePointerAtomic
    - [ ] Struct
    - [ ] PointerStruct
    - [ ] ArrayStruct
    - [ ] ArrayPointerStruct
    - [ ] SliceStruct
    - [ ] SlicePointerStruct
    - [ ] Complex
    - [ ] PointerComplex
    - [ ] ArrayComplex
    - [ ] ArrayPointerComplex
    - [ ] SliceComplex
    - [ ] SlicePointerComplex