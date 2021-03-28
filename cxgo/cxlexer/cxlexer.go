package cxlexer

import (
	"bufio"
	"bytes"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/globals"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/cx/cxgo/actions"
	"github.com/skycoin/cx/cxgo/cxgo"
	"github.com/skycoin/cx/cxgo/cxgo0"
	"github.com/skycoin/cx/cxgo/util/cxprof"
)

// re contains all regular expressions used for lexing
var re = struct {
	// comments
	comment           *regexp.Regexp
	multiCommentOpen  *regexp.Regexp
	multiCommentClose *regexp.Regexp

	// packages and structs
	pkg     *regexp.Regexp
	pkgName *regexp.Regexp
	str     *regexp.Regexp
	strName *regexp.Regexp

	// globals
	gl     *regexp.Regexp
	glName *regexp.Regexp

	// body open/close
	bodyOpen  *regexp.Regexp
	bodyClose *regexp.Regexp

	// imports
	imp     *regexp.Regexp
	impName *regexp.Regexp
}{
	comment:           regexp.MustCompile("//"),
	multiCommentOpen:  regexp.MustCompile(`/\*`),
	multiCommentClose: regexp.MustCompile(`\*/`),

	pkg:     regexp.MustCompile("package"),
	pkgName: regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`),
	str:     regexp.MustCompile("type"),
	strName: regexp.MustCompile(`(^|[\s])type\s+([_a-zA-Z][_a-zA-Z0-9]*)?\s`),

	gl:     regexp.MustCompile("var"),
	glName: regexp.MustCompile(`(^|[\s])var\s([_a-zA-Z][_a-zA-Z0-9]*)`),

	bodyOpen:  regexp.MustCompile("{"),
	bodyClose: regexp.MustCompile("}"),

	imp:     regexp.MustCompile("import"),
	impName: regexp.MustCompile(`(^|[\s])import\s+"([_a-zA-Z][_a-zA-Z0-9/-]*)"`),
}

// lg contains loggers
var lg = struct {
	l1 logrus.FieldLogger // packages and structs
	l2 logrus.FieldLogger // globals
	l3 logrus.FieldLogger // cxgo0
	l4 logrus.FieldLogger // parse
}{}

func SetLogger(log logrus.FieldLogger) {
	if log != nil {
		lg.l1 = log.WithField("i", 1).WithField("section", "packages/structs")
		lg.l2 = log.WithField("i", 2).WithField("section", "globals")
		lg.l3 = log.WithField("i", 3).WithField("section", "cxgo0")
		lg.l4 = log.WithField("i", 4).WithField("section", "parse")
	}
}

// Step0 performs a first pass for the CX cxgo. Globals, packages and
// custom types are added to `cxgo0.PRGRM0`.
func Step0(srcStrs, srcNames []string) int {
	var prePkg *ast.CXPackage
	parseErrs := 0

	// 1. Identify all the packages and structs
	func() {
		if lg.l1 != nil {
			_, stopL1 := cxprof.StartProfile(lg.l1)
			defer stopL1()
		}

		for srcI, src := range srcStrs {
			idPkgAndStructs(srcNames[srcI], strings.NewReader(src), &prePkg)
		}
	}()

	// 2. Identify all global variables
	//    We also identify packages again, so we know to what
	//    package we're going to add the variable declaration to.
	func() {
		if lg.l2 != nil {
			_, stopL2 := cxprof.StartProfile(lg.l2)
			defer stopL2()
		}

		for srcI, src := range srcStrs {
			idGlobVars(srcNames[srcI], strings.NewReader(src), &prePkg)
		}
	}()

	// 3. Parse sources into cxgo0
	func() {
		if lg.l3 != nil {
			_, stopL3 := cxprof.StartProfile(lg.l3)
			defer stopL3()
		}

		for srcI, src := range srcStrs {
			parseErrs += cxgo0Parse(srcNames[srcI], src)
		}
	}()

	return parseErrs
}

// idPkgAndStructs (1) identifies packages and structs contained within given file.
func idPkgAndStructs(filename string, r io.Reader, prePkg **ast.CXPackage) {
	if lg.l1 != nil {
		_, stopL1x := cxprof.StartProfile(lg.l1.WithField("src_file", filename))
		defer stopL1x()
	}

	var (
		s         = bufio.NewScanner(r)
		inComment = false
		lineN     = 0
	)

	for s.Scan() {
		line := s.Bytes()
		lineN++

		cl := makeCommentLocator(line)
		if skip := cl.skipLine(&inComment, true); skip {
			continue
		}

		// At this point we know that we are *not* in a comment

		// 1a. Identify all the packages
		if loc := re.pkg.FindIndex(line); loc != nil {
			if (cl.singleOpen != nil && cl.singleOpen[0] < loc[0]) ||
				(cl.multiOpen != nil && cl.multiOpen[0] < loc[0]) ||
				(cl.multiClose != nil && cl.multiClose[0] > loc[0]) {
				// then it's commented out
				continue
			}

			if match := re.pkgName.FindStringSubmatch(string(line)); match != nil {
				if pkg, err := cxgo0.PRGRM0.GetPackage(match[len(match)-1]); err != nil {
					// then it hasn't been added
					newPkg := ast.MakePackage(match[len(match)-1])
					cxgo0.PRGRM0.AddPackage(newPkg)
					*prePkg = newPkg
				} else {
					*prePkg = pkg
				}
			}
		}

		// 1b. Identify all the structs
		if loc := re.str.FindIndex(line); loc != nil {
			if (cl.singleOpen != nil && cl.singleOpen[0] < loc[0]) ||
				(cl.multiOpen != nil && cl.multiOpen[0] < loc[0]) ||
				(cl.multiClose != nil && cl.multiClose[0] > loc[0]) {
				// then it's commented out
				continue
			}

			if match := re.strName.FindStringSubmatch(string(line)); match != nil {
				if prePkg == nil {
					println(ast.CompilationError(filename, lineN), "No package defined")
				} else if _, err := cxgo0.PRGRM0.GetStruct(match[len(match)-1], (*prePkg).Name); err != nil {
					// then it hasn't been added
					strct := ast.MakeStruct(match[len(match)-1])
					(*prePkg).AddStruct(strct)
				}
			}
		}
	}
}

// idGlobVars (2) identifies global variables contained within given file.
// We also need to identify packages again, so we know which package to add
func idGlobVars(filename string, r io.Reader, prePkg **ast.CXPackage) {
	if lg.l2 != nil {
		_, stopL2x := cxprof.StartProfile(lg.l2.WithField("src_file", filename))
		defer stopL2x()
	}

	var (
		s         = bufio.NewScanner(r)
		inBlock   = 0
		inComment = false
	)

	for s.Scan() {
		line := s.Bytes()

		cl := makeCommentLocator(line)
		if skip := cl.skipLine(&inComment, false); skip {
			continue
		}

		// Identify all the package imports.
		if loc := re.imp.FindIndex(line); loc != nil {
			if cl.isLocationCommented(loc) {
				continue
			}

			if match := re.impName.FindStringSubmatch(string(line)); match != nil {
				pkgName := match[len(match)-1]
				// Checking if `pkgName` already exists and if it's not a standard library package.
				if _, err := cxgo0.PRGRM0.GetPackage(pkgName); err != nil && !constants.IsCorePackage(pkgName) {
					// _, sourceCode, srcNames := ParseArgsForCX([]string{fmt.Sprintf("%s%s", SRCPATH, pkgName)}, false)
					_, sourceCode, fileNames := ast.ParseArgsForCX([]string{filepath.Join(globals.SRCPATH, pkgName)}, false)
					ParseSourceCode(sourceCode, fileNames) // TODO @evanlinjin: Check return value.
				}
			}
		}

		// we search for packages at the same time, so we can know to what package to add the global
		if loc := re.pkg.FindIndex(line); loc != nil {
			if cl.isLocationCommented(loc) {
				continue
			}

			if match := re.pkgName.FindStringSubmatch(string(line)); match != nil {
				if pkg, err := cxgo0.PRGRM0.GetPackage(match[len(match)-1]); err != nil {
					// then it hasn't been added
					*prePkg = ast.MakePackage(match[len(match)-1])
					cxgo0.PRGRM0.AddPackage(*prePkg)
				} else {
					*prePkg = pkg
				}
			}
		}

		if locs := re.bodyOpen.FindAllIndex(line, -1); locs != nil {
			for _, loc := range locs {
				if !(cl.multiClose != nil && cl.multiClose[0] > loc[0]) {
					// then it's outside of a */, e.g. `*/ }`
					if (cl.singleOpen == nil && cl.multiOpen == nil && cl.multiClose == nil) ||
						(cl.singleOpen != nil && cl.singleOpen[0] > loc[0]) ||
						(cl.multiOpen != nil && cl.multiOpen[0] > loc[0]) ||
						(cl.multiClose != nil && cl.multiClose[0] < loc[0]) {
						// then we have an uncommented opening bracket
						inBlock++
					}
				}
			}
		}

		if locs := re.bodyClose.FindAllIndex(line, -1); locs != nil {
			for _, loc := range locs {
				if !(cl.multiClose != nil && cl.multiClose[0] > loc[0]) {
					// then it's outside of a */, e.g. `*/ {`
					if (cl.singleOpen == nil && cl.multiOpen == nil && cl.multiClose == nil) ||
						(cl.singleOpen != nil && cl.singleOpen[0] > loc[0]) ||
						(cl.multiOpen != nil && cl.multiOpen[0] > loc[0]) ||
						(cl.multiClose != nil && cl.multiClose[0] < loc[0]) {
						// then we have an uncommented closing bracket
						inBlock--
					}
				}
			}
		}

		// we could have this situation: {var local i32}
		// but we don't care about this, as the later passes will throw an error as it's invalid syntax

		if loc := re.pkg.FindIndex(line); loc != nil {
			if cl.isLocationCommented(loc) {
				continue
			}

			if match := re.pkgName.FindStringSubmatch(string(line)); match != nil {
				if pkg, err := cxgo0.PRGRM0.GetPackage(match[len(match)-1]); err != nil {
					// it should be already present
					panic(err)
				} else {
					*prePkg = pkg
				}
			}
		}

		// finally, if we read a "var" and we're in global scope, we add the global without any type
		// the type will be determined later on
		if loc := re.gl.FindIndex(line); loc != nil {
			if cl.isLocationCommented(loc) || inBlock != 0 {
				continue
			}

			if match := re.glName.FindStringSubmatch(string(line)); match != nil {
				if _, err := (*prePkg).GetGlobal(match[len(match)-1]); err != nil {
					// then it hasn't been added
					arg := ast.MakeArgument(match[len(match)-1], "", 0)
					arg.Offset = -1
					arg.Package = *prePkg
					(*prePkg).AddGlobal(arg)
				}
			}
		}
	}
}

// cxgo0Parse (3) parses the source file into cxgo0.
func cxgo0Parse(filename string, src string) int {
	if lg.l3 != nil {
		_, stopL3x := cxprof.StartProfile(lg.l3.WithField("src_file", filename))
		defer stopL3x()
	}

	if !strings.HasSuffix(src, "\n") {
		src += "\n"
	}

	cxgo0.CurrentFileName = filename
	return cxgo0.Parse(src)
}

// ParseSourceCode takes a group of files representing CX `sourceCode` and
// parses it into CX program structures for `AST`.
func ParseSourceCode(sourceCode []*os.File, fileNames []string) int {
	cxgo0.PRGRM0 = actions.AST

	// Copy the contents of the file pointers containing the CX source
	// code into sourceCodeCopy
	sourceCodeCopy := make([]string, len(sourceCode))
	for i, source := range sourceCode {
		tmp := bytes.NewBuffer(nil)
		io.Copy(tmp, source)
		sourceCodeCopy[i] = tmp.String()
	}

	// We need to traverse the elements by hierarchy first add all the
	// packages and structs at the same time then add globals, as these
	// can be of a custom type (and it could be imported) the signatures
	// of functions and methods are added in the cxgo0.y pass
	parseErrors := 0
	if len(sourceCode) > 0 {
		parseErrors = Step0(sourceCodeCopy, fileNames)
	}

	actions.AST.SetCurrentCxProgram()

	actions.AST = cxgo0.PRGRM0
	if globals.FoundCompileErrors || parseErrors > 0 {
		return constants.CX_COMPILATION_ERROR
	}

	// Adding global variables `OS_ARGS` to the `os` (operating system)
	// package.
	if osPkg, err := actions.AST.GetPackage(constants.OS_PKG); err == nil {
		if _, err := osPkg.GetGlobal(constants.OS_ARGS); err != nil {
			arg0 := ast.MakeArgument(constants.OS_ARGS, "", -1).AddType(constants.TypeNames[constants.TYPE_UNDEFINED])
			arg0.Package = osPkg

			arg1 := ast.MakeArgument(constants.OS_ARGS, "", -1).AddType(constants.TypeNames[constants.TYPE_STR])
			arg1 = actions.DeclarationSpecifiers(arg1, []int{0}, constants.DECL_BASIC)
			arg1 = actions.DeclarationSpecifiers(arg1, []int{0}, constants.DECL_SLICE)

			actions.DeclareGlobalInPackage(osPkg, arg0, arg1, nil, false)
		}
	}

	if lg.l4 != nil {
		_, stopL4 := cxprof.StartProfile(lg.l4)
		defer stopL4()
	}

	// The last pass of parsing that generates the actual output.
	for i, source := range sourceCodeCopy {
		// Because of an unknown reason, sometimes some CX programs
		// throw an error related to a premature EOF (particularly in Windows).
		// Adding a newline character solves this.
		source = source + "\n"
		actions.LineNo = 1
		b := bytes.NewBufferString(source)
		if len(fileNames) > 0 {
			actions.CurrentFile = fileNames[i]
		}

		func() {
			if lg.l4 != nil {
				_, stopL4x := cxprof.StartProfile(lg.l4.WithField("src_file", actions.CurrentFile))
				defer stopL4x()
			}
			parseErrors += cxgo.Parse(cxgo.NewLexer(b))
		}()
	}

	if globals.FoundCompileErrors || parseErrors > 0 {
		return constants.CX_COMPILATION_ERROR
	}

	return 0
}

type commentLocator struct {
	singleOpen []int // index of single-line comment
	multiOpen  []int // index of multi-line comment open
	multiClose []int // index of multi-line comment close
}

func makeCommentLocator(line []byte) commentLocator {
	return commentLocator{
		singleOpen: re.comment.FindIndex(line),
		multiOpen:  re.multiCommentOpen.FindIndex(line),
		multiClose: re.multiCommentClose.FindIndex(line),
	}
}

func (cl commentLocator) skipLine(isCommented *bool, skipWithoutMultiClose bool) (skip bool) {
	if *isCommented {
		// If no multi-line close is detected, this line is still commented.
		if cl.multiClose == nil {
			return true
		}
		// Multi-line comment closed.
		*isCommented = false
	}

	// Detect start of multi-line comment.
	if cl.multiOpen != nil && cl.multiClose == nil {
		*isCommented = true

		if skipWithoutMultiClose {
			return true
		}
	}

	// No skip.
	return false
}

func (cl commentLocator) isLocationCommented(loc []int) bool {
	return (cl.singleOpen != nil && cl.singleOpen[0] < loc[0]) ||
		(cl.multiOpen != nil && cl.multiOpen[0] < loc[0]) ||
		(cl.multiClose != nil && cl.multiClose[0] > loc[0])
}
