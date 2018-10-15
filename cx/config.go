package base

import (
	"os"
)

// global reference to our program
var PROGRAM *CXProgram

var CXPATH string = os.Getenv("CXPATH") + "/"
var BINPATH string = CXPATH + "bin/"
var PKGPATH string = CXPATH + "pkg/"
var SRCPATH string = CXPATH + "src/"
var COREPATH string

const MAIN_FUNC = "main"
const SYS_INIT_FUNC = "*init"
const MAIN_PKG = "main"
const NON_ASSIGN_PREFIX = "nonAssign"
const LOCAL_PREFIX = "*lcl"
const LABEL_PREFIX = "*lbl"
const CORE_MODULE = "core"
const ID_FN = "identity"
const INIT_FN = "initDef"
const SLICE_SIZE = 32
const MARK_SIZE = 1
const OBJECT_HEADER_SIZE = 9
const FORWARDING_ADDRESS_SIZE = 4
const OBJECT_SIZE = 4
const CALLSTACK_SIZE = 1000
const STACK_SIZE = 500000
const INIT_HEAP_SIZE = 500000
const NULL_ADDRESS = STACK_SIZE
const NULL_HEAP_ADDRESS_OFFSET = 4
const NULL_HEAP_ADDRESS = 0
const STR_HEADER_SIZE = 4
const TYPE_POINTER_SIZE = 4
const SLICE_HEADER_SIZE = 8
const MEMORY_SIZE = STACK_SIZE + INIT_HEAP_SIZE + TYPE_POINTER_SIZE

const MAX_UINT32 = ^uint32(0)
const MIN_UINT32 = 0
const MAX_INT32 = int(MAX_UINT32 >> 1)
const MIN_INT32 = -MAX_INT32 - 1

var BASIC_TYPES []string = []string{
	"bool", "str", "byte", "i32", "i64", "f32", "f64",
	"[]bool", "[]str", "[]byte", "[]i32", "[]i64", "[]f32", "[]f64",
}

const (
	DECL_POINTER = iota // 0
	DECL_ARRAY          // 1
	DECL_SLICE          // 2
	DECL_STRUCT         // 3
	DECL_BASIC          // 4
)

// what to write
const (
	PASSBY_VALUE = iota
	PASSBY_REFERENCE
)

const (
	DEREF_ARRAY = iota
	DEREF_FIELD
	DEREF_POINTER
	DEREF_DEREF
)

const (
	TYPE_AFF = iota
	TYPE_BOOL
	TYPE_BYTE
	TYPE_STR
	TYPE_F32
	TYPE_F64
	TYPE_I8
	TYPE_I16
	TYPE_I32
	TYPE_I64
	TYPE_UI8
	TYPE_UI16
	TYPE_UI32
	TYPE_UI64

	TYPE_THRESHOLD

	TYPE_UNDEFINED
	TYPE_CUSTOM
	TYPE_POINTER
	TYPE_IDENTIFIER
)

var TypeCounter int
var TypeCodes map[string]int = map[string]int{
	"identifier": TYPE_IDENTIFIER,
	"aff":        TYPE_AFF,
	"bool":       TYPE_BOOL,
	"byte":       TYPE_BYTE,
	"str":        TYPE_STR,
	"f32":        TYPE_F32,
	"f64":        TYPE_F64,
	"i8":         TYPE_I8,
	"i16":        TYPE_I16,
	"i32":        TYPE_I32,
	"i64":        TYPE_I64,
	"ui8":        TYPE_UI8,
	"ui16":       TYPE_UI16,
	"ui32":       TYPE_UI32,
	"ui64":       TYPE_UI64,
	"und":        TYPE_UNDEFINED,
}

var TypeNames map[int]string = map[int]string{
	TYPE_IDENTIFIER: "ident",
	TYPE_AFF:        "aff",
	TYPE_BOOL:       "bool",
	TYPE_BYTE:       "byte",
	TYPE_STR:        "str",
	TYPE_F32:        "f32",
	TYPE_F64:        "f64",
	TYPE_I8:         "i8",
	TYPE_I16:        "i16",
	TYPE_I32:        "i32",
	TYPE_I64:        "i64",
	TYPE_UI8:        "ui8",
	TYPE_UI16:       "ui16",
	TYPE_UI32:       "ui32",
	TYPE_UI64:       "ui64",
	TYPE_UNDEFINED:  "und",
}

// memory locations
const (
	MEM_STACK = iota
	MEM_HEAP
	MEM_DATA
)
