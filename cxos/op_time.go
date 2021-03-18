// +build os

package cxcore

import (
	"time"

	. "github.com/skycoin/cx/cx"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func opTimeUnixMilli(inputs []CXValue, outputs []CXValue) {
	outputs[0].Set_i64(makeTimestamp())
}

func opTimeUnixNano(inputs []CXValue, outputs []CXValue) {
    outputs[0].Set_i64(time.Now().UnixNano())
}

func opTimeSleep(inputs []CXValue, outputs []CXValue) {
	time.Sleep(time.Duration(inputs[0].Get_i32()) * time.Millisecond)
}
