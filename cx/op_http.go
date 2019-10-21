// +build base

package cxcore

import (
	"net"
	"net/http"
	"strings"

	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
)

func init() {
	// In this case we're adding the `URL` type to the `http` package.
	httpPkg := MakePackage("http")
	urlStrct := MakeStruct("URL")
	
	urlStrct.AddField(MakeArgument("Scheme", "", 0).AddType(TypeNames[TYPE_STR]))
	urlStrct.AddField(MakeArgument("Opaque", "", 0).AddType(TypeNames[TYPE_STR]))
	urlStrct.AddField(MakeArgument("Host", "", 0).AddType(TypeNames[TYPE_STR]))
	urlStrct.AddField(MakeArgument("Path", "", 0).AddType(TypeNames[TYPE_STR]))
	urlStrct.AddField(MakeArgument("RawPath", "", 0).AddType(TypeNames[TYPE_STR]))
	urlStrct.AddField(MakeArgument("ForceQuery", "", 0).AddType(TypeNames[TYPE_BOOL]))
	urlStrct.AddField(MakeArgument("RawQuery", "", 0).AddType(TypeNames[TYPE_STR]))
	urlStrct.AddField(MakeArgument("Fragment", "", 0).AddType(TypeNames[TYPE_STR]))
	
	httpPkg.AddStruct(urlStrct)

	// 95% sure that you will also need Golang's `net/http`'s `Request` struct.
	// If this is the case, you can add that structure like this:

	requestStrct := MakeStruct("Request")
	
	requestStrct.AddField(MakeArgument("Method", "", 0).AddType(TypeNames[TYPE_STR]))
	urlFld := MakeArgument("URL", "", 0).AddType(TypeNames[TYPE_CUSTOM])
	// Golang declares this one as a pointer. It can be done in CX using
	// `CXArgument.DeclarationSpecifiers`.
	urlFld.DeclarationSpecifiers = append(urlFld.DeclarationSpecifiers, DECL_STRUCT)
	urlFld.DeclarationSpecifiers = append(urlFld.DeclarationSpecifiers, DECL_POINTER)
	urlFld.IsPointer = true
	// urlFld.PassBy = PASSBY_REFERENCE
	urlFld.CustomType = urlStrct
	requestStrct.AddField(urlFld)
	// If `Request` is indeed needed, add the other fields. If you *need* `Header`,
	// you'll be in trouble (maybe), as `Header` is of type `map[string][]string`
	// and CX doesn't have maps. If you need this, implement it as an array (or a slice), where
	// the first element is X datum, the second element is Y datum, etc. For example:
	headerFld := MakeArgument("Header", "", 0).AddType(TypeNames[TYPE_STR]) // will be a slice of strings
	headerFld.DeclarationSpecifiers = append(headerFld.DeclarationSpecifiers, DECL_SLICE)
	headerFld.IsSlice = true
	headerFld.IsReference = true
	headerFld.IsArray = true
	headerFld.PassBy = PASSBY_REFERENCE
	headerFld.Lengths = []int{0} // 1D slice. If it was a 2D slice it'd be []int{0, 0}

	// Then we add the field
	requestStrct.AddField(headerFld)

	// And adding the `Request` structure to the `http` package.
	httpPkg.AddStruct(requestStrct)

	// Sorry, there ARE functions that handle all of these operations, but they're part of the
	// parser. These files are part of the `cxcore` package, where, until now, we didn't have
	// the need to add structs like this. So for now we have to do it manually.
	// For reference, check `cxgo/actions/declarations.go`, in particular `DeclarationSpecifiers()`.
	
	PROGRAM.AddPackage(httpPkg)
}

func opHTTPServe(prgrm *CXProgram) {
	expr := prgrm.GetExpr()

	fp := prgrm.GetFramePointer()
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	url := ReadStr(fp, inp1)

	l, err := net.Listen("tcp", url)
	if err != nil {
		writeString(fp, err.Error(), out1)
	}

	err = http.Serve(l, nil)
	if err != nil {
		writeString(fp, err.Error(), out1)
	}
}

func opHTTPNewRequest(prgrm *CXProgram) {
	// Use or remove as needed. In the following lines I
	// get the "http" package, then I get the "URL" struct
	// from it. If you wanted to return an URL, for example,
	// you'd use the offset provided by the frame pointer `fp`
	// (which is obtained below) + the offset of `out1` + the offset
	// of the field that you want to access.
	// In the case of a URL input, you'd do the same, actually, but `inp1`.
	
	// httpPkg, err := PROGRAM.GetPackage("http")
	// if err != nil {
	// 	panic(err)
	// }
	// urlType, err := httpPkg.GetStruct("URL")
	// if err != nil {
	// 	panic(err)
	// }

	// schemeFld, _ := urlType.GetField("Scheme")
	// opaqueFld, _ := urlType.GetField("Opaque")
	// hostFld, _ := urlType.GetField("Host")
	// Debug("SchemeOffset", schemeFld.Offset) // prints 0
	// Debug("OpaqueOffset", opaqueFld.Offset) // prints 4
	// Debug("HostOffset", hostFld.Offset) // prints 8
	///////

	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, inp3, out1 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Outputs[0]

	method := ReadStr(fp, inp1)
	url := ReadStr(fp, inp2)
	body := ReadStr(fp, inp3)

	req, err := http.NewRequest(method, url, strings.NewReader(body))

	if err != nil {
		writeString(fp, err.Error(), out1)
	}

	out1Offset := GetFinalOffset(fp, out1)
	byts := encoder.Serialize(req)
	WriteObject(out1Offset, byts)
}
