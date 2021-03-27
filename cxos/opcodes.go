// +build os

package cxos

import (
	. "github.com/skycoin/cx/cx"
	"github.com/skycoin/cx/cx/constants"
)

// op codes
const (
	// time
	OP_TIME_SLEEP = iota + constants.END_OF_CORE_OPS
	OP_TIME_UNIX_MILLI
	OP_TIME_UNIX_NANO

	// serialize
	OP_SERIAL_PROGRAM

	// os
	OP_OS_GET_WORKING_DIRECTORY
	OP_OS_LOG_FILE
	OP_OS_OPEN
	OP_OS_CREATE
	OP_OS_CLOSE
	OP_OS_SEEK
	OP_OS_READ_ALL_TEXT
	OP_OS_READ_STR
	OP_OS_READ_F64
	OP_OS_READ_F32
	OP_OS_READ_UI64
	OP_OS_READ_UI32
	OP_OS_READ_UI16
	OP_OS_READ_UI8
	OP_OS_READ_I64
	OP_OS_READ_I32
	OP_OS_READ_I16
	OP_OS_READ_I8
	OP_OS_READ_BOOL
	OP_OS_READ_F64_SLICE
	OP_OS_READ_F32_SLICE
	OP_OS_READ_UI64_SLICE
	OP_OS_READ_UI32_SLICE
	OP_OS_READ_UI16_SLICE
	OP_OS_READ_UI8_SLICE
	OP_OS_READ_I64_SLICE
	OP_OS_READ_I32_SLICE
	OP_OS_READ_I16_SLICE
	OP_OS_READ_I8_SLICE
	OP_OS_WRITE_STR
	OP_OS_WRITE_F64
	OP_OS_WRITE_F32
	OP_OS_WRITE_UI64
	OP_OS_WRITE_UI32
	OP_OS_WRITE_UI16
	OP_OS_WRITE_UI8
	OP_OS_WRITE_I64
	OP_OS_WRITE_I32
	OP_OS_WRITE_I16
	OP_OS_WRITE_I8
	OP_OS_WRITE_BOOL
	OP_OS_WRITE_F64_SLICE
	OP_OS_WRITE_F32_SLICE
	OP_OS_WRITE_UI64_SLICE
	OP_OS_WRITE_UI32_SLICE
	OP_OS_WRITE_UI16_SLICE
	OP_OS_WRITE_UI8_SLICE
	OP_OS_WRITE_I64_SLICE
	OP_OS_WRITE_I32_SLICE
	OP_OS_WRITE_I16_SLICE
	OP_OS_WRITE_I8_SLICE
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

	// regexp
	OP_REGEXP_COMPILE
	OP_REGEXP_MUST_COMPILE
	OP_REGEXP_FIND

	// cipher
	OP_CIPHER_GENERATE_KEY_PAIR

	END_OF_OS_OPS
)

func init() {
	// time
	Op_V2(OP_TIME_SLEEP, "time.Sleep", opTimeSleep, In(AI32), nil)
	Op_V2(OP_TIME_UNIX_MILLI, "time.UnixMilli", opTimeUnixMilli, nil, Out(AI64))
	Op_V2(OP_TIME_UNIX_NANO, "time.UnixNano", opTimeUnixNano, nil, Out(AI64))

	// http
	// Op(OP_HTTP_GET, "http.Get", opHttpGet, In(ASTR), Out(ASTR))

	// os
	Op_V2(OP_OS_GET_WORKING_DIRECTORY, "os.GetWorkingDirectory", opOsGetWorkingDirectory, nil, Out(ASTR))
	Op_V2(OP_OS_LOG_FILE, "os.LogFile", opOsLogFile, In(ABOOL), nil)
	Op_V2(OP_OS_OPEN, "os.Open", opOsOpen, In(ASTR), Out(AI32))
	Op_V2(OP_OS_CREATE, "os.Create", opOsCreate, In(ASTR), Out(AI32))
	Op_V2(OP_OS_CLOSE, "os.Close", opOsClose, In(AI32), Out(ABOOL))
	Op_V2(OP_OS_SEEK, "os.Seek", opOsSeek, In(AI32, AI64, AI32), Out(AI64))
	Op_V2(OP_OS_READ_STR, "os.ReadStr", opOsReadStr, In(AI32), Out(ASTR, ABOOL))
	Op_V2(OP_OS_READ_F64, "os.ReadF64", opOsReadF64, In(AI32), Out(AF64, ABOOL))
	Op_V2(OP_OS_READ_F32, "os.ReadF32", opOsReadF32, In(AI32), Out(AF32, ABOOL))
	Op_V2(OP_OS_READ_UI64, "os.ReadUI64", opOsReadUI64, In(AI32), Out(AUI64, ABOOL))
	Op_V2(OP_OS_READ_UI32, "os.ReadUI32", opOsReadUI32, In(AI32), Out(AUI32, ABOOL))
	Op_V2(OP_OS_READ_UI16, "os.ReadUI16", opOsReadUI16, In(AI32), Out(AUI16, ABOOL))
	Op_V2(OP_OS_READ_UI8, "os.ReadUI8", opOsReadUI8, In(AI32), Out(AUI8, ABOOL))
	Op_V2(OP_OS_READ_I64, "os.ReadI64", opOsReadI64, In(AI32), Out(AI64, ABOOL))
	Op_V2(OP_OS_READ_I32, "os.ReadI32", opOsReadI32, In(AI32), Out(AI32, ABOOL))
	Op_V2(OP_OS_READ_I16, "os.ReadI16", opOsReadI16, In(AI32), Out(AI16, ABOOL))
	Op_V2(OP_OS_READ_I8, "os.ReadI8", opOsReadI8, In(AI32), Out(AI8, ABOOL))
	Op_V2(OP_OS_READ_BOOL, "os.ReadBOOL", opOsReadBOOL, In(AI32), Out(ABOOL, ABOOL))
	Op_V2(OP_OS_READ_F64_SLICE, "os.ReadF64Slice", opOsReadF64Slice, In(AI32, Slice(constants.TYPE_F64), AUI64), Out(Slice(constants.TYPE_F64), ABOOL))
	Op_V2(OP_OS_READ_F32_SLICE, "os.ReadF32Slice", opOsReadF32Slice, In(AI32, Slice(constants.TYPE_F32), AUI64), Out(Slice(constants.TYPE_F32), ABOOL))
	Op_V2(OP_OS_READ_UI64_SLICE, "os.ReadUI64Slice", opOsReadUI64Slice, In(AI32, Slice(constants.TYPE_UI64), AUI64), Out(Slice(constants.TYPE_UI64), ABOOL))
	Op_V2(OP_OS_READ_UI32_SLICE, "os.ReadUI32Slice", opOsReadUI32Slice, In(AI32, Slice(constants.TYPE_UI32), AUI64), Out(Slice(constants.TYPE_UI32), ABOOL))
	Op_V2(OP_OS_READ_UI16_SLICE, "os.ReadUI16Slice", opOsReadUI16Slice, In(AI32, Slice(constants.TYPE_UI16), AUI64), Out(Slice(constants.TYPE_UI16), ABOOL))
	Op_V2(OP_OS_READ_UI8_SLICE, "os.ReadUI8Slice", opOsReadUI8Slice, In(AI32, Slice(constants.TYPE_UI8), AUI64), Out(Slice(constants.TYPE_UI8), ABOOL))
	Op_V2(OP_OS_READ_I64_SLICE, "os.ReadI64Slice", opOsReadI64Slice, In(AI32, Slice(constants.TYPE_I64), AUI64), Out(Slice(constants.TYPE_I64), ABOOL))
	Op_V2(OP_OS_READ_I32_SLICE, "os.ReadI32Slice", opOsReadI32Slice, In(AI32, Slice(constants.TYPE_I32), AUI64), Out(Slice(constants.TYPE_I32), ABOOL))
	Op_V2(OP_OS_READ_I16_SLICE, "os.ReadI16Slice", opOsReadI16Slice, In(AI32, Slice(constants.TYPE_I16), AUI64), Out(Slice(constants.TYPE_I16), ABOOL))
	Op_V2(OP_OS_READ_I8_SLICE, "os.ReadI8Slice", opOsReadI8Slice, In(AI32, Slice(constants.TYPE_I8), AUI64), Out(Slice(constants.TYPE_I8), ABOOL))
	Op_V2(OP_OS_READ_ALL_TEXT, "os.ReadAllText", opOsReadAllText, In(ASTR), Out(ASTR, ABOOL))
	Op_V2(OP_OS_WRITE_STR, "os.WriteStr", opOsWriteStr, In(AI32, ASTR), Out(ABOOL))
	Op_V2(OP_OS_WRITE_F64, "os.WriteF64", opOsWriteF64, In(AI32, AF64), Out(ABOOL))
	Op_V2(OP_OS_WRITE_F32, "os.WriteF32", opOsWriteF32, In(AI32, AF32), Out(ABOOL))
	Op_V2(OP_OS_WRITE_UI64, "os.WriteUI64", opOsWriteUI64, In(AI32, AUI64), Out(ABOOL))
	Op_V2(OP_OS_WRITE_UI32, "os.WriteUI32", opOsWriteUI32, In(AI32, AUI32), Out(ABOOL))
	Op_V2(OP_OS_WRITE_UI16, "os.WriteUI16", opOsWriteUI16, In(AI32, AUI16), Out(ABOOL))
	Op_V2(OP_OS_WRITE_UI8, "os.WriteUI8", opOsWriteUI8, In(AI32, AUI8), Out(ABOOL))
	Op_V2(OP_OS_WRITE_I64, "os.WriteI64", opOsWriteI64, In(AI32, AI64), Out(ABOOL))
	Op_V2(OP_OS_WRITE_I32, "os.WriteI32", opOsWriteI32, In(AI32, AI32), Out(ABOOL))
	Op_V2(OP_OS_WRITE_I16, "os.WriteI16", opOsWriteI16, In(AI32, AI16), Out(ABOOL))
	Op_V2(OP_OS_WRITE_I8, "os.WriteI8", opOsWriteI8, In(AI32, AI8), Out(ABOOL))
	Op_V2(OP_OS_WRITE_BOOL, "os.WriteBOOL", opOsWriteBOOL, In(AI32, ABOOL), Out(ABOOL))
	Op_V2(OP_OS_WRITE_F64_SLICE, "os.WriteF64Slice", opOsWriteF64Slice, In(AI32, Slice(constants.TYPE_F64)), Out(ABOOL))
	Op_V2(OP_OS_WRITE_F32_SLICE, "os.WriteF32Slice", opOsWriteF32Slice, In(AI32, Slice(constants.TYPE_F32)), Out(ABOOL))
	Op_V2(OP_OS_WRITE_UI64_SLICE, "os.WriteUI64Slice", opOsWriteUI64Slice, In(AI32, Slice(constants.TYPE_UI64)), Out(ABOOL))
	Op_V2(OP_OS_WRITE_UI32_SLICE, "os.WriteUI32Slice", opOsWriteUI32Slice, In(AI32, Slice(constants.TYPE_UI32)), Out(ABOOL))
	Op_V2(OP_OS_WRITE_UI16_SLICE, "os.WriteUI16Slice", opOsWriteUI16Slice, In(AI32, Slice(constants.TYPE_UI16)), Out(ABOOL))
	Op_V2(OP_OS_WRITE_UI8_SLICE, "os.WriteUI8Slice", opOsWriteUI8Slice, In(AI32, Slice(constants.TYPE_UI8)), Out(ABOOL))
	Op_V2(OP_OS_WRITE_I64_SLICE, "os.WriteI64Slice", opOsWriteI64Slice, In(AI32, Slice(constants.TYPE_I64)), Out(ABOOL))
	Op_V2(OP_OS_WRITE_I32_SLICE, "os.WriteI32Slice", opOsWriteI32Slice, In(AI32, Slice(constants.TYPE_I32)), Out(ABOOL))
	Op_V2(OP_OS_WRITE_I16_SLICE, "os.WriteI16Slice", opOsWriteI16Slice, In(AI32, Slice(constants.TYPE_I16)), Out(ABOOL))
	Op_V2(OP_OS_WRITE_I8_SLICE, "os.WriteI8Slice", opOsWriteI8Slice, In(AI32, Slice(constants.TYPE_I8)), Out(ABOOL))

	Op_V2(OP_OS_RUN, "os.Run", opOsRun, In(ASTR, AI32, AI32, ASTR), Out(AI32, AI32, ASTR))
	Op_V2(OP_OS_EXIT, "os.Exit", opOsExit, In(AI32), nil)

	// json
	Op_V2(OP_JSON_OPEN, "json.Open", opJsonOpen, In(ASTR), Out(AI32))
	Op_V2(OP_JSON_CLOSE, "json.Close", opJsonClose, In(AI32), Out(ABOOL))
	Op_V2(OP_JSON_TOKEN_MORE, "json.More", opJsonTokenMore, In(AI32), Out(ABOOL, ABOOL))
	Op_V2(OP_JSON_TOKEN_NEXT, "json.Next", opJsonTokenNext, In(AI32), Out(AI32, ABOOL))
	Op_V2(OP_JSON_TOKEN_TYPE, "json.Type", opJsonTokenType, In(AI32), Out(AI32, ABOOL))
	Op_V2(OP_JSON_TOKEN_DELIM, "json.Delim", opJsonTokenDelim, In(AI32), Out(AI32, ABOOL))
	Op_V2(OP_JSON_TOKEN_BOOL, "json.Bool", opJsonTokenBool, In(AI32), Out(ABOOL, ABOOL))
	Op_V2(OP_JSON_TOKEN_F64, "json.Float64", opJsonTokenF64, In(AI32), Out(AF64, ABOOL))
	Op_V2(OP_JSON_TOKEN_I64, "json.Int64", opJsonTokenI64, In(AI32), Out(AI64, ABOOL))
	Op_V2(OP_JSON_TOKEN_STR, "json.Str", opJsonTokenStr, In(AI32), Out(ASTR, ABOOL))

	// profile
	Op_V2(OP_START_CPU_PROFILE, "StartCPUProfile", opStartProfile, In(ASTR, AI32), nil)
	Op_V2(OP_STOP_CPU_PROFILE, "StopCPUProfile", opStopProfile, In(ASTR), nil)

	// regexp
	Op(OP_REGEXP_COMPILE, "regexp.Compile", opRegexpCompile, In(ASTR), Out(Struct("regexp", "Regexp", "r"), ASTR))
	Op(OP_REGEXP_MUST_COMPILE, "regexp.MustCompile", opRegexpMustCompile, In(ASTR), Out(Struct("regexp", "Regexp", "r")))
	Op(OP_REGEXP_FIND, "regexp.Regexp.Find", opRegexpFind, In(Struct("regexp", "Regexp", "r"), ASTR), Out(ASTR))

	// cipher
	Op_V2(OP_CIPHER_GENERATE_KEY_PAIR, "cipher.GenerateKeyPair", opCipherGenerateKeyPair, nil, Out(Struct("cipher", "PubKey", "pubKey"), Struct("cipher", "SecKey", "sec")))
}
