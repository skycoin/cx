package declaration_extraction

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type GlobalDeclaration struct {
	PackageID   string
	FileID      string
	StartOffset int
	Length      int
	LineNumber  int
	Name        string
}

type EnumDeclaration struct {
	PackageID   string
	FileID      string
	StartOffset int
	Length      int
	LineNumber  int
	Type        string
	Value       int
	Name        string
}

type StructDeclaration struct {
	PackageID   string
	FileID      string
	StartOffset int
	Length      int
	LineNumber  int
	Name        string
}

type FuncDeclaration struct {
	PackageID   string
	FileID      string
	StartOffset int
	Length      int
	LineNumber  int
	Name        string
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
			tmp.Name = string(line[match[2]+bytes : match[3]+bytes])
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
	reEnumDec := regexp.MustCompile(`[_a-zA-Z][_a-zA-Z0-9]*`)
	reEqual := regexp.MustCompile(`=`)

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)

	var EnumInit bool
	var inPrts int
	var Type string
	var Index int
	var bytes int
	var lineno int

	/* inPrts is supporting () with in the enum
	for example:

	const (
		a int = pointer(1) <---it would close off the enum here
	)

	*/

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

		if match := reEnumDec.FindAllIndex(line, -1); match != nil && inPrts == 1 && EnumInit {
			var tmp EnumDeclaration
			tmp.PackageID = pkg
			tmp.FileID = fileName
			tmp.StartOffset = match[0][0] + bytes
			tmp.Length = match[0][1] - match[0][0]

			if len(match) > 1 {
				tmp.Length = match[1][1] - match[0][0]
			}

			tmp.LineNumber = lineno
			fmt.Print(match)
			tmp.Name = string(source[match[0][0]+bytes : match[0][1]+bytes])

			if len(match) > 1 {
				Type = string(source[match[1][0]+bytes : match[1][1]+bytes])
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

	structLocs := reStruct.FindAllSubmatchIndex(source, -1)

	for i := range structLocs {
		var tmp StructDeclaration
		tmp.PackageID = pkg
		tmp.FileID = fileName
		tmp.StartOffset = structLocs[i][0]
		tmp.Length = structLocs[i][1] - structLocs[i][0]
		tmp.Name = string(source[structLocs[i][2]:structLocs[i][3]])
		StrctDec = append(StrctDec, tmp)
	}

	return StrctDec, err
}

func ExtractFuncs(source []byte, fileName string, pkg string) ([]FuncDeclaration, error) {

	var FuncDec []FuncDeclaration
	var err error

	reFunc := regexp.MustCompile(`func\s+([_a-zA-Z][_a-zA-Z0-9]*).*{`)

	funcLocs := reFunc.FindAllSubmatchIndex(source, -1)

	for i := range funcLocs {
		var tmp FuncDeclaration

		funcDec := string(bytes.Trim(source[funcLocs[i][0]:funcLocs[i][1]], "{"))
		funcTrimSpace := bytes.TrimSpace([]byte(funcDec))
		tmp.PackageID = pkg
		tmp.FileID = fileName
		tmp.StartOffset = funcLocs[i][0]
		tmp.Length = len(funcTrimSpace)
		tmp.Name = string(source[funcLocs[i][2]:funcLocs[i][3]])
		FuncDec = append(FuncDec, tmp)
	}

	return FuncDec, err
}

func ReDeclarationCheck(Glbl []GlobalDeclaration, Enum []EnumDeclaration, Strct []StructDeclaration, Func []FuncDeclaration) error {
	var err error

	for i := 0; i < len(Glbl); i++ {
		for j := i + 1; j < len(Glbl); j++ {
			if Glbl[i].Name == Glbl[j].Name {
				err = errors.New("global redeclared")
				return err
			}
		}
	}

	for i := 0; i < len(Enum); i++ {
		for j := i + 1; j < len(Enum); j++ {
			if Enum[i].Name == Enum[j].Name {
				err = errors.New("enum redeclared")
				return err
			}
		}
	}

	for i := 0; i < len(Strct); i++ {
		for j := i + 1; j < len(Strct); j++ {
			if Strct[i].Name == Strct[j].Name {
				err = errors.New("struct redeclared")
				return err
			}
		}
	}

	for i := 0; i < len(Func); i++ {
		for j := i + 1; j < len(Func); j++ {
			if Func[i].Name == Func[j].Name {
				err = errors.New("func redeclared")
				return err
			}
		}
	}

	err = nil
	return err
}

func GetDeclaration(source []byte, Glbl []GlobalDeclaration, Enum []EnumDeclaration, Strct []StructDeclaration, Func []FuncDeclaration) []string {
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
