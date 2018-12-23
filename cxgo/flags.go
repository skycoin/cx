package main

import (
	"flag"
	"os"
)

type cxCmdFlags struct {
	baseOutput          bool
	compileMode         bool
	compileOutput       string
	initProject         bool
	replMode            bool
	signalClientMode    bool
	signalClientID      int
	signalServerAddress string
	webMode             bool
	ideMode             bool
	flagMode            bool
	printHelp           bool
	printVersion        bool
	interpretMode       bool
}

func defaultCmdFlags() cxCmdOptions {
	return cxCmdFlags{
		baseOutput:          false,
		compileMode:         false,
		compileOutput:       "",
		newProject:          false,
		replMode:            false,
		signalClientMode:    false,
		signalClientID:      1,
		signalServerAddress: "localhost:7999",
		webMode:             false,
		ideMode:             false,
		printHelp:           false,
		printVersion:        false,
		interpretMode:       false,
	}
}

func registerFlags(options *cxCmdFlags) {
	args := os.Args
	if len(args) <= 1 {
		options.replMode = true
	}

	flag.BoolVar(&printVersion, "version", printVersion, "Print CX version")
	flag.BoolVar(&printVersion, "v", printVersion, "alias for -version")
	flag.BoolVar(&printHelp, "help", printHelp, "Print CX version")
	flag.BoolVar(&printHelp, "h", printHelp, "alias for -help")
	flag.BoolVar(&options.baseOutput, "base", options.baseOutput, "generate a 'out.cx.go' file with the transcompiled CX Base source code.")
	flag.BoolVar(&options.baseOutput, "b", options.baseOutput, "alias for -base")
	flag.BoolVar(&options.compileMode, "compile", options.compileMode, "generate a 'out' executable file of the program")
	flag.BoolVar(&options.compileMode, "c", options.compileMode, "alias for -compile")
	flag.BoolVar(&options.compileOutput, "compile-output", options.compileOutput, "Specifies the filename for the generated executable.")
	flag.StringVar(&options.compileOutput, "co", options.compileOutput, "alias for -compile-output")
	flag.BoolVar(&options.initProject, "new", options.initProject, "Creates a new project located at $CXPATH/src")
	flag.BoolVar(&options.initProject, "n", options.initProject, "alias for -new")
	flag.BoolVar(&options.replMode, "repl", options.replMode, "Loads source files into memory and starts a read-eval-print loop")
	flag.BoolVar(&options.replMode, "r", options.replMode, "alias for -repl")
	flag.BoolVar(&options.webMode, "web", options.webMode, "Start CX as a web service.")
	flag.BoolVar(&options.webMode, "w", options.webMode, "alias for -web")
	flag.BoolVar(&options.ideMode, "ide", options.ideMode, "Start CX as a web service, and Leaps service start also.")
	// viscript options
	flag.BoolVar(&options.signalClientMode, "signal-client", options.signalClientMode, "Run signal client")
	flag.IntVar(&options.signalClientID, "signal-client-id", options.signalClientID, "Id of signal client (default 1)")
	flag.StringVar(&options.signalSeverAddress, "signal-client-address", options.signalSeverAddress, "Address of signal server (default 'localhost:7999')")
}

func printHelp() {
	fmt.Printf(`Usage: cx [options] [source-files]

CX options:
-b, --base                        Generate a "out.cx.go" file with the transcompiled CX Base source code.
-c, --compile                     Generate a "out" executable file of the program.
-co, --compile-output FILENAME    Specifies the filename for the generated executable.
-h, --help                        Prints this message.
-n, --new                         Creates a new project located at $CXPATH/src
-r, --repl                        Loads source files into memory and starts a read-eval-print loop.
-w, --web                         Start CX as a web service.
-ide, --ide						  Start CX as a web service, and Leaps service start also.

Signal options:
-signal-client                   Run signal client
-signal-client-id UINT           Id of signal client (default 1)
-signal-server-address STRING    Address of signal server (default "localhost:7999")

Notes:
* Options --compile and --repl are mutually exclusive.
* Option --web makes every other flag to be ignored.
`)
}

func parseArgsForCX(args []string) (cxArgs []string, sourceCode []*os.File, fileNames []string) {
	for i, arg := range args {
		if len(arg) > 2 && arg[:2] == "++" {
			cxArgs = append(cxArgs, arg)
			continue
		}
		if !flagMode {

			fi, err := os.Stat(arg)
			_ = err

			if err != nil {
				panic(err)
			}

			switch mode := fi.Mode(); {
			case mode.IsDir():
				var fileList []string

				err := filepath.Walk(arg, func(path string, f os.FileInfo, err error) error {
					fileList = append(fileList, path)
					return nil
				})

				if err != nil {
					panic(err)
				}

				for _, path := range fileList {
					file, err := os.Open(path)

					if err != nil {
						panic(err)
					}

					fiName := file.Name()
					fiNameLen := len(fiName)

					if fiNameLen > 2 && fiName[fiNameLen-3:] == ".cx" {
						// only loading .cx files
						sourceCode = append(sourceCode, file)
						fileNames = append(fileNames, fiName)
					}
				}
			case mode.IsRegular():
				file, err := os.Open(arg)

				if err != nil {
					panic(err)
				}

				fileNames = append(fileNames, file.Name())
				sourceCode = append(sourceCode, file)
			}
		}
	}
}

func printVersion() {
	fmt.Println("CX version", VERSION)
}

func parseCxArgs(args []string) {
	a
}
