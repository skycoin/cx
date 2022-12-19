package fileloader

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/globals"
	"github.com/skycoin/cx/cx/packages"
	"github.com/skycoin/cx/cxparser/actions"
	"github.com/skycoin/cx/cxparser/util/profiling"
)

func LoadFiles(sourceCode []*os.File) (sourceCodeStrings [][]byte, fileNames []string, err error) {
	for _, source := range sourceCode {
		tmp := bytes.NewBuffer(nil)
		io.Copy(tmp, source)
		sourceCodeStrings = append(sourceCodeStrings, tmp.Bytes())
		fileNames = append(fileNames, source.Name())
	}

	reMultiCommentOpen := regexp.MustCompile(`/\*`)
	reMultiCommentClose := regexp.MustCompile(`\*/`)
	reComment := regexp.MustCompile("//")

	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)

	reImp := regexp.MustCompile("import")
	reImpName := regexp.MustCompile(`(^|[\s])import\s+"([_a-zA-Z][_a-zA-Z0-9/-]*)"`)

	profiling.StartProfile("1. packages")
	// 1. Identify all the packages and structs
	for srcI, srcStr := range sourceCodeStrings {
		srcName := fileNames[srcI]
		profiling.StartProfile(srcName)

		reader := bytes.NewReader(srcStr)
		scanner := bufio.NewScanner(reader)
		var commentedCode bool
		for scanner.Scan() {
			line := scanner.Bytes()

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
			// 1-a. Identify all the packages
			if loc := rePkg.FindIndex(line); loc != nil {
				if (commentLoc != nil && commentLoc[0] < loc[0]) ||
					(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
					(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
					// then it's commented out
					continue
				}

				if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
					if _, err := actions.AST.GetPackage(match[len(match)-1]); err != nil {
						// then it hasn't been added
						newPkg := ast.MakePackage(match[len(match)-1])
						pkgIdx := actions.AST.AddPackage(newPkg)
						newPkg, err = actions.AST.GetPackageFromArray(pkgIdx)
						if err != nil {
							panic(err)
						}
					}
				}
			}

		}
		profiling.StopProfile(srcName)
	} // for range sourceCodeStrings
	profiling.StopProfile("1. packages")

	profiling.StartProfile("2. imports")
	for i, source := range sourceCodeStrings {
		profiling.StartProfile(fileNames[i])

		var commentedCode bool

		scanner := bufio.NewScanner(bytes.NewReader(source))
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
					if _, err := actions.AST.GetPackage(pkgName); err != nil && !packages.IsDefaultPackage(pkgName) {
						// _, sourceCode, srcNames := ParseArgsForCX([]string{fmt.Sprintf("%s%s", SRCPATH, pkgName)}, false)
						_, sourceCode, _ := ast.ParseArgsForCX([]string{filepath.Join(globals.SRCPATH, pkgName)}, false)
						s, f, err := LoadFiles(sourceCode)
						if err != nil {
							return nil, nil, err
						}
						sourceCodeStrings = append(sourceCodeStrings, s...)
						fileNames = append(fileNames, f...)
					}
				}
			}

		}
		profiling.StopProfile(fileNames[i])
	}
	profiling.StopProfile("2. imports")

	return sourceCodeStrings, fileNames, nil
}
