// Copyright (c) 2014 The scanner Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner // import "github.com/skycoin/cx/goyacc/scanner/nquads"

import (
	"encoding/hex"
	"fmt"
	"path"
	"runtime"
	"strings"
	"testing"
	"unicode"
)

func caller(s string, va ...interface{}) {
	_, fn, fl, _ := runtime.Caller(2)
	fmt.Printf("caller: %s:%d: ", path.Base(fn), fl)
	fmt.Printf(s, va...)
	fmt.Println()
	_, fn, fl, _ = runtime.Caller(1)
	fmt.Printf("\tcallee: %s:%d: ", path.Base(fn), fl)
	fmt.Println()
}

func dbg(s string, va ...interface{}) {
	if s == "" {
		s = strings.Repeat("%v ", len(va))
	}
	_, fn, fl, _ := runtime.Caller(1)
	fmt.Printf("dbg %s:%d: ", path.Base(fn), fl)
	fmt.Printf(s, va...)
	fmt.Println()
}

func TODO(...interface{}) string {
	_, fn, fl, _ := runtime.Caller(1)
	return fmt.Sprintf("TODO: %s:%d:\n", path.Base(fn), fl)
}

func use(...interface{}) {}

func hd(b []byte) string { return hex.Dump(b) }

// ============================================================================

func TestScanner1(t *testing.T) {
	tab := []struct {
		src string
		ok  bool
		tok Token
		val string
	}{
		// 0
		{"_", true, ILLEGAL, "_"},
		{"_ ", true, ILLEGAL, "_"},
		{"_\n", true, ILLEGAL, "_"},
		{"__", true, ILLEGAL, "_"},
		{"__ ", true, ILLEGAL, "_"},

		// 5
		{"__\n", true, ILLEGAL, "_"},
		{"_:", true, ILLEGAL, "_"},
		{"_: ", true, ILLEGAL, "_"},
		{"_:\n", true, ILLEGAL, "_"},
		{"_:!", true, ILLEGAL, "_"},

		// 10
		{"_:0", true, LABEL, "0"},
		{" _:0", true, LABEL, "0"},
		{"\t_:0\t", true, LABEL, "0"},
		{"\n_:0\n", true, EOL, ""},
		{"\n\t_:0\t\n", true, EOL, ""},

		// 15
		{`"\t"`, true, STRING, "\t"},
		{`"\b"`, true, STRING, "\b"},
		{`"\n"`, true, STRING, "\n"},
		{`"\r"`, true, STRING, "\r"},
		{`"\f"`, true, STRING, "\f"},

		// 20
		{`"\""`, true, STRING, "\""},
		{`"\\"`, true, STRING, "\\"},
		{`"\'"`, true, STRING, "'"},
		{"_:subjec\\t1", true, LABEL, "subjec"},
		{"<http://one.example/subjec\\t1>", true, ILLEGAL, "<"},

		// 25
		{" \xef\xbb\xbf_:0", true, ILLEGAL, "\ufeff"}, // BOM not first
		{"\xef\xbb\xbf_:0", true, LABEL, "0"},         // BOM first
		{"_:0.x", true, LABEL, "0.x"},
		{"_:0x. ", true, LABEL, "0x"},
		{"_:0x.", true, LABEL, "0x"},

		// 30
		{"_:.x", true, ILLEGAL, "_"},
		{"_:0\u0080", true, LABEL, "0"},
		{"_:0.", true, LABEL, "0"},
		{"_:0.1", true, LABEL, "0.1"},
		{"_:0.1.", true, LABEL, "0.1"},

		// 35
		{"_:0.1..", true, LABEL, "0.1"},
		{"", true, EOF, ""},
		{"\"\x00\"", true, STRING, "\x00"},
		{`"\u0000"`, true, STRING, "\x00"},
	}

	for i, test := range tab {
		sc := New("test", []byte(test.src))
		tok, val := sc.Scan()
		errs := sc.Errors
		switch test.ok {
		case true:
			if len(errs) != 0 {
				t.Errorf("%d: errs %v", i, errs)
				break
			}

			if g, e := tok, test.tok; g != e {
				t.Errorf("%d: tok %v %v", i, g, e)
			}

			if g, e := val, test.val; g != e {
				t.Errorf("%d: val %q %q", i, g, e)
			}
		default:
			if len(errs) == 0 {
				t.Errorf("%d: errs %v %v", i, tok, val)
				break
			}
		}
	}
}

func TestScanner2(t *testing.T) {
	tab := []struct {
		src  string
		toks []Token
		vals []string
	}{
		// 0
		{"_:0.x _:1.y", []Token{LABEL, LABEL, EOF}, []string{"0.x", "1.y", ""}},
		{`_:0.x _:1.y
`, []Token{LABEL, LABEL, EOL, EOF}, []string{"0.x", "1.y", "", ""}},
		{"_:0.x .", []Token{LABEL, DOT, EOF}, []string{"0.x", ".", ""}},
		{`_:0.x .
`, []Token{LABEL, DOT, EOL, EOF}, []string{"0.x", ".", "", ""}},
		{"_:0.x .", []Token{LABEL, DOT, EOF}, []string{"0.x", ".", ""}},

		// 5
		{"_:0.", []Token{LABEL, DOT, EOF}, []string{"0", ".", ""}},
		{`_:0.x .
`, []Token{LABEL, DOT, EOL, EOF}, []string{"0.x", ".", "", ""}},
		{`_:0.x .
`, []Token{LABEL, DOT, EOL, EOF}, []string{"0.x", ".", "", ""}},
		{`_:0.
`, []Token{LABEL, DOT, EOL, EOF}, []string{"0", ".", "", ""}},
		{`_:0.
`, []Token{LABEL, DOT, EOL, EOF}, []string{"0", ".", "", ""}},

		// 10
		{`_:0.
`, []Token{LABEL, DOT, EOL, EOF}, []string{"0", ".", "", ""}},
		{`_:0.
`, []Token{LABEL, DOT, EOL, EOF}, []string{"0", ".", "", ""}},
		{"_:0\u0080.", []Token{LABEL, ILLEGAL, DOT, EOF}, []string{"0", "\u0080", ".", ""}},
		{"_:0.1", []Token{LABEL, EOF}, []string{"0.1", ""}},
		{"_:0.1.", []Token{LABEL, DOT, EOF}, []string{"0.1", ".", ""}},

		// 15
		{"_:0.1.x", []Token{LABEL, EOF}, []string{"0.1.x", ""}},
		{"_:0.1..", []Token{LABEL, DOT, DOT, EOF}, []string{"0.1", ".", ".", ""}},
		{"_:0.1..x", []Token{LABEL, EOF}, []string{"0.1..x", ""}},
		{`<http://example.org/property> _:anon.
`,
			[]Token{IRIREF, LABEL, DOT, EOL, EOF}, []string{"http://example.org/property", "anon", ".", "", ""}},
		{"<http://a.example/s> <http://a.example/p> \"\x00	&([]\" .\n",
			[]Token{IRIREF, IRIREF, STRING, DOT, EOL, EOF},
			[]string{"http://a.example/s", "http://a.example/p", "\x00\t\v\f\x0e&([]\u007f", ".", "", ""}},

		// 20
		{`<http://example/s> <http://example/p> "abc' .`,
			[]Token{IRIREF, IRIREF, ILLEGAL, EOF}, []string{"http://example/s", "http://example/p", "\"", ""}},
	}

	for i, test := range tab {
		b := []byte(test.src)
		//dbg("%d:\n%s", i, hd(b))
		sc := New("test", b)
		for j, tok := range test.toks {
			val := test.vals[j]
			gt, gv := sc.Scan()
			if g, e := gt, tok; g != e {
				t.Errorf("%d.%d: tok %v %v", i, j, g, e)
			}
			if g, e := gv, val; g != e {
				t.Errorf("%d.%d: val %q %q", i, j, g, e)
			}
		}
	}
}

// Must ignore surrogates.
func encodeRune(r rune) string {
	switch i := uint32(r); {
	case i <= rune1Max:
		return string([]byte{byte(r)})
	case i <= rune2Max:
		return string([]byte{t2 | byte(r>>6), tx | byte(r)&maskx})
	case i <= rune3Max:
		return string([]byte{t3 | byte(r>>12), tx | byte(r>>6)&maskx, tx | byte(r)&maskx})
	default:
		return string([]byte{t4 | byte(r>>18), tx | byte(r>>12)&maskx, tx | byte(r>>6)&maskx, tx | byte(r)&maskx})
	}
}

func TestLabel(t *testing.T) {
	for c := rune(0); c <= unicode.MaxRune; c++ {
		s := "_:" + encodeRune(c)
		sc := New("test", []byte(s))
		tok, val := sc.Scan()
		switch {
		case c >= '0' && c <= '9', checkPnCharsU(c):
			if g, e := tok, LABEL; g != e {
				t.Fatalf("tok: %q(%U) %v %v", c, c, g, e)
			}
			if g, e := val, s[2:]; g != e {
				t.Fatalf("val: %q(%U) %v %v", c, c, g, e)
			}
		default:
			if g, e := tok, ILLEGAL; g != e {
				t.Fatalf("should fail %q |% x|: %q(%U) %v %v", s, s, c, c, g, e)
			}
		}
	}

	for c := rune(0); c <= unicode.MaxRune; c++ {
		s := "_:0" + encodeRune(c)
		sc := New("test", []byte(s))
		tok, val := sc.Scan()
		switch {
		case checkPnChars(c):
			if g, e := tok, LABEL; g != e {
				t.Fatalf("tok: %q(%U) %v %v", c, c, g, e)
			}
			if g, e := val, s[2:]; g != e {
				t.Fatalf("val: %q(%U) %v %v", c, c, g, e)
			}
		default:
			if g, e := tok, LABEL; g != e {
				t.Fatalf("should fail %q |% x|: %q(%U) %v %v", s, s, c, c, g, e)
			}

			if g, e := val, s[2:3]; g != e {
				t.Fatalf("val: %q(%U) %v %v", c, c, g, e)
			}
		}
	}

	for c := rune(0); c <= unicode.MaxRune; c++ {
		s := "_:0." + encodeRune(c)
		sc := New("test", []byte(s))
		tok, val := sc.Scan()
		switch {
		case checkPnChars(c):
			if g, e := tok, LABEL; g != e {
				t.Fatalf("tok: %q(%U) %v %v", c, c, g, e)
			}
			if g, e := val, s[2:]; g != e {
				t.Fatalf("val: %q(%U) %v %v", c, c, g, e)
			}
		default:
			if g, e := tok, LABEL; g != e {
				t.Fatalf("should fail %q |% x|: %q(%U) %v %v", s, s, c, c, g, e)
			}

			if g, e := val, s[2:3]; g != e {
				t.Fatalf("val: %q(%U) %v %v", c, c, g, e)
			}
		}
	}
}

func ExampleScanner_Scan() {
	const src = `

<http://one.example/subject1> <http://one.example/predicate1> <http://one.example/object1> @us-EN <http://example.org/graph3> . # comments here
# or on a line by themselves
_:subject1 <http://an.example/predicate1> "object\u00411" "cafe\u0301 \'time" <http://example.org/graph1> .
_:subject2 <http://an.example/predicate2> "object\U000000422"  ^^ <http://example.com/literal> <http://example.org/graph5> .
_:0 _:01
 _:0 _:01

`
	sc := New("test", []byte(src))
	for {
		tok, val := sc.Scan()
		fmt.Printf("%s:%d:%d %v %q\n", sc.Fname, sc.Line, sc.Col, tok, val)
		if tok == EOF {
			break
		}
	}
	fmt.Printf("%v", sc.Errors)
	// Output:
	// test:2:0 EOL ""
	// test:3:1 IRIREF "http://one.example/subject1"
	// test:3:31 IRIREF "http://one.example/predicate1"
	// test:3:63 IRIREF "http://one.example/object1"
	// test:3:92 LANGTAG "us-EN"
	// test:3:99 IRIREF "http://example.org/graph3"
	// test:3:127 DOT "."
	// test:4:0 EOL ""
	// test:5:0 EOL ""
	// test:5:1 LABEL "subject1"
	// test:5:12 IRIREF "http://an.example/predicate1"
	// test:5:43 STRING "objectA1"
	// test:5:59 STRING "café 'time"
	// test:5:79 IRIREF "http://example.org/graph1"
	// test:5:107 DOT "."
	// test:6:0 EOL ""
	// test:6:1 LABEL "subject2"
	// test:6:12 IRIREF "http://an.example/predicate2"
	// test:6:43 STRING "objectB2"
	// test:6:64 DACCENT "^^"
	// test:6:67 IRIREF "http://example.com/literal"
	// test:6:96 IRIREF "http://example.org/graph5"
	// test:6:124 DOT "."
	// test:7:0 EOL ""
	// test:7:1 LABEL "0"
	// test:7:5 LABEL "01"
	// test:8:0 EOL ""
	// test:8:2 LABEL "0"
	// test:8:6 LABEL "01"
	// test:9:0 EOL ""
	// test:10:1 EOF ""
	// []
}
