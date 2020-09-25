package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"

	. "github.com/SkycoinProject/cx/cxgo/parser"
	. "github.com/SkycoinProject/cx/tests/test-lexer/neximport"
)

var filename string

func main0() {
	//lex := &Lexer{}

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("error: file IO error: " + err.Error())
		os.Exit(1)
	}
	source := string(src)
	b := bytes.NewBufferString(source)
	fmt.Printf("Beginning timer for new Lexer...\n")
	tm := time.Now()
	lex := NewLexer(b)
	ntok := 0
	for lex.Next() > 0 {
		ntok++
	}
	tel := time.Now().Sub(tm)
	times := tel.Seconds()
	rat := times
	fmt.Printf("Time taken (SSSS.MiS.McS): %04d.", int(math.Floor(times)))
	times -= math.Floor(times)
	times *= 1000
	fmt.Printf("%03d.", int(math.Floor(times)))
	times -= math.Floor(times)
	times *= 1000
	fmt.Printf("%03d\n", int(math.Floor(times)))
	fmt.Printf("Profile complete for new Lexer. NTOKS: %d\n\n", ntok)

	src, err = ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("error: file IO error: " + err.Error())
		os.Exit(1)
	}
	source = string(src)
	b = bytes.NewBufferString(source)
	fmt.Printf("Beginning timer for old Lexer...\n")
	tm = time.Now()
	lex2 := NewLexer_(b)
	ntok = 0
	for lex2.Next() > 0 {
		ntok++
	}
	tel = time.Now().Sub(tm)
	times2 := tel.Seconds()
	rat2 := times2
	fmt.Printf("Time taken (SSSS.MiS.McS): %04d.", int(math.Floor(times2)))
	times2 -= math.Floor(times2)
	times2 *= 1000
	fmt.Printf("%03d.", int(math.Floor(times2)))
	times2 -= math.Floor(times2)
	times2 *= 1000
	fmt.Printf("%03d\n", int(math.Floor(times2)))
	fmt.Printf("Profile complete for old Lexer. NTOKS: %d\n\n", ntok)
	times2 = rat2 - rat
	fmt.Printf("Total time more than new lexer: %04d.", int(math.Floor(times2)))
	times2 -= math.Floor(times2)
	times2 *= 1000
	fmt.Printf("%03d.", int(math.Floor(times2)))
	times2 -= math.Floor(times2)
	times2 *= 1000
	fmt.Printf("%03d\n", int(math.Floor(times2)))
	fmt.Printf("Ratio of old Lexer to new Lexer: %f\n", rat2/rat)
	fmt.Printf("Total lexer profile/comparison complete.\n")
}

func main1() {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("error: file IO error: " + err.Error())
		os.Exit(1)
	}
	source := string(src)
	b := bytes.NewBufferString(source)
	lex1 := NewLexer(b)

	src2, err2 := ioutil.ReadFile(filename)
	if err2 != nil {
		fmt.Printf("error: file IO error: " + err2.Error())
		os.Exit(1)
	}
	source2 := string(src2)
	b2 := bytes.NewBufferString(source2)
	lex2 := NewLexer_(b2)

	t1 := 65336
	t2 := 65336
	i := 0
	for t1 == t2 && t1 > 0 && t2 > 0 {
		t1 = lex1.Next()
		t2 = lex2.Next()
		i++
	}
	if t1 <= 0 || t2 <= 0 {
		fmt.Printf("Succeeded.\n")
	} else {
		fmt.Printf("Failed on step %d.\n", i)
		os.Exit(-1)
	}
}

func main() {
	filename = os.Args[1]
	fmt.Printf("\n\n\ninitializing...\n")
	/* sleep because test-toks.txt might've not been fully written yet */
	time.Sleep(time.Second * 3)
	fmt.Printf("running similarity regression test.\n")
	main1()
	fmt.Printf("running speed comparison/profile.\n")
	main0()
	fmt.Printf("All tests complete.\n\n\n")
}
