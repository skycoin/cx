package main

import (
	"fmt"
	"log"
	"os"

	"math/rand"
	"time"
	//"unsafe"
	
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
	byt := MakeType("byte")
	i32 := MakeType("i32")
	i64 := MakeType("i64")
	abyt := MakeType("[]byte")
	ai32 := MakeType("[]i32")
	ai64 := MakeType("[]i64")

	_ = byt
	_ = abyt
	_ = i64
	_ = ai32
	_ = ai64
	
	ident := MakeType("ident")

	num1 := encoder.SerializeAtomic(int32(0))
	num2 := encoder.SerializeAtomic(int32(1))
	// num2 := encoder.SerializeAtomic(int32(40))
	// num3 := encoder.SerializeAtomic(int64(40))
	num4 := encoder.SerializeAtomic(int32(33))
	num5 := encoder.SerializeAtomic(int32(66))
	bytA := []byte{byte(5), byte(50), byte(100)}
	oneByte := []byte{byte(70)}


	cxt := MakeContext()
	cxt.AddModule(MakeModule("Math"))
	cxt.AddModule(MakeModule("Stdio"))
	cxt.AddModule(MakeModule("OpenGl"))
	cxt.AddModule(MakeModule("DeepLearning"))
	
	cxt.AddModule(MakeModule("main"))
	if mod, err := cxt.GetCurrentModule(); err == nil {

		if imp, err := cxt.GetModule("Math"); err == nil {
			mod.AddImport(imp)
		}
		
		mod.AddDefinition(MakeDefinition("num1", &num1, i32))
		mod.AddDefinition(MakeDefinition("index", &num1, i32))
		mod.AddDefinition(MakeDefinition("index2", &num2, i32))
		//mod.AddDefinition(MakeDefinition("num2", &num2, i32))
		mod.AddDefinition(MakeDefinition("bytes", &bytA, abyt))
		mod.AddDefinition(MakeDefinition("byte", &oneByte, byt))

		mod.AddStruct(MakeStruct("Point"))
		if strct, err := mod.GetCurrentStruct(); err == nil {
			strct.AddField(MakeField("x", i32))
			strct.AddField(MakeField("y", i32))

			// initializing a struct
			emptyBytes := []byte{}
			// the adder will initialize them
			mod.AddDefinition(MakeDefinition("myPoint", &emptyBytes, MakeType("Point")))
			// and we can change their values like this; just replacing the ones we want
			mod.AddDefinition(MakeDefinition("myPoint.x", &num4, i32))
			mod.AddDefinition(MakeDefinition("myPoint.y", &num5, i32))
		}

		mod.AddFunction(MakeFunction("writeAByte"))
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddInput(MakeParameter("arr", abyt))
			fn.AddInput(MakeParameter("idx", i32))
			fn.AddInput(MakeParameter("val", byt))
			fn.AddOutput(MakeParameter("out", abyt))
		}

		mod.AddFunction(MakeFunction("readAByte"))
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddInput(MakeParameter("arr", abyt))
			fn.AddInput(MakeParameter("idx", i32))
			fn.AddOutput(MakeParameter("out", byt))
		}

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
			if addI32, err := cxt.GetFunction("addI32", mod.Name); err == nil {
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
			fn.AddOutput(MakeParameter("fx", i32))
		}
		
		mod.SelectFunction("main")

		if fn, err := cxt.GetCurrentFunction(); err == nil {
			fn.AddOutput(MakeParameter("outMain", byt))
			
			if double, err := cxt.GetFunction("double", mod.Name); err == nil {
				fn.AddExpression(MakeExpression("outMain", double))
				
				if expr, err := cxt.GetCurrentExpression(); err == nil {
					expr.AddArgument(MakeArgument(MakeValue("myPoint.y"), ident))
				}
			}

			if writeAByte, err := cxt.GetFunction("writeAByte", mod.Name); err == nil {
				fn.AddExpression(MakeExpression("foo", writeAByte))
				
				if expr, err := cxt.GetCurrentExpression(); err == nil {
					expr.AddArgument(MakeArgument(MakeValue("bytes"), ident))
					//expr.AddArgument(MakeArgument(MakeValue("index2"), ident))
					expr.AddArgument(MakeArgument(&num2, i32))
					expr.AddArgument(MakeArgument(MakeValue("byte"), ident))
				}
			}

			if readAByte, err := cxt.GetFunction("readAByte", mod.Name); err == nil {
				fn.AddExpression(MakeExpression("outMain", readAByte))
				
				if expr, err := cxt.GetCurrentExpression(); err == nil {
					expr.AddArgument(MakeArgument(MakeValue("bytes"), ident))
					expr.AddArgument(MakeArgument(MakeValue("index"), ident))
				}
			}

			// if solution, err := cxt.GetFunction("solution", mod.Name); err == nil {
			// 	fn.AddExpression(MakeExpression("outMain", solution))
				
			// 	if expr, err := cxt.GetCurrentExpression(); err == nil {
			// 		expr.AddArgument(MakeArgument(MakeValue("num1"), ident))
			// 	}
			// }
		}
	}

	dataIn := []int32{-10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	dataOut := make([]int32, len(dataIn))
	for i, in := range dataIn {
		dataOut[i] = in * in * in - (3 * in)
	}

	cxt.Run(true, 1)
	sCxt := Serialize(cxt)
	//cxt.PrintProgram(false)
	dsCxt := Deserialize(sCxt)
	//Deserialize(sCxt)
	fmt.Println("==================================")
	dsCxt.Run(true, 1)
	//dsCxt.PrintProgram(false)


	cxt.Run(false, -1)
	fmt.Println(cxt.Output.Value)
	dsCxt.Run(false, -1)
	fmt.Println(dsCxt.Output.Value)

	

	// evolvedProgram := cxt.EvolveSolution(dataIn, dataOut, 5, 2000)
	// evolvedProgram.PrintProgram(false)

	
	// //getting the simulated outputs
	// var error int32 = 0
	// for i, inp := range dataIn {
	// 	num1 := encoder.SerializeAtomic(inp)
	// 	if def, err := evolvedProgram.GetDefinition("num1"); err == nil {
	// 		def.Value = &num1
	// 	} else {
	// 		fmt.Println(err)
	// 	}

	// 	evolvedProgram.Reset()
	// 	evolvedProgram.Run(false, -1)

	// 	// getting the simulated output
	// 	var result int32
	// 	output := evolvedProgram.CallStack.Calls[0].State["outMain"].Value
	// 	encoder.DeserializeAtomic(*output, &result)

	// 	diff := result - dataOut[i]
	// 	fmt.Printf("Simulated #%d: %d\n", i, result)
	// 	fmt.Printf("Observed #%d: %d\n", i, dataOut[i])
	// 	if diff >= 0 {
	// 		error += diff
	// 	} else {
	// 		error += diff * -1
	// 	}
	// }
	// fmt.Println(error / int32(len(dataIn)))

	





	

	// fmt.Println("Interpreted")
	// for i := 0; i < 30; i++ {
	// 	start := time.Now()

	// 	for j := 0; j < 10000; j++ {
	// 		evolvedProgram.ResetTo(0)
	// 		evolvedProgram.Run(false, -1)
	// 	}
		
	// 	elapsed := time.Since(start)
	// 	//log.Printf("Interpreted took %s", elapsed)
	// 	fmt.Println(elapsed)
	// }

	// evolvedProgram.Compile()

	// fmt.Println("Compiled")
	// for i := 0; i < 30; i++ {
	// 	start := time.Now()

	// 	for j := 0; j < 10000; j++ {
	// 		evolvedProgram.ResetTo(0)
	// 		evolvedProgram.Run(false, -1)
	// 	}
		
	// 	elapsed := time.Since(start)
	// 	//log.Printf("Compiled took %s", elapsed)
	// 	fmt.Println(elapsed)
	// }

	

	// cxt.ResetTo(0)
	// cxt.Compile()
	// cxt.PrintProgram(false)
	// cxt.Run(true, -1)


	

	// arr := []int32{0, 11111, 222, 333333, -44, -55555}
	// sArr := encoder.Serialize(arr)
	// fmt.Println(sArr)
	// var dArr []int32
	// encoder.DeserializeRaw(sArr, &dArr)
	// fmt.Println(dArr)

	
	//Testing()


	// Testing compiling structs

	// for _, def := range evolvedProgram.CurrentModule.Definitions {
	// 	fmt.Println(def.Name)
	// }
	
	//evolvedProgram.Compile()
	//fmt.Println(evolvedProgram.Heap)
}
