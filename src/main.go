// Tasks
// Create a cxStruct             x
// Create a Module Struct        x
// Create a Function struct      x
// Create a DataType struct      x
// Create a Context struct       x
// Create an Expression struct   x
// Create basic functions        

// Need: Mechanism to execute functions/expressions               x
// Need: Mechanism to constraint data types
// Need: Mechanism for explicitly casting types                   

// Need: Mechanism for reflection                                 
// Need: Mechanism to determine "affordances" and "restrictions"  
// Need: Start with a null object and call affordances            

// Need: Mechanism to serialize structs                           
// Need: Mechanism to hash and pull data from program             

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

//type bar func()

// type affordance struct {
// 	Typ cxType
// 	Fn []byte //hmmmm
// }

func foo() func() {
	return func () {
		fmt.Println("hello")
	}
}

func main() {
	i32 := MakeType("i32")
	ident := MakeType("ident")
	
	// native functions only need a name to reference them
	addI32 := MakeFunction("addI32")

	num := []byte("num")
	
	// function which uses native function "addI32"
	double := MakeFunction("double").
		AddInput(MakeParameter("num", i32)).
		AddOutput(MakeParameter("", i32)).
		
		AddExpression(MakeExpression(addI32).
		AddArgument(MakeArgument(&num, ident)).
		AddArgument(MakeArgument(&num, ident))) // => (+ num num)


	// // Start of new way

	// // With the new way, it would be like this:

	// double := MakeFunction().AddName(MakeName("hello")).
	// 	AddInput(MakeParameter().AddName(MakeName("num")).AddType(MakeType("i32"))).
	// 	AddOutput(Ma)

	// double := MakeFunction().AddName(MakeName("hello")).
	// 	AddInput().AddName().
	// 	AddOutput(Ma)

	// double := MakeFunction().AddName(MakeName("hello")).
	// 	AddInput().AddInput().
	// 	AddOutput(Ma)


	// // We construct the program over a null object
	// context.AddFunction().AddName().AddInput().AddType().AddInput().AddType().AddExpression().AddFunction()
	// 	// Then we add the leaves, which are going to be applied to the first structure that makes sense OR see the other way which is compatible with this model
	// 	MakeName("sum").MakeName("num1").MakeName("num2")
	// //The makers are for: fn name, input1 name, input2 name

	// // This way is totally compatible because the makers are applied to the first non-null field that makes sense
	// context.AddFunction().MakeName("sum").AddInput().MakeName("num1").AddType()

	// // We can determine the object over we will apply the adder by using getters internally (the programmer can also use them)
	// // For example, we can do the following
	// context.AddFunction().AddInput().MakeName("num1")
	// // This MakeName would use a (cxt *cxContext) GetUnnamed() and add "num1" to the first unnamed "nameable" object
	
	// // In the case of MakeValue, we would use a GetUnvalued() method and it would add the value to the first unvalued "valuable" object

	// // There will only be MakeName and MakeValue.
	// // The adders will provide the logic of the program, and the programmer only needs to provide names for the definitions and what values they are going to hold

	// // In the end, we are only providing a program structure and values (names and values for these definitions; and we could actually omit the names)
	
	// // This would throw an error, because there is no appropriate object to apply the function_on
	// context.AddInput()


	// // The functions_of can give us, the programmers, what are the possible actions and what's the structure of the program
	// // For now, the most interesting part are the functions_on






	
	// double.GetAffordances() => Should it print


	// failed, a struct needs to have basic type fields
	////// fmt.Println(encoder.Serialize(*MakeExpression(addI32)))
	
	fmt.Println("...\n")

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
			AddArgument(MakeArgument(&inum1, ident)).
			
			Execute(mod.Definitions).Value)

	foo()()

	fmt.Println(MakeFunction("something").GetAffordances())

	var myfunc = foo()
	
	myfunc()
	
	fmt.Println("\n...")
}
