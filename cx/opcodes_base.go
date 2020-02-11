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
	OP_OS_READ_F32_SLICE
	OP_OS_READ_UI32_SLICE
	OP_OS_READ_UI16_SLICE
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

	// profile
	OP_START_CPU_PROFILE
	OP_STOP_CPU_PROFILE

	// object explorer
	OP_OBJ_QUERY

	END_OF_BASE_OPS
)

func init() {
	// time
	AddOpCode(OP_TIME_SLEEP, "time.Sleep",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_TIME_UNIX_MILLI, "time.UnixMilli",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_TIME_UNIX_NANO, "time.UnixNano",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})

	// os
	AddOpCode(OP_OS_GET_WORKING_DIRECTORY, "os.GetWorkingDirectory",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_OS_OPEN, "os.Open",
		[]*CXArgument{newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_OS_CLOSE, "os.Close",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_OS_SEEK, "os.Seek",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI64NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_OS_READ_F32, "os.ReadF32",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice), newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_OS_READ_UI32, "os.ReadUI32",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamUI32NotSlice), newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_OS_READ_UI16, "os.ReadUI16",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamUI16NotSlice), newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_OS_READ_F32_SLICE, "os.ReadF32Slice",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32Slice), newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_OS_READ_UI32_SLICE, "os.ReadUI32Slice",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamUI32Slice), newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_OS_READ_UI16_SLICE, "os.ReadUI16Slice",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamUI16Slice), newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_OS_READ_ALL_TEXT, "os.ReadAllText",
		[]*CXArgument{newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice), newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_OS_RUN, "os.Run",
		[]*CXArgument{newOpPar(opParamStrNotSlice), newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice), newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice), newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_OS_EXIT, "os.Exit",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{})

	// json
	AddOpCode(OP_JSON_OPEN, "json.Open",
		[]*CXArgument{newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_JSON_CLOSE, "json.Close",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_JSON_TOKEN_MORE, "json.More",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice), newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_JSON_TOKEN_NEXT, "json.Next",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_JSON_TOKEN_TYPE, "json.Type",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_JSON_TOKEN_DELIM, "json.Delim",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_JSON_TOKEN_BOOL, "json.Bool",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice), newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_JSON_TOKEN_F64, "json.Float64",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice), newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_JSON_TOKEN_I64, "json.Int64",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_JSON_TOKEN_STR, "json.Str",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice), newOpPar(opParamBoolNotSlice)})

	// profile
	AddOpCode(OP_START_CPU_PROFILE, "StartCPUProfile",
		[]*CXArgument{newOpPar(opParamStrNotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_STOP_CPU_PROFILE, "StopCPUProfile",
		[]*CXArgument{newOpPar(opParamStrNotSlice)},
		[]*CXArgument{})

	opcodeHandlerFinders = append(opcodeHandlerFinders, handleBaseOpcode)
}

func handleBaseOpcode(opCode int) opcodeHandler {
	switch opCode {
	// time
	case OP_TIME_SLEEP:
		return op_time_Sleep
	case OP_TIME_UNIX:
	case OP_TIME_UNIX_MILLI:
		return op_time_UnixMilli
	case OP_TIME_UNIX_NANO:
		return op_time_UnixNano

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
	case OP_OS_READ_F32_SLICE:
		return op_os_ReadF32Slice
	case OP_OS_READ_UI32_SLICE:
		return op_os_ReadUI32Slice
	case OP_OS_READ_UI16_SLICE:
		return op_os_ReadUI16Slice
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

	// profile
	case OP_START_CPU_PROFILE:
		return opStartProfile
	case OP_STOP_CPU_PROFILE:
		return opStopProfile
	}

	return nil
}
