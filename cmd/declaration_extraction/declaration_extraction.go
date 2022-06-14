package declaration_extraction

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

// All units for offset/length are in counted in bytes

type GlobalDeclaration struct {
	PackageID          string // name of package
	FileID             string // name of file
	StartOffset        int    // offset with in the file starting from 'var' of file
	Length             int    // length of entire declaration i.e. 'var [name] [type]'
	LineNumber         int    // line number of declaration
	GlobalVariableName string // name of variable being declared
}

type EnumDeclaration struct {
	PackageID        string // name of package
	FileID           string // name of file
	StartOffset      int    // offset with in the file starting from '[name of enum]' of file
	Length           int    // length of entire declaration i.e. '[name] [type]' or '[name]'
	LineNumber       int    // line number of declaration
	Type             string // type of enum
	Value            int    // value of enum
	EnumVariableName string // name of variable being declared
}

type StructDeclaration struct {
	PackageID          string // name of package
	FileID             string // name of file
	StartOffset        int    // offset with in the file starting from 'type'
	Length             int    // length of entire declaration i.e. 'type [name] [type]'
	LineNumber         int    // line number of declaration
	StructVariableName string // name of variable being declared
}

type FuncDeclaration struct {
	PackageID        string // name of package
	FileID           string // name of file
	StartOffset      int    // offset with in the file starting from 'func'
	Length           int    // length of entire declaration i.e. 'func [name] ()' or 'func [name] ([params])' or 'func [name] ([params]) [returns]'
	LineNumber       int    // line number of declaration
	FuncVariableName string // name of variable being declared
}

func ReplaceCommentsWithWhitespaces(source []byte) []byte {

	var sourceWithoutComments []byte = source

	reComment := regexp.MustCompile(`//.*|/\*[\s\S]*?\*/`)

	comments := reComment.FindAllIndex(source, -1)

	for i := range comments {
		for loc := comments[i][0]; loc < comments[i][1]; loc++ {
			if unicode.IsSpace(rune(sourceWithoutComments[loc])) {
				continue
			}
			sourceWithoutComments[loc] = byte(' ')
		}
	}

	return sourceWithoutComments
}

func ExtractPackages(source []byte) string {
	rePkgName := regexp.MustCompile(`(^|[\s])package[ \t]+([_a-zA-Z][_a-zA-Z0-9]*)`)

	srcStr := rePkgName.FindString(string(source))
	if srcStr != "" {
		srcStr = strings.Split(srcStr, " ")[1]
	}

	return srcStr
}

func ExtractGlobals(source []byte, fileName string, pkg string) ([]GlobalDeclaration, error) {

	var GlblDec []GlobalDeclaration
	var err error

	//Regexs
	reGlbl := regexp.MustCompile(`var\s([_a-zA-Z][_a-zA-Z0-9]*)\s+[\[_a-zA-Z][\]_a-zA-Z0-9]*`)
	reBodyOpen := regexp.MustCompile("{")
	reBodyClose := regexp.MustCompile("}")

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)

	var bytes int
	var lineno int
	var inBlock int

	for scanner.Scan() {
		line := scanner.Bytes()
		lineno++

		if locs := reBodyOpen.FindAllIndex(line, -1); locs != nil {
			inBlock++
		}
		if locs := reBodyClose.FindAllIndex(line, -1); locs != nil {
			inBlock--
		}

		if match := reGlbl.FindSubmatchIndex(line); match != nil && inBlock == 0 {
			var tmp GlobalDeclaration
			tmp.PackageID = pkg
			tmp.FileID = fileName
			tmp.StartOffset = bytes + match[0]
			tmp.Length = match[1] - match[0]
			tmp.LineNumber = lineno
			tmp.GlobalVariableName = string(source[match[2]+bytes : match[3]+bytes])
			GlblDec = append(GlblDec, tmp)
		}
		bytes += len(line) + 2
	}
	return GlblDec, err

}

func ExtractEnums(source []byte, fileName string, pkg string) ([]EnumDeclaration, error) {

	var EnumDec []EnumDeclaration
	var err error

	//Regexes
	reEnumInit := regexp.MustCompile(`const\s+\(`)
	rePrtsOpen := regexp.MustCompile(`\(`)
	rePrtsClose := regexp.MustCompile(`\)`)
	reEnumDec := regexp.MustCompile(`([_a-zA-Z][_a-zA-Z0-9]*)\s+([_a-zA-Z][_a-zA-Z0-9]*)|([_a-zA-Z][_a-zA-Z0-9]*)`)

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)

	var EnumInit bool
	var inPrts int
	var Type string
	var Index int
	var bytes int
	var lineno int

	for scanner.Scan() {
		line := scanner.Bytes()
		lineno++

		if locs := reEnumInit.FindAllIndex(line, -1); locs != nil {
			EnumInit = true
			inPrts++
			bytes += len(line) + 2
			continue
		}

		if locs := rePrtsOpen.FindAllIndex(line, -1); locs != nil && EnumInit {
			inPrts++
		}
		if locs := rePrtsClose.FindAllIndex(line, -1); locs != nil && EnumInit {
			inPrts--
		}

		if inPrts == 0 {
			EnumInit = false
			Type = ""
			Index = 0
		}

		if match := reEnumDec.FindSubmatchIndex(line); match != nil && inPrts == 1 && EnumInit {

			var tmp EnumDeclaration

			tmp.PackageID = pkg
			tmp.FileID = fileName
			tmp.StartOffset = match[0] + bytes
			tmp.Length = match[1] - match[0]
			tmp.LineNumber = lineno

			tmp.EnumVariableName = string(source[match[6]+bytes : match[7]+bytes])

			if match[2] != -1 {
				Type = string(source[match[4]+bytes : match[5]+bytes])
				tmp.EnumVariableName = string(source[match[2]+bytes : match[3]+bytes])
			}

			tmp.Type = Type

			tmp.Value = Index
			EnumDec = append(EnumDec, tmp)
			Index++
		}

		bytes += len(line) + 2
	}

	return EnumDec, err

}

func ExtractStructs(source []byte, fileName string, pkg string) ([]StructDeclaration, error) {

	var StrctDec []StructDeclaration
	var err error

	reStruct := regexp.MustCompile(`type\s+([_a-zA-Z][_a-zA-Z0-9]*)\s+[_a-zA-Z][_a-zA-Z0-9]*`)

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)

	var bytes int
	var lineno int

	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()
		lineno++

		if match := reStruct.FindSubmatchIndex(line); match != nil {

			var tmp StructDeclaration

			tmp.PackageID = pkg
			tmp.FileID = fileName

			tmp.StartOffset = match[0] + bytes
			tmp.Length = match[1] - match[0]
			tmp.StructVariableName = string(source[match[2]+bytes : match[3]+bytes])

			tmp.LineNumber = lineno

			StrctDec = append(StrctDec, tmp)

		}
		bytes += len(line) + 2
	}

	return StrctDec, err
}

func ExtractFuncs(source []byte, fileName string, pkg string) ([]FuncDeclaration, error) {

	var FuncDec []FuncDeclaration
	var err error

	reFunc := regexp.MustCompile(`func\s+([_a-zA-Z]\w*)\s+\(.*\)\s+\S+\w+|func\s+([_a-zA-Z]\w*)\s+\(.*\)`)

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)

	var bytes int
	var lineno int

	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()
		lineno++

		if match := reFunc.FindSubmatchIndex(line); match != nil {

			var tmp FuncDeclaration

			tmp.PackageID = pkg
			tmp.FileID = fileName

			tmp.StartOffset = match[0] + bytes
			tmp.Length = match[1] - match[0]

			tmp.FuncVariableName = string(source[match[4]+bytes : match[5]+bytes])

			if match[2] != -1 {
				tmp.FuncVariableName = string(source[match[2]+bytes : match[3]+bytes])
			}

			tmp.LineNumber = lineno

			FuncDec = append(FuncDec, tmp)

		}
		bytes += len(line) + 2
	}

	return FuncDec, err
}

func ReDeclarationCheck(Glbl []GlobalDeclaration, Enum []EnumDeclaration, Strct []StructDeclaration, Func []FuncDeclaration) error {
	var err error

	for i := 0; i < len(Glbl); i++ {
		for j := i + 1; j < len(Glbl); j++ {
			if Glbl[i].GlobalVariableName == Glbl[j].GlobalVariableName {
				err = errors.New("global redeclared")
				return err
			}
		}
	}

	for i := 0; i < len(Enum); i++ {
		for j := i + 1; j < len(Enum); j++ {
			if Enum[i].EnumVariableName == Enum[j].EnumVariableName {
				err = errors.New("enum redeclared")
				return err
			}
		}
	}

	for i := 0; i < len(Strct); i++ {
		for j := i + 1; j < len(Strct); j++ {
			if Strct[i].StructVariableName == Strct[j].StructVariableName {
				err = errors.New("struct redeclared")
				return err
			}
		}
	}

	for i := 0; i < len(Func); i++ {
		for j := i + 1; j < len(Func); j++ {
			if Func[i].FuncVariableName == Func[j].FuncVariableName {
				err = errors.New("func redeclared")
				return err
			}
		}
	}

	err = nil
	return err
}

func GetDeclarations(source []byte, Glbl []GlobalDeclaration, Enum []EnumDeclaration, Strct []StructDeclaration, Func []FuncDeclaration) []string {
	var declarations []string

	for i := 0; i < len(Glbl); i++ {
		declarations = append(declarations, string(source[Glbl[i].StartOffset:Glbl[i].StartOffset+Glbl[i].Length]))
	}

	for i := 0; i < len(Enum); i++ {
		declarations = append(declarations, string(source[Enum[i].StartOffset:Enum[i].StartOffset+Enum[i].Length]))
	}

	for i := 0; i < len(Strct); i++ {
		declarations = append(declarations, string(source[Strct[i].StartOffset:Strct[i].StartOffset+Strct[i].Length]))
	}

	for i := 0; i < len(Func); i++ {
		declarations = append(declarations, string(source[Func[i].StartOffset:Func[i].StartOffset+Func[i].Length]))
	}

	return declarations
}

func ExtractAllDeclarations(source []*os.File) ([]GlobalDeclaration, []EnumDeclaration, []StructDeclaration, []FuncDeclaration, error) {

	GlobalsCh := make(chan []GlobalDeclaration, len(source))
	EnumsCh := make(chan []EnumDeclaration, len(source))
	StructsCh := make(chan []StructDeclaration, len(source))
	FuncsCh := make(chan []FuncDeclaration, len(source))
	ErrCh := make(chan error, len(source))

	var wg sync.WaitGroup

	for i := range source {

		wg.Add(1)

		currentFile := source[i]

		go func(curFile *os.File, GlblsCh chan<- []GlobalDeclaration, EnmsCh chan<- []EnumDeclaration, StrctsCh chan<- []StructDeclaration, FncsCh chan<- []FuncDeclaration, ErCh chan<- error, WG *sync.WaitGroup) {

			src, err := io.ReadAll(curFile)
			fileName := curFile.Name()

			if err != nil {
				ErrCh <- err
			}

			replaceComments := ReplaceCommentsWithWhitespaces(src)
			pkg := ExtractPackages(replaceComments)

			var extractWg sync.WaitGroup

			extractWg.Add(4)

			go func(ch chan<- []GlobalDeclaration, errch chan<- error, rplcComments []byte, filename string, pckg string, ewg *sync.WaitGroup) {

				defer ewg.Done() // <====== Doesn't Work

				globals, err := ExtractGlobals(rplcComments, filename, pckg)

				if err != nil {
					errch <- err
				}

				ch <- globals // <====== Works

			}(GlblsCh, ErCh, replaceComments, fileName, pkg, &extractWg)

			go func(ch chan<- []EnumDeclaration, errch chan<- error, rplcComments []byte, filename string, pckg string, ewg *sync.WaitGroup) {

				defer ewg.Done() // <====== Doesn't Work

				enums, err := ExtractEnums(rplcComments, filename, pckg)

				if err != nil {
					errch <- err
				}

				ch <- enums // <====== Works

			}(EnmsCh, ErCh, replaceComments, fileName, pkg, &extractWg)

			go func(ch chan<- []StructDeclaration, errch chan<- error, rplcComments []byte, filename string, pckg string, ewg *sync.WaitGroup) {

				defer ewg.Done() // <====== Doesn't Work

				structs, err := ExtractStructs(rplcComments, filename, pckg)

				if err != nil {
					errch <- err
				}

				ch <- structs // <====== Works

			}(StrctsCh, ErCh, replaceComments, fileName, pkg, &extractWg)

			go func(ch chan<- []FuncDeclaration, errch chan<- error, rplcComments []byte, filename string, pckg string, ewg *sync.WaitGroup) {

				defer ewg.Done() // <====== Doesn't Work

				funcs, err := ExtractFuncs(rplcComments, filename, pckg)

				if err != nil {
					errch <- err
				}

				ch <- funcs // <====== Works

			}(FncsCh, ErCh, replaceComments, fileName, pkg, &extractWg)

			extractWg.Wait()

			WG.Done() // <====== Works here

		}(currentFile, GlobalsCh, EnumsCh, StructsCh, FuncsCh, ErrCh, &wg)

	}

	// ISSUES AND FINDINGS
	//---------------------

	// 1. Waitgroup.done in nested goroutine working causing
	// the outer waitgroup to be waiting infinitely

	// 2. channels are recieving values in the goroutines

	// 3. tried passing values as variables and pointers

	// 4. tried renaming the passed variables

	// 5. tried placing waitgroup.done in different positions
	// with and without defer

	// 6. tried using errgroup.Group which by the way would be better
	// error handling since you can return error instead of
	// sending to channels

	// 7. tried testing each goroutine to see which one isn't excuting
	// all are writing to channel

	// 8. There is always one inner goroutine that doesn't wg.Done

	wg.Wait()

	close(GlobalsCh)
	close(EnumsCh)
	close(StructsCh)
	close(FuncsCh)

	Globals := <-GlobalsCh
	Enums := <-EnumsCh
	Structs := <-StructsCh
	Funcs := <-FuncsCh

	if err := ReDeclarationCheck(Globals, Enums, Structs, Funcs); err != nil {
		ErrCh <- err
	}

	err := <-ErrCh

	return Globals, Enums, Structs, Funcs, err
}
