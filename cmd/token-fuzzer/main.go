package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"

	. "github.com/skycoin/cx/cxparser/cxparsingcompletor"
)

var tbits int = 4
var tchunks int = 100

func lower(ch rune) rune     { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func isLetter(ch rune) bool  { return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' }
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }
func isHex(ch rune) bool     { return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f' }

func randomString(length int, notIdent bool) string {
	result := ""
	l := length - 1
	for l >= 0 {
		num := rand.Intn(128)
		if l == length-1 {
			if isDecimal(rune(num)) {
				num = 0
			}
		}
		if notIdent {
			if num >= 32 && num < 128 && num != '"' && num != '\\' && num != '\n' {
				result += string([]byte{byte(num)})
				l--
			}
		} else {
			if isLetter(rune(num)) || isDecimal(rune(num)) {
				result += string([]byte{byte(num)})
				l--
			}
		}
	}
	return result
}

var tok2ideal map[int]string = map[int]string{
	BYTE_LITERAL:    "BYTE_LITERAL",
	BOOLEAN_LITERAL: "BOOLEAN_LITERAL",
	INT_LITERAL:     "INT_LITERAL",
	LONG_LITERAL:    "LONG_LITERAL",
	FLOAT_LITERAL:   "FLOAT_LITERAL",
	DOUBLE_LITERAL:  "DOUBLE_LITERAL",
	FUNC:            "func",
	OP:              "",
	LPAREN:          "(",
	RPAREN:          ")",
	LBRACE:          "{",
	RBRACE:          "}",
	LBRACK:          "[",
	RBRACK:          "]",
	IDENTIFIER:      "IDENTIFIER",
	VAR:             "var",
	COMMA:           ",",
	PERIOD:          ".",
	COMMENT:         "COMMENT",
	STRING_LITERAL:  "STRING_LITERAL",
	PACKAGE:         "package",
	IF:              "if",
	ELSE:            "else",
	FOR:             "for",
	TYPSTRUCT:       "",
	STRUCT:          "struct",
	SEMICOLON:       ";",
	NEWLINE:         "\n",
	ASSIGN:          "=",
	CASSIGN:         ":=",
	IMPORT:          "import",
	RETURN:          "return",
	GOTO:            "goto",
	GT_OP:           ">",
	LT_OP:           "<",
	GTEQ_OP:         ">=",
	LTEQ_OP:         "<=",
	EQUAL:           "",
	COLON:           ":",
	NEW:             "new",
	EQUALWORD:       "",
	GTHANWORD:       "",
	LTHANWORD:       "",
	GTHANEQ:         "",
	LTHANEQ:         "",
	UNEQUAL:         "",
	AND:             "",
	OR:              "",
	ADD_OP:          "+",
	SUB_OP:          "-",
	MUL_OP:          "*",
	DIV_OP:          "/",
	MOD_OP:          "%",
	REF_OP:          "&",
	NEG_OP:          "!",
	AFFVAR:          "",
	PLUSPLUS:        "",
	MINUSMINUS:      "",
	REMAINDER:       "",
	LEFTSHIFT:       "<<",
	RIGHTSHIFT:      ">>",
	EXP:             "",
	NOT:             "",
	BITXOR_OP:       "^",
	BITOR_OP:        "|",
	BITCLEAR_OP:     "&^",
	PLUSEQ:          "",
	MINUSEQ:         "",
	MULTEQ:          "",
	DIVEQ:           "",
	REMAINDEREQ:     "",
	EXPEQ:           "",
	LEFTSHIFTEQ:     "",
	RIGHTSHIFTEQ:    "",
	BITANDEQ:        "",
	BITXOREQ:        "",
	BITOREQ:         "",

	DEC_OP:       "--",
	INC_OP:       "++",
	PTR_OP:       "",
	LEFT_OP:      "<<",
	RIGHT_OP:     ">>",
	GE_OP:        ">=",
	LE_OP:        "<=",
	EQ_OP:        "==",
	NE_OP:        "!=",
	AND_OP:       "&&",
	OR_OP:        "||",
	ADD_ASSIGN:   "+=",
	AND_ASSIGN:   "&=",
	LEFT_ASSIGN:  "<<=",
	MOD_ASSIGN:   "%=",
	MUL_ASSIGN:   "*=",
	DIV_ASSIGN:   "/=",
	OR_ASSIGN:    "|=",
	RIGHT_ASSIGN: ">>=",
	SUB_ASSIGN:   "-=",
	XOR_ASSIGN:   "^=",
	BOOL:         "bool",
	F32:          "f32",
	F64:          "f64",
	I8:           "i8",
	I16:          "i16",
	I32:          "i32",
	I64:          "i64",
	STR:          "str",
	UI8:          "ui8",
	UI16:         "ui16",
	UI32:         "ui32",
	UI64:         "ui64",
	UNION:        "union",
	ENUM:         "enum",
	CONST:        "const",
	CASE:         "case",
	DEFAULT:      "default",
	SWITCH:       "switch",
	BREAK:        "break",
	CONTINUE:     "continue",
	TYPE:         "type",

	/* Types */
	BASICTYPE: "",
	/* Removers */
	REM:     ":rem",
	DEF:     "def",
	EXPR:    "",
	FIELD:   "field",
	CLAUSES: "clauses",
	OBJECT:  "",
	OBJECTS: "",
	/* Debugging */
	DSTACK:   ":ds",
	DPROGRAM: ":dp",
	DSTATE:   ":dl",
	/* Affordances */
	AFF:   "aff",
	CAFF:  ":aff",
	TAG:   "",
	INFER: "#",
	VALUE: "",
	/* Pointers */
	ADDR:                   "",
	SHORT_LITERAL:          "SHORT_LITERAL",
	UNSIGNED_SHORT_LITERAL: "UNSIGNED_SHORT_LITERAL",
	UNSIGNED_INT_LITERAL:   "UNSIGNED_INT_LITERAL",
	UNSIGNED_LONG_LITERAL:  "UNSIGNED_LONG_LITERAL",
	UNSIGNED_BYTE_LITERAL:  "UNSIGNED_BYTE_LITERAL",
}

func tok2Ideal0(j int) string {
	tok := tok2ideal[j]
	switch tok {
	case "BYTE_LITERAL":
		return strconv.FormatInt(rand.Int63()%254-127, 10) + "B"
	case "SHORT_LITERAL":
		return strconv.FormatInt(rand.Int63()%65534-32767, 10) + "H"
	case "INT_LITERAL":
		return strconv.FormatInt(rand.Int63()%((1<<32)-2)-(1<<31)+1, 10)
	case "LONG_LITERAL":
		return strconv.FormatInt(rand.Int63()-(1<<62), 10) + "L"
	case "UNSIGNED_BYTE_LITERAL":
		return strconv.FormatInt(rand.Int63()%255, 10) + "UB"
	case "UNSIGNED_SHORT_LITERAL":
		return strconv.FormatInt(rand.Int63()%65535, 10) + "UH"
	case "UNSIGNED_INT_LITERAL":
		return strconv.FormatInt(rand.Int63()%((1<<32)-1), 10) + "U"
	case "UNSIGNED_LONG_LITERAL":
		return strconv.FormatInt(rand.Int63(), 10) + "UL"
	case "FLOAT_LITERAL":
		return strconv.FormatFloat(float64(rand.Float32()), 'f', 14, 32)
	case "DOUBLE_LITERAL":
		return strconv.FormatFloat(rand.NormFloat64(), 'f', 32, 64) + "D"
	case "STRING_LITERAL":
		return "\"" + randomString(int(rand.Intn(100))+10, true) + "\""
	case "IDENTIFIER":
		return randomString(int(rand.Intn(32))+5, false)
	case "COMMENT":
		return "/* " + randomString(int(rand.Intn(200)), false) + " */"
	default:
		return tok
	}
}

func tok2Ideal(j int) string {
	return tok2Ideal0(j) + " "
}

func main() {
	var outfile string = "tokens.txt"
	commandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	commandLine.StringVar(&outfile, "o", outfile, "file to output to")
	commandLine.IntVar(&tbits, "b", tbits, "how to chunk the generator")
	commandLine.IntVar(&tchunks, "c", tchunks, "how many chunks to perform")
	commandLine.Parse(os.Args[1:])
	if tbits > 8 || tbits < 0 {
		fmt.Printf("invalid tbits: %d\n", tbits)
		os.Exit(-1)
	}
	var src []string = []string{""}
	var sidx int = 0
	var tenbits int = 1
	for i := 0; i < tbits; i++ {
		tenbits *= 10
	}
	if tchunks < 0 || tchunks*tenbits > 1000000000 {
		fmt.Printf("invalid tchunks or tchunk size: %d, %d\n", tchunks, tchunks*tenbits)
		os.Exit(-1)
	}
	fmt.Printf("Token fuzzer configured. Chunk size: %d, totalchunks: %d\n", tenbits, tchunks)
	for i := 0; i < tchunks*tenbits; i++ {
		var j int
		for tok2Ideal0(j) == "" {
			j = int(rand.Int31()) % 65536
		}
		src[sidx] += tok2Ideal(j)
		if i%tenbits == tenbits-1 {
			//fmt.Printf("%d\n", i+1)
			fmt.Printf("[BLOCK %05d] len: %d bytes\n", (i+1)/tenbits, len(src[sidx]))
			src = append(src, "")
			sidx++
		}
	}
	fmt.Printf("Collecting...\n")
	var src2 string
	for i := 0; i <= sidx; i++ {
		src2 += src[i]
	}
	ioutil.WriteFile(outfile, []byte(src2), 0644)
	fmt.Printf("Completed. Total size: %d bytes.\n", len(src2))
}
