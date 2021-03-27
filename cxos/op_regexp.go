// +build os

package cxos

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"regexp"

	"github.com/jinzhu/copier"
)

var regexps map[string]*regexp.Regexp = make(map[string]*regexp.Regexp, 0)

func init() {
	regexpPkg := ast.MakePackage("regexp")
	regexpStrct := ast.MakeStruct("Regexp")

	regexpStrct.AddField(ast.MakeArgument("exp", "", 0).AddType(constants.TypeNames[constants.TYPE_STR]).AddPackage(regexpPkg))

	regexpPkg.AddStruct(regexpStrct)

	ast.PROGRAM.AddPackage(regexpPkg)
}

// regexpCompile is a helper function for `opRegexpMustCompile` and
// `opRegexpCompile`. `regexpCompile` compiles a `regexp.Regexp` structure
// and adds it to global `regexps`. It also writes CX structure `regexp.Regexp`.
func regexpCompile(expr *ast.CXExpression, fp int) error {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]

	// Extracting regular expression to work with, contained in `inp1`.
	exp := ast.ReadStr(fp, inp1)

	// Output structure `Regexp`.
	reg := ast.CXArgument{}
	err := copier.Copy(&reg, out1)
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
	ast.WriteString(fp, exp, &reg)

	// Storing `Regexp` instance.
	regexps[exp], err = regexp.Compile(exp)

	return err
}

// opRegexpMustCompile is a wrapper for golang's `regexp`'s `MustCompile`.
func opRegexpMustCompile(expr *ast.CXExpression, fp int) {
	err := regexpCompile(expr, fp)

	if err != nil {
		println(err.Error())
		panic(constants.CX_RUNTIME_ERROR)
	}

}

// opRegexpCompile is a wrapper for golang's `regexp`'s `MustCompile`.
func opRegexpCompile(expr *ast.CXExpression, fp int) {
	// We're only interested in `out2`, which represents the
	// returned error.
	out2 := expr.Outputs[1]

	err := regexpCompile(expr, fp)

	// Writing error message to `out2`.
	if err != nil {
		ast.WriteString(fp, err.Error(), out2)
	}
}

// opRegexpCompile is a wrapper for golang's `regexp`'s `MustCompile`.
func opRegexpFind(expr *ast.CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	// Output structure `Regexp`.
	reg := ast.CXArgument{}
	err := copier.Copy(&reg, inp1)
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
	exp := ast.ReadStr(fp, &reg)
	r := regexps[exp]

	ast.WriteString(fp, string(r.Find([]byte(ast.ReadStr(fp, inp2)))), out1)
}
