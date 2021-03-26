package globals

import (
	"os"
	"path/filepath"
	"strings"
)

var CxProgramPath string = ""

func SetWorkingDir(filename string)  {
		filename = filepath.FromSlash(filename)
		i := strings.LastIndexByte(filename, os.PathSeparator)
		if i == -1 {
		i = 0
	}
		CxProgramPath = filename[:i]
}

func GetWorkDir(filename string) string {
	return CxProgramPath
}