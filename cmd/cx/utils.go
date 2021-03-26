package main

import "os"

//todo find out why program halt
func parseCmdFlags(options cxCmdFlags, args []string) {

	parseFlags(&options, args)

	// Checking if CXPATH is set, either by setting an environment variable
	// or by setting the `--cxpath` flag.
	GetCXPath(options)

	//checkHelp check command line argumenets
	//$ cx help
	if checkHelp(args) {
		commandLine.PrintDefaults()
		return
	}

	// Does the user want to print the command-line help?
	//options.printHelp works when flags are provided.
	//$ cx --vesion
	if options.printHelp {
		printHelp()
		return
	}

	// Does the user want to print CX's version?
	if options.printVersion {
		printVersion()
		return
	}

	//checkversion check command line argumenets
	//$ cx version
	if checkversion(args) {
		printVersion()
		return
	}

	// User wants to print CX env
	if options.printEnv {
		printEnv()
		return
	}

	//checkenv check command line argumenets
	//$ cx
	if checkenv(args) {
		printEnv()
		return
	}
}

func printlexerandast(args []string, options cxCmdFlags, cxArgs []string, sourceCode []*os.File, fileNames []string) {

	if checkAST(args) {
		printProgramAST(options, cxArgs, sourceCode)
		return
	}

	if options.tokenizeMode {
		printTokenize(options, fileNames)
		return
	}

	if checktokenizeMode(args) {
		printTokenize(options, fileNames)
		return
	}

}
