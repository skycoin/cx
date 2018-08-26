package base

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"io/ioutil"
	"math/rand"
	// "runtime"
	"time"
)

func callsEqual(call1, call2 *CXCall) bool {
	if call1.Line != call2.Line ||
		len(call1.State) != len(call2.State) ||
		call1.Operator != call2.Operator ||
		call1.ReturnAddress != call2.ReturnAddress ||
		call1.Package != call2.Package {
		return false
	}

	for k, v := range call1.State {
		if call2.State[k] != v {
			return false
		}
	}

	return true
}

func saveStep(call *CXCall) {
	lenCallStack := len(call.Program.CallStack)
	newStep := MakeCallStack(lenCallStack)
	// newStep := make([][]CXCall, 0)

	if len(call.Program.Steps) < 1 {
		// First call, copy everything
		for i, call := range call.Program.CallStack {
			newStep[i] = *MakeCallCopy(&call, call.Package, call.Program)
		}

		call.Program.Steps = append(call.Program.Steps, newStep)
		return
	}

	lastStep := call.Program.Steps[len(call.Program.Steps)-1]
	lenLastStep := len(lastStep)

	smallerLen := 0
	if lenLastStep < lenCallStack {
		smallerLen = lenLastStep
	} else {
		smallerLen = lenCallStack
	}

	// Everytime a call changes, we need to make a hard copy of it
	// If the call doesn't change, we keep saving a pointer to it

	for i, call := range call.Program.CallStack[:smallerLen] {
		if callsEqual(&call, &lastStep[i]) {
			// if they are equal
			// append reference
			newStep[i] = lastStep[i]
		} else {
			newStep[i] = *MakeCallCopy(&call, call.Package, call.Program)
		}
	}

	// sizes can be different. if this is the case, we hard copy the rest
	for i, call := range call.Program.CallStack[smallerLen:] {
		newStep[i+smallerLen] = *MakeCallCopy(&call, call.Package, call.Program)
	}

	call.Program.Steps = append(call.Program.Steps, newStep)
	return
}

// It "un-runs" a program
func (prgrm *CXProgram) Reset() {
	prgrm.CallStack = MakeCallStack(0)
	prgrm.Steps = make([][]CXCall, 0)
	prgrm.Outputs = make([]*CXArgument, 0)
	//prgrm.ProgramSteps = nil
}

func (prgrm *CXProgram) ResetTo(stepNumber int) {
	// if no steps, we do nothing. the program will run from step 0
	if len(prgrm.Steps) > 0 {
		if stepNumber > len(prgrm.Steps) {
			stepNumber = len(prgrm.Steps) - 1
		}
		reqStep := prgrm.Steps[stepNumber]

		newStep := MakeCallStack(len(reqStep))

		var lastCall *CXCall
		for j, call := range reqStep {
			newCall := *MakeCallCopy(&call, call.Package, call.Program)
			newCall.ReturnAddress = lastCall
			lastCall = &newCall
			newStep[j] = newCall
		}

		prgrm.CallStack = newStep
		prgrm.Steps = prgrm.Steps[:stepNumber]
	}
}

// func (prgrm *CXProgram) UnRun(nCalls int) {
// 	if len(prgrm.Steps) > 0 && nCalls > 0 {
// 		if nCalls > len(prgrm.Steps) {
// 			nCalls = len(prgrm.Steps) - 1
// 		}

// 		reqStep := prgrm.Steps[len(prgrm.Steps)-nCalls]

// 		newStep := MakeCallStack(len(reqStep))

// 		var lastCall *CXCall
// 		for j, call := range reqStep {
// 			newCall := *MakeCallCopy(&call, call.Package, call.Program)
// 			newCall.ReturnAddress = lastCall
// 			lastCall = &newCall
// 			newStep[j] = newCall
// 		}

// 		prgrm.CallStack = newStep
// 		prgrm.Steps = prgrm.Steps[:len(prgrm.Steps)-nCalls]
// 	}
// }

func (prgrm *CXProgram) UnRun(nCalls int) {
	if nCalls >= 0 || prgrm.CallCounter < 0 {
		return
	}

	call := &prgrm.CallStack[prgrm.CallCounter]
	
	for c := nCalls; c < 0; c++ {
		if call.Line >= c {
			// then we stay in this call counter
			call.Line += c
			c -= c
		} else {
			
			if prgrm.CallCounter == 0 {
				call.Line = 0
				return
			}
			c += call.Line
			call.Line = 0
			prgrm.CallCounter--
			call = &prgrm.CallStack[prgrm.CallCounter]
		}
	}
}

// Compiling from CXGO to CX Base
func (prgrm *CXProgram) Compile(withProfiling bool) {
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

var prgrm = MakeContext();var mod *CXModule;var imp *CXModule;var fn *CXFunction;var op *CXFunction;var expr *CXExpression;var strct *CXStruct;var fld *CXField;var arg *CXArgument;var tag string = "";

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
		program.WriteString(`package main;import (. github.com/skycoin/cx/src/base"; "runtime";);var prgrm = MakeContext();var mod *CXModule;var imp *CXModule;var fn *CXFunction;var op *CXFunction;var expr *CXExpression;var strct *CXStruct;var arg *CXArgument;var tag string = "";func main () {runtime.LockOSThread();`)
	}

	for _, mod := range prgrm.Packages {
		program.WriteString(fmt.Sprintf(`mod = MakeModule("%s");prgrm.AddModule(mod);%s`, mod.Name, asmNL))
		for _, imp := range mod.Imports {
			program.WriteString(fmt.Sprintf(`imp, _ = prgrm.GetModule("%s");mod.AddImport(imp);%s`, imp.Name, asmNL))
		}

		for _, fn := range mod.Functions {
			isUsed := false
			if fn.Name != MAIN_FUNC {
				for _, mod := range prgrm.Packages {
					for _, chkFn := range mod.Functions {
						for _, expr := range chkFn.Expressions {
							if expr.Operator.Name == fn.Name {
								isUsed = true
								break
							}
						}
						if isUsed {
							break
						}
					}
				}
			} else {
				isUsed = true
			}

			if !isUsed {
				continue
			}

			program.WriteString(fmt.Sprintf(`fn = MakeFunction(%#v);mod.AddFunction(fn);%s`, fn.Name, asmNL))

			for _, inp := range fn.Inputs {
				program.WriteString(fmt.Sprintf(`fn.AddInput(MakeParameter("%s", "%s"));%s`, inp.Name, inp.Typ, asmNL))
			}

			for _, out := range fn.Outputs {
				program.WriteString(fmt.Sprintf(`fn.AddOutput(MakeParameter("%s", "%s"));%s`, out.Name, out.Typ, asmNL))
			}

			// var optExpressions []*CXExpression
			// for _, expr := range fn.Expressions {
			// 	if expr.Operator.Name == "identity" {
			// 		var nonAssignIdent string
			// 		encoder.DeserializeRaw(*expr.Inputs[0].Value, &nonAssignIdent)

			// 		for _, idExpr := range fn.Expressions {
			// 			for i, out := range idExpr.Outputs {
			// 				if out.Name == nonAssignIdent {
			// 					idExpr.Outputs[i] = expr.Outputs[0]
			// 					break
			// 				}
			// 			}
			// 		}
			// 		continue
			// 	}
			// 	optExpressions = append(optExpressions, expr)
			// }

			// fn.Expressions = optExpressions

			//for _, expr := range optExpressions {
			for _, expr := range fn.Expressions {
				var tagStr string
				if expr.Label != "" {
					tagStr = fmt.Sprintf(`expr.Label = "%s";`, expr.Label)
				}
				program.WriteString(fmt.Sprintf(`op, _ = prgrm.GetFunction("%s", "%s");expr = MakeExpression(op);expr.FileLine = %d;fn.AddExpression(expr);%s%s`,
					expr.Operator.Name, expr.Operator.Package.Name, expr.FileLine, tagStr, asmNL))

				for _, arg := range expr.Inputs {
					program.WriteString(fmt.Sprintf(`expr.AddArgument(MakeArgument(&%#v, "%s"));%s`, *arg.Value, arg.Typ, asmNL))
				}

				for _, outName := range expr.Outputs {
					program.WriteString(fmt.Sprintf(`expr.AddOutputName("%s");%s`, outName.Name, asmNL))
				}
			}

		}

		for _, strct := range mod.Structs {
			program.WriteString(fmt.Sprintf(`strct = MakeStruct("%s");mod.AddStruct(strct);%s`, strct.Name, asmNL))
			for _, fld := range strct.Fields {
				program.WriteString(fmt.Sprintf(`fld = MakeField("%s", "%s");strct.AddField(fld);%s`, fld.Name, fld.Typ, asmNL))
			}
		}

		for _, def := range mod.Globals {
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

	program.WriteString(`prgrm.Run(false, -1);}`)
	ioutil.WriteFile(fmt.Sprintf("o.go"), []byte(program.String()), 0644)
}

var isTesting bool
var isErrorPresent bool

// func checkNative(opName string, expr *CXExpression, call *CXCall, argsCopy *[]*CXArgument, exc *bool, excError *error) {
// 	var err error
// 	switch opName {
// 	// case "serialize": serialize_program(expr, call)
// 	// case "deserialize":
// 	// 	// it only prints the deserialized program for now
// 	// 	Deserialize((*argsCopy)[0].Value).PrintProgram(false)
// 	case "evolve":
// 		var fnName string
// 		var fnBag string
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &fnName)
// 		encoder.DeserializeRaw(*(*argsCopy)[1].Value, &fnBag)

// 		var inps []float64
// 		encoder.DeserializeRaw(*(*argsCopy)[2].Value, &inps)

// 		var outs []float64
// 		encoder.DeserializeRaw(*(*argsCopy)[3].Value, &outs)

// 		var numberExprs int32
// 		encoder.DeserializeRaw(*(*argsCopy)[4].Value, &numberExprs)
// 		var iterations int32
// 		encoder.DeserializeRaw(*(*argsCopy)[5].Value, &iterations)
// 		var epsilon float64
// 		encoder.DeserializeRaw(*(*argsCopy)[6].Value, &epsilon)

// 		// err = call.Program.Evolve(fnName, fnBag, inps, outs, int(numberExprs), int(iterations), epsilon, expr, call)
// 		// flow control
// 	case "jmp":
// 		jmp(expr, &call.Program.Stacks[0], 0, call)

// 	case "baseGoTo":
// 		err = baseGoTo(call, (*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
// 	case "goTo":
// 		err = goTo(call, (*argsCopy)[0])
// 		// I/O functions
// 	case "bool.print":
// 		var val bool
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
// 		fmt.Println(val)
// 	case "str.print":
// 		var val string
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
// 		fmt.Println(val)
// 	case "byte.print":
// 		fmt.Println((*(*argsCopy)[0].Value)[0])
// 	case "i32.print":
// 		var val int32
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
// 		fmt.Println(val)
// 	case "i64.print":
// 		var val int64
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
// 		fmt.Println(val)
// 	case "f32.print":
// 		var val float32
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
// 		fmt.Println(val)
// 	case "f64.print":
// 		var val float64
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
// 		fmt.Println(val)
// 	case "[]bool.print":
// 		var val []int32
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
// 		fmt.Print("[")
// 		for i, v := range val {
// 			if v == 0 {
// 				fmt.Print("false")
// 			} else {
// 				fmt.Print("true")
// 			}
// 			if i != len(val)-1 {
// 				fmt.Print(" ")
// 			}
// 		}
// 		fmt.Print("]")
// 		fmt.Println()
// 	case "[]byte.print":
// 		var val []byte
// 		//fmt.Println(*(*argsCopy)[0].Value)
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
// 		fmt.Println(val)
// 	case "[]str.print":
// 		var val []string
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
// 		fmt.Println(val)
// 	case "[]i32.print":
// 		var val []int32
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
// 		fmt.Println(val)
// 	case "[]i64.print":
// 		var val []int64
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
// 		fmt.Println(val)
// 	case "[]f32.print":
// 		var val []float32
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
// 		fmt.Println(val)
// 	case "[]f64.print":
// 		var val []float64
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &val)
// 		fmt.Println(val)
// 		// identity functions
// 	case "str.id", "bool.id", "byte.id", "i32.id", "i64.id", "f32.id", "f64.id", "[]bool.id", "[]byte.id", "[]str.id", "[]i32.id", "[]i64.id", "[]f32.id", "[]f64.id":
// 		assignOutput(0, *(*argsCopy)[0].Value, (*argsCopy)[0].Typ, expr, call)
// 	case "identity":
// 		assignOutput(0, *(*argsCopy)[0].Value, (*argsCopy)[0].Typ, expr, call)
// 		// identity((*argsCopy)[0], expr, call)
// 		// cast functions
// 	case "[]byte.str", "byte.str", "bool.str", "i32.str", "i64.str", "f32.str", "f64.str":
// 		err = castToStr((*argsCopy)[0], expr, call)
// 	case "str.[]byte":
// 		err = castToByteA((*argsCopy)[0], expr, call)
// 	case "i32.byte", "i64.byte", "f32.byte", "f64.byte":
// 		err = castToByte((*argsCopy)[0], expr, call)
// 	case "byte.i32", "i64.i32", "f32.i32", "f64.i32":
// 		err = castToI32((*argsCopy)[0], expr, call)
// 	case "byte.i64", "i32.i64", "f32.i64", "f64.i64":
// 		err = castToI64((*argsCopy)[0], expr, call)
// 	case "byte.f32", "i32.f32", "i64.f32", "f64.f32":
// 		err = castToF32((*argsCopy)[0], expr, call)
// 	case "byte.f64", "i32.f64", "i64.f64", "f32.f64":
// 		err = castToF64((*argsCopy)[0], expr, call)
// 	case "[]i32.[]byte", "[]i64.[]byte", "[]f32.[]byte", "[]f64.[]byte":
// 		err = castToByteA((*argsCopy)[0], expr, call)
// 	case "[]byte.[]i32", "[]i64.[]i32", "[]f32.[]i32", "[]f64.[]i32":
// 		err = castToI32A((*argsCopy)[0], expr, call)
// 	case "[]byte.[]i64", "[]i32.[]i64", "[]f32.[]i64", "[]f64.[]i64":
// 		err = castToI64A((*argsCopy)[0], expr, call)
// 	case "[]byte.[]f32", "[]i32.[]f32", "[]i64.[]f32", "[]f64.[]f32":
// 		err = castToF32A((*argsCopy)[0], expr, call)
// 	case "[]byte.[]f64", "[]i32.[]f64", "[]i64.[]f64", "[]f32.[]f64":
// 		err = castToF64A((*argsCopy)[0], expr, call)
// 		// logical operators
// 	case "and":
// 		err = and((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "or":
// 		err = or((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "not":
// 		err = not((*argsCopy)[0], expr, call)
// 		// relational operators
// 	case "bool.eq":
// 		err = eqBool((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "bool.uneq":
// 		err = uneqBool((*argsCopy)[0], (*argsCopy)[1], expr, call)

// 		// undefined type operators
// 	case "gt":
// 		err = gtUnd((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "len":
// 		err = lenUnd((*argsCopy)[0], expr, call)

// 	case "i32.lt":
// 		err = ltI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.gt":
// 		err = gtI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.eq":
// 		err = eqI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.uneq":
// 		err = uneqI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.lteq":
// 		err = lteqI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.gteq":
// 		err = gteqI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.lt":
// 		err = ltI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.gt":
// 		err = gtI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.eq":
// 		err = eqI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.uneq":
// 		err = uneqI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.lteq":
// 		err = lteqI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.gteq":
// 		err = gteqI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f32.lt":
// 		err = ltF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f32.gt":
// 		err = gtF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f32.eq":
// 		err = eqF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f32.uneq":
// 		err = uneqF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f32.lteq":
// 		err = lteqF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f32.gteq":
// 		err = gteqF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f64.lt":
// 		err = ltF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f64.gt":
// 		err = gtF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f64.eq":
// 		err = eqF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f64.uneq":
// 		err = uneqF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f64.lteq":
// 		err = lteqF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f64.gteq":
// 		err = gteqF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "str.lt":
// 		err = ltStr((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "str.gt":
// 		err = gtStr((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "str.eq":
// 		err = eqStr((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "str.uneq":
// 		err = uneqStr((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "str.lteq":
// 		err = lteqStr((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "str.gteq":
// 		err = gteqStr((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "byte.lt":
// 		err = ltByte((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "byte.gt":
// 		err = gtByte((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "byte.eq":
// 		err = eqByte((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "byte.uneq":
// 		err = uneqByte((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "byte.lteq":
// 		err = lteqByte((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "byte.gteq":
// 		err = gteqByte((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 		// io functions
// 	case "str.read":
// 		err = readStr(expr, call)
// 	case "i32.read":
// 		err = readI32(expr, call)
// 		// struct operations
// 	case "initDef":
// 		err = initDef((*argsCopy)[0], expr, call)
// 		// arithmetic functions
// 	case "i32.add":
// 		err = addI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.mul":
// 		err = mulI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.sub":
// 		err = subI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.div":
// 		err = divI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.pow":
// 		err = powI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.sqrt":
// 		err = sqrtI32((*argsCopy)[0], expr, call)
// 	case "i64.add":
// 		err = addI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.mul":
// 		err = mulI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.sub":
// 		err = subI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.div":
// 		err = divI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.pow":
// 		err = powI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.sqrt":
// 		err = sqrtI64((*argsCopy)[0], expr, call)
// 	case "f32.add":
// 		err = addF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f32.mul":
// 		err = mulF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f32.sub":
// 		err = subF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f32.div":
// 		err = divF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f32.pow":
// 		err = powF32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f32.sqrt":
// 		err = sqrtF32((*argsCopy)[0], expr, call)
// 	case "f32.cos":
// 		err = cosF32((*argsCopy)[0], expr, call)
// 	case "f32.sin":
// 		err = sinF32((*argsCopy)[0], expr, call)
// 	case "f64.add":
// 		err = addF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f64.mul":
// 		err = mulF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f64.sub":
// 		err = subF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f64.div":
// 		err = divF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f64.pow":
// 		err = powF64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "f64.sqrt":
// 		err = sqrtF64((*argsCopy)[0], expr, call)
// 	case "f64.cos":
// 		err = cosF64((*argsCopy)[0], expr, call)
// 	case "f64.sin":
// 		err = sinF64((*argsCopy)[0], expr, call)
// 	case "i32.abs":
// 		err = absI32((*argsCopy)[0], expr, call)
// 	case "i64.abs":
// 		err = absI64((*argsCopy)[0], expr, call)
// 	case "f32.abs":
// 		err = absF32((*argsCopy)[0], expr, call)
// 	case "f64.abs":
// 		err = absF64((*argsCopy)[0], expr, call)
// 	case "i32.mod":
// 		err = modI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.mod":
// 		err = modI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 		// bitwise operators
// 	case "i32.bitand":
// 		err = andI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.bitor":
// 		err = orI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.bitxor":
// 		err = xorI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.bitclear":
// 		err = andNotI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.bitshl":
// 		err = shiftLeftI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i32.bitshr":
// 		err = shiftRightI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.bitand":
// 		err = andI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.bitor":
// 		err = orI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.bitxor":
// 		err = xorI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.bitclear":
// 		err = andNotI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.bitshl":
// 		err = shiftLeftI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.bitshr":
// 		err = shiftRightI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 		// make functions
// 	case "[]bool.make":
// 		err = makeArray("[]bool", (*argsCopy)[0], expr, call)
// 	case "[]byte.make":
// 		err = makeArray("[]byte", (*argsCopy)[0], expr, call)
// 	case "[]str.make":
// 		err = makeArray("[]str", (*argsCopy)[0], expr, call)
// 	case "[]i32.make":
// 		err = makeArray("[]i32", (*argsCopy)[0], expr, call)
// 	case "[]i64.make":
// 		err = makeArray("[]i64", (*argsCopy)[0], expr, call)
// 	case "[]f32.make":
// 		err = makeArray("[]f32", (*argsCopy)[0], expr, call)
// 	case "[]f64.make":
// 		err = makeArray("[]f64", (*argsCopy)[0], expr, call)
// 		// array functions
// 	case "[]bool.read":
// 		err = readBoolA((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]bool.write":
// 		err = writeBoolA((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
// 	case "[]byte.read":
// 		err = readByteA((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]byte.write":
// 		err = writeByteA((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
// 	case "[]str.read":
// 		err = readStrA((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]str.write":
// 		err = writeStrA((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
// 	case "[]i32.read":
// 		err = readI32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]i32.write":
// 		err = writeI32A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
// 	case "[]i64.read":
// 		err = readI64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]i64.write":
// 		err = writeI64A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
// 	case "[]f32.read":
// 		err = readF32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]f32.write":
// 		err = writeF32A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
// 	case "[]f64.read":
// 		err = readF64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]f64.write":
// 		err = writeF64A((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
// 	case "[]bool.len":
// 		err = lenBoolA((*argsCopy)[0], expr, call)
// 	case "[]byte.len":
// 		err = lenByteA((*argsCopy)[0], expr, call)
// 	case "str.len":
// 		err = lenStr((*argsCopy)[0], expr, call)
// 	case "[]str.len":
// 		err = lenStrA((*argsCopy)[0], expr, call)
// 	case "[]i32.len":
// 		err = lenI32A((*argsCopy)[0], expr, call)
// 	case "[]i64.len":
// 		err = lenI64A((*argsCopy)[0], expr, call)
// 	case "[]f32.len":
// 		err = lenF32A((*argsCopy)[0], expr, call)
// 	case "[]f64.len":
// 		err = lenF64A((*argsCopy)[0], expr, call)
// 		// concatenation functions
// 	case "str.concat":
// 		err = concatStr((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]byte.append":
// 		err = appendByteA((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]bool.append":
// 		err = appendBoolA((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]str.append":
// 		err = appendStrA((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]i32.append":
// 		err = appendI32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]i64.append":
// 		err = appendI64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]f32.append":
// 		err = appendF32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]f64.append":
// 		err = appendF64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]byte.concat":
// 		err = concatByteA((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]bool.concat":
// 		err = concatBoolA((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]str.concat", "aff.concat":
// 		err = concatStrA((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]i32.concat":
// 		err = concatI32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]i64.concat":
// 		err = concatI64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]f32.concat":
// 		err = concatF32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]f64.concat":
// 		err = concatF64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 		// copy functions
// 	case "[]byte.copy":
// 		err = copyByteA((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]bool.copy":
// 		err = copyBoolA((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]str.copy":
// 		err = copyStrA((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]i32.copy":
// 		err = copyI32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]i64.copy":
// 		err = copyI64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]f32.copy":
// 		err = copyF32A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "[]f64.copy":
// 		err = copyF64A((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 		// time functions
// 	case "sleep":
// 		err = sleep((*argsCopy)[0])
// 		// utilitiy functions
// 	case "i32.rand":
// 		err = randI32((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "i64.rand":
// 		err = randI64((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 		// meta functions

// 	case "aff.query":
// 		err = aff_query((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
// 	case "aff.execute":
// 		err = aff_execute((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
// 	case "aff.index":
// 		err = aff_index((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "aff.name":
// 		err = aff_name((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "aff.print":
// 		err = aff_print((*argsCopy)[0], call)
// 	case "aff.len":
// 		err = aff_len((*argsCopy)[0], expr, call)

// 	case "rem.expr":
// 		err = rem_expr((*argsCopy)[0], call.Operator)
// 	case "rem.arg":
// 		err = rem_arg((*argsCopy)[0], call.Operator)
// 	case "add.expr":
// 		err = add_expr((*argsCopy)[0], (*argsCopy)[1], call)
// 		// testing functions
// 	case "halt":
// 		var msg string
// 		encoder.DeserializeRaw(*(*argsCopy)[0].Value, &msg)
// 		fmt.Println(msg)
// 		call.Line++
// 		*exc = true
// 		*excError = errors.New(fmt.Sprintf("%s: %d: call to halt", expr.FileName, expr.FileLine))
// 	case "test.start":
// 		isTesting = true
// 	case "test.stop":
// 		isTesting = false
// 	case "test.error":
// 		//fmt.Println(isErrorPresent)
// 		err = test_error((*argsCopy)[0], isErrorPresent, expr)
// 		isErrorPresent = false
// 		// case "test.bool", "test.byte", "test.str", "test.i32", "test.i64", "test.f32", "test.f64", "test.[]bool", "test.[]byte", "test.[]str", "test.[]i32", "test.[]f32", "test.[]f64":
// 	case "test":
// 		err = test_value((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr)
// 		// multi dimensional array functions
// 	case "mdim.append":
// 		err = mdim_append((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "mdim.read":
// 		err = mdim_read((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "mdim.write":
// 		err = mdim_write((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
// 	case "mdim.len":
// 		err = mdim_len((*argsCopy)[0], expr, call)
// 	case "mdim.make":
// 		err = mdim_make((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 		// custom types functions
// 	case "cstm.append":
// 		err = cstm_append((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "cstm.read":
// 		err = cstm_read((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "cstm.write":
// 		err = cstm_write((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], expr, call)
// 	case "cstm.len":
// 		err = cstm_len((*argsCopy)[0], expr, call)
// 	case "cstm.make":
// 		err = cstm_make((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "cstm.serialize":
// 		err = cstm_serialize((*argsCopy)[0], expr, call)
// 	case "cstm.deserialize":
// 		err = cstm_deserialize((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 		// Time
// 	case "time.Unix":
// 		time_Unix(expr, call)
// 	case "time.UnixMilli":
// 		time_UnixMilli(expr, call)
// 	case "time.UnixNano":
// 		time_UnixNano(expr, call)
// 		// Runtime
// 	case "runtime.LockOSThread":
// 		runtime.LockOSThread()
// 		// GLText
// 	case "gltext.LoadTrueType":
// 		err = gltext_LoadTrueType((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3], (*argsCopy)[4], (*argsCopy)[5])
// 	case "gltext.Printf":
// 		err = gltext_Printf((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3])
// 		// OpenGL
// 	case "gl.Init":
// 		err = gl_Init()
// 	case "gl.CreateProgram":
// 		err = gl_CreateProgram(expr, call)
// 	case "gl.LinkProgram":
// 		err = gl_LinkProgram((*argsCopy)[0])
// 	case "gl.Clear":
// 		err = gl_Clear((*argsCopy)[0])
// 	case "gl.UseProgram":
// 		err = gl_UseProgram((*argsCopy)[0])
// 	case "gl.BindBuffer":
// 		err = gl_BindBuffer((*argsCopy)[0], (*argsCopy)[1])
// 	case "gl.BindVertexArray":
// 		err = gl_BindVertexArray((*argsCopy)[0])
// 	case "gl.EnableVertexAttribArray":
// 		err = gl_EnableVertexAttribArray((*argsCopy)[0])
// 	case "gl.VertexAttribPointer":
// 		err = gl_VertexAttribPointer((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3], (*argsCopy)[4])
// 	case "gl.DrawArrays":
// 		err = gl_DrawArrays((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
// 	case "gl.GenBuffers":
// 		err = gl_GenBuffers((*argsCopy)[0], (*argsCopy)[1])
// 	case "gl.BufferData":
// 		err = gl_BufferData((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3])
// 	case "gl.GenVertexArrays":
// 		err = gl_GenVertexArrays((*argsCopy)[0], (*argsCopy)[1])
// 	case "gl.CreateShader":
// 		err = gl_CreateShader((*argsCopy)[0], expr, call)
// 	case "gl.Strs":
// 		err = gl_Strs((*argsCopy)[0], (*argsCopy)[1])
// 	case "gl.Free":
// 		err = gl_Free((*argsCopy)[0])
// 	case "gl.ShaderSource":
// 		err = gl_ShaderSource((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
// 	case "gl.CompileShader":
// 		err = gl_CompileShader((*argsCopy)[0])
// 	case "gl.GetShaderiv":
// 		err = gl_GetShaderiv((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
// 	case "gl.AttachShader":
// 		err = gl_AttachShader((*argsCopy)[0], (*argsCopy)[1])
// 	case "gl.LoadIdentity":
// 		err = gl_LoadIdentity()
// 	case "gl.MatrixMode":
// 		err = gl_MatrixMode((*argsCopy)[0])
// 	case "gl.EnableClientState":
// 		err = gl_EnableClientState((*argsCopy)[0])
// 	case "gl.PushMatrix":
// 		err = gl_PushMatrix()
// 	case "gl.PopMatrix":
// 		err = gl_PopMatrix()
// 	case "gl.Rotatef":
// 		err = gl_Rotatef((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3])
// 	case "gl.Translatef":
// 		err = gl_Translatef((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
// 	case "gl.NewTexture":
// 		err = gl_NewTexture((*argsCopy)[0], expr, call)

// 	case "gl.BindTexture":
// 		err = gl_BindTexture((*argsCopy)[0], (*argsCopy)[1])
// 	case "gl.Color3f":
// 		err = gl_Color3f((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
// 	case "gl.Color4f":
// 		err = gl_Color4f((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3])
// 	case "gl.Begin":
// 		err = gl_Begin((*argsCopy)[0])
// 	case "gl.End":
// 		err = gl_End()
// 	case "gl.Normal3f":
// 		err = gl_Normal3f((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
// 	case "gl.TexCoord2f":
// 		err = gl_TexCoord2f((*argsCopy)[0], (*argsCopy)[1])
// 	case "gl.Vertex2f":
// 		err = gl_Vertex2f((*argsCopy)[0], (*argsCopy)[1])
// 	case "gl.Vertex3f":
// 		err = gl_Vertex3f((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
// 	case "gl.Hint":
// 		err = gl_Hint((*argsCopy)[0], (*argsCopy)[1])
// 	case "gl.Ortho":
// 		err = gl_Ortho((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3], (*argsCopy)[4], (*argsCopy)[5])
// 	case "gl.Viewport":
// 		err = gl_Viewport((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3])

// 	case "gl.Enable":
// 		err = gl_Enable((*argsCopy)[0])
// 	case "gl.Disable":
// 		err = gl_Enable((*argsCopy)[0])
// 	case "gl.ClearColor":
// 		err = gl_ClearColor((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3])
// 	case "gl.ClearDepth":
// 		err = gl_ClearDepth((*argsCopy)[0])
// 	case "gl.DepthFunc":
// 		err = gl_DepthFunc((*argsCopy)[0])
// 	case "gl.DepthMask":
// 		err = gl_DepthMask((*argsCopy)[0])
// 	case "gl.BlendFunc":
// 		err = gl_BlendFunc((*argsCopy)[0], (*argsCopy)[1])
// 	case "gl.TexCoord2d":
// 		err = gl_TexCoord2d((*argsCopy)[0], (*argsCopy)[1])
// 	case "gl.Lightfv":
// 		err = gl_Lightfv((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
// 	case "gl.TexEnvi":
// 		err = gl_TexEnvi((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
// 	case "gl.Scalef":
// 		err = gl_Scalef((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
// 	case "gl.Frustum":
// 		err = gl_Frustum((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3], (*argsCopy)[4], (*argsCopy)[5])
// 		// GLFW
// 	case "glfw.Init":
// 		err = glfw_Init()
// 	case "glfw.WindowHint":
// 		err = glfw_WindowHint((*argsCopy)[0], (*argsCopy)[1])
// 	case "glfw.CreateWindow":
// 		err = glfw_CreateWindow((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2], (*argsCopy)[3])
// 	case "glfw.MakeContextCurrent":
// 		err = glfw_MakeContextCurrent((*argsCopy)[0])
// 	case "glfw.ShouldClose":
// 		err = glfw_ShouldClose((*argsCopy)[0], expr, call)
// 	case "glfw.SetShouldClose":
// 		err = glfw_SetShouldClose((*argsCopy)[0], (*argsCopy)[0])
// 	case "glfw.PollEvents":
// 		err = glfw_PollEvents()
// 	case "glfw.SwapBuffers":
// 		err = glfw_SwapBuffers((*argsCopy)[0])
// 	case "glfw.GetKey":
// 		err = glfw_GetKey((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "glfw.GetTime":
// 		err = glfw_GetTime(expr, call)
// 	case "glfw.GetFramebufferSize":
// 		err = glfw_GetFramebufferSize((*argsCopy)[0], expr, call)
// 	case "glfw.SetKeyCallback":
// 		err = glfw_SetKeyCallback((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "glfw.SetMouseButtonCallback":
// 		err = glfw_SetMouseButtonCallback((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "glfw.SetCursorPosCallback":
// 		err = glfw_SetCursorPosCallback((*argsCopy)[0], (*argsCopy)[1], expr, call)
// 	case "glfw.GetCursorPos":
// 		err = glfw_GetCursorPos((*argsCopy)[0], expr, call)
// 	case "glfw.SetInputMode":
// 		err = glfw_SetInputMode((*argsCopy)[0], (*argsCopy)[1], (*argsCopy)[2])
// 		// Operating System
// 	case "os.Create":
// 		err = os_Create((*argsCopy)[0])
// 	case "os.ReadFile":
// 		err = os_ReadFile((*argsCopy)[0], expr, call)
// 	case "os.WriteFile":
// 		err = os_WriteFile((*argsCopy)[0], (*argsCopy)[1])
// 	case "os.Open":
// 		err = os_Open((*argsCopy)[0])
// 	case "os.Write":
// 		err = os_Write((*argsCopy)[0], (*argsCopy)[1])
// 	case "os.Close":
// 		err = os_Close((*argsCopy)[0])
// 	case "os.GetWorkingDirectory":
// 		err = os_GetWorkingDirectory(expr, call)
// 	case "":
// 	}

// 	// there was an error and we'll report line number and err msg
// 	if err != nil {
// 		*exc = true
// 		*excError = errors.New(fmt.Sprintf("%s: %d: %s", expr.FileName, expr.FileLine, err))
// 	}
// }

func (prgrm *CXProgram) RunInterpreted(withDebug bool, nCalls int) error {
	rand.Seed(time.Now().UTC().UnixNano())
	if prgrm.Terminated {
		// user wants to re-run the program
		prgrm.Terminated = false
		prgrm.CallCounter = 0
	}

	var callCounter int = 0
	// we are going to do this if the CallStack is empty
	// if prgrm.CallStack != nil && len(prgrm.CallStack) > 0 {
	if prgrm.CallStack != nil && prgrm.CallCounter > 0 {
		// we resume the program
		var lastCall *CXCall
		var err error

		var untilEnd = false
		if nCalls < 1 {
			nCalls = 1 // so the for loop executes
			untilEnd = true
		}

		for !prgrm.Terminated && nCalls > 0 && prgrm.CallCounter > 0 {
			// lastCall = &prgrm.CallStack[len(prgrm.CallStack) - 1]
			lastCall = &prgrm.CallStack[prgrm.CallCounter]
			err = lastCall.icall(withDebug, 1, callCounter)
			if err != nil {
				return err
			}
			if !untilEnd {
				nCalls = nCalls - 1
			}
		}
	} else if prgrm.CallCounter == 0 {
		// initialization and checking
		if mod, err := prgrm.SelectPackage(MAIN_PKG); err == nil {

			if fn, err := mod.SelectFunction(SYS_INIT_FUNC); err == nil {
				// *init function
				state := make([]*CXArgument, 0, 20)
				mainCall := MakeCall(fn, state, nil, mod, mod.Program)

				// prgrm.CallStack = append(prgrm.CallStack, mainCall)
				prgrm.CallStack[prgrm.CallCounter] = mainCall
				// prgrm.CallCounter++

				var err error

				for !prgrm.Terminated {
					call := &prgrm.CallStack[prgrm.CallCounter]
					err = call.icall(withDebug, 1, callCounter)
					if err != nil {
						return err
					}
				}

				// we reset call state
				prgrm.Terminated = false
				prgrm.CallCounter = 0
			} else {
				return err
			}

			if fn, err := mod.SelectFunction(MAIN_FUNC); err == nil {
				// main function
				state := make([]*CXArgument, 0, 20)
				mainCall := MakeCall(fn, state, nil, mod, mod.Program)

				// prgrm.CallStack = append(prgrm.CallStack, mainCall)
				prgrm.CallStack[prgrm.CallCounter] = mainCall
				// prgrm.CallCounter++

				var lastCall *CXCall
				var err error

				var untilEnd = false
				if nCalls < 1 {
					nCalls = 1 // so the for loop executes
					untilEnd = true
				}

				for !prgrm.Terminated && nCalls > 0 {
					// lastCall = &prgrm.CallStack[len(prgrm.CallStack) - 1]
					lastCall = &prgrm.CallStack[prgrm.CallCounter]
					err = lastCall.icall(withDebug, 1, callCounter)

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

func (call *CXCall) icall(withDebug bool, nCalls, callCounter int) error {
	// for _, arg := range call.State {
	// 	fmt.Println("exec.arg", arg.Name, arg.Value, arg.Fields)
	// 	if len(arg.Fields) > 0 {
	// 		fmt.Println("exec.flds", arg.Fields[0].Value)
	// 	}
	// }
	// fmt.Println()

	//  add a counter here to pause
	if nCalls > 0 && callCounter >= nCalls {
		return nil
	}
	callCounter++

	//saveStep(call)
	if withDebug {
		PrintCallStack(call.Program.CallStack)
	}

	// exceptions
	var exc bool
	var excError error

	if call.Line >= len(call.Operator.Expressions) || call.Line < 0 {
		// popping the stack
		// call.Program.CallStack = call.Program.CallStack[:len(call.Program.CallStack) - 1]
		call.Program.CallCounter--

		numOutputs := len(call.Operator.Outputs)
		for i, out := range call.Operator.Outputs {
			found := true
			for _, def := range call.State {
				/////////// throw error if output was not defined, or handle outputs from last expression
				if out.Name == def.Name {
					flds := def.Fields

					if call.ReturnAddress != nil {
						retName := call.ReturnAddress.Operator.Expressions[call.ReturnAddress.Line-1].Outputs[i].Name
						found := false
						for _, retDef := range call.ReturnAddress.State {
							if len(flds) > 0 {
								if sameFields(def.Fields, retDef.Fields) && retDef.Name == retName {
									retDef.Fields[len(retDef.Fields)-1].Value = def.Fields[len(def.Fields)-1].Value
									found = true
									break
								}

							} else {
								if retDef.Name == retName {
									retDef.Value = def.Value
									found = true
									break
								}
							}
						}

						if !found {
							if len(flds) > 0 {
								// arg := MakeArgument(outName)
								// arg.Fields = flds
								// arg.Fields[len(flds) - 1].AddValue(&output).AddType(typ)
								// call.State = append(call.State, arg)

								def.Name = retName
								call.ReturnAddress.State = append(call.ReturnAddress.State, def)
							} else {
								def.Name = retName
								call.ReturnAddress.State = append(call.ReturnAddress.State, def)
							}

							def.Name = retName
							call.ReturnAddress.State = append(call.ReturnAddress.State, def)
						}

						found = true
						// break
						if i == numOutputs {
							return call.ReturnAddress.icall(withDebug, nCalls, callCounter)
						}
					} else {
						// no return address. should only be for main
						call.Program.Terminated = true
						call.Program.Outputs = append(call.Program.Outputs, def)
					}
				}
			}

			// this isn't complete yet
			if !found {
				return errors.New(fmt.Sprintf("'%s' output(s) not specified", call.Operator.Name))
			}
		}

		if call.ReturnAddress != nil {
			return call.ReturnAddress.icall(withDebug, nCalls, callCounter)
		} else {
			// no return address. should only be for main
			call.Program.Terminated = true
			//call.Program.Outputs = append(call.Program.Outputs, def)
		}
	} else {
		fn := call.Operator

		if expr, err := fn.GetExpression(call.Line); err == nil {
			if expr.Operator == nil {
				// then it's a declaration
				call.State = append(call.State, expr.Outputs[0])
				call.Line++
				return call.icall(withDebug, nCalls, callCounter)
			}

			// getting arguments
			argsRefs, _ := expr.GetInputs()
			argsCopy := make([]*CXArgument, len(argsRefs))

			if len(argsRefs) != len(expr.Operator.Inputs) {
				if len(argsRefs) == 1 {
					return errors.New(fmt.Sprintf("%s: %d: %s: expected %d inputs; %d was provided",
						expr.FileName, expr.FileLine, expr.Operator.Name, len(expr.Operator.Inputs), len(argsRefs)))
				} else {
					return errors.New(fmt.Sprintf("%s: %d: %s: expected %d inputs; %d were provided",
						expr.FileName, expr.FileLine, OpNames[expr.Operator.OpCode], len(expr.Operator.Inputs), len(argsRefs)))
				}
			}

			// we don't want to modify by reference, we need to make copies
			for i := 0; i < len(argsRefs); i++ {
				if argsRefs[i].Typ == "ident" || len(argsRefs[i].Fields) > 0 {
					var lookingFor string
					// encoder.DeserializeRaw(*argsRefs[i].Value, &lookingFor)

					lookingFor = getFQDN(argsRefs[i])

					// if len(argsRefs[i].Fields) > 0 {
					// 	for _, fld := range argsRefs[i].Fields {
					// 		fmt.Println("this one has it", fld.Name, fld.Value)
					// 	}
					// }

					if arg, err := resolveIdent(lookingFor, argsRefs[i], call); err == nil {
						argsCopy[i] = arg
						// Extracting array values
						// if len(argsRefs[i].Indexes) > 0 {
						// 	if val, err := getValueFromArray(arg, getArrayIndex(argsRefs[i], call)); err == nil {
						// 		arg.Value = &val
						// 	} else {
						// 		panic(err)
						// 	}
						// }
					} else {
						return errors.New(fmt.Sprintf("%d: %s", expr.FileLine, err.Error()))
					}
				} else {
					argsCopy[i] = argsRefs[i]
				}

				// checking if arguments types match with expressions required types
				if len(expr.Operator.Inputs) > 0 &&
					expr.Operator.Inputs[i].Typ !=
						argsCopy[i].Typ &&
					expr.Operator.Inputs[i].Typ != "" &&
					expr.Operator.Inputs[i].Type != TYPE_UNDEFINED {
					return errors.New(fmt.Sprintf("%s: %d: %s: input %d is type '%s'; expected type '%s'\n",
						expr.FileName, expr.FileLine, expr.Operator.Name, i+1, argsCopy[i].Typ, expr.Operator.Inputs[i].Typ))
				}
			}

			var opName string
			var isNative bool
			if expr.Operator != nil {
				if expr.Operator.IsNative {
					isNative = true
					opName = OpNames[expr.Operator.OpCode]
				} else {
					opName = expr.Operator.Name
				}
			} else {
				opName = "id" // return the same
			}

			if _, ok := NATIVE_FUNCTIONS[opName]; ok {
				isNative = true
			}

			// check if struct array function
			if isNative {
				// checkNative(opName, expr, call, &argsCopy, &exc, &excError)
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
								fmt.Printf("%s:\t\t%s\n", def.Name, PrintValue(def.Name, def.Value, def.Typ, call.Program))
							} else {
								fmt.Println(def.Name)
								PrintValue(def.Name, def.Value, def.Typ, call.Program)
							}
						}
					}
					fmt.Println()
					fmt.Printf("%s() Inputs:\n", expr.Operator.Name)
					for i, arg := range argsCopy {
						fmt.Printf("%d: %s\n", i, PrintValue("", arg.Value, arg.Typ, call.Program))
					}
					fmt.Println()
					return excError
				}

				call.Line++
				return call.icall(withDebug, nCalls, callCounter)
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
						fmt.Printf("%s:\t\t%s\n", def.Name, PrintValue(def.Name, def.Value, def.Typ, call.Program))
					}
					fmt.Println()
					fmt.Printf("%s() Arguments:\n", expr.Operator.Name)
					for i, arg := range argsCopy {
						fmt.Printf("%d: %s\n", i, PrintValue("", arg.Value, arg.Typ, call.Program))
					}
					fmt.Println()
					return excError
				}

				// call.Line++ // once the subcall finishes, call next line
				call.Program.CallStack[call.Program.CallCounter].Line++
				if argDefs, err := argsToDefs(argsCopy, expr.Operator.Inputs, expr.Operator.Outputs, call.Package, call.Program); err == nil {
					// fmt.Println("not a native function", expr.Operator.Name)
					subcall := MakeCall(expr.Operator, argDefs, call, call.Package, call.Program)

					// call.Program.CallStack = append(call.Program.CallStack, subcall)
					call.Program.CallCounter++
					call.Program.CallStack[call.Program.CallCounter] = subcall

					return subcall.icall(withDebug, nCalls, callCounter)
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

func (prgrm *CXProgram) ToCall () *CXExpression {
	for c := prgrm.CallCounter - 1; c >= 0; c-- {
		if prgrm.CallStack[c].Line + 1 >= len(prgrm.CallStack[c].Operator.Expressions) {
			// then it'll also return from this function call; continue
			continue
		}
		return prgrm.CallStack[c].Operator.Expressions[prgrm.CallStack[c].Line + 1]
		// prgrm.CallStack[c].Operator.Expressions[prgrm.CallStack[prgrm.CallCounter-1].Line + 1]
	}
	// error
	return &CXExpression{Operator: MakeFunction("")}
	// panic("")
}

func (prgrm *CXProgram) RunCompiled(nCalls int) error {
	// prgrm.PrintProgram()
	rand.Seed(time.Now().UTC().UnixNano())

	var untilEnd bool
	if nCalls == 0 {
		untilEnd = true
	}
	
	if mod, err := prgrm.SelectPackage(MAIN_PKG); err == nil {
		// initializing program resources
		// prgrm.Stacks = append(prgrm.Stacks, MakeStack(1024))

		if prgrm.CallStack[0].Operator == nil {
			// then the program is just starting and we need to run the SYS_INIT_FUNC
			if fn, err := mod.SelectFunction(SYS_INIT_FUNC); err == nil {
				// *init function
				mainCall := MakeCall(fn, nil, nil, mod, mod.Program)
				prgrm.CallStack[0] = mainCall
				prgrm.Stacks[0].StackPointer = fn.Size

				var err error

				for !prgrm.Terminated {
					call := &prgrm.CallStack[prgrm.CallCounter]
					err = call.ccall(prgrm)
					if err != nil {
						return err
					}
				}
				// we reset call state
				prgrm.Terminated = false
				prgrm.CallCounter = 0
				prgrm.CallStack[0].Operator = nil
			} else {
				return err
			}
		}

		if fn, err := mod.SelectFunction(MAIN_FUNC); err == nil {
			if len(fn.Expressions) < 1 {
				return nil
			}

			if prgrm.CallStack[0].Operator == nil {
				// main function
				mainCall := MakeCall(fn, nil, nil, mod, mod.Program)
				// initializing program resources
				prgrm.CallStack[0] = mainCall
				// prgrm.Stacks = append(prgrm.Stacks, MakeStack(1024))
				prgrm.Stacks[0].StackPointer = fn.Size

				prgrm.Terminated = false
			}

			var err error

			for !prgrm.Terminated && (untilEnd || nCalls != 0) {
				call := &prgrm.CallStack[prgrm.CallCounter]

				if !untilEnd {
					var inName string
					var toCallName string
					var toCall *CXExpression

					if call.Line >= call.Operator.Length && prgrm.CallCounter == 0 {
						prgrm.Terminated = true
						prgrm.CallStack[0].Operator = nil
						prgrm.CallCounter = 0
						fmt.Println("in:terminated")
						return err
					}

					if call.Line >= call.Operator.Length && prgrm.CallCounter != 0 {
						toCall = prgrm.ToCall()
						// toCall = prgrm.CallStack[prgrm.CallCounter-1].Operator.Expressions[prgrm.CallStack[prgrm.CallCounter-1].Line + 1]
						inName = prgrm.CallStack[prgrm.CallCounter-1].Operator.Name
					} else {
						toCall = call.Operator.Expressions[call.Line]
						inName = call.Operator.Name
					}

					if toCall.Operator == nil {
						// then it's a declaration
						toCallName = "declaration"
					} else if toCall.Operator.IsNative {
						toCallName = OpNames[toCall.Operator.OpCode]
					} else {
						if toCall.Operator.Name != "" {
							toCallName = toCall.Operator.Package.Name + "." + toCall.Operator.Name
						} else {
							// then it's the end of the program got from nested function calls
							prgrm.Terminated = true
							prgrm.CallStack[0].Operator = nil
							prgrm.CallCounter = 0
							fmt.Println("in:terminated")
							return err
						}
					}
					
					fmt.Printf("in:%s, expr#:%d, calling:%s()\n", inName, call.Line + 1, toCallName)
					
					nCalls--
				}
				
				err = call.ccall(prgrm)
				if err != nil {
					return err
				}
			}

			if prgrm.Terminated {
				prgrm.Terminated = false
				prgrm.CallCounter = 0
				prgrm.CallStack[0].Operator = nil
			}

			// debugging memory
			// fmt.Println("prgrm.Stack", prgrm.Stacks[0].Stack)
			// fmt.Println("prgrm.Heap", prgrm.Heap)
			// fmt.Println("prgrm.Data", prgrm.Data)
			return err
		} else {
			return err
		}
	} else {
		return err
	}
}

func (call *CXCall) ccall(prgrm *CXProgram) error {
	// CX is still single-threaded, so only one stack
	if call.Line >= call.Operator.Length {
		/*
		   popping the stack
		*/
		// going back to the previous call
		prgrm.CallCounter--
		if prgrm.CallCounter < 0 {
			// then the program finished
			prgrm.Terminated = true
		} else {
			// copying the outputs to the previous stack frame
			returnAddr := &prgrm.CallStack[prgrm.CallCounter]
			returnOp := returnAddr.Operator
			returnLine := returnAddr.Line
			returnFP := returnAddr.FramePointer
			fp := call.FramePointer

			expr := returnOp.Expressions[returnLine]
			for i, out := range expr.Outputs {
				WriteMemory(
					&prgrm.Stacks[0],
					GetFinalOffset(&prgrm.Stacks[0], returnFP, out, MEM_WRITE),
					out,
					ReadMemory(
						&prgrm.Stacks[0],
						GetFinalOffset(&prgrm.Stacks[0], fp, call.Operator.Outputs[i], MEM_READ),
						call.Operator.Outputs[i]))
			}

			// return the stack pointer to its previous state
			prgrm.Stacks[0].StackPointer = call.FramePointer
			// we'll now execute the next command
			prgrm.CallStack[prgrm.CallCounter].Line++
			// calling the actual command
			prgrm.CallStack[prgrm.CallCounter].ccall(prgrm)
		}
	} else {
		/*
		   continue with call operator's execution
		*/
		fn := call.Operator
		expr := fn.Expressions[call.Line]
		// if it's a native, then we just process the arguments with execNative
		if expr.Operator == nil {
			// then it's a declaration
			call.Line++
		} else if expr.Operator.IsNative {
			execNative(prgrm)
			call.Line++
		} else {
			/*
			   It was not a native, so we need to create another call
			   with the current expression's operator
			*/
			// we're going to use the next call in the callstack
			prgrm.CallCounter++
			newCall := &prgrm.CallStack[prgrm.CallCounter]
			// setting the new call
			newCall.Operator = expr.Operator
			newCall.Line = 0
			newCall.FramePointer = prgrm.Stacks[0].StackPointer
			// the stack pointer is moved to create room for the next call
			// prgrm.Stacks[0].StackPointer += fn.Size
			prgrm.Stacks[0].StackPointer += newCall.Operator.Size

			fp := call.FramePointer
			newFP := newCall.FramePointer

			// wiping next stack frame (removing garbage)
			for c := 0; c < expr.Operator.Size; c++ {
				prgrm.Stacks[0].Stack[newFP+c] = 0
			}

			for i, inp := range expr.Inputs {
				var byts []byte
				// finalOffset := inp.Offset
				finalOffset := GetFinalOffset(&prgrm.Stacks[0], fp, inp, MEM_READ)
				// finalOffset := fp + inp.Offset

				// if inp.Indexes != nil {
				// 	finalOffset = GetFinalOffset(&prgrm.Stacks[0], fp, inp)
				// }
				if inp.IsReference {
					byts = encoder.Serialize(int32(finalOffset))
				} else {
					switch inp.MemoryWrite {
					case MEM_STACK:
						byts = prgrm.Stacks[0].Stack[finalOffset : finalOffset+inp.TotalSize]
					case MEM_DATA:
						byts = prgrm.Data[finalOffset : finalOffset+inp.TotalSize]
					case MEM_HEAP:
						byts = prgrm.Heap.Heap[finalOffset : finalOffset+inp.TotalSize]
					default:
						panic("implement the other mem types")
					}
				}
				
				// writing inputs to new stack frame
				WriteToStack(
					&prgrm.Stacks[0],
					GetFinalOffset(&prgrm.Stacks[0], newFP, newCall.Operator.Inputs[i], MEM_WRITE),
					byts)
			}
		}
	}
	return nil
}
