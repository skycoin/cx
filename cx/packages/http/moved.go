package http

// TODO, add function to register these op codes for usage

/*
func LoadOpCodeTables() {
	httpPkg, err := ast.PROGRAM.GetPackage("http")
	if err != nil {
		panic(err)
	}

	RegisterOpCode(constants.OP_HTTP_SERVE, "http.Serve", opHTTPServe, In(ast.ConstCxArg_STR), Out(ast.ConstCxArg_STR))
	RegisterOpCode(constants.OP_HTTP_LISTEN_AND_SERVE, "http.ListenAndServe", opHTTPListenAndServe, In(ast.ConstCxArg_STR), Out(ast.ConstCxArg_STR))
	RegisterOpCode(constants.OP_HTTP_NEW_REQUEST, "http.NewRequest", opHTTPNewRequest, In(ast.ConstCxArg_STR, ast.ConstCxArg_STR, ast.ConstCxArg_STR), Out(ast.ConstCxArg_STR))
	RegisterOpCode(constants.OP_HTTP_DO, "http.Do", opHTTPDo, In(ast.ConstCxArg_UND_TYPE), Out(ast.ConstCxArg_UND_TYPE, ast.ConstCxArg_STR))
	//RegisterOpCode(constants.OP_DMSG_DO, "http.DmsgDo", opDMSGDo, In(ast.ConstCxArg_UND_TYPE), Out(ast.ConstCxArg_STR))

	RegisterOpCode(constants.OP_HTTP_HANDLE, "http.Handle", opHTTPHandle,
		In(
			ast.ConstCxArg_STR,
			ParamEx(ParamData{TypCode: constants.TYPE_FUNC, Pkg: httpPkg, inputs: In(ast.MakeArgument("ResponseWriter", "", -1).AddType(constants.TypeNames[constants.TYPE_STR]), ast.Pointer(Struct("http", "Request", "r")))})),
		Out())

	RegisterOpCode(constants.OP_HTTP_CLOSE, "http.Close", opHTTPClose, nil, nil)

}
 */

//opcodes.go

//const (
/*
	OP_HTTP_SERVE
	OP_HTTP_LISTEN_AND_SERVE
	OP_HTTP_NEW_REQUEST
	OP_HTTP_DO
	OP_HTTP_HANDLE
	OP_HTTP_CLOSE
 */


