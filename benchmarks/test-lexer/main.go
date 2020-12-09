package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/skycoin/cx/benchmarks/test-lexer/oldnex"
	"github.com/skycoin/cx/cxgo/parser"
)

func main0(filename string) {
	// Obtain file bytes.
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("error: file IO error: " + err.Error())
		os.Exit(1)
	}

	// Benchmark new lexer.
	fmt.Printf("Beginning timer for new Lexer...\n")
	start1 := time.Now()

	lex1 := parser.NewLexer(bytes.NewReader(src))
	nTok1 := 0
	for lex1.Next() > 0 {
		nTok1++
	}

	elapsed1 := time.Since(start1)
	printElapsed(elapsed1)
	fmt.Printf("Profile complete for new Lexer. nToks: %d\n\n", nTok1)

	// Benchmark old lexer.
	fmt.Printf("Beginning timer for old Lexer...\n")
	start2 := time.Now()

	lex2 := oldnex.NewLexer(bytes.NewReader(src))
	nTok2 := 0
	for lex2.Next() > 0 {
		nTok2++
	}

	elapsed2 := time.Since(start2)
	printElapsed(elapsed2)
	fmt.Printf("Profile complete for old Lexer. nToks: %d\n\n", nTok1)

	diff := elapsed2 - elapsed1
	fmt.Print("Total time more than new lexer: ")
	printElapsed(diff)
	fmt.Printf("Ratio of old Lexer to new Lexer: %f\n", float64(elapsed2)/float64(elapsed1))
	fmt.Printf("Total lexer profile/comparison complete.\n")
}

func main1(filename string) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("error: file IO error: " + err.Error())
		os.Exit(1)
	}
	source := string(src)
	b := bytes.NewBufferString(source)
	lex1 := parser.NewLexer(b)

	src2, err2 := ioutil.ReadFile(filename)
	if err2 != nil {
		fmt.Printf("error: file IO error: " + err2.Error())
		os.Exit(1)
	}
	source2 := string(src2)
	b2 := bytes.NewBufferString(source2)
	lex2 := oldnex.NewLexer(b2)

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
	filename := os.Args[1]
	fmt.Printf("\n\n\ninitializing...\n")
	/* sleep because test-toks.txt might've not been fully written yet */
	time.Sleep(time.Second * 3)
	fmt.Printf("running similarity regression test.\n")
	main1(filename)
	fmt.Printf("running speed comparison/profile.\n")
	main0(filename)
	fmt.Printf("All tests complete.\n\n\n")
}

func printElapsed(elapsed time.Duration) {
	u := elapsed.Microseconds()

	us := u%1000
	u /= 1000

	ms := u%1000
	u /= 1000

	s := u%1000
	u /= 1000

	fmt.Printf("Elapsed (SSSS.MiS.McS): %04d.%03d.%03d\n", s, ms, us)
}