package webapi

import (
	"unicode"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/types"
)

// ProgramMetaResp is a program meta data response.
type ProgramMetaResp struct {
	UsedHeapMemory types.Pointer `json:"used_heap_memory"`
	FreeHeapMemory types.Pointer `json:"free_heap_memory"`
	StackSize      types.Pointer `json:"stack_size"`
	CallStackSize  int           `json:"call_stack_size"`
}

func extractProgramMeta(pg *ast.CXProgram) ProgramMetaResp {
	return ProgramMetaResp{
		UsedHeapMemory: pg.Heap.Pointer - pg.Heap.StartsAt,
		FreeHeapMemory: pg.Heap.Size - pg.Heap.StartsAt,
		StackSize:      pg.Stack.Size,
		CallStackSize:  len(pg.CallStack),
	}
}

// ExportedSymbol represents a single exported cx symbol.
type ExportedSymbol struct {
	Name      string      `json:"name"`
	Signature interface{} `json:"signature,omitempty"`
	Type      types.Code  `json:"type"`
	TypeName  string      `json:"type_name"`
}

// ExportedSymbolsResp is the exported symbols response.
type ExportedSymbolsResp struct {
	Functions []ExportedSymbol `json:"functions"`
	Structs   []ExportedSymbol `json:"structs"`
	Globals   []ExportedSymbol `json:"globals"`
}

func extractExportedSymbols(pkg *ast.CXPackage) ExportedSymbolsResp {
	resp := ExportedSymbolsResp{
		Functions: make([]ExportedSymbol, 0, len(pkg.Functions)),
		Structs:   make([]ExportedSymbol, 0, len(pkg.Structs)),
		Globals:   make([]ExportedSymbol, 0, len(pkg.Globals)),
	}

	for _, f := range pkg.Functions {
		if isExported(f.Name) {
			resp.Functions = append(resp.Functions, displayCXFunction(pkg, f))
		}
	}

	for _, s := range pkg.Structs {
		if isExported(s.Name) {
			resp.Structs = append(resp.Structs, displayCXStruct(s))
		}
	}

	for _, g := range pkg.Globals {
		if isExported(g.ArgDetails.Name) {
			resp.Globals = append(resp.Globals, displayCXGlobal(g))
		}
	}

	return resp
}

func displayCXFunction(pkg *ast.CXPackage, f *ast.CXFunction) ExportedSymbol {
	return ExportedSymbol{
		Name:      f.Name,
		Signature: ast.SignatureStringOfFunction(pkg, f),
		Type:      types.FUNC,
		TypeName:  types.FUNC.Name(),
	}
}

func displayCXStruct(s *ast.CXStruct) ExportedSymbol {
	return ExportedSymbol{
		Name:      s.Name,
		Signature: ast.SignatureStringOfStruct(s),
		Type:      types.STRUCT,
		TypeName:  "struct",
	}
}

func displayCXGlobal(a *ast.CXArgument) ExportedSymbol {
	return ExportedSymbol{
		Name:      a.ArgDetails.Name,
		Signature: nil,
		Type:      a.Type,
		TypeName:  a.Type.Name(),
	}
}

func isExported(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsUpper(rune(s[0]))
}
