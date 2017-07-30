// CAUTION: Generated file - DO NOT EDIT.

package main

import __yyfmt__ "fmt"

import (
	. "github.com/skycoin/cx/src/base"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"strings"
)

var cxt = MakeContext()

type yySymType struct {
	yys int
	i32 int32
	f64 float64
	//    str string
	tok      string
	fun      *CXFunction
	params   []*CXParameter
	param    *CXParameter
	args     []*CXArgument
	arg      *CXArgument
	outNames []string

	cxt *CXContext
	mod *CXModule
}

type yyXError struct {
	state, xsym int
}

const (
	yyDefault = 57362
	yyEofCode = 57344
	COMMA     = 57358
	COMMENT   = 57359
	FLOAT     = 57347
	FUNC      = 57348
	IDENT     = 57354
	INT       = 57346
	KEYWORD   = 57355
	LBRACE    = 57352
	LPAREN    = 57350
	OP        = 57349
	RBRACE    = 57353
	RPAREN    = 57351
	STRING    = 57360
	TYP       = 57356
	VAR       = 57357
	yyErrCode = 57345

	yyMaxDepth = 200
	yyTabOfs   = -34
)

var (
	yyXLAT = map[int]int{
		57358: 0,  // COMMA (21x)
		57344: 1,  // $end (20x)
		59:    2,  // ';' (20x)
		10:    3,  // '\n' (20x)
		57348: 4,  // FUNC (20x)
		57355: 5,  // KEYWORD (20x)
		57357: 6,  // VAR (20x)
		57354: 7,  // IDENT (17x)
		57351: 8,  // RPAREN (16x)
		57353: 9,  // RBRACE (14x)
		57356: 10, // TYP (6x)
		57363: 11, // arg (4x)
		57347: 12, // FLOAT (4x)
		57346: 13, // INT (4x)
		57349: 14, // OP (4x)
		57360: 15, // STRING (4x)
		57352: 16, // LBRACE (3x)
		57350: 17, // LPAREN (3x)
		57374: 18, // param (3x)
		57364: 19, // args (2x)
		57366: 20, // def (2x)
		57367: 21, // expr (2x)
		57369: 22, // fun (2x)
		57373: 23, // outNames (2x)
		57375: 24, // params (2x)
		57361: 25, // $@1 (1x)
		57365: 26, // cxtAdder (1x)
		57368: 27, // fnAdder (1x)
		57370: 28, // input (1x)
		57371: 29, // line (1x)
		57372: 30, // modAdder (1x)
		57376: 31, // term (1x)
		57362: 32, // $default (0x)
		42:    33, // '*' (0x)
		43:    34, // '+' (0x)
		44:    35, // ',' (0x)
		45:    36, // '-' (0x)
		47:    37, // '/' (0x)
		61:    38, // '=' (0x)
		57359: 39, // COMMENT (0x)
		57345: 40, // error (0x)
	}

	yySymNames = []string{
		"COMMA",
		"$end",
		"';'",
		"'\\n'",
		"FUNC",
		"KEYWORD",
		"VAR",
		"IDENT",
		"RPAREN",
		"RBRACE",
		"TYP",
		"arg",
		"FLOAT",
		"INT",
		"OP",
		"STRING",
		"LBRACE",
		"LPAREN",
		"param",
		"args",
		"def",
		"expr",
		"fun",
		"outNames",
		"params",
		"$@1",
		"cxtAdder",
		"fnAdder",
		"input",
		"line",
		"modAdder",
		"term",
		"$default",
		"'*'",
		"'+'",
		"','",
		"'-'",
		"'/'",
		"'='",
		"COMMENT",
		"error",
	}

	yyTokenLiteralStrings = map[int]string{}

	yyReductions = map[int]struct{ xsym, components int }{
		0:  {0, 1},
		1:  {28, 0},
		2:  {28, 2},
		3:  {29, 1},
		4:  {29, 1},
		5:  {29, 1},
		6:  {18, 2},
		7:  {24, 1},
		8:  {24, 3},
		9:  {24, 0},
		10: {11, 1},
		11: {11, 1},
		12: {11, 1},
		13: {11, 1},
		14: {11, 4},
		15: {19, 0},
		16: {19, 1},
		17: {19, 3},
		18: {23, 1},
		19: {23, 3},
		20: {21, 6},
		21: {27, 0},
		22: {27, 1},
		23: {27, 2},
		24: {20, 5},
		25: {25, 0},
		26: {22, 12},
		27: {30, 1},
		28: {30, 1},
		29: {30, 2},
		30: {30, 2},
		31: {26, 2},
		32: {31, 1},
		33: {31, 1},
	}

	yyXErrors = map[yyXError]string{}

	yyParseTab = [58][]uint8{
		// 0
		{1: 33, 33, 33, 33, 33, 33, 28: 35},
		{1: 34, 45, 46, 41, 44, 40, 20: 42, 22: 43, 26: 39, 29: 36, 38, 37},
		{1: 32, 32, 32, 32, 32, 32},
		{1: 31, 31, 31, 31, 31, 31},
		{1: 30, 30, 30, 41, 30, 40, 20: 91, 22: 90},
		// 5
		{1: 29, 29, 29, 29, 29, 29},
		{7: 86},
		{7: 48},
		{1: 7, 7, 7, 7, 7, 7},
		{1: 6, 6, 6, 6, 6, 6},
		// 10
		{7: 47},
		{1: 2, 2, 2, 2, 2, 2},
		{1: 1, 1, 1, 1, 1, 1},
		{1: 3, 3, 3, 3, 3, 3},
		{17: 49},
		// 15
		{25, 7: 50, 25, 18: 51, 24: 52},
		{10: 85},
		{27, 8: 27},
		{53, 8: 54},
		{7: 50, 18: 84},
		// 20
		{17: 55},
		{25, 7: 50, 25, 18: 51, 24: 56},
		{53, 8: 57},
		{16: 9, 25: 58},
		{16: 59},
		// 25
		{7: 60, 9: 13, 21: 62, 23: 61, 27: 63},
		{16, 14: 16},
		{66, 14: 67},
		{7: 12, 9: 12},
		{7: 60, 9: 65, 21: 64, 23: 61},
		// 30
		{7: 11, 9: 11},
		{1: 8, 8, 8, 8, 8, 8},
		{7: 83},
		{7: 68},
		{17: 69},
		// 35
		{19, 7: 73, 19, 10: 74, 75, 71, 70, 15: 72, 19: 76},
		{24, 24, 24, 24, 24, 24, 24, 8: 24, 24},
		{23, 23, 23, 23, 23, 23, 23, 8: 23, 23},
		{22, 22, 22, 22, 22, 22, 22, 8: 22, 22},
		{21, 21, 21, 21, 21, 21, 21, 8: 21, 21},
		// 40
		{16: 80},
		{18, 8: 18, 18},
		{77, 8: 78},
		{7: 73, 10: 74, 79, 71, 70, 15: 72},
		{7: 14, 9: 14},
		// 45
		{17, 8: 17, 17},
		{19, 7: 73, 9: 19, 74, 75, 71, 70, 15: 72, 19: 81},
		{77, 9: 82},
		{20, 20, 20, 20, 20, 20, 20, 8: 20, 20},
		{15, 14: 15},
		// 50
		{26, 8: 26},
		{28, 8: 28},
		{10: 87},
		{14: 88},
		{7: 73, 10: 74, 89, 71, 70, 15: 72},
		// 55
		{1: 10, 10, 10, 10, 10, 10},
		{1: 5, 5, 5, 5, 5, 5},
		{1: 4, 4, 4, 4, 4, 4},
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
	const yyError = 40

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
	case 6:
		{
			yyVAL.param = MakeParameter(yyS[yypt-1].tok, MakeType(yyS[yypt-0].tok))
		}
	case 7:
		{
			var params []*CXParameter
			params = append(params, yyS[yypt-0].param)
			yyVAL.params = params
		}
	case 8:
		{
			yyS[yypt-2].params = append(yyS[yypt-2].params, yyS[yypt-0].param)
			yyVAL.params = yyS[yypt-2].params
		}
	case 9:
		{
			yyVAL.params = nil
		}
	case 10:
		{
			val := encoder.SerializeAtomic(yyS[yypt-0].i32)
			yyVAL.arg = MakeArgument(&val, MakeType("i32"))
		}
	case 11:
		{
			val := encoder.Serialize(yyS[yypt-0].f64)
			yyVAL.arg = MakeArgument(&val, MakeType("f64"))
		}
	case 12:
		{
			var str string
			str = strings.TrimPrefix(yyS[yypt-0].tok, "\"")
			str = strings.TrimSuffix(str, "\"")

			val := []byte(str)
			yyVAL.arg = MakeArgument(&val, MakeType("str"))
		}
	case 13:
		{
			val := []byte(yyS[yypt-0].tok)
			yyVAL.arg = MakeArgument(&val, MakeType("ident"))
		}
	case 14:
		{
			var val []byte
			for _, arg := range yyS[yypt-1].args {
				val = append(val, *arg.Value...)
			}
			yyVAL.arg = MakeArgument(&val, MakeType(yyS[yypt-3].tok))
		}
	case 16:
		{
			var args []*CXArgument
			args = append(args, yyS[yypt-0].arg)
			yyVAL.args = args
		}
	case 17:
		{
			yyS[yypt-2].args = append(yyS[yypt-2].args, yyS[yypt-0].arg)
			yyVAL.args = yyS[yypt-2].args
		}
	case 18:
		{
			outNames := make([]string, 1)
			outNames[0] = yyS[yypt-0].tok
			yyVAL.outNames = outNames
		}
	case 19:
		{
			yyS[yypt-2].outNames = append(yyS[yypt-2].outNames, yyS[yypt-0].tok)
			yyVAL.outNames = yyS[yypt-2].outNames
		}
	case 20:
		{
			if mod, err := cxt.GetCurrentModule(); err == nil {
				if fn, err := cxt.GetCurrentFunction(); err == nil {
					if op, err := cxt.GetFunction(yyS[yypt-3].tok, mod.Name); err == nil {
						expr := MakeExpression(op)

						for _, outName := range yyS[yypt-5].outNames {
							expr.AddOutputName(outName)
						}

						fn.AddExpression(expr)

						if expr, err := fn.GetCurrentExpression(); err == nil {
							for _, arg := range yyS[yypt-1].args {
								expr.AddArgument(arg)
							}
						}

					} else {
						panic(err)
					}
				}

			} else {
				panic(err)
			}
		}
	case 24:
		{
			if mod, err := cxt.GetCurrentModule(); err == nil {
				mod.AddDefinition(MakeDefinition(yyS[yypt-3].tok, yyS[yypt-0].arg.Value, MakeType(yyS[yypt-2].tok)))
			}
		}
	case 25:
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
	case 31:
		{
			cxt.AddModule(MakeModule(yyS[yypt-0].tok))
			yyVAL.cxt = cxt
		}

	}

	if yyEx != nil && yyEx.Reduced(r, exState, &yyVAL) {
		return -1
	}
	goto yystack /* stack new state and value */
}
