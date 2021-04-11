// +build regexp

package regexp

//TODO: Get method of assigning the enums dynamically

//TODO: Add init function to register the call back functions

/*
	// regexp
	OP_REGEXP_COMPILE
	OP_REGEXP_MUST_COMPILE
	OP_REGEXP_FIND
*/


/*
	// regexp
	RegisterOpCode(OP_REGEXP_COMPILE, "regexp.Compile", opRegexpCompile, In(ast.ConstCxArg_STR), Out(Struct("regexp", "Regexp", "r"), ast.ConstCxArg_STR))
	RegisterOpCode(OP_REGEXP_MUST_COMPILE, "regexp.MustCompile", opRegexpMustCompile, In(ast.ConstCxArg_STR), Out(Struct("regexp", "Regexp", "r")))
	RegisterOpCode(OP_REGEXP_FIND, "regexp.Regexp.Find", opRegexpFind, In(Struct("regexp", "Regexp", "r"), ast.ConstCxArg_STR), Out(ast.ConstCxArg_STR))

func init() {

 */