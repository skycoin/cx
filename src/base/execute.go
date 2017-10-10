package base

import (
	"bytes"
	"io/ioutil"
	"runtime"

	"fmt"
	"errors"
	"strings"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func argsToDefs (args []*CXArgument, inputs []*CXParameter, mod *CXModule, cxt *CXProgram) ([]*CXDefinition, error) {
	if len(inputs) == len(args) {
		defs := make([]*CXDefinition, len(args), len(args) + 10)
		for i, arg := range args {
			defs[i] = &CXDefinition{
				Name: inputs[i].Name,
				Typ: arg.Typ,
				Value: arg.Value,
				Module: mod,
				Context: cxt,
			}
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
			switch def.Typ {
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
	cxt.Outputs = make([]*CXDefinition, 0)
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
	switch arg.Typ {
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
        "runtime"
	"runtime/pprof"
);

var cxt = MakeContext();var mod *CXModule;var imp *CXModule;var fn *CXFunction;var op *CXFunction;var expr *CXExpression;var strct *CXStruct;var arg *CXArgument;var tag string = "";

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile 'file'")
var memprofile = flag.String("memprofile", "", "write memory profile to 'file'")

func main () {
	runtime.LockOSThread()

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
		program.WriteString(`package main;import (. github.com/skycoin/cx/src/base"; "runtime";);var cxt = MakeContext();var mod *CXModule;var imp *CXModule;var fn *CXFunction;var op *CXFunction;var expr *CXExpression;var strct *CXStruct;var arg *CXArgument;var tag string = "";func main () {runtime.LockOSThread();`)
	}
	
	

	for _, mod := range cxt.Modules {
		program.WriteString(fmt.Sprintf(`mod = MakeModule("%s");cxt.AddModule(mod);%s`, mod.Name, asmNL))
		for _, imp := range mod.Imports {
			program.WriteString(fmt.Sprintf(`imp, _ = cxt.GetModule("%s");mod.AddImport(imp);%s`, imp.Name, asmNL))
		}

		for _, fn := range mod.Functions {
			program.WriteString(fmt.Sprintf(`fn = MakeFunction(%#v);mod.AddFunction(fn);%s`, fn.Name, asmNL))

			for _, inp := range fn.Inputs {
				program.WriteString(fmt.Sprintf(`fn.AddInput(MakeParameter("%s", "%s"));%s`, inp.Name, inp.Typ, asmNL))
			}

			for _, out := range fn.Outputs {
				program.WriteString(fmt.Sprintf(`fn.AddOutput(MakeParameter("%s", "%s"));%s`, out.Name, out.Typ, asmNL))
			}

			for _, expr := range fn.Expressions {
				var tagStr string
				if expr.Tag != "" {
					tagStr = fmt.Sprintf(`expr.Tag = tag;tag = "";`)
				}
				program.WriteString(fmt.Sprintf(`op, _ = cxt.GetFunction("%s", "%s");expr = MakeExpression(op);expr.FileLine = %d;fn.AddExpression(expr);%s%s`,
					expr.Operator.Name, expr.Operator.Module.Name, expr.FileLine, tagStr, asmNL))

				for _, arg := range expr.Arguments {
					program.WriteString(fmt.Sprintf(`expr.AddArgument(MakeArgument(&%#v, "%s"));%s`, *arg.Value, arg.Typ, asmNL))
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
			program.WriteString(fmt.Sprintf(`mod.AddDefinition(MakeDefinition("%s", &%#v, "%s"));%s`, def.Name, *def.Value, def.Typ, asmNL))
		}
	}

	program.WriteString(`
if *memprofile != "" {
        f, err := os.Create(*memprofile)
        if err != nil {
            log.Fatal("could not create memory profile: ", err)
        }
        runtime.GC() // get up-to-date statistics
        if err := pprof.WriteHeapProfile(f); err != nil {
            log.Fatal("could not write memory profile: ", err)
        }
        f.Close()
    }
`)

	program.WriteString(`cxt.Run(false, -1);}`)
	ioutil.WriteFile(fmt.Sprintf("o.go"), []byte(program.String()), 0644)
}


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
				state := make([]*CXDefinition, 0, 20)
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

func checkNative (opName string, expr *CXExpression, call *CXCall, argsCopy *[]*CXArgument, exc *bool, excError *error) {
	switch opName {
	case "serialize":
		//Serialize(call.Context, expr, call)
		serialize_program(expr, call)
	case "deserialize":
		// it only prints the deserialized program for now
		Deserialize((*argsCopy)[0].Value).PrintProgram(false)
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

		if err := call.Context.Evolve(fnName, fnBag, inps, outs, int(numberExprs), int(iterations), epsilon, expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// flow control
	case "baseGoTo":
		if err := baseGoTo(call, (*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2]); err == nil {

		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "goTo":
		if err := goTo(call, (*argsCopy)[0]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
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
	case "str.print":
		fmt.Println(string(*(*argsCopy)[0].Value))
	case "byte.print":
		fmt.Println((*(*argsCopy)[0].Value)[0])
	case "i32.print":
		var val int32
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "i64.print":
		var val int64
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "f32.print":
		var val float32
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "f64.print":
		var val float64
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
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
	case "[]byte.print":
		fmt.Println(*(*argsCopy)[0].Value)
	case "[]i32.print":
		var val []int32
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "[]i64.print":
		var val []int64
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "[]f32.print":
		var val []float32
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "[]f64.print":
		var val []float64
		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
		fmt.Println(val)
		// identity functions
	case "str.id", "bool.id", "byte.id", "i32.id", "i64.id", "f32.id", "f64.id", "[]bool.id", "[]byte.id", "[]i32.id", "[]i64.id", "[]f32.id", "[]f64.id":
		found := false
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = (*argsCopy)[0].Value
				found = true
				break
			}
		}
		if !found {
			call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, (*argsCopy)[0].Value, (*argsCopy)[0].Typ))
		}
		// cast functions
	case "[]byte.str":
		if err := castToStr((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "str.[]byte":
		if err := castToByteA((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i32.byte", "i64.byte", "f32.byte", "f64.byte":
		if err := castToByte((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "byte.i32", "i64.i32", "f32.i32", "f64.i32":
		if err := castToI32((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "byte.i64", "i32.i64", "f32.i64", "f64.i64":
		if err := castToI64((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "byte.f32", "i32.f32", "i64.f32", "f64.f32":
		if err := castToF32((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "byte.f64", "i32.f64", "i64.f64", "f32.f64":
		if err := castToF64((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i32.[]byte", "[]i64.[]byte", "[]f32.[]byte", "[]f64.[]byte":
		if err := castToByteA((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]byte.[]i32", "[]i64.[]i32", "[]f32.[]i32", "[]f64.[]i32":
		if err := castToI32A((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]byte.[]i64", "[]i32.[]i64", "[]f32.[]i64", "[]f64.[]i64":
		if err := castToI64A((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]byte.[]f32", "[]i32.[]f32", "[]i64.[]f32", "[]f64.[]f32":
		if err := castToF32A((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]byte.[]f64", "[]i32.[]f64", "[]i64.[]f64", "[]f32.[]f64":
		if err := castToF64A((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// logical operators
	case "and":
		if err := and((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "or":
		if err := or((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "not":
		if err := not((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// relational operators
	case "i32.lt":
		if err := ltI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i32.gt":
		if err := gtI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i32.eq":
		if err := eqI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i32.lteq":
		if err := lteqI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i32.gteq":
		if err := gteqI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.lt":
		if err := ltI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.gt":
		if err := gtI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.eq":
		if err := eqI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.lteq":
		if err := lteqI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.gteq":
		if err := gteqI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f32.lt":
		if err := ltF32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f32.gt":
		if err := gtF32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f32.eq":
		if err := eqF32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f32.lteq":
		if err := lteqF32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f32.gteq":
		if err := gteqF32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f64.lt":
		if err := ltF64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f64.gt":
		if err := gtF64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f64.eq":
		if err := eqF64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f64.lteq":
		if err := lteqF64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f64.gteq":
		if err := gteqF64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "str.lt":
		if err := ltStr((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "str.eq":
		if err := eqStr((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "str.lteq":
		if err := lteqStr((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "str.gteq":
		if err := gteqStr((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "byte.lt":
		if err := ltByte((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "byte.gt":
		if err := gtByte((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "byte.eq":
		if err := eqByte((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "byte.lteq":
		if err := lteqByte((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "byte.gteq":
		if err := gteqByte((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// struct operations
	case "initDef":
		if err := initDef((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// arithmetic functions
	case "i32.add":
		if err := addI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i32.mul":
		if err := mulI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i32.sub":
		if err := subI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i32.div":
		if err := divI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.add":
		if err := addI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.mul":
		if err := mulI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.sub":
		if err := subI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.div":
		if err := divI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f32.add":
		if err := addF32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f32.mul":
		if err := mulF32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f32.sub":
		if err := subF32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f32.div":
		if err := divF32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f64.add":
		if err := addF64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f64.mul":
		if err := mulF64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f64.sub":
		if err := subF64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "f64.div":
		if err := divF64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i32.mod":
		if err := modI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.mod":
		if err := modI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// bitwise operators
	case "i32.bitand":
		if err := andI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i32.bitor":
		if err := orI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i32.bitxor":
		if err := xorI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i32.bitclear":
		if err := andNotI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.bitand":
		if err := andI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.bitor":
		if err := orI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.bitxor":
		if err := xorI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.bitclear":
		if err := andNotI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// make functions
	case "[]bool.make":
		if err := makeArray("[]bool", (*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]byte.make":
		if err := makeArray("[]byte", (*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i32.make":
		if err := makeArray("[]i32", (*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i64.make":
		if err := makeArray("[]i64", (*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]f32.make":
		if err := makeArray("[]f32", (*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]f64.make":
		if err := makeArray("[]f64", (*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// array functions
	case "[]bool.read":
		if err := readBoolA((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]bool.write":
		if err := writeBoolA((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]byte.read":
		if err := readByteA((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]byte.write":
		if err := writeByteA((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i32.read":
		if err := readI32A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i32.write":
		if err := writeI32A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i64.read":
		if err := readI64A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i64.write":
		if err := writeI64A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]f32.read":
		if err := readF32A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]f32.write":
		if err := writeF32A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]f64.read":
		if err := readF64A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]f64.write":
		if err := writeF64A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]bool.len":
		if err := lenBoolA((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]byte.len":
		if err := lenByteA((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i32.len":
		if err := lenI32A((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i64.len":
		if err := lenI64A((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]f32.len":
		if err := lenF32A((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]f64.len":
		if err := lenF64A((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// concatenation functions
	case "str.concat":
		if err := concatStr((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]byte.append":
		if err := appendByteA((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]bool.append":
		if err := appendBoolA((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i32.append":
		if err := appendI32A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i64.append":
		if err := appendI64A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]f32.append":
		if err := appendF32A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]f64.append":
		if err := appendF64A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]byte.concat":
		if err := concatByteA((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]bool.concat":
		if err := concatBoolA((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i32.concat":
		if err := concatI32A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i64.concat":
		if err := concatI64A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]f32.concat":
		if err := concatF32A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]f64.concat":
		if err := concatF64A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// copy functions
	case "[]byte.copy":
		if err := copyByteA((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]bool.copy":
		if err := copyBoolA((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i32.copy":
		if err := copyI32A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]i64.copy":
		if err := copyI64A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]f32.copy":
		if err := copyF32A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "[]f64.copy":
		if err := copyF64A((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// time functions
	case "sleep":
		if err := sleep((*argsCopy)[0]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// utilitiy functions
	case "i32.rand":
		if err := randI32((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "i64.rand":
		if err := randI64((*argsCopy)[0], (*argsCopy)[1], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// meta functions
	case "setClauses":
		if err := setClauses((*argsCopy)[0], call.Operator.Module); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "addObject":
		if err := addObject((*argsCopy)[0], call.Operator.Module); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "setQuery":
		if err := setQuery((*argsCopy)[0], call.Operator.Module); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "remObject":
		if err := remObject((*argsCopy)[0], call.Operator.Module); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "remObjects":
		if err := remObjects(call.Operator.Module); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "remExpr":
		if err := remExpr((*argsCopy)[0], call.Operator); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "remArg":
		if err := remArg((*argsCopy)[0], call.Operator); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// case "addArg":
		// 	if err := addArg((*argsCopy)[0], (*argsCopy)[1], expr, call, call.Operator); err == nil {
		// 	} else {
		// 		*exc = true
		// 		*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		// 	}
	case "addExpr":
		if err := addExpr((*argsCopy)[0], (*argsCopy)[1], call.Operator, expr.Line); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "affExpr":
		if err := affExpr((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], call.Operator, expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// debugging functions
	case "halt":
		fmt.Println(string(*(*argsCopy)[0].Value))
		*exc = true
		*excError = errors.New(fmt.Sprintf("%d: call to halt", expr.FileLine))
		// Runtime
	case "runtime.LockOSThread":
		runtime.LockOSThread()
		// if err := runtime_LockOSThread(expr, call); err == nil {
		// } else {
		// 	*exc = true
		// 	*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		// }
		// OpenGL
	case "gl.Init":
		if err := gl_Init(expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.CreateProgram":
		if err := gl_CreateProgram(expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.LinkProgram":
		if err := gl_LinkProgram((*argsCopy)[0]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.Clear":
		if err := gl_Clear((*argsCopy)[0]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.UseProgram":
		if err := gl_UseProgram((*argsCopy)[0]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.BindBuffer":
		if err := gl_BindBuffer((*argsCopy)[0], (*argsCopy)[1]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.BindVertexArray":
		if err := gl_BindVertexArray((*argsCopy)[0]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.EnableVertexAttribArray":
		if err := gl_EnableVertexAttribArray((*argsCopy)[0]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.VertexAttribPointer":
		if err := gl_VertexAttribPointer((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3], (*argsCopy)[4]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.DrawArrays":
		if err := gl_DrawArrays((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.GenBuffers":
		if err := gl_GenBuffers((*argsCopy)[0], (*argsCopy)[1]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.BufferData":
		if err := gl_BufferData((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.GenVertexArrays":
		if err := gl_GenVertexArrays((*argsCopy)[0], (*argsCopy)[1]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.CreateShader":
		if err := gl_CreateShader((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.Strs":
		if err := gl_Strs((*argsCopy)[0], (*argsCopy)[1]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.Free":
		if err := gl_Free((*argsCopy)[0]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.ShaderSource":
		if err := gl_ShaderSource((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.CompileShader":
		if err := gl_CompileShader((*argsCopy)[0]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.GetShaderiv":
		if err := gl_GetShaderiv((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "gl.AttachShader":
		if err := gl_AttachShader((*argsCopy)[0], (*argsCopy)[1]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
		// GLFW
	case "glfw.Init":
		if err := glfw_Init(); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "glfw.WindowHint":
		if err := glfw_WindowHint((*argsCopy)[0], (*argsCopy)[1]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "glfw.CreateWindow":
		if err := glfw_CreateWindow((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "glfw.MakeContextCurrent":
		if err := glfw_MakeContextCurrent((*argsCopy)[0]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "glfw.ShouldClose":
		if err := glfw_ShouldClose((*argsCopy)[0], expr, call); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "glfw.PollEvents":
		if err := glfw_PollEvents(); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
	case "glfw.SwapBuffers":
		if err := glfw_SwapBuffers((*argsCopy)[0]); err == nil {
		} else {
			*exc = true
			*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
		}
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

		// we're looking for the output name in the state, and then we look in return's state for that
		for i, out := range call.Operator.Outputs {
			found := true
			for _, def := range call.State {
				/////////// throw error if output was not defined, or handle outputs from last expression
				if out.Name == def.Name {
					if call.ReturnAddress != nil {
						retName := call.ReturnAddress.Operator.Expressions[call.ReturnAddress.Line - 1].OutputNames[i].Name
						found := false
						for _, retDef := range call.ReturnAddress.State {
							if retDef.Name == retName {
								retDef.Value = def.Value
								found = true
								break
							}
						}
						if !found {
							def.Name = retName
							call.ReturnAddress.State = append(call.ReturnAddress.State, def)
						}

						found = true
						// break
						return call.ReturnAddress.call(withDebug, nCalls, callCounter)
					} else {
						// no return address. should only be for main
						call.Context.Terminated = true
						call.Context.Outputs = append(call.Context.Outputs, def)
					}
				}
			}

			// this isn't complete yet
			if !found {
				return errors.New(fmt.Sprintf("'%s' output(s) not specified", call.Operator.Name))
			}
		}

		if call.ReturnAddress != nil {
			return call.ReturnAddress.call(withDebug, nCalls, callCounter)
		} else {
			// no return address. should only be for main
			call.Context.Terminated = true
			//call.Context.Outputs = append(call.Context.Outputs, def)
		}
	} else {
		fn := call.Operator
		
		if expr, err := fn.GetExpression(call.Line); err == nil {
			
			// getting arguments
			argsRefs, _ := expr.GetArguments()

			argsCopy := make([]*CXArgument, len(argsRefs))
			//argNames := make([]string, len(argsRefs))

			// exceptions
			var exc bool
			var excError error
			
			if len(argsRefs) != len(expr.Operator.Inputs) {
				return errors.New(fmt.Sprintf("%d: %s: expected %d arguments; %d were provided",
					expr.FileLine, expr.Operator.Name, len(expr.Operator.Inputs), len(argsRefs)))
			}
			
			// we don't want to modify by reference, we need to make copies
			for i := 0; i < len(argsRefs); i++ {
				if argsRefs[i].Typ == "ident" {
					lookingFor := string(*argsRefs[i].Value)

					var resolvedIdent *CXDefinition

					identParts := strings.Split(lookingFor, ".")

					if len(identParts) > 1 {
						if mod, err := call.Context.GetModule(identParts[0]); err == nil {
							// then it's an external definition or struct
							isImported := false
							for _, imp := range call.Operator.Module.Imports {
								if imp.Name == identParts[0] {
									isImported = true
									break
								}
							}
							if isImported {
								if def, err := mod.GetDefinition(concat(identParts[1:]...)); err == nil {
									resolvedIdent = def
								}
							} else {
								return errors.New(fmt.Sprintf("Module '%s' not imported", mod.Name))
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
								break
							}
						}

						if !local {
							mod := call.Operator.Module
							if def, err := mod.GetDefinition(lookingFor); err == nil {
								resolvedIdent = def
							}
						}
					}

					if resolvedIdent == nil {
						return errors.New(fmt.Sprintf("%d: '%s' is undefined", expr.FileLine, lookingFor))
					}
					argsCopy[i] = MakeArgument(resolvedIdent.Value, resolvedIdent.Typ)
				} else {
					argsCopy[i] = argsRefs[i]
				}

				// checking if arguments types match with expressions required types
				if len(expr.Operator.Inputs) > 0 && expr.Operator.Inputs[i].Typ != argsCopy[i].Typ {
					fmt.Println()
					// panic(fmt.Sprintf("%s, line #%d: %s argument #%d is type '%s'; expected type '%s'",
					return errors.New(fmt.Sprintf("%d: %s: argument %d is type '%s'; expected type '%s'\n",
						expr.FileLine, expr.Operator.Name, i+1, argsCopy[i].Typ, expr.Operator.Inputs[i].Typ))
				}
			}

			var opName string
			if expr.Operator != nil {
				opName = expr.Operator.Name
			} else {
				opName = "id" // return the same
			}

			isNative := false
			if _, ok := NATIVE_FUNCTIONS[opName]; ok {
				isNative = true
			}

			if isNative {
				checkNative(opName, expr, call, &argsCopy, &exc, &excError)

				if exc {
					fmt.Println()
					fmt.Println("Call's State:")
					for _, def := range call.State {
						fmt.Printf("%s:\t\t%s\n", def.Name, PrintValue(def.Value, def.Typ))
					}
					fmt.Println()
					fmt.Printf("%s() Arguments:\n", expr.Operator.Name)
					for i, arg := range argsCopy {
						fmt.Printf("%d: %s\n", i, PrintValue(arg.Value, arg.Typ))
					}
					fmt.Println()
					return excError
				}
				
				call.Line++
				return call.call(withDebug, nCalls, callCounter)
			} else {
				// operator was not a native function
				call.Line++ // once the subcall finishes, call next line
				if argDefs, err := argsToDefs(argsCopy, expr.Operator.Inputs, call.Module, call.Context); err == nil {
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
