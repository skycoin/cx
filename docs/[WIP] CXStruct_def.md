# CXStruct Definition

## CXStruct Usage
These have the same layout format as a “struct” definition
1. A struct
2. the variable layout of a function
3. the inputs to a function
4. the outputs of a function

## CXStruct
- num fields 
- field string
- field offset 
- field size
- field type
- TypeID
    - CXAtomic
    - CXAtomicPointer
    - CXStruct
    - CXStructPointer
    - CXArray
    - CXArrayPointer
    - etc
- TypeField
    - If CXAtomic, then CXAtomic type
    - if StructPointer, then the pointer id
    - If CXAtomicPointer, then the type of CXAtomic
    - etc

