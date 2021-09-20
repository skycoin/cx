package cxparsing

import (
	"bufio"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	globals2 "github.com/skycoin/cx/cx/globals"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cxparser/actions"
	constants2 "github.com/skycoin/cx/cxparser/constants"
	cxpartialparsing "github.com/skycoin/cx/cxparser/cxpartialparsing"
	"github.com/skycoin/cx/cxparser/globals"
	"github.com/skycoin/cx/cxparser/util/profiling"
)

/**
- take input source code
- parse packages, structs, package imports, globals

*/

// Preliminarystage performs a first pass for the CX cxgo. Globals, packages and
// custom types are added to `cxpartialparsing.Program`.
// takes source strings in one package
// TODO: takes SourcePackage
func Preliminarystage(srcStrs, srcNames []string) int {
	var prePkg *ast.CXPackage
	parseErrors := 0

	profiling.StartProfile("1. packages/structs")
	// 1. Identify all the packages and structs
	for srcI, srcStr := range srcStrs {
		srcName := srcNames[srcI]
		profiling.StartProfile(srcName)
		prePkg, _ = ParsePackages(srcStr, srcName)
		ParseStructs(srcStr, srcName, prePkg)

		profiling.StopProfile(srcName)
	} // for range srcStrs
	profiling.StopProfile("1. packages/structs")

	profiling.StartProfile("2. imports")

	// 1. Identify package imports
	for i, source := range srcStrs {
		ParsePackageImports(source, srcNames[i], prePkg)
	}
	profiling.StopProfile("2. imports")

	profiling.StartProfile("3. globals")
	// 3. Identify all global variables
	//    We also identify packages again, so we know to what
	//    package we're going to add the variable declaration to.
	for i, source := range srcStrs {
		ParseGlobalVariables(source, srcNames[i], prePkg)
	}
	profiling.StopProfile("3. globals")

	profiling.StartProfile("4. cxpartialparsing")

	for i, source := range srcStrs {
		profiling.StartProfile(srcNames[i])
		source = source + "\n"
		if len(srcNames) > 0 {
			cxpartialparsing.CurrentFileName = srcNames[i]
		}
		/*
			passone
		*/
		parseErrors += Passone(source)
		profiling.StopProfile(srcNames[i])
	}

	profiling.StopProfile("4. cxpartialparsing")
	return parseErrors
}

// ParsePackages - stage1
func ParsePackages(srcStr, srcName string) (*ast.CXPackage, int) {
	var prePkg *ast.CXPackage
	parseErrors := 0

	// for parsing comments
	reMultiCommentOpen := regexp.MustCompile(`/\*`)
	reMultiCommentClose := regexp.MustCompile(`\*/`)
	reComment := regexp.MustCompile("//")

	// for parsing packages
	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)

	profiling.StartProfile("1. packages/structs")
	// 1. Identify all the packages and structs
	profiling.StartProfile(srcName)

	reader := strings.NewReader(srcStr)
	scanner := bufio.NewScanner(reader)
	var commentedCode bool
	var lineno = 0
	for scanner.Scan() {
		line := scanner.Bytes()
		lineno++

		// Identify whether we are in a comment or not.
		commentLoc := reComment.FindIndex(line)
		multiCommentOpenLoc := reMultiCommentOpen.FindIndex(line)
		multiCommentCloseLoc := reMultiCommentClose.FindIndex(line)
		if commentedCode && multiCommentCloseLoc != nil {
			commentedCode = false
		}
		if commentedCode {
			continue
		}
		if multiCommentOpenLoc != nil && !commentedCode && multiCommentCloseLoc == nil {
			commentedCode = true
			continue
		}

		// At this point we know that we are *not* in a comment

		// 1a. Identify all the packages
		if loc := rePkg.FindIndex(line); loc != nil {
			if (commentLoc != nil && commentLoc[0] < loc[0]) ||
				(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
				(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
				// then it's commented out
				continue
			}

			if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
				if pkg, err := cxpartialparsing.Program.GetPackage(match[len(match)-1]); err != nil {
					// then it hasn't been added
					pkgName := match[len(match)-1]
					newPkg := ast.MakePackage(pkgName)
					cxpartialparsing.Program.AddPackage(newPkg)
					prePkg = newPkg
					return prePkg, parseErrors
				} else {
					return pkg, parseErrors
				}
			}
		}
	}
	profiling.StopProfile(srcName)
	profiling.StopProfile("1. packages/structs")
	return prePkg, parseErrors
}

// ParseStructs - stage2
func ParseStructs(srcStr, srcName string, prePkg *ast.CXPackage) int {
	parseErrors := 0

	// for parsing comments
	reMultiCommentOpen := regexp.MustCompile(`/\*`)
	reMultiCommentClose := regexp.MustCompile(`\*/`)
	reComment := regexp.MustCompile("//")

	// for parsing structs
	reStrct := regexp.MustCompile("type")
	reStrctName := regexp.MustCompile(`(^|[\s])type\s+([_a-zA-Z][_a-zA-Z0-9]*)?\s`)

	profiling.StartProfile("1. packages/structs")
	// 1. Identify all the packages and structs
	profiling.StartProfile(srcName)

	reader := strings.NewReader(srcStr)
	scanner := bufio.NewScanner(reader)
	var commentedCode bool
	var lineno = 0
	for scanner.Scan() {
		line := scanner.Bytes()
		lineno++

		// Identify whether we are in a comment or not.
		commentLoc := reComment.FindIndex(line)
		multiCommentOpenLoc := reMultiCommentOpen.FindIndex(line)
		multiCommentCloseLoc := reMultiCommentClose.FindIndex(line)
		if commentedCode && multiCommentCloseLoc != nil {
			commentedCode = false
		}
		if commentedCode {
			continue
		}
		if multiCommentOpenLoc != nil && !commentedCode && multiCommentCloseLoc == nil {
			commentedCode = true
			continue
		}

		// At this point we know that we are *not* in a comment

		// 1b. Identify all the structs
		if loc := reStrct.FindIndex(line); loc != nil {
			if (commentLoc != nil && commentLoc[0] < loc[0]) ||
				(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
				(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
				// then it's commented out
				continue
			}

			if match := reStrctName.FindStringSubmatch(string(line)); match != nil {
				if prePkg == nil {
					println(ast.CompilationError(srcName, lineno),
						"No package defined")
				} else if _, err := cxpartialparsing.Program.GetStruct(match[len(match)-1], prePkg.Name); err != nil {
					// then it hasn't been added
					strct := ast.MakeStruct(match[len(match)-1])
					prePkg.AddStruct(strct)
				}
			}
		}
		profiling.StopProfile(srcName)
	} // for range srcStrs
	profiling.StopProfile("1. packages/structs")

	profiling.StopProfile("3. cxpartialparsing")
	return parseErrors
}

func ParseGlobalVariables(source, srcName string, prePkg *ast.CXPackage) {

	reMultiCommentOpen := regexp.MustCompile(`/\*`)
	reMultiCommentClose := regexp.MustCompile(`\*/`)
	reComment := regexp.MustCompile("//")

	// for parsing packages
	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)

	reGlbl := regexp.MustCompile("var")
	reGlblName := regexp.MustCompile(`(^|[\s])var\s([_a-zA-Z][_a-zA-Z0-9]*)`)

	reBodyOpen := regexp.MustCompile("{")
	reBodyClose := regexp.MustCompile("}")

	profiling.StartProfile(srcName)
	// inBlock needs to be 0 to guarantee that we're in the global scope
	var inBlock int
	var commentedCode bool

	scanner := bufio.NewScanner(strings.NewReader(source))
	for scanner.Scan() {
		line := scanner.Bytes()

		// we need to ignore function bodies
		// it'll also ignore struct declaration's bodies, but this doesn't matter
		commentLoc := reComment.FindIndex(line)

		multiCommentOpenLoc := reMultiCommentOpen.FindIndex(line)
		multiCommentCloseLoc := reMultiCommentClose.FindIndex(line)

		if commentedCode && multiCommentCloseLoc != nil {
			commentedCode = false
		}

		if commentedCode {
			continue
		}

		if multiCommentOpenLoc != nil && !commentedCode && multiCommentCloseLoc == nil {
			commentedCode = true
			// continue
		}

		if locs := reBodyOpen.FindAllIndex(line, -1); locs != nil {
			for _, loc := range locs {
				if !(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
					// then it's outside of a */, e.g. `*/ }`
					if (commentLoc == nil && multiCommentOpenLoc == nil && multiCommentCloseLoc == nil) ||
						(commentLoc != nil && commentLoc[0] > loc[0]) ||
						(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] > loc[0]) ||
						(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] < loc[0]) {
						// then we have an uncommented opening bracket
						inBlock++
					}
				}
			}
		}

		if locs := reBodyClose.FindAllIndex(line, -1); locs != nil {
			for _, loc := range locs {
				if !(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
					if (commentLoc == nil && multiCommentOpenLoc == nil && multiCommentCloseLoc == nil) ||
						(commentLoc != nil && commentLoc[0] > loc[0]) ||
						(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] > loc[0]) ||
						(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] < loc[0]) {
						// then we have an uncommented closing bracket
						inBlock--
					}
				}
			}
		}

		// we could have this situation: {var local i32}
		// but we don't care about this, as the later passes will throw an error as it's invalid syntax

		if loc := rePkg.FindIndex(line); loc != nil {
			if (commentLoc != nil && commentLoc[0] < loc[0]) ||
				(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
				(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
				// then it's commented out
				continue
			}

			if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
				if pkg, err := cxpartialparsing.Program.GetPackage(match[len(match)-1]); err != nil {
					// it should be already present
					panic(err)
				} else {
					prePkg = pkg
				}
			}
		}

		// finally, if we read a "var" and we're in global scope, we add the global without any type
		// the type will be determined later on
		if loc := reGlbl.FindIndex(line); loc != nil {
			if (commentLoc != nil && commentLoc[0] < loc[0]) ||
				(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
				(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) || inBlock != 0 {
				// then it's commented out or inside a block
				continue
			}

			if match := reGlblName.FindStringSubmatch(string(line)); match != nil {
				if _, err := prePkg.GetGlobal(match[len(match)-1]); err != nil {
					// then it hasn't been added
					arg := ast.MakeArgument(match[len(match)-1], "", 0)
					arg.Offset = types.InvalidPointer
					arg.ArgDetails.Package = prePkg
					prePkg.AddGlobal(arg)
				}
			}
		}
	}
	profiling.StopProfile(srcName)
}

func ParsePackageImports(source, srcName string, prePkg *ast.CXPackage) {

	reMultiCommentOpen := regexp.MustCompile(`/\*`)
	reMultiCommentClose := regexp.MustCompile(`\*/`)
	reComment := regexp.MustCompile("//")

	reImp := regexp.MustCompile("import")
	reImpName := regexp.MustCompile(`(^|[\s])import\s+"([_a-zA-Z][_a-zA-Z0-9/-]*)"`)
	profiling.StartProfile(srcName)

	var commentedCode bool

	scanner := bufio.NewScanner(strings.NewReader(source))
	for scanner.Scan() {
		line := scanner.Bytes()

		// we need to ignore function bodies
		// it'll also ignore struct declaration's bodies, but this doesn't matter
		commentLoc := reComment.FindIndex(line)

		multiCommentOpenLoc := reMultiCommentOpen.FindIndex(line)
		multiCommentCloseLoc := reMultiCommentClose.FindIndex(line)

		if commentedCode && multiCommentCloseLoc != nil {
			commentedCode = false
		}

		if commentedCode {
			continue
		}

		if multiCommentOpenLoc != nil && !commentedCode && multiCommentCloseLoc == nil {
			commentedCode = true
			// continue
		}

		// TODO: move this out to its own function
		// Identify all the package imports.
		if loc := reImp.FindIndex(line); loc != nil {
			if (commentLoc != nil && commentLoc[0] < loc[0]) ||
				(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
				(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
				// then it's commented out
				continue
			}

			if match := reImpName.FindStringSubmatch(string(line)); match != nil {
				pkgName := match[len(match)-1]
				// Checking if `pkgName` already exists and if it's not a standard library package.
				if _, err := cxpartialparsing.Program.GetPackage(pkgName); err != nil && !constants2.IsCorePackage(pkgName) {
					// _, sourceCode, srcNames := ParseArgsForCX([]string{fmt.Sprintf("%s%s", SRCPATH, pkgName)}, false)
					_, sourceCode, fileNames := ast.ParseArgsForCX([]string{filepath.Join(globals2.SRCPATH, pkgName)}, false)
					ParseSourceCode(sourceCode, fileNames)
				}
			}
		}

	}
	profiling.StopProfile(srcName)
}

func AddInitFunction(prgrm *ast.CXProgram) error {
	mainPkg, err := prgrm.GetPackage(constants.MAIN_PKG)
	if err != nil {
		return err
	}

	initFn := ast.MakeFunction(constants.SYS_INIT_FUNC, actions.CurrentFile, actions.LineNo)
	mainPkg.AddFunction(initFn)

	//Init Expressions
	actions.FunctionDeclaration(initFn, nil, nil, globals.SysInitExprs)

	if _, err := mainPkg.SelectFunction(constants.MAIN_FUNC); err != nil {
		return err
	}
	return nil
}
