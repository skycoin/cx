package main

import (
	"fmt"
	"log"
	"os"
	
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

func main() {
	i8 := MakeType("i8")
	i32 := MakeType("i32")
	i64 := MakeType("i64")

	_ = i8
	_ = i64
	
	ident := MakeType("ident")

	num1 := encoder.SerializeAtomic(int32(25))
	num2 := encoder.SerializeAtomic(int32(40))
	//inum1 := []byte("num1")
	//inum2 := []byte("num2")

	cxt := MakeContext().AddModule(MakeModule("main").
		AddDefinition(MakeDefinition("num1", &num1, i32)).
		AddDefinition(MakeDefinition("num2", &num2, i32)))


	addI32 := MakeFunction("addI32")
	double := MakeFunction("double").
		AddInput(MakeParameter("n", i32)).
		AddOutput(MakeParameter("out", i32)).
		AddExpression(
		MakeExpression("notOut", addI32).
			AddArgument(MakeArgument(MakeValue("n"), ident)).
			AddArgument(MakeArgument(MakeValue("n"), ident))).
		AddExpression(
		MakeExpression("out", addI32).
			AddArgument(MakeArgument(MakeValue("notOut"), ident)).
			AddArgument(MakeArgument(MakeValue("num1"), ident)))

	quad := MakeFunction("quad").
		AddInput(MakeParameter("n", i32)).
		AddOutput(MakeParameter("out", i32)).
		AddExpression(MakeExpression("n", double).
		  AddArgument(MakeArgument(MakeValue("n"), ident))).
		AddExpression(MakeExpression("out", double).
		AddArgument(MakeArgument(MakeValue("n"), ident)))

	sumQuads := MakeFunction("sumQuads").
		AddInput(MakeParameter("n1", i32)).
		AddInput(MakeParameter("n2", i32)).
		AddOutput(MakeParameter("result", i32)).
		AddExpression(MakeExpression("quad1", quad).
		AddArgument(MakeArgument(MakeValue("n1"), ident))).
		AddExpression(MakeExpression("quad2", quad).
		AddArgument(MakeArgument(MakeValue("n2"), ident))).
		AddExpression(MakeExpression("result", addI32).
		AddArgument(MakeArgument(MakeValue("quad1"), ident)).
		AddArgument(MakeArgument(MakeValue("quad2"), ident)))

	main := MakeFunction("main").
		AddExpression(MakeExpression("outputMain", sumQuads).
		AddArgument(MakeArgument(MakeValue("num1"), ident)).
		//AddArgument(MakeArgument(&num2, i32)).
		AddArgument(MakeArgument(&num2, i32)))




	
	



	if mod, err := cxt.GetCurrentModule(); err == nil {
		// native functions only need a name to reference them
		
		mod.AddFunction(addI32)
		mod.AddFunction(double)
		mod.AddFunction(quad)
		mod.AddFunction(sumQuads)
		mod.AddFunction(main)
	} else {
		fmt.Println(err)
	}
	
	// Generating random program
	//cxt := RandomProgram(200)
	//cxt.PrintProgram(false)

	// Call stack testing
	cxt.Run(false, 2000)
	fmt.Println(len(cxt.Steps))
	fmt.Println("...")
	//PrintCallStack(cxt.Steps[len(cxt.Steps) - 1])
	PrintCallStack(cxt.CallStack)
	//cxt.Run(false, 200)
	fmt.Println(len(cxt.Steps))

	// for _, step := range cxt.Steps {
	// 	//fmt.Println(step)
	// 	//PrintCallStack(step)
	// 	fmt.Printf("%p\n", step[len(step) - 1].Context)
	// }
	//PrintCallStack(cxt.CallStack)	


	fmt.Println("...")

	
	// copy of program at about the middle of its execution
	cxtCopy := MakeContextCopy(cxt, 1000)
	// if defs, err := cxtCopy.GetCurrentDefinitions(); err == nil {
	// 	defs["num1"].Value = &num2
	// }
	//cxt.Run(false, 200)
	// fmt.Println("...")

	//PrintCallStack(cxtCopy.Steps[len(cxtCopy.Steps) - 1])
	PrintCallStack(cxtCopy.CallStack)
	//fmt.Println(len(cxtCopy.Steps))
	//fmt.Println(len(cxt.Steps))
	cxtCopy.Run(false, 2000)
	//fmt.Println(len(cxtCopy.Steps))
	// fmt.Println("...")
	fmt.Println(len(cxt.Steps))
	fmt.Println(len(cxtCopy.Steps))
	// PrintCallStack(cxtCopy.Steps[len(cxtCopy.Steps) - 1])
	fmt.Println("checking calls contexts")
	// for _, step := range cxtCopy.Steps {
	// 	fmt.Printf("%p\n", step[len(step) - 1].Context)
	// }

	for i := 0; i < len(cxtCopy.Steps); i++ {
		fmt.Printf("...%d\n", i)
		//PrintCallStack(cxt.Steps[i])
		PrintCallStack(cxtCopy.Steps[i])
	}
}
