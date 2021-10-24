package cxparsering

import (
	"bufio"
	"bytes"
	"io"
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

// PackageMeta contains the meta data for a package
type PackageMeta struct {
	packageName     string
	InputsProcessed bool //starts false
	Parsed          bool //starts false
	Parents         []*PackageMeta
	Imports         []*PackageMeta
	FileNames       []string
	SourceCodes     []string
}

// ParseList contains a list of packages parsed
var ParseList []*PackageMeta

// ParseIndex is the index of the package being parsed
var ParseIndex int = 0
var packageIndex = 0

// Parse parses a program
// example Parse([`package main ...`])

///two fiunctions
// 1. is load files, in package, to string
// 2. extract imports from  package; return list of packages in files/stgrings
// 3. parse package step

// ExtractImportedPackages takes in files, returns packages imported
func ExtractImportedPackages(pkgMeta PackageMeta) {
	importedPkgs := GetPackageImports(pkgMeta.packageName)
	for _, importedPkgName := range importedPkgs {
		for _, parsedPkg := range ParseList {
			if parsedPkg.packageName == importedPkgName {
				// package already exists
				// check if the existing pkg has imported this package
				// if not, add it
				if packageExists(pkgMeta.packageName, parsedPkg.Imports) != nil {
					panic("import cycle")
				}
				// then no need to create a new one
				pkgMeta.Imports = append(pkgMeta.Imports, parsedPkg)
			}
		}

		// parse files in package
		_, sourceCode, fileNames := ast.ParseArgsForCX([]string{filepath.Join(globals2.SRCPATH, importedPkgName)}, false)
		sourceCodeStrings := make([]string, len(sourceCode))
		for i, source := range sourceCode {
			tmp := bytes.NewBuffer(nil)
			io.Copy(tmp, source)
			sourceCodeStrings[i] = tmp.String()
		}

		importedPkgMeta := PackageMeta{
			packageName: importedPkgName,
			FileNames:   fileNames,
			SourceCodes: sourceCodeStrings,
		}

		pkgMeta.Imports = append(pkgMeta.Imports, &importedPkgMeta)
		ParseList = append(ParseList, &importedPkgMeta)
		PreParse(sourceCodeStrings, fileNames)
	}
}

func packageExists(pkgName string, pkgList []*PackageMeta) *PackageMeta {
	for _, pkg := range pkgList {
		if pkg.packageName == pkgName {
			return pkg
		}
	}
	return nil
}

// PreParse extracts all imports/dependency tree and load the files in to memory
func PreParse(srcStrs, srcNames []string) {
	for _, sourceCode := range srcStrs {
		pkgName := GetPackageName(sourceCode)
		if pkgFound := packageExists(pkgName, ParseList); pkgFound != nil {
			ExtractImportedPackages(*pkgFound)
		}
		pkgMeta := PackageMeta{
			packageName: pkgName,
			FileNames:   srcNames,
			SourceCodes: srcStrs,
		}
		ParseList = append(ParseList, &pkgMeta)
		ExtractImportedPackages(pkgMeta)
	}
}

// GetPackageName returns the package with the given name
func GetPackageName(sourceCode string) string {
	// for parsing comments
	reMultiCommentOpen := regexp.MustCompile(`/\*`)
	reMultiCommentClose := regexp.MustCompile(`\*/`)
	reComment := regexp.MustCompile("//")

	// for parsing packages
	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)

	reader := strings.NewReader(sourceCode)
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
				pkgName := match[len(match)-1]
				return pkgName
			}
		}
	}
	return ""
}

// GetPackageImports returns imported packages in a given source code
func GetPackageImports(sourceCode string) []string {
	var importedPkgs []string
	reMultiCommentOpen := regexp.MustCompile(`/\*`)
	reMultiCommentClose := regexp.MustCompile(`\*/`)
	reComment := regexp.MustCompile("//")

	reImp := regexp.MustCompile("import")
	reImpName := regexp.MustCompile(`(^|[\s])import\s+"([_a-zA-Z][_a-zA-Z0-9/-]*)"`)

	var commentedCode bool

	scanner := bufio.NewScanner(strings.NewReader(sourceCode))
	for scanner.Scan() {
		line := scanner.Bytes()

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
				importedPkgs = append(importedPkgs, pkgName)
			}
		}

	}
	return importedPkgs
}

// preliminarystage performs a first pass for the CX cxgo. Globals, packages and
// custom types are added to `cxpartialparsing.Program`.
func Preliminarystage(srcStrs, srcNames []string) int {
	var prePkg *ast.CXPackage
	parseErrors := 0

	reMultiCommentOpen := regexp.MustCompile(`/\*`)
	reMultiCommentClose := regexp.MustCompile(`\*/`)
	reComment := regexp.MustCompile("//")

	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)
	reStrct := regexp.MustCompile("type")
	reStrctName := regexp.MustCompile(`(^|[\s])type\s+([_a-zA-Z][_a-zA-Z0-9]*)?\s`)

	reGlbl := regexp.MustCompile("var")
	reGlblName := regexp.MustCompile(`(^|[\s])var\s([_a-zA-Z][_a-zA-Z0-9]*)`)

	reBodyOpen := regexp.MustCompile("{")
	reBodyClose := regexp.MustCompile("}")

	reImp := regexp.MustCompile("import")
	reImpName := regexp.MustCompile(`(^|[\s])import\s+"([_a-zA-Z][_a-zA-Z0-9/-]*)"`)

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
					if pkg, err := cxpartialparsing.Program.GetPackage(match[len(match)-1]); err != nil {
						// then it hasn't been added
						newPkg := ast.MakePackage(match[len(match)-1])
						cxpartialparsing.Program.AddPackage(newPkg)
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
						println(ast.CompilationError(srcName, lineno),
							"No package defined")
					} else if _, err := cxpartialparsing.Program.GetStruct(match[len(match)-1], prePkg.Name); err != nil {
						// then it hasn't been added
						strct := ast.MakeStruct(match[len(match)-1])
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
					if _, err := cxpartialparsing.Program.GetPackage(pkgName); err != nil && !constants2.IsCorePackage(pkgName) {
						// _, sourceCode, srcNames := ParseArgsForCX([]string{fmt.Sprintf("%s%s", SRCPATH, pkgName)}, false)
						_, sourceCode, fileNames := ast.ParseArgsForCX([]string{filepath.Join(globals2.SRCPATH, pkgName)}, false)
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
					if pkg, err := cxpartialparsing.Program.GetPackage(match[len(match)-1]); err != nil {
						// then it hasn't been added
						prePkg = ast.MakePackage(match[len(match)-1])
						cxpartialparsing.Program.AddPackage(prePkg)
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
		profiling.StopProfile(srcNames[i])
	}
	profiling.StopProfile("2. globals")

	profiling.StartProfile("3. cxpartialparsing")

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

	profiling.StopProfile("3. cxpartialparsing")
	return parseErrors
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
