package cxcore

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	// "github.com/skycoin/dmsg/cipher"
	// "github.com/skycoin/dmsg/disc"

	// dmsghttp "github.com/SkycoinProject/dmsg-http"
	"github.com/SkycoinProject/skycoin/src/cipher/encoder"

	"github.com/jinzhu/copier"
)

func init() {
	// In this case we're adding the `URL` type to the `http` package.
	httpPkg := MakePackage("http")
	urlStrct := MakeStruct("URL")

	urlStrct.AddField(MakeArgument("Scheme", "", 0).AddType(TypeNames[TYPE_STR]).AddPackage(httpPkg))
	urlStrct.AddField(MakeArgument("Opaque", "", 0).AddType(TypeNames[TYPE_STR]).AddPackage(httpPkg))
	urlStrct.AddField(MakeArgument("Host", "", 0).AddType(TypeNames[TYPE_STR]).AddPackage(httpPkg))
	urlStrct.AddField(MakeArgument("Path", "", 0).AddType(TypeNames[TYPE_STR]).AddPackage(httpPkg))
	urlStrct.AddField(MakeArgument("RawPath", "", 0).AddType(TypeNames[TYPE_STR]).AddPackage(httpPkg))
	urlStrct.AddField(MakeArgument("ForceQuery", "", 0).AddType(TypeNames[TYPE_BOOL]).AddPackage(httpPkg))
	urlStrct.AddField(MakeArgument("RawQuery", "", 0).AddType(TypeNames[TYPE_STR]).AddPackage(httpPkg))
	urlStrct.AddField(MakeArgument("Fragment", "", 0).AddType(TypeNames[TYPE_STR]).AddPackage(httpPkg))

	httpPkg.AddStruct(urlStrct)

	// 95% sure that you will also need Golang's `net/http`'s `Request` struct.
	// If this is the case, you can add that structure like this:

	requestStrct := MakeStruct("Request")

	requestStrct.AddField(MakeArgument("Method", "", 0).AddType(TypeNames[TYPE_STR]).AddPackage(httpPkg))
	urlFld := MakeArgument("URL", "", 0).AddType(TypeNames[TYPE_CUSTOM]).AddPackage(httpPkg)
	// Golang declares this one as a pointer. It can be done in CX using
	// `CXArgument.DeclarationSpecifiers`.
	urlFld.DeclarationSpecifiers = append(urlFld.DeclarationSpecifiers, DECL_STRUCT)
	urlFld.DeclarationSpecifiers = append(urlFld.DeclarationSpecifiers, DECL_POINTER)
	urlFld.IsPointer = true
	urlFld.Size = TYPE_POINTER_SIZE
	urlFld.TotalSize = TYPE_POINTER_SIZE
	// urlFld.PassBy = PASSBY_REFERENCE
	urlFld.CustomType = urlStrct
	requestStrct.AddField(urlFld)

	// If `Request` is indeed needed, add the other fields. If you *need* `Header`,
	// you'll be in trouble (maybe), as `Header` is of type `map[string][]string`
	// and CX doesn't have maps. If you need this, implement it as an array (or a slice), where
	// the first element is X datum, the second element is Y datum, etc. For example:
	headerFld := MakeArgument("Header", "", 0).AddType(TypeNames[TYPE_STR]).AddPackage(httpPkg) // will be a slice of strings
	headerFld.DeclarationSpecifiers = append(headerFld.DeclarationSpecifiers, DECL_SLICE)
	headerFld.DeclarationSpecifiers = append(headerFld.DeclarationSpecifiers, DECL_SLICE)
	headerFld.IsSlice = true
	headerFld.IsReference = true
	headerFld.IsArray = true
	headerFld.PassBy = PASSBY_REFERENCE
	headerFld.Lengths = []int{0, 0} // 1D slice. If it was a 2D slice it'd be []int{0, 0}

	// Then we add the field
	requestStrct.AddField(headerFld)

	requestStrct.AddField(MakeArgument("Body", "", 0).AddType(TypeNames[TYPE_STR]).AddPackage(httpPkg))

	// And adding the `Request` structure to the `http` package.
	httpPkg.AddStruct(requestStrct)

	// Sorry, there ARE functions that handle all of these operations, but they're part of the
	// parser. These files are part of the `cxcore` package, where, until now, we didn't have
	// the need to add structs like this. So for now we have to do it manually.
	// For reference, check `cxgo/actions/declarations.go`, in particular `DeclarationSpecifiers()`.

	// Mapping http.Response struct
	responseStruct := MakeStruct("Response")
	responseStruct.AddField(MakeArgument("Status", "", 0).AddType(TypeNames[TYPE_STR]).AddPackage(httpPkg))
	responseStruct.AddField(MakeArgument("StatusCode", "", 0).AddType(TypeNames[TYPE_I32]).AddPackage(httpPkg))
	responseStruct.AddField(MakeArgument("Proto", "", 0).AddType(TypeNames[TYPE_STR]).AddPackage(httpPkg))
	responseStruct.AddField(MakeArgument("ProtoMajor", "", 0).AddType(TypeNames[TYPE_I32]).AddPackage(httpPkg))
	responseStruct.AddField(MakeArgument("ProtoMinor", "", 0).AddType(TypeNames[TYPE_I32]).AddPackage(httpPkg))
	//TODO Header Header - not sure if headerFld used for http.Request can be used here
	//TODO Body io.ReadCloser
	responseStruct.AddField(MakeArgument("ContentLength", "", 0).AddType(TypeNames[TYPE_I64]).AddPackage(httpPkg))
	transferEncodingFld := MakeArgument("TransferEncoding", "", 0).AddType(TypeNames[TYPE_STR]).AddPackage(httpPkg)
	transferEncodingFld.DeclarationSpecifiers = append(transferEncodingFld.DeclarationSpecifiers, DECL_SLICE)
	transferEncodingFld.IsSlice = true
	transferEncodingFld.IsReference = true
	transferEncodingFld.IsArray = true
	transferEncodingFld.PassBy = PASSBY_REFERENCE
	transferEncodingFld.Lengths = []int{0}
	responseStruct.AddField(transferEncodingFld)
	urlStrct.AddField(MakeArgument("Close", "", 0).AddType(TypeNames[TYPE_BOOL]).AddPackage(httpPkg))
	urlStrct.AddField(MakeArgument("Uncompressed", "", 0).AddType(TypeNames[TYPE_BOOL]).AddPackage(httpPkg))
	//TODO Trailer Header
	//TODO Request *Request
	//TODO TLS *tls.ConnectionState

	httpPkg.AddStruct(responseStruct)

	PROGRAM.AddPackage(httpPkg)
}

func opHTTPHandle(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()
	inp1, inp2 := expr.Inputs[0], expr.Inputs[1]

	_ = inp1

	// Getting handler function.
	handlerPkg, err := prgrm.GetPackage(inp2.Package.Name)
	if err != nil {
		panic(err)
	}
	handlerFn, err := handlerPkg.GetFunction(inp2.Name)
	if err != nil {
		panic(err)
	}

	http.HandleFunc(ReadStr(fp, inp1), func(w http.ResponseWriter, r *http.Request) {
		tmpExpr := CXExpression{Operator: handlerFn}

		callFP := fp + PROGRAM.CallStack[PROGRAM.CallCounter].Operator.Size

		writeHTTPRequest(callFP, handlerFn.Inputs[1], r)

		i1Off := callFP + handlerFn.Inputs[0].Offset
		i1Size := handlerFn.Inputs[0].TotalSize
		i2Off := callFP + handlerFn.Inputs[1].Offset
		i2Size := handlerFn.Inputs[1].TotalSize

		i1 := make([]byte, i1Size)
		i2 := make([]byte, i1Size)

		copy(i1, PROGRAM.Memory[i1Off:i1Off+i1Size])
		copy(i2, PROGRAM.Memory[i2Off:i2Off+i2Size])

		Debug("i1", i1)
		Debug("i2", i2)

		PROGRAM.ccallback(&tmpExpr, handlerFn.Name, handlerPkg.Name, [][]byte{i1, i2})
		fmt.Fprintf(w, ReadStr(callFP, handlerFn.Inputs[0]))
	})
}

func opHTTPListenAndServe(prgrm *CXProgram) {
	expr := prgrm.GetExpr()

	fp := prgrm.GetFramePointer()
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	url := ReadStr(fp, inp1)

	err := http.ListenAndServe(url, nil)
	writeString(fp, err.Error(), out1)
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
	urlString := ReadStr(fp, inp2)
	body := ReadStr(fp, inp3)

	// workaround for http.NewRequest starts here
	// schemeIndex := strings.Index(urlString, "://")

	// scheme := urlString[:schemeIndex]
	// schemeIndex += 3
	// noSchemeURL := urlString[schemeIndex:]
	// fmt.Println(scheme)
	// fmt.Println(noSchemeURL)
	// path := ""
	// pathIndex := strings.Index(noSchemeURL, "/")
	// if pathIndex > -1 {
	// 	path = noSchemeURL[pathIndex:]
	// 	noSchemeURL = noSchemeURL[:pathIndex]
	// }
	// fmt.Println(path)

	// u := &url.URL{
	// 	Scheme: scheme,
	// 	Host:   noSchemeURL,
	// 	Path:   path,
	// }
	// fmt.Println("url done")
	// fmt.Println(u)
	// // u, err := url.Parse(urlString)
	// // if err != nil {
	// // 	writeString(fp, err.Error(), out1)
	// // }

	// var br io.Reader
	// br = bytes.NewBuffer([]byte(body))
	// fmt.Println("buffer done")
	// rc, ok := br.(io.ReadCloser)
	// if !ok && br != nil {
	// 	rc = ioutil.NopCloser(br)
	// }

	// req := &http.Request{
	// 	Method:     method,
	// 	URL:        u,
	// 	Proto:      "HTTP/1.1",
	// 	ProtoMajor: 1,
	// 	ProtoMinor: 1,
	// 	// Header:     make(Header),
	// 	Body: rc,
	// 	Host: u.Host,
	// }
	// fmt.Println("req done")
	// req.WithContext(context.Background())
	// fmt.Println("req with ctx done")

	//above is an alternative for following 3 lines of code that fail due to URL
	req, err := http.NewRequest(method, urlString, bytes.NewBuffer([]byte(body)))
	if err != nil {
		writeString(fp, err.Error(), out1)
	}

	// // new request ends with the following, bellow is just a workaround
	// out1Offset := GetFinalOffset(fp, out1)
	// byts := encoder.Serialize(req)
	// WriteObject(out1Offset, byts)

	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	fmt.Println("client done")
	resp, err := netClient.Do(req)
	if err != nil {
		fmt.Println("err on do")
		fmt.Println(err)
		writeString(fp, err.Error(), out1)
	}
	resp1 := *resp // dereference to exclude pointer issue
	fmt.Println(resp1)

	// TODO issue with returning response,
	// the line where resp is serialized (byts := encoder.Serialize(resp)) throws following error, adding response object content for context:
	// &{404 Not Found 404 HTTP/1.1 1 1 map[Content-Length:[19] Content-Type:[text/plain; charset=utf-8] Date:[Sun, 27 Oct 2019 22:10:14 GMT] X-Content-Type-Options:[nosniff]] 0xc0000c8700 19 [] false false map[] 0xc000175400 <nil>}
	// 2019/10/27 23:10:14 invalid type int
	// error: examples/http-serve-and-request-mine.cx:8, CX_RUNTIME_ERROR, invalid type int

	out1Offset := GetFinalOffset(fp, out1)
	byts := encoder.Serialize(resp1)
	WriteObject(out1Offset, byts)

	// req.Header["Timestamp"] := time.Now().UnixNano()

	// place this on response once handler is supported
	// req.Header["Blockchain"] :=
	// req.Header["NodeID"] :=
	// req.Header["BlockHeight"] :=
	// req.Header["HeadBlockHash"] :=
	// req.Header["HeadBlockSeq"] :=
	// req.Header["HeadBlockTime"] :=
	// req.Header["RequestId"] :=
	// req.Header["ResponseError"] :=

	// url = fmt.Sprintf("http://127.0.0.1:%d/api/v1/injectTransaction", options.port + 420)
	// dataMap = make(map[string]interface{}, 0)
	// dataMap["rawtx"] = respBody["encoded_transaction"]

	// jsonStr, err = json.Marshal(dataMap)
	// if err != nil {
	// 	panic(err)
	// }

	// req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-CSRF-Token", csrfToken)
	// req.Header.Set("Content-Type", "application/json")

	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// body, err = ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
}

func writeHTTPRequest(fp int, param *CXArgument, request *http.Request) {
	req := CXArgument{}
	copier.Copy(&req, param)
	// req.DereferenceOperations = append(req.DereferenceOperations, DEREF_POINTER)

	httpPkg, err := PROGRAM.GetPackage("http")
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

	derefUrlFld := CXArgument{}
	copier.Copy(&derefUrlFld, urlFld)
	
	derefUrlFld.DereferenceOperations = append(derefUrlFld.DereferenceOperations, DEREF_POINTER)

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

	accessMethod := []*CXArgument{methodFld}
	accessURLScheme := []*CXArgument{&derefUrlFld, schemeFld}
	accessURLHost := []*CXArgument{&derefUrlFld, hostFld}
	accessURLPath := []*CXArgument{&derefUrlFld, pathFld}
	accessURLRawPath := []*CXArgument{&derefUrlFld, rawPathFld}
	accessURLForceQuery := []*CXArgument{&derefUrlFld, forceQueryFld}

	req.Fields = accessMethod
	writeString(fp, request.Method, &req)
	Debug("method", request.Method, fp, req.Offset)
	Debug("memory", PROGRAM.Memory[fp:fp+8])
	req.Fields = accessURLScheme
	writeString(fp, request.URL.Scheme, &req)
	req.Fields = accessURLHost
	writeString(fp, request.URL.Host, &req)
	req.Fields = accessURLPath
	writeString(fp, request.URL.Path, &req)
	req.Fields = accessURLRawPath
	writeString(fp, request.URL.RawPath, &req)
	req.Fields = accessURLForceQuery
	WriteMemory(GetFinalOffset(fp, &req), FromBool(request.URL.ForceQuery))
}

func opHTTPDo(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, out1, out2 := expr.Inputs[0], expr.Outputs[0], expr.Outputs[1]
	//TODO read req from the inputs
	// reqByts := ReadMemory(GetFinalOffset(fp, inp1), inp1)

	_ = out1
	_ = out2

	// var methodOffset int32
	// encoder.DeserializeAtomic(reqByts[:TYPE_POINTER_SIZE], &methodOffset)
	
	// var urlStrctOffset int32
	// encoder.DeserializeAtomic(reqByts[TYPE_POINTER_SIZE:TYPE_POINTER_SIZE * 2], &urlStrctOffset)

	req := CXArgument{}
	copier.Copy(&req, inp1)

	httpPkg, err := PROGRAM.GetPackage("http")
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

	derefUrlFld := CXArgument{}
	copier.Copy(&derefUrlFld, urlFld)
	
	derefUrlFld.DereferenceOperations = append(derefUrlFld.DereferenceOperations, DEREF_POINTER)

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

	accessMethod := []*CXArgument{methodFld}
	accessURLScheme := []*CXArgument{&derefUrlFld, schemeFld}
	accessURLHost := []*CXArgument{&derefUrlFld, hostFld}
	accessURLPath := []*CXArgument{&derefUrlFld, pathFld}
	accessURLRawPath := []*CXArgument{&derefUrlFld, rawPathFld}
	accessURLForceQuery := []*CXArgument{&derefUrlFld, forceQueryFld}

	request := http.Request{}
	url := url.URL{}
	request.URL = &url
	
	req.Fields = accessMethod
	request.Method = ReadStr(fp, &req)
	req.Fields = accessURLScheme
	url.Scheme = ReadStr(fp, &req)
	req.Fields = accessURLHost
	url.Host = ReadStr(fp, &req)
	req.Fields = accessURLPath
	url.Path = ReadStr(fp, &req)
	req.Fields = accessURLRawPath
	url.RawPath = ReadStr(fp, &req)
	req.Fields = accessURLForceQuery
	url.ForceQuery = ReadBool(fp, &req)

	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	response, err := netClient.Do(&request)
	if err != nil {
		writeString(fp, err.Error(), out2)
		return
	}

	resp := CXArgument{}
	copier.Copy(&resp, out1)

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
	// transferEncodingFld, err := responseType.GetField("TransferEncoding")
	// if err != nil {
	// 	panic(err)
	// }

	accessStatus := []*CXArgument{statusFld}
	accessStatusCode := []*CXArgument{statusCodeFld}
	accessProto := []*CXArgument{protoFld}
	accessProtoMajor := []*CXArgument{protoMajorFld}
	accessProtoMinor := []*CXArgument{protoMinorFld}
	accessContentLength := []*CXArgument{contentLengthFld}
	// accessTransferEncoding := []*CXArgument{transferEncodingFld}

	resp.Fields = accessStatus
	writeString(fp, response.Status, &resp)
	resp.Fields = accessStatusCode
	WriteMemory(GetFinalOffset(fp, &resp), FromI32(int32(response.StatusCode)))
	resp.Fields = accessProto
	writeString(fp, response.Proto, &resp)
	resp.Fields = accessProtoMajor
	WriteMemory(GetFinalOffset(fp, &resp), FromI32(int32(response.ProtoMajor)))
	resp.Fields = accessProtoMinor
	WriteMemory(GetFinalOffset(fp, &resp), FromI32(int32(response.ProtoMinor)))
	resp.Fields = accessContentLength
	WriteMemory(GetFinalOffset(fp, &resp), FromI64(int64(response.ContentLength)))

	// out1Offset := GetFinalOffset(fp, out1)
	// byts := encoder.Serialize(resp)
	// WriteObject(out1Offset, byts)
}

func opDMSGDo(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	var req http.Request
	byts1 := ReadMemory(GetFinalOffset(fp, inp1), inp1)
	err := encoder.DeserializeRawExact(byts1, &req)
	if err != nil {
		writeString(fp, err.Error(), out1)
	}

	// cPK, cSK := cipher.GenerateKeyPair()
	// dmsgD := disc.NewHTTP("http://dmsg.discovery.skywire.skycoin.com")
	// c := dmsghttp.DMSGClient(dmsgD, cPK, cSK)

	// resp, err := c.Do(&req)
	// if err != nil {
	// 	writeString(fp, err.Error(), out1)
	// }
}
