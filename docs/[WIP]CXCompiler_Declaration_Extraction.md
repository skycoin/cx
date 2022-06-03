# CXCompiler: Declaration Extraction
Goals:
- [ ] Function to extract all global declarations from source file
- [ ] Function to extract all alias and enum declarations from source file
- [ ] Function to extract all struct declarations from source files
- [ ] Function to extract all function declarations from source file
- [ ] Function for redeclaration checks

## Global Declaration
Extract all globals declarations but not assignments to globals

### Global Declaration Struct
- PackageID
- FileID
- StartOffset
- Length
- Name

Data Output:
`An array of GlobalDeclarationStructs`

## Aliases And Enums Declaration
Extract all enums and declarations declarations 

### Enums Struct
- PackageID
- FileID
- StartOffset
- Length
- Type
- Value
- Name

Data Output:
`An array of EnumsDeclarationStructs`

## Struct Declaration
Extract all struct declarations but not including its body

### Struct Declaration Struct
- PackageID
- FileID
- StartOffset
- Length
- Name

Data Output:
`An array of StructDeclarationStructs`


## Function Declaration
Extract all function declarations but not including its body

### Function Declaration Struct
- PackageID
- FileID
- StartOffset
- Length
- Name

Data Output:
`An array of FunctionDeclarationStructs`

## Redeclaration Checks:
- Verify that no redeclarations are occurring
- Verify that same type has not been declared twice 
- Verify that same function has not been declared twice
- Verify that enums and aliases do not conflict

---
## Follow-up milestones:
- [ ] Use Goroutine to extract declarations (One goroutine per file)
- [ ] Add the data on a channel instead of adding on an array
    - One channel for global declaration
    - One channel for aliases and enums declaration
    - One channel for struct declaration
    - One channel for function declaration

---
Questions: 

    Do we parse here or wait till the next step?

    The structs are not defined until next step. i.e. a global variable that is a struct type that is defined or imported from another package; the struct id/type does not exist yet until this stage is completed, so cannot parse at this stage.