package declaration_extraction

import (
	"bufio"
	"fmt"
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

func extractGlbl(source *os.File, fileName string) []Declaration {

	var GlblDec []Declaration
	var pkgID string

	scanner := bufio.NewScanner(source)

	//Regexs
	reMultiCommentOpen := regexp.MustCompile(`/\*`)
	reMultiCommentClose := regexp.MustCompile(`\*/`)
	reComment := regexp.MustCompile("//")

	reBodyOpen := regexp.MustCompile("{")
	reBodyClose := regexp.MustCompile("}")

	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)

	reGlbl := regexp.MustCompile("var")
	reGlblName := regexp.MustCompile(`(^|[\s])var\s([_a-zA-Z][_a-zA-Z0-9]*)`)

	var inBlock int
	var commentedCode bool

	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()

		commentLoc := reComment.FindIndex(line)

		multiCommentOpenLoc := reMultiCommentOpen.FindIndex(line)
		multiCommentCloseLoc := reMultiCommentClose.FindIndex(line)

		//Check for commented code
		if commentedCode && multiCommentCloseLoc != nil {
			commentedCode = false
		}

		if commentedCode {
			continue
		}

		if multiCommentOpenLoc != nil && !commentedCode && multiCommentCloseLoc != nil {
			commentedCode = true
		}

		//Check if in block of code
		if locs := reBodyOpen.FindAllIndex(line, -1); locs != nil {
			for _, loc := range locs {
				if !(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
					if (commentLoc == nil && multiCommentOpenLoc == nil && multiCommentCloseLoc == nil) ||
						(commentLoc != nil && commentLoc[0] > loc[0]) ||
						(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] > loc[0]) ||
						(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] < loc[0]) {
						inBlock++
					}
				}
			}
		}

		if locs := reBodyClose.FindAllIndex(line, -1); locs != nil {
			for _, loc := range locs {
				if !(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
					if (commentLoc == nil && multiCommentOpenLoc == nil && multiCommentCloseLoc == nil) ||
						(commentLoc != nil && commentLoc[0] > loc[0]) ||
						(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] > loc[0]) ||
						(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] < loc[0]) {
						inBlock--
					}
				}
			}
		}

		// we search for packages at the same time, so we can know to what package to add the global
		if loc := rePkg.FindIndex(line); loc != nil {
			if (commentLoc != nil && commentLoc[0] < loc[0]) ||
				(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
				(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
				// then it's commented out
				continue
			}

			if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
				tmp := strings.Split(match[0], "package")[1]
				pkgID = strings.TrimSpace(tmp)
			}
		}

		if loc := reGlbl.FindIndex(line); loc != nil {
			if commentLoc != nil && commentLoc[0] < loc[0] ||
				(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
				(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0] || inBlock != 0) {
				continue
			}
			if match := reGlblName.FindStringSubmatch(string(line)); match != nil {
				var tmp Declaration
				tmp.PackageID = pkgID
				tmp.FileID = fileName
				tmp.StartOffset = reGlblName.FindIndex(line)[0]
				tmp.Length = len(match)
				name := strings.Split(match[0], "var")[1]
				tmp.Name = strings.TrimSpace(name)
				GlblDec = append(GlblDec, tmp)
			}
		}
	}

	return GlblDec

}

func extractStrct(source *os.File, fileName string) []Declaration {
	var StrctDec []Declaration
	var pkgID string
	var commentedCode bool

	scanner := bufio.NewScanner(source)

	//Regexs
	reMultiCommentOpen := regexp.MustCompile(`/\*`)
	reMultiCommentClose := regexp.MustCompile(`\*/`)
	reComment := regexp.MustCompile("//")

	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)

	reStrct := regexp.MustCompile("type")
	reStrctName := regexp.MustCompile(`(^|[\s])type\s+([_a-zA-Z][_a-zA-Z0-9]*)?\s`)

	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()

		commentLoc := reComment.FindIndex(line)

		multiCommentOpenLoc := reMultiCommentOpen.FindIndex(line)
		multiCommentCloseLoc := reMultiCommentClose.FindIndex(line)

		//Check for commented code
		if commentedCode && multiCommentCloseLoc != nil {
			commentedCode = false
		}

		if commentedCode {
			continue
		}

		if multiCommentOpenLoc != nil && !commentedCode && multiCommentCloseLoc != nil {
			commentedCode = true
		}

		// we search for packages at the same time
		if loc := rePkg.FindIndex(line); loc != nil {
			if (commentLoc != nil && commentLoc[0] < loc[0]) ||
				(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
				(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
				// then it's commented out
				continue
			}

			if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
				tmp := strings.Split(match[0], "package")[1]
				pkgID = strings.TrimSpace(tmp)
			}
		}

		// 1-b. Identify all the structs
		if loc := reStrct.FindIndex(line); loc != nil {
			if (commentLoc != nil && commentLoc[0] < loc[0]) ||
				(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
				(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
				// then it's commented out
				continue
			}

			if match := reStrctName.FindStringSubmatch(string(line)); match != nil {
				var tmp Declaration
				tmp.PackageID = pkgID
				tmp.FileID = fileName
				tmp.StartOffset = reStrctName.FindIndex(line)[0]
				tmp.Length = len(match)
				name := strings.Split(match[0], "type")[1]
				tmp.Name = strings.TrimSpace(name)
				StrctDec = append(StrctDec, tmp)
			}
		}
	}
	return StrctDec
}

func extractFunc(source *os.File, fileName string) []Declaration {
	var FuncDec []Declaration
	var pkgID string
	var commentedCode bool

	scanner := bufio.NewScanner(source)

	//Regexs
	reMultiCommentOpen := regexp.MustCompile(`/\*`)
	reMultiCommentClose := regexp.MustCompile(`\*/`)
	reComment := regexp.MustCompile("//")

	rePkg := regexp.MustCompile("package")
	rePkgName := regexp.MustCompile(`(^|[\s])package\s+([_a-zA-Z][_a-zA-Z0-9]*)`)

	reFunc := regexp.MustCompile("func")
	reFuncName := regexp.MustCompile(`(^|[\s])func\s+([_a-zA-Z][_a-zA-Z0-9]*)?\s`)

	//Reading code line by line
	for scanner.Scan() {
		line := scanner.Bytes()

		commentLoc := reComment.FindIndex(line)

		multiCommentOpenLoc := reMultiCommentOpen.FindIndex(line)
		multiCommentCloseLoc := reMultiCommentClose.FindIndex(line)

		//Check for commented code
		if commentedCode && multiCommentCloseLoc != nil {
			commentedCode = false
		}

		if commentedCode {
			continue
		}

		if multiCommentOpenLoc != nil && !commentedCode && multiCommentCloseLoc != nil {
			commentedCode = true
		}

		// we search for packages at the same time
		if loc := rePkg.FindIndex(line); loc != nil {
			if (commentLoc != nil && commentLoc[0] < loc[0]) ||
				(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
				(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
				// then it's commented out
				continue
			}

			if match := rePkgName.FindStringSubmatch(string(line)); match != nil {
				tmp := strings.Split(match[0], "package")[1]
				pkgID = strings.TrimSpace(tmp)
			}
		}

		// 1-b. Identify all the structs
		if loc := reFunc.FindIndex(line); loc != nil {
			if (commentLoc != nil && commentLoc[0] < loc[0]) ||
				(multiCommentOpenLoc != nil && multiCommentOpenLoc[0] < loc[0]) ||
				(multiCommentCloseLoc != nil && multiCommentCloseLoc[0] > loc[0]) {
				// then it's commented out
				continue
			}

			if match := reFuncName.FindStringSubmatch(string(line)); match != nil {
				var tmp Declaration
				tmp.PackageID = pkgID
				tmp.FileID = fileName
				tmp.StartOffset = reFuncName.FindIndex(line)[0]
				tmp.Length = len(match)
				name := strings.Split(match[0], "func")[1]
				tmp.Name = strings.TrimSpace(name)
				FuncDec = append(FuncDec, tmp)

				fmt.Printf("%+v\n", name)
			}
		}
	}
	return FuncDec
}
