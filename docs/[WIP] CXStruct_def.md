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
    - CXAtomicI32
    - CXAtomicI32Pointer
    - CXStruct
    - CXStructPointer
    - CXArrayPointer
    - CXArrayAtomic
    - CXArrayStruct
    - CXMap
    - CXComplexType
    - etc
- Meta enum 
    - if FieldType is CXAtomicI32Pointer, CXAtomicI32
    - if FieldType is CXStructPointer, CXStruct
    - if FieldType is ArrayAtomic, AtomicType
    
Note:
- Implement SizeOf() method for CXStructField, panics if asked for type for a complex/non fixed size, only sizeof simple types
- Make one struct for complex type.
- glob all complex types which will just be cxargs
- From complex types, move types one by one to FieldTypes
