// +build os

package cxos

import (
	. "github.com/skycoin/cx/cx"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/packages/http"
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
	RegisterOpCode(OP_TIME_SLEEP, "time.Sleep", opTimeSleep, http.In(ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_TIME_UNIX_MILLI, "time.UnixMilli", opTimeUnixMilli, nil, http.Out(ast.ConstCxArg_I64))
	RegisterOpCode(OP_TIME_UNIX_NANO, "time.UnixNano", opTimeUnixNano, nil, http.Out(ast.ConstCxArg_I64))

	// http
	// RegisterOpCode(OP_HTTP_GET, "http.Get", opHttpGet, In(ConstCxArg_STR), Out(ConstCxArg_STR))

	// os
	RegisterOpCode(OP_OS_GET_WORKING_DIRECTORY, "os.GetWorkingDirectory", opOsGetWorkingDirectory, nil, http.Out(ast.ConstCxArg_STR))
	RegisterOpCode(OP_OS_LOG_FILE, "os.LogFile", opOsLogFile, http.In(ast.ConstCxArg_BOOL), nil)
	RegisterOpCode(OP_OS_OPEN, "os.Open", opOsOpen, http.In(ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_OS_CREATE, "os.Create", opOsCreate, http.In(ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_OS_CLOSE, "os.Close", opOsClose, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_SEEK, "os.Seek", opOsSeek, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I64, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I64))
	RegisterOpCode(OP_OS_READ_STR, "os.ReadStr", opOsReadStr, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_STR, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_F64, "os.ReadF64", opOsReadF64, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_F64, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_F32, "os.ReadF32", opOsReadF32, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_F32, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_UI64, "os.ReadUI64", opOsReadUI64, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_UI64, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_UI32, "os.ReadUI32", opOsReadUI32, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_UI32, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_UI16, "os.ReadUI16", opOsReadUI16, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_UI16, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_UI8, "os.ReadUI8", opOsReadUI8, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_UI8, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_I64, "os.ReadI64", opOsReadI64, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I64, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_I32, "os.ReadI32", opOsReadI32, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_I16, "os.ReadI16", opOsReadI16, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I16, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_I8, "os.ReadI8", opOsReadI8, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCXArg_I8, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_BOOL, "os.ReadBOOL", opOsReadBOOL, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_BOOL, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_F64_SLICE, "os.ReadF64Slice", opOsReadF64Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_F64), ast.ConstCxArg_UI64), http.Out(ast.Slice(constants.TYPE_F64), ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_F32_SLICE, "os.ReadF32Slice", opOsReadF32Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_F32), ast.ConstCxArg_UI64), http.Out(ast.Slice(constants.TYPE_F32), ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_UI64_SLICE, "os.ReadUI64Slice", opOsReadUI64Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI64), ast.ConstCxArg_UI64), http.Out(ast.Slice(constants.TYPE_UI64), ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_UI32_SLICE, "os.ReadUI32Slice", opOsReadUI32Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI32), ast.ConstCxArg_UI64), http.Out(ast.Slice(constants.TYPE_UI32), ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_UI16_SLICE, "os.ReadUI16Slice", opOsReadUI16Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI16), ast.ConstCxArg_UI64), http.Out(ast.Slice(constants.TYPE_UI16), ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_UI8_SLICE, "os.ReadUI8Slice", opOsReadUI8Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI8), ast.ConstCxArg_UI64), http.Out(ast.Slice(constants.TYPE_UI8), ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_I64_SLICE, "os.ReadI64Slice", opOsReadI64Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I64), ast.ConstCxArg_UI64), http.Out(ast.Slice(constants.TYPE_I64), ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_I32_SLICE, "os.ReadI32Slice", opOsReadI32Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I32), ast.ConstCxArg_UI64), http.Out(ast.Slice(constants.TYPE_I32), ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_I16_SLICE, "os.ReadI16Slice", opOsReadI16Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I16), ast.ConstCxArg_UI64), http.Out(ast.Slice(constants.TYPE_I16), ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_I8_SLICE, "os.ReadI8Slice", opOsReadI8Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I8), ast.ConstCxArg_UI64), http.Out(ast.Slice(constants.TYPE_I8), ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_READ_ALL_TEXT, "os.ReadAllText", opOsReadAllText, http.In(ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_STR, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_STR, "os.WriteStr", opOsWriteStr, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_F64, "os.WriteF64", opOsWriteF64, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_F64), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_F32, "os.WriteF32", opOsWriteF32, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_F32), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_UI64, "os.WriteUI64", opOsWriteUI64, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_UI64), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_UI32, "os.WriteUI32", opOsWriteUI32, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_UI32), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_UI16, "os.WriteUI16", opOsWriteUI16, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_UI16), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_UI8, "os.WriteUI8", opOsWriteUI8, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_UI8), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_I64, "os.WriteI64", opOsWriteI64, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I64), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_I32, "os.WriteI32", opOsWriteI32, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_I16, "os.WriteI16", opOsWriteI16, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_I16), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_I8, "os.WriteI8", opOsWriteI8, http.In(ast.ConstCxArg_I32, ast.ConstCXArg_I8), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_BOOL, "os.WriteBOOL", opOsWriteBOOL, http.In(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_F64_SLICE, "os.WriteF64Slice", opOsWriteF64Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_F64)), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_F32_SLICE, "os.WriteF32Slice", opOsWriteF32Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_F32)), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_UI64_SLICE, "os.WriteUI64Slice", opOsWriteUI64Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI64)), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_UI32_SLICE, "os.WriteUI32Slice", opOsWriteUI32Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI32)), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_UI16_SLICE, "os.WriteUI16Slice", opOsWriteUI16Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI16)), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_UI8_SLICE, "os.WriteUI8Slice", opOsWriteUI8Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI8)), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_I64_SLICE, "os.WriteI64Slice", opOsWriteI64Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I64)), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_I32_SLICE, "os.WriteI32Slice", opOsWriteI32Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I32)), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_I16_SLICE, "os.WriteI16Slice", opOsWriteI16Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I16)), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_OS_WRITE_I8_SLICE, "os.WriteI8Slice", opOsWriteI8Slice, http.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I8)), http.Out(ast.ConstCxArg_BOOL))

	RegisterOpCode(OP_OS_RUN, "os.Run", opOsRun, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_STR))
	RegisterOpCode(OP_OS_EXIT, "os.Exit", opOsExit, http.In(ast.ConstCxArg_I32), nil)

	// json
	RegisterOpCode(OP_JSON_OPEN, "json.Open", opJsonOpen, http.In(ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_I32))
	RegisterOpCode(OP_JSON_CLOSE, "json.Close", opJsonClose, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_JSON_TOKEN_MORE, "json.More", opJsonTokenMore, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_BOOL, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_JSON_TOKEN_NEXT, "json.Next", opJsonTokenNext, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_JSON_TOKEN_TYPE, "json.Type", opJsonTokenType, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_JSON_TOKEN_DELIM, "json.Delim", opJsonTokenDelim, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_JSON_TOKEN_BOOL, "json.Bool", opJsonTokenBool, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_BOOL, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_JSON_TOKEN_F64, "json.Float64", opJsonTokenF64, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_F64, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_JSON_TOKEN_I64, "json.Int64", opJsonTokenI64, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_I64, ast.ConstCxArg_BOOL))
	RegisterOpCode(OP_JSON_TOKEN_STR, "json.Str", opJsonTokenStr, http.In(ast.ConstCxArg_I32), http.Out(ast.ConstCxArg_STR, ast.ConstCxArg_BOOL))

	// profile
	RegisterOpCode(OP_START_CPU_PROFILE, "StartCPUProfile", opStartProfile, http.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32), nil)
	RegisterOpCode(OP_STOP_CPU_PROFILE, "StopCPUProfile", opStopProfile, http.In(ast.ConstCxArg_STR), nil)

	// regexp
	RegisterOpCode(OP_REGEXP_COMPILE, "regexp.Compile", opRegexpCompile, http.In(ast.ConstCxArg_STR), http.Out(Struct("regexp", "Regexp", "r"), ast.ConstCxArg_STR))
	RegisterOpCode(OP_REGEXP_MUST_COMPILE, "regexp.MustCompile", opRegexpMustCompile, http.In(ast.ConstCxArg_STR), http.Out(Struct("regexp", "Regexp", "r")))
	RegisterOpCode(OP_REGEXP_FIND, "regexp.Regexp.Find", opRegexpFind, http.In(Struct("regexp", "Regexp", "r"), ast.ConstCxArg_STR), http.Out(ast.ConstCxArg_STR))

	// cipher
	RegisterOpCode(OP_CIPHER_GENERATE_KEY_PAIR, "cipher.GenerateKeyPair", opCipherGenerateKeyPair, nil, http.Out(Struct("cipher", "PubKey", "pubKey"), Struct("cipher", "SecKey", "sec")))
}
