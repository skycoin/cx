// +build regexp

package regexp

import (
	"regexp"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"

	"github.com/jinzhu/copier"
)

var regexps map[string]*regexp.Regexp = make(map[string]*regexp.Regexp, 0)

// regexpCompile is a helper function for `opRegexpMustCompile` and
// `opRegexpCompile`. `regexpCompile` compiles a `regexp.Regexp` structure
// and adds it to global `regexps`. It also writes CX structure `regexp.Regexp`.
func regexpCompile(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) error {
	// Extracting regular expression to work with, contained in `inp0`.
	exp := inputs[0].Get_str(prgrm)

	// Output structure `Regexp`.
	reg := ast.CXArgument{}
	err := copier.Copy(&reg, outputs[0].TypeSignature)
	if err != nil {
		panic(err)
	}

	// Extracting CX `regexp` package.
	regexpPkg, err := prgrm.GetPackage("regexp")
	if err != nil {
		panic(err)
	}

	// Extracting `regexp`'s Regexp structure.
	regexpType, err := regexpPkg.GetStruct(prgrm, "Regexp")
	if err != nil {
		panic(err)
	}

	// Extracting `regexp.Regexp`'s `exp` field.
	expFld, err := regexpType.GetField(prgrm, "exp")
	if err != nil {
		panic(err)
	}

	// Writing the regex provided by the user to `reg`.
	// This allows us to know what `Regexp` instance the user wants to use
	// in other parts of CX code.
	// TODO: I don't know what would happen if the user uses the same regex
	// in two parts of a CX program. They'll be using the same instance
	// internally.
	expFldIdx := prgrm.AddCXArgInArray(expFld)
	accessExp := []ast.CXArgumentIndex{expFldIdx}
	reg.Fields = accessExp
	types.Write_str(prgrm, prgrm.Memory, ast.GetFinalOffset(prgrm, outputs[0].FramePointer, &reg, nil), exp)
	// Storing `Regexp` instance.
	regexps[exp], err = regexp.Compile(exp)

	return err
}

// opRegexpMustCompile is a wrapper for golang's `regexp`'s `MustCompile`.
func opRegexpMustCompile(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	err := regexpCompile(prgrm, inputs, outputs)

	if err != nil {
		println(err.Error())
		panic(constants.CX_RUNTIME_ERROR)
	}

}

// opRegexpCompile is a wrapper for golang's `regexp`'s `MustCompile`.
func opRegexpCompile(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	// We're only interested in `out1`, which represents the
	// returned error.
	err := regexpCompile(prgrm, inputs, outputs)

	// Writing error message to `out1`.
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	outputs[1].Set_str(prgrm, errStr)
}

// opRegexpCompile is a wrapper for golang's `regexp`'s `MustCompile`.
func opRegexpFind(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	// Output structure `Regexp`.
	reg := ast.CXArgument{}
	err := copier.Copy(&reg, inputs[0].TypeSignature)
	if err != nil {
		panic(err)
	}

	// Extracting CX `regexp` package.
	regexpPkg, err := prgrm.GetPackage("regexp")
	if err != nil {
		panic(err)
	}

	// Extracting `regexp`'s Regexp structure.
	regexpType, err := regexpPkg.GetStruct(prgrm, "Regexp")
	if err != nil {
		panic(err)
	}

	// Extracting `regexp.Regexp`'s `exp` field.
	expFld, err := regexpType.GetField(prgrm, "exp")
	if err != nil {
		panic(err)
	}

	// Getting corresponding `Regexp` instance.
	expFldIdx := prgrm.AddCXArgInArray(expFld)
	accessExp := []ast.CXArgumentIndex{expFldIdx}
	reg.Fields = accessExp
	exp := types.Read_str(prgrm.Memory, ast.GetFinalOffset(prgrm, inputs[0].FramePointer, &reg, nil))
	r := regexps[exp]

	outputs[0].Set_str(prgrm, string(r.Find([]byte(inputs[1].Get_str(prgrm)))))
}
