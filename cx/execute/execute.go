package execute

import (
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"math/rand"
	"time"
)

// Only called in this file
// TODO: What does this do? Is it named poorly?
func ToCall(cxprogram *ast.CXProgram) *ast.CXExpression {
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

func RunCxAst(cxprogram *ast.CXProgram, untilEnd bool, nCalls *int, untilCall int) error {
	defer ast.RuntimeError()
	var err error

	var inputs []ast.CXValue
	var outputs []ast.CXValue
	for !cxprogram.Terminated && (untilEnd || *nCalls != 0) && cxprogram.CallCounter > untilCall {
		call := &cxprogram.CallStack[cxprogram.CallCounter]

		// checking if enough memory in stack
		if cxprogram.StackPointer > constants.STACK_SIZE {
			panic(constants.STACK_OVERFLOW_ERROR)
		}

		if !untilEnd {
			var inName string
			var toCallName string
			var toCall *ast.CXExpression

			if call.Line >= call.Operator.Length && cxprogram.CallCounter == 0 {
				cxprogram.Terminated = true
				cxprogram.CallStack[0].Operator = nil
				cxprogram.CallCounter = 0
				fmt.Println("in:terminated")
				return err
			}

			if call.Line >= call.Operator.Length && cxprogram.CallCounter != 0 {
				toCall = ToCall(cxprogram)
				// toCall = cxprogram.CallStack[cxprogram.CallCounter-1].Operator.Expressions[cxprogram.CallStack[cxprogram.CallCounter-1].Line + 1]
				inName = cxprogram.CallStack[cxprogram.CallCounter-1].Operator.Name
			} else {
				toCall = call.Operator.Expressions[call.Line]
				inName = call.Operator.Name
			}

			if toCall.Operator == nil {
				// then it's a declaration
				toCallName = "declaration"
			} else if toCall.Operator.IsAtomic {
				toCallName = ast.OpNames[toCall.Operator.OpCode]
			} else {
				if toCall.Operator.Name != "" {
					toCallName = toCall.Operator.Package.Name + "." + toCall.Operator.Name
				} else {
					// then it's the end of the program got from nested function calls
					cxprogram.Terminated = true
					cxprogram.CallStack[0].Operator = nil
					cxprogram.CallCounter = 0
					fmt.Println("in:terminated")
					return err
				}
			}

			fmt.Printf("in:%s, expr#:%d, calling:%s()\n", inName, call.Line+1, toCallName)
			*nCalls--
		}

		err = call.Ccall(cxprogram, &inputs, &outputs)
		if err != nil {
			return err
		}
	}

	return nil
}

// RunCompiled ...
func RunCompiled(cxprogram *ast.CXProgram, nCalls int, args []string) error {
	_, err := cxprogram.SetCurrentCxProgram()
	if err != nil {
		panic(err)
	}
	cxprogram.EnsureMinimumHeapSize()
	rand.Seed(time.Now().UTC().UnixNano())

	var untilEnd bool
	if nCalls == 0 {
		untilEnd = true
	}
	mod, err := cxprogram.SelectPackage(constants.MAIN_PKG)
	if err == nil {
		// initializing program resources
		// cxprogram.Stacks = append(cxprogram.Stacks, MakeStack(1024))

		var inputs []ast.CXValue
		var outputs []ast.CXValue
		if cxprogram.CallStack[0].Operator == nil {
			// then the program is just starting and we need to run the SYS_INIT_FUNC
			if fn, err := mod.SelectFunction(constants.SYS_INIT_FUNC); err == nil {
				// *init function
				mainCall := MakeCall(fn)
				cxprogram.CallStack[0] = mainCall
				cxprogram.StackPointer = fn.Size

				var err error

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
			} else {
				return err
			}
		}

		if fn, err := mod.SelectFunction(constants.MAIN_FUNC); err == nil {
			if len(fn.Expressions) < 1 {
				return nil
			}

			if cxprogram.CallStack[0].Operator == nil {
				// main function
				mainCall := MakeCall(fn)
				mainCall.FramePointer = cxprogram.StackPointer
				// initializing program resources
				cxprogram.CallStack[0] = mainCall

				// cxprogram.Stacks = append(cxprogram.Stacks, MakeStack(1024))
				cxprogram.StackPointer += fn.Size

				// feeding os.Args
				if osPkg, err := ast.PROGRAM.SelectPackage(constants.OS_PKG); err == nil {
					argsOffset := 0
					if osGbl, err := osPkg.GetGlobal(constants.OS_ARGS); err == nil {
						for _, arg := range args {
							argBytes := encoder.Serialize(arg)
							argOffset := ast.AllocateSeq(len(argBytes) + constants.OBJECT_HEADER_SIZE)

							var header = make([]byte, constants.OBJECT_HEADER_SIZE)
							ast.WriteMemI32(header, 5, int32(encoder.Size(arg)+constants.OBJECT_HEADER_SIZE))
							obj := append(header, argBytes...)

							ast.WriteMemory(argOffset, obj)

							var argOffsetBytes [4]byte
							ast.WriteMemI32(argOffsetBytes[:], 0, int32(argOffset))
							argsOffset = ast.WriteToSlice(argsOffset, argOffsetBytes[:])
						}
						ast.WriteI32(ast.GetFinalOffset(0, osGbl), int32(argsOffset))
					}
				}
				cxprogram.Terminated = false
			}

			if err = RunCxAst(cxprogram, untilEnd, &nCalls, -1); err != nil {
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

			return err
		}
		return err

	}
	return err

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

