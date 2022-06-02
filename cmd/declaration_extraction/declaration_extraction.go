package declaration_extraction

import (
	"bufio"
	"bytes"
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

// func isDeclared(decrArr []Declaration, dec Declaration) bool {
// 	for _, curDec := range decrArr {
// 		if curDec == dec {
// 			return true
// 		}
// 	}
// 	return false
// }

func rmComment(source []byte) []byte {

	var src []byte
	// Regexs
	reMultiComment := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	reComment := regexp.MustCompile(`//[^\n\r]*`)

	//Replace contents between /* */ with ""
	src = reMultiComment.ReplaceAll(source, []byte(""))
	//Replace contents after // with ""
	src = reComment.ReplaceAll(src, []byte(""))

	return src
}

func extractPkg(source []byte) string {
	rePkgName := regexp.MustCompile(`(^|[\s])package[ ]+([_a-zA-Z][_a-zA-Z0-9]*)`)

	srcStr := rePkgName.FindString(string(source))
	if srcStr != "" {
		srcStr = strings.Split(srcStr, " ")[1]
	}

	return srcStr
}

func extractGlbl(source *os.File) ([]Declaration, error) {

	var GlblDec []Declaration
	var pkgID string
	src, err := io.ReadAll(source)

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

	return GlblDec, err

}

func extractEnum(source *os.File) ([]EnumDeclaration, error) {
	var EnumDec []EnumDeclaration
	var pkgID string
	src, err := io.ReadAll(source)

	CmtRmd := rmComment(src)
	pkgID = extractPkg(CmtRmd)

	//Regexes
	reEnumInit := regexp.MustCompile(`const`)
	rePrtsOpen := regexp.MustCompile(`\(`)
	rePrtsClose := regexp.MustCompile(`\)`)
	reEnumDec := regexp.MustCompile(`([_a-zA-Z][_a-zA-Z0-9]*)`)
	// reEqual := regexp.MustCompile(`=`)

	reader := bytes.NewReader(CmtRmd)
	scanner := bufio.NewScanner(reader)

	var EnumInit bool
	var inPrts int
	var Type string
	var Index int
	var lineNo int

	for scanner.Scan() {
		line := scanner.Bytes()
		if locs := reEnumInit.FindAllIndex(line, -1); locs != nil {
			EnumInit = true
		}

		if locs := rePrtsOpen.FindAllIndex(line, -1); locs != nil && EnumInit {
			inPrts++
		}
		if locs := rePrtsClose.FindAllIndex(line, -1); locs != nil && EnumInit {
			inPrts--
		}

		if EnumInit && inPrts > 0 {
			lineNo++
		}

		if inPrts == 0 {
			EnumInit = false
			Type = ""
			Index = 0
			lineNo = 0
		}

		if match := reEnumDec.FindAll(line, -1); match != nil && inPrts == 1 && EnumInit && lineNo > 1 {
			var tmp EnumDeclaration
			index := reEnumDec.FindIndex(line)
			tmp.PackageID = string(pkgID)
			tmp.FileID = source.Name()
			tmp.StartOffset = index[1]
			tmp.Length = len(match)
			tmp.Name = string(match[0][0])
			tmp.Type = Type
			tmp.Value = Index
			EnumDec = append(EnumDec, tmp)
			Index++

			fmt.Printf("[%v]", string(match[0][1]))
			fmt.Println(index)
		}
		fmt.Println(inPrts)
		fmt.Println(EnumInit)
	}

	return EnumDec, err

}

func extractStrct(source *os.File) ([]Declaration, error) {
	var StrctDec []Declaration
	var pkgID string
	src, err := io.ReadAll(source)

	CmtRmd := rmComment(src)
	pkgID = extractPkg(CmtRmd)

	reStrctName := regexp.MustCompile(`[\s]*type\s+([_a-zA-Z][_a-zA-Z0-9]*)\s[_a-zA-Z][_a-zA-Z0-9]*`)

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
	return StrctDec, err
}

func extractFunc(source *os.File) ([]Declaration, error) {
	var FuncDec []Declaration
	var pkgID string
	src, err := io.ReadAll(source)

	CmtRmd := rmComment(src)
	pkgID = extractPkg(CmtRmd)

	reFuncName := regexp.MustCompile(`func\s+([_a-zA-Z][_a-zA-Z0-9]*)?`)

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
	return FuncDec, err
}
