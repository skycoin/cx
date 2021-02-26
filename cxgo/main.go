package main

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

const VERSION = "0.7.1"

var (
	logger          = logging.MustGetLogger("CX")

)

func getJSON(url string, target interface{}) error {
	r, err := apiClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// optionTokenize checks if the user wants to use CX to generate the lexer tokens
func optionTokenize(options cxCmdFlags, fileNames []string) {
	var r *os.File
	var w *os.File
	var err error

	if len(fileNames) == 0 {
		r = os.Stdin
	} else {
		sourceFilename := fileNames[0]
		if len(fileNames) > 1 {
			fmt.Fprintln(os.Stderr, "Multiple source files detected. Ignoring all except", sourceFilename)
		}
		r, err = cxcore.CXOpenFile(sourceFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading:", sourceFilename, err)
			return
		}
		defer r.Close()
	}

	if options.compileOutput == "" {
		w = os.Stdout
	} else {
		tokenFilename := options.compileOutput
		w, err = cxcore.CXCreateFile(tokenFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error writing:", tokenFilename, err)
			return
		}
		defer w.Close()
	}

	parser.Tokenize(r, w)
}

func cleanupAndExit(exitCode int) {
	StopCPUProfile(profile)
	os.Exit(exitCode)
}

func parseProgram(options cxCmdFlags, fileNames []string, sourceCode []*os.File) (bool, []byte, []byte) {
	profile := StartCPUProfile("parse")
	defer StopCPUProfile(profile)

	defer DumpMEMProfile("parse")

	StartProfile("parse")
	defer StopProfile("parse")

	actions.PRGRM = cxcore.MakeProgram()
	corePkgsPrgrm, err := cxcore.GetProgram()
	if err != nil {
		panic(err)
	}
	actions.PRGRM.Packages = corePkgsPrgrm.Packages

	if options.webMode {
		ServiceMode()
		return false, nil, nil
	}

	// TODO @evanlinjin: Do we need this? What is the 'leaps' command?
	if options.ideMode {
		IdeServiceMode()
		ServiceMode()
		return false, nil, nil
	}

	// TODO @evanlinjin: We do not need a persistent mode?
	if options.webPersistentMode {
		go ServiceMode()
		PersistentServiceMode()
		return false, nil, nil
	}

	// TODO @evanlinjin: This is a separate command now.
	if options.tokenizeMode {
		optionTokenize(options, fileNames)
		return false, nil, nil
	}

	// var bcPrgrm *CXProgram
	var sPrgrm []byte
	// In case of a CX chain, we need to temporarily store the blockchain code heap elsewhere,
	// so we can then add it after the transaction code's data segment.
	var bcHeap []byte
	if options.transactionMode || options.broadcastMode {
		chainStatePrelude(&sPrgrm, &bcHeap, actions.PRGRM) // TODO: refactor injection logic
	}

	// Parsing all the source code files sent as CLI arguments to CX.
	cxgo.ParseSourceCode(sourceCode, fileNames)

	// setting project's working directory
	if !options.replMode && len(sourceCode) > 0 {
		cxgo0.PRGRM0.Path = determineWorkDir(sourceCode[0].Name())
	}

	// Checking if a main package exists. If not, create and add it to `PRGRM`.
	if _, err := actions.PRGRM.GetFunction(cxcore.MAIN_FUNC, cxcore.MAIN_PKG); err != nil {
		initMainPkg(actions.PRGRM)
	}
	// Setting what function to start in if using the REPL.
	actions.ReplTargetFn = cxcore.MAIN_FUNC

	// Adding *init function that initializes all the global variables.
	cxgo.AddInitFunction(actions.PRGRM)

	actions.LineNo = 0

	if cxcore.FoundCompileErrors {
		cleanupAndExit(cxcore.CX_COMPILATION_ERROR)
	}

	return true, bcHeap, sPrgrm
}

func runProgram(options cxCmdFlags, cxArgs []string, sourceCode []*os.File, bcHeap []byte, sPrgrm []byte) {
	StartProfile("run")
	defer StopProfile("run")

	if options.replMode || len(sourceCode) == 0 {
		actions.PRGRM.SelectProgram()
		repl()
		return
	}

	// If it's a CX chain transaction, we need to add the heap extracted
	// from the retrieved CX chain program state.
	if options.transactionMode || options.broadcastMode {
		mergeBlockchainHeap(bcHeap, sPrgrm) // TODO: refactor injection logic
	}

	if options.blockchainMode {
		panic("blockchainMode is moved to the github.com/skycoin/cx-chains repo")
	} else if options.broadcastMode {
		panic("broadcastMode is moved to the github.com/skycoin/cx-chains repo")

	} else {
		// Normal run of a CX program.
		err := actions.PRGRM.RunCompiled(0, cxArgs)
		if err != nil {
			panic(err)
		}

		if cxcore.AssertFailed() {
			os.Exit(cxcore.CX_ASSERT)
		}
	}
}

func Run(args []string) {
	runtime.LockOSThread()
	runtime.GOMAXPROCS(2)

	options := defaultCmdFlags()
	parseFlags(&options, args)

	// Checking if CXPATH is set, either by setting an environment variable
	// or by setting the `--cxpath` flag.
	checkCXPathSet(options)

	if checkHelp(args) {
		commandLine.PrintDefaults()
		return

	}

	// Does the user want to print the command-line help?
	if options.printHelp {
		printHelp()
		return
	}

	// Does the user want to print CX's version?
	if options.printVersion {
		printVersion()
		return
	}

	// User wants to print CX env
	if options.printEnv {
		printEnv()
		return
	}

	if options.initialHeap != "" {
		cxcore.INIT_HEAP_SIZE = parseMemoryString(options.initialHeap)
	}
	if options.maxHeap != "" {
		cxcore.MAX_HEAP_SIZE = parseMemoryString(options.maxHeap)
		if cxcore.MAX_HEAP_SIZE < cxcore.INIT_HEAP_SIZE {
			// Then MAX_HEAP_SIZE overrides INIT_HEAP_SIZE's value.
			cxcore.INIT_HEAP_SIZE = cxcore.MAX_HEAP_SIZE
		}
	}
	if options.stackSize != "" {
		cxcore.STACK_SIZE = parseMemoryString(options.stackSize)
		actions.DataOffset = cxcore.STACK_SIZE
	}
	if options.minHeapFreeRatio != float64(0) {
		cxcore.MIN_HEAP_FREE_RATIO = float32(options.minHeapFreeRatio)
	}
	if options.maxHeapFreeRatio != float64(0) {
		cxcore.MAX_HEAP_FREE_RATIO = float32(options.maxHeapFreeRatio)
	}

	// options, file pointers, filenames
	cxArgs, sourceCode, fileNames := cxcore.ParseArgsForCX(commandLine.Args(), true)

	// Propagate some options out to other packages.
	parser.DebugLexer = options.debugLexer // in package parser
	DebugProfileRate = options.debugProfile
	DebugProfile = DebugProfileRate > 0

	if run, bcHeap, sPrgrm := parseProgram(options, fileNames, sourceCode); run {
		runProgram(options, cxArgs, sourceCode, bcHeap, sPrgrm)
	}
}

// Used for the -heap-initial, -heap-max and -stack-size flags.
// This function parses, for example, "1M" to 1048576 (the corresponding number of bytes)
// Possible suffixes are: G or g (gigabytes), M or m (megabytes), K or k (kilobytes)

//IS NOT CALLED BY ANYTHING
/*
func parseMemoryString(s string) int {
	suffix := s[len(s)-1]
	_, notSuffix := strconv.ParseFloat(string(suffix), 64)

	if notSuffix == nil {
		// then we don't have a suffix
		num, err := strconv.ParseInt(s, 10, 64)

		if err != nil {
			// malformed size
			return -1
		}

		return int(num)
	} else {
		// then we have a suffix
		num, err := strconv.ParseFloat(s[:len(s)-1], 64)

		if err != nil {
			// malformed size
			return -1
		}

		// The user can use suffixes to give as input gigabytes, megabytes or kilobytes.
		switch suffix {
		case 'G', 'g':
			return int(num * 1073741824)
		case 'M', 'm':
			return int(num * 1048576)
		case 'K', 'k':
			return int(num * 1024)
		default:
			return -1
		}
	}
}
*/
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

type SourceCode struct {
	Code string
}

func ServiceMode() {
	host := ":5336"

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./dist")))
	mux.Handle("/program/", api2.NewAPI("/program", actions.PRGRM))
	mux.HandleFunc("/eval", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var b []byte
		var err error
		if b, err = ioutil.ReadAll(r.Body); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var source SourceCode
		if err := json.Unmarshal(b, &source); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if err := r.ParseForm(); err == nil {
			fmt.Fprintf(w, "%s", Eval(source.Code+"\n"))
		}
	})

	if listener, err := net.Listen("tcp", host); err == nil {
		fmt.Println("Starting CX web service on http://127.0.0.1:5336/")
		http.Serve(listener, mux)
	}
}

func IdeServiceMode() {
	// Leaps's host address
	ideHost := "localhost:5335"

	// Working directory for ide
	sharedPath := fmt.Sprintf("%s/src/github.com/skycoin/cx", os.Getenv("GOPATH"))

	// Start Leaps
	// cmd = `leaps -address localhost:5335 $GOPATH/src/skycoin/cx`
	cmnd := exec.Command("leaps", "-address", ideHost, sharedPath)

	// Just leave start command
	cmnd.Start()
}

func PersistentServiceMode() {
	fmt.Println("Start persistent for service mode!")

	fi := bufio.NewReader(os.Stdin)

	for {
		var inp string
		var ok bool

		printPrompt()

		if inp, ok = readline(fi); ok {
			if isJSON(inp) {
				var err error
				client := &http.Client{}
				body := bytes.NewBufferString(inp)
				req, err := http.NewRequest("GET", "http://127.0.0.1:5336/eval", body)
				if err != nil {
					fmt.Println(err)
					return
				}

				if resp, err := client.Do(req); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(resp.Status)
				}
			}
		}
	}
}

func determineWorkDir(filename string) string {
	filename = filepath.FromSlash(filename)

	i := strings.LastIndexByte(filename, os.PathSeparator)
	if i == -1 {
		i = 0
	}
	return filename[:i]
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

func repl() {
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

// checkCXPathSet checks if the user has set the environment variable
// `CXPATH`. If not, CX creates a workspace at $HOME/cx, along with $HOME/cx/bin,
// $HOME/cx/pkg and $HOME/cx/src
func checkCXPathSet(options cxCmdFlags) {
	// Determining the filepath of the directory where the user
	// started the `cx` command.
	_, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// cxcore.COREPATH = filepath.Dir(ex) // TODO @evanlinjin: Not used.

	CXPATH := ""
	if os.Getenv("CXPATH") != "" {
		CXPATH = os.Getenv("CXPATH")
	}
	// `options.cxpath` overrides `os.Getenv("CXPATH")`
	if options.cxpath != "" {
		CXPATH, err = filepath.Abs(options.cxpath)
		if err != nil {
			panic(err)
		}
	}
	if os.Getenv("CXPATH") == "" && options.cxpath == "" {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}

		CXPATH = usr.HomeDir + "/cx/"
	}

	cxcore.BINPATH = filepath.Join(CXPATH, "bin/")
	cxcore.PKGPATH = filepath.Join(CXPATH, "pkg/")
	cxcore.SRCPATH = filepath.Join(CXPATH, "src/")

	// Creating directories in case they do not exist.
	if _, err := cxcore.CXStatFile(CXPATH); os.IsNotExist(err) {
		cxcore.CXMkdirAll(CXPATH, 0755)
	}
	if _, err := cxcore.CXStatFile(cxcore.BINPATH); os.IsNotExist(err) {
		cxcore.CXMkdirAll(cxcore.BINPATH, 0755)
	}
	if _, err := cxcore.CXStatFile(cxcore.PKGPATH); os.IsNotExist(err) {
		cxcore.CXMkdirAll(cxcore.PKGPATH, 0755)
	}
	if _, err := cxcore.CXStatFile(cxcore.SRCPATH); os.IsNotExist(err) {
		cxcore.CXMkdirAll(cxcore.SRCPATH, 0755)
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

