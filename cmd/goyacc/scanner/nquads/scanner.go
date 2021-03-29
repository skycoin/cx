// Copyright (c) 2014 The scanner Authors. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in the
// LICENSE file.
//
// CAUTION: If this file is a Go source file (*.go), it was generated
// automatically by '$ golex' from a *.l file - DO NOT EDIT in that case!

package scanner // import "github.com/skycoin/cx/cmd/goyacc/scanner/nquads"

import (
	"strconv"
)

// Scan scans the next token and returns the token and its value if applicable.
// The source end is indicated by EOF.
//
// If the returned token is ILLEGAL, the literal string is the offending
// character or number/string/char literal.
func (s *Scanner) Scan() (tok Token, st string) {
	//defer func() { dbg("%s:%d:%d %v %q :%d:%d s.i %d: %#x", s.Fname, s.Line, s.Col, tok, st, s.NLine, s.NCol, s.i, s.c) }()
	c := s.c

yystate0:

	s.val = s.val[:0]
	s.i0, s.Line, s.Col = s.i, s.NLine, s.NCol
	i := s.i
	if i > 0 {
		i--
	}
	c0, n0 := decodeRune(s.src[i:])
	if c < 0 {
		s.i0++
		return EOF, ""
	}

	goto yystart1

	goto yystate1 // silence unused label error
yystate1:
	c = s.next()
yystart1:
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate4
	case c == '#':
		goto yystate15
	case c == '.':
		goto yystate16
	case c == '<':
		goto yystate17
	case c == '@':
		goto yystate28
	case c == '\n' || c == '\r':
		goto yystate3
	case c == '\t' || c == ' ':
		goto yystate2
	case c == '^':
		goto yystate32
	case c == '_':
		goto yystate34
	}

yystate2:
	c = s.next()
	switch {
	default:
		goto yyrule1
	case c == '\t' || c == ' ':
		goto yystate2
	}

yystate3:
	c = s.next()
	switch {
	default:
		goto yyrule6
	case c == '\n' || c == '\r':
		goto yystate3
	}

yystate4:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate5
	case c == '\\':
		goto yystate6
	case c >= '\x00' && c <= '\t' || c == '\v' || c == '\f' || c >= '\x0e' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate4
	}

yystate5:
	// c = s.next()
	goto yyrule9

yystate6:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '"' || c == '\'' || c == '\\' || c == 'b' || c == 'f' || c == 'n' || c == 'r' || c == 't':
		goto yystate4
	case c == 'U':
		goto yystate7
	case c == 'u':
		goto yystate11
	}

yystate7:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate8
	}

yystate8:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate9
	}

yystate9:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate10
	}

yystate10:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate11
	}

yystate11:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate12
	}

yystate12:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate13
	}

yystate13:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate14
	}

yystate14:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate4
	}

yystate15:
	c = s.next()
	switch {
	default:
		goto yyrule2
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate15
	}

yystate16:
	// c = s.next()
	goto yyrule3

yystate17:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '!' || c >= '#' && c <= ';' || c == '=' || c >= '?' && c <= '[' || c == ']' || c == '_' || c >= 'a' && c <= 'z' || c >= '~' && c <= 'ÿ':
		goto yystate17
	case c == '>':
		goto yystate18
	case c == '\\':
		goto yystate19
	}

yystate18:
	// c = s.next()
	goto yyrule7

yystate19:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'U':
		goto yystate20
	case c == 'u':
		goto yystate24
	}

yystate20:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate21
	}

yystate21:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate22
	}

yystate22:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate23
	}

yystate23:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate24
	}

yystate24:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate25
	}

yystate25:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate26
	}

yystate26:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate27
	}

yystate27:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate17
	}

yystate28:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate29
	}

yystate29:
	c = s.next()
	switch {
	default:
		goto yyrule8
	case c == '-':
		goto yystate30
	case c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate29
	}

yystate30:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate31:
	c = s.next()
	switch {
	default:
		goto yyrule8
	case c == '-':
		goto yystate30
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate31
	}

yystate32:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '^':
		goto yystate33
	}

yystate33:
	// c = s.next()
	goto yyrule4

yystate34:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == ':':
		goto yystate35
	}

yystate35:
	// c = s.next()
	goto yyrule5

yyrule1: // [ \t]+

	goto yystate0
yyrule2: // #.*

	goto yystate0
yyrule3: // \.
	{
		return DOT, "."
	}
yyrule4: // "^^"
	{
		return DACCENT, "^^"
	}
yyrule5: // {blank_node_label}
	{

		if s.c < 0 {
			return ILLEGAL, string(c0)
		}
		var v []rune
		s.i--
		switch r, n := decodeRune(s.src[s.i:]); {
		case checkPnCharsU(r), r >= '0' && r <= '9':
			s.i += n
			s.NCol += n
			v = append(v, r)
		default:
			s.next()
			return ILLEGAL, string(c0)
		}
	loop:
		for {
			switch r, n := decodeRune(s.src[s.i:]); {
			case checkPnChars(r), r == '.':
				s.i += n
				s.NCol += n
				v = append(v, r)
			default:
				if v[len(v)-1] != '.' {
					s.next()
					break loop
				}
				for v[len(v)-1] == '.' {
					v = v[:len(v)-1]
					s.i--
					s.NCol--
				}
				s.next()
				break loop
			}
		}

		if s.NCol > 0 {
			s.NCol--
		}
		return LABEL, string(v)
	}
yyrule6: // {eol}
	{
		return EOL, ""
	}
yyrule7: // {iriref}
	{

		val, err := strconv.Unquote(`"` + string(s.val) + `"`)
		if err != nil {
			s.err(err.Error())
		}
		return IRIREF, val[1 : len(val)-1]
	}
yyrule8: // {langtag}
	{

		return LANGTAG, string(s.val[1:])
	}
yyrule9: // {string_literal_quote}
	{

		// \' needs special preprocessing.
		var c, prev byte
		v := s.val
		for i := 1; i < len(v)-2; i, prev = i+1, c {
			if c = v[i]; prev != '\\' && c == '\\' && v[i+1] == '\'' {
				v = append(v[:i], v[i+1:]...)
			}
		}
		val, err := strconv.Unquote(string(v))
		if err != nil {
			s.err(err.Error())
		}
		return STRING, val
	}
	panic("unreachable")

	goto yyabort // silence unused label error

yyabort: // no lexem recognized
	//dbg("yyabort s.i %d, c0 %q(%#x), n0 %d", s.i, c0, c0, n0)
	s.i += n0 - 1
	s.next()
	return ILLEGAL, string(c0)

}
