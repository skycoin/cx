package stage1

import (
	"bufio"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var mainPackageDir string

type SourceFile struct {
	FileName   string
	SourceCode []byte
}

type Package struct {
	Name        string
	PackageDir  string
	SourceFiles []SourceFile
	// PackageImports
}

// ParsePackages parses and returns packages from all import paths
func ParsePackages(sourceCode []*os.File, sourceFiles []string) ([]*Package, error) {
	var packages []*Package = []*Package{}
	// read source codes from copy of sourceCodes
	sourceCodeStrings := make([]string, len(sourceCode))
	for i, source := range sourceCode {
		tmp := bytes.NewBuffer(nil)
		io.Copy(tmp, source)
		sourceCodeStrings[i] = tmp.String()
	}
	if len(sourceCodeStrings) == 0 {
		return nil, nil
	}

	// comment parser
	reMultiCommentOpen := regexp.MustCompile(`/\*`)
	reMultiCommentClose := regexp.MustCompile(`\*/`)
	reComment := regexp.MustCompile("//")

	// package name parser
	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)

	// package import parser
	reImp := regexp.MustCompile("import")
	reImpName := regexp.MustCompile(`(^|[\s])import\s+"([_a-zA-Z][_a-zA-Z0-9/-]*)"`)

	// parse package imports dirs
	for srcI, srcStr := range sourceCodeStrings {
		srcName := sourceFiles[srcI]

		reader := strings.NewReader(srcStr)
		scanner := bufio.NewScanner(reader)

		var commentedCode bool
		var lineno = 0
		// var currentPackage *Package
		for scanner.Scan() {
			line := scanner.Bytes()
			lineno++

			// identify whether we are in a comment or not.
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

			// 1a. Identify all the packages
			if loc := rePkg.FindIndex(line); loc != nil {
				if (commentLoc != nil && commentLoc[0] < loc[0]) ||
					(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
					(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
					// then it's commented out
					continue
				}

				if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
					// check if package exists
					pkgName := match[len(match)-1]
					var exists bool
					for _, pkg := range packages {
						if pkg.Name == pkgName {
							exists = true
							// currentPackage = pkg
							break
						}
					}
					if !exists {
						pkg := Package{
							Name:       match[len(match)-1],
							PackageDir: srcName,
						}
						packages = append(packages, &pkg)
						if mainPackageDir == "" && pkgName == "main" {
							absPath, err := filepath.Abs(srcName)
							if err != nil {
								return nil, err
							} else {
								mainPackageDir = filepath.Dir(absPath)
							}
						}
						// currentPackage = &pkg
					}

				}
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
					pkgPath := filepath.Join(mainPackageDir, pkgName)
					cxFiles, err := FindCxFiles(pkgPath, ".cx")
					// panic(pkgPath)
					if err != nil {
						return nil, err
					}
					var pkgSourceCodes = []*os.File{}
					for _, cxFile := range cxFiles {
						pkgSourceCode, err := os.Open(cxFile)
						if err != nil {
							return nil, err
						}
						pkgSourceCodes = append(pkgSourceCodes, pkgSourceCode)
					}
					pkgPackages, err := ParsePackages(pkgSourceCodes, cxFiles)
					if err != nil {
						return nil, err
					}

					// TODO: check if it's not a standard library package.
					// TODO: read all files in a given package
					// Checking if `pkgName` already exists

					for _, pkgPackage := range pkgPackages {
						exists := false
						for _, pkg := range packages {
							if pkg.Name == pkgPackage.Name {
								exists = true
								break
							}
						}
						if !exists {

							packages = append(packages, pkgPackage)
						}
					}
				}
			}
		}
	}

	// call ParsePackage for files in package dirs

	//
	return packages, nil
}

func FindCxFiles(root, ext string) ([]string, error) {
	var files []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ext {
			files = append(files, s)
		}
		return nil
	})
	return files, nil
}
