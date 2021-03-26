package actions

import (
	"fmt"
	"time"

	"github.com/skycoin/cx/cx"
)

/*
./cxparser/stage2/cxparser.y:234:			actions.Stepping(int($2), int($3), true)
./cxparser/stage2/cxparser.y:238:			actions.Stepping(int($2), 0, false)
./cxparser/stage2/cxparser.go:1693:			actions.Stepping(int(yyS[yypt-1].i32), int(yyS[yypt-0].i32), true)
./cxparser/stage2/cxparser.go:1697:			actions.Stepping(int(yyS[yypt-0].i32), 0, false)
./cxparser/actions/interactive.go:10:func Stepping(steps int, delay int, withDelay bool) {
*/

//DELETE THIS, only calls RunCompiled
func SteppingNoDelay(steps int) {
	if steps < 0 {
		panic("error, should not be negative")
	}
	//steps=0 for running infinitely
	if err := PRGRM.RunCompiled(steps, nil); err != nil {
		fmt.Println(err)
	}
}

//DELETE THIS, only calls RunCompiled
func SteppingWithDelay(steps int, delay int, withDelay bool) {
	if steps < 0 {
		panic("error, should not be negative")
	}
	for i := 0; i < steps; i++ {
		time.Sleep(time.Duration(int32(delay)) * time.Second)
		err := PRGRM.RunCompiled(1, nil)
		if PRGRM.Terminated {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

//delete this function, its retarded
//Used Twice; in cxparser/stage2/cxparser.go
func Stepping(steps int, delay int, withDelay bool) {
	if !withDelay {
		SteppingNoDelay(steps)
	} else {
		SteppingWithDelay(steps, delay, withDelay)
	}
}

func Selector(ident string, selTyp int) string {
	switch selTyp {
	case SELECT_TYP_PKG:
		var previousModule *cxcore.CXPackage
		if mod, err := PRGRM.GetCurrentPackage(); err == nil {
			previousModule = mod
		} else {
			fmt.Println("a current package does not exist")
		}
		if _, err := PRGRM.SelectPackage(ident); err == nil {
			//fmt.Println(fmt.Sprintf("== Changed to package '%s' ==", mod.Name))
		} else {
			fmt.Println(err)
		}

		ReplTargetMod = ident
		ReplTargetStrct = ""
		ReplTargetFn = ""

		return previousModule.Name
	case SELECT_TYP_FUNC:
		var previousFunction *cxcore.CXFunction
		if fn, err := PRGRM.GetCurrentFunction(); err == nil {
			previousFunction = fn
		} else {
			fmt.Println("A current function does not exist")
		}
		if _, err := PRGRM.SelectFunction(ident); err == nil {
			//fmt.Println(fmt.Sprintf("== Changed to function '%s' ==", fn.Name))
		} else {
			fmt.Println(err)
		}

		ReplTargetMod = ""
		ReplTargetStrct = ""
		ReplTargetFn = ident

		return previousFunction.Name
	case SELECT_TYP_STRCT:
		var previousStruct *cxcore.CXStruct
		if fn, err := PRGRM.GetCurrentStruct(); err == nil {
			previousStruct = fn
		} else {
			fmt.Println("A current struct does not exist")
		}
		if _, err := PRGRM.SelectStruct(ident); err == nil {
			//fmt.Println(fmt.Sprintf("== Changed to struct '%s' ==", fn.Name))
		} else {
			fmt.Println(err)
		}

		ReplTargetStrct = ident
		ReplTargetMod = ""
		ReplTargetFn = ""

		return previousStruct.Name
	}

	panic("")

}
