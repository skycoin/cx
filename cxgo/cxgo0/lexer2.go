package cxgo0

import (
	"fmt"
	"io"
)

var CurrentFileName string

func (yylex Lexer) Error(e string) {
	if inREPL {
		fmt.Printf("syntax error: %s\n", e)
	} else {
		fmt.Printf("%s:%d: syntax error: %s\n", CurrentFileName, yylex.l+1, e)
	}

	yylex.stop()
}

func (yylex *Lexer) Stop() {
	yylex.stop()
}

func (yylex *Lexer) Lex(lval *yySymType) int {
	yylex.next()
	lval.scancopy(yylex.tok)
	lineNo = lval.line
	return lval.yys
}

func NewLexer(rdr io.Reader) *Lexer {
	lx := &Lexer{}
	lx.init(rdr, func(l, c int, msg string) {
		fmt.Printf("[%d:%d] %s\n", l, c, msg)
	})
	return lx
}

func (lval *yySymType) scancopy(tok *yySymType) {
	lval.argument = tok.argument
	lval.arguments = tok.arguments
	lval.bool = tok.bool
	lval.expression = tok.expression
	lval.expressions = tok.expressions
	lval.f32 = tok.f32
	lval.f64 = tok.f64
	lval.function = tok.function
	lval.i = tok.i
	lval.i16 = tok.i16
	lval.i32 = tok.i32
	lval.i64 = tok.i64
	lval.i8 = tok.i8
	lval.ints = tok.ints
	lval.line = tok.line
	lval.string = tok.string
	lval.stringA = tok.stringA
	lval.tok = tok.tok
	lval.ui16 = tok.ui16
	lval.ui32 = tok.ui32
	lval.ui64 = tok.ui64
	lval.ui8 = tok.ui8
	lval.yys = tok.yys
}

//Unused
//looks like a copy of cxparser/cxgo/lexer.go?
//should definatel be in constants or something
/*
var tokenNames = map[int]string{
	BYTE_LITERAL:    "BYTE_LITERAL",
	BOOLEAN_LITERAL: "BOOLEAN_LITERAL",
	INT_LITERAL:     "INT_LITERAL",
	LONG_LITERAL:    "LONG_LITERAL",
	FLOAT_LITERAL:   "FLOAT_LITERAL",
	DOUBLE_LITERAL:  "DOUBLE_LITERAL",
	FUNC:            "FUNC",
	OP:              "OP",
	LPAREN:          "LPAREN",
	RPAREN:          "RPAREN",
	LBRACE:          "LBRACE",
	RBRACE:          "RBRACE",
	LBRACK:          "LBRACK",
	RBRACK:          "RBRACK",
	IDENTIFIER:      "IDENTIFIER",
	VAR:             "VAR",
	COMMA:           "COMMA",
	PERIOD:          "PERIOD",
	COMMENT:         "COMMENT",
	STRING_LITERAL:  "STRING_LITERAL",
	PACKAGE:         "PACKAGE",
	IF:              "IF",
	ELSE:            "ELSE",
	FOR:             "FOR",
	TYPSTRUCT:       "TYPSTRUCT",
	STRUCT:          "STRUCT",
	SEMICOLON:       "SEMICOLON",
	NEWLINE:         "NEWLINE",
	ASSIGN:          "ASSIGN",
	CASSIGN:         "CASSIGN",
	IMPORT:          "IMPORT",
	RETURN:          "RETURN",
	GOTO:            "GOTO",
	GT_OP:           "GT_OP",
	LT_OP:           "LT_OP",
	GTEQ_OP:         "GTEQ_OP",
	LTEQ_OP:         "LTEQ_OP",
	EQUAL:           "EQUAL",
	COLON:           "COLON",
	NEW:             "NEW",
	EQUALWORD:       "EQUALWORD",
	GTHANWORD:       "GTHANWORD",
	LTHANWORD:       "LTHANWORD",
	GTHANEQ:         "GTHANEQ",
	LTHANEQ:         "LTHANEQ",
	UNEQUAL:         "UNEQUAL",
	AND:             "AND",
	OR:              "OR",
	ADD_OP:          "ADD_OP",
	SUB_OP:          "SUB_OP",
	MUL_OP:          "MUL_OP",
	DIV_OP:          "DIV_OP",
	MOD_OP:          "MOD_OP",
	REF_OP:          "REF_OP",
	NEG_OP:          "NEG_OP",
	AFFVAR:          "AFFVAR",
	PLUSPLUS:        "PLUSPLUS",
	MINUSMINUS:      "MINUSMINUS",
	REMAINDER:       "REMAINDER",
	LEFTSHIFT:       "LEFTSHIFT",
	RIGHTSHIFT:      "RIGHTSHIFT",
	EXP:             "EXP",
	NOT:             "NOT",
	BITXOR_OP:       "BITXOR_OP",
	BITOR_OP:        "BITOR_OP",
	BITCLEAR_OP:     "BITCLEAR_OP",
	PLUSEQ:          "PLUSEQ",
	MINUSEQ:         "MINUSEQ",
	MULTEQ:          "MULTEQ",
	DIVEQ:           "DIVEQ",
	REMAINDEREQ:     "REMAINDEREQ",
	EXPEQ:           "EXPEQ",
	LEFTSHIFTEQ:     "LEFTSHIFTEQ",
	RIGHTSHIFTEQ:    "RIGHTSHIFTEQ",
	BITANDEQ:        "BITANDEQ",
	BITXOREQ:        "BITXOREQ",
	BITOREQ:         "BITOREQ",

	DEC_OP:       "DEC_OP",
	INC_OP:       "INC_OP",
	PTR_OP:       "PTR_OP",
	LEFT_OP:      "LEFT_OP",
	RIGHT_OP:     "RIGHT_OP",
	GE_OP:        "GE_OP",
	LE_OP:        "LE_OP",
	EQ_OP:        "EQ_OP",
	NE_OP:        "NE_OP",
	AND_OP:       "AND_OP",
	OR_OP:        "OR_OP",
	ADD_ASSIGN:   "ADD_ASSIGN",
	AND_ASSIGN:   "AND_ASSIGN",
	LEFT_ASSIGN:  "LEFT_ASSIGN",
	MOD_ASSIGN:   "MOD_ASSIGN",
	MUL_ASSIGN:   "MUL_ASSIGN",
	DIV_ASSIGN:   "DIV_ASSIGN",
	OR_ASSIGN:    "OR_ASSIGN",
	RIGHT_ASSIGN: "RIGHT_ASSIGN",
	SUB_ASSIGN:   "SUB_ASSIGN",
	XOR_ASSIGN:   "XOR_ASSIGN",
	BOOL:         "BOOL",
	F32:          "F32",
	F64:          "F64",
	I8:           "I8",
	I16:          "I16",
	I32:          "I32",
	I64:          "I64",
	STR:          "STR",
	UI8:          "UI8",
	UI16:         "UI16",
	UI32:         "UI32",
	UI64:         "UI64",
	UNION:        "UNION",
	ENUM:         "ENUM",
	CONST:        "CONST",
	CASE:         "CASE",
	DEFAULT:      "DEFAULT",
	SWITCH:       "SWITCH",
	BREAK:        "BREAK",
	CONTINUE:     "CONTINUE",
	TYPE:         "TYPE",

	// Types 
	BASICTYPE: "BASICTYPE",
	// Selectors 
	SPACKAGE: "SPACKAGE",
	SSTRUCT:  "SSTRUCT",
	SFUNC:    "SFUNC",
	// Removers 
	REM:     "REM",
	DEF:     "DEF",
	EXPR:    "EXPR",
	FIELD:   "FIELD",
	CLAUSES: "CLAUSES",
	OBJECT:  "OBJECT",
	OBJECTS: "OBJECTS",
	// Stepping 
	STEP:  "STEP",
	PSTEP: "PSTEP",
	TSTEP: "TSTEP",
	// Debugging 
	DSTACK:   "DSTACK",
	DPROGRAM: "DPROGRAM",
	DSTATE:   "DSTATE",
	// Affordances 
	AFF:   "AFF",
	CAFF:  "CAFF",
	TAG:   "TAG",
	INFER: "INFER",
	VALUE: "VALUE",
	// Pointers 
	ADDR: "ADDR",
}
*/