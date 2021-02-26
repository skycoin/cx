package main
//put the repo command stuff here

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	cxcore "github.com/skycoin/cx/cx"
	"github.com/skycoin/cx/cxgo/actions"
	api2 "github.com/skycoin/cx/cxgo/api"
	"github.com/skycoin/cx/cxgo/cxgo"
	"github.com/skycoin/cx/cxgo/cxgo0"
	"github.com/skycoin/cx/cxgo/parser"

	"github.com/skycoin/skycoin/src/util/logging"
)


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
