package globals

import "github.com/skycoin/cx/cx"

var (
	// OpNames ...
	OpNames = map[int]string{}

	// OpCodes ...
	OpCodes = map[string]int{}

	// Versions ...
	OpVersions = map[int]int{}
)

var (
	OpcodeHandlers    []cxcore.OpcodeHandler
	OpcodeHandlers_V2 []cxcore.OpcodeHandler_V2
)

