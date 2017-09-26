package base

import (
	"bytes"
	"io/ioutil"
	
	"fmt"
	"errors"
	"strings"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func argsToDefs (args []*CXArgument, names []string, mod *CXModule, cxt *CXProgram) ([]*CXDefinition, error) {
	if len(names) == len(args) {
		defs := make([]*CXDefinition, 0)
		for i, arg := range args {
			defs = append(defs, &CXDefinition{
				Name: names[i],
				Typ: arg.Typ,
				Value: arg.Value,
				Module: mod,
				Context: cxt,
			})
			// defs[names[i]] = &CXDefinition{
			// 	Name: names[i],
			// 	Typ: arg.Typ,
			// 	Value: arg.Value,
			// 	Module: mod,
			// 	Context: cxt,
			// }
		}
		return defs, nil
	} else {
		return nil, errors.New("Not enough definition names provided")
	}
}

func PrintCallStack (callStack []*CXCall) {
	for i, call := range callStack {
		tabs := strings.Repeat("___", i)
		if tabs == "" {
			fmt.Printf("%sfn:%s ln:%d, \tlocals: ", tabs, call.Operator.Name, call.Line)
		} else {
			fmt.Printf("↓%sfn:%s ln:%d, \tlocals: ", tabs, call.Operator.Name, call.Line)
		}

		lenState := len(call.State)
		idx := 0
		for _, def := range call.State {
			if def.Name == "_" || (len(def.Name) > len(NON_ASSIGN_PREFIX) && def.Name[:len(NON_ASSIGN_PREFIX)] == NON_ASSIGN_PREFIX) {
				continue
			}
			var valI32 int32
			var valI64 int64
			var valF32 float32
			var valF64 float64
			switch def.Typ.Name {
			case "i32":
				encoder.DeserializeRaw(*def.Value, &valI32)
				if idx == lenState - 1 {
					fmt.Printf("%s: %d", def.Name, valI32)
				} else {
					fmt.Printf("%s: %d, ", def.Name, valI32)
				}
			case "i64":
				encoder.DeserializeRaw(*def.Value, &valI64)
				if idx == lenState - 1 {
					fmt.Printf("%s: %d", def.Name, valI64)
				} else {
					fmt.Printf("%s: %d, ", def.Name, valI64)
				}
			case "f32":
				encoder.DeserializeRaw(*def.Value, &valF32)
				if idx == lenState - 1 {
					fmt.Printf("%s: %f", def.Name, valF32)
				} else {
					fmt.Printf("%s: %f, ", def.Name, valF32)
				}
			case "f64":
				encoder.DeserializeRaw(*def.Value, &valF64)
				if idx == lenState - 1 {
					fmt.Printf("%s: %f", def.Name, valF64)
				} else {
					fmt.Printf("%s: %f, ", def.Name, valF64)
				}
			case "byte":
				if idx == lenState - 1 {
					fmt.Printf("%s: %d", def.Name, (*def.Value)[0])
				} else {
					fmt.Printf("%s: %d, ", def.Name, (*def.Value)[0])
				}
			case "[]byte":
				var val []byte
				encoder.DeserializeRaw(*def.Value, &val)
				if idx == lenState - 1 {
					fmt.Printf("%s: %v", def.Name, val)
				} else {
					fmt.Printf("%s: %v, ", def.Name, val)
				}
			case "[]i32":
				var val []int32
				encoder.DeserializeRaw(*def.Value, &val)
				if idx == lenState - 1 {
					fmt.Printf("%s: %v", def.Name, val)
				} else {
					fmt.Printf("%s: %v, ", def.Name, val)
				}
			case "[]i64":
				var val []int64
				encoder.DeserializeRaw(*def.Value, &val)
				if idx == lenState - 1 {
					fmt.Printf("%s: %v", def.Name, val)
				} else {
					fmt.Printf("%s: %v, ", def.Name, val)
				}
			case "[]f32":
				var val []float32
				encoder.DeserializeRaw(*def.Value, &val)
				if idx == lenState - 1 {
					fmt.Printf("%s: %v", def.Name, val)
				} else {
					fmt.Printf("%s: %v, ", def.Name, val)
				}
			case "[]f64":
				var val []float64
				encoder.DeserializeRaw(*def.Value, &val)
				if idx == lenState - 1 {
					fmt.Printf("%s: %v", def.Name, val)
				} else {
					fmt.Printf("%s: %v, ", def.Name, val)
				}
			}
			
			idx++
		}
		fmt.Println()
	}
}

func callsEqual (call1, call2 *CXCall) bool {
	if call1.Line != call2.Line ||
		len(call1.State) != len(call2.State) ||
		call1.Operator != call2.Operator ||
		call1.ReturnAddress != call2.ReturnAddress ||
		call1.Module != call2.Module {
		return false
	}

	for k, v := range call1.State {
		if call2.State[k] != v {
			return false
		}
	}

	return true
}

func saveStep (call *CXCall) {
	lenCallStack := len(call.Context.CallStack.Calls)
	newStep := MakeCallStack(lenCallStack)

	if len(call.Context.Steps) < 1 {
		// First call, copy everything
		for i, call := range call.Context.CallStack.Calls {
			newStep.Calls[i] = MakeCallCopy(call, call.Module, call.Context)
		}

		call.Context.Steps = append(call.Context.Steps, newStep)
		return
	}
	
	lastStep := call.Context.Steps[len(call.Context.Steps) - 1]
	lenLastStep := len(lastStep.Calls)
	
	smallerLen := 0
	if lenLastStep < lenCallStack {
		smallerLen = lenLastStep
	} else {
		smallerLen = lenCallStack
	}
	
	// Everytime a call changes, we need to make a hard copy of it
	// If the call doesn't change, we keep saving a pointer to it

	for i, call := range call.Context.CallStack.Calls[:smallerLen] {
		if callsEqual(call, lastStep.Calls[i]) {
			// if they are equal
			// append reference
			newStep.Calls[i] = lastStep.Calls[i]
		} else {
			newStep.Calls[i] = MakeCallCopy(call, call.Module, call.Context)
		}
	}

	// sizes can be different. if this is the case, we hard copy the rest
	for i, call := range call.Context.CallStack.Calls[smallerLen:] {
		newStep.Calls[i + smallerLen] = MakeCallCopy(call, call.Module, call.Context)
	}
	
	call.Context.Steps = append(call.Context.Steps, newStep)
	return
}

// It "un-runs" a program
func (cxt *CXProgram) Reset() {
	cxt.CallStack = MakeCallStack(0)
	cxt.Steps = make([]*CXCallStack, 0)
	//cxt.ProgramSteps = nil
}

func (cxt *CXProgram) ResetTo(stepNumber int) {
	// if no steps, we do nothing. the program will run from step 0
	if len(cxt.Steps) > 0 {
		if stepNumber > len(cxt.Steps) {
			stepNumber = len(cxt.Steps) - 1
		}
		reqStep := cxt.Steps[stepNumber]

		newStep := MakeCallStack(len(reqStep.Calls))
		
		var lastCall *CXCall
		for j, call := range reqStep.Calls {
			newCall := MakeCallCopy(call, call.Module, call.Context)
			newCall.ReturnAddress = lastCall
			lastCall = newCall
			newStep.Calls[j] = newCall
		}

		cxt.CallStack = newStep
		cxt.Steps = cxt.Steps[:stepNumber]
	}
}

func (cxt *CXProgram) UnRun (nCalls int) {
	if len(cxt.Steps) > 0 && nCalls > 0 {
		if nCalls > len(cxt.Steps) {
			nCalls = len(cxt.Steps) - 1
		}

		reqStep := cxt.Steps[len(cxt.Steps) - nCalls]

		newStep := MakeCallStack(len(reqStep.Calls))
		
		var lastCall *CXCall
		for j, call := range reqStep.Calls {
			newCall := MakeCallCopy(call, call.Module, call.Context)
			newCall.ReturnAddress = lastCall
			lastCall = newCall
			newStep.Calls[j] = newCall
		}

		cxt.CallStack = newStep
		cxt.Steps = cxt.Steps[:len(cxt.Steps) - nCalls]
	}
}

func replPrintEvaluation (arg *CXArgument) {
	fmt.Printf(">> ")
	switch arg.Typ.Name {
	case "str":
		fmt.Printf("%#v\n", string(*arg.Value))
	case "bool":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		if val == 0 {
			fmt.Printf("false\n")
		} else {
			fmt.Printf("true\n")
		}
	case "byte":
		fmt.Printf("%#v\n", *arg.Value)
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "[]byte":
		var val []byte
		encoder.DeserializeRaw(*arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)
		fmt.Printf("%#v\n", val)
	default:
		fmt.Printf("")
	}
}

func beforeDot (str string ) (beforeDot string, afterDot string) {
	beforeDot = ""
	afterDot = ""
	beforeDotCounter := 0
	foundDot := false


	for i, letter := range str {
		if letter == '.' {
			foundDot = true
			beforeDot = str[:i]
			beforeDotCounter = i
			break
		}
	}

	if !foundDot {
		beforeDot = str
	}
	
	if foundDot {
		afterDot = str[beforeDotCounter + 1:] // ignore the dot
	}

	return beforeDot, afterDot
}

func GetIdentParts (str string) []string {
	var identParts []string
	before, after := beforeDot(str) // => "Math", "myPoint.x"
	identParts = append(identParts, before)
	for after != "" {
		before, after = beforeDot(after)
		identParts = append(identParts, before)
	}
	return identParts
}

func getAllocSize (typ string) int {
	// I need to do arrays too
	switch typ {
	case "byte":
		return 1
	case "i32":
		return 4
	case "i64":
		return 8
	case "f32":
		return 4
	case "f64":
		return 8
	default:
		return -1
	}
}

// Compiling from CXGO to CX Base
func (cxt *CXProgram) Compile (withProfiling bool) {
	var asmNL string = "\n"
	var program bytes.Buffer

	if withProfiling {
		program.WriteString(`package main;

import (
. "github.com/skycoin/cx/src/base"
"os"
	"log"
	"flag"
	"runtime/pprof"
);

var cxt = MakeContext();var mod *CXModule;var imp *CXModule;var fn *CXFunction;var op *CXFunction;var expr *CXExpression;var strct *CXStruct;var arg *CXArgument;var tag string = "";

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile 'file'")

func main () {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
`)
	} else {
		program.WriteString(`package main;import (. "github.com/skycoin/cx/src/base";);var cxt = MakeContext();var mod *CXModule;var imp *CXModule;var fn *CXFunction;var op *CXFunction;var expr *CXExpression;var strct *CXStruct;var arg *CXArgument;var tag string = "";func main () {`)
	}
	
	

	for _, mod := range cxt.Modules {
		program.WriteString(fmt.Sprintf(`mod = MakeModule("%s");cxt.AddModule(mod);%s`, mod.Name, asmNL))
		for _, imp := range mod.Imports {
			program.WriteString(fmt.Sprintf(`imp, _ = cxt.GetModule("%s");mod.AddImport(imp);%s`, imp.Name, asmNL))
		}

		for _, fn := range mod.Functions {
			program.WriteString(fmt.Sprintf(`fn = MakeFunction(%#v);mod.AddFunction(fn);%s`, fn.Name, asmNL))

			for _, inp := range fn.Inputs {
				program.WriteString(fmt.Sprintf(`fn.AddInput(MakeParameter("%s", MakeType("%s")));%s`, inp.Name, inp.Typ.Name, asmNL))
			}

			for _, out := range fn.Outputs {
				program.WriteString(fmt.Sprintf(`fn.AddOutput(MakeParameter("%s", MakeType("%s")));%s`, out.Name, out.Typ.Name, asmNL))
			}

			for _, expr := range fn.Expressions {
				var tagStr string
				if expr.Tag != "" {
					tagStr = fmt.Sprintf(`expr.Tag = tag;tag = "";`)
				}
				program.WriteString(fmt.Sprintf(`op, _ = cxt.GetFunction("%s", "%s");expr = MakeExpression(op);expr.FileLine = %d;fn.AddExpression(expr);%s%s`,
					expr.Operator.Name, expr.Operator.Module.Name, expr.FileLine, tagStr, asmNL))

				for _, arg := range expr.Arguments {
					program.WriteString(fmt.Sprintf(`expr.AddArgument(MakeArgument(&%#v, MakeType("%s")));%s`, *arg.Value, arg.Typ.Name, asmNL))
				}
				
				for _, outName := range expr.OutputNames {
					program.WriteString(fmt.Sprintf(`expr.AddOutputName("%s");%s`, outName.Name, asmNL))
				}
			}
		}

		for _, strct := range mod.Structs {
			program.WriteString(fmt.Sprintf(`strct = MakeStruct("%s");mod.AddStruct(strct);%s`, strct.Name, asmNL))
		}

		for _, def := range mod.Definitions {
			program.WriteString(fmt.Sprintf(`mod.AddDefinition(MakeDefinition("%s", &%#v, MakeType("%s")));%s`, def.Name, *def.Value, def.Typ.Name, asmNL))
		}
	}

	program.WriteString(`cxt.Run(false, -1);}`)
	ioutil.WriteFile(fmt.Sprintf("o.go"), []byte(program.String()), 0644)
}

// func (cxt *CXProgram) Compile () *CXProgram {
// 	allocs := make(map[string]*CXArgument, 0)
// 	heap := *cxt.Heap
// 	heapCounter := 0

// 	for _, mod := range cxt.Modules {
// 		for _, fn := range mod.Functions {

// 			// using struct fields
// 			for _, expr := range fn.Expressions {


// 				for i, outName := range expr.OutputNames {
// 					allocSize := getAllocSize(expr.Operator.Outputs[i].Typ.Name)
// 					outName.Offset = heapCounter
// 					heapCounter = heapCounter + allocSize
// 				}
				
				
// 				for argIdx, arg := range expr.Arguments {
// 					if arg.Typ.Name == "ident" {
// 						identParts := GetIdentParts(string(*arg.Value))

// 						//var def *CXDefinition
// 						if len(identParts) == 1 {
// 							// it's a current module's definition
// 							// TODO: compile later
// 						} else {
// 							// is the first part of the identifier a module, if yes:
// 							if identMod, err := cxt.GetModule(identParts[0]); err == nil {
// 								// identParts[1] will always be a definition if before is a module
// 								if len(identParts) == 2 {
// 									// then we are referring to the struct itself or a normal ident
// 									// TODO: compile later
// 								} else {
// 									// we're referring to a struct field then
// 									if len(identParts) > 3 {
// 										// nested structs
// 										// TODO: compile later
// 									} else {
// 										if def, err := identMod.GetDefinition(concat(identParts[1], ".", identParts[2])); err == nil {
// 											var heapArg *CXArgument
// 											defSize := len(*def.Value)
// 											if found, ok := allocs[concat(identParts[0], ".", identParts[1], ".", identParts[2])]; ok {
// 												heapArg = found
// 											} else {
// 												heapArg = &CXArgument{
// 													Typ: def.Typ,
// 													//Value: arg.Value, //irrelevant now
// 													Offset: heapCounter,
// 													Size: defSize,
// 												}
// 												allocs[concat(identParts[0], ".", identParts[1], ".", identParts[2])] = heapArg
// 											}
											
// 											// replacing argument
// 											expr.Arguments[argIdx] = heapArg

// 											heap = append(heap, *def.Value...)
// 											heapCounter = heapCounter + defSize
// 										} else {
// 											// this means it's a local definition (or it doesn't exist, but Run() takes care of this)
											
// 										}
// 									}
// 								}
// 							} else {
// 								if len(identParts) > 2 {
// 									fmt.Println("identParts > 2")
// 									// nested structs
// 									// TODO: compile later
// 								} else {
// 									if def, err := mod.GetDefinition(concat(identParts[0], ".", identParts[1])); err == nil {
// 										//def = mod.Definitions[concat(identParts[0], ".", identParts[1])]
// 										var heapArg *CXArgument
// 										defSize := len(*def.Value)
// 										if found, ok := allocs[concat(mod.Name, ".", identParts[0], ".", identParts[1])]; ok {
// 											heapArg = found
// 										} else {
// 											heapArg = &CXArgument{
// 												Typ: def.Typ,
// 												Offset: heapCounter,
// 												Size: defSize,
// 											}
// 											// mod.Name identParts[0]
// 											allocs[concat(mod.Name, ".", identParts[0], ".", identParts[1])] = heapArg
// 										}
										
// 										// replacing argument
// 										expr.Arguments[argIdx] = heapArg

// 										heap = append(heap, *def.Value...)
// 										heapCounter = heapCounter + defSize
// 									}
// 								}
// 							}
// 						}

// 						// checking if first part is a module
						
// 					}
// 				}
// 			}
// 		}
// 	}
// 	cxt.Heap = &heap
// 	return cxt
// }


//var biggestOutputNumber int = 0

func (cxt *CXProgram) Run (withDebug bool, nCalls int) error {
	if cxt.Terminated {
		// user wants to re-run the program
		cxt.Terminated = false
	}

	var callCounter int = 0
	// we are going to do this if the CallStack is empty
	if cxt.CallStack != nil && len(cxt.CallStack.Calls) > 0 {
		// we resume the program
		var lastCall *CXCall
		var err error

		var untilEnd = false
		if nCalls < 1 {
			nCalls = 1 // so the for loop executes
			untilEnd = true
		}

		for !cxt.Terminated && nCalls > 0 {
			lastCall = cxt.CallStack.Calls[len(cxt.CallStack.Calls) - 1]
			err = lastCall.call(withDebug, 1, callCounter)
			if err != nil {
				return err
			}
			if !untilEnd {
				nCalls = nCalls - 1
			}
		}
	} else {
		// initialization and checking
		if mod, err := cxt.SelectModule("main"); err == nil {
			if fn, err := mod.SelectFunction("main"); err == nil {
				// main function
				state := make([]*CXDefinition, 0)
				mainCall := MakeCall(fn, state, nil, mod, mod.Context)
				
				cxt.CallStack.Calls = append(cxt.CallStack.Calls, mainCall)

				//return mainCall.call(withDebug, nCalls, callCounter)

				var lastCall *CXCall
				var err error

				var untilEnd = false
				if nCalls < 1 {
					nCalls = 1 // so the for loop executes
					untilEnd = true
				}
				
				for !cxt.Terminated && nCalls > 0 {
					lastCall = cxt.CallStack.Calls[len(cxt.CallStack.Calls) - 1]
					err = lastCall.call(withDebug, 1, callCounter)
					if err != nil {
						return err
					}
					if !untilEnd {
						nCalls = nCalls - 1
					}
				}
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func checkNative (opName string, expr *CXExpression, call *CXCall, values, argsCopy *[]*CXArgument, exc *bool, excError *error) {
	switch opName {
	case "serialize":
		sProgram := Serialize(call.Context)
		*values = append(*values, MakeArgument(sProgram, MakeType("[]byte")))
	case "deserialize":
		// it only prints the deserialized program for now
		Deserialize((*argsCopy)[0].Value).PrintProgram(false)
		*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
	case "evolve":
		fnName := string(*(*argsCopy)[0].Value)
		fnBag := string(*(*argsCopy)[1].Value)
		
		var inps []float64
		encoder.DeserializeRaw(*(*argsCopy)[2].Value, &inps)
		
		var outs []float64
		encoder.DeserializeRaw(*(*argsCopy)[3].Value, &outs)

		var numberExprs int32
		encoder.DeserializeRaw(*(*argsCopy)[4].Value, &numberExprs)
		var iterations int32
		encoder.DeserializeRaw(*(*argsCopy)[5].Value, &iterations)
		var epsilon float64
		encoder.DeserializeRaw(*(*argsCopy)[6].Value, &epsilon)

		if evolutionErr, err := call.Context.Evolve(fnName, fnBag, inps, outs, int(numberExprs), int(iterations), epsilon); err == nil {
			val := encoder.Serialize(evolutionErr)
			*values = append(*values, MakeArgument(&val, MakeType("f64")))
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("f64"), MakeType("f64")))
		}
		// flow control
	case "baseGoTo":
		if val, err := baseGoTo(call, (*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "goTo":
		if val, err := goTo(call, (*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
		// I/O functions
	case "bool.print":
		var val int32
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		if val == 0 {
			fmt.Println("false")
		} else {
			fmt.Println("true")
		}
		*values = append(*values, (*argsCopy)[0])
	case "str.print":
		fmt.Println(string(*(*argsCopy)[0].Value))
		*values = append(*values, (*argsCopy)[0])
	case "byte.print":
		fmt.Println((*(*argsCopy)[0].Value)[0])
		*values = append(*values, (*argsCopy)[0])
	case "i32.print":
		var val int32
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
		*values = append(*values, (*argsCopy)[0])
	case "i64.print":
		var val int64
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
		*values = append(*values, (*argsCopy)[0])
	case "f32.print":
		var val float32
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
		*values = append(*values, (*argsCopy)[0])
	case "f64.print":
		var val float64
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
		*values = append(*values, (*argsCopy)[0])
	case "[]bool.print":
		var val []int32
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Print("[")
		for i, v := range val {
			if v == 0 {
				fmt.Print("false")
			} else {
				fmt.Print("true")
			}
			if i != len(val) -1 {
				fmt.Print(" ")
			}
		}
		fmt.Print("]")
		fmt.Println()
		*values = append(*values, (*argsCopy)[0])
	case "[]byte.print":
		fmt.Println(*(*argsCopy)[0].Value)
		*values = append(*values, (*argsCopy)[0])
	case "[]i32.print":
		var val []int32
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
		*values = append(*values, (*argsCopy)[0])
	case "[]i64.print":
		var val []int64
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
		*values = append(*values, (*argsCopy)[0])
	case "[]f32.print":
		var val []float32
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
		*values = append(*values, (*argsCopy)[0])
	case "[]f64.print":
		var val []float64
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
		*values = append(*values, (*argsCopy)[0])
		// identity functions
	case "str.id", "bool.id", "byte.id", "i32.id", "i64.id", "f32.id", "f64.id", "[]bool.id", "[]byte.id", "[]i32.id", "[]i64.id", "[]f32.id", "[]f64.id":
		*values = append(*values, (*argsCopy)[0])
		// cast functions
	case "[]byte.str":
		if val, err := castToStr((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("str"), MakeType("str")))
		}
	case "str.[]byte":
		if val, err := castToByteA((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]byte"), MakeType("[]byte")))
		}
	case "i32.byte", "i64.byte", "f32.byte", "f64.byte":
		if val, err := castToByte((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("byte"), MakeType("byte")))
		}
	case "byte.i32", "i64.i32", "f32.i32", "f64.i32":
		if val, err := castToI32((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "byte.i64", "i32.i64", "f32.i64", "f64.i64":
		if val, err := castToI64((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i64"), MakeType("i64")))
		}
	case "byte.f32", "i32.f32", "i64.f32", "f64.f32":
		if val, err := castToF32((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("f32"), MakeType("f32")))
		}
	case "byte.f64", "i32.f64", "i64.f64", "f32.f64":
		if val, err := castToF64((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("f64"), MakeType("f64")))
		}
	case "[]i32.[]byte", "[]i64.[]byte", "[]f32.[]byte", "[]f64.[]byte":
		if val, err := castToByteA((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]byte"), MakeType("[]byte")))
		}
	case "[]byte.[]i32", "[]i64.[]i32", "[]f32.[]i32", "[]f64.[]i32":
		if val, err := castToI32A((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]i32"), MakeType("[]i32")))
		}
	case "[]byte.[]i64", "[]i32.[]i64", "[]f32.[]i64", "[]f64.[]i64":
		if val, err := castToI64A((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]i64"), MakeType("[]i64")))
		}
	case "[]byte.[]f32", "[]i32.[]f32", "[]i64.[]f32", "[]f64.[]f32":
		if val, err := castToF32A((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]f32"), MakeType("[]f32")))
		}
	case "[]byte.[]f64", "[]i32.[]f64", "[]i64.[]f64", "[]f32.[]f64":
		if val, err := castToF64A((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]f64"), MakeType("[]f64")))
		}
		// logical operators
	case "and":
		if val, err := and((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "or":
		if val, err := or((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "not":
		if val, err := not((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
		// relational operators
	case "i32.lt":
		if val, err := ltI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "i32.gt":
		if val, err := gtI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "i32.eq":
		if val, err := eqI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "i32.lteq":
		if val, err := lteqI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "i32.gteq":
		if val, err := gteqI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "i64.lt":
		if val, err := ltI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "i64.gt":
		if val, err := gtI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "i64.eq":
		if val, err := eqI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "i64.lteq":
		if val, err := lteqI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "i64.gteq":
		if val, err := gteqI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "f32.lt":
		if val, err := ltF32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "f32.gt":
		if val, err := gtF32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "f32.eq":
		if val, err := eqF32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "f32.lteq":
		if val, err := lteqF32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "f32.gteq":
		if val, err := gteqF32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "f64.lt":
		if val, err := ltF64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "f64.gt":
		if val, err := gtF64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "f64.eq":
		if val, err := eqF64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "f64.lteq":
		if val, err := lteqF64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "f64.gteq":
		if val, err := gteqF64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "str.lt":
		if val, err := ltStr((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "str.eq":
		if val, err := eqStr((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "str.lteq":
		if val, err := lteqStr((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "str.gteq":
		if val, err := gteqStr((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "byte.lt":
		if val, err := ltByte((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "byte.gt":
		if val, err := gtByte((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "byte.eq":
		if val, err := eqByte((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "byte.lteq":
		if val, err := lteqByte((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "byte.gteq":
		if val, err := gteqByte((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
		// struct operations
	case "initDef":
		if val, err := initDef((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
		// arithmetic functions
	case "i32.add":
		if val, err := addI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "i32.mul":
		if val, err := mulI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "i32.sub":
		if val, err := subI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "i32.div":
		if val, err := divI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "i64.add":
		if val, err := addI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i64"), MakeType("i64")))
		}
	case "i64.mul":
		if val, err := mulI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i64"), MakeType("i64")))
		}
	case "i64.sub":
		if val, err := subI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i64"), MakeType("i64")))
		}
	case "i64.div":
		if val, err := divI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i64"), MakeType("i64")))
		}
	case "f32.add":
		if val, err := addF32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("f32"), MakeType("f32")))
		}
	case "f32.mul":
		if val, err := mulF32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("f32"), MakeType("f32")))
		}
	case "f32.sub":
		if val, err := subF32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("f32"), MakeType("f32")))
		}
	case "f32.div":
		if val, err := divF32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("f32"), MakeType("f32")))
		}
	case "f64.add":
		if val, err := addF64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("f64"), MakeType("f64")))
		}
	case "f64.mul":
		if val, err := mulF64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("f64"), MakeType("f64")))
		}
	case "f64.sub":
		if val, err := subF64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("f64"), MakeType("f64")))
		}
	case "f64.div":
		if val, err := divF64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("f64"), MakeType("f64")))
		}
	case "i32.mod":
		if val, err := modI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "i64.mod":
		if val, err := modI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i64"), MakeType("i64")))
		}
		// bitwise operators
	case "i32.and":
		if val, err := andI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "i32.or":
		if val, err := orI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "i32.xor":
		if val, err := xorI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "i32.andNot":
		if val, err := andNotI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "i64.and":
		if val, err := andI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i64"), MakeType("i64")))
		}
	case "i64.or":
		if val, err := orI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i64"), MakeType("i64")))
		}
	case "i64.xor":
		if val, err := xorI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i64"), MakeType("i64")))
		}
	case "i64.andNot":
		if val, err := andNotI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i64"), MakeType("i64")))
		}
		// make functions
	case "[]bool.make":
		if val, err := makeArray("[]bool", (*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]bool"), MakeType("[]bool")))
		}
	case "[]byte.make":
		if val, err := makeArray("[]byte", (*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]byte"), MakeType("[]byte")))
		}
	case "[]i32.make":
		if val, err := makeArray("[]i32", (*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]i32"), MakeType("[]i32")))
		}
	case "[]i64.make":
		if val, err := makeArray("[]i64", (*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]i64"), MakeType("[]i64")))
		}
	case "[]f32.make":
		if val, err := makeArray("[]f32", (*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]f32"), MakeType("[]f32")))
		}
	case "[]f64.make":
		if val, err := makeArray("[]f64", (*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]f64"), MakeType("[]f64")))
		}
		// array functions
	case "[]bool.read":
		if val, err := readBoolA((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "[]bool.write":
		if val, err := writeBoolA((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]bool"), MakeType("[]bool")))
		}
	case "[]byte.read":
		if val, err := readByteA((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("byte"), MakeType("byte")))
		}
	case "[]byte.write":
		if val, err := writeByteA((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]byte"), MakeType("[]byte")))
		}
	case "[]i32.read":
		if val, err := readI32A((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "[]i32.write":
		if val, err := writeI32A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]i32"), MakeType("[]i32")))
		}
	case "[]i64.read":
		if val, err := readI64A((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i64"), MakeType("i64")))
		}
	case "[]i64.write":
		if val, err := writeI64A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]i64"), MakeType("[]i64")))
		}
	case "[]f32.read":
		if val, err := readF32A((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("f32"), MakeType("f32")))
		}
	case "[]f32.write":
		if val, err := writeF32A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]f32"), MakeType("[]f32")))
		}
	case "[]f64.read":
		if val, err := readF64A((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("f64"), MakeType("f64")))
		}
	case "[]f64.write":
		if val, err := writeF64A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("[]f64"), MakeType("[]f64")))
		}
	case "[]bool.len":
		if val, err := lenBoolA((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "[]byte.len":
		if val, err := lenByteA((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "[]i32.len":
		if val, err := lenI32A((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "[]i64.len":
		if val, err := lenI64A((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "[]f32.len":
		if val, err := lenF32A((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "[]f64.len":
		if val, err := lenF64A((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
		// time functions
	case "sleep":
		if val, err := sleep((*argsCopy)[0]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
		// utilitiy functions
	case "i32.rand":
		if val, err := randI32((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i32"), MakeType("i32")))
		}
	case "i64.rand":
		if val, err := randI64((*argsCopy)[0], (*argsCopy)[1]); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("i64"), MakeType("i64")))
		}
		// meta functions
	case "setClauses":
		if val, err := setClauses((*argsCopy)[0], call.Operator.Module); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("str"), MakeType("str")))
		}
	case "addObject":
		if val, err := addObject((*argsCopy)[0], call.Operator.Module); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("str"), MakeType("str")))
		}
	case "setQuery":
		if val, err := setQuery((*argsCopy)[0], call.Operator.Module); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("str"), MakeType("str")))
		}
	case "remObject":
		if val, err := remObject((*argsCopy)[0], call.Operator.Module); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("str"), MakeType("str")))
		}
	case "remObjects":
		if val, err := remObjects(call.Operator.Module); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "remExpr":
		if val, err := remExpr((*argsCopy)[0], call.Operator); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "remArg":
		if val, err := remArg((*argsCopy)[0], call.Operator); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
		// case "addArg":
		// 	if val, err := addArg((*argsCopy)[0], (*argsCopy)[1], call.Operator); err == nil {
		// 		*values = append(*values, val)
		// 	} else {
		// 		*exc = true
		// 		*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		// 		*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		// 	}
	case "addExpr":
		if val, err := addExpr((*argsCopy)[0], (*argsCopy)[1], call.Operator, expr.Line); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
	case "affExpr":
		if val, err := affExpr((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], call.Operator); err == nil {
			*values = append(*values, val)
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
			*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
		}
		// debugging functions
	case "halt":
		fmt.Println(string(*(*argsCopy)[0].Value))
		*exc = true
		*excError = errors.New(fmt.Sprintf("%d: call to halt", expr.FileLine))
		*values = append(*values, MakeArgument(MakeDefaultValue("bool"), MakeType("bool")))
	case "":
	}
}

func (call *CXCall) call (withDebug bool, nCalls, callCounter int) error {
	//  add a counter here to pause
	if nCalls > 0 && callCounter >= nCalls {
		return nil
	}
	callCounter++

	//saveStep(call)
	if withDebug {
		PrintCallStack(call.Context.CallStack.Calls)
	}
	
	if call.Line >= len(call.Operator.Expressions) || call.Line < 0 {
		// popping the stack
		call.Context.CallStack.Calls = call.Context.CallStack.Calls[:len(call.Context.CallStack.Calls) - 1]
		//outName := call.Operator.Output.Name


		outNames := make([]string, len(call.Operator.Outputs))
		//outNames := make([]string, biggestOutputNumber)
		for i, out := range call.Operator.Outputs {
			outNames[i] = out.Name
		}

		// this one is for returning result
		//output := call.State[outName]

		outputs := make([]*CXDefinition, len(call.Operator.Outputs))

		for i, outName := range outNames {
			//fmt.Printf("%s : %v\n", outName, call.State[outName])
			//fmt.Println(call.State[outName])

			//outputs[i] = call.State[outName]
			for _, stateDef := range call.State {
				if stateDef.Name == outName {
					outputs[i] = stateDef
				}
			}

			// if call.State[outName] == nil {
			// 	call.Context.PrintProgram(false)
			// }
		}

		allOutsNil := true
		for _, out := range outputs {
			if out != nil {
				allOutsNil = false
				break
			}
		}

		if allOutsNil && len(outputs) > 0 {
			outNames := call.Operator.Expressions[len(call.Operator.Expressions) - 1].OutputNames
			for i, outName := range outNames {
				for _, stateDef := range call.State {
					if stateDef.Name == outName.Name {
						outputs[i] = stateDef
					}
				}
			}
		}

		for i, out := range outputs {
			if out.Typ.Name != call.Operator.Outputs[i].Typ.Name {
				panic(fmt.Sprintf("output var '%s' is of type '%s'; function '%s' requires output of type '%s'",
					out.Name, out.Typ.Name, call.Operator.Name, call.Operator.Outputs[i].Typ.Name))
			}
		}
		
		//checking if output var has the same type as the required output
		// if output.Typ.Name != call.Operator.Output.Typ.Name {
		// 	panic(fmt.Sprintf("output var '%s' is of type '%s'; function '%s' requires output of type '%s'",
		// 		output.Name, output.Typ.Name, call.Operator.Name, call.Operator.Output.Typ.Name))
		// }

		if call.ReturnAddress != nil {
			//returnName := call.ReturnAddress.Operator.Expressions[call.ReturnAddress.Line - 1].OutputName

			returnNames := make([]string, len(call.ReturnAddress.Operator.Expressions[call.ReturnAddress.Line - 1].OutputNames))
			for i, out := range call.ReturnAddress.Operator.Expressions[call.ReturnAddress.Line - 1].OutputNames {
				returnNames[i] = out.Name
			}

			if len(outputs) > 0 {
				for i, returnName := range returnNames {
					def := MakeDefinition(returnName, outputs[i].Value, outputs[i].Typ)
					def.Module = call.Module
					def.Context = call.Context

					idx := -1
					for i, stateDef := range call.ReturnAddress.State {
						if stateDef.Name == returnName {
							idx = i
							break
						}
					}
					if idx < 0 {
						call.ReturnAddress.State = append(call.ReturnAddress.State, def)
					} else {
						call.ReturnAddress.State[idx] = def
					}
				}
			}

			return call.ReturnAddress.call(withDebug, nCalls, callCounter)
		} else {
			// no return address. should only be for main
			call.Context.Terminated = true
			call.Context.Outputs = outputs
		}
	} else {
		fn := call.Operator
		
		if expr, err := fn.GetExpression(call.Line); err == nil {
			
			// getting arguments
			//outName := expr.OutputName
			argsRefs, _ := expr.GetArguments()

			argsCopy := make([]*CXArgument, len(argsRefs))
			argNames := make([]string, len(argsRefs))

			// exceptions
			var exc bool
			var excError error
			
			if len(argNames) != len(expr.Operator.Inputs) {
				//fmt.Println("delete-me")
				// fmt.Println(argsRefs)
				// for _, ref := range argsRefs {
				// 	fmt.Println(string(*ref.Value))
				// }
				// fmt.Println(len(argNames))
				return errors.New(fmt.Sprintf("%d: %s: expected %d arguments; %d were provided",
					expr.FileLine, expr.Operator.Name, len(expr.Operator.Inputs), len(argNames)))
			}
			
			for i, inp := range expr.Operator.Inputs {
				argNames[i] = inp.Name
			}
			
			// we don't want to modify by reference, we need to make copies
			for i := 0; i < len(argsRefs); i++ {
				if argsRefs[i].Offset > -1 {
					offset := argsRefs[i].Offset
					size := argsRefs[i].Size
					val := (*call.Context.Heap)[offset:offset+size]
					argsRefs[i].Value = &val
					continue
				}
				if argsRefs[i].Typ.Name == "ident" {
					lookingFor := string(*argsRefs[i].Value)

					var resolvedIdent *CXDefinition

					identParts := GetIdentParts(lookingFor)
					if len(identParts) > 1 {
						if mod, err := call.Context.GetModule(identParts[0]); err == nil {
							// then it's an external definition or struct
							if def, err := mod.GetDefinition(concat(identParts[1:]...)); err == nil {
								resolvedIdent = def
							}
						} else {
							// then it's a global struct
							mod := call.Operator.Module
							if def, err := mod.GetDefinition(concat(identParts[:]...)); err == nil {
								resolvedIdent = def
							} else {
								// then it's a local struct
								for _, stateDef := range call.State {
									if stateDef.Name == lookingFor {
										resolvedIdent = stateDef
										break
									}
								}
							}
						}
					} else {
						// then it's a local or global definition
						local := false
						for _, stateDef := range call.State {
							if stateDef.Name == lookingFor {
								local = true
								resolvedIdent = stateDef
							}
						}

						if !local {
							mod := call.Operator.Module
							if def, err := mod.GetDefinition(lookingFor); err == nil {
								resolvedIdent = def
							}
						}
					}
					
					//if (local == nil && global == nil) {
					if resolvedIdent == nil {
						panic(fmt.Sprintf("'%s' is undefined", lookingFor))
					}
					argsCopy[i] = MakeArgument(resolvedIdent.Value, resolvedIdent.Typ)
				} else {
					argsCopy[i] = argsRefs[i]
				}

				// checking if arguments types match with expressions required types
				if len(expr.Operator.Inputs) > 0 && expr.Operator.Inputs[i].Typ.Name != argsCopy[i].Typ.Name {
					// panic(fmt.Sprintf("%s, line #%d: %s argument #%d is type '%s'; expected type '%s'",
					fmt.Printf("%d: %s: argument %d is type '%s'; expected type '%s'\n",
						expr.FileLine, expr.Operator.Name, i+1, argsCopy[i].Typ.Name, expr.Operator.Inputs[i].Typ.Name)
				}
			}

			// checking if native or not
			var values []*CXArgument
			var opName string
			if expr.Operator != nil {
				opName = expr.Operator.Name
			} else {
				opName = "id" // return the same
			}
			
			//fmt.Println(opName)
			
			checkNative(opName, expr, call, &values, &argsCopy, &exc, &excError)
			
			if len(values) > 0 {

				if exc {
					fmt.Println()
					fmt.Println("Call's State:")
					for _, def := range call.State {
						fmt.Printf("%s:\t\t%s\n", def.Name, PrintValue(def.Value, def.Typ.Name))
					}
					fmt.Println()
					fmt.Printf("%s() Arguments:\n", expr.Operator.Name)
					for i, arg := range argsCopy {
						fmt.Printf("%d: %s\n", i, PrintValue(arg.Value, arg.Typ.Name))
					}
					fmt.Println()
					return excError
				}
				
				// operator was a native function
				for i, outName := range expr.OutputNames {
					if withDebug {
						replPrintEvaluation(values[i])
					}
					def := MakeDefinition(outName.Name, values[i].Value, values[i].Typ)
					def.Module = call.Module
					def.Context = call.Context

					// if outName.Offset > -1 {
					// 	call.State[outName.Offset] = def
					// } else {
					
					// }
					
					found := false
					for i, stateDef := range call.State {
						if stateDef.Name == outName.Name {
							found = true
							call.State[i] = def
						}
					}
					if !found {
						call.State = append(call.State, def)
					}
				}
				call.Line++
				return call.call(withDebug, nCalls, callCounter)
			} else {
				// operator was not a native function
				call.Line++ // once the subcall finishes, call next line
				if argDefs, err := argsToDefs(argsCopy, argNames, call.Module, call.Context); err == nil {
					subcall := MakeCall(expr.Operator, argDefs, call, call.Module, call.Context)
					call.Context.CallStack.Calls = append(call.Context.CallStack.Calls, subcall)
					return subcall.call(withDebug, nCalls, callCounter)
				} else {
					fmt.Println(err)
				}
			}
		} else {
			fmt.Println(err)
		}
	}
	return nil
}
