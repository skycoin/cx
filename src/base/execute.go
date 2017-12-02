package base

import (
	"bytes"
	"io/ioutil"
	"runtime"

	"strconv"
	"fmt"
	"errors"
	"strings"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func argsToDefs (args []*CXArgument, inputs []*CXParameter, outputs []*CXParameter, mod *CXModule, cxt *CXProgram) ([]*CXDefinition, error) {
	if len(inputs) == len(args) {
		defs := make([]*CXDefinition, len(args) + len(outputs), len(args) + len(outputs) + 10)
		for i, arg := range args {
			defs[i] = &CXDefinition{
				Name: inputs[i].Name,
				Typ: arg.Typ,
				Value: arg.Value,
				Module: mod,
				Context: cxt,
			}
		}
		for i, out := range outputs {
			var zeroValue []byte
			isBasic := false
			for _, basic := range BASIC_TYPES {
				if basic == out.Typ {
					zeroValue = *MakeDefaultValue(basic)
					isBasic = true
					break
				}
			}
			if !isBasic {
				var err error
				if zeroValue, err = ResolveStruct(out.Typ, cxt); err != nil {
					return nil, err
				}
			}
			defs[i+len(args)] = &CXDefinition{
				Name: out.Name,
				Typ: out.Typ,
				Value: zeroValue,
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
				encoder.DeserializeRaw(def.Value, &valI32)
				if idx == lenState - 1 {
					fmt.Printf("%s: %d", def.Name, valI32)
				} else {
					fmt.Printf("%s: %d, ", def.Name, valI32)
				}
			case "i64":
				encoder.DeserializeRaw(def.Value, &valI64)
				if idx == lenState - 1 {
					fmt.Printf("%s: %d", def.Name, valI64)
				} else {
					fmt.Printf("%s: %d, ", def.Name, valI64)
				}
			case "f32":
				encoder.DeserializeRaw(def.Value, &valF32)
				if idx == lenState - 1 {
					fmt.Printf("%s: %f", def.Name, valF32)
				} else {
					fmt.Printf("%s: %f, ", def.Name, valF32)
				}
			case "f64":
				encoder.DeserializeRaw(def.Value, &valF64)
				if idx == lenState - 1 {
					fmt.Printf("%s: %f", def.Name, valF64)
				} else {
					fmt.Printf("%s: %f, ", def.Name, valF64)
				}
			case "byte":
				if idx == lenState - 1 {
					fmt.Printf("%s: %d", def.Name, def.Value[0])
				} else {
					fmt.Printf("%s: %d, ", def.Name, def.Value[0])
				}
			case "[]byte":
				var val []byte
				encoder.DeserializeRaw(def.Value, &val)
				if idx == lenState - 1 {
					fmt.Printf("%s: %v", def.Name, val)
				} else {
					fmt.Printf("%s: %v, ", def.Name, val)
				}
			case "[]i32":
				var val []int32
				encoder.DeserializeRaw(def.Value, &val)
				if idx == lenState - 1 {
					fmt.Printf("%s: %v", def.Name, val)
				} else {
					fmt.Printf("%s: %v, ", def.Name, val)
				}
			case "[]i64":
				var val []int64
				encoder.DeserializeRaw(def.Value, &val)
				if idx == lenState - 1 {
					fmt.Printf("%s: %v", def.Name, val)
				} else {
					fmt.Printf("%s: %v, ", def.Name, val)
				}
			case "[]f32":
				var val []float32
				encoder.DeserializeRaw(def.Value, &val)
				if idx == lenState - 1 {
					fmt.Printf("%s: %v", def.Name, val)
				} else {
					fmt.Printf("%s: %v, ", def.Name, val)
				}
			case "[]f64":
				var val []float64
				encoder.DeserializeRaw(def.Value, &val)
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
		var val string
		encoder.DeserializeRaw(arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "bool":
		var val int32
		encoder.DeserializeRaw(arg.Value, &val)
		if val == 0 {
			fmt.Printf("false\n")
		} else {
			fmt.Printf("true\n")
		}
	case "byte":
		fmt.Printf("%#v\n", arg.Value)
	case "i32":
		var val int32
		encoder.DeserializeRaw(arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "i64":
		var val int64
		encoder.DeserializeRaw(arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "f32":
		var val float32
		encoder.DeserializeRaw(arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "f64":
		var val float64
		encoder.DeserializeRaw(arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "[]byte":
		var val []byte
		encoder.DeserializeRaw(arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(arg.Value, &val)
		fmt.Printf("%#v\n", val)
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(arg.Value, &val)
		fmt.Printf("%#v\n", val)
	default:
		fmt.Printf("")
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
        "runtime"
	"runtime/pprof"
);

var cxt = MakeContext();var mod *CXModule;var imp *CXModule;var fn *CXFunction;var op *CXFunction;var expr *CXExpression;var strct *CXStruct;var fld *CXField;var arg *CXArgument;var tag string = "";

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
					tagStr = fmt.Sprintf(`expr.Tag = "%s";`, expr.Tag)
				}
				program.WriteString(fmt.Sprintf(`op, _ = cxt.GetFunction("%s", "%s");expr = MakeExpression(op);expr.FileLine = %d;fn.AddExpression(expr);%s%s`,
					expr.Operator.Name, expr.Operator.Module.Name, expr.FileLine, tagStr, asmNL))

				for _, arg := range expr.Arguments {
					program.WriteString(fmt.Sprintf(`expr.AddArgument(MakeArgument(&%#v, "%s"));%s`, arg.Value, arg.Typ, asmNL))
				}
				
				for _, outName := range expr.OutputNames {
					program.WriteString(fmt.Sprintf(`expr.AddOutputName("%s");%s`, outName.Name, asmNL))
				}
			}
		}

		for _, strct := range mod.Structs {
			program.WriteString(fmt.Sprintf(`strct = MakeStruct("%s");mod.AddStruct(strct);%s`, strct.Name, asmNL))
			for _, fld := range strct.Fields {
				// here here
				program.WriteString(fmt.Sprintf(`fld = MakeField("%s", "%s");strct.AddField(fld);%s`, fld.Name, fld.Typ, asmNL))
			}
		}

		for _, def := range mod.Definitions {
			program.WriteString(fmt.Sprintf(`mod.AddDefinition(MakeDefinition("%s", &%#v, "%s"));%s`, def.Name, def.Value, def.Typ, asmNL))
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

var isTesting bool
var isErrorPresent bool

func checkNative (opName string, expr *CXExpression, call *CXCall, argsCopy *[]*CXArgument, exc *bool, excError *error) {
	var err error
	switch opName {
	case "serialize": serialize_program(expr, call)
	case "deserialize":
		// it only prints the deserialized program for now
		Deserialize((*argsCopy)[0].Value).PrintProgram(false)
	case "evolve":
		var fnName string
		var fnBag string
		encoder.DeserializeRaw((*argsCopy)[0].Value, &fnName)
		encoder.DeserializeRaw((*argsCopy)[1].Value, &fnBag)
		
		var inps []float64
		encoder.DeserializeRaw((*argsCopy)[2].Value, &inps)
		
		var outs []float64
		encoder.DeserializeRaw((*argsCopy)[3].Value, &outs)

		var numberExprs int32
		encoder.DeserializeRaw((*argsCopy)[4].Value, &numberExprs)
		var iterations int32
		encoder.DeserializeRaw((*argsCopy)[5].Value, &iterations)
		var epsilon float64
		encoder.DeserializeRaw((*argsCopy)[6].Value, &epsilon)

		err = call.Context.Evolve(fnName, fnBag, inps, outs, int(numberExprs), int(iterations), epsilon, expr, call)
		// flow control
	case "baseGoTo": err = baseGoTo(call, (*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
	case "goTo": err = goTo(call, (*argsCopy)[0])
		// I/O functions
	case "bool.print":
		var val int32
		encoder.DeserializeRaw((*argsCopy)[0].Value, &val)
		if val == 0 {
			fmt.Println("false")
		} else {
			fmt.Println("true")
		}
	case "str.print":
		var val string
		encoder.DeserializeRaw((*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "byte.print":
		fmt.Println((*(*argsCopy)[0].Value)[0])
	case "i32.print":
		var val int32
		encoder.DeserializeRaw((*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "i64.print":
		var val int64
		encoder.DeserializeRaw((*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "f32.print":
		var val float32
		encoder.DeserializeRaw((*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "f64.print":
		var val float64
		encoder.DeserializeRaw((*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "[]bool.print":
		var val []int32
		encoder.DeserializeRaw((*argsCopy)[0].Value, &val)
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
		var val []byte
		encoder.DeserializeRaw((*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "[]str.print":
		var val []string
		encoder.DeserializeRaw((*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "[]i32.print":
		var val []int32
		encoder.DeserializeRaw((*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "[]i64.print":
		var val []int64
		encoder.DeserializeRaw((*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "[]f32.print":
		var val []float32
		encoder.DeserializeRaw((*argsCopy)[0].Value, &val)
		fmt.Println(val)
	case "[]f64.print":
		var val []float64
		encoder.DeserializeRaw((*argsCopy)[0].Value, &val)
		fmt.Println(val)
		// identity functions
	case "str.id", "bool.id", "byte.id", "i32.id", "i64.id", "f32.id", "f64.id", "[]bool.id", "[]byte.id", "[]str.id", "[]i32.id", "[]i64.id", "[]f32.id", "[]f64.id": assignOutput(0, *(*argsCopy)[0].Value, (*argsCopy)[0].Typ, expr, call)
	case "identity": identity((*argsCopy)[0], expr, call)
		// cast functions
	case "[]byte.str", "byte.str", "bool.str", "i32.str", "i64.str", "f32.str", "f64.str": err = castToStr((*argsCopy)[0], expr, call)
	case "str.[]byte": err = castToByteA((*argsCopy)[0], expr, call)
	case "i32.byte", "i64.byte", "f32.byte", "f64.byte": err = castToByte((*argsCopy)[0], expr, call)
	case "byte.i32", "i64.i32", "f32.i32", "f64.i32": err = castToI32((*argsCopy)[0], expr, call)
	case "byte.i64", "i32.i64", "f32.i64", "f64.i64": err = castToI64((*argsCopy)[0], expr, call)
	case "byte.f32", "i32.f32", "i64.f32", "f64.f32": err = castToF32((*argsCopy)[0], expr, call)
	case "byte.f64", "i32.f64", "i64.f64", "f32.f64": err = castToF64((*argsCopy)[0], expr, call)
	case "[]i32.[]byte", "[]i64.[]byte", "[]f32.[]byte", "[]f64.[]byte": err = castToByteA((*argsCopy)[0], expr, call)
	case "[]byte.[]i32", "[]i64.[]i32", "[]f32.[]i32", "[]f64.[]i32": err = castToI32A((*argsCopy)[0], expr, call)
	case "[]byte.[]i64", "[]i32.[]i64", "[]f32.[]i64", "[]f64.[]i64": err = castToI64A((*argsCopy)[0], expr, call)
	case "[]byte.[]f32", "[]i32.[]f32", "[]i64.[]f32", "[]f64.[]f32": err = castToF32A((*argsCopy)[0], expr, call)
	case "[]byte.[]f64", "[]i32.[]f64", "[]i64.[]f64", "[]f32.[]f64": err = castToF64A((*argsCopy)[0], expr, call)
		// logical operators
	case "and": err = and((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "or": err = or((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "not": err = not((*argsCopy)[0], expr, call)
		// relational operators
	case "i32.lt": err = ltI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i32.gt": err = gtI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i32.eq": err = eqI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i32.lteq": err = lteqI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i32.gteq": err = gteqI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.lt": err = ltI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.gt": err = gtI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.eq": err = eqI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.lteq": err = lteqI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.gteq": err = gteqI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f32.lt": err = ltF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f32.gt": err = gtF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f32.eq": err = eqF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f32.lteq": err = lteqF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f32.gteq": err = gteqF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f64.lt": err = ltF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f64.gt": err = gtF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f64.eq": err = eqF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f64.lteq": err = lteqF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f64.gteq": err = gteqF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "str.lt": err = ltStr((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "str.gt": err = gtStr((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "str.eq": err = eqStr((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "str.lteq": err = lteqStr((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "str.gteq": err = gteqStr((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "byte.lt": err = ltByte((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "byte.gt": err = gtByte((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "byte.eq": err = eqByte((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "byte.lteq": err = lteqByte((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "byte.gteq": err = gteqByte((*argsCopy)[0], (*argsCopy)[1], expr, call)
		// io functions
	case "str.read": err = readStr(expr, call)
	case "i32.read": err = readI32(expr, call)
		// struct operations
	case "initDef": err = initDef((*argsCopy)[0], expr, call)
		// arithmetic functions
	case "i32.add": err = addI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i32.mul": err = mulI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i32.sub": err = subI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i32.div": err = divI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.add": err = addI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.mul": err = mulI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.sub": err = subI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.div": err = divI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f32.add": err = addF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f32.mul": err = mulF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f32.sub": err = subF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f32.div": err = divF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f32.cos": err = cosF32((*argsCopy)[0], expr, call)
	case "f32.sin": err = sinF32((*argsCopy)[0], expr, call)
	case "f64.add": err = addF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f64.mul": err = mulF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f64.sub": err = subF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f64.div": err = divF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "f64.cos": err = cosF64((*argsCopy)[0], expr, call)
	case "f64.sin": err = sinF64((*argsCopy)[0], expr, call)
	case "i32.abs": err = absI32((*argsCopy)[0], expr, call)
	case "i64.abs": err = absI64((*argsCopy)[0], expr, call)
	case "f32.abs": err = absF32((*argsCopy)[0], expr, call)
	case "f64.abs": err = absF64((*argsCopy)[0], expr, call)
	case "i32.mod": err = modI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.mod": err = modI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
		// bitwise operators
	case "i32.bitand": err = andI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i32.bitor": err = orI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i32.bitxor": err = xorI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i32.bitclear": err = andNotI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i32.bitshl": err = shiftLeftI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i32.bitshr": err = shiftRightI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.bitand": err = andI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.bitor": err = orI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.bitxor": err = xorI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.bitclear": err = andNotI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.bitshl": err = shiftLeftI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.bitshr": err = shiftRightI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
		// make functions
	case "[]bool.make": err = makeArray("[]bool", (*argsCopy)[0], expr, call)
	case "[]byte.make": err = makeArray("[]byte", (*argsCopy)[0], expr, call)
	case "[]str.make": err = makeArray("[]str", (*argsCopy)[0], expr, call)
	case "[]i32.make": err = makeArray("[]i32", (*argsCopy)[0], expr, call)
	case "[]i64.make": err = makeArray("[]i64", (*argsCopy)[0], expr, call)
	case "[]f32.make": err = makeArray("[]f32", (*argsCopy)[0], expr, call)
	case "[]f64.make": err = makeArray("[]f64", (*argsCopy)[0], expr, call)
		// array functions
	case "[]bool.read": err = readBoolA((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]bool.write": err = writeBoolA((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
	case "[]byte.read": err = readByteA((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]byte.write": err = writeByteA((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
	case "[]str.read": err = readStrA((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]str.write": err = writeStrA((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
	case "[]i32.read": err = readI32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]i32.write": err = writeI32A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
	case "[]i64.read": err = readI64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]i64.write": err = writeI64A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
	case "[]f32.read": err = readF32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]f32.write": err = writeF32A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
	case "[]f64.read": err = readF64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]f64.write": err = writeF64A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
	case "[]bool.len": err = lenBoolA((*argsCopy)[0], expr, call)
	case "[]byte.len": err = lenByteA((*argsCopy)[0], expr, call)
	case "str.len": err = lenStr((*argsCopy)[0], expr, call)
	case "[]str.len": err = lenStrA((*argsCopy)[0], expr, call)
	case "[]i32.len": err = lenI32A((*argsCopy)[0], expr, call)
	case "[]i64.len": err = lenI64A((*argsCopy)[0], expr, call)
	case "[]f32.len": err = lenF32A((*argsCopy)[0], expr, call)
	case "[]f64.len": err = lenF64A((*argsCopy)[0], expr, call)
		// concatenation functions
	case "str.concat": err = concatStr((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]byte.append": err = appendByteA((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]bool.append": err = appendBoolA((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]str.append": err = appendStrA((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]i32.append": err = appendI32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]i64.append": err = appendI64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]f32.append": err = appendF32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]f64.append": err = appendF64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]byte.concat": err = concatByteA((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]bool.concat": err = concatBoolA((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]str.concat", "aff.concat": err = concatStrA((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]i32.concat": err = concatI32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]i64.concat": err = concatI64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]f32.concat": err = concatF32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]f64.concat": err = concatF64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
		// copy functions
	case "[]byte.copy": err = copyByteA((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]bool.copy": err = copyBoolA((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]str.copy": err = copyStrA((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]i32.copy": err = copyI32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]i64.copy": err = copyI64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]f32.copy": err = copyF32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "[]f64.copy": err = copyF64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
		// time functions
	case "sleep": err = sleep((*argsCopy)[0])
		// utilitiy functions
	case "i32.rand": err = randI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "i64.rand": err = randI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
		// meta functions
		
	case "aff.query": err = aff_query((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
	case "aff.execute": err = aff_execute((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
	case "aff.index": err = aff_index((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "aff.name": err = aff_name((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "aff.print": err = aff_print((*argsCopy)[0], call)
	case "aff.len": err = aff_len((*argsCopy)[0], expr, call)

	case "rem.expr": err = rem_expr((*argsCopy)[0], call.Operator)
	case "rem.arg": err = rem_arg((*argsCopy)[0], call.Operator)
	case "add.expr": err = add_expr((*argsCopy)[0], (*argsCopy)[1], call)
		// debugging functions
	case "halt":
		var msg string
		encoder.DeserializeRaw((*argsCopy)[0].Value, &msg)
		fmt.Println(msg)
		call.Line++
		*exc = true
		*excError = errors.New(fmt.Sprintf("%d: call to halt", expr.FileLine))
	case "test.start": isTesting = true
	case "test.stop": isTesting = false
	case "test.error":
		//fmt.Println(isErrorPresent)
		err = test_error((*argsCopy)[0], isErrorPresent, expr)
		isErrorPresent = false
		case "test.bool", "test.byte", "test.str", "test.i32", "test.i64", "test.f32", "test.f64", "test.[]bool", "test.[]byte", "test.[]str", "test.[]i32", "test.[]f32", "test.[]f64":
		err = test_value((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr)
		// custom types functions
	case "cstm.append": err = cstm_append((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "cstm.read": err = cstm_read((*argsCopy)[0], (*argsCopy)[1], expr, call)
	case "cstm.write": err = cstm_write((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
	case "cstm.len": err = cstm_len((*argsCopy)[0], expr, call)
	case "cstm.make": err = cstm_make((*argsCopy)[0], (*argsCopy)[1], expr, call)
		// Runtime
	case "runtime.LockOSThread": runtime.LockOSThread()
		// OpenGL
	case "gl.Init": err = gl_Init()
	case "gl.CreateProgram": err = gl_CreateProgram(expr, call)
	case "gl.LinkProgram": err = gl_LinkProgram((*argsCopy)[0])
	case "gl.Clear": err = gl_Clear((*argsCopy)[0])
	case "gl.UseProgram": err = gl_UseProgram((*argsCopy)[0])
	case "gl.BindBuffer": err = gl_BindBuffer((*argsCopy)[0], (*argsCopy)[1])
	case "gl.BindVertexArray": err = gl_BindVertexArray((*argsCopy)[0])
	case "gl.EnableVertexAttribArray": err = gl_EnableVertexAttribArray((*argsCopy)[0])
	case "gl.VertexAttribPointer": err = gl_VertexAttribPointer((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3], (*argsCopy)[4])
	case "gl.DrawArrays": err = gl_DrawArrays((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
	case "gl.GenBuffers": err = gl_GenBuffers((*argsCopy)[0], (*argsCopy)[1])
	case "gl.BufferData": err = gl_BufferData((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3])
	case "gl.GenVertexArrays": err = gl_GenVertexArrays((*argsCopy)[0], (*argsCopy)[1])
	case "gl.CreateShader": err = gl_CreateShader((*argsCopy)[0], expr, call)
	case "gl.Strs": err = gl_Strs((*argsCopy)[0], (*argsCopy)[1])
	case "gl.Free": err = gl_Free((*argsCopy)[0])
	case "gl.ShaderSource": err = gl_ShaderSource((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
	case "gl.CompileShader": err = gl_CompileShader((*argsCopy)[0])
	case "gl.GetShaderiv": err = gl_GetShaderiv((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
	case "gl.AttachShader": err = gl_AttachShader((*argsCopy)[0], (*argsCopy)[1])
	case "gl.LoadIdentity": err = gl_LoadIdentity()
	case "gl.MatrixMode": err = gl_MatrixMode((*argsCopy)[0])
	case "gl.EnableClientState": err = gl_EnableClientState((*argsCopy)[0])
	case "gl.PushMatrix": err = gl_PushMatrix()
	case "gl.PopMatrix": err = gl_PopMatrix()
	case "gl.Rotatef": err = gl_Rotatef((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3])
	case "gl.Translatef": err = gl_Translatef((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
	case "gl.NewTexture": err = gl_NewTexture((*argsCopy)[0], expr, call)

	case "gl.BindTexture": err = gl_BindTexture((*argsCopy)[0], (*argsCopy)[1])
	case "gl.Color3f": err = gl_Color3f((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
	case "gl.Color4f": err = gl_Color4f((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3])
	case "gl.Begin": err = gl_Begin((*argsCopy)[0])
	case "gl.End": err = gl_End()
	case "gl.Normal3f": err = gl_Normal3f((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
	case "gl.TexCoord2f": err = gl_TexCoord2f((*argsCopy)[0], (*argsCopy)[1])
	case "gl.Vertex2f": err = gl_Vertex2f((*argsCopy)[0], (*argsCopy)[1])
	case "gl.Vertex3f": err = gl_Vertex3f((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
	case "gl.Hint": err = gl_Hint((*argsCopy)[0], (*argsCopy)[1])
	case "gl.Ortho": err = gl_Ortho((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3], (*argsCopy)[4], (*argsCopy)[5])
	case "gl.Viewport": err = gl_Viewport((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3])

	case "gl.Enable": err = gl_Enable((*argsCopy)[0])
	case "gl.Disable": err = gl_Enable((*argsCopy)[0])
	case "gl.ClearColor": err = gl_ClearColor((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3])
	case "gl.ClearDepth": err = gl_ClearDepth((*argsCopy)[0])
	case "gl.DepthFunc": err = gl_DepthFunc((*argsCopy)[0])
	case "gl.DepthMask": err = gl_DepthMask((*argsCopy)[0])
	case "gl.BlendFunc": err = gl_BlendFunc((*argsCopy)[0], (*argsCopy)[1])
	case "gl.TexCoord2d": err = gl_TexCoord2d((*argsCopy)[0], (*argsCopy)[1])
	case "gl.Lightfv": err = gl_Lightfv((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
	case "gl.TexEnvi": err = gl_TexEnvi((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
	case "gl.Scalef": err = gl_Scalef((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
	case "gl.Frustum": err = gl_Frustum((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3], (*argsCopy)[4], (*argsCopy)[5])
		// GLFW
	case "glfw.Init": err = glfw_Init()
	case "glfw.WindowHint": err = glfw_WindowHint((*argsCopy)[0], (*argsCopy)[1])
	case "glfw.CreateWindow": err = glfw_CreateWindow((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3])
	case "glfw.MakeContextCurrent": err = glfw_MakeContextCurrent((*argsCopy)[0])
	case "glfw.ShouldClose": err = glfw_ShouldClose((*argsCopy)[0], expr, call)
	case "glfw.PollEvents": err = glfw_PollEvents()
	case "glfw.SwapBuffers": err = glfw_SwapBuffers((*argsCopy)[0])
	case "glfw.GetTime": err = glfw_GetTime(expr, call)
	case "glfw.GetFramebufferSize": err = glfw_GetFramebufferSize((*argsCopy)[0], expr, call)
	case "glfw.SetKeyCallback": err = glfw_SetKeyCallback((*argsCopy)[0], (*argsCopy)[1], expr, call)
		// Operating System
	case "os.Create": err = os_Create((*argsCopy)[0])
	case "os.Open": err = os_Open((*argsCopy)[0])
	case "os.Close": err = os_Close((*argsCopy)[0])
	case "os.GetWorkingDirectory": err = os_GetWorkingDirectory(expr, call)
	case "":
	}

	// there was an error and we'll report line number and err msg
	if err != nil {
		*exc = true
		*excError = errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err))
	}
}

func resolveStructField (fld string, val *[]byte, strct *CXStruct) ([]byte, string, int32, int32) {
	var offset int32 = 0
	for _, f := range strct.Fields {

		var fldType string
		
		isArray := false
		isBasic := false
		if f.Typ[:2] == "[]" {
			isArray = true
			for _, basic := range BASIC_TYPES {
				if basic == f.Typ[2:] {
					isBasic = true
					break
				}
			}
		} else {
			for _, basic := range BASIC_TYPES {
				if basic == f.Typ {
					isBasic = true
					break
				}
			}
		}

		if isBasic {
			fldType = f.Typ
		} else {
			if isArray {
				fldType = "[]"
			} else {
				fldType = "struct"
			}
		}
		
		if f.Name == fld {
			var size int32
			
			switch fldType {
			case "byte":
				size = 1
			case "bool", "i32", "f32":
				size = 4
			case "i64", "f64":
				size = 8
			case "[]str":
				var noElms int32
				encoder.DeserializeAtomic((*val)[offset:offset+4], &noElms)

				noSize := (*val)[offset+4:]
				
				var subOffset int32
				for c := 0; c < int(noElms); c++ {
					var strSize int32
					encoder.DeserializeRaw(noSize[subOffset:subOffset+4], &strSize)
					subOffset += strSize + 4
				}
				size = subOffset

				return (*val)[offset:offset+size + 4], f.Typ, offset, size + 4
			case "str", "[]byte":
				var arrOffset int32
				encoder.DeserializeAtomic((*val)[offset:offset+4], &arrOffset)
				size = arrOffset

				return (*val)[offset:offset+size + 4], f.Typ, offset, size + 4
			case "[]bool", "[]i32", "[]f32":
				var arrOffset int32
				encoder.DeserializeAtomic((*val)[offset:offset+4], &arrOffset)
				size = arrOffset
				
				return (*val)[offset:offset+(size * 4) + 4], f.Typ, offset, (size * 4) + 4
			case "[]i64", "[]f64":
				var arrOffset int32
				encoder.DeserializeAtomic((*val)[offset:offset+4], &arrOffset)
				size = arrOffset
				
				return (*val)[offset:offset+(size * 8) + 4], f.Typ, offset, (size * 8) + 4
			case "[]":
				if strct, err := strct.Context.GetStruct(f.Typ[2:], strct.Module.Name); err == nil {
					lastFld := strct.Fields[len(strct.Fields) - 1]
					instances := (*val)[offset+4:]

					var upperBound int32
					var size int32
					encoder.DeserializeAtomic((*val)[offset:offset + 4], &size)
					
					if size == 0 {
						return (*val)[offset:offset+4], f.Typ, offset, 4
					}

					for c := int32(0); c < size; c++ {
						subArray := instances[upperBound:]
						_, _, off, size := resolveStructField(lastFld.Name, &subArray, strct)
						
						upperBound = upperBound + off + size
					}

					return (*val)[offset:offset + upperBound + 4], f.Typ, offset, upperBound + 4
				}
			case "struct":
				if strct, err := strct.Context.GetStruct(f.Typ, strct.Module.Name); err == nil {
					lastFld := strct.Fields[len(strct.Fields) - 1]

					instances := (*val)[offset:]
					_, _, off, size := resolveStructField(lastFld.Name, &instances, strct)
					
					return (*val)[offset:offset + off + size], f.Typ, offset, off + size
				}
			}
			return (*val)[offset:offset+size], f.Typ, offset, size
		}
		
		switch fldType {
		case "byte":
			offset += 1
		case "bool", "i32", "f32":
			offset += 4
		case "i64", "f64":
			offset += 8
		case "[]str":
			var noElms int32
			encoder.DeserializeAtomic((*val)[offset:offset+4], &noElms)

			noSize := (*val)[offset+4:]

			var subOffset int32
			for c := 0; c < int(noElms); c++ {
				var strSize int32
				encoder.DeserializeRaw(noSize[subOffset:subOffset+4], &strSize)
				subOffset += strSize + 4
			}
			offset += subOffset + 4
		case "str", "[]byte":
			var arrOffset int32
			encoder.DeserializeAtomic((*val)[offset:offset+4], &arrOffset)
			offset += arrOffset + 4
		case "[]bool", "[]i32", "[]f32":
			var arrOffset int32
			encoder.DeserializeAtomic((*val)[offset:offset+4], &arrOffset)
			
			offset += (arrOffset * 4) + 4
		case "[]i64", "[]f64":
			var arrOffset int32
			encoder.DeserializeAtomic((*val)[offset:offset+4], &arrOffset)

			offset += (arrOffset * 8) + 4
		case "[]":
			if strct, err := strct.Context.GetStruct(f.Typ[2:], strct.Module.Name); err == nil {
				instances := (*val)[offset+4:]
				lastFld := strct.Fields[len(strct.Fields) - 1]
				
				var upperBound int32
				
				var size int32
				encoder.DeserializeAtomic((*val)[offset:offset+4], &size)

				// we don't need this. if size == 0, the loop won't execute
				// and we'll return lowerBound(0) + 4 = 4
				// if size == 0 {
				// 	offset += 4
				// }
				
				for c := int32(0); c < size; c++ {
					subArray := instances[upperBound:]
					_, _, off, size := resolveStructField(lastFld.Name, &subArray, strct)

					upperBound = upperBound + off + size
				}
				offset += upperBound + 4
			}
		case "struct":
			if strct, err := strct.Context.GetStruct(f.Typ, strct.Module.Name); err == nil {
				lastFld := strct.Fields[len(strct.Fields) - 1]

				instances := (*val)[offset:]
				_, _, off, size := resolveStructField(lastFld.Name, &instances, strct)

				offset += off + size
			}
		}
	}
	
	return nil, "", 0, 0
}

func resolveArrayIndex (index int, val *[]byte, typ string) ([]byte, string) {
	switch typ {
	case "[]byte":
		return (*val)[index+4:(index+1)+4], "byte"
	case "[]bool":
		return (*val)[(index+1)*4:(index+2)*4], "bool"
	case "[]i32":
		return (*val)[(index+1)*4:(index+2)*4], "i32"
	case "[]i64":
		return (*val)[((index)*8)+4:((index+1)*8)+4], "i64"
	case "[]f32":
		return (*val)[(index+1)*4:(index+2)*4], "f32"
	case "[]f64":
		return (*val)[((index)*8)+4:((index+1)*8)+4], "f64"
	}
	
	return nil, ""
}

func resolveIdent (lookingFor string, call *CXCall) (*CXArgument, error) {
	var resolvedIdent *CXDefinition
	
	isStructFld := false
	isArray := false

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
				return nil, errors.New(fmt.Sprintf("module '%s' was not imported or does not exist", mod.Name))
			}
		} else {
			// then it's a global struct
			mod := call.Operator.Module
			//if def, err := mod.GetDefinition(concat(identParts[:]...)); err == nil {
			if def, err := mod.GetDefinition(identParts[0]); err == nil {
				isStructFld = true
				//resolvedIdent = def
				if strct, err := mod.Context.GetStruct(def.Typ, mod.Name); err == nil {
					byts, typ, _, _ := resolveStructField(identParts[1], def.Value, strct)
					arg := MakeArgument(&byts, typ)
					return arg, nil
					
				} else {
					return nil, err
				}
			} else {
				// then it's a local struct
				isStructFld = true

				for _, stateDef := range call.State {
					if stateDef.Name == identParts[0] {
						if strct, err := mod.Context.GetStruct(stateDef.Typ, mod.Name); err == nil {
							byts, typ, _, _ := resolveStructField(identParts[1], stateDef.Value, strct)
							arg := MakeArgument(&byts, typ)
							return arg, nil
							
						} else {
							return nil, err
						}
					}
				}
			}
		}
	} else {
		// then it's a local or global definition
		local := false
		arrayParts := strings.Split(lookingFor, "[")
		if len(arrayParts) > 1 {
			lookingFor = arrayParts[0]
		}
		for _, stateDef := range call.State {
			if stateDef.Name == arrayParts[0] {
				local = true
				resolvedIdent = stateDef
				break
			}
		}

		if len(arrayParts) > 1 && local {
			if idx, err := strconv.ParseInt(arrayParts[1], 10, 64); err == nil {
				isArray = true
				byts, typ := resolveArrayIndex(int(idx), resolvedIdent.Value, resolvedIdent.Typ)
				arg := MakeArgument(&byts, typ)
				return arg, nil
			} else {
				//excError = err
				return nil, err
			}
		}

		if !local {
			mod := call.Operator.Module
			if def, err := mod.GetDefinition(lookingFor); err == nil {
				resolvedIdent = def
			}
		}
	}

	if resolvedIdent == nil && !isStructFld && !isArray {
		return nil, errors.New(fmt.Sprintf("'%s' is undefined", lookingFor))
	}
	
	if resolvedIdent != nil && !isStructFld && !isArray {
		// if it was a struct field, we already created the argument above for efficiency reasons
		// the same goes to arrays in the form ident[index]
		arg := MakeArgument(resolvedIdent.Value, resolvedIdent.Typ)
		return arg, nil
	}
	return nil, errors.New(fmt.Sprintf("identifier '%s' could not be resolved", lookingFor))
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

	// exceptions
	var exc bool
	var excError error

	if call.Line >= len(call.Operator.Expressions) || call.Line < 0 {
		// popping the stack
		call.Context.CallStack.Calls = call.Context.CallStack.Calls[:len(call.Context.CallStack.Calls) - 1]
		numOutputs := len(call.Operator.Outputs)
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
						if i == numOutputs {
							return call.ReturnAddress.call(withDebug, nCalls, callCounter)
						}
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

			if len(argsRefs) != len(expr.Operator.Inputs) {
				
				if len(argsRefs) == 1 {
					return errors.New(fmt.Sprintf("%d: %s: expected %d arguments; %d was provided",
						expr.FileLine, expr.Operator.Name, len(expr.Operator.Inputs), len(argsRefs)))
				} else {
					return errors.New(fmt.Sprintf("%d: %s: expected %d arguments; %d were provided",
						expr.FileLine, expr.Operator.Name, len(expr.Operator.Inputs), len(argsRefs)))
				}
			}
			
			// we don't want to modify by reference, we need to make copies
			for i := 0; i < len(argsRefs); i++ {
				// if argsRefs[i].Typ == "str" {
				// 	fmt.Println(argsRefs[i].Value)
				// }

				
				if argsRefs[i].Typ == "ident" {
					var lookingFor string
					encoder.DeserializeRaw(argsRefs[i].Value, &lookingFor)
					if arg, err := resolveIdent(lookingFor, call); err == nil {
						argsCopy[i] = arg
					} else {
						return errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err.Error()))
					}
				} else {
					argsCopy[i] = argsRefs[i]
				}

				// checking if arguments types match with expressions required types
				if len(expr.Operator.Inputs) > 0 &&
					expr.Operator.Inputs[i].Typ !=
					argsCopy[i].Typ {
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

			// check if struct array function
			//fmt.Println("here", opName)

			if isNative {
				checkNative(opName, expr, call, &argsCopy, &exc, &excError)
				if exc && isTesting {
					isErrorPresent = true
				}
				if exc && !isTesting {
					fmt.Println()
					fmt.Println("Call's State:")
					for _, def := range call.State {
						isBasic := false
						for _, basic := range BASIC_TYPES {
							if basic == def.Typ {
								isBasic = true
								break
							}
						}

						if len(def.Name) > len(NON_ASSIGN_PREFIX) && def.Name[:len(NON_ASSIGN_PREFIX)] != NON_ASSIGN_PREFIX {
							if isBasic {
								fmt.Printf("%s:\t\t%s\n", def.Name, PrintValue(def.Name, def.Value, def.Typ, call.Context))
							} else {
								fmt.Println(def.Name)
								PrintValue(def.Name, def.Value, def.Typ, call.Context)
							}
						}
					}
					fmt.Println()
					fmt.Printf("%s() Arguments:\n", expr.Operator.Name)
					for i, arg := range argsCopy {
						fmt.Printf("%d: %s\n", i, PrintValue("", arg.Value, arg.Typ, call.Context))
					}
					fmt.Println()
					return excError
				}
				
				call.Line++
				return call.call(withDebug, nCalls, callCounter)
			} else {
				// operator was not a native function
				if exc && isTesting {
					isErrorPresent = true
					//fmt.Println(excError)
				}
				if exc && !isTesting {
					fmt.Println()
					fmt.Println("Call's State:")
					for _, def := range call.State {
						fmt.Printf("%s:\t\t%s\n", def.Name, PrintValue(def.Name, def.Value, def.Typ, call.Context))
					}
					fmt.Println()
					fmt.Printf("%s() Arguments:\n", expr.Operator.Name)
					for i, arg := range argsCopy {
						fmt.Printf("%d: %s\n", i, PrintValue("", arg.Value, arg.Typ, call.Context))
					}
					fmt.Println()
					return excError
				}
				
				call.Line++ // once the subcall finishes, call next line
				if argDefs, err := argsToDefs(argsCopy, expr.Operator.Inputs, expr.Operator.Outputs, call.Module, call.Context); err == nil {
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
