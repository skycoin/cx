%{

// Copyright (c) 2014 The scanner Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// 
// CAUTION: If this file is a Go source file (*.go), it was generated
// automatically by '$ golex' from a *.l file - DO NOT EDIT in that case!

// Package scanner implements a scanner for yacc[0] source text with actions
// written in Go. It takes a []byte as source which can then be tokenized
// through repeated calls to the Scan method.
//
// Changelog
//
// 2014-12-18: Support %precedence for bison compatibility[1].
//
// 2014-11-17: Allow dashes in symbol names for bison compatibility[1].
//
// Links
// 
// Referenced from above:
// 
//  [0]: http://pubs.opengroup.org/onlinepubs/009695399/utilities/yacc.html
//  [1]: http://www.gnu.org/software/bison/manual/html_node/Symbols.html
package scanner

import (
	"errors"
	"fmt"
	"go/token"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Token int

const (
	_INITIAL = iota
	_GO
	_S1
	_S2
)

const (
	_ Token = 0xE000+iota

	// ------------------------------------------- Go mode or shared tokens

	// Special tokens
	ILLEGAL
	EOF
	COMMENT
	
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	IMAG   // 123.45i
	CHAR   // 'a'
	STRING // "abc"
	
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %
	
	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^
	
	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=
	
	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=
	
	LAND  // &&
	LOR   // ||
	ARROW // <-
	INC   // ++
	DEC   // --
	
	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !
	
	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...
	
	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .
	
	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :
	
	// Keywords
	BREAK
	CASE
	CHAN
	CONST
	CONTINUE
	
	DEFAULT
	DEFER
	ELSE
	FALLTHROUGH
	FOR
	
	FUNC
	GO
	GOTO
	IF
	IMPORT
	
	INTERFACE
	MAP
	PACKAGE
	RANGE
	RETURN
	
	SELECT
	STRUCT
	SWITCH
	GO_TYPE
	VAR

	DLR_DLR      // $$
	DLR_NUM      // $1
	DLR_TAG_DLR  // $<tag>$
	DLR_TAG_NUM  // $<tag>2

	// yacc mode tokens
	C_IDENTIFIER // IDENTIFIER ':'
	ERR_VERBOSE  // %error-verbose
	IDENTIFIER   // [a-zA-Z_][a-zA-Z0-9_.-]*
	LCURL        // %{
	LEFT         // %left
	MARK         // %%
	NONASSOC     // %nonassoc
	PREC         // %prec
	RCURL        // %}
	RIGHT        // %right
	START        // %start
	TOKEN        // %token
	TYPE         // %type
	UNION        // %union
	PRECEDENCE   // %precedence
)

var ts = map[Token]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",
	
	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	IMAG:   "IMAG",
	CHAR:   "CHAR",
	STRING: "STRING",
	
	ADD: "ADD",
	SUB: "SUB",
	MUL: "MUL",
	QUO: "QUO",
	REM: "REM",
	
	AND:     "AND",
	OR:      "OR",
	XOR:     "XOR",
	SHL:     "SHL",
	SHR:     "SHR",
	AND_NOT: "AND_NOT",
	
	ADD_ASSIGN: "ADD_ASSIGN",
	SUB_ASSIGN: "SUB_ASSIGN",
	MUL_ASSIGN: "MUL_ASSIGN",
	QUO_ASSIGN: "QUO_ASSIGN",
	REM_ASSIGN: "REM_ASSIGN",
	
	AND_ASSIGN:     "AND_ASSIGN",
	OR_ASSIGN:      "OR_ASSIGN",
	XOR_ASSIGN:     "XOR_ASSIGN",
	SHL_ASSIGN:     "SHL_ASSIGN",
	SHR_ASSIGN:     "SHR_ASSIGN",
	AND_NOT_ASSIGN: "AND_NOT_ASSIGN",
	
	LAND:  "LAND",
	LOR:   "LOR",
	ARROW: "ARROW",
	INC:   "INC",
	DEC:   "DEC",
	
	EQL:    "EQL",
	LSS:    "LSS",
	GTR:    "GTR",
	ASSIGN: "ASSIGN",
	NOT:    "NOT",
	
	NEQ:      "NEQ",
	LEQ:      "LEQ",
	GEQ:      "GEQ",
	DEFINE:   "DEFINE",
	ELLIPSIS: "ELLIPSIS",
	
	LPAREN: "LPAREN",
	LBRACK: "LBRACK",
	LBRACE: "LBRACE",
	COMMA:  "COMMA",
	PERIOD: "PERIOD",
	
	RPAREN:    "RPAREN",
	RBRACK:    "RBRACK",
	RBRACE:    "RBRACE",
	SEMICOLON: "SEMICOLON",
	COLON:      "COLON",
	
	BREAK:    "BREAK",
	CASE:     "CASE",
	CHAN:     "CHAN",
	CONST:    "CONST",
	CONTINUE: "CONTINUE",
	
	DEFAULT:     "DEFAULT",
	DEFER:       "DEFER",
	ELSE:        "ELSE",
	FALLTHROUGH: "FALLTHROUGH",
	FOR:         "FOR",
	
	FUNC:   "FUNC",
	GO:     "GO",
	GOTO:   "GOTO",
	IF:     "IF",
	IMPORT: "IMPORT",
	
	INTERFACE: "INTERFACE",
	MAP:       "MAP",
	PACKAGE:   "PACKAGE",
	RANGE:     "RANGE",
	RETURN:    "RETURN",
	
	SELECT:  "SELECT",
	STRUCT:  "STRUCT",
	SWITCH:  "SWITCH",
	GO_TYPE: "GO_TYPE",
	VAR:     "VAR",

	DLR_DLR:     "DLR_DLR",
	DLR_NUM:     "DLR_NUM",
	DLR_TAG_DLR: "DLR_TAG_DLR",
	DLR_TAG_NUM: "DLR_TAG_NUM",

	// --------------------------------------------- yacc/bison mode tokens

	C_IDENTIFIER: "C_IDENTIFIER",
	IDENTIFIER:   "IDENTIFIER",
	LCURL:        "LCURL",
	LEFT:         "LEFT",
	MARK:         "MARK",

	NONASSOC:     "NONASSOC",
	PREC:         "PREC",
	PRECEDENCE:   "PRECEDENCE",
	RCURL:        "RCURL",
	RIGHT:        "RIGHT",

	START:        "START",
	TOKEN:        "TOKEN",
	TYPE:         "TYPE",
	UNION:        "UNION",
	ERR_VERBOSE:  "ERR_VERBOSE",
}

// String implements fmt.Stringer
func (i Token) String() string {
	if s := ts[i]; s != "" {
		return s
	}

	return fmt.Sprintf("Token(%d)", int(i))
}

// A Scanner holds the scanner's internal state while processing a given text.
type Scanner struct {
	Col     int     // Starting column of the last scanned token.
	Errors  []error // List of accumulated errors.
	FName   string  // File name (reported) of the scanned source.
	Line    int     // Starting line of the last scanned token.
	NCol    int     // Starting column (reported) for the next scanned token.
	NLine   int     // Starting line (reported) for the next scanned token.
	c       int
	file    *token.File
	i       int
	i0      int
	sc      int
	src     []byte
	ssc     int     // saved state condition
	val     []byte
}

// New returns a newly created Scanner and set its FName to fname
func New(fset *token.FileSet, fname string, src []byte) (s *Scanner) {
	if len(src) > 2 && src[0] == 0xEF && src[1] == 0xBB && src[2] == 0xBF {
		src = src[3:]
	}
	s = &Scanner{
		FName: fname,
		src:   src,
		NLine: 1,
		NCol:  0,
	}
	s.file = fset.AddFile(fname, -1, len(src))
	s.next()
	return
}

func (s *Scanner) next() int {
	if s.c != 0 {
		s.val = append(s.val, byte(s.c))
	}
	s.c = 0
	if s.i < len(s.src) {
		s.c = int(s.src[s.i])
		s.i++
	}
	switch s.c {
	case '\n':
		s.NLine++
		s.NCol = 0
		s.file.AddLine(s.i)
	default:
		s.NCol++
	}
	return s.c
}

// Pos returns the starting offset of the last scanned token.
func (s *Scanner) Pos() int {
	return s.i0
}

func (s *Scanner) err(format string, arg ...interface{}) {
	err := fmt.Errorf(fmt.Sprintf("%s:%d:%d ", s.FName, s.Line, s.Col)+format, arg...)
	s.Errors = append(s.Errors, err)
}

// Error implements yyLexer.
func (s *Scanner) Error(msg string) {
	switch msg {
	case "syntax error":
		s.err(msg)
	default:
		s.Errors = append(s.Errors, errors.New(msg))
	}
}

// Mode allows to switch the scanner mode from scanning yacc tokens to scanning
// rule action tokens (Go tokens). Initially the yacc mode is active.
func (s *Scanner) Mode(yacc bool) {
	switch yacc {
	case true:
		s.sc = _INITIAL
	case false:
		s.sc = _GO
	}
}

// Scan works like ScanRaw but recognizes also yacc's C_IDENTIFIER token (in
// yacc mode).
func (s *Scanner) Scan() (tok Token, lval interface{}, num int) {
	tok, lval, num = s.ScanRaw()
	if s.sc != _INITIAL || tok != IDENTIFIER {
		return
	}

	idLine, idCol := s.Line, s.Col
	i, nl, nc, c := s.i, s.NLine, s.NCol, s.c
	i0 := s.i0
	tok2, lit, _ := s.ScanRaw()
	s.i0 = i0
	s.Line, s.Col = idLine, idCol
	if tok2 == ILLEGAL && lit.(string) == ":" {
		return C_IDENTIFIER, lval, 0	
	}

	s.i, s.NLine, s.NCol, s.c = i, nl, nc, c
	return
}

// ScanRaw scans the next token and returns the token and its value if
// applicable.  The source end is indicated by EOF.
// 
// If the returned token is IDENT, INT, FLOAT, IMAG, CHAR, STRING or COMMENT,
// lval has has the corresponding value - not the string representation of the
// value. However, numeric literals which overflow the corresponding Go
// predeclared types are returned as string.
//
// If the returned token is DLR_NUM or DLR_TAG_DLR, num has the corresponding
// number (int) or lval has the corresponding tag (string).
//
// If the returned token is DLR_TAG_NUM, lval has the corresponding tag (string)
// and num has the corresponding number.
// 
// If the returned token is ILLEGAL, the literal string is the offending
// character or number/string/char literal.
func (s *Scanner) ScanRaw() (tok Token, lval interface{}, num int) {
	//defer func() { fmt.Printf("%s(%d) %v\n", tok, int(tok), lval) }()
	c0, c := s.c, s.c
%}

int_lit         {decimal_lit}|{octal_lit}|{hex_lit}
decimal_lit     [1-9][0-9]*
octal_lit       0[0-7]*
hex_lit         0[xX][0-9a-fA-F]+

float_lit       {D}"."{D}?{E}?|{D}{E}|"."{D}{E}?
D        	[0-9]+
E        	[eE][-+]?[0-9]+

imaginary_ilit  {D}i
imaginary_lit   {float_lit}i
utf8_h2         [\xC2-\xDF]
utf8_h3         [\xE0-\xEF]
utf8_h4         [\xF0-\xF4]
utf8_cont       [\x80-\xBF]
utf8_2          {utf8_h2}{utf8_cont}
utf8_3          {utf8_h3}{utf8_cont}{utf8_cont}
utf8_4          {utf8_h4}{utf8_cont}{utf8_cont}{utf8_cont}
non_ascii       {utf8_2}|{utf8_3}|{utf8_4}

a		[aA]
b		[bB]
c		[cC]
d		[dD]
e		[eE]
f		[fF]
g		[gG]
h		[hH]
i		[iI]
k		[kK]
l		[lL]
n		[nN]
o		[oO]
p		[pP]
r		[rR]
s		[sS]
t		[tT]
u		[uU]
v		[vV]
y		[yY]

errorverbose	{e}{r}{r}{o}{r}-{v}{e}{r}{b}{o}{s}{e}
left		{l}{e}{f}{t}
nonassoc	{n}{o}{n}{a}{s}{s}{o}{c}
prec		{p}{r}{e}{c}
precedence      {p}{r}{e}{c}{e}{d}{e}{n}{c}{e}
right		{r}{i}{g}{h}{t}
start		{s}{t}{a}{r}{t}
token		{t}{o}{k}{e}{n}
type		{t}{y}{p}{e}
union		{u}{n}{i}{o}{n}

%yyc c
%yyn c = s.next()
%yyt s.sc

%x _GO _S1 _S2

%%
			s.val = s.val[:0]
			s.i0, s.Line, s.Col, c0 = s.i, s.NLine, s.NCol, c
			
<*>\0			s.i0++
			return EOF, lval, 0

<INITIAL,_GO>[ \t\n\r]+

<INITIAL,_GO>\/\*([^*]|\*+[^*/])*\*+\/
			|
<INITIAL,_GO>\/\/.*	return COMMENT, string(s.val), 0

<_GO>$-?{D}
	v := s.val[1:]
	neg := v[0] == '-'
	if neg {
		v = v[1:]
	}
	n, err := strconv.ParseUint(string(v), 0, 32)
	if err != nil {
		fmt.Printf("%q\n", v)
		return ILLEGAL, nil, 0
	}

	num = int(n)
	if neg {
		num *= -1
	}
	return DLR_NUM, lval, num

<_GO>$<[a-zA-Z_][a-zA-Z0-9_.]*>\$
	return DLR_TAG_DLR, string(s.val[2:len(s.val)-2]), 0

<_GO>$<[a-zA-Z_][a-zA-Z0-9_.]*>-?{D}
	a := strings.Split(string(s.val[2:]), ">")
	v := a[1]
	neg := v[0] == '-'
	if neg {
		v = v[1:]
	}
	n, err := strconv.ParseUint(string(v), 0, 32)
	if err != nil {
		fmt.Printf("%q\n", v)
		return ILLEGAL, nil, 0
	}

	num = int(n)
	if neg {
		num *= -1
	}
	return DLR_TAG_NUM, a[0], num

<_GO>"$$"		return DLR_DLR, lval, 0
<_GO>"!"		return NOT, lval, 0
<_GO>"!="		return NEQ, lval, 0
<_GO>"%"		return REM, lval, 0
<_GO>"%="		return REM_ASSIGN, lval, 0
<_GO>"&"		return AND, lval, 0
<_GO>"&&"		return LAND, lval, 0
<_GO>"&="		return AND_ASSIGN, lval, 0
<_GO>"&^"		return AND_NOT, lval, 0
<_GO>"&^="		return AND_NOT_ASSIGN, lval, 0
<_GO>"("		return LPAREN, lval, 0
<_GO>")"		return RPAREN, lval, 0
<_GO>"*"		return MUL, lval, 0
<_GO>"*="		return MUL_ASSIGN, lval, 0
<_GO>"+"		return ADD, lval, 0
<_GO>"++"		return INC, lval, 0
<_GO>"+="		return ADD_ASSIGN, lval, 0
<_GO>","		return COMMA, lval, 0
<_GO>"-"		return SUB, lval, 0
<_GO>"--"		return DEC, lval, 0
<_GO>"-="		return SUB_ASSIGN, lval, 0
<_GO>"."		return PERIOD, lval, 0
<_GO>"..."		return ELLIPSIS, lval, 0
<_GO>"/"		return QUO, lval, 0
<_GO>"/="		return QUO_ASSIGN, lval, 0
<_GO>":"		return COLON, lval, 0
<_GO>":="		return DEFINE, lval, 0
<_GO>";"		return SEMICOLON, lval, 0
<_GO>"<"		return LSS, lval, 0
<_GO>"<-"		return ARROW, lval, 0
<_GO>"<<"		return SHL, lval, 0
<_GO>"<<="		return SHL_ASSIGN, lval, 0
<_GO>"<="		return LEQ, lval, 0
<_GO>"="		return ASSIGN, lval, 0
<_GO>"=="		return EQL, lval, 0
<_GO>">"		return GTR, lval, 0
<_GO>">="		return GEQ, lval, 0
<_GO>">>"		return SHR, lval, 0
<_GO>">>="		return SHR_ASSIGN, lval, 0
<_GO>"["		return LBRACK, lval, 0
<_GO>"]"		return RBRACK, lval, 0
<_GO>"^"		return XOR, lval, 0
<_GO>"^="		return XOR_ASSIGN, lval, 0
<_GO>"{"		return LBRACE, lval, 0
<_GO>"|"		return OR, lval, 0
<_GO>"|="		return OR_ASSIGN, lval, 0
<_GO>"||"		return LOR, lval, 0
<_GO>"}"		return RBRACE, lval, 0

<_GO>break		return BREAK, lval, 0
<_GO>case		return CASE, lval, 0
<_GO>chan		return CHAN, lval, 0
<_GO>const		return CONST, lval, 0
<_GO>continue		return CONTINUE, lval, 0
<_GO>default		return DEFAULT, lval, 0
<_GO>defer		return DEFER, lval, 0
<_GO>else		return ELSE, lval, 0
<_GO>fallthrough	return FALLTHROUGH, lval, 0
<_GO>for		return FOR, lval, 0
<_GO>func		return FUNC, lval, 0
<_GO>go			return GO, lval, 0
<_GO>goto		return GOTO, lval, 0
<_GO>if			return IF, lval, 0
<_GO>import		return IMPORT, lval, 0
<_GO>interface		return INTERFACE, lval, 0
<_GO>map		return MAP, lval, 0
<_GO>package		return PACKAGE, lval, 0
<_GO>range		return RANGE, lval, 0
<_GO>return		return RETURN, lval, 0
<_GO>select		return SELECT, lval, 0
<_GO>struct		return STRUCT, lval, 0
<_GO>switch		return SWITCH, lval, 0
<_GO>type		return GO_TYPE, lval, 0
<_GO>var		return VAR, lval, 0

<_GO>{imaginary_ilit}	return s.int(IMAG)
<_GO>{imaginary_lit}	return s.float(IMAG)
<INITIAL,_GO>{int_lit}	return s.int(INT)
<_GO>{float_lit}	return s.float(FLOAT)

<INITIAL,_GO>\"		s.ssc, s.sc = s.sc, _S1
<INITIAL,_GO>`		s.ssc, s.sc = s.sc, _S2

<INITIAL,_GO>''			|
<INITIAL,_GO>'(\\.)?[^']*	return CHAR, string(s.val), 0

<INITIAL,_GO>'(\\.)?[^']*'	
			s.ssc = s.sc
			if tok, lval, _ = s.str(""); tok != STRING {
				return
			}
			s.i0++
			return CHAR, []rune(lval.(string))[0], 0

<_S1>(\\.|[^\\"])*\"	return s.str("\"")
<_S2>([^`]|\n)*`	return s.str("`")

"%{"			return LCURL, lval, 0
"%}"			return RCURL, lval, 0
%%			return MARK, lval, 0
%{errorverbose}		return ERR_VERBOSE, lval, 0
%{left}			return LEFT, lval, 0
%{nonassoc}		return NONASSOC, lval, 0
%{prec}			return PREC, lval, 0
%{precedence}		return PRECEDENCE, lval, 0
%{right}		return RIGHT, lval, 0
%{start}		return START, lval, 0
%{token}		return TOKEN, lval, 0
%{type}			return TYPE, lval, 0
%{union}		return UNION, lval, 0
,			return COMMA, lval, 0

<_GO>[a-zA-Z_][a-zA-Z0-9_]*
	if c >= '\xC2' && c <= '\xF4' {
		s.i--
		s.NCol--
		for rune := rune(1) ; rune >= 0; rune = s.getRune(true, false) {
		tok, lval = IDENT, string(s.src[s.i0-1:s.i])
		}
		s.next()
		return
	}

	return IDENT, string(s.val), 0

<INITIAL>[a-zA-Z_][a-zA-Z0-9_.-]*
	if c >= '\xC2' && c <= '\xF4' {
		s.i--
		s.NCol--
		for rune := rune(1) ; rune >= 0; rune = s.getRune(true, true) {
		tok, lval = IDENTIFIER, string(s.src[s.i0-1:s.i])
		}
		s.next()
		return
	}

	return IDENTIFIER, string(s.val), 0

<_GO>{non_ascii}
	s.i = s.i0-1
	if rune := s.getRune(false, false); rune < 0 {
		_, sz := utf8.DecodeRune(s.src[s.i:])
		s.i += sz
		s.next()
		s.err("expected unicode lettter, got %U", rune)
		return ILLEGAL, string(-rune), 0
	}

	for rune := rune(1) ; rune >= 0; rune = s.getRune(true, false) {
	}
	s.next()
	return IDENT, string(s.src[s.i0-1:s.i-1]), 0

<INITIAL>{non_ascii}
	s.i = s.i0-1
	if rune := s.getRune(false, false); rune < 0 {
		_, sz := utf8.DecodeRune(s.src[s.i:])
		s.i += sz
		s.next()
		s.err("expected unicode lettter, got %U", rune)
		return ILLEGAL, string(-rune), 0
	}

	for rune := rune(1) ; rune >= 0; rune = s.getRune(true, true) {
	}
	s.next()
	return IDENTIFIER, string(s.src[s.i0-1:s.i-1]), 0

%%
	s.next()
	return ILLEGAL, string(c0), 0
}

func (s *Scanner) getRune(acceptDigits, acceptDot bool) (r rune) {
	var sz int
	if r, sz = utf8.DecodeRune(s.src[s.i:]); sz != 0 && 
		(r == '_' || unicode.IsLetter(r) || (acceptDigits && unicode.IsDigit(r)) || (acceptDot && r == '.')) {
		s.i += sz
		s.NCol += sz
		return
	}

	return -r
}

func (s *Scanner) str(pref string) (tok Token, lval interface{}, num int) {
	s.sc = s.ssc
	ss := pref + string(s.val)
	ss, err := strconv.Unquote(ss)
	if err != nil {
		s.err("string literal %q: %v", ss, err)
		return ILLEGAL, ss, 0
	}

	s.i0--
	return STRING, ss, 0
}

func (s *Scanner) int(tk Token) (tok Token, lval interface{}, num int) {
	tok = tk
	if tok == IMAG {
		s.val = s.val[:len(s.val)-1]
	}
	n, err := strconv.ParseUint(string(s.val), 0, 64)
	if err != nil {
		lval = string(s.val)
	} else if tok == IMAG {
		lval = complex(0, float64(n))
	} else {
		lval = n
	}
	return
}

func (s *Scanner) float(tk Token) (tok Token, lval interface{}, num int) {
	tok = tk
	if tok == IMAG {
		s.val = s.val[:len(s.val)-1]
	}
	n, err := strconv.ParseFloat(string(s.val), 64)
	if err != nil {
		lval = string(s.val)
	} else if tok == IMAG {
		lval = complex(0, n)
	} else {
		lval = n
	}
	return
}
