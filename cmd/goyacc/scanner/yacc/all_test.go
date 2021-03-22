// Copyright (c) 2014 The scanner Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner // import "github.com/skycoin/cx/goyacc/scanner/yacc"

import (
	"bytes"
	"fmt"
	"go/token"
	"path"
	"runtime"
	"testing"
	"unicode"
)

func dbg(s string, va ...interface{}) {
	_, fn, fl, _ := runtime.Caller(1)
	fmt.Printf("%s:%d: ", path.Base(fn), fl)
	fmt.Printf(s, va...)
	fmt.Println()
}

type row struct {
	src string
	tok Token
	lit interface{}
	num int
}

func testTokens(t *testing.T, yacc bool, table []row) {
	for i, test := range table {
		s := New(token.NewFileSet(), "testTokens", []byte(test.src))
		s.Mode(yacc)
		tok, lit, num := s.Scan()
		if g, e, g2, e2, g3, e3 := tok, test.tok, lit, test.lit, num, test.num; g != e || g2 != e2 || g3 != e3 {
			t.Errorf("%d: %s(%d) %s(%d)", i, g, g, e, e)
			t.Errorf("%d: %T(%#v) %T(%#v)", i, g2, g2, e2, e2)
			t.Fatalf("%d: %T(%#v) %T(%#v)", i, g3, g3, e3, e3)
		}
	}
}

func TestGoTokens(t *testing.T) {
	testTokens(t, false, []row{
		{"@", ILLEGAL, "@", 0}, // 0
		{"%{", REM, nil, 0},
		{"%left", REM, nil, 0},
		{"%%", REM, nil, 0},
		{"%nonassoc", REM, nil, 0},

		{"%prec", REM, nil, 0}, // 5
		{"%}", REM, nil, 0},
		{"%right", REM, nil, 0},
		{"%start", REM, nil, 0},
		{"%token", REM, nil, 0},

		{"%type", REM, nil, 0}, // 10
		{"%union", REM, nil, 0},
		{"", EOF, nil, 0},
		{"//", COMMENT, "//", 0},
		{"// ", COMMENT, "// ", 0},

		{"/**/ ", COMMENT, "/**/", 0}, // 15
		{"/***/ ", COMMENT, "/***/", 0},
		{"/** */ ", COMMENT, "/** */", 0},
		{"/* **/ ", COMMENT, "/* **/", 0},
		{"/* * */ ", COMMENT, "/* * */", 0},

		{"a", IDENT, "a", 0}, // 20
		{"ab", IDENT, "ab", 0},
		{"1", INT, uint64(1), 0},
		{"12", INT, uint64(12), 0},
		{`""`, STRING, "", 0},

		{`"1"`, STRING, "1", 0}, // 25
		{`"12"`, STRING, "12", 0},
		{"``", STRING, "", 0},
		{"`1`", STRING, "1", 0},
		{"`12`", STRING, "12", 0},

		{"'@'", CHAR, int32('@'), 0}, // 30
		{"a ", IDENT, "a", 0},
		{"ab ", IDENT, "ab", 0},
		{"1 ", INT, uint64(1), 0},
		{"12 ", INT, uint64(12), 0},

		{`"" `, STRING, "", 0}, // 35
		{`"1" `, STRING, "1", 0},
		{`"12" `, STRING, "12", 0},
		{"`` ", STRING, "", 0},
		{"`1` ", STRING, "1", 0},

		{"`12` ", STRING, "12", 0}, // 40
		{"'@' ", CHAR, int32('@'), 0},
		{" a", IDENT, "a", 0},
		{" ab", IDENT, "ab", 0},
		{" 1", INT, uint64(1), 0},

		{" 12", INT, uint64(12), 0}, // 45
		{` ""`, STRING, "", 0},
		{` "1"`, STRING, "1", 0},
		{` "12"`, STRING, "12", 0},
		{" ``", STRING, "", 0},

		{" `1`", STRING, "1", 0}, // 50
		{" `12`", STRING, "12", 0},
		{" '@'", CHAR, int32('@'), 0},
		{" a ", IDENT, "a", 0},
		{" ab ", IDENT, "ab", 0},

		{" 1 ", INT, uint64(1), 0}, // 55
		{" 12 ", INT, uint64(12), 0},
		{` "" `, STRING, "", 0},
		{` "1" `, STRING, "1", 0},
		{` "12" `, STRING, "12", 0},

		{" `` ", STRING, "", 0}, // 60
		{" `1` ", STRING, "1", 0},
		{" `12` ", STRING, "12", 0},
		{" '@' ", CHAR, int32('@'), 0},
		{"f1234567890", IDENT, "f1234567890", 0},

		{"bár", IDENT, "bár", 0}, // 65
		{"bára", IDENT, "bára", 0},
		{"123", INT, uint64(123), 0},
		{"4e6", FLOAT, 4000000., 0},
		{"42i", IMAG, 42i, 0},

		{"'@'", CHAR, int32(64), 0}, // 70
		{`"foo"`, STRING, "foo", 0},
		{"`foo`", STRING, "foo", 0},
		{"+", ADD, nil, 0},
		{"-", SUB, nil, 0},

		{"*", MUL, nil, 0}, // 75
		{"/", QUO, nil, 0},
		{"%", REM, nil, 0},
		{"&", AND, nil, 0},
		{"|", OR, nil, 0},

		{"^", XOR, nil, 0}, // 80
		{"<<", SHL, nil, 0},
		{">>", SHR, nil, 0},
		{"&^", AND_NOT, nil, 0},
		{"+=", ADD_ASSIGN, nil, 0},

		{"-=", SUB_ASSIGN, nil, 0}, // 85
		{"*=", MUL_ASSIGN, nil, 0},
		{"/=", QUO_ASSIGN, nil, 0},
		{"%=", REM_ASSIGN, nil, 0},
		{"&=", AND_ASSIGN, nil, 0},

		{"|=", OR_ASSIGN, nil, 0}, // 90
		{"^=", XOR_ASSIGN, nil, 0},
		{"<<=", SHL_ASSIGN, nil, 0},
		{">>=", SHR_ASSIGN, nil, 0},
		{"&^=", AND_NOT_ASSIGN, nil, 0},

		{"&&", LAND, nil, 0}, // 95
		{"||", LOR, nil, 0},
		{"<-", ARROW, nil, 0},
		{"++", INC, nil, 0},
		{"--", DEC, nil, 0},

		{"==", EQL, nil, 0}, // 100
		{"<", LSS, nil, 0},
		{">", GTR, nil, 0},
		{"=", ASSIGN, nil, 0},
		{"!", NOT, nil, 0},

		{"!=", NEQ, nil, 0}, // 105
		{"<=", LEQ, nil, 0},
		{">=", GEQ, nil, 0},
		{":=", DEFINE, nil, 0},
		{"...", ELLIPSIS, nil, 0},

		{"(", LPAREN, nil, 0}, // 110
		{"[", LBRACK, nil, 0},
		{"{", LBRACE, nil, 0},
		{",", COMMA, nil, 0},
		{".", PERIOD, nil, 0},

		{")", RPAREN, nil, 0}, // 115
		{"]", RBRACK, nil, 0},
		{"}", RBRACE, nil, 0},
		{";", SEMICOLON, nil, 0},
		{":", COLON, nil, 0},

		{"break", BREAK, nil, 0}, // 120
		{"case", CASE, nil, 0},
		{"chan", CHAN, nil, 0},
		{"const", CONST, nil, 0},
		{"continue", CONTINUE, nil, 0},

		{"default", DEFAULT, nil, 0}, // 125
		{"defer", DEFER, nil, 0},
		{"else", ELSE, nil, 0},
		{"fallthrough", FALLTHROUGH, nil, 0},
		{"for", FOR, nil, 0},

		{"func", FUNC, nil, 0}, // 130
		{"go", GO, nil, 0},
		{"goto", GOTO, nil, 0},
		{"if", IF, nil, 0},
		{"import", IMPORT, nil, 0},

		{"interface", INTERFACE, nil, 0}, // 135
		{"map", MAP, nil, 0},
		{"package", PACKAGE, nil, 0},
		{"range", RANGE, nil, 0},
		{"return", RETURN, nil, 0},

		{"select", SELECT, nil, 0}, // 140
		{"struct", STRUCT, nil, 0},
		{"switch", SWITCH, nil, 0},
		{"type", GO_TYPE, nil, 0},
		{"var", VAR, nil, 0},

		{"$$", DLR_DLR, nil, 0}, // 145
		{"$-2", DLR_NUM, nil, -2},
		{"$-1", DLR_NUM, nil, -1},
		{"$0", DLR_NUM, nil, 0},
		{"$1", DLR_NUM, nil, 1},

		{"$<>$", ILLEGAL, "$", 0}, // 150
		{"$<a>$", DLR_TAG_DLR, "a", 0},
		{"$<a_b>$", DLR_TAG_DLR, "a_b", 0},
		{"$<a.b>$", DLR_TAG_DLR, "a.b", 0},
		{"$<abc>$", DLR_TAG_DLR, "abc", 0},

		{"$<>1", ILLEGAL, "$", 0}, // 155
		{"$<a>-2", DLR_TAG_NUM, "a", -2},
		{"$<a_b>-1", DLR_TAG_NUM, "a_b", -1},
		{"$<a.b>0", DLR_TAG_NUM, "a.b", 0},
		{"$<abc>1", DLR_TAG_NUM, "abc", 1},

		{`'\u0061'`, CHAR, 'a', 0},     // 160
		{`'\U00000061'`, CHAR, 'a', 0}, // 160
	})
}

func TestYaccTokens(t *testing.T) {
	testTokens(t, true, []row{
		{"@", ILLEGAL, "@", 0}, // 0
		{"%{", LCURL, nil, 0},
		{"%left", LEFT, nil, 0},
		{"%%", MARK, nil, 0},
		{"%nonassoc", NONASSOC, nil, 0},

		{"%prec", PREC, nil, 0}, // 5
		{"%}", RCURL, nil, 0},
		{"%right", RIGHT, nil, 0},
		{"%start", START, nil, 0},
		{"%token", TOKEN, nil, 0},

		{"%type", TYPE, nil, 0}, // 10
		{"%union", UNION, nil, 0},
		{"", EOF, nil, 0},
		{"//", COMMENT, "//", 0},
		{"// ", COMMENT, "// ", 0},

		{"/**/ ", COMMENT, "/**/", 0}, // 15
		{"/***/ ", COMMENT, "/***/", 0},
		{"/** */ ", COMMENT, "/** */", 0},
		{"/* **/ ", COMMENT, "/* **/", 0},
		{"/* * */ ", COMMENT, "/* * */", 0},

		{"a", IDENTIFIER, "a", 0}, // 20
		{"ab", IDENTIFIER, "ab", 0},
		{"1", INT, uint64(1), 0},
		{"12", INT, uint64(12), 0},
		{`""`, STRING, "", 0},

		{`"1"`, STRING, "1", 0}, // 25
		{`"12"`, STRING, "12", 0},
		{"``", STRING, "", 0},
		{"`1`", STRING, "1", 0},
		{"`12`", STRING, "12", 0},

		{"'@'", CHAR, int32('@'), 0}, // 30
		{"a ", IDENTIFIER, "a", 0},
		{"ab ", IDENTIFIER, "ab", 0},
		{"1 ", INT, uint64(1), 0},
		{"12 ", INT, uint64(12), 0},

		{`"" `, STRING, "", 0}, // 35
		{`"1" `, STRING, "1", 0},
		{`"12" `, STRING, "12", 0},
		{"`` ", STRING, "", 0},
		{"`1` ", STRING, "1", 0},

		{"`12` ", STRING, "12", 0}, // 40
		{"'@' ", CHAR, int32('@'), 0},
		{" a", IDENTIFIER, "a", 0},
		{" ab", IDENTIFIER, "ab", 0},
		{" 1", INT, uint64(1), 0},

		{" 12", INT, uint64(12), 0}, // 45
		{` ""`, STRING, "", 0},
		{` "1"`, STRING, "1", 0},
		{` "12"`, STRING, "12", 0},
		{" ``", STRING, "", 0},

		{" `1`", STRING, "1", 0}, // 50
		{" `12`", STRING, "12", 0},
		{" '@'", CHAR, int32('@'), 0},
		{" a ", IDENTIFIER, "a", 0},
		{" ab ", IDENTIFIER, "ab", 0},

		{" 1 ", INT, uint64(1), 0}, // 55
		{" 12 ", INT, uint64(12), 0},
		{` "" `, STRING, "", 0},
		{` "1" `, STRING, "1", 0},
		{` "12" `, STRING, "12", 0},

		{" `` ", STRING, "", 0}, // 60
		{" `1` ", STRING, "1", 0},
		{" `12` ", STRING, "12", 0},
		{" '@' ", CHAR, int32('@'), 0},
		{"f1234567890", IDENTIFIER, "f1234567890", 0},

		{"bár", IDENTIFIER, "bár", 0}, // 65
		{"bára", IDENTIFIER, "bára", 0},
		{"123", INT, uint64(123), 0},
		{"4e6", INT, uint64(4), 0},
		{"42i", INT, uint64(42), 0},

		{"'@'", CHAR, int32(64), 0}, // 70
		{`"foo"`, STRING, "foo", 0},
		{"`foo`", STRING, "foo", 0},
		{"+", ILLEGAL, "+", 0},
		{"-", ILLEGAL, "-", 0},

		{"*", ILLEGAL, "*", 0}, // 75
		{"/", ILLEGAL, "/", 0},
		{"%", ILLEGAL, "%", 0},
		{"&", ILLEGAL, "&", 0},
		{"|", ILLEGAL, "|", 0},

		{"^", ILLEGAL, "^", 0}, // 80
		{"<<", ILLEGAL, "<", 0},
		{">>", ILLEGAL, ">", 0},
		{"&^", ILLEGAL, "&", 0},
		{"+=", ILLEGAL, "+", 0},

		{"-=", ILLEGAL, "-", 0}, // 85
		{"*=", ILLEGAL, "*", 0},
		{"/=", ILLEGAL, "/", 0},
		{"%=", ILLEGAL, "%", 0},
		{"&=", ILLEGAL, "&", 0},

		{"|=", ILLEGAL, "|", 0}, // 90
		{"^=", ILLEGAL, "^", 0},
		{"<<=", ILLEGAL, "<", 0},
		{">>=", ILLEGAL, ">", 0},
		{"&^=", ILLEGAL, "&", 0},

		{"&&", ILLEGAL, "&", 0}, // 95
		{"||", ILLEGAL, "|", 0},
		{"<-", ILLEGAL, "<", 0},
		{"++", ILLEGAL, "+", 0},
		{"--", ILLEGAL, "-", 0},

		{"==", ILLEGAL, "=", 0}, // 100
		{"<", ILLEGAL, "<", 0},
		{">", ILLEGAL, ">", 0},
		{"=", ILLEGAL, "=", 0},
		{"!", ILLEGAL, "!", 0},

		{"!=", ILLEGAL, "!", 0}, // 105
		{"<=", ILLEGAL, "<", 0},
		{">=", ILLEGAL, ">", 0},
		{":=", ILLEGAL, ":", 0},
		{"...", ILLEGAL, ".", 0},

		{"(", ILLEGAL, "(", 0}, // 110
		{"[", ILLEGAL, "[", 0},
		{"{", ILLEGAL, "{", 0},
		{",", COMMA, nil, 0},
		{".", ILLEGAL, ".", 0},

		{")", ILLEGAL, ")", 0}, // 115
		{"]", ILLEGAL, "]", 0},
		{"}", ILLEGAL, "}", 0},
		{";", ILLEGAL, ";", 0},
		{":", ILLEGAL, ":", 0},

		{"break", IDENTIFIER, "break", 0}, // 120
		{"case", IDENTIFIER, "case", 0},
		{"chan", IDENTIFIER, "chan", 0},
		{"const", IDENTIFIER, "const", 0},
		{"continue", IDENTIFIER, "continue", 0},

		{"default", IDENTIFIER, "default", 0}, // 125
		{"defer", IDENTIFIER, "defer", 0},
		{"else", IDENTIFIER, "else", 0},
		{"fallthrough", IDENTIFIER, "fallthrough", 0},
		{"for", IDENTIFIER, "for", 0},

		{"func", IDENTIFIER, "func", 0}, // 130
		{"go", IDENTIFIER, "go", 0},
		{"goto", IDENTIFIER, "goto", 0},
		{"if", IDENTIFIER, "if", 0},
		{"import", IDENTIFIER, "import", 0},

		{"interface", IDENTIFIER, "interface", 0}, // 135
		{"map", IDENTIFIER, "map", 0},
		{"package", IDENTIFIER, "package", 0},
		{"range", IDENTIFIER, "range", 0},
		{"return", IDENTIFIER, "return", 0},

		{"select", IDENTIFIER, "select", 0}, // 140
		{"struct", IDENTIFIER, "struct", 0},
		{"switch", IDENTIFIER, "switch", 0},
		{"type", IDENTIFIER, "type", 0},
		{"var", IDENTIFIER, "var", 0},

		// ----

		{"a.foo", IDENTIFIER, "a.foo", 0}, // 145
		{"b.fooára", IDENTIFIER, "b.fooára", 0},
		{"ab.foo", IDENTIFIER, "ab.foo", 0},
		{"a.foo ", IDENTIFIER, "a.foo", 0},
		{"ab.foo ", IDENTIFIER, "ab.foo", 0},

		{" a.foo", IDENTIFIER, "a.foo", 0}, // 150
		{" ab.foo", IDENTIFIER, "ab.foo", 0},
		{" a.foo ", IDENTIFIER, "a.foo", 0},
		{" ab.foo ", IDENTIFIER, "ab.foo", 0},
		{"f1234567890.foo", IDENTIFIER, "f1234567890.foo", 0},

		{"b.fooár", IDENTIFIER, "b.fooár", 0}, // 155
		{"$$", ILLEGAL, "$", 0},
		{"$-2", ILLEGAL, "$", 0},
		{"$-1", ILLEGAL, "$", 0},
		{"$0", ILLEGAL, "$", 0},

		{"$1", ILLEGAL, "$", 0}, // 160
		{"$<>$", ILLEGAL, "$", 0},
		{"$<a>$", ILLEGAL, "$", 0},
		{"$<a_b>$", ILLEGAL, "$", 0},
		{"$<a.b>$", ILLEGAL, "$", 0},

		{"$<abc>$", ILLEGAL, "$", 0}, // 165
		{"$<>1", ILLEGAL, "$", 0},
		{"$<a>-2", ILLEGAL, "$", 0},
		{"$<a_b>-1", ILLEGAL, "$", 0},
		{"$<a.b>0", ILLEGAL, "$", 0},

		{"$<abc>1", ILLEGAL, "$", 0}, // 170

		// --

		{"a:", C_IDENTIFIER, "a", 0},
		{"a: ", C_IDENTIFIER, "a", 0},
		{"a :", C_IDENTIFIER, "a", 0},
		{"a : ", C_IDENTIFIER, "a", 0},

		{" a:", C_IDENTIFIER, "a", 0}, // 175
		{" a: ", C_IDENTIFIER, "a", 0},
		{" a :", C_IDENTIFIER, "a", 0},
		{" a : ", C_IDENTIFIER, "a", 0},
		{"ab:", C_IDENTIFIER, "ab", 0},

		{"ab: ", C_IDENTIFIER, "ab", 0}, // 180
		{"ab :", C_IDENTIFIER, "ab", 0},
		{"ab : ", C_IDENTIFIER, "ab", 0},
		{" ab:", C_IDENTIFIER, "ab", 0},
		{" ab: ", C_IDENTIFIER, "ab", 0},

		{" ab :", C_IDENTIFIER, "ab", 0}, // 185
		{" ab : ", C_IDENTIFIER, "ab", 0},
		{"a.b:", C_IDENTIFIER, "a.b", 0},
		{"a.b: ", C_IDENTIFIER, "a.b", 0},
		{"a.b :", C_IDENTIFIER, "a.b", 0},

		{"a.b : ", C_IDENTIFIER, "a.b", 0}, // 190
		{" a.b:", C_IDENTIFIER, "a.b", 0},
		{" a.b: ", C_IDENTIFIER, "a.b", 0},
		{" a.b :", C_IDENTIFIER, "a.b", 0},
		{" a.b : ", C_IDENTIFIER, "a.b", 0},
	})
}

func TestBug(t *testing.T) {
	tab := []struct {
		src  string
		toks []Token
	}{
		{`%left`, []Token{LEFT}},
		{`%left %left`, []Token{LEFT, LEFT}},
		{`%left 'a' %left`, []Token{LEFT, CHAR, LEFT}},
		{`foo`, []Token{IDENTIFIER}},
		{`foo bar`, []Token{IDENTIFIER, IDENTIFIER}},
		{`foo bar baz`, []Token{IDENTIFIER, IDENTIFIER, IDENTIFIER}},
		{`%token <ival> DREG VREG`, []Token{TOKEN, ILLEGAL, IDENTIFIER, ILLEGAL, IDENTIFIER, IDENTIFIER}},
	}

	for i, test := range tab {
		s := New(token.NewFileSet(), "TestBug", []byte(test.src))
		s.Mode(true)
		for j, etok := range test.toks {
			tok, _, _ := s.Scan()
			if g, e := tok, etok; g != e {
				t.Errorf("%d.%d: %s(%d) %s(%d)", i, j, g, g, e, e)
			}
		}
	}
}

func TestChar(t *testing.T) {
	var buf bytes.Buffer
	for r := rune(0); r <= unicode.MaxRune; r++ {
		if r >= 0xd800 && r <= 0xdfff {
			continue
		}

		buf.WriteString(fmt.Sprintf("'\\U%08x'\n", r))
	}
	s := New(token.NewFileSet(), "TestChar", buf.Bytes())
	for r := rune(0); r <= unicode.MaxRune; r++ {
		if r >= 0xd800 && r <= 0xdfff {
			continue
		}

		tok, val, _ := s.Scan()
		if len(s.Errors) != 0 {
			t.Fatalf("%#x err: %v", r, s.Errors)
		}

		if g, e := tok, CHAR; g != e {
			t.Fatalf("%#x tok: %v %v", r, g, e)
		}

		if g, e := val, r; g != e {
			t.Fatalf("%#x val: %T(%v) %T(%v)", r, g, g, e, e)
		}
	}
	tok, _, _ := s.Scan()
	if len(s.Errors) != 0 {
		t.Fatalf("err: %v", s.Errors)
	}

	if g, e := tok, EOF; g != e {
		t.Fatalf("tok: %v %v", g, e)
	}
}
