package cxcore

import (
	"bytes"
	"net/http"

	"github.com/amherag/skycoin/src/cipher/encoder"
)

func opHTTPServe(prgrm *CXProgram) {
	expr := prgrm.GetExpr()

	fp := prgrm.GetFramePointer()
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	_ = ReadMemory(fp, inp1)

	// FIXME - need to extract net.Listener from input...
	// FIXME - figure out what to do with handler...
	err := http.Serve(nil, nil)
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
	body := ReadMemory(fp, inp3)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	if err != nil {
		writeString(fp, err.Error(), out1)
	}

	out1Offset := GetFinalOffset(fp, out1)
	byts := encoder.Serialize(req)
	WriteObject(out1Offset, byts)
}
