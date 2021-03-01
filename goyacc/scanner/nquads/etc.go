// Copyright (c) 2014 The scanner Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package scanner implements a scanner for N-Quads[0] source text. It takes a
// []byte as source which can then be tokenized through repeated calls to the
// Scan method.
//
// Links
//
// Referenced from elsewhere:
//
//  [0]: http://www.w3.org/TR/n-quads/
//  [1]: http://www.w3.org/TR/n-quads/#grammar-production-BLANK_NODE_LABEL
//  [2]: http://www.w3.org/TR/n-quads/#grammar-production-EOL
//  [3]: http://www.w3.org/TR/n-quads/#grammar-production-IRIREF
//  [4]: http://www.w3.org/TR/n-quads/#grammar-production-LANGTAG
//  [5]: http://www.w3.org/TR/n-quads/#grammar-production-STRING_LITERAL_QUOTE
package scanner // import "github.com/skycoin/cx/goyacc/scanner/nquads"

import (
	"errors"
	"fmt"
	"math"
	"unicode"
)

// Token is the type of the token identifier returned by Scan().
type Token int

// Values of type Token.
const (
	_ = Token(0xE000) + iota // Unicode user space

	ILLEGAL // Returned when no token was recognized.
	EOF     // Returned after all source was consumed.

	DOT     // .
	DACCENT // ^^
	LABEL   // [1]
	EOL     // [2]
	IRIREF  // [3]
	LANGTAG // [4]
	STRING  // [5]
)

var ts = map[Token]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	LABEL:   "LABEL",
	EOL:     "EOL",
	IRIREF:  "IRIREF",
	LANGTAG: "LANGTAG",
	STRING:  "STRING",
	DOT:     "DOT",
	DACCENT: "DACCENT",
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
	Col    int     // Starting column of the last scanned token.
	Errors []error // List of accumulated errors.
	Fname  string  // File name (reported) of the scanned source.
	Line   int     // Starting line of the last scanned token.
	NCol   int     // Starting column (reported) for the next scanned token.
	NLine  int     // Starting line (reported) for the next scanned token.
	c      int
	i      int
	i0     int
	sc     int
	src    []byte
	val    []byte
}

// New returns a newly created Scanner with fname as the name of the source.
func New(fname string, src []byte) (s *Scanner) {
	if len(src) > 2 && src[0] == 0xEF && src[1] == 0xBB && src[2] == 0xBF {
		src = src[3:]
	}
	s = &Scanner{
		src:   src,
		NLine: 1,
		NCol:  0,
		Fname: fname,
	}
	s.next()
	return
}

func (s *Scanner) back() {
	s.NCol--
	s.i--
	s.c = int(s.src[s.i])
}

func (s *Scanner) next() int {
	if s.i <= len(s.src) && s.c >= 0 {
		s.val = append(s.val, byte(s.c))
	}
	s.c = -1
	if s.i < len(s.src) {
		s.c = int(s.src[s.i])
		s.i++
		switch s.c {
		case '\n':
			s.NLine++
			s.NCol = 0
			if s.i == len(s.src) {
				s.NCol = 1
			}
		default:
			s.NCol++
		}
	}
	return s.c
}

// Pos returns the starting offset of the last scanned token.
func (s *Scanner) Pos() int {
	return s.i0
}

func (s *Scanner) err(format string, arg ...interface{}) {
	err := fmt.Errorf(fmt.Sprintf("%s:%d:%d ", s.Fname, s.Line, s.Col)+format, arg...)
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

const (
	runeError = math.MaxInt32

	t1 = 0x00 // 0000 0000
	tx = 0x80 // 1000 0000
	t2 = 0xC0 // 1100 0000
	t3 = 0xE0 // 1110 0000
	t4 = 0xF0 // 1111 0000
	t5 = 0xF8 // 1111 1000

	maskx = 0x3F // 0011 1111
	mask2 = 0x1F // 0001 1111
	mask3 = 0x0F // 0000 1111
	mask4 = 0x07 // 0000 0111

	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1
)

// A customized version of strconv.DecodeRune is needed, which does not turn
// illegal encoding into \uFFFD as that's unfortunately accepted in some parts
// of the N-Quads grammar (PN_CHARS_BASE).
func decodeRune(s []byte) (r rune, size int) {
	n := len(s)
	if n < 1 {
		return 0, 0
	}
	c0 := s[0]

	// 1-byte, 7-bit sequence?
	if c0 < tx {
		return rune(c0), 1
	}

	// unexpected continuation byte?
	if c0 < t2 {
		return runeError, 1
	}

	// need first continuation byte
	if n < 2 {
		return runeError, 1
	}
	c1 := s[1]
	if c1 < tx || t2 <= c1 {
		return runeError, 1
	}

	// 2-byte, 11-bit sequence?
	if c0 < t3 {
		r = rune(c0&mask2)<<6 | rune(c1&maskx)
		if r <= rune1Max {
			return runeError, 1
		}
		return r, 2
	}

	// need second continuation byte
	if n < 3 {
		return runeError, 1
	}
	c2 := s[2]
	if c2 < tx || t2 <= c2 {
		return runeError, 1
	}

	// 3-byte, 16-bit sequence?
	if c0 < t4 {
		r = rune(c0&mask3)<<12 | rune(c1&maskx)<<6 | rune(c2&maskx)
		if r <= rune2Max {
			return runeError, 1
		}
		return r, 3
	}

	// need third continuation byte
	if n < 4 {
		return runeError, 1
	}
	c3 := s[3]
	if c3 < tx || t2 <= c3 {
		return runeError, 1
	}

	// 4-byte, 21-bit sequence?
	if c0 < t5 {
		r = rune(c0&mask4)<<18 | rune(c1&maskx)<<12 | rune(c2&maskx)<<6 | rune(c3&maskx)
		if r <= rune3Max || unicode.MaxRune < r {
			return runeError, 1
		}
		return r, 4
	}

	// error
	return runeError, 1
}

var pnTab = []rune{
	'A', 'Z', // 0
	'a', 'z', // 1
	0x00C0, 0x00D6, // 2
	0x00D8, 0x00F6, // 3
	0x00F8, 0x02FF, // 4
	0x0370, 0x037D, // 5
	0x037F, 0x1FFF, // 6
	0x200C, 0x200D, // 7
	0x2070, 0x218F, // 8
	0x2C00, 0x2FEF, // 9
	0x3001, 0xD7FF, // 10
	0xF900, 0xFDCF, // 11
	0xFDF0, 0xFFFD, // 12
	0x10000, 0xEFFFF, // 13, last PN_CHARS_BASE
	'_', '_', // 14
	':', ':', // 15, last PN_CHARS_U
	'-', '-', // 16
	'0', '9', // 17
	0x00B7, 0x00B7, // 18
	0x0300, 0x036F, // 19
	0x203F, 0x2040, // 20, last PN_CHARS
}

func check(r rune, tab []rune) bool {
	for i := 0; i < len(tab); i += 2 {
		if r >= tab[i] && r <= tab[i+1] {
			return true
		}
	}
	return false
}

func checkPnCharsU(r rune) bool {
	return check(r, pnTab[:2*16])
}

func checkPnChars(r rune) bool {
	return check(r, pnTab)
}
