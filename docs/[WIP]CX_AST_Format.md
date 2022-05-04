# CX AST Format Specification

The CX AST, or abstract syntax tree, is the intermediate representation of a CX program which the CX runtime executes. It consists of multiple flavors of nodes, each one broken down in the following subsections.

## CXArgument
Note: CXArgument is to be deprecated after the new CXStruct def will be fully implemented.

The object is composed of the following fields:
* Name - the symbol of the CXArgument.
* Package - the CXPackage this CXArgument resides in.
* Lengths - an array of integers which specify the dimensions of a multidimensional array (first element is the outermost array length, etc).
* DereferenceOperations - an array of integers which specify, in bytes, how much in program memory this CXArgument is offsetted from its base CXArgument. For example, for a struct field, this is the struct field’s offset in bytes. For an array element in an array of an array, this is the first array’s offset times the second array’s size times the second array’s element type size, plus the offset in the second array times the second array’s element type size.
* DeclarationSpecifiers - an array of integer constants which specify the type of the CXArgument. For type `[][]ui8`, this is `[]int{TYPE_SLICE, TYPE_SLICE, TYPE_UI8}`.
* Indexes - an array of CXArguments which correspond to how to dereference the current CXArgument. For example, for an array dereference with index “i”, i is also a CXArgument. So the array dereference CXArgument stores i in Indexes.
* Fields - an array of CXArguments which store, in order, the field dereferences of structs or nested structs. So this would be `struct_name.field_of_first_struct.field_of_this_struct`.
* Inputs, Outputs  - the input and output parameters of a CXArgument of type `TYPE_FUNC`. Usually used for function calls.
* Type - the specific integer constant that reflects the type of this CXArgument (one of `TYPE_SLICE`, `TYPE_UI8`, etc.)
* PointerTargetType -
* Size - the size of the underlying basic type.
* TotalSize - total size of an array (only used for Array). 
* Offset - the location in the program memory this CXArgument resides in.
* PassBy - an int constant representing how the variable is passed - pass by value, or pass by reference.
* StructType - a CXStruct representing the custom type this variable might be (if struct).
* IsSlice, IsStruct - name says it all.
* IsInnerArg - if this is a package global, is this CXArgument representing the actual global variable from that package?
* IsLocalDeclaration - is this CXArgument bounded in lifetime by scope (as opposed to being a global variable)?
* IsShortDeclaration - is this CXArgument the result of a `CASSIGN` operation (`:=`)?
* IsInnerReference - is this CXArgument a reference to the field or element of another CXArgument? (`&array[0]` or `&struct.field`)
* PreviouslyDeclared - used by compiler to check if this variable has been declared yet or not, or if there are duplicate declarations.
* DoesEscape - should this variable be kept alive after the scope ends? (for example, a function returning a pointer to a variable created in the function should keep that variable alive after the scope ends, hence, this should be set to true then).

## CXAtomicOperator
CXAtomicOperator have the following fields:
* Inputs - input arguments of the atomic operator.
* Outputs - output arguments of the atomic operator.
* Operator - the operator of the expression.
* Function - the function where the atomic operator expression belongs.
* Package - the package where it belongs. 
* Label - label used for goto expressions.
* ThenLines - used for if-else and goto/jump expressions.
* ElseLines - used for if-else and goto/jump expressions.

## CXLine
CXLine is a NOP(No operation). Only used for information about the expression.
CXLine have the following fields:
* FileName - Used for debugging.
* LineNumber - Used for debugging.
* LineStr - Complete string line of the expression.

Note that in CXArgument, its line and filename fields are removed so for error handling, walk program backward to last CXLine to get line number and file name.

## CXExpression
CXOperation have the following fields:
* Index - an int32 which determines its index in the array of CXAtomicOperator, CXArgument, or CXLine.
* Type - Determines which type or which array it is located.
  * The types are:
    * 0-reserved
    * 1-CXAtomicOperator
    * 2-CXArgument
    * 3-CXLine
    * 4-CXStructDef - to be implemented.
    * 5-CXFunctionDef - to be implemented.
    * 6-CXModuleDef - to be implemented.
    * 7-CXGlobalDef/CXModuleGlobalDef - to be implemented.
* ExpressionType
  * The types are:
    * CXEXPR_METHOD_CALL
    * CXEXPR_STRUCT_LITERAL
    * CXEXPR_ARRAY_LITERAL
    * CXEXPR_SCOPE_NEW
    * CXEXPR_SCOPE_DEL

## CXFunction
[CXFunctions](https://github.com/skycoin/cx/blob/develop/cx/cxfunction.go) are the operators in CX, whether they be custom operations, internal library operations, opcodes, or functions defined in CX. The Inputs and Outputs of the CXFunction differ from the CXExpression Inputs and Outputs because the ones the CXFunction have are parameters, not arguments. So CXExpressions store CXArguments that represent the actual input and output data during evaluation, while CXFunctions store inputs and outputs corresponding to: what kind of data types they accept; what kind of named parameters to declare; or otherwise, just type.

They have the following fields:
* Name - name of the function, if named.
* Package - the package this CXFunction resides in.
* IsNative - is this function native to CX? (i32.add, etc)
* OpCode - if IsNative, then this is non-zero and set to the operation it correlates to.
* Inputs, Outputs - Input and output parameters to the CXFunction.
* Operations - all operations in the function.
* Length - number of expressions in the function.
* Size - size, in memory, of the function.
* FileName, FileLine - used for debugging.
* ListOfPointers - used by the garbage collector. These are CXArguments that are the root pointers of the object trees in the heap. 
* CurrentExpression - used by the REPL and parser when checking function validity and processing expressions at compile time.

Functions/Methods (first two are functions, rest are methods): 
* MakeFunction - same as other makers for the previous objects.
* MakeNativeFunction - used to create Opcodes.
* GetExpressions, GetExpressionByLabel, GetExpressionByLine - getters for exactly what they sound like. For GetExpressionByLine, and for other “Line” things, line just means “get the expression in Expressions at index ‘line’”.
* GetCurrentExpression - returns CurrentExpression if not nil, otherwise first element in Expressions if not nil, otherwise an error.
* AddInput, RemoveInput, AddOutput, RemoveOutput - same as CXExpression.
* AddExpression - adds expression to the CXFunction’s Expressions. Also sets the expression’s Package to Package, Function to the caller, and sets CurrentExpression to the input expression, and increments Length.
* RemoveExpression - Not sure when this is used, but does the same thing, but does not unset expression’s Package or Function fields, nor unsets CurrentExpression, nor decrements Length.
* SelectExpression - takes in an integer, and tries to return the element with that index from Expressions. Throws error if Expressions is nil or length zero. If index is out of bounds, it either returns the first or last expression from Expressions, depending on which end the index goes out of bounds at. Also sets CurrentExpression to the grabbed expression.

## CXStruct

[CXStructs](https://github.com/skycoin/cx/blob/develop/cx/cxstruct.go) are the custom struct types in CX, albeit all struct types in CX are currently custom types. They have a Name, a Package they reside in, an array of CXArguments which represent the field names and types, and a Size, which is the precomputed CXStruct’s size in memory. They have a MakeStruct which takes in a name and outputs a blank CXStruct with said name, and also GetFields(), GetField(name), and more:
* AddField - tries to add the CXArgument input to Fields. If the input’s name shares the name of another field, then CX panics. Otherwise, the input’s Offset is set to the previous field’s Offset (or zero if there are no other fields), plus the last field’s * TotalSize (or zero, if nonexistent). Size is incremented by the size of the input.
* RemoveField - same thing as previous, but forgets to update the offset, Size, and other things.

## CXPackage

[CXPackages](https://github.com/skycoin/cx/blob/develop/cx/cxpackage.go) are just like Go packages, and store globals and other things:
* Name - name of the package.
* Imports - CXPackages which this CXPackage imports.
* Functions, Structs, Globals - exactly what you think they are. Globals has CXArguments.
* CurrentFunction, CurrentStruct - same as the case with CXFunction CurrentExpression.

## CXProgram / Runtime

The CX Runtime consists mainly of a [CXProgram](https://github.com/skycoin/cx/blob/develop/cx/cxprogram.go) object, which itself is composed of:
* An array of bytes - the program memory (code, data, stack, heap segments)
* A callstack of CXCalls, consisting of the operator (a CXFunction) and frame pointer
* Several state registers:
  * Stacksize and Heapsize
  * Heapstartsat 
  * Stackpointer and Heappointer
  * Callcounter
* Additionally, a CXProgram stores two arrays of CXArguments - Inputs and Outputs - which represent OS arguments and outputs.
* A Version string and a CurrentPackage, used by the REPL and by the compiler
* The full array of Packages (CXPackage) in the program (essentially the AST/IR)
* Terminated - is the program terminated?
* CurrentPackage - the currently active package in the program. Used by REPL, Compiler.
* AtomicOps  - Array of AtomicOp.
* CxArgs - Array of CXArg.
* CXLines - Array of CXLine.
* CXPackages - Array of CXPackage.
* CXFunctions - Arrays of CXFunction.
