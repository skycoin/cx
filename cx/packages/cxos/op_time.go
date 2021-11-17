// +build cxos

package cxos

import (
	"time"

	"github.com/skycoin/cx/cx/ast"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func opTimeUnixMilli(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_i64(prgrm, makeTimestamp())
}

func opTimeUnixNano(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_i64(prgrm, time.Now().UnixNano())
}

func opTimeSleep(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	time.Sleep(time.Duration(inputs[0].Get_i32(prgrm)) * time.Millisecond)
}
