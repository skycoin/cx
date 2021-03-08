// +build base

package cxcore

import (
	"time"

	. "github.com/skycoin/cx/cx"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func opTimeUnixMilli(expr *CXExpression, fp int) {
	WriteI64(GetFinalOffset(fp, expr.Outputs[0]), makeTimestamp())
}

func opTimeUnixNano(expr *CXExpression, fp int) {
	WriteI64(GetFinalOffset(fp, expr.Outputs[0]), time.Now().UnixNano())
}

func opTimeSleep(expr *CXExpression, fp int) {
	time.Sleep(time.Duration(ReadI32(fp, expr.Inputs[0])) * time.Millisecond)
}
