package declaration_extraction

import (
	"bufio"
	"bytes"
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

func extractGlbl(fileName string) ([]Declaration, error) {

	var GlblDec []Declaration
	var pkgID string
	src, err := os.ReadFile(fileName)

	CmtRmd := rmComment(src)
	pkgID = extractPkg(CmtRmd)

	//Regexs
	reGlbl := regexp.MustCompile(`var\s([_a-zA-Z][_a-zA-Z0-9]*)[ ]+([_a-zA-Z][_a-zA-Z0-9]*)`)
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

func extractStrct(fileName string) ([]Declaration, error) {
	var StrctDec []Declaration
	var pkgID string
	src, err := os.ReadFile(fileName)

	CmtRmd := rmComment(src)
	pkgID = extractPkg(CmtRmd)

	reStrctName := regexp.MustCompile(`type\s+([_a-zA-Z][_a-zA-Z0-9]*)[ ]+([_a-zA-Z][_a-zA-Z0-9]*)`)

	reader := bytes.NewReader(CmtRmd)
	scanner := bufio.NewScanner(reader)
	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()

		if match := reStrctName.Find(line); match != nil {
			var tmp Declaration
			tmp.PackageID = pkgID
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

func extractFunc(fileName string) ([]Declaration, error) {
	var FuncDec []Declaration
	var pkgID string
	src, err := os.ReadFile(fileName)

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
