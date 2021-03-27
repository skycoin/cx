package cxparser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/skycoin/cx/cxgo/globals"

	cxcore "github.com/skycoin/cx/cx"
	"github.com/skycoin/cx/cxgo/actions"
	"github.com/skycoin/cx/cxgo/cxgo"
	"github.com/skycoin/cx/cxgo/cxgo0"
	"github.com/skycoin/cx/cxgo/util/profiling"
)

// ParseSourceCode takes a group of files representing CX `sourceCode` and
// parses it into CX program structures for `PRGRM`.
func ParseSourceCode(sourceCode []*os.File, fileNames []string) {

	//local
	cxgo0.PRGRM0 = actions.PRGRM

	programsourceCode := make([]string, len(sourceCode))
	for i, source := range sourceCode {
		tmp := bytes.NewBuffer(nil)
		io.Copy(tmp, source)
		programsourceCode[i] = tmp.String()

	}

	/*
		We need to traverse the elements by hierarchy first add all the
		packages and structs at the same time then add globals, as these
		can be of a custom type (and it could be imported) the signatures
		of functions and methods are added in the cxgo0.y pass
	*/

	parseErrors := 0
	if len(sourceCode) > 0 {
		parseErrors = lexer(programsourceCode, fileNames)
	}

	//package level program
	actions.PRGRM.SetCurrentCxProgram()

	actions.PRGRM = cxgo0.PRGRM0

	if cxcore.FoundCompileErrors || parseErrors > 0 {
		profiling.CleanupAndExit(cxcore.CX_COMPILATION_ERROR)
	}

	// Adding global variables `OS_ARGS` to the `os` (operating system)
	// package.
	if osPkg, err := actions.PRGRM.GetPackage(cxcore.OS_PKG); err == nil {
		if _, err := osPkg.GetGlobal(cxcore.OS_ARGS); err != nil {
			arg0 := cxcore.MakeArgument(cxcore.OS_ARGS, "", -1).AddType(cxcore.TypeNames[cxcore.TYPE_UNDEFINED])
			arg0.Package = osPkg

			arg1 := cxcore.MakeArgument(cxcore.OS_ARGS, "", -1).AddType(cxcore.TypeNames[cxcore.TYPE_STR])
			arg1 = actions.DeclarationSpecifiers(arg1, []int{0}, cxcore.DECL_BASIC)
			arg1 = actions.DeclarationSpecifiers(arg1, []int{0}, cxcore.DECL_SLICE)

			actions.DeclareGlobalInPackage(osPkg, arg0, arg1, nil, false)
		}
	}

	profiling.StartProfile("4. parse")
	// The last pass of parsing that generates the actual output.
	for i, source := range programsourceCode {
		// Because of an unkown reason, sometimes some CX programs
		// throw an error related to a premature EOF (particularly in Windows).
		// Adding a newline character solves this.
		source = source + "\n"
		actions.LineNo = 1
		b := bytes.NewBufferString(source)
		if len(fileNames) > 0 {
			actions.CurrentFile = fileNames[i]
		}
		profiling.StartProfile(actions.CurrentFile)
		parseErrors += cxgo.Parse(cxgo.NewLexer(b))
		profiling.StopProfile(actions.CurrentFile)
	}
	profiling.StopProfile("4. parse")

	if cxcore.FoundCompileErrors || parseErrors > 0 {
		profiling.CleanupAndExit(cxcore.CX_COMPILATION_ERROR)
	}
}

/* function lexer performs a first pass for the CX cxgo.
Globals, packages and custom types are parsed added to `cxgo0.PRGRM0`.
*/
func lexer(programSourceCode, programFileNames []string) int {

	fmt.Println("###############")
	fmt.Println(programSourceCode)

	fmt.Println("###############")

	fmt.Println(programFileNames)

	var prePkg *cxcore.CXPackage
	parseErrors := 0

	re := newRegulaEexpression()

	lexstruct(programSourceCode, programFileNames, prePkg, re)

	lexglobal(programSourceCode, programFileNames, prePkg, re)

	profiling.StartProfile("1. packages/structs")

	profiling.StartProfile("3. cxgo0")
	// cxgo0.Parse(allSC)
	for i, source := range programSourceCode {
		profiling.StartProfile(programFileNames[i])
		source = source + "\n"
		if len(programFileNames) > 0 {
			cxgo0.CurrentFileName = programFileNames[i]
		}
		//Parse calls yyParse
		parseErrors += cxgo0.Parse(source)
		profiling.StopProfile(programFileNames[i])
	}
	profiling.StopProfile("3. cxgo0")
	return parseErrors
}

func AddInitFunction(prgrm *cxcore.CXProgram) error {
	mainPkg, err := prgrm.GetPackage(cxcore.MAIN_PKG)
	if err != nil {
		return err
	}

	initFn := cxcore.MakeFunction(cxcore.SYS_INIT_FUNC, actions.CurrentFile, actions.LineNo)
	mainPkg.AddFunction(initFn)

	//Init Expressions
	actions.FunctionDeclaration(initFn, nil, nil, globals.SysInitExprs)

	if _, err := mainPkg.SelectFunction(cxcore.MAIN_FUNC); err != nil {
		return err
	}
	return nil
}

func lexstruct(programSourceCode []string, programFileNames []string, prePkg *cxcore.CXPackage, re RegularExpression) {

	// 1. Identify all the packages and structs
	for srcI, srcStr := range programSourceCode {

		srcName := programFileNames[srcI]
		profiling.StartProfile(srcName)

		reader := strings.NewReader(srcStr)
		scanner := bufio.NewScanner(reader)
		var commentedCode bool
		var lineno = 0
		for scanner.Scan() {
			line := scanner.Bytes()
			lineno++

			// Identify whether we are in a comment or not.
			commentLoc := re.reComment.FindIndex(line)
			multiCommentOpenLoc := re.reMultiCommentOpen.FindIndex(line)
			multiCommentCloseLoc := re.reMultiCommentClose.FindIndex(line)
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
			if loc := re.rePackage.FindIndex(line); loc != nil {
				if (commentLoc != nil && commentLoc[0] < loc[0]) ||
					(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
					(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
					// then it's commented out
					continue
				}

				if match := re.rePackageName.FindStringSubmatch(string(line)); match != nil {
					if pkg, err := cxgo0.PRGRM0.GetPackage(match[len(match)-1]); err != nil {
						// then it hasn't been added
						newPkg := cxcore.MakePackage(match[len(match)-1])
						cxgo0.PRGRM0.AddPackage(newPkg)
						prePkg = newPkg
					} else {
						prePkg = pkg
					}
				}
			}

			// 1b. Identify all the structs
			if loc := re.reStruct.FindIndex(line); loc != nil {
				if (commentLoc != nil && commentLoc[0] < loc[0]) ||
					(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
					(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
					// then it's commented out
					continue
				}

				if match := re.reStructName.FindStringSubmatch(string(line)); match != nil {
					if prePkg == nil {
						println(cxcore.CompilationError(srcName, lineno),
							"No package defined")
					} else if _, err := cxgo0.PRGRM0.GetStruct(match[len(match)-1], prePkg.Name); err != nil {
						// then it hasn't been added
						strct := cxcore.MakeStruct(match[len(match)-1])
						prePkg.AddStruct(strct)
					}
				}
			}
		}
		profiling.StopProfile(srcName)
	} // for range srcStrs
	profiling.StopProfile("1. packages/structs")

}

func lexglobal(programSourceCode []string, programFileNames []string, prePkg *cxcore.CXPackage, re RegularExpression) {

	profiling.StartProfile("2. globals")
	// 2. Identify all global variables
	//    We also identify packages again, so we know to what
	//    package we're going to add the variable declaration to.
	for i, source := range programSourceCode {
		profiling.StartProfile(programFileNames[i])
		// inBlock needs to be 0 to guarantee that we're in the global scope
		var inBlock int
		var commentedCode bool

		scanner := bufio.NewScanner(strings.NewReader(source))
		for scanner.Scan() {
			line := scanner.Bytes()

			// we need to ignore function bodies
			// it'll also ignore struct declaration's bodies, but this doesn't matter
			commentLoc := re.reComment.FindIndex(line)

			multiCommentOpenLoc := re.reMultiCommentOpen.FindIndex(line)
			multiCommentCloseLoc := re.reMultiCommentClose.FindIndex(line)

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

			// Identify all the package imports.
			if loc := re.reImport.FindIndex(line); loc != nil {
				if (commentLoc != nil && commentLoc[0] < loc[0]) ||
					(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
					(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
					// then it's commented out
					continue
				}

				if match := re.reImportName.FindStringSubmatch(string(line)); match != nil {
					pkgName := match[len(match)-1]
					// Checking if `pkgName` already exists and if it's not a standard library package.
					if _, err := cxgo0.PRGRM0.GetPackage(pkgName); err != nil && !cxcore.IsCorePackage(pkgName) {
						// _, sourceCode, srcNames := ParseArgsForCX([]string{fmt.Sprintf("%s%s", SRCPATH, pkgName)}, false)
						_, sourceCode, fileNames := cxcore.ParseArgsForCX([]string{filepath.Join(cxcore.SRCPATH, pkgName)}, false)
						ParseSourceCode(sourceCode, fileNames)
					}
				}
			}

			// we search for packages at the same time, so we can know to what package to add the global
			if loc := re.rePackage.FindIndex(line); loc != nil {
				if (commentLoc != nil && commentLoc[0] < loc[0]) ||
					(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
					(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
					// then it's commented out
					continue
				}

				if match := re.rePackageName.FindStringSubmatch(string(line)); match != nil {
					if pkg, err := cxgo0.PRGRM0.GetPackage(match[len(match)-1]); err != nil {
						// then it hasn't been added
						prePkg = cxcore.MakePackage(match[len(match)-1])
						cxgo0.PRGRM0.AddPackage(prePkg)
					} else {
						prePkg = pkg
					}
				}
			}

			if locs := re.reBodyOpen.FindAllIndex(line, -1); locs != nil {
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

			if locs := re.reBodyClose.FindAllIndex(line, -1); locs != nil {
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

			if loc := re.rePackage.FindIndex(line); loc != nil {
				if (commentLoc != nil && commentLoc[0] < loc[0]) ||
					(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
					(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
					// then it's commented out
					continue
				}

				if match := re.rePackageName.FindStringSubmatch(string(line)); match != nil {
					if pkg, err := cxgo0.PRGRM0.GetPackage(match[len(match)-1]); err != nil {
						// it should be already present
						panic(err)
					} else {
						prePkg = pkg
					}
				}
			}

			// finally, if we read a "var" and we're in global scope, we add the global without any type
			// the type will be determined later on
			if loc := re.reGlobal.FindIndex(line); loc != nil {
				if (commentLoc != nil && commentLoc[0] < loc[0]) ||
					(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
					(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) || inBlock != 0 {
					// then it's commented out or inside a block
					continue
				}

				if match := re.reGlobalName.FindStringSubmatch(string(line)); match != nil {
					if _, err := prePkg.GetGlobal(match[len(match)-1]); err != nil {
						// then it hasn't been added
						arg := cxcore.MakeArgument(match[len(match)-1], "", 0)
						arg.Offset = -1
						arg.Package = prePkg
						prePkg.AddGlobal(arg)
					}
				}
			}
		}
		profiling.StopProfile(programFileNames[i])
	}
	profiling.StopProfile("2. globals")
}
