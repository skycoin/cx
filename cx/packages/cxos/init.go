// +build cxos

package cxos

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/opcodes"
)


func RegisterPackage() {
	// time
	opcodes.RegisterFunction("time.Sleep", opTimeSleep, opcodes.In(ast.ConstCxArg_I32), nil)
	opcodes.RegisterFunction("time.UnixMilli", opTimeUnixMilli, nil, opcodes.Out(ast.ConstCxArg_I64))
	opcodes.RegisterFunction("time.UnixNano", opTimeUnixNano, nil, opcodes.Out(ast.ConstCxArg_I64))

	// os
	opcodes.RegisterFunction("os.GetWorkingDirectory", opOsGetWorkingDirectory, nil, opcodes.Out(ast.ConstCxArg_STR))
	opcodes.RegisterFunction("os.LogFile", opOsLogFile, opcodes.In(ast.ConstCxArg_BOOL), nil)
	opcodes.RegisterFunction("os.Open", opOsOpen, opcodes.In(ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_I32))
	opcodes.RegisterFunction("os.Create", opOsCreate, opcodes.In(ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_I32))
	opcodes.RegisterFunction("os.Close", opOsClose, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.Seek", opOsSeek, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_I64, ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I64))
	opcodes.RegisterFunction("os.ReadStr", opOsReadStr, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_STR, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadF64", opOsReadF64, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_F64, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadF32", opOsReadF32, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_F32, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadUI64", opOsReadUI64, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_UI64, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadUI32", opOsReadUI32, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_UI32, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadUI16", opOsReadUI16, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_UI16, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadUI8", opOsReadUI8, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_UI8, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadI64", opOsReadI64, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I64, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadI32", opOsReadI32, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadI16", opOsReadI16, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I16, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadI8", opOsReadI8, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCXArg_I8, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadBOOL", opOsReadBOOL, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_BOOL, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadF64Slice", opOsReadF64Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_F64), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(constants.TYPE_F64), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadF32Slice", opOsReadF32Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_F32), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(constants.TYPE_F32), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadUI64Slice", opOsReadUI64Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI64), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(constants.TYPE_UI64), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadUI32Slice", opOsReadUI32Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI32), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(constants.TYPE_UI32), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadUI16Slice", opOsReadUI16Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI16), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(constants.TYPE_UI16), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadUI8Slice", opOsReadUI8Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI8), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(constants.TYPE_UI8), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadI64Slice", opOsReadI64Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I64), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(constants.TYPE_I64), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadI32Slice", opOsReadI32Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I32), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(constants.TYPE_I32), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadI16Slice", opOsReadI16Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I16), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(constants.TYPE_I16), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadI8Slice", opOsReadI8Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I8), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(constants.TYPE_I8), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.ReadAllText", opOsReadAllText, opcodes.In(ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_STR, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteStr", opOsWriteStr, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteF64", opOsWriteF64, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_F64), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteF32", opOsWriteF32, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_F32), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteUI64", opOsWriteUI64, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_UI64), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteUI32", opOsWriteUI32, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_UI32), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteUI16", opOsWriteUI16, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_UI16), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteUI8", opOsWriteUI8, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_UI8), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteI64", opOsWriteI64, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_I64), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteI32", opOsWriteI32, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteI16", opOsWriteI16, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_I16), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteI8", opOsWriteI8, opcodes.In(ast.ConstCxArg_I32, ast.ConstCXArg_I8), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteBOOL", opOsWriteBOOL, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteF64Slice", opOsWriteF64Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_F64)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteF32Slice", opOsWriteF32Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_F32)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteUI64Slice", opOsWriteUI64Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI64)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteUI32Slice", opOsWriteUI32Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI32)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteUI16Slice", opOsWriteUI16Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI16)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteUI8Slice", opOsWriteUI8Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_UI8)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteI64Slice", opOsWriteI64Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I64)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteI32Slice", opOsWriteI32Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I32)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteI16Slice", opOsWriteI16Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I16)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("os.WriteI8Slice", opOsWriteI8Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(constants.TYPE_I8)), opcodes.Out(ast.ConstCxArg_BOOL))

	opcodes.RegisterFunction("os.Run", opOsRun, opcodes.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_STR))
	opcodes.RegisterFunction("os.Exit", opOsExit, opcodes.In(ast.ConstCxArg_I32), nil)

	// json
	opcodes.RegisterFunction("json.Open", opJsonOpen, opcodes.In(ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_I32))
	opcodes.RegisterFunction("json.Close", opJsonClose, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("json.More", opJsonTokenMore, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_BOOL, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("json.Next", opJsonTokenNext, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("json.Type", opJsonTokenType, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("json.Delim", opJsonTokenDelim, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("json.Bool", opJsonTokenBool, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_BOOL, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("json.Float64", opJsonTokenF64, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_F64, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("json.Int64", opJsonTokenI64, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I64, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction("json.Str", opJsonTokenStr, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_STR, ast.ConstCxArg_BOOL))

	// profile
	opcodes.RegisterFunction("StartCPUProfile", opStartProfile, opcodes.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32), nil)
	opcodes.RegisterFunction("StopCPUProfile", opStopProfile, opcodes.In(ast.ConstCxArg_STR), nil)
}
