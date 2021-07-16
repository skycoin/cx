package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/execute"
	"github.com/skycoin/cx/cx/types"

	"github.com/skycoin/skycoin/src/cipher/encoder"

	"github.com/jinzhu/copier"
)

func opHTTPHandle(inputs []ast.CXValue, outputs []ast.CXValue) {

	//step 3  : specify the input and outout parameters of Handle function.
	urlstring, functionnamestring := inputs[0].Arg, inputs[1].Arg
	fp := inputs[0].FramePointer

	// Getting handler function.
	handlerPkg, err := ast.PROGRAM.GetPackage(functionnamestring.ArgDetails.Package.Name)

	if err != nil {
		panic(err)
	}
	handlerFn, err := handlerPkg.GetFunction(functionnamestring.ArgDetails.Name)
	if err != nil {
		panic(err)
	}

	http.HandleFunc(types.Read_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, urlstring)), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		callFP := fp + ast.PROGRAM.CallStack[ast.PROGRAM.CallCounter].Operator.Size

		ast.PROGRAM.CallCounter++
		// PROGRAM.StackPointer += handlerFn.Size
		ast.PROGRAM.CallStack[ast.PROGRAM.CallCounter].Operator = handlerFn
		ast.PROGRAM.CallStack[ast.PROGRAM.CallCounter].Line = 0
		ast.PROGRAM.CallStack[ast.PROGRAM.CallCounter].FramePointer = ast.PROGRAM.StackPointer
		writeHTTPRequest(callFP, handlerFn.Inputs[1], r)
		// PROGRAM.StackPointer -= handlerFn.Size
		ast.PROGRAM.CallCounter--

		i1Off := callFP + handlerFn.Inputs[0].Offset
		i1Size := handlerFn.Inputs[0].TotalSize
		i2Off := callFP + handlerFn.Inputs[1].Offset
		i2Size := handlerFn.Inputs[1].TotalSize

		i1 := make([]byte, i1Size)
		i2 := make([]byte, i1Size)

		copy(i1, ast.PROGRAM.Memory[i1Off:i1Off+i1Size])
		copy(i2, ast.PROGRAM.Memory[i2Off:i2Off+i2Size])

		//PROGRAM.Callback(handlerFn, [][]byte{i1, i2})
		execute.Callback(ast.PROGRAM, handlerFn, [][]byte{i1, i2})
		fmt.Fprint(w, types.Read_str(ast.PROGRAM.Memory, ast.GetFinalOffset(callFP, handlerFn.Inputs[0])))
	})
}

var server *http.Server

func opHTTPClose(inputs []ast.CXValue, outputs []ast.CXValue) {
	server.Close()
}

func opHTTPListenAndServe(inputs []ast.CXValue, outputs []ast.CXValue) {

	//step 3  : specify the input and outout parameters of HTTPServe function.

	url := inputs[0].Get_str()

	server = &http.Server{Addr: url}

	err := server.ListenAndServe()
	outputs[0].Set_str(err.Error())
}

func opHTTPServe(inputs []ast.CXValue, outputs []ast.CXValue) {

	//step 3  : specify the imput and out of HTTPServe function.

	url := inputs[0].Get_str()

	l, err := net.Listen("tcp", url)
	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	err = http.Serve(l, nil)
	if err != nil {
		errStr = err.Error()
	}

	outputs[0].Set_str(errStr)
}

func opHTTPNewRequest(inputs []ast.CXValue, outputs []ast.CXValue) {
	// TODO: This whole OP needs rewriting/finishing.
	// Seems more a prototype.
	stringmethod, stringurl, stringbody, errorstring := inputs[0].Arg, inputs[1].Arg, inputs[2].Arg, outputs[0].Arg

	fp := inputs[0].FramePointer

	//this is an alternative for following 3 lines of code that fail due to URL
	method := types.Read_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, stringmethod))
	urlString := types.Read_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, stringurl))
	body := types.Read_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, stringbody))

	//above is an alternative for following 3 lines of code that fail due to URL
	req, err := http.NewRequest(method, urlString, bytes.NewBuffer([]byte(body)))
	if err != nil {
		types.Write_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, errorstring), err.Error())
	}

	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := netClient.Do(req)
	if err != nil {
		types.Write_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, errorstring), err.Error())
	}
	resp1 := *resp // dereference to exclude pointer issue

	// TODO issue with returning response,
	// the line where resp is serialized (byts := encoder.Serialize(resp)) throws following error, adding response object content for context:
	// &{404 Not Found 404 HTTP/1.1 1 1 map[Content-Length:[19] Content-Type:[text/plain; charset=utf-8] Date:[Sun, 27 Oct 2019 22:10:14 GMT] X-Content-Type-Options:[nosniff]] 0xc0000c8700 19 [] false false map[] 0xc000175400 <nil>}
	// 2019/10/27 23:10:14 invalid type int
	// error: examples/http-serve-and-request-mine.cx:8, CX_RUNTIME_ERROR, invalid type int

	out1Offset := ast.GetFinalOffset(fp, errorstring)

	// TODO: Used `Response.Status` for now, to avoid getting an error.
	// This will be rewritten as the whole operator is unfinished.
	byts := encoder.Serialize(resp1.Status)
	types.Write_obj(ast.PROGRAM.Memory, out1Offset, byts)
}

func writeHTTPRequest(fp types.Pointer, param *ast.CXArgument, request *http.Request) {
	req := ast.CXArgument{}
	err := copier.Copy(&req, param)
	if err != nil {
		panic(err)
	}

	httpPkg, err := ast.PROGRAM.GetPackage("http")
	if err != nil {
		panic(err)
	}

	urlType, err := httpPkg.GetStruct("URL")
	if err != nil {
		panic(err)
	}
	requestType, err := httpPkg.GetStruct("Request")
	if err != nil {
		panic(err)
	}

	methodFld, err := requestType.GetField("Method")
	if err != nil {
		panic(err)
	}

	bodyFld, err := requestType.GetField("Body")
	if err != nil {
		panic(err)
	}

	urlFld, err := requestType.GetField("URL")
	if err != nil {
		panic(err)
	}

	derefURLFld := ast.CXArgument{}
	err = copier.Copy(&derefURLFld, urlFld)
	if err != nil {
		panic(err)
	}

	derefURLFld.DereferenceOperations = append(derefURLFld.DereferenceOperations, constants.DEREF_POINTER)

	schemeFld, err := urlType.GetField("Scheme")
	if err != nil {
		panic(err)
	}
	hostFld, err := urlType.GetField("Host")
	if err != nil {
		panic(err)
	}
	pathFld, err := urlType.GetField("Path")
	if err != nil {
		panic(err)
	}
	rawPathFld, err := urlType.GetField("RawPath")
	if err != nil {
		panic(err)
	}
	forceQueryFld, err := urlType.GetField("ForceQuery")
	if err != nil {
		panic(err)
	}

	// Creating empty `http.Request` object on heap.
	types.AllocWrite_obj_data(ast.PROGRAM.Memory, make([]byte, requestType.Size))

	req.DereferenceOperations = append(req.DereferenceOperations, constants.DEREF_POINTER)

	// Creating empty `http.URL` object on heap.
	req.Fields = []*ast.CXArgument{urlFld}
	types.AllocWrite_obj_data(ast.PROGRAM.Memory, make([]byte, urlType.Size))

	req.Fields = []*ast.CXArgument{methodFld}
	types.Write_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &req), request.Method)

	req.Fields = []*ast.CXArgument{bodyFld}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	types.Write_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &req), string(body))
	req.Fields = []*ast.CXArgument{&derefURLFld, schemeFld}

	types.Write_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &req), request.URL.Scheme)
	req.Fields = []*ast.CXArgument{&derefURLFld, hostFld}

	types.Write_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &req), request.URL.Host)
	req.Fields = []*ast.CXArgument{&derefURLFld, pathFld}

	types.Write_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &req), request.URL.Path)
	req.Fields = []*ast.CXArgument{&derefURLFld, rawPathFld}

	types.Write_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &req), request.URL.RawPath)
	req.Fields = []*ast.CXArgument{&derefURLFld, forceQueryFld}

	types.Write_bool(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &req),request.URL.ForceQuery)
}

func opHTTPDo(inputs []ast.CXValue, outputs []ast.CXValue) {

	reqstruct, respstruct, errorstring := inputs[0].Arg, outputs[0].Arg, outputs[1].Arg
	fp := inputs[0].FramePointer

	//TODO read req from the inputs
	// reqByts := ReadMemory(GetFinalOffset(fp, inp1), inp1)

	req := ast.CXArgument{}
	err := copier.Copy(&req, reqstruct)
	if err != nil {
		panic(err)
	}

	httpPkg, err := ast.PROGRAM.GetPackage("http")
	if err != nil {
		panic(err)
	}
	urlType, err := httpPkg.GetStruct("URL")
	if err != nil {
		panic(err)
	}
	requestType, err := httpPkg.GetStruct("Request")
	if err != nil {
		panic(err)
	}

	methodFld, err := requestType.GetField("Method")
	if err != nil {
		panic(err)
	}

	urlFld, err := requestType.GetField("URL")
	if err != nil {
		panic(err)
	}

	derefURLFld := ast.CXArgument{}
	err = copier.Copy(&derefURLFld, urlFld)
	if err != nil {
		panic(err)
	}

	derefURLFld.DereferenceOperations = append(derefURLFld.DereferenceOperations, constants.DEREF_POINTER)

	schemeFld, err := urlType.GetField("Scheme")
	if err != nil {
		panic(err)
	}
	hostFld, err := urlType.GetField("Host")
	if err != nil {
		panic(err)
	}
	pathFld, err := urlType.GetField("Path")
	if err != nil {
		panic(err)
	}
	rawPathFld, err := urlType.GetField("RawPath")
	if err != nil {
		panic(err)
	}
	forceQueryFld, err := urlType.GetField("ForceQuery")
	if err != nil {
		panic(err)
	}

	request := http.Request{}
	url := url.URL{}
	request.URL = &url

	req.Fields = []*ast.CXArgument{methodFld}
	request.Method = types.Read_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &req))

	req.Fields = []*ast.CXArgument{&derefURLFld, schemeFld}
	url.Scheme = types.Read_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &req))

	req.Fields = []*ast.CXArgument{&derefURLFld, hostFld}
	url.Host = types.Read_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &req))

	req.Fields = []*ast.CXArgument{&derefURLFld, pathFld}
	url.Path = types.Read_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &req))

	req.Fields = []*ast.CXArgument{&derefURLFld, rawPathFld}
	url.RawPath = types.Read_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &req))

	req.Fields = []*ast.CXArgument{&derefURLFld, forceQueryFld}
	url.ForceQuery = types.Read_bool(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &req))

	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	response, err := netClient.Do(&request)
	if err != nil {
		types.Write_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, errorstring), err.Error())
		return
	}

	resp := ast.CXArgument{}
	err = copier.Copy(&resp, respstruct)
	if err != nil {
		panic(err)
	}

	responseType, err := httpPkg.GetStruct("Response")
	if err != nil {
		panic(err)
	}

	statusFld, err := responseType.GetField("Status")
	if err != nil {
		panic(err)
	}
	statusCodeFld, err := responseType.GetField("StatusCode")
	if err != nil {
		panic(err)
	}
	protoFld, err := responseType.GetField("Proto")
	if err != nil {
		panic(err)
	}
	protoMajorFld, err := responseType.GetField("ProtoMajor")
	if err != nil {
		panic(err)
	}
	protoMinorFld, err := responseType.GetField("ProtoMinor")
	if err != nil {
		panic(err)
	}
	contentLengthFld, err := responseType.GetField("ContentLength")
	if err != nil {
		panic(err)
	}
	bodyFld, err := responseType.GetField("Body")
	if err != nil {
		panic(err)
	}


	resp.Fields = []*ast.CXArgument{statusFld}
	types.Write_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &resp), response.Status)

	resp.Fields = []*ast.CXArgument{statusCodeFld}
	types.Write_i32(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &resp), int32(response.StatusCode))

	resp.Fields = []*ast.CXArgument{protoFld}
	types.Write_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &resp), response.Proto)

	resp.Fields = []*ast.CXArgument{protoMajorFld}
	types.Write_i32(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &resp), int32(response.ProtoMajor))

	resp.Fields = []*ast.CXArgument{protoMinorFld}
	types.Write_i32(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &resp), int32(response.ProtoMinor))

	resp.Fields = []*ast.CXArgument{contentLengthFld}
	types.Write_i64(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &resp), response.ContentLength)

	resp.Fields = []*ast.CXArgument{bodyFld}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	types.Write_str(ast.PROGRAM.Memory, ast.GetFinalOffset(fp, &resp), string(body))
}

/*
func opDMSGDo(inputs []ast.CXValue, outputs []ast.CXValue) {
	inp1, out1 := inputs[0].Arg, outputs[0].Arg
    fp := inputs[0].FramePointer

    var req http.Request
	byts1 := ast.ReadMemory(ast.GetFinalOffset(fp, inp1), inp1)
	err := encoder.DeserializeRawExact(byts1, &req)
	if err != nil {
		ast.WriteString(fp, err.Error(), out1)
	}
}
*/
