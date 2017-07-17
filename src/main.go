package main

import (
	"fmt"
	"log"
	"os"

	"math/rand"
	"time"
	"unsafe"
	
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

	evolvedProgram := cxt.EvolveSolution(dataIn, dataOut, 5, 5)
	evolvedProgram.PrintProgram(false)

	//getting the simulated outputs
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
		output := evolvedProgram.CallStack.Calls[0].State["outMain"].Value
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
	

	//copy := MakeContextCopy(evolvedProgram, 8)
	//copy.Run(true, -1)
	
	evolvedProgram.ResetTo(0)
	evolvedProgram.Run(true, -1)

	//foo := encoder.SerializeAtomic(int32(0))
	foo := encoder.Serialize("hello world, I'm happy to be here")
	bar := []byte("hello")
	fmt.Println(encoder.Size(foo))
	fmt.Println(unsafe.Sizeof(foo))
	fmt.Println(encoder.Size(bar))
	fmt.Println(unsafe.Sizeof(bar))

	Testing()
}
