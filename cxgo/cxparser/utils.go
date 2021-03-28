package cxparser

import (
	"bufio"
	"path/filepath"
	"strings"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	globals2 "github.com/skycoin/cx/cx/globals"
	"github.com/skycoin/cx/cxgo/actions"
	"github.com/skycoin/cx/cxgo/cxgo0"
	"github.com/skycoin/cx/cxgo/globals"
	"github.com/skycoin/cx/cxgo/util/profiling"
)

func lexstruct(programSourceCode []string, programFileNames []string, cxpackage *ast.CXPackage, re RegularExpression) {

	//Identify structs

	for sourceIndex, sourcestring := range programSourceCode {

		srcName := programFileNames[sourceIndex]

		profiling.StartProfile(srcName)

		reader := strings.NewReader(sourcestring)

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

			// Identify all the structs
			if loc := re.reStruct.FindIndex(line); loc != nil {

				if (commentLoc != nil && commentLoc[0] < loc[0]) ||
					(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
					(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
					// then it's commented out
					continue
				}

				if match := re.reStructName.FindStringSubmatch(string(line)); match != nil {

					if cxpackage == nil {

						println(ast.CompilationError(srcName, lineno),
							"No package defined")

					} else if _, err := cxgo0.PRGRM0.GetStruct(match[len(match)-1], cxpackage.Name); err != nil {

						// then it hasn't been added
						strct := ast.MakeStruct(match[len(match)-1])
						cxpackage.AddStruct(strct)
					}
				}
			}
		}
		profiling.StopProfile(srcName)
	} // for range srcStrs
	profiling.StopProfile("1. packages/structs")

}

func lexpackages(programSourceCode []string, programFileNames []string, cxpackage *ast.CXPackage, re RegularExpression) {

	// 1. Identify all the packages and structs
	for sourceIndex, srcStr := range programSourceCode {

		srcName := programFileNames[sourceIndex]
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
						newPkg := ast.MakePackage(match[len(match)-1])
						cxgo0.PRGRM0.AddPackage(newPkg)
						cxpackage = newPkg
					} else {
						cxpackage = pkg
					}
				}
			}
		}
		profiling.StopProfile(srcName)
	} // for range srcStrs
	profiling.StopProfile("1. packages/structs")

}

func lexglobal(programSourceCode []string, programFileNames []string, cxpackage *ast.CXPackage, re RegularExpression) {

	profiling.StartProfile("2. globals")
	// 2. Identify all global variables
	//    We also identify packages again, so we know to what
	//    package we're going to add the variable declaration to.
	for sourceindex, sourcestring := range programSourceCode {
		profiling.StartProfile(programFileNames[sourceindex])
		// inBlock needs to be 0 to guarantee that we're in the global scope
		var inBlock int
		var commentedCode bool

		scanner := bufio.NewScanner(strings.NewReader(sourcestring))
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
					if _, err := cxgo0.PRGRM0.GetPackage(pkgName); err != nil && !constants.IsCorePackage(pkgName) {
						// _, sourceCode, srcNames := ParseArgsForCX([]string{fmt.Sprintf("%s%s", SRCPATH, pkgName)}, false)
						_, sourceCode, fileNames := ast.ParseArgsForCX([]string{filepath.Join(globals2.SRCPATH, pkgName)}, false)
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
						cxpackage = ast.MakePackage(match[len(match)-1])
						cxgo0.PRGRM0.AddPackage(cxpackage)
					} else {
						cxpackage = pkg
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
						cxpackage = pkg
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
					if _, err := cxpackage.GetGlobal(match[len(match)-1]); err != nil {
						// then it hasn't been added
						arg := ast.MakeArgument(match[len(match)-1], "", 0)
						arg.Offset = -1
						arg.Package = cxpackage
						cxpackage.AddGlobal(arg)
					}
				}
			}
		}
		profiling.StopProfile(programFileNames[sourceindex])
	}
	profiling.StopProfile("2. globals")
}

func AddInitFunction(program *ast.CXProgram) error {

	mainPackage, err := program.GetPackage(constants.MAIN_PKG)
	if err != nil {
		return err
	}

	initFn := ast.MakeFunction(constants.SYS_INIT_FUNC, actions.CurrentFile, actions.LineNo)
	mainPackage.AddFunction(initFn)

	//Init Expressions
	actions.FunctionDeclaration(initFn, nil, nil, globals.SysInitExprs)

	if _, err := mainPackage.SelectFunction(constants.MAIN_FUNC); err != nil {
		return err
	}
	return nil
}
