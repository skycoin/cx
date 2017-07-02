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
	
	// native functions only need a name to reference them
	addI32 := MakeFunction("addI32")
	
	// function which uses native function "addI32"
	double := MakeFunction("double").
		  AddInput(MakeParameter("num", i32)).
		  AddOutput(MakeParameter("", i32)).
		
		  AddExpression(MakeExpression(addI32).
		                AddArgument(MakeArgument(MakeValue("num"), ident)).
				AddArgument(MakeArgument(MakeValue("num"), ident)))


	
	// double.GetAffordances() => Should it print


	// failed, a struct needs to have basic type fields
	////// fmt.Println(encoder.Serialize(*MakeExpression(addI32)))

	num1 := encoder.SerializeAtomic(int32(25))
	num2 := encoder.SerializeAtomic(int32(40))
	inum1 := []byte("num1")
	//inum2 := []byte("num2")

	mod := MakeModule("main").
		AddDefinition(MakeDefinition("num1", &num1, i32)).
		AddDefinition(MakeDefinition("num2", &num2, i32)).
		AddFunction(double)

	
	fmt.Println(
		MakeExpression(addI32).
			AddArgument(MakeArgument(&inum1, ident)).
			AddArgument(MakeArgument(&num2, i32)).
			
			Execute(mod.Definitions).Value)

	fmt.Println(
		MakeExpression(double).
			AddArgument(MakeArgument(&num1, i32)).
			
			Execute(mod.Definitions).Value)
	
	fmt.Println(MakeContext().AddModule(MakeModule("Math")).AddModule(MakeModule("StdLib")))

	fmt.Println("\nTesting Selectors")

	fmt.Println(MakeContext().AddModule(MakeModule("Math")).AddModule(MakeModule("StdLib")).
		CurrentModule)
	fmt.Println(MakeContext().AddModule(MakeModule("Math")).AddModule(MakeModule("StdLib")).
		SelectModule("Math"))

	

	
	fmt.Println("\nTesting Affordances")
	fmt.Println("\n---Context Affordances")
	
	cxt := MakeContext()
	
	affs := cxt.GetAffordances()
	
	PrintAffordances(affs)
	
	fmt.Println(cxt.Modules)
	affs[0].ApplyAffordance() // This should add a module

	fmt.Println(cxt.Modules[0].Name)
	PrintAffordances(cxt.GetAffordances())

	fmt.Println("\n---Module Affordances")
	
	cxt.AddModule(MakeModule("Math"))
	imp := MakeModule("StdLib")
	cxt.AddModule(imp)
	//cxt.AddModule(MakeModule("Bytes"))
	//cxt.AddModule(MakeModule("HTTP"))
	affs = cxt.SelectModule("Math").
		AddStruct(MakeStruct("List")).
		AddStruct(MakeStruct("Complex")).
		AddImport(imp).
		GetAffordances()
	
	PrintAffordances(affs)
	fmt.Println(cxt.GetCurrentModule().Definitions)
	affs[0].ApplyAffordance()
	affs[1].ApplyAffordance()
	fmt.Println(cxt.GetCurrentModule().Definitions["def2"].Typ)

	fmt.Println("Functions before:")
	fmt.Println(cxt.GetCurrentModule().Functions)
	fmt.Println("Functions after:")
	affs[len(affs) - 2].ApplyAffordance() // Adding a function
	affs[len(affs) - 2].ApplyAffordance() // This would redefine the last function

	affs = cxt.GetCurrentModule().GetAffordances()
	affs[len(affs) - 2].ApplyAffordance()
	cxt.SelectModule("StdLib").
		AddStruct(MakeStruct("Stream"))
	affs = cxt.GetCurrentModule().GetAffordances()
	affs[len(affs) - 2].ApplyAffordance() // Adding a function
	//affs = cxt.GetCurrentFunction().
	cxt.SelectModule("Math")

	fmt.Println(cxt.GetCurrentModule().Functions)
	fmt.Println(cxt.GetCurrentFunction())

	fmt.Println("\n---Function Affordances")
	affs = cxt.SelectFunction("fn5").GetAffordances()
	PrintAffordances(affs)
	affs[0].ApplyAffordance() // Adds an input
	affs[6].ApplyAffordance()
	// fmt.Println(cxt.GetCurrentFunction().Expressions[0].Operator.Module)
	// fmt.Println(cxt.GetCurrentFunction().Expressions[0].Operator.Name)
	// fmt.Println(cxt.GetCurrentFunction().Expressions[0].Arguments)

	// cxt.GetCurrentModule().AddDefinition(MakeDefinition("num1", nil, i32))
	// cxt.GetCurrentModule().AddDefinition(MakeDefinition("num2", nil, i32))
	// cxt.GetCurrentModule().AddDefinition(MakeDefinition("sum", nil, i32))
	cxt.GetCurrentModule().AddDefinition(MakeDefinition("hugs", nil, i32))
	cxt.GetCurrentModule().AddDefinition(MakeDefinition("items", nil, MakeType("StdLib.Stream")))
	cxt.GetCurrentModule().AddDefinition(MakeDefinition("views", nil, i64))

	cxt.GetCurrentFunction().AddInput(MakeParameter("num1", i32))
	cxt.GetCurrentFunction().AddInput(MakeParameter("num2", i32))
	cxt.GetCurrentFunction().AddOutput(MakeParameter("sum", i32))

	cxt.SelectFunction("fn10")
	cxt.GetCurrentFunction().AddInput(MakeParameter("items", MakeType("StdLib.Stream")))
	FilterAffordances(cxt.GetCurrentFunction().GetAffordances(),
		"AddExpression")[0].ApplyAffordance()
	
	
	expr, err := cxt.GetCurrentFunction().SelectExpression(100)
	if err == nil {
		affs = expr.GetAffordances()
		FilterAffordances(affs, "AddArgument")[0].ApplyAffordance()
		affs = expr.GetAffordances()
		FilterAffordances(affs, "AddArgument")[0].ApplyAffordance()
		affs = expr.GetAffordances()
		FilterAffordances(affs, "AddArgument")[0].ApplyAffordance()
		affs = expr.GetAffordances()
		FilterAffordances(affs, "AddArgument")[0].ApplyAffordance()
	} else {
		fmt.Println(err)
	}

	FilterAffordances(cxt.SelectStruct("List").GetAffordances(),
		"AddField")[0].ApplyAffordance()
	FilterAffordances(cxt.GetCurrentStruct().GetAffordances(),
		"AddField")[1].ApplyAffordance()

	fmt.Println("\n...\n")
	//cxt.PrintProgram(true)

	fmt.Println("\nNew Program\n")

	// Generating random program
	cxt = RandomProgram(200)
	cxt.PrintProgram(false)
}
