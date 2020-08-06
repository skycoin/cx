// +build base

package cxcore

import (
	"regexp"

	"github.com/jinzhu/copier"

	cxcore "github.com/SkycoinProject/cx/cx"
)

var regexps map[string]*regexp.Regexp = make(map[string]*regexp.Regexp, 0)

func init() {
	regexpPkg := cxcore.MakePackage("regexp")
	regexpStrct := cxcore.MakeStruct("Regexp")

	regexpStrct.AddField(cxcore.MakeArgument("exp", "", 0).AddType(cxcore.TypeNames[cxcore.TYPE_STR]).AddPackage(regexpPkg))

	regexpPkg.AddStruct(regexpStrct)

	cxcore.PROGRAM.AddPackage(regexpPkg)
}

// regexpCompile is a helper function for `opRegexpMustCompile` and
// `opRegexpCompile`. `regexpCompile` compiles a `regexp.Regexp` structure
// and adds it to global `regexps`. It also writes CX structure `regexp.Regexp`.
func regexpCompile(prgrm *cxcore.CXProgram) error {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]

	// Extracting regular expression to work with, contained in `inp1`.
	exp := cxcore.ReadStr(fp, inp1)

	// Output structure `Regexp`.
	reg := cxcore.CXArgument{}
	err := copier.Copy(&reg, out1)
	if err != nil {
		panic(err)
	}

	// Extracting CX `regexp` package.
	regexpPkg, err := cxcore.PROGRAM.GetPackage("regexp")
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
	accessExp := []*cxcore.CXArgument{expFld}
	reg.Fields = accessExp
	cxcore.WriteString(fp, exp, &reg)

	// Storing `Regexp` instance.
	regexps[exp], err = regexp.Compile(exp)

	return err
}

// opRegexpMustCompile is a wrapper for golang's `regexp`'s `MustCompile`.
func opRegexpMustCompile(prgrm *cxcore.CXProgram) {
	err := regexpCompile(prgrm)

	if err != nil {
		println(err.Error())
		panic(cxcore.CX_RUNTIME_ERROR)
	}

}

// opRegexpCompile is a wrapper for golang's `regexp`'s `MustCompile`.
func opRegexpCompile(prgrm *cxcore.CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	// We're only interested in `out2`, which represents the
	// returned error.
	out2 := expr.Outputs[1]

	err := regexpCompile(prgrm)

	// Writing error message to `out2`.
	if err != nil {
		cxcore.WriteString(fp, err.Error(), out2)
	}
}

// opRegexpCompile is a wrapper for golang's `regexp`'s `MustCompile`.
func opRegexpFind(prgrm *cxcore.CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	// Output structure `Regexp`.
	reg := cxcore.CXArgument{}
	err := copier.Copy(&reg, inp1)
	if err != nil {
		panic(err)
	}

	// Extracting CX `regexp` package.
	regexpPkg, err := cxcore.PROGRAM.GetPackage("regexp")
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
	accessExp := []*cxcore.CXArgument{expFld}
	reg.Fields = accessExp
	exp := cxcore.ReadStr(fp, &reg)
	r := regexps[exp]

	cxcore.WriteString(fp, string(r.Find([]byte(cxcore.ReadStr(fp, inp2)))), out1)
}
