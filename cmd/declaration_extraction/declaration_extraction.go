package declaration_extraction

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type Declaration struct {
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

func rmComment(source []byte) []byte {

	var src []byte
	// Regexs
	reMultiComment := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	reComment := regexp.MustCompile(`//[.]*`)

	// Possible issue if there are infront for multiline comments the offset maybe a bit off
	// For Example:
	// |/*
	// |	Multi line comment
	// |
	// |	*/func main () { <--- offset should be 3 but when extracted is 0 since there are no more comments
	// |
	// |}
	// |

	//Replace contents between /* */ with ""
	src = reMultiComment.ReplaceAll(source, []byte(""))
	//Replace contents after // with ""
	src = reComment.ReplaceAll(src, []byte(""))

	return src
}

func extractPkg(source []byte) string {
	rePkgName := regexp.MustCompile(`(^|[\s])package[ \t]+([_a-zA-Z][_a-zA-Z0-9]*)`)

	srcStr := rePkgName.FindString(string(source))
	if srcStr != "" {
		srcStr = strings.Split(srcStr, " ")[1]
	}

	return srcStr
}

func extractGlbl(source *os.File) []Declaration {

	var GlblDec []Declaration
	var pkgID string
	src, err := io.ReadAll(source)

	if err != nil {
		fmt.Println("Error reading...")
	}

	CmtRmd := rmComment(src)
	pkgID = extractPkg(CmtRmd)

	//Regexs
	reGlbl := regexp.MustCompile(`var\s([_a-zA-Z][_a-zA-Z0-9]*)\s[_a-zA-Z][_a-zA-Z0-9]*`)
	reBodyOpen := regexp.MustCompile("{")
	reBodyClose := regexp.MustCompile("}")

	reader := bytes.NewReader(CmtRmd)
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
			var tmp Declaration
			tmp.PackageID = string(pkgID)
			tmp.FileID = source.Name()
			tmp.StartOffset = reGlbl.FindIndex(line)[0]
			name := strings.Split(string(match), " ")[1]
			tmp.Length = len(match)
			tmp.Name = name
			GlblDec = append(GlblDec, tmp)
		}
	}

	return GlblDec

}

func extractEnum(source *os.File) []EnumDeclaration {
	var EnumDec []EnumDeclaration
	var pkgID string
	src, err := io.ReadAll(source)

	if err != nil {
		fmt.Println("Error reading...")
	}

	CmtRmd := rmComment(src)
	pkgID = extractPkg(CmtRmd)

	//Regexes
	reEnumInit := regexp.MustCompile(`const\s+\(`)
	rePrtsOpen := regexp.MustCompile(`\(`)
	rePrtsClose := regexp.MustCompile(`\)`)
	reEnumDec := regexp.MustCompile(`([_a-zA-Z][_a-zA-Z0-9]*)[ \t]+([_a-zA-Z][_a-zA-Z0-9]*)|([_a-zA-Z][_a-zA-Z0-9]*)`)

	reader := bytes.NewReader(CmtRmd)
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
			tmp.PackageID = string(pkgID)
			tmp.FileID = source.Name()
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

			fmt.Printf("[%v]", slice)
			// fmt.Println(index)
		}
		fmt.Println(inPrts)
		fmt.Println(EnumInit)
	}

	return EnumDec

}

func extractStrct(source *os.File) []Declaration {
	var StrctDec []Declaration
	var pkgID string
	src, err := io.ReadAll(source)

	if err != nil {
		fmt.Println("Error reading...")
	}

	CmtRmd := rmComment(src)
	pkgID = extractPkg(CmtRmd)

	reStrctName := regexp.MustCompile(`type\s+([_a-zA-Z][_a-zA-Z0-9]*)\s[_a-zA-Z][_a-zA-Z0-9]*`)

	reader := bytes.NewReader(CmtRmd)
	scanner := bufio.NewScanner(reader)
	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()

		if match := reStrctName.Find(line); match != nil {
			var tmp Declaration
			tmp.PackageID = pkgID
			tmp.FileID = source.Name()
			tmp.StartOffset = reStrctName.FindIndex(line)[0]
			name := strings.Split(string(match), " ")[1]
			tmp.Length = len(match)
			tmp.Name = name
			StrctDec = append(StrctDec, tmp)
		}
	}
	return StrctDec
}

func extractFunc(source *os.File) []Declaration {
	var FuncDec []Declaration
	var pkgID string
	src, err := io.ReadAll(source)

	if err != nil {
		fmt.Println("Error reading...")
	}

	CmtRmd := rmComment(src)
	pkgID = extractPkg(CmtRmd)

	//the offset is for the whole declaration including parameters "func main(para1, para2)"?

	reFuncName := regexp.MustCompile(`func\s+([_a-zA-Z][_a-zA-Z0-9]*)\s+(\(.*\))`)

	reader := bytes.NewReader(CmtRmd)
	scanner := bufio.NewScanner(reader)
	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()

		if match := reFuncName.Find(line); match != nil {
			var tmp Declaration
			tmp.PackageID = pkgID
			tmp.FileID = source.Name()
			tmp.StartOffset = reFuncName.FindIndex(line)[0]
			name := strings.Split(string(match), " ")[1]
			tmp.Length = len(match)
			tmp.Name = name
			FuncDec = append(FuncDec, tmp)
		}
	}
	return FuncDec
}

func reDeclearCheck(Glbl []Declaration, Enum []EnumDeclaration, Strct []Declaration, Func []Declaration) []error {
	var err []error

	for i := range Glbl {
		for j := i; j < len(Glbl); j++ {
			if Glbl[i].Name == Glbl[j].Name {
				err = append(err, errors.New(`redecleared global`))
				break
			}
		}
	}

	for i := range Enum {
		for j := i; j < len(Enum); j++ {
			if Enum[i].Name == Enum[j].Name &&
				Enum[i].Type == Enum[j].Type &&
				Enum[i].Value == Enum[j].Value {
				err = append(err, errors.New(`redecleared enum`))
				break
			}
		}
	}

	for i := range Strct {
		for j := i; j < len(Strct); j++ {
			if Strct[i].Name == Strct[j].Name {
				err = append(err, errors.New(`redecleared struct`))
				break
			}
		}
	}

	for i := range Func {
		for j := i; j < len(Func); j++ {
			if Func[i].Name == Func[j].Name {
				err = append(err, errors.New(`redecleared function`))
				break
			}
		}
	}

	return err

}
