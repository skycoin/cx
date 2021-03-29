package globals

import (
	"github.com/skycoin/cx/cx/ast"
)

/*
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
 */

//What does this do?
//This is where intializers get pushed, but only used 4 times
//is a global program attribute, so shouldnt be here or in actions
var SysInitExprs []*ast.CXExpression

//cxgo/actions/declarations
//globals.SysInitExprs = append(globals.SysInitExprs, initializer...)
