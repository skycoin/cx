// +build http

package http

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/opcodes"
	"github.com/skycoin/cx/cx/types"
)

func RegisterPackage() {
	httpPkg := ast.MakePackage("http")
	urlStrct := ast.MakeStruct("URL")

	urlStrct.AddField(ast.MakeArgument("Scheme", "", 0).AddType(types.STR).AddPackage(httpPkg))
	urlStrct.AddField(ast.MakeArgument("Opaque", "", 0).AddType(types.STR).AddPackage(httpPkg))
	urlStrct.AddField(ast.MakeArgument("Host", "", 0).AddType(types.STR).AddPackage(httpPkg))
	urlStrct.AddField(ast.MakeArgument("Path", "", 0).AddType(types.STR).AddPackage(httpPkg))
	urlStrct.AddField(ast.MakeArgument("RawPath", "", 0).AddType(types.STR).AddPackage(httpPkg))
	urlStrct.AddField(ast.MakeArgument("ForceQuery", "", 0).AddType(types.BOOL).AddPackage(httpPkg))
	urlStrct.AddField(ast.MakeArgument("RawQuery", "", 0).AddType(types.STR).AddPackage(httpPkg))
	urlStrct.AddField(ast.MakeArgument("Fragment", "", 0).AddType(types.STR).AddPackage(httpPkg))

	httpPkg.AddStruct(urlStrct)

	requestStrct := ast.MakeStruct("Request")

	requestStrct.AddField(ast.MakeArgument("Method", "", 0).AddType(types.STR).AddPackage(httpPkg))
	urlFld := ast.MakeArgument("URL", "", 0).AddType(types.CUSTOM).AddPackage(httpPkg)
	urlFld.DeclarationSpecifiers = append(urlFld.DeclarationSpecifiers, constants.DECL_STRUCT)
	urlFld.DeclarationSpecifiers = append(urlFld.DeclarationSpecifiers, constants.DECL_POINTER)
	urlFld.IsPointer = true
	urlFld.Size = types.POINTER_SIZE
	urlFld.TotalSize = types.POINTER_SIZE
	urlFld.CustomType = urlStrct
	requestStrct.AddField(urlFld)

	headerFld := ast.MakeArgument("Header", "", 0).AddType(types.STR).AddPackage(httpPkg) // will be a slice of strings
	headerFld.DeclarationSpecifiers = append(headerFld.DeclarationSpecifiers, constants.DECL_SLICE)
	headerFld.DeclarationSpecifiers = append(headerFld.DeclarationSpecifiers, constants.DECL_SLICE)
	headerFld.IsSlice = true
	headerFld.IsReference = true
	// headerFld.IsArray = true
	headerFld.PassBy = constants.PASSBY_REFERENCE
	headerFld.Lengths = []types.Pointer{0, 0}

	requestStrct.AddField(headerFld)

	requestStrct.AddField(ast.MakeArgument("Body", "", 0).AddType(types.STR).AddPackage(httpPkg))

	httpPkg.AddStruct(requestStrct)

	// Mapping http.Response struct
	responseStruct := ast.MakeStruct("Response")
	responseStruct.AddField(ast.MakeArgument("Status", "", 0).AddType(types.STR).AddPackage(httpPkg))
	responseStruct.AddField(ast.MakeArgument("StatusCode", "", 0).AddType(types.I32).AddPackage(httpPkg))
	responseStruct.AddField(ast.MakeArgument("Proto", "", 0).AddType(types.STR).AddPackage(httpPkg))
	responseStruct.AddField(ast.MakeArgument("ProtoMajor", "", 0).AddType(types.I32).AddPackage(httpPkg))
	responseStruct.AddField(ast.MakeArgument("ProtoMinor", "", 0).AddType(types.I32).AddPackage(httpPkg))
	responseStruct.AddField(ast.MakeArgument("Body", "", 0).AddType(types.STR).AddPackage(httpPkg))
	//TODO Header Header - not sure if headerFld used for http.Request can be used here
	//TODO Body io.ReadCloser
	responseStruct.AddField(ast.MakeArgument("ContentLength", "", 0).AddType(types.I64).AddPackage(httpPkg))
	transferEncodingFld := ast.MakeArgument("TransferEncoding", "", 0).AddType(types.STR).AddPackage(httpPkg)
	transferEncodingFld.DeclarationSpecifiers = append(transferEncodingFld.DeclarationSpecifiers, constants.DECL_SLICE)
	transferEncodingFld.IsSlice = true
	transferEncodingFld.IsReference = true
	// transferEncodingFld.IsArray = true
	transferEncodingFld.PassBy = constants.PASSBY_REFERENCE
	transferEncodingFld.Lengths = []types.Pointer{0}
	responseStruct.AddField(transferEncodingFld)
	urlStrct.AddField(ast.MakeArgument("Close", "", 0).AddType(types.BOOL).AddPackage(httpPkg))
	urlStrct.AddField(ast.MakeArgument("Uncompressed", "", 0).AddType(types.BOOL).AddPackage(httpPkg))
	//TODO Trailer Header
	//TODO Request *Request
	//TODO TLS *tls.ConnectionState

	httpPkg.AddStruct(responseStruct)

	ast.PROGRAM.AddPackage(httpPkg)

	opcodes.RegisterFunction("http.Serve", opHTTPServe, opcodes.In(ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_STR))
	opcodes.RegisterFunction("http.ListenAndServe", opHTTPListenAndServe, opcodes.In(ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_STR))
	opcodes.RegisterFunction("http.NewRequest", opHTTPNewRequest, opcodes.In(ast.ConstCxArg_STR, ast.ConstCxArg_STR, ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_STR))
	opcodes.RegisterFunction("http.Do", opHTTPDo, opcodes.In(ast.ConstCxArg_UND_TYPE), opcodes.Out(ast.ConstCxArg_UND_TYPE, ast.ConstCxArg_STR))
	//opcodes.RegisterFunction("http.DmsgDo", opDMSGDo, opcodes.In(ast.ConstCxArg_UND_TYPE), opcodes.Out(ast.ConstCxArg_STR))
	opcodes.RegisterFunction("http.Handle", opHTTPHandle,
		opcodes.In(
			ast.ConstCxArg_STR,
			ast.Func(httpPkg, opcodes.In(ast.MakeArgument("ResponseWriter", "", -1).AddType(types.STR), ast.Pointer(ast.Struct("http", "Request", "r"))), nil)),
		nil)
	opcodes.RegisterFunction("http.Close", opHTTPClose, nil, nil)

}
