package cxgo

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	cxcore "github.com/skycoin/cx/cx"
	"github.com/skycoin/cx/cxgo/actions"
	"github.com/skycoin/cx/cxgo/cxgo0"
	"github.com/skycoin/cx/cxgo/parser"
	"github.com/skycoin/cx/cxgo/profiling"
)

// ParseSourceCode takes a group of files representing CX `sourceCode` and
// parses it into CX program structures for `PRGRM`.
func ParseSourceCode(sourceCode []*os.File, fileNames []string) {
	cxgo0.PRGRM0 = actions.PRGRM

	// Copy the contents of the file pointers containing the CX source
	// code into sourceCodeCopy
	sourceCodeCopy := make([]string, len(sourceCode))
	for i, source := range sourceCode {
		tmp := bytes.NewBuffer(nil)
		io.Copy(tmp, source)
		sourceCodeCopy[i] = string(tmp.Bytes())
	}

	// We need to traverse the elements by hierarchy first add all the
	// packages and structs at the same time then add globals, as these
	// can be of a custom type (and it could be imported) the signatures
	// of functions and methods are added in the cxgo0.y pass
	parseErrors := 0
	if len(sourceCode) > 0 {
		parseErrors = lexerStep0(sourceCodeCopy, fileNames)
	}

	actions.PRGRM.SelectProgram()

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
	for i, source := range sourceCodeCopy {
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
		parseErrors += parser.Parse(parser.NewLexer(b))
		profiling.StopProfile(actions.CurrentFile)
	}
	profiling.StopProfile("4. parse")

	if cxcore.FoundCompileErrors || parseErrors > 0 {
		profiling.CleanupAndExit(cxcore.CX_COMPILATION_ERROR)
	}
}

// lexerStep0 performs a first pass for the CX parser. Globals, packages and
// custom types are added to `cxgo0.PRGRM0`.
func lexerStep0(srcStrs, srcNames []string) int {
	var prePkg *cxcore.CXPackage
	parseErrors := 0

	reMultiCommentOpen := regexp.MustCompile(`/\*`)
	reMultiCommentClose := regexp.MustCompile(`\*/`)
	reComment := regexp.MustCompile("//")

	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile("(^|[\\s])package\\s+([_a-zA-Z][_a-zA-Z0-9]*)")
	reStrct := regexp.MustCompile("type")
	reStrctName := regexp.MustCompile("(^|[\\s])type\\s+([_a-zA-Z][_a-zA-Z0-9]*)?\\s")

	reGlbl := regexp.MustCompile("var")
	reGlblName := regexp.MustCompile("(^|[\\s])var\\s([_a-zA-Z][_a-zA-Z0-9]*)")

	reBodyOpen := regexp.MustCompile("{")
	reBodyClose := regexp.MustCompile("}")

	reImp := regexp.MustCompile("import")
	reImpName := regexp.MustCompile("(^|[\\s])import\\s+\"([_a-zA-Z][_a-zA-Z0-9/-]*)\"")

	profiling.StartProfile("1. packages/structs")
	// 1. Identify all the packages and structs
	for srcI, srcStr := range srcStrs {
		srcName := srcNames[srcI]
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
			if loc := reStrct.FindIndex(line); loc != nil {
				if (commentLoc != nil && commentLoc[0] < loc[0]) ||
					(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
					(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
					// then it's commented out
					continue
				}

				if match := reStrctName.FindStringSubmatch(string(line)); match != nil {
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

	profiling.StartProfile("2. globals")
	// 2. Identify all global variables
	//    We also identify packages again, so we know to what
	//    package we're going to add the variable declaration to.
	for i, source := range srcStrs {
		profiling.StartProfile(srcNames[i])
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
					if _, err := cxgo0.PRGRM0.GetPackage(pkgName); err != nil && !cxcore.IsCorePackage(pkgName) {
						// _, sourceCode, srcNames := ParseArgsForCX([]string{fmt.Sprintf("%s%s", SRCPATH, pkgName)}, false)
						_, sourceCode, fileNames := cxcore.ParseArgsForCX([]string{filepath.Join(cxcore.SRCPATH, pkgName)}, false)
						ParseSourceCode(sourceCode, fileNames)
					}
				}
			}

			// we search for packages at the same time, so we can know to what package to add the global
			if loc := rePkg.FindIndex(line); loc != nil {
				if (commentLoc != nil && commentLoc[0] < loc[0]) ||
					(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
					(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
					// then it's commented out
					continue
				}

				if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
					if pkg, err := cxgo0.PRGRM0.GetPackage(match[len(match)-1]); err != nil {
						// then it hasn't been added
						prePkg = cxcore.MakePackage(match[len(match)-1])
						cxgo0.PRGRM0.AddPackage(prePkg)
					} else {
						prePkg = pkg
					}
				}
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
						arg := cxcore.MakeArgument(match[len(match)-1], "", 0)
						arg.Offset = -1
						arg.Package = prePkg
						prePkg.AddGlobal(arg)
					}
				}
			}
		}
		profiling.StopProfile(srcNames[i])
	}
	profiling.StopProfile("2. globals")

	profiling.StartProfile("3. cxgo0")
	// cxgo0.Parse(allSC)
	for i, source := range srcStrs {
		profiling.StartProfile(srcNames[i])
		source = source + "\n"
		if len(srcNames) > 0 {
			cxgo0.CurrentFileName = srcNames[i]
		}
		parseErrors += cxgo0.Parse(source)
		profiling.StopProfile(srcNames[i])
	}
	profiling.StopProfile("3. cxgo0")
	return parseErrors
}

func AddInitFunction(prgrm *cxcore.CXProgram) {
	mainPkg, err := prgrm.GetPackage(cxcore.MAIN_PKG)
	if err != nil {
		panic(err)
	}

	initFn := cxcore.MakeFunction(cxcore.SYS_INIT_FUNC, actions.CurrentFile, actions.LineNo)
	mainPkg.AddFunction(initFn)

	actions.FunctionDeclaration(initFn, nil, nil, actions.SysInitExprs)
	if _, err := mainPkg.SelectFunction(cxcore.MAIN_FUNC); err != nil {
		panic(err)
	}
}
