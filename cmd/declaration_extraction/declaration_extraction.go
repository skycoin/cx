package declaration_extraction

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"strings"
	"unicode"
)

type GlobalDeclaration struct {
	PackageID   string
	FileID      string
	StartOffset int
	Length      int
	Name        string
}

type EnumDeclaration struct {
	PackageID   string
	FileID      string
	StartOffset int
	Length      int
	Name        string
	Type        string
	Value       int
}

type StructDeclaration struct {
	PackageID   string
	FileID      string
	StartOffset int
	Length      int
	Name        string
}

type FuncDeclaration struct {
	PackageID   string
	FileID      string
	StartOffset int
	Length      int
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

	var sourceWithoutBody []byte = source

	var inBlock int

	for i := range source {
		if source[i] == byte('{') {
			inBlock++
		}
		if source[i] == byte('}') {
			inBlock--
		}
		if inBlock != 0 {
			sourceWithoutBody[i] = byte(' ')
		}
	}

	GlblLocs := reGlbl.FindAllSubmatchIndex(sourceWithoutBody, -1)

	for i := range GlblLocs {
		var tmp GlobalDeclaration
		tmp.PackageID = pkg
		tmp.FileID = fileName
		tmp.StartOffset = GlblLocs[i][0]
		tmp.Length = GlblLocs[i][1] - GlblLocs[i][0]
		tmp.Name = string(source[GlblLocs[i][2]:GlblLocs[i][3]])
		GlblDec = append(GlblDec, tmp)
	}

	return GlblDec, err

}

func ExtractEnums(source []byte, fileName string, pkg string) ([]EnumDeclaration, error) {

	var EnumDec []EnumDeclaration
	var err error

	//Regexes
	reEnumInit := regexp.MustCompile(`const\s+\(`)
	reEnumDec := regexp.MustCompile(`([_a-zA-Z][_a-zA-Z0-9]*)\s+([_a-zA-Z][_a-zA-Z0-9]*)|([_a-zA-Z][_a-zA-Z0-9]*)`)

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)

	var EnumInit bool
	var inPrts int
	var Type string
	var Index int

	/* inPrts is supporting () with in the enum
	for example:

	const (
		a int = pointer(1) <---it would close off the enum here
	)

	*/

	for scanner.Scan() {
		line := scanner.Bytes()
		if locs := reEnumInit.FindAllIndex(line, -1); locs != nil {
			EnumInit = true
			inPrts++
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

		if match := reEnumDec.Find(line); match != nil && inPrts == 1 && EnumInit {
			reEnum := regexp.MustCompile(string(match))
			var tmp EnumDeclaration
			slice := strings.Split(string(match), " ")
			tmp.PackageID = pkg
			tmp.FileID = fileName
			tmp.StartOffset = reEnum.FindIndex(source)[1]
			tmp.Length = len(match)
			tmp.Name = string(slice[0])
			if len(slice) > 1 {
				Type = slice[1]
			}
			tmp.Type = Type
			tmp.Value = Index
			EnumDec = append(EnumDec, tmp)
			Index++
		}
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

	reFunc := regexp.MustCompile(`func\s+([_a-zA-Z][_a-zA-Z0-9]*)\s+\(.*\)`)

	funcLocs := reFunc.FindAllSubmatchIndex(source, -1)

	for i := range funcLocs {
		var tmp FuncDeclaration
		tmp.PackageID = pkg
		tmp.FileID = fileName
		tmp.StartOffset = funcLocs[i][0]
		tmp.Length = funcLocs[i][1] - funcLocs[i][0]
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

func GetDeclaration(source []byte, offset int, length int) []byte {
	return source[offset : offset+length]
}
