package cxgo

import (
	cxcore "github.com/skycoin/cx/cx"
)

var Initialized bool

func init() {
	InitCXCore()
}

func InitCXCore() {
	if !Initialized {
		cxcore.LoadOpCodeTables()
		Initialized = true
	}
}
