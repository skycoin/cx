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

	// Possible issue if there are infront for multiline comments the offset maybe a bit off
	// For Example:
	// |/*
	// |	Multi line comment
	// |
	// |	*/func main () { <--- offset should be 3 but when extracted is 0 since there are no more comments
	// |
	// |}
	// |

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
	reGlbl := regexp.MustCompile(`var\s([_a-zA-Z][_a-zA-Z0-9]*)\s[_a-zA-Z][_a-zA-Z0-9]*`)
	reBodyOpen := regexp.MustCompile("{")
	reBodyClose := regexp.MustCompile("}")

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)

	var inBlock int

	for scanner.Scan() {
		line := scanner.Bytes()

		if locs := reBodyOpen.FindAllIndex(line, -1); locs != nil {
			inBlock++
		}
		if locs := reBodyClose.FindAllIndex(line, -1); locs != nil {
			inBlock--
		}

		if match := reGlbl.Find(line); match != nil && inBlock == 0 {
			var tmp GlobalDeclaration
			tmp.PackageID = pkg
			tmp.FileID = fileName
			tmp.StartOffset = reGlbl.FindIndex(line)[0]
			name := strings.Split(string(match), " ")[1]
			tmp.Length = len(match)
			tmp.Name = name
			GlblDec = append(GlblDec, tmp)
		}
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
			var tmp EnumDeclaration
			slice := strings.Split(string(match), " ")
			tmp.PackageID = pkg
			tmp.FileID = fileName
			tmp.StartOffset = reEnumDec.FindIndex(line)[0]
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

	reStrctName := regexp.MustCompile(`type\s+([_a-zA-Z][_a-zA-Z0-9]*)\s+[_a-zA-Z][_a-zA-Z0-9]*`)

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)
	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()

		if match := reStrctName.Find(line); match != nil {
			var tmp StructDeclaration
			tmp.PackageID = pkg
			tmp.FileID = fileName
			tmp.StartOffset = reStrctName.FindIndex(line)[0]
			name := strings.Split(string(match), " ")[1]
			tmp.Length = len(match)
			tmp.Name = name
			StrctDec = append(StrctDec, tmp)
		}
	}
	return StrctDec, err
}

func ExtractFuncs(source []byte, fileName string, pkg string) ([]FuncDeclaration, error) {

	var FuncDec []FuncDeclaration
	var err error

	//the offset is for the whole declaration including parameters "func main(para1, para2)"?

	reFuncName := regexp.MustCompile(`func\s+([_a-zA-Z][_a-zA-Z0-9]*)\s+(\(.*\))`)

	reader := bytes.NewReader(source)
	scanner := bufio.NewScanner(reader)
	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()

		if match := reFuncName.Find(line); match != nil {
			var tmp FuncDeclaration
			tmp.PackageID = pkg
			tmp.FileID = fileName
			tmp.StartOffset = reFuncName.FindIndex(line)[0]
			name := strings.Split(string(match), " ")[1]
			tmp.Length = len(match)
			tmp.Name = name
			FuncDec = append(FuncDec, tmp)
		}
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
