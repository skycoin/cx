// +build base opengl

package cxcore

// op codes
const (
	// time
	OP_TIME_SLEEP = iota + END_OF_CORE_OPS
	OP_TIME_UNIX
	OP_TIME_UNIX_MILLI
	OP_TIME_UNIX_NANO

	// serialize
	OP_SERIAL_PROGRAM

	// os
	OP_OS_GET_WORKING_DIRECTORY
	OP_OS_OPEN
	OP_OS_CLOSE
	OP_OS_SEEK
	OP_OS_READ_ALL_TEXT
	OP_OS_READ_F32
	OP_OS_READ_UI32
	OP_OS_READ_UI16
	OP_OS_RUN
	OP_OS_EXIT

	// json
	OP_JSON_OPEN
	OP_JSON_CLOSE
	OP_JSON_TOKEN_MORE
	OP_JSON_TOKEN_NEXT
	OP_JSON_TOKEN_TYPE
	OP_JSON_TOKEN_DELIM
	OP_JSON_TOKEN_BOOL
	OP_JSON_TOKEN_F64
	OP_JSON_TOKEN_I64
	OP_JSON_TOKEN_STR

	// http
	// OP_HTTP_GET

	// object explorer
	OP_OBJ_QUERY

	END_OF_BASE_OPS
)

func init() {
	// time
	AddOpCode(OP_TIME_SLEEP, "time.Sleep",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_TIME_UNIX_MILLI, "time.UnixMilli",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_TIME_UNIX_NANO, "time.UnixNano",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_I64, false)})

	// http
	// AddOpCode(OP_HTTP_GET, "http.Get",
	// 	[]*CXArgument{newOpPar(TYPE_STR, false)},
	// 	[]*CXArgument{newOpPar(TYPE_STR, false)})

	// os
	AddOpCode(OP_OS_GET_WORKING_DIRECTORY, "os.GetWorkingDirectory",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_OS_OPEN, "os.Open",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_OS_CLOSE, "os.Close",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_OS_SEEK, "os.Seek",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_I64, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_OS_READ_F32, "os.ReadF32",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_OS_READ_UI32, "os.ReadUI32",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_OS_READ_UI16, "os.ReadUI16",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_OS_READ_ALL_TEXT, "os.ReadAllText",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_OS_RUN, "os.Run",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_STR, false)})
	AddOpCode(OP_OS_EXIT, "os.Exit",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})

	// json
	AddOpCode(OP_JSON_OPEN, "json.Open",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_JSON_CLOSE, "json.Close",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_JSON_TOKEN_MORE, "json.More",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false), newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_JSON_TOKEN_NEXT, "json.Next",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_JSON_TOKEN_TYPE, "json.Type",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_JSON_TOKEN_DELIM, "json.Delim",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_JSON_TOKEN_BOOL, "json.Bool",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false), newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_JSON_TOKEN_F64, "json.Float64",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_JSON_TOKEN_I64, "json.Int64",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_JSON_TOKEN_STR, "json.Str",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_BOOL, false)})

	// exec
	handleOpcode := func(opCode int) opcodeHandler {
		switch opCode {
		// time
		case OP_TIME_SLEEP:
			return op_time_Sleep
		case OP_TIME_UNIX:
		case OP_TIME_UNIX_MILLI:
			return op_time_UnixMilli
		case OP_TIME_UNIX_NANO:
			return op_time_UnixNano

		// http
		// case OP_HTTP_GET:
		// 	return op_http_get

		// os
		case OP_OS_GET_WORKING_DIRECTORY:
			return op_os_GetWorkingDirectory
		case OP_OS_OPEN:
			return op_os_Open
		case OP_OS_CLOSE:
			return op_os_Close
		case OP_OS_SEEK:
			return op_os_Seek
		case OP_OS_READ_F32:
			return op_os_ReadF32
		case OP_OS_READ_UI32:
			return op_os_ReadUI32
		case OP_OS_READ_UI16:
			return op_os_ReadUI16
		case OP_OS_READ_ALL_TEXT:
			return op_os_ReadAllText
		case OP_OS_RUN:
			return op_os_Run
		case OP_OS_EXIT:
			return op_os_Exit

		// json
		case OP_JSON_OPEN:
			return opJSONOpen
		case OP_JSON_CLOSE:
			return opJSONClose
		case OP_JSON_TOKEN_MORE:
			return opJSONTokenMore
		case OP_JSON_TOKEN_NEXT:
			return opJSONTokenNext
		case OP_JSON_TOKEN_TYPE:
			return opJSONTokenType
		case OP_JSON_TOKEN_DELIM:
			return opJSONTokenDelim
		case OP_JSON_TOKEN_BOOL:
			return opJSONTokenBool
		case OP_JSON_TOKEN_F64:
			return opJSONTokenF64
		case OP_JSON_TOKEN_I64:
			return opJSONTokenI64
		case OP_JSON_TOKEN_STR:
			return opJSONTokenStr
		}

		return nil
	}

	opcodeHandlerFinders = append(opcodeHandlerFinders, handleOpcode)
}
