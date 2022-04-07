# CXStruct Definition

## CXStruct Usage
These have the same layout format as a “struct” definition
1. A struct
2. the variable layout of a function
3. the inputs to a function
4. the outputs of a function

## CXStruct
- NumFields int
- Fields []CXStructField

## CXStructField
- FieldNameStringId  int
- FieldOffset int
- FieldType enum
    - CXAtomic
    - CXAtomicPointer
    - CXStruct
    - CXStructPointer
    - CXArrayPointer
    - CXArrayAtomic
    - CXArrayStruct
    - CXMap
    - CXComplexType
    - etc
- FieldMeta enum 
    - if FieldType is a cx atomic pointer, enum of the type, is struct pointer meta will be struct field, if array struct pointer will be type of struct id, array to atomic will be the type atomic type
    - if FieldType is CXAtomicPointer, CXAtomic
    - if FieldType is CXStructPointer, CXStruct
    - if FieldType is ArrayAtomic, AtomicType
    
Note:
- Implement SizeOf() method for CXStructField, panics if asked for type for a complex/non fixed size, only sizeof simple types
- Make one struct for complex type.
- glob all complex types which will just be cxargs
- From complex types, move types one by one to FieldTypes
