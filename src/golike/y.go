// CAUTION: Generated file - DO NOT EDIT.

package main

import __yyfmt__ "fmt"

import (
	"fmt"
	. "github.com/skycoin/cx/src/base"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"strings"
)

var cxt = MakeContext()
var lineNo int = 0

type yySymType struct {
	yys int
	i64 int64
	f64 float64

	tok string

	param  *CXParameter
	params []*CXParameter

	arg  *CXArgument
	args []*CXArgument

	expr  *CXExpression
	exprs []*CXExpression

	fld  *CXField
	flds []*CXField

	names []string
}

type yyXError struct {
	state, xsym int
}

const (
	yyDefault  = 57389
	yyEofCode  = 57344
	ASSIGN     = 57367
	BOOLEAN    = 57348
	BYTE       = 57361
	BYTEA      = 57362
	CASSIGN    = 57368
	COMMA      = 57374
	COMMENT    = 57375
	ELSE       = 57379
	F32        = 57359
	F32A       = 57365
	F64        = 57360
	F64A       = 57366
	FLOAT      = 57347
	FUNC       = 57349
	GTEQ       = 57372
	GTHAN      = 57369
	I32        = 57357
	I32A       = 57363
	I64        = 57358
	I64A       = 57364
	IDENT      = 57355
	IF         = 57378
	INT        = 57346
	LBRACE     = 57353
	LPAREN     = 57351
	LTEQ       = 57371
	LTHAN      = 57370
	OP         = 57350
	PACKAGE    = 57377
	RBRACE     = 57354
	RPAREN     = 57352
	STR        = 57356
	STRING     = 57376
	STRUCT     = 57382
	TYPSTRUCT  = 57381
	UNEXPECTED = 57383
	VAR        = 57373
	WHILE      = 57380
	yyErrCode  = 57345

	yyMaxDepth = 200
	yyTabOfs   = -62
)

var (
	yyXLAT = map[int]int{
		57355: 0,  // IDENT (64x)
		57373: 1,  // VAR (54x)
		57354: 2,  // RBRACE (50x)
		57378: 3,  // IF (41x)
		57380: 4,  // WHILE (41x)
		57374: 5,  // COMMA (37x)
		57344: 6,  // $end (31x)
		57349: 7,  // FUNC (31x)
		57377: 8,  // PACKAGE (31x)
		57381: 9,  // TYPSTRUCT (31x)
		57383: 10, // UNEXPECTED (31x)
		57352: 11, // RPAREN (30x)
		57353: 12, // LBRACE (25x)
		57367: 13, // ASSIGN (17x)
		57368: 14, // CASSIGN (17x)
		57382: 15, // STRUCT (16x)
		57361: 16, // BYTE (14x)
		57362: 17, // BYTEA (14x)
		57359: 18, // F32 (14x)
		57365: 19, // F32A (14x)
		57360: 20, // F64 (14x)
		57366: 21, // F64A (14x)
		57357: 22, // I32 (14x)
		57363: 23, // I32A (14x)
		57358: 24, // I64 (14x)
		57364: 25, // I64A (14x)
		57356: 26, // STR (14x)
		57347: 27, // FLOAT (11x)
		57346: 28, // INT (11x)
		57376: 29, // STRING (11x)
		57410: 30, // typeSpec (10x)
		57391: 31, // arglessExpr (8x)
		57393: 32, // argsExpr (8x)
		57398: 33, // expr (8x)
		57401: 34, // ifExpr (8x)
		57404: 35, // opExpr (8x)
		57405: 36, // outNames (8x)
		57411: 37, // whileExpr (8x)
		57390: 38, // arg (7x)
		57399: 39, // exprs (4x)
		57351: 40, // LPAREN (4x)
		57407: 41, // param (4x)
		57392: 42, // args (3x)
		57394: 43, // assignOp (3x)
		57408: 44, // params (3x)
		57395: 45, // defAssign (2x)
		57384: 46, // $@1 (1x)
		57385: 47, // $@2 (1x)
		57386: 48, // $@3 (1x)
		57387: 49, // $@4 (1x)
		57388: 50, // $@5 (1x)
		57396: 51, // defDecl (1x)
		57379: 52, // ELSE (1x)
		57397: 53, // elseExpr (1x)
		57400: 54, // funcDecl (1x)
		57402: 55, // line (1x)
		57403: 56, // lines (1x)
		57406: 57, // packageDecl (1x)
		57409: 58, // structDecl (1x)
		57389: 59, // $default (0x)
		57348: 60, // BOOLEAN (0x)
		57375: 61, // COMMENT (0x)
		57345: 62, // error (0x)
		57372: 63, // GTEQ (0x)
		57369: 64, // GTHAN (0x)
		57371: 65, // LTEQ (0x)
		57370: 66, // LTHAN (0x)
		57350: 67, // OP (0x)
	}

	yySymNames = []string{
		"IDENT",
		"VAR",
		"RBRACE",
		"IF",
		"WHILE",
		"COMMA",
		"$end",
		"FUNC",
		"PACKAGE",
		"TYPSTRUCT",
		"UNEXPECTED",
		"RPAREN",
		"LBRACE",
		"ASSIGN",
		"CASSIGN",
		"STRUCT",
		"BYTE",
		"BYTEA",
		"F32",
		"F32A",
		"F64",
		"F64A",
		"I32",
		"I32A",
		"I64",
		"I64A",
		"STR",
		"FLOAT",
		"INT",
		"STRING",
		"typeSpec",
		"arglessExpr",
		"argsExpr",
		"expr",
		"ifExpr",
		"opExpr",
		"outNames",
		"whileExpr",
		"arg",
		"exprs",
		"LPAREN",
		"param",
		"args",
		"assignOp",
		"params",
		"defAssign",
		"$@1",
		"$@2",
		"$@3",
		"$@4",
		"$@5",
		"defDecl",
		"ELSE",
		"elseExpr",
		"funcDecl",
		"line",
		"lines",
		"packageDecl",
		"structDecl",
		"$default",
		"BOOLEAN",
		"COMMENT",
		"error",
		"GTEQ",
		"GTHAN",
		"LTEQ",
		"LTHAN",
		"OP",
	}

	yyTokenLiteralStrings = map[int]string{}

	yyReductions = map[int]struct{ xsym, components int }{
		0:  {0, 1},
		1:  {56, 0},
		2:  {56, 2},
		3:  {55, 1},
		4:  {55, 1},
		5:  {55, 1},
		6:  {55, 1},
		7:  {55, 1},
		8:  {43, 1},
		9:  {43, 1},
		10: {30, 1},
		11: {30, 1},
		12: {30, 1},
		13: {30, 1},
		14: {30, 1},
		15: {30, 1},
		16: {30, 1},
		17: {30, 1},
		18: {30, 1},
		19: {30, 1},
		20: {30, 1},
		21: {30, 2},
		22: {57, 2},
		23: {45, 0},
		24: {45, 2},
		25: {51, 4},
		26: {46, 0},
		27: {58, 7},
		28: {47, 0},
		29: {54, 12},
		30: {41, 2},
		31: {44, 0},
		32: {44, 1},
		33: {44, 3},
		34: {32, 6},
		35: {31, 4},
		36: {31, 4},
		37: {48, 0},
		38: {37, 6},
		39: {49, 0},
		40: {34, 7},
		41: {53, 0},
		42: {50, 0},
		43: {53, 5},
		44: {35, 1},
		45: {35, 1},
		46: {33, 1},
		47: {33, 1},
		48: {33, 1},
		49: {39, 0},
		50: {39, 1},
		51: {39, 2},
		52: {38, 1},
		53: {38, 1},
		54: {38, 1},
		55: {38, 1},
		56: {38, 4},
		57: {42, 0},
		58: {42, 1},
		59: {42, 3},
		60: {36, 1},
		61: {36, 3},
	}

	yyXErrors = map[yyXError]string{}

	yyParseTab = [107][]uint16{
		// 0
		{1: 61, 6: 61, 61, 61, 61, 61, 56: 63},
		{1: 71, 6: 62, 73, 70, 72, 69, 51: 65, 54: 68, 64, 57: 67, 66},
		{1: 60, 6: 60, 60, 60, 60, 60},
		{1: 59, 6: 59, 59, 59, 59, 59},
		{1: 58, 6: 58, 58, 58, 58, 58},
		// 5
		{1: 57, 6: 57, 57, 57, 57, 57},
		{1: 56, 6: 56, 56, 56, 56, 56},
		{1: 55, 6: 55, 55, 55, 55, 55},
		{168},
		{165},
		// 10
		{159},
		{74},
		{40: 75},
		{77, 5: 31, 11: 31, 41: 78, 44: 76},
		{5: 94, 11: 93},
		// 15
		{15: 90, 83, 84, 81, 87, 82, 88, 79, 85, 80, 86, 89, 30: 91},
		{2: 30, 5: 30, 11: 30},
		{52, 52, 52, 52, 52, 52, 52, 52, 52, 52, 52, 52, 52, 52, 52},
		{51, 51, 51, 51, 51, 51, 51, 51, 51, 51, 51, 51, 51, 51, 51},
		{50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50},
		// 20
		{49, 49, 49, 49, 49, 49, 49, 49, 49, 49, 49, 49, 49, 49, 49},
		{48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48},
		{47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47},
		{46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46},
		{45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45},
		// 25
		{44, 44, 44, 44, 44, 44, 44, 44, 44, 44, 44, 44, 44, 44, 44},
		{43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43},
		{42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42},
		{92},
		{2: 32, 5: 32, 11: 32},
		// 30
		{41, 41, 41, 41, 41, 41, 41, 41, 41, 41, 41, 41, 41, 41, 41},
		{40: 96},
		{77, 41: 95},
		{2: 29, 5: 29, 11: 29},
		{77, 5: 31, 11: 31, 41: 78, 44: 97},
		// 35
		{5: 94, 11: 98},
		{12: 34, 47: 99},
		{12: 100},
		{103, 104, 13, 106, 105, 31: 107, 108, 112, 109, 111, 102, 110, 39: 101},
		{103, 104, 158, 106, 105, 31: 107, 108, 129, 109, 111, 102, 110},
		// 40
		{5: 152, 13: 143, 144, 43: 151},
		{5: 2, 13: 2, 2, 40: 148},
		{141},
		{25, 15: 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 48: 136},
		{23, 15: 23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 49: 113},
		// 45
		{18, 18, 18, 18, 18},
		{17, 17, 17, 17, 17},
		{16, 16, 16, 16, 16},
		{15, 15, 15, 15, 15},
		{14, 14, 14, 14, 14},
		// 50
		{12, 12, 12, 12, 12},
		{118, 15: 90, 83, 84, 81, 87, 82, 88, 79, 85, 80, 86, 89, 116, 115, 117, 119, 38: 114},
		{12: 126},
		{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
		{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		// 55
		{8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8},
		{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7},
		{12: 120},
		{118, 2: 5, 5: 5, 15: 90, 83, 84, 81, 87, 82, 88, 79, 85, 80, 86, 89, 116, 115, 117, 119, 38: 122, 42: 121},
		{2: 123, 5: 124},
		// 60
		{2: 4, 5: 4, 11: 4},
		{6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6},
		{118, 15: 90, 83, 84, 81, 87, 82, 88, 79, 85, 80, 86, 89, 116, 115, 117, 119, 38: 125},
		{2: 3, 5: 3, 11: 3},
		{103, 104, 13, 106, 105, 31: 107, 108, 112, 109, 111, 102, 110, 39: 127},
		// 65
		{103, 104, 128, 106, 105, 31: 107, 108, 129, 109, 111, 102, 110},
		{21, 21, 21, 21, 21, 52: 131, 130},
		{11, 11, 11, 11, 11},
		{22, 22, 22, 22, 22},
		{12: 20, 50: 132},
		// 70
		{12: 133},
		{103, 104, 13, 106, 105, 31: 107, 108, 112, 109, 111, 102, 110, 39: 134},
		{103, 104, 135, 106, 105, 31: 107, 108, 129, 109, 111, 102, 110},
		{19, 19, 19, 19, 19},
		{118, 15: 90, 83, 84, 81, 87, 82, 88, 79, 85, 80, 86, 89, 116, 115, 117, 119, 38: 137},
		// 75
		{12: 138},
		{103, 104, 13, 106, 105, 31: 107, 108, 112, 109, 111, 102, 110, 39: 139},
		{103, 104, 140, 106, 105, 31: 107, 108, 129, 109, 111, 102, 110},
		{24, 24, 24, 24, 24},
		{15: 90, 83, 84, 81, 87, 82, 88, 79, 85, 80, 86, 89, 30: 142},
		// 80
		{39, 39, 39, 39, 39, 13: 143, 144, 43: 145, 45: 146},
		{54, 15: 54, 54, 54, 54, 54, 54, 54, 54, 54, 54, 54, 54, 54, 54, 54},
		{53, 15: 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53},
		{118, 15: 90, 83, 84, 81, 87, 82, 88, 79, 85, 80, 86, 89, 116, 115, 117, 119, 38: 147},
		{26, 26, 26, 26, 26},
		// 85
		{38, 38, 38, 38, 38, 6: 38, 38, 38, 38, 38},
		{118, 5: 5, 11: 5, 15: 90, 83, 84, 81, 87, 82, 88, 79, 85, 80, 86, 89, 116, 115, 117, 119, 38: 122, 42: 149},
		{5: 124, 11: 150},
		{27, 27, 27, 27, 27},
		{154},
		// 90
		{153},
		{5: 1, 13: 1, 1},
		{40: 155},
		{118, 5: 5, 11: 5, 15: 90, 83, 84, 81, 87, 82, 88, 79, 85, 80, 86, 89, 116, 115, 117, 119, 38: 122, 42: 156},
		{5: 124, 11: 157},
		// 95
		{28, 28, 28, 28, 28},
		{1: 33, 6: 33, 33, 33, 33, 33},
		{15: 36, 46: 160},
		{15: 161},
		{12: 162},
		// 100
		{77, 2: 31, 5: 31, 41: 78, 44: 163},
		{2: 164, 5: 94},
		{1: 35, 6: 35, 35, 35, 35, 35},
		{15: 90, 83, 84, 81, 87, 82, 88, 79, 85, 80, 86, 89, 30: 166},
		{1: 39, 6: 39, 39, 39, 39, 39, 13: 143, 144, 43: 145, 45: 167},
		// 105
		{1: 37, 6: 37, 37, 37, 37, 37},
		{1: 40, 6: 40, 40, 40, 40, 40},
	}
)

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyLexerEx interface {
	yyLexer
	Reduced(rule, state int, lval *yySymType) bool
}

func yySymName(c int) (s string) {
	x, ok := yyXLAT[c]
	if ok {
		return yySymNames[x]
	}

	if c < 0x7f {
		return __yyfmt__.Sprintf("'%c'", c)
	}

	return __yyfmt__.Sprintf("%d", c)
}

func yylex1(yylex yyLexer, lval *yySymType) (n int) {
	n = yylex.Lex(lval)
	if n <= 0 {
		n = yyEofCode
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("\nlex %s(%#x %d), lval: %+v\n", yySymName(n), n, n, lval)
	}
	return n
}

func yyParse(yylex yyLexer) int {
	const yyError = 62

	yyEx, _ := yylex.(yyLexerEx)
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, 200)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yyerrok := func() {
		if yyDebug >= 2 {
			__yyfmt__.Printf("yyerrok()\n")
		}
		Errflag = 0
	}
	_ = yyerrok
	yystate := 0
	yychar := -1
	var yyxchar int
	var yyshift int
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
		var ok bool
		if yyxchar, ok = yyXLAT[yychar]; !ok {
			yyxchar = len(yySymNames) // > tab width
		}
	}
	if yyDebug >= 4 {
		var a []int
		for _, v := range yyS[:yyp+1] {
			a = append(a, v.yys)
		}
		__yyfmt__.Printf("state stack %v\n", a)
	}
	row := yyParseTab[yystate]
	yyn = 0
	if yyxchar < len(row) {
		if yyn = int(row[yyxchar]); yyn != 0 {
			yyn += yyTabOfs
		}
	}
	switch {
	case yyn > 0: // shift
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		yyshift = yyn
		if yyDebug >= 2 {
			__yyfmt__.Printf("shift, and goto state %d\n", yystate)
		}
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	case yyn < 0: // reduce
	case yystate == 1: // accept
		if yyDebug >= 2 {
			__yyfmt__.Println("accept")
		}
		goto ret0
	}

	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			if yyDebug >= 1 {
				__yyfmt__.Printf("no action for %s in state %d\n", yySymName(yychar), yystate)
			}
			msg, ok := yyXErrors[yyXError{yystate, yyxchar}]
			if !ok {
				msg, ok = yyXErrors[yyXError{yystate, -1}]
			}
			if !ok && yyshift != 0 {
				msg, ok = yyXErrors[yyXError{yyshift, yyxchar}]
			}
			if !ok {
				msg, ok = yyXErrors[yyXError{yyshift, -1}]
			}
			if yychar > 0 {
				ls := yyTokenLiteralStrings[yychar]
				if ls == "" {
					ls = yySymName(yychar)
				}
				if ls != "" {
					switch {
					case msg == "":
						msg = __yyfmt__.Sprintf("unexpected %s", ls)
					default:
						msg = __yyfmt__.Sprintf("unexpected %s, %s", ls, msg)
					}
				}
			}
			if msg == "" {
				msg = "syntax error"
			}
			yylex.Error(msg)
			Nerrs++
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				row := yyParseTab[yyS[yyp].yys]
				if yyError < len(row) {
					yyn = int(row[yyError]) + yyTabOfs
					if yyn > 0 { // hit
						if yyDebug >= 2 {
							__yyfmt__.Printf("error recovery found error shift in state %d\n", yyS[yyp].yys)
						}
						yystate = yyn /* simulate a shift of "error" */
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery failed\n")
			}
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yySymName(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}

			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	r := -yyn
	x0 := yyReductions[r]
	x, n := x0.xsym, x0.components
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= n
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	exState := yystate
	yystate = int(yyParseTab[yyS[yyp].yys][x]) + yyTabOfs
	/* reduction by production r */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce using rule %v (%s), and goto state %d\n", r, yySymNames[x], yystate)
	}

	switch r {
	case 7:
		{
			yylex.Error(fmt.Sprintf("error in line #%d: unexpected character '%s'", lineNo, yyS[yypt-0].tok))
		}
	case 10:
		{
			yyVAL.tok = yyS[yypt-0].tok
		}
	case 11:
		{
			yyVAL.tok = yyS[yypt-0].tok
		}
	case 12:
		{
			yyVAL.tok = yyS[yypt-0].tok
		}
	case 13:
		{
			yyVAL.tok = yyS[yypt-0].tok
		}
	case 14:
		{
			yyVAL.tok = yyS[yypt-0].tok
		}
	case 15:
		{
			yyVAL.tok = yyS[yypt-0].tok
		}
	case 16:
		{
			yyVAL.tok = yyS[yypt-0].tok
		}
	case 17:
		{
			yyVAL.tok = yyS[yypt-0].tok
		}
	case 18:
		{
			yyVAL.tok = yyS[yypt-0].tok
		}
	case 19:
		{
			yyVAL.tok = yyS[yypt-0].tok
		}
	case 20:
		{
			yyVAL.tok = yyS[yypt-0].tok
		}
	case 21:
		{
			yyVAL.tok = yyS[yypt-0].tok
		}
	case 22:
		{
			cxt.AddModule(MakeModule(yyS[yypt-0].tok))
		}
	case 23:
		{
			yyVAL.arg = nil
		}
	case 24:
		{
			yyVAL.arg = yyS[yypt-0].arg
		}
	case 25:
		{
			if mod, err := cxt.GetCurrentModule(); err == nil {
				var val *CXArgument
				if yyS[yypt-0].arg == nil {
					var zeroVal []byte
					switch yyS[yypt-1].tok {
					case "byte":
						zeroVal = []byte{byte(0)}
					case "i32":
						zeroVal = encoder.Serialize(int32(0))
					case "i64":
						zeroVal = encoder.Serialize(int64(0))
					case "f32":
						zeroVal = encoder.Serialize(float32(0))
					case "f64":
						zeroVal = encoder.Serialize(float64(0))
					case "[]byte":
						zeroVal = []byte{byte(0)}
					case "[]i32":
						zeroVal = encoder.Serialize([]int32{0})
					case "[]i64":
						zeroVal = encoder.Serialize([]int64{0})
					case "[]f32":
						zeroVal = encoder.Serialize([]float32{0})
					case "[]f64":
						zeroVal = encoder.Serialize([]float64{0})
					}
					val = MakeArgument(&zeroVal, MakeType(yyS[yypt-1].tok))
				} else {
					switch yyS[yypt-1].tok {
					case "byte":
						var ds int64
						encoder.DeserializeRaw(*yyS[yypt-0].arg.Value, &ds)
						//new := encoder.Serialize(byte(ds))
						new := []byte{byte(ds)}
						val = MakeArgument(&new, MakeType("byte"))
					case "i32":
						var ds int64
						encoder.DeserializeRaw(*yyS[yypt-0].arg.Value, &ds)
						new := encoder.Serialize(int32(ds))
						val = MakeArgument(&new, MakeType("i32"))
					case "i64": /* stays the same */
					case "f32":
						var ds float64
						encoder.DeserializeRaw(*yyS[yypt-0].arg.Value, &ds)
						new := encoder.Serialize(float32(ds))
						val = MakeArgument(&new, MakeType("f32"))
					case "f64":
						val = yyS[yypt-0].arg
					case "[]byte":
						/* var ds []int64 */
						/* encoder.DeserializeRaw(*$4.Value, &ds) */
						/* new := make([]byte, len(ds)) */
						/* for i, val := range ds { */
						/*         new[i] = byte(val) */
						/*     } */
						/* val = MakeArgument(&new, MakeType("[]byte")) */
						val = yyS[yypt-0].arg
					case "[]i32":
						/* var ds []int64 */
						/* encoder.DeserializeRaw(*$4.Value, &ds) */
						/* new := make([]int32, len(ds)) */
						/* for i, val := range ds { */
						/*         new[i] = int32(val) */
						/*     } */
						/* sNew := encoder.Serialize(new) */
						/* val = MakeArgument(&sNew, MakeType("[]i32")) */
						val = yyS[yypt-0].arg
					case "[]i64":
						val = yyS[yypt-0].arg
					case "[]f32":
						/* var ds []float64 */
						/* encoder.DeserializeRaw(*$4.Value, &ds) */
						/* new := make([]float32, len(ds)) */
						/* for i, val := range ds { */
						/*         new[i] = float32(val) */
						/*     } */
						/* sNew := encoder.Serialize(new) */
						/* val = MakeArgument(&sNew, MakeType("[]f32")) */
						val = yyS[yypt-0].arg
					case "[]f64": /* stays the same */
						val = yyS[yypt-0].arg
					}
					//val = $4
				}

				mod.AddDefinition(MakeDefinition(yyS[yypt-2].tok, val.Value, MakeType(yyS[yypt-1].tok)))
			}
		}
	case 26:
		{
			if mod, err := cxt.GetCurrentModule(); err == nil {
				strct := MakeStruct(yyS[yypt-0].tok)
				mod.AddStruct(strct)
			}
		}
	case 27:
		{
			if strct, err := cxt.GetCurrentStruct(); err == nil {
				for _, fld := range yyS[yypt-1].params {
					fldFromParam := MakeField(fld.Name, fld.Typ)
					strct.AddField(fldFromParam)
				}
			}
		}
	case 28:
		{
			if mod, err := cxt.GetCurrentModule(); err == nil {
				mod.AddFunction(MakeFunction(yyS[yypt-6].tok))
				if fn, err := mod.GetCurrentFunction(); err == nil {
					for _, inp := range yyS[yypt-4].params {
						fn.AddInput(inp)
					}
					for _, out := range yyS[yypt-1].params {
						fn.AddOutput(out)
					}
				}
			}
		}
	case 30:
		{
			yyVAL.param = MakeParameter(yyS[yypt-1].tok, MakeType(yyS[yypt-0].tok))
		}
	case 31:
		{
			yyVAL.params = nil
		}
	case 32:
		{
			var params []*CXParameter
			params = append(params, yyS[yypt-0].param)
			yyVAL.params = params
		}
	case 33:
		{
			yyS[yypt-2].params = append(yyS[yypt-2].params, yyS[yypt-0].param)
			yyVAL.params = yyS[yypt-2].params
		}
	case 34:
		{
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := cxt.GetCurrentFunction(); err == nil {
					if op, err := cxt.GetFunction(yyS[yypt-3].tok, mod.Name); err == nil {
						expr := MakeExpression(op)

						for _, outName := range yyS[yypt-5].names {
							expr.AddOutputName(outName)
						}

						fn.AddExpression(expr)

						if expr, err := fn.GetCurrentExpression(); err == nil {
							for _, arg := range yyS[yypt-1].args {
								expr.AddArgument(arg)
							}
							yyVAL.expr = expr
						}

					} else {
						panic(err)
					}
				}

			} else {
				panic(err)
			}
		}
	case 35:
		{
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := cxt.GetCurrentFunction(); err == nil {
					if op, err := cxt.GetFunction(yyS[yypt-3].tok, mod.Name); err == nil {
						expr := MakeExpression(op)

						fn.AddExpression(expr)

						if expr, err := fn.GetCurrentExpression(); err == nil {
							for _, arg := range yyS[yypt-1].args {
								expr.AddArgument(arg)
							}
							yyVAL.expr = expr
						}

					} else {
						panic(err)
					}
				}

			} else {
				panic(err)
			}
		}
	case 36:
		{
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := cxt.GetCurrentFunction(); err == nil {
					if yyS[yypt-0].arg == nil {
						if op, err := cxt.GetFunction("initDef", mod.Name); err == nil {
							expr := MakeExpression(op)
							fn.AddExpression(expr)
							expr.AddOutputName(yyS[yypt-2].tok)
							typ := []byte(yyS[yypt-1].tok)
							arg := MakeArgument(&typ, MakeType("str"))
							expr.AddArgument(arg)

							if strct, err := cxt.GetStruct(yyS[yypt-1].tok, mod.Name); err == nil {
								for _, fld := range strct.Fields {
									expr := MakeExpression(op)
									fn.AddExpression(expr)
									expr.AddOutputName(fmt.Sprintf("%s.%s", yyS[yypt-2].tok, fld.Name))
									typ := []byte(fld.Typ.Name)
									arg := MakeArgument(&typ, MakeType("str"))
									expr.AddArgument(arg)
								}
							}
						}
					} else {
						switch yyS[yypt-1].tok {
						case "byte":
							var ds int64
							encoder.DeserializeRaw(*yyS[yypt-0].arg.Value, &ds)
							new := []byte{byte(ds)}
							val := MakeArgument(&new, MakeType("byte"))

							if op, err := cxt.GetFunction("idByte", mod.Name); err == nil {
								expr := MakeExpression(op)
								fn.AddExpression(expr)
								expr.AddOutputName(yyS[yypt-2].tok)
								expr.AddArgument(val)
							}
						case "i32":
							var ds int64
							encoder.DeserializeRaw(*yyS[yypt-0].arg.Value, &ds)
							new := encoder.Serialize(int32(ds))
							val := MakeArgument(&new, MakeType("i32"))

							if op, err := cxt.GetFunction("idI32", mod.Name); err == nil {
								expr := MakeExpression(op)
								fn.AddExpression(expr)
								expr.AddOutputName(yyS[yypt-2].tok)
								expr.AddArgument(val)
							}
						case "f32":
							var ds float64
							encoder.DeserializeRaw(*yyS[yypt-0].arg.Value, &ds)
							new := encoder.Serialize(float32(ds))
							val := MakeArgument(&new, MakeType("f32"))

							if op, err := cxt.GetFunction("idF32", mod.Name); err == nil {
								expr := MakeExpression(op)
								fn.AddExpression(expr)
								expr.AddOutputName(yyS[yypt-2].tok)
								expr.AddArgument(val)
							}
						default:
							val := yyS[yypt-0].arg
							var getFn string
							switch yyS[yypt-1].tok {
							case "i64":
								getFn = "idI64"
							case "f64":
								getFn = "idF64"
							case "[]byte":
								getFn = "idByteA"
							case "[]i32":
								getFn = "idI32A"
							case "[]i64":
								getFn = "idI64A"
							case "[]f32":
								getFn = "idF32A"
							case "[]f64":
								getFn = "idF64A"
							}

							if op, err := cxt.GetFunction(getFn, mod.Name); err == nil {
								expr := MakeExpression(op)
								fn.AddExpression(expr)
								expr.AddOutputName(yyS[yypt-2].tok)
								expr.AddArgument(val)
							}
						}

					}
				}
			}
		}
	case 37:
		{
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					if goToFn, err := cxt.GetFunction("goTo", mod.Name); err == nil {
						expr := MakeExpression(goToFn)
						fn.AddExpression(expr)
					}
				}
			}
		}
	case 38:
		{
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {

					goToExpr := fn.Expressions[len(fn.Expressions)-len(yyS[yypt-1].exprs)-1]

					elseLines := encoder.Serialize(int32(len(yyS[yypt-1].exprs) + 2))
					thenLines := encoder.Serialize(int32(1))

					predVal := *yyS[yypt-3].arg.Value

					goToExpr.AddArgument(MakeArgument(&predVal, yyS[yypt-3].arg.Typ))
					goToExpr.AddArgument(MakeArgument(&thenLines, MakeType("i32")))
					goToExpr.AddArgument(MakeArgument(&elseLines, MakeType("i32")))

					if goToFn, err := cxt.GetFunction("goTo", mod.Name); err == nil {
						goToExpr := MakeExpression(goToFn)
						fn.AddExpression(goToExpr)

						elseLines := encoder.Serialize(int32(1))
						thenLines := encoder.Serialize(int32(-len(yyS[yypt-1].exprs)))

						goToExpr.AddArgument(MakeArgument(&predVal, yyS[yypt-3].arg.Typ))
						goToExpr.AddArgument(MakeArgument(&thenLines, MakeType("i32")))
						goToExpr.AddArgument(MakeArgument(&elseLines, MakeType("i32")))
					}

				}
			}
		}
	case 39:
		{
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					if goToFn, err := cxt.GetFunction("goTo", mod.Name); err == nil {
						expr := MakeExpression(goToFn)
						fn.AddExpression(expr)
					}
				}
			}
		}
	case 40:
		{
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {

					var goToExpr *CXExpression
					if len(yyS[yypt-0].exprs) > 0 {
						goToExpr = fn.Expressions[len(fn.Expressions)-2-len(yyS[yypt-2].exprs)-len(yyS[yypt-0].exprs)]
					} else {
						goToExpr = fn.Expressions[len(fn.Expressions)-2-len(yyS[yypt-2].exprs)]
					}

					elseLines := encoder.Serialize(int32(len(yyS[yypt-2].exprs) + 2))
					thenLines := encoder.Serialize(int32(1))

					predVal := *yyS[yypt-4].arg.Value

					goToExpr.AddArgument(MakeArgument(&predVal, yyS[yypt-4].arg.Typ))
					goToExpr.AddArgument(MakeArgument(&thenLines, MakeType("i32")))
					goToExpr.AddArgument(MakeArgument(&elseLines, MakeType("i32")))
				}
			}
		}
	case 41:
		{
			exprs := make([]*CXExpression, 0)
			yyVAL.exprs = exprs
		}
	case 42:
		{
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					if goToFn, err := cxt.GetFunction("goTo", mod.Name); err == nil {
						expr := MakeExpression(goToFn)
						fn.AddExpression(expr)
					}
				}
			}
		}
	case 43:
		{
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {

					goToExpr := fn.Expressions[len(fn.Expressions)-1-len(yyS[yypt-1].exprs)]

					elseLines := encoder.Serialize(int32(0))
					thenLines := encoder.Serialize(int32(len(yyS[yypt-1].exprs) + 1))

					predVal := []byte{1}

					goToExpr.AddArgument(MakeArgument(&predVal, MakeType("byte")))
					goToExpr.AddArgument(MakeArgument(&thenLines, MakeType("i32")))
					goToExpr.AddArgument(MakeArgument(&elseLines, MakeType("i32")))
				}
			}

			yyVAL.exprs = yyS[yypt-1].exprs
		}
	case 49:
		{
			exprs := make([]*CXExpression, 0)
			yyVAL.exprs = exprs
		}
	case 50:
		{
			exprs := make([]*CXExpression, 1)
			exprs[0] = yyS[yypt-0].expr
			yyVAL.exprs = exprs

		}
	case 51:
		{
			yyS[yypt-1].exprs = append(yyS[yypt-1].exprs, yyS[yypt-0].expr)
			yyVAL.exprs = yyS[yypt-1].exprs
		}
	case 52:
		{
			val := encoder.SerializeAtomic(yyS[yypt-0].i64)
			yyVAL.arg = MakeArgument(&val, MakeType("i64"))
		}
	case 53:
		{
			val := encoder.Serialize(yyS[yypt-0].f64)
			yyVAL.arg = MakeArgument(&val, MakeType("f64"))
		}
	case 54:
		{
			var str string
			str = strings.TrimPrefix(yyS[yypt-0].tok, "\"")
			str = strings.TrimSuffix(str, "\"")

			val := []byte(str)
			yyVAL.arg = MakeArgument(&val, MakeType("str"))
		}
	case 55:
		{
			val := []byte(yyS[yypt-0].tok)
			yyVAL.arg = MakeArgument(&val, MakeType("ident"))
		}
	case 56:
		{
			switch yyS[yypt-3].tok {
			case "[]byte":
				vals := make([]byte, len(yyS[yypt-1].args))
				for i, arg := range yyS[yypt-1].args {
					var val int64
					encoder.DeserializeRaw(*arg.Value, &val)
					vals[i] = byte(val)
				}
				sVal := encoder.Serialize(vals)
				yyVAL.arg = MakeArgument(&sVal, MakeType("[]i32"))
			case "[]i32":
				vals := make([]int32, len(yyS[yypt-1].args))
				for i, arg := range yyS[yypt-1].args {
					var val int64
					encoder.DeserializeRaw(*arg.Value, &val)
					vals[i] = int32(val)
				}
				sVal := encoder.Serialize(vals)
				yyVAL.arg = MakeArgument(&sVal, MakeType("[]i32"))
			case "[]i64":
				vals := make([]int64, len(yyS[yypt-1].args))
				for i, arg := range yyS[yypt-1].args {
					var val int64
					encoder.DeserializeRaw(*arg.Value, &val)
					vals[i] = val
				}
				sVal := encoder.Serialize(vals)
				yyVAL.arg = MakeArgument(&sVal, MakeType("[]i64"))
			case "[]f32":
				vals := make([]float32, len(yyS[yypt-1].args))
				for i, arg := range yyS[yypt-1].args {
					var val float64
					encoder.DeserializeRaw(*arg.Value, &val)
					vals[i] = float32(val)
				}
				sVal := encoder.Serialize(vals)
				yyVAL.arg = MakeArgument(&sVal, MakeType("[]f32"))
			case "[]f64":
				vals := make([]float64, len(yyS[yypt-1].args))
				for i, arg := range yyS[yypt-1].args {
					var val float64
					encoder.DeserializeRaw(*arg.Value, &val)
					vals[i] = val
				}
				sVal := encoder.Serialize(vals)
				yyVAL.arg = MakeArgument(&sVal, MakeType("[]f64"))
			}
		}
	case 58:
		{
			var args []*CXArgument
			args = append(args, yyS[yypt-0].arg)
			yyVAL.args = args
		}
	case 59:
		{
			yyS[yypt-2].args = append(yyS[yypt-2].args, yyS[yypt-0].arg)
			yyVAL.args = yyS[yypt-2].args
		}
	case 60:
		{
			outNames := make([]string, 1)
			outNames[0] = yyS[yypt-0].tok
			yyVAL.names = outNames
		}
	case 61:
		{
			yyS[yypt-2].names = append(yyS[yypt-2].names, yyS[yypt-0].tok)
			yyVAL.names = yyS[yypt-2].names
		}

	}

	if yyEx != nil && yyEx.Reduced(r, exState, &yyVAL) {
		return -1
	}
	goto yystack /* stack new state and value */
}
