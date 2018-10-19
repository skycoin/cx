package base

/*
  Root Program
*/

type CXProgram struct {
        Packages                        []*CXPackage
        CurrentPackage                  *CXPackage

        Inputs                          []*CXArgument
        Outputs                         []*CXArgument

        CallStack                       []CXCall
        CallCounter                     int

        Memory                          []byte
        HeapPointer                     int
        StackPointer                    int

        HeapStartsAt                    int

        Terminated                      bool

        Path                            string
}

type CXCall struct {
        Operator                        *CXFunction
        Line                            int
        FramePointer                    int

        Package                         *CXPackage
        Program                         *CXProgram
}

/*
  Packages
*/

type CXPackage struct {
        Name                            string
        Imports                         []*CXPackage
        Functions                       []*CXFunction
        Structs                         []*CXStruct
        Globals                         []*CXArgument

        CurrentFunction                 *CXFunction
        CurrentStruct                   *CXStruct
        Program                         *CXProgram
}

/*
  Structs
*/

type CXStruct struct {
        Name                            string
        Fields                          []*CXArgument
        Size                            int

        Package                         *CXPackage
        Program                         *CXProgram
}

/*
  Functions
*/

type CXFunction struct {
        Name                            string
        Inputs                          []*CXArgument
        Outputs                         []*CXArgument
        Expressions                     []*CXExpression
        Size                            int // automatic memory size
        Length                          int // number of expressions, pre-computed for performance

        ListOfPointers                  []*CXArgument

        IsNative                        bool
        OpCode                          int

        CurrentExpression               *CXExpression
        Package                         *CXPackage
        Program                         *CXProgram
}

type CXExpression struct {
        Operator                        *CXFunction
        Inputs                          []*CXArgument
        Outputs                         []*CXArgument
        // debugging
        Line                            int
        FileLine                        int
        FileName                        string

        // used for jmp statements
        Label                           string
        ThenLines                       int
        ElseLines                       int

        IsMethodCall                    bool
        IsStructLiteral                 bool
        IsArrayLiteral                  bool

        Function                        *CXFunction
        Package                         *CXPackage
        Program                         *CXProgram
}

type CXConstant struct {
        // native constants. only used for pre-packaged constants (e.g. math package's PI)
        // these fields are used to feed WritePrimary
        Type                            int
        Value                           []byte
}

type CXArgument struct {
        Name                            string
        Type                            int
        CustomType                      *CXStruct
        Size                            int // size of underlaying basic type
        TotalSize                       int // total size of an array, performance reasons

        Offset                          int

        IndirectionLevels               int
        DereferenceLevels               int
        DereferenceOperations           []int // offset by array index, struct field, pointer
        DeclarationSpecifiers           []int // used to determine finalSize

        IsSlice                         bool
        IsArray                         bool
        IsArrayFirst                    bool // and then dereference
        IsPointer                       bool
        IsReference                     bool

        IsDereferenceFirst              bool // and then array
        IsStruct                        bool
        IsRest                          bool // pkg.var <- var is rest
        IsLocalDeclaration              bool
        IsShortDeclaration              bool

        PassBy                          int  // pass by value or reference
        DoesEscape                      bool

        Lengths                         []int // declared lengths at compile time
        Indexes                         []*CXArgument
        Fields                          []*CXArgument // strct.fld1.fld2().fld3

        FileLine                        int
        FileName                        string
        
        Package                         *CXPackage
        Program                         *CXProgram
}
