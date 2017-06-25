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

	fmt.Println("\n...")
}
