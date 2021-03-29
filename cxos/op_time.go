// +build os

package cxos

import (
	"github.com/skycoin/cx/cx/ast"
	"time"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func opTimeUnixMilli(inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_i64(makeTimestamp())
}

func opTimeUnixNano(inputs []ast.CXValue, outputs []ast.CXValue) {
    outputs[0].Set_i64(time.Now().UnixNano())
}

func opTimeSleep(inputs []ast.CXValue, outputs []ast.CXValue) {
	time.Sleep(time.Duration(inputs[0].Get_i32()) * time.Millisecond)
}
