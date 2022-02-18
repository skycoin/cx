// +build cxos

package cxos

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/opcodes"
	"github.com/skycoin/cx/cx/types"
)

func RegisterPackage(prgrm *ast.CXProgram) {
	// time
	opcodes.RegisterFunction(prgrm, "time.Sleep", opTimeSleep, opcodes.In(ast.ConstCxArg_I32), nil)
	opcodes.RegisterFunction(prgrm, "time.UnixMilli", opTimeUnixMilli, nil, opcodes.Out(ast.ConstCxArg_I64))
	opcodes.RegisterFunction(prgrm, "time.UnixNano", opTimeUnixNano, nil, opcodes.Out(ast.ConstCxArg_I64))

	// os
	opcodes.RegisterFunction(prgrm, "os.GetWorkingDirectory", opOsGetWorkingDirectory, nil, opcodes.Out(ast.ConstCxArg_STR))
	opcodes.RegisterFunction(prgrm, "os.LogFile", opOsLogFile, opcodes.In(ast.ConstCxArg_BOOL), nil)
	opcodes.RegisterFunction(prgrm, "os.Open", opOsOpen, opcodes.In(ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_I32))
	opcodes.RegisterFunction(prgrm, "os.Create", opOsCreate, opcodes.In(ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_I32))
	opcodes.RegisterFunction(prgrm, "os.Close", opOsClose, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.Seek", opOsSeek, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_I64, ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I64))
	opcodes.RegisterFunction(prgrm, "os.ReadStr", opOsReadStr, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_STR, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadF64", opOsReadF64, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_F64, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadF32", opOsReadF32, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_F32, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadUI64", opOsReadUI64, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_UI64, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadUI32", opOsReadUI32, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_UI32, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadUI16", opOsReadUI16, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_UI16, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadUI8", opOsReadUI8, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_UI8, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadI64", opOsReadI64, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I64, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadI32", opOsReadI32, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadI16", opOsReadI16, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I16, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadI8", opOsReadI8, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCXArg_I8, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadBOOL", opOsReadBOOL, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_BOOL, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadF64Slice", opOsReadF64Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.F64), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(types.F64), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadF32Slice", opOsReadF32Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.F32), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(types.F32), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadUI64Slice", opOsReadUI64Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.UI64), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(types.UI64), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadUI32Slice", opOsReadUI32Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.UI32), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(types.UI32), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadUI16Slice", opOsReadUI16Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.UI16), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(types.UI16), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadUI8Slice", opOsReadUI8Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.UI8), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(types.UI8), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadI64Slice", opOsReadI64Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.I64), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(types.I64), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadI32Slice", opOsReadI32Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.I32), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(types.I32), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadI16Slice", opOsReadI16Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.I16), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(types.I16), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadI8Slice", opOsReadI8Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.I8), ast.ConstCxArg_UI64), opcodes.Out(ast.Slice(types.I8), ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.ReadAllText", opOsReadAllText, opcodes.In(ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_STR, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteStr", opOsWriteStr, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteF64", opOsWriteF64, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_F64), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteF32", opOsWriteF32, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_F32), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteUI64", opOsWriteUI64, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_UI64), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteUI32", opOsWriteUI32, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_UI32), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteUI16", opOsWriteUI16, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_UI16), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteUI8", opOsWriteUI8, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_UI8), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteI64", opOsWriteI64, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_I64), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteI32", opOsWriteI32, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteI16", opOsWriteI16, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_I16), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteI8", opOsWriteI8, opcodes.In(ast.ConstCxArg_I32, ast.ConstCXArg_I8), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteBOOL", opOsWriteBOOL, opcodes.In(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteF64Slice", opOsWriteF64Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.F64)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteF32Slice", opOsWriteF32Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.F32)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteUI64Slice", opOsWriteUI64Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.UI64)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteUI32Slice", opOsWriteUI32Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.UI32)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteUI16Slice", opOsWriteUI16Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.UI16)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteUI8Slice", opOsWriteUI8Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.UI8)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteI64Slice", opOsWriteI64Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.I64)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteI32Slice", opOsWriteI32Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.I32)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteI16Slice", opOsWriteI16Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.I16)), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "os.WriteI8Slice", opOsWriteI8Slice, opcodes.In(ast.ConstCxArg_I32, ast.Slice(types.I8)), opcodes.Out(ast.ConstCxArg_BOOL))

	opcodes.RegisterFunction(prgrm, "os.Run", opOsRun, opcodes.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_I32, ast.ConstCxArg_I32, ast.ConstCxArg_STR))
	opcodes.RegisterFunction(prgrm, "os.Exit", opOsExit, opcodes.In(ast.ConstCxArg_I32), nil)

	// json
	opcodes.RegisterFunction(prgrm, "json.Open", opJsonOpen, opcodes.In(ast.ConstCxArg_STR), opcodes.Out(ast.ConstCxArg_I32))
	opcodes.RegisterFunction(prgrm, "json.Close", opJsonClose, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "json.More", opJsonTokenMore, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_BOOL, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "json.Next", opJsonTokenNext, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "json.Type", opJsonTokenType, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "json.Delim", opJsonTokenDelim, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I32, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "json.Bool", opJsonTokenBool, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_BOOL, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "json.Float64", opJsonTokenF64, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_F64, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "json.Int64", opJsonTokenI64, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_I64, ast.ConstCxArg_BOOL))
	opcodes.RegisterFunction(prgrm, "json.Str", opJsonTokenStr, opcodes.In(ast.ConstCxArg_I32), opcodes.Out(ast.ConstCxArg_STR, ast.ConstCxArg_BOOL))

	// profile
	opcodes.RegisterFunction(prgrm, "StartCPUProfile", opStartProfile, opcodes.In(ast.ConstCxArg_STR, ast.ConstCxArg_I32), nil)
	opcodes.RegisterFunction(prgrm, "StopCPUProfile", opStopProfile, opcodes.In(ast.ConstCxArg_STR), nil)
}
