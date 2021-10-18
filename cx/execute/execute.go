package execute

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

func getLastLine(cxprogram *ast.CXProgram) *ast.CXExpression {
	for c := cxprogram.CallCounter - 1; c >= 0; c-- {
		if cxprogram.CallStack[c].Line+1 >= len(cxprogram.CallStack[c].Operator.Expressions) {
			// then it'll also return from this function call; continue
			continue
		}
		return cxprogram.CallStack[c].Operator.Expressions[cxprogram.CallStack[c].Line+1]
		// cxprogram.CallStack[c].Operator.Expressions[cxprogram.CallStack[cxprogram.CallCounter-1].Line + 1]
	}
	// error
	return &ast.CXExpression{Operator: ast.MakeFunction("", "", -1)}
	// panic("")
}

func RunCxAst(cxprogram *ast.CXProgram, untilEnd bool, maxOps int, untilCall types.Pointer) error {
	defer ast.RuntimeError()

	var inputs []ast.CXValue
	var outputs []ast.CXValue
	var opCount int = 0
	for !cxprogram.Terminated && (untilEnd || (opCount < maxOps)) && (!untilCall.IsValid() || cxprogram.CallCounter > untilCall) {
		call := &cxprogram.CallStack[cxprogram.CallCounter]

		// checking if enough memory in stack
		if cxprogram.Stack.Pointer > constants.STACK_SIZE {
			panic(constants.STACK_OVERFLOW_ERROR)
		}

		if !untilEnd {
			var inName string
			var toCallName string
			var toCall *ast.CXExpression

			if call.Line >= call.Operator.LineCount && cxprogram.CallCounter == 0 {
				cxprogram.Terminated = true
				cxprogram.CallStack[0].Operator = nil
				cxprogram.CallCounter = 0
				fmt.Println("in:terminated")
				return nil
			}

			if call.Line >= call.Operator.LineCount && cxprogram.CallCounter != 0 {
				toCall = getLastLine(cxprogram)
				// toCall = cxprogram.CallStack[cxprogram.CallCounter-1].Operator.Expressions[cxprogram.CallStack[cxprogram.CallCounter-1].Line + 1]
				inName = cxprogram.CallStack[cxprogram.CallCounter-1].Operator.Name
			} else {
				toCall = call.Operator.Expressions[call.Line]
				inName = call.Operator.Name
			}

			if toCall.Operator == nil {
				// then it's a declaration
				toCallName = "declaration"
			} else if toCall.Operator.IsBuiltIn() {
				toCallName = ast.OpNames[toCall.Operator.AtomicOPCode]
			} else {
				if toCall.Operator.Name != "" {
					toCallName = toCall.Operator.Package.Name + "." + toCall.Operator.Name
				} else {
					// then it's the end of the program got from nested function calls
					cxprogram.Terminated = true
					cxprogram.CallStack[0].Operator = nil
					cxprogram.CallCounter = 0
					fmt.Println("in:terminated")
					return nil
				}
			}

			fmt.Printf("in:%s, expr#:%d, calling:%s()\n", inName, call.Line+1, toCallName)
			opCount++
		}

		err := call.Ccall(cxprogram, &inputs, &outputs)
		if err != nil {
			return err
		}
	}

	return nil
}

// RunCompiled ...
func RunCompiled(cxprogram *ast.CXProgram, maxOps int, args []string) error {
	_, err := cxprogram.SetCurrentCxProgram()
	if err != nil {
		panic(err)
	}

	cxprogram.EnsureMinimumHeapSize()
	rand.Seed(time.Now().UTC().UnixNano())

	var untilEnd bool
	if maxOps == 0 {
		untilEnd = true
	}

	mod, err := cxprogram.SelectPackage(constants.MAIN_PKG)
	if err != nil {
		return err
	}
	// initializing program resources
	// cxprogram.Stacks = append(cxprogram.Stacks, MakeStack(1024))

	var inputs []ast.CXValue
	var outputs []ast.CXValue
	if cxprogram.CallStack[0].Operator == nil {
		// then the program is just starting and we need to run the SYS_INIT_FUNC
		fn, err := mod.SelectFunction(constants.SYS_INIT_FUNC)
		if err != nil {
			return err
		}

		// *init function
		mainCall := MakeCall(fn)
		cxprogram.CallStack[0] = mainCall
		cxprogram.Stack.Pointer = fn.Size

		for !cxprogram.Terminated {
			call := &cxprogram.CallStack[cxprogram.CallCounter]
			err = call.Ccall(cxprogram, &inputs, &outputs)
			if err != nil {
				return err
			}
		}
		// we reset call state
		cxprogram.Terminated = false
		cxprogram.CallCounter = 0
		cxprogram.CallStack[0].Operator = nil
	}

	fn, err := mod.SelectFunction(constants.MAIN_FUNC)
	if err != nil {
		return err
	}

	if len(fn.Expressions) < 1 {
		return nil
	}

	if cxprogram.CallStack[0].Operator == nil {
		// main function
		mainCall := MakeCall(fn)
		mainCall.FramePointer = cxprogram.Stack.Pointer
		// initializing program resources
		cxprogram.CallStack[0] = mainCall

		// cxprogram.Stacks = append(cxprogram.Stacks, MakeStack(1024))
		cxprogram.Stack.Pointer += fn.Size

		// feeding os.Args
		if osPkg, err := ast.PROGRAM.SelectPackage(constants.OS_PKG); err == nil {
			argsOffset := types.Pointer(0)
			if osGbl, err := osPkg.GetGlobal(constants.OS_ARGS); err == nil {
				for _, arg := range args {
					argOffset := types.AllocWrite_obj_data(cxprogram.Memory, []byte(arg))

					var argOffsetBytes [types.POINTER_SIZE]byte
					types.Write_ptr(argOffsetBytes[:], 0, argOffset)
					argsOffset = ast.WriteToSlice(argsOffset, argOffsetBytes[:])
				}
				types.Write_ptr(ast.PROGRAM.Memory, ast.GetFinalOffset(0, osGbl), argsOffset)
			}
		}
		cxprogram.Terminated = false
	}

	if err = RunCxAst(cxprogram, untilEnd, maxOps, types.InvalidPointer); err != nil {
		return err
	}

	if cxprogram.Terminated {
		cxprogram.Terminated = false
		cxprogram.CallCounter = 0
		cxprogram.CallStack[0].Operator = nil
	}

	// debugging memory
	// if len(cxprogram.Memory) < 2000 {
	// 	fmt.Println("cxprogram.Memory", cxprogram.Memory)
	// }

	return nil
}

func MakeCall(op *ast.CXFunction) ast.CXCall {
	return ast.CXCall{
		Operator:     op,
		Line:         0,
		FramePointer: 0,
		// Package:       pkg,
		// Program:       prgrm,
	}
}
