// +build os

package cxos

import (
	"time"

	"github.com/skycoin/cx/cx"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func opTimeUnixMilli(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	outputs[0].Set_i64(makeTimestamp())
}

func opTimeUnixNano(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
    outputs[0].Set_i64(time.Now().UnixNano())
}

func opTimeSleep(inputs []cxcore.CXValue, outputs []cxcore.CXValue) {
	time.Sleep(time.Duration(inputs[0].Get_i32()) * time.Millisecond)
}
