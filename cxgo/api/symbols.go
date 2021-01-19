package api

import (
	"unicode"

	cxcore "github.com/skycoin/cx/cx"
)

// ProgramMetaResp is a program meta data response.
type ProgramMetaResp struct {
	UsedHeapMemory int `json:"used_heap_memory"`
	FreeHeapMemory int `json:"free_heap_memory"`
	StackSize      int `json:"stack_size"`
	CallStackSize  int `json:"call_stack_size"`
}

func extractProgramMeta(pg *cxcore.CXProgram) ProgramMetaResp {
	return ProgramMetaResp{
		UsedHeapMemory: pg.HeapPointer - pg.HeapStartsAt,
		FreeHeapMemory: pg.HeapSize - pg.HeapStartsAt,
		StackSize:      pg.StackSize,
		CallStackSize:  len(pg.CallStack),
	}
}

// ExportedSymbol represents a single exported cx symbol.
type ExportedSymbol struct {
	Name      string      `json:"name"`
	Signature interface{} `json:"signature,omitempty"`
	Type      int         `json:"type"`
	TypeName  string      `json:"type_name"`
}

// ExportedSymbolsResp is the exported symbols response.
type ExportedSymbolsResp struct {
	Functions []ExportedSymbol `json:"functions"`
	Structs   []ExportedSymbol `json:"structs"`
	Globals   []ExportedSymbol `json:"globals"`
}

func extractExportedSymbols(pkg *cxcore.CXPackage) ExportedSymbolsResp {
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
		if isExported(g.Name) {
			resp.Globals = append(resp.Globals, displayCXGlobal(g))
		}
	}

	return resp
}

func displayCXFunction(pkg *cxcore.CXPackage, f *cxcore.CXFunction) ExportedSymbol {
	return ExportedSymbol{
		Name:      f.Name,
		Signature: cxcore.SignatureStringOfFunction(pkg, f),
		Type:      cxcore.TYPE_FUNC,
		TypeName:  cxcore.TypeNames[cxcore.TYPE_FUNC],
	}
}

func displayCXStruct(s *cxcore.CXStruct) ExportedSymbol {
	return ExportedSymbol{
		Name:      s.Name,
		Signature: cxcore.SignatureStringOfStruct(s),
		Type:      cxcore.TYPE_CUSTOM,
		TypeName:  "struct",
	}
}

func displayCXGlobal(a *cxcore.CXArgument) ExportedSymbol {
	return ExportedSymbol{
		Name:      a.Name,
		Signature: nil,
		Type:      a.Type,
		TypeName:  cxcore.TypeNames[a.Type],
	}
}

func isExported(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsUpper(rune(s[0]))
}
