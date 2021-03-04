package main

//put the repo command stuff here

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	cxcore "github.com/skycoin/cx/cx"
	"github.com/skycoin/cx/cxgo/actions"
	"github.com/skycoin/cx/cxgo/cxgo"
	"github.com/skycoin/cx/cxgo/cxgo0"
	"github.com/skycoin/cx/cxgo/parser"
)

func unsafeEval(code string) (out string) {
	var lexer *parser.Lexer
	defer func() {
		if r := recover(); r != nil {
			out = fmt.Sprintf("%v", r)
			lexer.Stop()
		}
	}()

	// storing strings sent to standard output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	actions.LineNo = 0

	actions.PRGRM = cxcore.MakeProgram()
	cxgo0.PRGRM0 = actions.PRGRM

	cxgo0.Parse(code)

	actions.PRGRM = cxgo0.PRGRM0

	lexer = parser.NewLexer(bytes.NewBufferString(code))
	parser.Parse(lexer)
	//yyParse(lexer)

	cxgo.AddInitFunction(actions.PRGRM)

	if err := actions.PRGRM.RunCompiled(0, nil); err != nil {
		actions.PRGRM = cxcore.MakeProgram()
		return fmt.Sprintf("%s", err)
	}

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old // restoring the real stdout
	out = <-outC

	actions.PRGRM = cxcore.MakeProgram()
	return out
}

func Eval(code string) string {
	runtime.GOMAXPROCS(2)
	ch := make(chan string, 1)

	var result string

	go func() {
		result = unsafeEval(code)
		ch <- result
	}()

	timer := time.NewTimer(20 * time.Second)
	defer timer.Stop()

	select {
	case <-ch:
		return result
	case <-timer.C:
		actions.PRGRM = cxcore.MakeProgram()
		return "Timed out."
	}
}

func Repl() {
	fmt.Println("CX", VERSION)
	fmt.Println("More information about CX is available at http://cx.skycoin.com/ and https://github.com/skycoin/cx/")

	cxcore.InREPL = true

	// fi := bufio.NewReader(os.NewFile(0, "stdin"))
	fi := bufio.NewReader(os.Stdin)
	// scanner := bufio.NewScanner(os.Stdin)

	for {
		var inp string
		var ok bool

		printPrompt()

		if inp, ok = readline(fi); ok {
			if actions.ReplTargetFn != "" {
				inp = fmt.Sprintf(":func %s {\n%s\n}\n", actions.ReplTargetFn, inp)
			}
			if actions.ReplTargetMod != "" {
				inp = fmt.Sprintf("%s", inp)
			}
			if actions.ReplTargetStrct != "" {
				inp = fmt.Sprintf(":struct %s {\n%s\n}\n", actions.ReplTargetStrct, inp)
			}

			b := bytes.NewBufferString(inp)

			parser.Parse(parser.NewLexer(b))
			//yyParse(NewLexer(b))
		} else {
			if actions.ReplTargetFn != "" {
				actions.ReplTargetFn = ""
				continue
			}

			if actions.ReplTargetStrct != "" {
				actions.ReplTargetStrct = ""
				continue
			}

			if actions.ReplTargetMod != "" {
				actions.ReplTargetMod = ""
				continue
			}

			fmt.Printf("\nBye!\n")
			break
		}
	}
}

func printPrompt() {
	if actions.ReplTargetMod != "" {
		fmt.Println(fmt.Sprintf(":package %s ...", actions.ReplTargetMod))
		fmt.Printf("* ")
	} else if actions.ReplTargetFn != "" {
		fmt.Println(fmt.Sprintf(":func %s {...", actions.ReplTargetFn))
		fmt.Printf("\t* ")
	} else if actions.ReplTargetStrct != "" {
		fmt.Println(fmt.Sprintf(":struct %s {...", actions.ReplTargetStrct))
		fmt.Printf("\t* ")
	} else {
		fmt.Printf("* ")
	}
}

// ----------------------------------------------------------------
//                     Utility functions

func readline(fi *bufio.Reader) (string, bool) {
	s, err := fi.ReadString('\n')

	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)

	for _, ch := range s {
		if ch == rune(4) {
			err = io.EOF
			break
		}
	}

	if err != nil {
		return "", false
	}

	return s, true
}
