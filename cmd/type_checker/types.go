package type_checker

import "github.com/skycoin/cx/cx/types"

// Types Map
var TypesMap map[string]types.Code = map[string]types.Code{

	"bool": types.BOOL,

	"i8":  types.I8,
	"i16": types.I16,
	"i32": types.I32,
	"i64": types.I64,

	"ui8":  types.UI8,
	"ui16": types.UI16,
	"ui32": types.UI32,
	"ui64": types.UI64,

	"f32": types.F32,
	"f64": types.F64,

	"str": types.STR,
	"aff": types.AFF,
}
