package cxcore

import (
	"net"
	"net/http"
	"strings"

	"github.com/amherag/skycoin/src/cipher/encoder"
)

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
