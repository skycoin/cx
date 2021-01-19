package actions

import (
	"fmt"
	"time"

	. "github.com/skycoin/cx/cx"
)

func Stepping(steps int, delay int, withDelay bool) {
	if withDelay {
		if steps == 0 {
			// Maybe nothing for now
		} else {
			if steps < 0 {
				PRGRM.UnRun(steps)
			} else {
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
		}
	} else {
		if steps == 0 {
			// we run until halt or end of program;
			if err := PRGRM.RunCompiled(0, nil); err != nil {
				fmt.Println(err)
			}
		} else {
			if steps < 0 {
				// nCalls := steps * -1
				// PRGRM.UnRun(int(nCalls))
				PRGRM.UnRun(steps)
			} else {
				PRGRM.RunCompiled(steps, nil)
				// err := PRGRM.RunInterpreted(dStack, int(steps))
				// if err != nil {
				// 	fmt.Println(err)
				// }
			}
		}
	}
}

func Selector(ident string, selTyp int) string {
	switch selTyp {
	case SELECT_TYP_PKG:
		var previousModule *CXPackage
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
		var previousFunction *CXFunction
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
		var previousStruct *CXStruct
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
