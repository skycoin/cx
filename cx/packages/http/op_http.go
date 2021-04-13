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
	"github.com/skycoin/cx/cx/helper"

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

	http.HandleFunc(ast.ReadStr(fp, urlstring), func(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprint(w, ast.ReadStr(callFP, handlerFn.Inputs[0]))
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
	method := ast.ReadStr(fp, stringmethod)
	urlString := ast.ReadStr(fp, stringurl)
	body := ast.ReadStr(fp, stringbody)

	//above is an alternative for following 3 lines of code that fail due to URL
	req, err := http.NewRequest(method, urlString, bytes.NewBuffer([]byte(body)))
	if err != nil {
		ast.WriteString(fp, err.Error(), errorstring)
	}

	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := netClient.Do(req)
	if err != nil {
		ast.WriteString(fp, err.Error(), errorstring)
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
	ast.WriteObject(out1Offset, byts)
}

func writeHTTPRequest(fp int, param *ast.CXArgument, request *http.Request) {
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

	accessMethod := []*ast.CXArgument{methodFld}
	accessBody := []*ast.CXArgument{bodyFld}
	accessURL := []*ast.CXArgument{urlFld}
	accessURLScheme := []*ast.CXArgument{&derefURLFld, schemeFld}
	accessURLHost := []*ast.CXArgument{&derefURLFld, hostFld}
	accessURLPath := []*ast.CXArgument{&derefURLFld, pathFld}
	accessURLRawPath := []*ast.CXArgument{&derefURLFld, rawPathFld}
	accessURLForceQuery := []*ast.CXArgument{&derefURLFld, forceQueryFld}

	// Creating empty `http.Request` object on heap.
	reqOff := ast.WriteObjectData(make([]byte, requestType.Size))
	reqOffByts := encoder.SerializeAtomic(int32(reqOff))
	ast.WriteMemory(ast.GetFinalOffset(fp, &req), reqOffByts)

	req.DereferenceOperations = append(req.DereferenceOperations, constants.DEREF_POINTER)

	// Creating empty `http.URL` object on heap.
	req.Fields = accessURL
	urlOff := ast.WriteObjectData(make([]byte, urlType.Size))
	urlOffByts := encoder.SerializeAtomic(int32(urlOff))
	ast.WriteMemory(ast.GetFinalOffset(fp, &req), urlOffByts)

	req.Fields = accessMethod
	ast.WriteString(fp, request.Method, &req)

	req.Fields = accessBody
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	ast.WriteString(fp, string(body), &req)
	req.Fields = accessURLScheme
	ast.WriteString(fp, request.URL.Scheme, &req)
	req.Fields = accessURLHost
	ast.WriteString(fp, request.URL.Host, &req)
	req.Fields = accessURLPath
	ast.WriteString(fp, request.URL.Path, &req)
	req.Fields = accessURLRawPath
	ast.WriteString(fp, request.URL.RawPath, &req)
	req.Fields = accessURLForceQuery
	ast.WriteMemory(ast.GetFinalOffset(fp, &req), helper.FromBool(request.URL.ForceQuery))
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

	accessMethod := []*ast.CXArgument{methodFld}
	accessURLScheme := []*ast.CXArgument{&derefURLFld, schemeFld}
	accessURLHost := []*ast.CXArgument{&derefURLFld, hostFld}
	accessURLPath := []*ast.CXArgument{&derefURLFld, pathFld}
	accessURLRawPath := []*ast.CXArgument{&derefURLFld, rawPathFld}
	accessURLForceQuery := []*ast.CXArgument{&derefURLFld, forceQueryFld}

	request := http.Request{}
	url := url.URL{}
	request.URL = &url

	req.Fields = accessMethod
	request.Method = ast.ReadStr(fp, &req)
	req.Fields = accessURLScheme
	url.Scheme = ast.ReadStr(fp, &req)
	req.Fields = accessURLHost
	url.Host = ast.ReadStr(fp, &req)
	req.Fields = accessURLPath
	url.Path = ast.ReadStr(fp, &req)
	req.Fields = accessURLRawPath
	url.RawPath = ast.ReadStr(fp, &req)
	req.Fields = accessURLForceQuery
	url.ForceQuery = ast.ReadBool(fp, &req)

	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	response, err := netClient.Do(&request)
	if err != nil {
		ast.WriteString(fp, err.Error(), errorstring)
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

	accessStatus := []*ast.CXArgument{statusFld}
	accessStatusCode := []*ast.CXArgument{statusCodeFld}
	accessProto := []*ast.CXArgument{protoFld}
	accessProtoMajor := []*ast.CXArgument{protoMajorFld}
	accessProtoMinor := []*ast.CXArgument{protoMinorFld}
	accessContentLength := []*ast.CXArgument{contentLengthFld}
	accessBody := []*ast.CXArgument{bodyFld}

	resp.Fields = accessStatus
	ast.WriteString(fp, response.Status, &resp)
	resp.Fields = accessStatusCode
	ast.WriteMemory(ast.GetFinalOffset(fp, &resp), helper.FromI32(int32(response.StatusCode)))
	resp.Fields = accessProto
	ast.WriteString(fp, response.Proto, &resp)
	resp.Fields = accessProtoMajor
	ast.WriteMemory(ast.GetFinalOffset(fp, &resp), helper.FromI32(int32(response.ProtoMajor)))
	resp.Fields = accessProtoMinor
	ast.WriteMemory(ast.GetFinalOffset(fp, &resp), helper.FromI32(int32(response.ProtoMinor)))
	resp.Fields = accessContentLength
	ast.WriteMemory(ast.GetFinalOffset(fp, &resp), helper.FromI64(int64(response.ContentLength)))
	resp.Fields = accessBody
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	ast.WriteString(fp, string(body), &resp)
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
