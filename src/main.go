package main

import (
	"fmt"
	"log"
	"os"

	"math/rand"
	"time"
	
	"github.com/skycoin/skycoin/src/cipher/encoder"
	. "github.com/skycoin/cx/src/base"
)

func init() {
	log.SetOutput(os.Stdout)
}

func allClear(e []error) bool {
	for i := 1; i < len(e); i++ {
		if e[i] != nil {
			return false
		}
	}
	return true
}

func dbg(elt string) {
	fmt.Println(elt)
}

func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max - min) + min
}

func main() {
	i8 := MakeType("i8")
	i32 := MakeType("i32")
	i64 := MakeType("i64")

	_ = i8
	_ = i64
	
	ident := MakeType("ident")

	num1 := encoder.SerializeAtomic(int32(25))
	num2 := encoder.SerializeAtomic(int32(40))

	cxt := MakeContext().AddModule(MakeModule("main"))
	if mod, err := cxt.GetCurrentModule(); err == nil {
		mod.AddDefinition(MakeDefinition("num1", &num1, i32))
		mod.AddDefinition(MakeDefinition("num2", &num2, i32))

		mod.AddFunction(MakeFunction("addI32"))
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddInput(MakeParameter("n1", i32))
			fn.AddInput(MakeParameter("n2", i32))
			fn.AddOutput(MakeParameter("out", i32))
		}
		
		mod.AddFunction(MakeFunction("subI32"))
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddInput(MakeParameter("n1", i32))
			fn.AddInput(MakeParameter("n2", i32))
			fn.AddOutput(MakeParameter("out", i32))
		}

		mod.AddFunction(MakeFunction("mulI32"))
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddInput(MakeParameter("n1", i32))
			fn.AddInput(MakeParameter("n2", i32))
			fn.AddOutput(MakeParameter("out", i32))
		}
		
		mod.AddFunction(MakeFunction("main"))
		mod.AddFunction(MakeFunction("double"))

		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddInput(MakeParameter("n", i32))
			fn.AddOutput(MakeParameter("out", i32))
			if addI32, err := cxt.GetFunction("mulI32", mod.Name); err == nil {
				fn.AddExpression(MakeExpression("out", addI32))

				if expr, err := cxt.GetCurrentExpression(); err == nil {
					expr.AddArgument(MakeArgument(MakeValue("n"), ident))
					expr.AddArgument(MakeArgument(MakeValue("n"), ident))
				}
			}
		}

		mod.AddFunction(MakeFunction("solution"))

		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddInput(MakeParameter("x", i32))
			fn.AddOutput(MakeParameter("f(x)", i32))
		}
		
		mod.SelectFunction("main")

		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddOutput(MakeParameter("outMain", i32))
			if solution, err := cxt.GetFunction("solution", mod.Name); err == nil {
				fn.AddExpression(MakeExpression("outMain", solution))
				
				if expr, err := cxt.GetCurrentExpression(); err == nil {
					expr.AddArgument(MakeArgument(MakeValue("num2"), ident))
				}
			}
		}
	}

	dataIn := []int32{-10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	dataOut := make([]int32, len(dataIn))
	for i, in := range dataIn {
		dataOut[i] = in * in * in - (3 * in)
	}

	evolvedProgram := cxt.EvolveSolution(dataIn, dataOut, 5, 10000)
	evolvedProgram.PrintProgram(false)

	// getting the simulated outputs
	var error int32 = 0
	for i, inp := range dataIn {
		num1 := encoder.SerializeAtomic(inp)
		if def, err := evolvedProgram.GetDefinition("num1"); err == nil {
			def.Value = &num1
		} else {
			fmt.Println(err)
		}

		evolvedProgram.Reset()
		evolvedProgram.Run(false, -1)

		// getting the simulated output
		var result int32
		output := evolvedProgram.CallStack[0].State["outMain"].Value
		encoder.DeserializeAtomic(*output, &result)

		diff := result - dataOut[i]
		fmt.Printf("Simulated #%d: %d\n", i, result)
		fmt.Printf("Observed #%d: %d\n", i, dataOut[i])
		if diff >= 0 {
			error += diff
		} else {
			error += diff * -1
		}
	}
	fmt.Println(error / int32(len(dataIn)))
	
	

	//fmt.Println(len(cxt.ProgramSteps))

	// cxtCopy := MakeContext()
	// cxtCopy2 := MakeContext()
	// cxtCopy3 := MakeContext()
	
	// fmt.Println("Before copying")
	//cxt.PrintProgram(false)
	// cxtCopy.PrintProgram(false)

	//fmt.Println(len(cxt.ProgramSteps))
	
	// for i := 0; i < 15; i++ {
	// 	cxt.ProgramSteps[i].Action(cxtCopy)
	// }

	// for i := 0; i < 5; i++ {
	// 	cxt.ProgramSteps[i].Action(cxtCopy2)
	// }
	
	// for i := 0; i < 10; i++ {
	// 	cxt.ProgramSteps[i].Action(cxtCopy3)
	// }
	
	// for _, step := range cxt.ProgramSteps {
	// 	step.Action(cxtCopy)
	// }

	// fmt.Println("After copying")
	// //cxt.PrintProgram(false)
	// cxtCopy.PrintProgram(false)
	// cxtCopy2.PrintProgram(false)
	// cxtCopy3.PrintProgram(false)

	// fmt.Println("...")
	
	// if mainFn, err := cxt.SelectFunction("main"); err == nil {
	// 	fmt.Println(len(mainFn.Expressions))
	// }

	// if mainFn, err := cxtCopy.SelectFunction("main"); err == nil {
	// 	fmt.Println(len(mainFn.Expressions))
	// }

	// fmt.Println("...")
	// for _, fn := range cxt.CurrentModule.Functions {
	// 	fmt.Printf("%s %d\n", fn.Name, len(fn.Expressions))
	// }
	// fmt.Println("...")
	// for _, fn := range cxtCopy.CurrentModule.Functions {
	// 	fmt.Printf("%s %d\n", fn.Name, len(fn.Expressions))
	// }
	

	// addI32 := MakeFunction("addI32")
	// double := MakeFunction("double").
	// 	AddInput(MakeParameter("n", i32)).
	// 	AddOutput(MakeParameter("out", i32)).
	// 	AddExpression(
	// 	MakeExpression("notOut", addI32).
	// 		AddArgument(MakeArgument(MakeValue("n"), ident)).
	// 		AddArgument(MakeArgument(MakeValue("n"), ident))).
	// 	AddExpression(
	// 	MakeExpression("out", addI32).
	// 		AddArgument(MakeArgument(MakeValue("notOut"), ident)).
	// 		AddArgument(MakeArgument(MakeValue("num1"), ident)))

	// quad := MakeFunction("quad").
	// 	AddInput(MakeParameter("n", i32)).
	// 	AddOutput(MakeParameter("out", i32)).
	// 	AddExpression(MakeExpression("n", double).
	// 	  AddArgument(MakeArgument(MakeValue("n"), ident))).
	// 	AddExpression(MakeExpression("out", double).
	// 	AddArgument(MakeArgument(MakeValue("n"), ident)))

	// sumQuads := MakeFunction("sumQuads").
	// 	AddInput(MakeParameter("n1", i32)).
	// 	AddInput(MakeParameter("n2", i32)).
	// 	AddOutput(MakeParameter("result", i32)).
	// 	AddExpression(MakeExpression("quad1", quad).
	// 	AddArgument(MakeArgument(MakeValue("n1"), ident))).
	// 	AddExpression(MakeExpression("quad2", quad).
	// 	AddArgument(MakeArgument(MakeValue("n2"), ident))).
	// 	AddExpression(MakeExpression("result", addI32).
	// 	AddArgument(MakeArgument(MakeValue("quad1"), ident)).
	// 	AddArgument(MakeArgument(MakeValue("quad2"), ident)))

	// main := MakeFunction("main").
	// 	AddExpression(MakeExpression("outputMain", sumQuads).
	// 	AddArgument(MakeArgument(MakeValue("num1"), ident)).
	// 	//AddArgument(MakeArgument(&num2, i32)).
	// 	AddArgument(MakeArgument(&num2, i32)))




	
	



	// if mod, err := cxt.GetCurrentModule(); err == nil {
	// 	// native functions only need a name to reference them
		
	// 	mod.AddFunction(addI32)
	// 	mod.AddFunction(double)
	// 	mod.AddFunction(quad)
	// 	mod.AddFunction(sumQuads)
	// 	mod.AddFunction(main)
	// } else {
	// 	fmt.Println(err)
	// }
	
	// // Generating random program
	// //cxt := RandomProgram(200)
	// //cxt.PrintProgram(false)

	// // Call stack testing
	// cxt.Run(false, 2000)
	// fmt.Println(len(cxt.Steps))
	// fmt.Println("...")
	// //PrintCallStack(cxt.Steps[len(cxt.Steps) - 1])
	// PrintCallStack(cxt.CallStack)
	// //cxt.Run(false, 200)
	// fmt.Println(len(cxt.Steps))

	// // for _, step := range cxt.Steps {
	// // 	//fmt.Println(step)
	// // 	//PrintCallStack(step)
	// // 	fmt.Printf("%p\n", step[len(step) - 1].Context)
	// // }
	// //PrintCallStack(cxt.CallStack)	


	// fmt.Println("...")

	
	// // copy of program at about the middle of its execution
	// cxtCopy := MakeContextCopy(cxt, 1000)
	// // if defs, err := cxtCopy.GetCurrentDefinitions(); err == nil {
	// // 	defs["num1"].Value = &num2
	// // }
	// //cxt.Run(false, 200)
	// // fmt.Println("...")

	// //PrintCallStack(cxtCopy.Steps[len(cxtCopy.Steps) - 1])
	// PrintCallStack(cxtCopy.CallStack)
	// //fmt.Println(len(cxtCopy.Steps))
	// //fmt.Println(len(cxt.Steps))
	// cxtCopy.Run(false, 2000)
	// //fmt.Println(len(cxtCopy.Steps))
	// // fmt.Println("...")
	// fmt.Println(len(cxt.Steps))
	// fmt.Println(len(cxtCopy.Steps))
	// // PrintCallStack(cxtCopy.Steps[len(cxtCopy.Steps) - 1])
	// fmt.Println("checking calls contexts")
	// // for _, step := range cxtCopy.Steps {
	// // 	fmt.Printf("%p\n", step[len(step) - 1].Context)
	// // }

	// for i := 0; i < len(cxtCopy.Steps); i++ {
	// 	fmt.Printf("...%d\n", i)
	// 	//PrintCallStack(cxt.Steps[i])
	// 	PrintCallStack(cxtCopy.Steps[i])
	// }
}
