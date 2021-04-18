package ast

// CXArgument is used to define local variables, global variables,
// literals (strings, numbers), inputs and outputs to function
// calls. All of the fields in this structure are determined at
// compile time.
type CXArgument struct {
	// Lengths is used if the `CXArgument` defines an array or a
	// slice. The number of dimensions for the array/slice is
	// equal to `len(Lengths)`, while the contents of `Lengths`
	// define the sizes of each dimension. In the case of a slice,
	// `Lengths` only determines the number of dimensions and the
	// sizes are all equal to 0 (these 0s are not used for any
	// computation).
	Lengths []int
	// DereferenceOperations is a slice of integers where each
	// integer corresponds a `DEREF_*` constant (for example
	// `DEREF_ARRAY`, `DEREF_POINTER`.). A dereference is a
	// process where we consider the bytes at `Offset : Offset +
	// TotalSize` as an address in memory, and we use that address
	// to find the desired value (the referenced
	// value).
	DereferenceOperations []int
	// DeclarationSpecifiers is a slice of integers where each
	// integer corresponds a `DECL_*` constant (for example
	// `DECL_ARRAY`, `DECL_POINTER`.). Declarations are used to
	// create complex types such as `[5][]*Point` (an array of 5
	// slices of pointers to struct instances of type
	// `Point`).
	DeclarationSpecifiers []int
	// Indexes stores what indexes we want to access from the
	// `CXArgument`. A non-nil `Indexes` means that the
	// `CXArgument` is an index or a slice. The elements of
	// `Indexes` can be any `CXArgument` (for example, literals
	// and variables).
	Indexes []*CXArgument
	// Fields stores what fields are being accessed from the
	// `CXArgument` and in what order. Whenever a `DEREF_FIELD` in
	// `DereferenceOperations` is found, we consume a field from
	// `Field` to determine the new offset to the desired
	// value.
	Fields []*CXArgument
	// Inputs defines the input parameters of a first-class
	// function. The `CXArgument` is of type `TYPE_FUNC` if
	// `ProgramInput` is non-nil.
	Inputs []*CXArgument
	// Outputs defines the output parameters of a first-class
	// function. The `CXArgument` is of type `TYPE_FUNC` if
	// `ProgramOutput` is non-nil.
	Outputs []*CXArgument
	// Type defines what's the basic or primitev type of the
	// `CXArgument`. `Type` can be equal to any of the `TYPE_*`
	// constants (e.g. `TYPE_STR`, `TYPE_I32`).
	Type int
	// Size determines the size of the basic type. For example, if
	// the `CXArgument` is of type `TYPE_CUSTOM` (i.e. a
	// user-defined type or struct) and the size of the struct
	// representing the custom type is 10 bytes, then `Size == 10`.
	Size int
	// TotalSize represents how many bytes are referenced by the
	// `CXArgument` in total. For example, if the `CXArgument`
	// defines an array of 5 struct instances of size 10 bytes,
	// then `TotalSize == 50`.
	TotalSize int
	// Offset defines a relative memory offset (used in
	// conjunction with the frame pointer), in the case of local
	// variables, or it could define an absolute memory offset, in
	// the case of global variables and literals. It is used by
	// the CX virtual machine to find the bytes that represent the
	// value of the `CXArgument`.
	Offset int
	// IndirectionLevels
	IndirectionLevels int
	DereferenceLevels int
	PassBy            int // pass by value or reference

	ArgDetails *CXArgumentDebug

	CustomType *CXStruct
	IsSlice    bool
	// IsArray                      bool
	IsPointer                    bool
	IsReference                  bool
	IsStruct                     bool
	IsRest                       bool // pkg.var <- var is rest
	IsLocalDeclaration           bool
	IsShortAssignmentDeclaration bool // variables defined with :=
	IsInnerReference             bool // for example: &slice[0] or &struct.field
	PreviouslyDeclared           bool
	DoesEscape                   bool
}
