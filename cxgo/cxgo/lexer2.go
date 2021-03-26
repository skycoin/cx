package cxgo

import (
	"fmt"
	"io"

	"github.com/skycoin/cx/cxgo/actions"
)

var DebugLexer bool

func (yylex Lexer) Error(msg string) {
	yylex.stop()
	yylex.errorf(msg)
}

func (yylex *Lexer) Lex(lval *yySymType) int {
	yylex.next()
	lval.scancopy(yylex.tok)
	actions.LineNo = lval.line
	return lval.yys
}

func (yylex *Lexer) Next() int {
	yylex.next()
	//fmt.Println(TokenName(yylex.tok.yys))
	return yylex.tok.yys
}

// func (yylex *Lexer) Stop() {
// 	// yylex.stop() //bug???? # https://github.com/skycoin/cx/issues/529
// }

func NewLexer(rdr io.Reader) *Lexer {
	lx := &Lexer{}
	lx.init(rdr, func(l, c int, msg string) {
		fmt.Printf("[%d:%d] %s\n", l, c, msg)
	})
	return lx
}

func (lval *yySymType) scancopy(tok *yySymType) {
	lval.ReturnExpressions = tok.ReturnExpressions
	lval.SelectStatement = tok.SelectStatement
	lval.SelectStatements = tok.SelectStatements
	lval.argument = tok.argument
	lval.arguments = tok.arguments
	lval.arrayArguments = tok.arrayArguments
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

//Warning Unused
//is duplicated in cxparser/cxgo0/lexer.go, also unused


