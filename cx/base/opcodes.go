// +build base

package cxcore

import (
	. "github.com/SkycoinProject/cx/cx"
)

// op codes
const (
	// time
	OP_TIME_SLEEP = iota + END_OF_CORE_OPS
	OP_TIME_UNIX_MILLI
	OP_TIME_UNIX_NANO

	// serialize
	OP_SERIAL_PROGRAM

	// os
	OP_OS_GET_WORKING_DIRECTORY
	OP_OS_LOG_FILE
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

	// http
	// OP_HTTP_GET

	// object explorer
	OP_OBJ_QUERY

	END_OF_BASE_OPS
)

func init() {
	// time
	Op(OP_TIME_SLEEP, "time.Sleep", opTimeSleep, In(AI32), nil)
	Op(OP_TIME_UNIX_MILLI, "time.UnixMilli", opTimeUnixMilli, nil, Out(AI64))
	Op(OP_TIME_UNIX_NANO, "time.UnixNano", opTimeUnixNano, nil, Out(AI64))

	// http
	// Op(OP_HTTP_GET, "http.Get", opHttpGet, In(ASTR), Out(ASTR))

	// os
	Op(OP_OS_GET_WORKING_DIRECTORY, "os.GetWorkingDirectory", opOsGetWorkingDirectory, nil, Out(ASTR))
	Op(OP_OS_LOG_FILE, "os.LogFile", opOsLogFile, In(ABOOL), nil)
	Op(OP_OS_OPEN, "os.Open", opOsOpen, In(ASTR), Out(AI32))
	Op(OP_OS_CLOSE, "os.Close", opOsClose, In(AI32), Out(ABOOL))
	Op(OP_OS_SEEK, "os.Seek", opOsSeek, In(AI32, AI64, AI32), Out(AI64))
	Op(OP_OS_READ_F32, "os.ReadF32", opOsReadF32, In(AI32), Out(AF32, ABOOL))
	Op(OP_OS_READ_UI32, "os.ReadUI32", opOsReadUI32, In(AI32), Out(AUI32, ABOOL))
	Op(OP_OS_READ_UI16, "os.ReadUI16", opOsReadUI16, In(AI32), Out(AUI16, ABOOL))
	Op(OP_OS_READ_F32_SLICE, "os.ReadF32Slice", opOsReadF32Slice, In(AI32, AI32), Out(Slice(TYPE_F32), ABOOL))
	Op(OP_OS_READ_UI32_SLICE, "os.ReadUI32Slice", opOsReadUI32Slice, In(AI32, AI32), Out(Slice(TYPE_UI32), ABOOL))
	Op(OP_OS_READ_UI16_SLICE, "os.ReadUI16Slice", opOsReadUI16Slice, In(AI32, AI32), Out(Slice(TYPE_UI16), ABOOL))
	Op(OP_OS_READ_ALL_TEXT, "os.ReadAllText", opOsReadAllText, In(ASTR), Out(ASTR, ABOOL))
	Op(OP_OS_RUN, "os.Run", opOsRun, In(ASTR, AI32, AI32, ASTR), Out(AI32, AI32, ASTR))
	Op(OP_OS_EXIT, "os.Exit", opOsExit, In(AI32), nil)

	// json
	Op(OP_JSON_OPEN, "json.Open", opJsonOpen, In(ASTR), Out(AI32))
	Op(OP_JSON_CLOSE, "json.Close", opJsonClose, In(AI32), Out(ABOOL))
	Op(OP_JSON_TOKEN_MORE, "json.More", opJsonTokenMore, In(AI32), Out(ABOOL, ABOOL))
	Op(OP_JSON_TOKEN_NEXT, "json.Next", opJsonTokenNext, In(AI32), Out(AI32, ABOOL))
	Op(OP_JSON_TOKEN_TYPE, "json.Type", opJsonTokenType, In(AI32), Out(AI32, ABOOL))
	Op(OP_JSON_TOKEN_DELIM, "json.Delim", opJsonTokenDelim, In(AI32), Out(AI32, ABOOL))
	Op(OP_JSON_TOKEN_BOOL, "json.Bool", opJsonTokenBool, In(AI32), Out(ABOOL, ABOOL))
	Op(OP_JSON_TOKEN_F64, "json.Float64", opJsonTokenF64, In(AI32), Out(AF64, ABOOL))
	Op(OP_JSON_TOKEN_I64, "json.Int64", opJsonTokenI64, In(AI32), Out(AI64, ABOOL))
	Op(OP_JSON_TOKEN_STR, "json.Str", opJsonTokenStr, In(AI32), Out(ASTR, ABOOL))

	// profile
	Op(OP_START_CPU_PROFILE, "StartCPUProfile", opStartProfile, In(ASTR, AI32), nil)
	Op(OP_STOP_CPU_PROFILE, "StopCPUProfile", opStopProfile, In(ASTR), nil)
}
