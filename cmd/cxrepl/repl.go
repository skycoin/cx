package repl

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
	"github.com/skycoin/cx/cxgo/cxparser"
)

const VERSION = "0.8.0"

var ReplTargetFn string = ""
var ReplTargetStrct string = ""
var ReplTargetMod string = ""

func unsafeEval(code string) (out string) {
	var lexer *cxgo.Lexer
	defer func() {
		if r := recover(); r != nil {
			out = fmt.Sprintf("%v", r)
			// lexer.Stop()
			return
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

	lexer = cxgo.NewLexer(bytes.NewBufferString(code))
	cxgo.Parse(lexer)
	//yyParse(lexer)
	err := cxparser.AddInitFunction(actions.PRGRM)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
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
			if ReplTargetFn != "" {
				inp = fmt.Sprintf(":func %s {\n%s\n}\n", ReplTargetFn, inp)
			}
			if ReplTargetMod != "" {
				inp = fmt.Sprintf("%s", inp)
			}
			if ReplTargetStrct != "" {
				inp = fmt.Sprintf(":struct %s {\n%s\n}\n", ReplTargetStrct, inp)
			}

			b := bytes.NewBufferString(inp)

			cxgo.Parse(cxgo.NewLexer(b))
			//yyParse(NewLexer(b))
		} else {
			if ReplTargetFn != "" {
				ReplTargetFn = ""
				continue
			}

			if ReplTargetStrct != "" {
				ReplTargetStrct = ""
				continue
			}

			if ReplTargetMod != "" {
				ReplTargetMod = ""
				continue
			}

			fmt.Printf("\nBye!\n")
			break
		}
	}
}

func printPrompt() {
	if ReplTargetMod != "" {
		fmt.Println(fmt.Sprintf(":package %s ...", ReplTargetMod))
		fmt.Printf("* ")
	} else if ReplTargetFn != "" {
		fmt.Println(fmt.Sprintf(":func %s {...", ReplTargetFn))
		fmt.Printf("\t* ")
	} else if ReplTargetStrct != "" {
		fmt.Println(fmt.Sprintf(":struct %s {...", ReplTargetStrct))
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
