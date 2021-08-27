// +build regexp

package regexp

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
	"regexp"

	"github.com/jinzhu/copier"
)

var regexps map[string]*regexp.Regexp = make(map[string]*regexp.Regexp, 0)

// regexpCompile is a helper function for `opRegexpMustCompile` and
// `opRegexpCompile`. `regexpCompile` compiles a `regexp.Regexp` structure
// and adds it to global `regexps`. It also writes CX structure `regexp.Regexp`.
func regexpCompile(inputs []ast.CXValue, outputs []ast.CXValue) error {
	// Extracting regular expression to work with, contained in `inp0`.
	exp := inputs[0].Get_str()

	// Output structure `Regexp`.
	reg := ast.CXArgument{}
	err := copier.Copy(&reg, outputs[0].Arg)
	if err != nil {
		panic(err)
	}

	// Extracting CX `regexp` package.
	regexpPkg, err := ast.PROGRAM.GetPackage("regexp")
	if err != nil {
		panic(err)
	}

	// Extracting `regexp`'s Regexp structure.
	regexpType, err := regexpPkg.GetStruct("Regexp")
	if err != nil {
		panic(err)
	}

	// Extracting `regexp.Regexp`'s `exp` field.
	expFld, err := regexpType.GetField("exp")
	if err != nil {
		panic(err)
	}

	// Writing the regex provided by the user to `reg`.
	// This allows us to know what `Regexp` instance the user wants to use
	// in other parts of CX code.
	// TODO: I don't know what would happen if the user uses the same regex
	// in two parts of a CX program. They'll be using the same instance
	// internally.
	accessExp := []*ast.CXArgument{expFld}
	reg.Fields = accessExp
	types.Write_str(ast.PROGRAM.Memory, ast.GetFinalOffset(outputs[0].FramePointer, &reg), exp)
	// Storing `Regexp` instance.
	regexps[exp], err = regexp.Compile(exp)

	return err
}

// opRegexpMustCompile is a wrapper for golang's `regexp`'s `MustCompile`.
func opRegexpMustCompile(inputs []ast.CXValue, outputs []ast.CXValue) {
	err := regexpCompile(inputs, outputs)

	if err != nil {
		println(err.Error())
		panic(constants.CX_RUNTIME_ERROR)
	}

}

// opRegexpCompile is a wrapper for golang's `regexp`'s `MustCompile`.
func opRegexpCompile(inputs []ast.CXValue, outputs []ast.CXValue) {
	// We're only interested in `out1`, which represents the
	// returned error.
	err := regexpCompile(inputs, outputs)

	// Writing error message to `out1`.
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	outputs[1].Set_str(errStr)
}

// opRegexpCompile is a wrapper for golang's `regexp`'s `MustCompile`.
func opRegexpFind(inputs []ast.CXValue, outputs []ast.CXValue) {
	// Output structure `Regexp`.
	reg := ast.CXArgument{}
	err := copier.Copy(&reg, inputs[0].Arg)
	if err != nil {
		panic(err)
	}

	// Extracting CX `regexp` package.
	regexpPkg, err := ast.PROGRAM.GetPackage("regexp")
	if err != nil {
		panic(err)
	}

	// Extracting `regexp`'s Regexp structure.
	regexpType, err := regexpPkg.GetStruct("Regexp")
	if err != nil {
		panic(err)
	}

	// Extracting `regexp.Regexp`'s `exp` field.
	expFld, err := regexpType.GetField("exp")
	if err != nil {
		panic(err)
	}

	// Getting corresponding `Regexp` instance.
	accessExp := []*ast.CXArgument{expFld}
	reg.Fields = accessExp
	exp := types.Read_str(ast.PROGRAM.Memory, ast.GetFinalOffset(inputs[0].FramePointer, &reg))
	r := regexps[exp]

	outputs[0].Set_str(string(r.Find([]byte(inputs[1].Get_str()))))
}
