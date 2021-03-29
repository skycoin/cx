package cxgo0

import (
	"fmt"
	"github.com/skycoin/cx/cx/constants"

	//"github.com/skycoin/cx/cxgo/globals"
	"io"
	"os"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	l, c      int    //line and column numbers
	b, r, e   int    //used for buffer mechanics
	buf       []byte //buffer
	scan      io.Reader
	ch        rune //most recently read character
	chw       int  //width of character, in bytes
	errh      func(l, c int, msg string)
	ioerr     error
	bsize     uint //buffer size
	nlsemi    bool //whether or not newline is SCOLON
	eof       bool //eof
	crash     bool //used for crash behaviour
	colbefore bool //used for colon keywords

	tok *yySymType //symbol read. soon to be depracated for fully new cxgo
}

const sentinel = utf8.RuneSelf
const LexerBufferMin = 12
const LexerBufferMax = 20
const readCountMax = 10

/* the following code is provided */
func (s *Lexer) init(r io.Reader, errh func(l, c int, msg string)) {
	s.scan = r
	s.l, s.c, s.b, s.r, s.e = 0, 0, -1, 0, 0
	s.buf = make([]byte, 1<<LexerBufferMin)
	s.buf[0] = sentinel
	s.ch = ' '
	s.errh = errh
	s.chw = -1
	s.bsize = LexerBufferMin
}

func (s *Lexer) start() { s.b = s.r - s.chw }
func (s *Lexer) stop()  { s.b = -1 }
func (s *Lexer) segment() []byte {
	return s.buf[s.b : s.r-s.chw]
}

func (s *Lexer) errorf(msg string) {
	s.errh(s.l+1, s.c+1, msg)
	//panic("")
}

func (s *Lexer) rewind() {
	// ok to verify precondition - rewind is rarely called
	if s.b < 0 {
		panic("no active segment")
	}
	s.c -= s.r - s.b
	s.r = s.b
	s.nextch()
}

func (s *Lexer) nextch() {
redo:
	s.c += int(s.chw)
	if s.ch == '\n' {
		s.l++
		s.c = 0
	}

	//first test for ASCII
	if s.ch = rune(s.buf[s.r]); s.ch < sentinel {
		s.r++
		s.chw = 1
		if s.ch == 0 {
			s.errorf("NUL")
			goto redo
		}
		return
	}

	for s.e-s.r < utf8.UTFMax && !utf8.FullRune(s.buf[s.r:s.e]) && s.ioerr == nil {
		s.fill()
	}

	//EOF
	if s.r == s.e || s.ioerr == io.EOF {
		if s.ioerr != io.EOF {
			s.errorf("IO ProgramError: " + s.ioerr.Error())
			s.ioerr = nil
		}
		s.ch = -1
		s.chw = 0
		return
	}

	s.ch, s.chw = utf8.DecodeRune(s.buf[s.r:s.e])
	s.r += s.chw
	if s.ch == utf8.RuneError && s.chw == 1 {
		s.errorf("invalid UTF-8 encoding!")
		goto redo
	}

	//WATCH OUT FOR BOM
	if s.ch == 0xfeff {
		if s.l > 0 || s.c > 0 {
			s.errorf("invalid UFT-8 byte-order mark in middle of file")
		}
		goto redo
	}
}

func (s *Lexer) fill() {
	b := s.r
	if s.b >= 0 {
		b = s.b
		s.b = 0
	}
	content := s.buf[b:s.e]
	if len(content)*2 > len(s.buf) {
		s.bsize++
		if s.bsize > LexerBufferMax {
			s.bsize = LexerBufferMax
		}
		s.buf = make([]byte, 1<<s.bsize)
		copy(s.buf, content)
	} else if b > 0 {
		copy(s.buf, content)
	}
	s.r -= b
	s.e -= b

	for i := 0; i < readCountMax; i++ {
		var n int
		n, s.ioerr = s.scan.Read(s.buf[s.e : len(s.buf)-1])
		if n < 0 {
			panic("negative read!") //invalid io.Reader
		}
		if n > 0 || s.ioerr != nil {
			s.e += n
			s.buf[s.e] = sentinel
			return
		}
	}

	s.buf[s.e] = sentinel
	s.ioerr = io.ErrNoProgress
}

func (s *Lexer) next() {
	nlsemi := s.nlsemi
	s.nlsemi = false
	if s.eof {
		if s.crash {
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
		s.tok.yys = -1
		return
	}
redonext:
	s.stop()
	s.tok = &yySymType{}
	s.tok.line = s.l + 1
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' && !nlsemi || s.ch == '\r' {
		s.nextch()
	}

	s.start()
	if isLetter(s.ch) || s.ch >= utf8.RuneSelf && s.atIdentChar(true) {
		s.nextch()
		s.ident()
		return
	}

	switch s.ch {
	case -1, 0:
		s.eof = true
		s.tok.yys = -1
		if s.crash {
			os.Exit(constants.CX_COMPILATION_ERROR)
		}
	case '\n':
		s.nextch()
		if nlsemi {
			s.tok.yys = SEMICOLON
		}
		s.nlsemi = false
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.number(false)
	case '"':
		s.stdString()
	case '`':
		s.rawString()
	case '(':
		s.nextch()
		s.tok.yys = LPAREN
	case ')':
		s.nextch()
		s.nlsemi = true
		s.tok.yys = RPAREN
	case '[':
		s.nextch()
		s.tok.yys = LBRACK
	case ']':
		s.nextch()
		s.nlsemi = true
		s.tok.yys = RBRACK
	case '{':
		s.nextch()
		s.tok.yys = LBRACE
	case '}':
		s.nextch()
		s.nlsemi = true
		s.tok.yys = RBRACE
	case ',':
		s.nextch()
		s.tok.yys = COMMA
	case ';':
		s.nextch()
		s.tok.yys = SEMICOLON
	case ':':
		s.nextch()
		if s.ch == '=' {
			s.nextch()
			s.tok.yys = CASSIGN
			s.tok.tok = ":="
			break
		} else if isLetter(s.ch) {
			s.colbefore = true
			s.ident() /* might be :ds */
			if !s.colbefore {
				s.rewind()
				s.nextch()
			} else {
				s.colbefore = false
				return
			}
		}
		s.tok.yys = COLON
	case '.':
		s.nextch()
		if isDecimal(s.ch) {
			s.number(true)
			break
		}
		s.tok.yys = PERIOD
	case '=':
		s.nextch()
		if s.ch == '=' {
			s.nextch()
			s.tok.yys = EQ_OP
			break
		}
		s.tok.yys = ASSIGN
		s.tok.tok = "="
	case '!':
		s.nextch()
		if s.ch == '=' {
			s.nextch()
			s.tok.yys = NE_OP
			break
		}
		s.tok.yys = NEG_OP
		s.tok.tok = "!"
	case '#':
		s.nextch()
		s.tok.yys = INFER
	case '+':
		s.nextch()
		s.tok.yys = ADD_OP
		s.tok.tok = "+"
		if s.ch == '=' {
			s.nextch()
			s.tok.yys = ADD_ASSIGN
			s.tok.tok = "+="
		} else if s.ch == '+' {
			s.nextch()
			s.nlsemi = true
			s.tok.yys = INC_OP
			s.tok.tok = ""
		}
	case '-':
		s.nextch()
		s.tok.yys = SUB_OP
		s.tok.tok = "-"
		if s.ch == '=' {
			s.nextch()
			s.tok.yys = SUB_ASSIGN
			s.tok.tok = "-="
		} else if s.ch == '-' {
			s.nextch()
			s.nlsemi = true
			s.tok.yys = DEC_OP
			s.tok.tok = ""
		}
	case '^':
		s.nextch()
		s.tok.yys = BITXOR_OP
		if s.ch == '=' {
			s.nextch()
			s.tok.tok = "^="
			s.tok.yys = XOR_ASSIGN
		}
	case '&':
		s.nextch()
		if s.ch == '^' {
			s.nextch()
			s.tok.yys = BITCLEAR_OP
			break
		} else if s.ch == '=' {
			s.nextch()
			s.tok.yys = AND_ASSIGN
			s.tok.tok = "&="
			break
		} else if s.ch == '&' {
			s.nextch()
			s.tok.yys = AND_OP
			break
		}
		s.tok.yys = REF_OP
		s.tok.tok = "&"
	case '*':
		s.nextch()
		if s.ch == '=' {
			s.nextch()
			s.tok.yys = MUL_ASSIGN
			s.tok.tok = "*="
			break
		}
		s.tok.yys = MUL_OP
		s.tok.tok = "*"
	case '/':
		s.nextch()
		if s.ch == '=' {
			s.nextch()
			s.tok.yys = DIV_ASSIGN
			s.tok.tok = "/="
			break
		} else if s.ch == '/' {
			s.nextch()
			s.lineComment()
			goto redonext
		} else if s.ch == '*' {
			s.nextch()
			s.fullComment()
			goto redonext
		}
		s.tok.yys = DIV_OP
		s.tok.tok = "/"
	case '%':
		s.nextch()
		if s.ch == '=' {
			s.nextch()
			s.tok.yys = MOD_ASSIGN
			s.tok.tok = "%="
			break
		}
		s.tok.yys = MOD_OP
		s.tok.tok = "%"
	case '|':
		s.nextch()
		if s.ch == '|' {
			s.nextch()
			s.tok.yys = OR_OP
			break
		} else if s.ch == '=' {
			s.nextch()
			s.tok.yys = OR_ASSIGN
			s.tok.tok = "|="
			break
		}
		s.tok.yys = BITOR_OP
	case '>':
		s.nextch()
		if s.ch == '>' {
			s.nextch()
			if s.ch == '=' {
				s.nextch()
				s.tok.yys = RIGHT_ASSIGN
				s.tok.tok = ">>="
				break
			}
			s.tok.yys = RIGHT_OP
			break
		} else if s.ch == '=' {
			s.nextch()
			s.tok.yys = GTEQ_OP
			s.tok.tok = ">="
			break
		}
		s.tok.yys = GT_OP
		s.tok.tok = ">"
	case '<':
		s.nextch()
		if s.ch == '<' {
			s.nextch()
			if s.ch == '=' {
				s.nextch()
				s.tok.yys = LEFT_ASSIGN
				s.tok.tok = "<<="
				break
			}
			s.tok.yys = LEFT_OP
			break
		} else if s.ch == '=' {
			s.nextch()
			s.tok.yys = LTEQ_OP
			s.tok.tok = "<="
			break
		}
		s.tok.yys = LT_OP
		s.tok.tok = "<"
	default:
		s.errorf(fmt.Sprintf("invalid character %#U", s.ch))
		s.nextch()
		goto redonext
	}
	if s.ch == -1 {

		s.eof = true
	}
	//fmt.Printf("%s\n", tokenName(s.tok.yys))
}

func (s *Lexer) ident() {
	// accelerate common case (7bit ASCII)
	for isLetter(s.ch) || isDecimal(s.ch) {
		s.nextch()
	}

	// general case
	if s.ch >= utf8.RuneSelf {
		for s.atIdentChar(false) {
			s.nextch()
		}
	}

	// possibly a keyword
	lit := s.segment()
	s.tok = &yySymType{}
	if len(lit) >= 2 {
		if tok := KeywordMap[string(lit)]; tok != 0 {
			switch tok {
			case IDENTIFIER,
				BOOL, STR,
				I8, I16, I32, I64,
				UI8, UI16, UI32, UI64,
				F32, F64, AFF,
				RETURN, BREAK, CONTINUE, BOOLEAN_LITERAL:
				s.nlsemi = true
				s.tok.tok = getTypeName(tok)
			}
			s.tok.yys = tok
			if tok == BOOLEAN_LITERAL {
				if string(lit) == "true" {
					s.tok.bool = true
				} else {
					s.tok.bool = false
				}
			}
			return
		} else {
			if s.colbefore {
				/* then it was actually COLON IDENTIFIER. */
				s.colbefore = false
				return
			}
		}
	}

	s.nlsemi = true
	s.tok.yys = IDENTIFIER
	s.tok.tok = string(lit)
}

func getTypeName(tok int) string {
	switch tok {
	case BOOL:
		return "bool"
	case AFF:
		return "aff"
	case F32:
		return "f32"
	case F64:
		return "f64"
	case I8:
		return "i8"
	case I16:
		return "i16"
	case I32:
		return "i32"
	case I64:
		return "i64"
	case UI8:
		return "ui8"
	case UI16:
		return "ui16"
	case UI32:
		return "ui32"
	case UI64:
		return "ui64"
	default:
		return ""
	}
}

func (s *Lexer) atIdentChar(first bool) bool {
	switch {
	case s.ch == -1 || s.ch == 0:
		return false
	case unicode.IsLetter(s.ch) || s.ch == '_':
	case unicode.IsDigit(s.ch) && !first:
		s.errorf(fmt.Sprintf("identifier cannot begin with digit %#U", s.ch))
		return false
	case s.ch >= utf8.RuneSelf:
		s.errorf(fmt.Sprintf("invalid character %#U in identifier", s.ch))
		return false
	default:
		return false
	}
	return true
}

func lower(ch rune) rune     { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func isLetter(ch rune) bool  { return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' }
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }
func isHex(ch rune) bool     { return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f' }

func (s *Lexer) number(seenPoint bool) {
	ok := true
	base := 10        // number base
	prefix := rune(0) // one of 0 (decimal), '0' (0-octal), 'x', 'o', or 'b'
	invalid := -1     // index of invalid digit in literal, or < 0
	kind := yyDefault
	// integer part
	if !seenPoint {
		if s.ch == '0' {
			s.nextch()
			switch lower(s.ch) {
			case 'x':
				s.nextch()
				base, prefix = 16, 'x'
			default:
				base, prefix = 8, '0'
			}
		}
		if base == 16 {
			for isHex(s.ch) {
				s.nextch()
			}
		} else {
			for isDecimal(s.ch) {
				s.nextch()
			}
		}
		if s.ch == '.' {
			kind = FLOAT_LITERAL
			if prefix == 'o' || prefix == 'b' {
				s.errorf(fmt.Sprintf("invalid radix point in %s literal", baseName(base)))
				ok = false
			}
			s.nextch()
			seenPoint = true
		}
	}

	// fractional part
	if seenPoint {
		kind = FLOAT_LITERAL
		if base == 16 {
			for isHex(s.ch) {
				s.nextch()
			}
		} else {
			for isDecimal(s.ch) {
				s.nextch()
			}
		}
	}

	// exponent
	if e := lower(s.ch); e == 'e' || e == 'p' {
		if ok {
			switch {
			case e == 'e' && prefix != 0 && prefix != '0':
				s.errorf(fmt.Sprintf("%q exponent requires decimal mantissa", s.ch))
				ok = false
			case e == 'p' && prefix != 'x':
				s.errorf(fmt.Sprintf("%q exponent requires hexadecimal mantissa", s.ch))
				ok = false
			}
		}
		s.nextch()
		kind = FLOAT_LITERAL
		if s.ch == '+' || s.ch == '-' {
			s.nextch()
		}
		if base == 16 {
			for isHex(s.ch) {
				s.nextch()
			}
		} else {
			for isDecimal(s.ch) {
				s.nextch()
			}
		}
	} else if prefix == 'x' && kind == FLOAT_LITERAL && ok {
		s.errorf("hexadecimal mantissa requires a 'p' exponent")
		ok = false
	}

	yylex := string(s.segment())

	// suffixes
	switch s.ch {
	case 'B':
		kind = BYTE_LITERAL
	case 'H':
		kind = SHORT_LITERAL
	case 'L':
		kind = LONG_LITERAL
	case 'F':
		kind = FLOAT_LITERAL
	case 'D':
		kind = DOUBLE_LITERAL
	case 'U':
		s.nextch()
		switch s.ch {
		case 'B':
			kind = UNSIGNED_BYTE_LITERAL
		case 'H':
			kind = UNSIGNED_SHORT_LITERAL
		case 'L':
			kind = UNSIGNED_LONG_LITERAL
		default:
			kind = UNSIGNED_INT_LITERAL
		}
	default:
		if kind != FLOAT_LITERAL {
			kind = INT_LITERAL
		}
	}

	if kind != INT_LITERAL && kind != UNSIGNED_INT_LITERAL && !(kind == FLOAT_LITERAL && s.ch != 'F') {
		s.nextch()
	}

	if kind != FLOAT_LITERAL && kind != DOUBLE_LITERAL && invalid >= 0 && ok {
		s.errorf(fmt.Sprintf("invalid digit in %s literal", baseName(base)))
		ok = false
	}

	s.tok = &yySymType{}
	switch kind {
	case INT_LITERAL:
		result, err := strconv.ParseInt(yylex, base, 32)
		if err != nil {
			s.errorf("invalid int literal: " + yylex)
		}
		s.tok.i32 = int32(result)
	case FLOAT_LITERAL:
		result, err := strconv.ParseFloat(yylex, 32)
		if err != nil {
			s.errorf("invalid float literal: " + yylex)
		}
		s.tok.f32 = float32(result)
	case DOUBLE_LITERAL:
		result, err := strconv.ParseFloat(yylex, 64)
		if err != nil {
			s.errorf("invalid double literal: " + yylex)
		}
		s.tok.f64 = result
	case SHORT_LITERAL:
		result, err := strconv.ParseInt(yylex, base, 16)
		if err != nil {
			s.errorf("invalid short literal: " + yylex)
		}
		s.tok.i16 = int16(result)
	case LONG_LITERAL:
		result, err := strconv.ParseInt(yylex, base, 64)
		if err != nil {
			s.errorf("invalid long literal: " + yylex)
		}
		s.tok.i64 = result
	case BYTE_LITERAL:
		result, err := strconv.ParseUint(yylex, base, 8)
		if err != nil {
			s.errorf("invalid byte literal: " + yylex)
		}
		s.tok.i8 = int8(result)
	case UNSIGNED_INT_LITERAL:
		result, err := strconv.ParseUint(yylex, base, 32)
		if err != nil {
			s.errorf("invalid unsigned int literal: " + yylex)
		}
		s.tok.ui32 = uint32(result)
	case UNSIGNED_BYTE_LITERAL:
		result, err := strconv.ParseUint(yylex, base, 8)
		if err != nil {
			s.errorf("invalid unsigned byte literal: " + yylex)
		}
		s.tok.ui8 = uint8(result)
	case UNSIGNED_LONG_LITERAL:
		result, err := strconv.ParseUint(yylex, base, 64)
		if err != nil {
			s.errorf("invalid unsigned long literal: " + yylex)
		}
		s.tok.ui64 = result
	case UNSIGNED_SHORT_LITERAL:
		result, err := strconv.ParseUint(yylex, base, 16)
		if err != nil {
			s.errorf("invalid unsigned short literal: " + yylex)
		}
		s.tok.ui16 = uint16(result)
	}
	s.tok.yys = kind
	s.nlsemi = true

	//s.bad = !ok // correct s.bad
}

func baseName(base int) string {
	switch base {
	case 2:
		return "binary"
	case 8:
		return "octal"
	case 10:
		return "decimal"
	case 16:
		return "hexadecimal"
	}
	panic("invalid base")
}

func (s *Lexer) fullComment() {
	for s.ch >= 0 {
		for s.ch == '*' {
			s.nextch()
			if s.ch == '/' {
				s.nextch()
				return
			}
		}
		s.nextch()
	}
	s.errorf("comment not terminated")
}

func (s *Lexer) lineComment() {
	// don't consume '\n' - needed for nlsemi logic
	for s.ch >= 0 && s.ch != '\n' {
		s.nextch()
	}
}

func (s *Lexer) rawString() {
	s.stop()
	s.nextch()
	s.start()
	for {
		if s.ch == '`' {
			s.tok.tok = string(s.segment())
			s.nextch()
			break
		}
		if s.ch < 0 {
			s.errorf("string not terminated")
			break
		}
		s.nextch()
	}
	// We leave CRs in the string since they are part of the
	// literal (even though they are not part of the literal
	// value).

	s.tok.yys = STRING_LITERAL
	s.nlsemi = true
}

func (s *Lexer) stdString() {
	s.nextch()

	for {
		if s.ch == '"' {
			s.nextch()
			break
		}
		if s.ch == '\\' {
			s.nextch()
			if !s.escape('"') {

			}
			continue
		}
		if s.ch == '\n' {
			//TODO: ISSUE #152 (old #133a) THIS FIXES IT.
			//s.errorf("newline in string")
			//break
			s.crash = true
			s.nextch()
		}
		if s.ch < 0 {
			s.errorf("string not terminated")
			break
		}
		s.nextch()
	}

	var err error
	s.tok.yys = STRING_LITERAL
	s.tok.tok, err = strconv.Unquote(string(s.segment()))
	if err != nil {
		s.errorf("invalid string literal: " + err.Error())
	}
	s.nlsemi = true
}

func (s *Lexer) escape(quote rune) bool {
	var n int
	var base, max uint32

	switch s.ch {
	case quote, 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\':
		s.nextch()
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		s.nextch()
		n, base, max = 2, 16, 255
	case 'u':
		s.nextch()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		s.nextch()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		if s.ch < 0 {
			return true // complain in caller about EOF
		}
		s.errorf("unknown escape")
		return false
	}

	var x uint32
	for i := n; i > 0; i-- {
		if s.ch < 0 {
			return true // complain in caller about EOF
		}
		d := base
		if isDecimal(s.ch) {
			d = uint32(s.ch) - '0'
		} else if 'a' <= lower(s.ch) && lower(s.ch) <= 'f' {
			d = uint32(lower(s.ch)) - 'a' + 10
		}
		if d >= base {
			s.errorf(fmt.Sprintf("invalid character %q in %s escape", s.ch, baseName(int(base))))
			return false
		}
		// d < base
		x = x*base + d
		s.nextch()
	}

	if x > max && base == 8 {
		s.errorf(fmt.Sprintf("octal escape value %d > 255", x))
		return false
	}

	if x > max || 0xD800 <= x && x < 0xE000 /* surrogate range */ {
		s.errorf(fmt.Sprintf("escape is invalid Unicode code point %#U", x))
		return false
	}

	return true
}
