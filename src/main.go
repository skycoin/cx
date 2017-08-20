package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	//"bytes"

	"math/rand"
	"time"
	//"unsafe"

	// "github.com/mndrix/golog"
	// "github.com/mndrix/golog/read"
	// "github.com/mndrix/golog/term"
	
	"github.com/skycoin/skycoin/src/cipher/encoder"
	//. "github.com/skycoin/cx/src/base"
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

func readline (fi *bufio.Reader) (string, bool) {
	s, err := fi.ReadString('\n')
	if err != nil {
		return "", false
	}
	return s, true
}

func main() {
	fmt.Printf("Byte: %d\n", len([]byte{0}))
	
	var i32 int32
	si32 := encoder.Serialize(i32)
	fmt.Printf("Int32: %d\n", len(si32))

	var i64 int64
	si64 := encoder.Serialize(i64)
	fmt.Printf("Int64: %d\n", len(si64))

	var f32 float32
	sf32 := encoder.Serialize(f32)
	fmt.Printf("Float32: %d\n", len(sf32))

	var f64 float64
	sf64 := encoder.Serialize(f64)
	fmt.Printf("Float64: %d\n", len(sf64))
	

	
	// // creating the machine
	// m := golog.NewInteractiveMachine()

	// fi := bufio.NewReader(os.NewFile(0, "stdin"))

	// initializing golog's database
	// fmt.Println("Initializing")
	// for {
	// 	var inp string
	// 	var ok bool

	// 	fmt.Printf("?- ")
	// 	if inp, ok = readline(fi); ok {
	// 		b := bytes.NewBufferString(inp)
	// 		m = m.Consult(b)
	// 	} else {
	// 		break
	// 	}
	// }

	// fmt.Println("Querying")
	// for {
	// 	var inp string
	// 	var ok bool

	// 	fmt.Printf("?- ")
	// 	if inp, ok = readline(fi); ok {
	// 		//b := bytes.NewBufferString(inp)

	// 		goal, _ := read.Term(inp)
			
	// 		variables := term.Variables(goal)
	// 		answers := m.ProveAll(goal)


	// 		if len(answers) == 0 {
	// 			fmt.Println("no")
	// 			continue
	// 		}

	// 		if variables.Size() == 0 {
	// 			fmt.Println("yes")
	// 			continue
	// 		}

			
	// 		//fmt.Println(answers)
	// 	} else {
	// 		break
	// 	}
	// }
	
	// byt := MakeType("byte")
	// i32 := MakeType("i32")
	// i64 := MakeType("i64")
	// abyt := MakeType("[]byte")
	// ai32 := MakeType("[]i32")
	// ai64 := MakeType("[]i64")

	// _ = byt
	// _ = abyt
	// _ = i64
	// _ = ai32
	// _ = ai64
	
	// ident := MakeType("ident")

	// num1 := encoder.SerializeAtomic(int32(0))
	// num2 := encoder.SerializeAtomic(int32(1))
	// // num2 := encoder.SerializeAtomic(int32(40))
	// // num3 := encoder.SerializeAtomic(int64(40))
	// num4 := encoder.SerializeAtomic(int32(33))
	// num5 := encoder.SerializeAtomic(int32(66))
	// bytA := []byte{byte(5), byte(50), byte(100)}
	// oneByte := []byte{byte(70)}


	// cxt := MakeContext()
	// cxt.AddModule(MakeModule("Math"))
	// cxt.AddModule(MakeModule("Stdio"))
	// cxt.AddModule(MakeModule("OpenGl"))
	// cxt.AddModule(MakeModule("DeepLearning"))
	
	// cxt.AddModule(MakeModule("main"))
	// if mod, err := cxt.GetCurrentModule(); err == nil {

	// 	if imp, err := cxt.GetModule("Math"); err == nil {
	// 		mod.AddImport(imp)
	// 	}
		
	// 	mod.AddDefinition(MakeDefinition("num1", &num1, i32))
	// 	mod.AddDefinition(MakeDefinition("index", &num1, i32))
	// 	mod.AddDefinition(MakeDefinition("index2", &num2, i32))
	// 	//mod.AddDefinition(MakeDefinition("num2", &num2, i32))
	// 	mod.AddDefinition(MakeDefinition("bytes", &bytA, abyt))
	// 	mod.AddDefinition(MakeDefinition("byte", &oneByte, byt))

	// 	mod.AddStruct(MakeStruct("Point"))
	// 	if strct, err := mod.GetCurrentStruct(); err == nil {
	// 		strct.AddField(MakeField("x", i32))
	// 		strct.AddField(MakeField("y", i32))

	// 		// initializing a struct
	// 		emptyBytes := []byte{}
	// 		// the adder will initialize them
	// 		mod.AddDefinition(MakeDefinition("myPoint", &emptyBytes, MakeType("Point")))
	// 		// and we can change their values like this; just replacing the ones we want
	// 		mod.AddDefinition(MakeDefinition("myPoint.x", &num4, i32))
	// 		mod.AddDefinition(MakeDefinition("myPoint.y", &num5, i32))
	// 	}

	// 	mod.AddFunction(MakeFunction("writeByteA"))
	// 	if fn, err := cxt.GetCurrentFunction(); err == nil {
	// 		fn.AddInput(MakeParameter("arr", abyt))
	// 		fn.AddInput(MakeParameter("idx", i32))
	// 		fn.AddInput(MakeParameter("val", byt))
	// 		fn.AddOutput(MakeParameter("out", abyt))
	// 	}

	// 	mod.AddFunction(MakeFunction("readByteA"))
	// 	if fn, err := cxt.GetCurrentFunction(); err == nil {
	// 		fn.AddInput(MakeParameter("arr", abyt))
	// 		fn.AddInput(MakeParameter("idx", i32))
	// 		fn.AddOutput(MakeParameter("out", byt))
	// 	}

	// 	mod.AddFunction(MakeFunction("addI32"))
	// 	if fn, err := cxt.GetCurrentFunction(); err == nil {
	// 		fn.AddInput(MakeParameter("n1", i32))
	// 		fn.AddInput(MakeParameter("n2", i32))
	// 		fn.AddOutput(MakeParameter("out", i32))
	// 	}
		
	// 	mod.AddFunction(MakeFunction("subI32"))
	// 	if fn, err := cxt.GetCurrentFunction(); err == nil {
	// 		fn.AddInput(MakeParameter("n1", i32))
	// 		fn.AddInput(MakeParameter("n2", i32))
	// 		fn.AddOutput(MakeParameter("out", i32))
	// 	}

	// 	mod.AddFunction(MakeFunction("mulI32"))
	// 	if fn, err := cxt.GetCurrentFunction(); err == nil {
	// 		fn.AddInput(MakeParameter("n1", i32))
	// 		fn.AddInput(MakeParameter("n2", i32))
	// 		fn.AddOutput(MakeParameter("out", i32))
	// 	}
		
	// 	mod.AddFunction(MakeFunction("main"))
	// 	mod.AddFunction(MakeFunction("double"))

	// 	if fn, err := cxt.GetCurrentFunction(); err == nil {
	// 		fn.AddInput(MakeParameter("n", i32))
	// 		fn.AddOutput(MakeParameter("out", i32))
	// 		fn.AddOutput(MakeParameter("error", i32))
	// 		if addI32, err := cxt.GetFunction("addI32", mod.Name); err == nil {
	// 			fn.AddExpression(MakeExpression(addI32))

	// 			if expr, err := cxt.GetCurrentExpression(); err == nil {
	// 				expr.AddOutputName("out")
	// 				expr.AddArgument(MakeArgument(MakeValue("n"), ident))
	// 				expr.AddArgument(MakeArgument(MakeValue("n"), ident))
	// 			}
	// 		}

	// 		if mulI32, err := cxt.GetFunction("mulI32", mod.Name); err == nil {
	// 			fn.AddExpression(MakeExpression(mulI32))

	// 			if expr, err := cxt.GetCurrentExpression(); err == nil {
	// 				expr.AddOutputName("error")
	// 				expr.AddArgument(MakeArgument(MakeValue("n"), ident))
	// 				expr.AddArgument(MakeArgument(MakeValue("n"), ident))
	// 			}
	// 		}
	// 	}

	// 	mod.AddFunction(MakeFunction("solution"))

	// 	if fn, err := cxt.GetCurrentFunction(); err == nil {
	// 		fn.AddInput(MakeParameter("x", i32))
	// 		fn.AddOutput(MakeParameter("fx", i32))

	// 		// if double, err := cxt.GetFunction("double", mod.Name); err == nil {
	// 		// 	fn.AddExpression(MakeExpression(double))

	// 		// 	if expr, err := cxt.GetCurrentExpression(); err == nil {
	// 		// 		expr.AddOutputName("fx")
	// 		// 		expr.AddOutputName("outMain")
					
	// 		// 		expr.AddArgument(MakeArgument(MakeValue("myPoint.22y"), ident))
	// 		// 	}
	// 		// }
			
	// 	}
		
	// 	mod.SelectFunction("main")

	// 	if fn, err := cxt.GetCurrentFunction(); err == nil {
	// 		fn.AddOutput(MakeParameter("outMain", i32))
			
	// 		// if double, err := cxt.GetFunction("double", mod.Name); err == nil {
	// 		// 	fn.AddExpression(MakeExpression(double))
				
	// 		// 	if expr, err := cxt.GetCurrentExpression(); err == nil {
	// 		// 		expr.AddOutputName("err")
	// 		// 		expr.AddOutputName("outMain")
					
	// 		// 		expr.AddArgument(MakeArgument(MakeValue("myPoint.y"), ident))
	// 		// 	}
	// 		// }

	// 		// if writeByteA, err := cxt.GetFunction("writeByteA", mod.Name); err == nil {
	// 		// 	fn.AddExpression(MakeExpression(writeByteA))
				
	// 		// 	if expr, err := cxt.GetCurrentExpression(); err == nil {
	// 		// 		expr.AddOutputName("foo")
	// 		// 		expr.AddArgument(MakeArgument(MakeValue("bytes"), ident))
	// 		// 		//expr.AddArgument(MakeArgument(MakeValue("index2"), ident))
	// 		// 		expr.AddArgument(MakeArgument(&num2, i32))
	// 		// 		expr.AddArgument(MakeArgument(MakeValue("byte"), ident))
	// 		// 	}
	// 		// }

	// 		// if readByteA, err := cxt.GetFunction("readByteA", mod.Name); err == nil {
	// 		// 	fn.AddExpression(MakeExpression(readByteA))
				
	// 		// 	if expr, err := cxt.GetCurrentExpression(); err == nil {
	// 		// 		expr.AddOutputName("outMain")
	// 		// 		expr.AddArgument(MakeArgument(MakeValue("bytes"), ident))
	// 		// 		expr.AddArgument(MakeArgument(MakeValue("index"), ident))
	// 		// 	}
	// 		// }

	// 		if solution, err := cxt.GetFunction("solution", mod.Name); err == nil {
	// 			fn.AddExpression(MakeExpression(solution))
				
	// 			if expr, err := cxt.GetCurrentExpression(); err == nil {
	// 				expr.AddOutputName("outMain")
	// 				expr.AddArgument(MakeArgument(MakeValue("num1"), ident))
	// 			}
	// 		}
	// 	}
	// }
	
	// dataIn := []int32{-10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// dataOut := make([]int32, len(dataIn))
	// for i, in := range dataIn {
	// 	dataOut[i] = in * in * in - (3 * in)
	// }

	// //cxt.Run(true, -1)
	// //sCxt := Serialize(cxt)
	// //cxt.PrintProgram(false)

	// // dsCxt := Deserialize(sCxt)
	// // //Deserialize(sCxt)
	// // fmt.Println("===================================")
	// // dsCxt.Run(true, -1)
	// // //dsCxt.PrintProgram(false)


	// // cxt.Run(false, -1)
	// // fmt.Println(cxt.Outputs[0].Value)
	// // dsCxt.Run(false, -1)
	// // fmt.Println(dsCxt.Outputs[0].Value)

	

	// //evolvedProgram := cxt.Evolve("solution", dataIn, dataOut, 10, 2, 0.001)
	// //evolvedProgram.PrintProgram(false)

	// encoded := encoder.Serialize(float64(1234))
	// fmt.Println(encoded)
	// var decoded float64
	// encoder.DeserializeRaw(encoded, &decoded)
	// fmt.Println(decoded)

	
	// // //getting the simulated outputs
	// // var error int32 = 0
	// // for i, inp := range dataIn {
	// // 	num1 := encoder.SerializeAtomic(inp)
	// // 	if def, err := evolvedProgram.GetDefinition("num1"); err == nil {
	// // 		def.Value = &num1
	// // 	} else {
	// // 		fmt.Println(err)
	// // 	}

	// // 	evolvedProgram.Reset()
	// // 	evolvedProgram.Run(false, -1)

	// // 	// getting the simulated output
	// // 	var result int32
	// // 	output := evolvedProgram.CallStack.Calls[0].State["outMain"].Value
	// // 	encoder.DeserializeAtomic(*output, &result)

	// // 	diff := result - dataOut[i]
	// // 	fmt.Printf("Simulated #%d: %d\n", i, result)
	// // 	fmt.Printf("Observed #%d: %d\n", i, dataOut[i])
	// // 	if diff >= 0 {
	// // 		error += diff
	// // 	} else {
	// // 		error += diff * -1
	// // 	}
	// // }
	// // fmt.Println(error / int32(len(dataIn)))

	





	

	// // fmt.Println("Interpreted")
	// // for i := 0; i < 30; i++ {
	// // 	start := time.Now()

	// // 	for j := 0; j < 10000; j++ {
	// // 		cxt.ResetTo(0)
	// // 		cxt.Run(false, -1)
	// // 	}
		
	// // 	elapsed := time.Since(start)
	// // 	fmt.Println(elapsed)
	// // }

	// // cxt.Compile()

	// // fmt.Println("Compiled")
	// // for i := 0; i < 30; i++ {
	// // 	start := time.Now()

	// // 	for j := 0; j < 10000; j++ {
	// // 		cxt.ResetTo(0)
	// // 		cxt.Run(false, -1)
	// // 	}
		
	// // 	elapsed := time.Since(start)
	// // 	fmt.Println(elapsed)
	// // }

	

	// // cxt.ResetTo(0)
	// // cxt.Compile()
	// // cxt.PrintProgram(false)
	// // cxt.Run(true, -1)


	

	// // arr := []int32{0, 11111, 222, 333333, -44, -55555}
	// // sArr := encoder.Serialize(arr)
	// // fmt.Println(sArr)
	// // var dArr []int32
	// // encoder.DeserializeRaw(sArr, &dArr)
	// // fmt.Println(dArr)

	
	// //Testing()


	// // Testing compiling structs

	// // for _, def := range evolvedProgram.CurrentModule.Definitions {
	// // 	fmt.Println(def.Name)
	// // }
	
	// //evolvedProgram.Compile()
	// //fmt.Println(evolvedProgram.Heap)
}
