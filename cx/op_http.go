// +build base

package cxcore

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"time"

	// "github.com/skycoin/dmsg/cipher"
	// "github.com/skycoin/dmsg/disc"

	// dmsghttp "github.com/SkycoinProject/dmsg-http"
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
	headerFld.Size = 4
	headerFld.TotalSize = 4

	// Then we add the field
	requestStrct.AddField(headerFld)

	Debug("strctSize", requestStrct.Size)

	// And adding the `Request` structure to the `http` package.
	httpPkg.AddStruct(requestStrct)

	// Debug("urlOffset", methodFld.Offset)

	// Sorry, there ARE functions that handle all of these operations, but they're part of the
	// parser. These files are part of the `cxcore` package, where, until now, we didn't have
	// the need to add structs like this. So for now we have to do it manually.
	// For reference, check `cxgo/actions/declarations.go`, in particular `DeclarationSpecifiers()`.

	// Mapping http.Response struct
	responseStruct := MakeStruct("Response")
	responseStruct.AddField(MakeArgument("Status", "", 0).AddType(TypeNames[TYPE_STR]))
	responseStruct.AddField(MakeArgument("StatusCode", "", 0).AddType(TypeNames[TYPE_I16]))
	responseStruct.AddField(MakeArgument("Proto", "", 0).AddType(TypeNames[TYPE_STR]))
	responseStruct.AddField(MakeArgument("ProtoMajor", "", 0).AddType(TypeNames[TYPE_I16]))
	responseStruct.AddField(MakeArgument("ProtoMinor", "", 0).AddType(TypeNames[TYPE_I16]))
	//TODO Header Header - not sure if headerFld used for http.Request can be used here
	//TODO Body io.ReadCloser
	responseStruct.AddField(MakeArgument("ContentLength", "", 0).AddType(TypeNames[TYPE_I64]))
	transferEncodingFld := MakeArgument("TransferEncoding", "", 0).AddType(TypeNames[TYPE_STR])
	transferEncodingFld.DeclarationSpecifiers = append(headerFld.DeclarationSpecifiers, DECL_SLICE)
	transferEncodingFld.IsSlice = true
	transferEncodingFld.IsReference = true
	transferEncodingFld.IsArray = true
	transferEncodingFld.PassBy = PASSBY_REFERENCE
	transferEncodingFld.Lengths = []int{0}
	responseStruct.AddField(transferEncodingFld)
	urlStrct.AddField(MakeArgument("Close", "", 0).AddType(TypeNames[TYPE_BOOL]))
	urlStrct.AddField(MakeArgument("Uncompressed", "", 0).AddType(TypeNames[TYPE_BOOL]))
	//TODO Trailer Header
	//TODO Request *Request
	//TODO TLS *tls.ConnectionState

	httpPkg.AddStruct(responseStruct)

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

func opHTTPDo(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	//TODO read req from the inputs
	// var req http.Request
	reqByts := ReadMemory(GetFinalOffset(fp, inp1), inp1)

	_ = out1

	var methodOffset int32
	encoder.DeserializeAtomic(reqByts[:TYPE_POINTER_SIZE], &methodOffset)
	
	Debug("method", ReadStrOffset(methodOffset))

	// var urlStrctOffset int32
	// encoder.DeserializeAtomic(reqByts[TYPE_POINTER_SIZE:TYPE_POINTER_SIZE * 2], &urlStrctOffset)

	// urlByts := PROGRAM.Memory[urlStrctOffset:urlStrctOffset+TYPE_POINTER_SIZE]

	// var schemeOffset int32
	// encoder.DeserializeAtomic(urlByts[:TYPE_POINTER_SIZE], &schemeOffset)

	// Debug("url.scheme", urlByts)
	// Debug("url.scheme", ReadStrOffset(schemeOffset))
	
	// err := encoder.DeserializeRawExact(byts1, &req)
	// if err != nil {
	// 	writeString(fp, err.Error(), out1)
	// }
	// Debug("path", req.URL.Path)%
	// fmt.Printf("%+v\n", req)
	// writeString(fp, req.URL.Path, out1) // This is just for testing, confirming if request is red correctly should be removed
	// return

	// var netClient = &http.Client{
	// 	Timeout: time.Second * 30,
	// }
	// resp, err := netClient.Do(&req)
	// if err != nil {
	// 	writeString(fp, err.Error(), out1)
	// }

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
