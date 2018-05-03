package main

import (
	"fmt"
	"os"
	"time"
	"strings"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	. "github.com/skycoin/cx/src/interpreted"
)



func warnf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Stderr.Sync()
}

func binaryOp (op string, arg1, arg2 *CXArgument, line int) *CXArgument {
	var opName string
	var typArg1 string
	// var typArg2 string
	// _ = typArg2

	if (len(arg1.Typ) > len("ident.") && arg1.Typ[:len("ident.")] == "ident.") {
		arg1.Typ = "ident"
	}

	if (len(arg2.Typ) > len("ident.") && arg2.Typ[:len("ident.")] == "ident.") {
		arg2.Typ = "ident"
	}
	
	if arg1.Typ == "ident" {
		var identName string
		encoder.DeserializeRaw(*arg1.Value, &identName)

		if typ, err := GetIdentType(identName, line, fileName, cxt); err == nil {
			typArg1 = typ
		} else {
			fmt.Println(err)
		}
	} else {
		typArg1 = arg1.Typ
	}

	switch op {
	case "+":
		opName = fmt.Sprintf("%s.add", typArg1)
	case "-":
		opName = fmt.Sprintf("%s.sub", typArg1)
	case "*":
		opName = fmt.Sprintf("%s.mul", typArg1)
	case "/":
		opName = fmt.Sprintf("%s.div", typArg1)
	case "%":
		opName = fmt.Sprintf("%s.mod", typArg1)
	case ">":
		opName = fmt.Sprintf("%s.gt", typArg1)
	case "<":
		opName = fmt.Sprintf("%s.lt", typArg1)
	case "<=":
		opName = fmt.Sprintf("%s.lteq", typArg1)
	case ">=":
		opName = fmt.Sprintf("%s.gteq", typArg1)
	case "<<":
		opName = fmt.Sprintf("%s.bitshl", typArg1)
	case ">>":
		opName = fmt.Sprintf("%s.bitshr", typArg1)
	case "**":
		opName = fmt.Sprintf("%s.pow", typArg1)
	case "&":
		opName = fmt.Sprintf("%s.bitand", typArg1)
	case "|":
		opName = fmt.Sprintf("%s.bitor", typArg1)
	case "^":
		opName = fmt.Sprintf("%s.bitxor", typArg1)
	case "&^":
		opName = fmt.Sprintf("%s.bitclear", typArg1)
	case "&&":
		opName = "and"
	case "||":
		opName = "or"
	case "==":
		opName = fmt.Sprintf("%s.eq", typArg1)
	case "!=":
		opName = fmt.Sprintf("%s.uneq", typArg1)
	}

	if fn, err := cxt.GetCurrentFunction(); err == nil {
		if op, err := cxt.GetFunction(opName, CORE_MODULE); err == nil {
			expr := MakeExpression(op)
			if !replMode {
				expr.FileLine = line
				expr.FileName = fileName
			}
			fn.AddExpression(expr)
			expr.AddLabel(tag)
			tag = ""
			expr.AddInput(arg1)
			expr.AddInput(arg2)

			outName := MakeGenSym(NON_ASSIGN_PREFIX)
			byteName := encoder.Serialize(outName)
			
			expr.AddOutputName(outName)
			return MakeArgument(&byteName, "ident")
		}
	}
	return nil
}

func unaryOp (op string, arg1 *CXArgument, line int) *CXArgument {
	var opName string
	var typArg1 string

	if arg1.Typ == "ident" {
		var identName string
		encoder.DeserializeRaw(*arg1.Value, &identName)

		if typ, err := GetIdentType(identName, line, fileName, cxt); err == nil {
			typArg1 = typ
		} else {
			fmt.Println(err)
		}
	} else {
		typArg1 = arg1.Typ
	}

	switch op {
	case "++":
		opName = fmt.Sprintf("%s.add", typArg1)
	case "--":
		opName = fmt.Sprintf("%s.sub", typArg1)
	}
	
	if fn, err := cxt.GetCurrentFunction(); err == nil {
		if op, err := cxt.GetFunction(opName, CORE_MODULE); err == nil {
			expr := MakeExpression(op)
			if !replMode {
				expr.FileLine = line
				expr.FileName = fileName
			}
			fn.AddExpression(expr)
			expr.AddLabel(tag)
			tag = ""

			
			expr.AddInput(arg1)

			// var one *CXArgument

			switch typArg1 {
			case "i32":
				sOne := encoder.Serialize(int32(1))
				expr.AddInput(MakeArgument(&sOne, "i32"))
			case "i64":
				sOne := encoder.Serialize(int64(1))
				expr.AddInput(MakeArgument(&sOne, "i64"))
			case "f32":
				sOne := encoder.Serialize(float32(1))
				expr.AddInput(MakeArgument(&sOne, "f32"))
			case "f64":
				sOne := encoder.Serialize(float64(1))
				expr.AddInput(MakeArgument(&sOne, "f64"))
			}

			var outName string
			if arg1.Typ == "ident" {
				encoder.DeserializeRaw(*arg1.Value, &outName)
			} else {
				outName = MakeGenSym(NON_ASSIGN_PREFIX)
			}
			
			byteName := encoder.Serialize(outName)
			
			expr.AddOutputName(outName)
			return MakeArgument(&byteName, "ident")
		}
	}
	return nil
}

func Import (name string) {
	impName := strings.TrimPrefix(name, "\"")
	impName = strings.TrimSuffix(impName, "\"")
	if imp, err := cxt.GetModule(impName); err == nil {
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			mod.AddImport(imp)
		}
	}
}

const (
	// affordance element
	AFF_FUNC = iota
	AFF_PKG
	AFF_STRCT
	AFF_EXPR
)

const (
	// affordance type
	AFF_TYP1 = iota
	AFF_TYP2
	AFF_TYP3
	AFF_TYP4
)

func Affordance (affElt int, affTyp int, ident string, lbl string, idx int32) {
	switch affElt {
	case AFF_FUNC:
		switch affTyp {
		case AFF_TYP1:
			if mod, err := cxt.GetCurrentPackage(); err == nil {
				if fn, err := cxt.GetFunction(ident, mod.Name); err == nil {
					affs := fn.GetAffordances()
					for i, aff := range affs {
						fmt.Printf("(%d)\t%s\n", i, aff.Description)
					}
				}
			}
		case AFF_TYP2:
			if mod, err := cxt.GetCurrentPackage(); err == nil {
				if fn, err := cxt.GetFunction(ident, mod.Name); err == nil {
					affs := fn.GetAffordances()
					affs[idx].ApplyAffordance()
				}
			}
		case AFF_TYP3:
			if mod, err := cxt.GetCurrentPackage(); err == nil {
				if fn, err := cxt.GetFunction(ident, mod.Name); err == nil {
					affs := fn.GetAffordances()
					filter := strings.TrimPrefix(lbl, "\"")
					filter = strings.TrimSuffix(filter, "\"")
					affs = FilterAffordances(affs, filter)
					for i, aff := range affs {
						fmt.Printf("(%d)\t%s\n", i, aff.Description)
					}
				}
			}
		case AFF_TYP4:
			if mod, err := cxt.GetCurrentPackage(); err == nil {
				if fn, err := cxt.GetFunction(ident, mod.Name); err == nil {
					affs := fn.GetAffordances()
					filter := strings.TrimPrefix(lbl, "\"")
					filter = strings.TrimSuffix(filter, "\"")
					affs = FilterAffordances(affs, filter)
					affs[idx].ApplyAffordance()
				}
			}
		}
	case AFF_PKG:
		switch affTyp {
		case AFF_TYP1:
			if mod, err := cxt.GetModule(ident); err == nil {
				affs := mod.GetAffordances()
				for i, aff := range affs {
					fmt.Printf("(%d)\t%s\n", i, aff.Description)
				}
			}
		case AFF_TYP2:
			if mod, err := cxt.GetModule(ident); err == nil {
				affs := mod.GetAffordances()
				affs[idx].ApplyAffordance()
			}
		case AFF_TYP3:
			if mod, err := cxt.GetModule(ident); err == nil {
				affs := mod.GetAffordances()
				filter := strings.TrimPrefix(lbl, "\"")
				filter = strings.TrimSuffix(filter, "\"")
				affs = FilterAffordances(affs, filter)
				for i, aff := range affs {
					fmt.Printf("(%d)\t%s\n", i, aff.Description)
				}
			}
		case AFF_TYP4:
			if mod, err := cxt.GetModule(ident); err == nil {
				affs := mod.GetAffordances()
				filter := strings.TrimPrefix(lbl, "\"")
				filter = strings.TrimSuffix(filter, "\"")
				affs = FilterAffordances(affs, filter)
				affs[idx].ApplyAffordance()
			}
		}
	case AFF_STRCT:
		switch affTyp {
		case AFF_TYP1:
			if mod, err := cxt.GetCurrentPackage(); err == nil {
				if strct, err := cxt.GetStruct(ident, mod.Name); err == nil {
					affs := strct.GetAffordances()
					for i, aff := range affs {
						fmt.Printf("(%d)\t%s\n", i, aff.Description)
					}
				}
			}
		case AFF_TYP2:
			if mod, err := cxt.GetCurrentPackage(); err == nil {
				if strct, err := cxt.GetStruct(ident, mod.Name); err == nil {
					affs := strct.GetAffordances()
					affs[idx].ApplyAffordance()
				}
			}
		case AFF_TYP3:
			if mod, err := cxt.GetCurrentPackage(); err == nil {
				if strct, err := cxt.GetStruct(ident, mod.Name); err == nil {
					affs := strct.GetAffordances()
					filter := strings.TrimPrefix(lbl, "\"")
					filter = strings.TrimSuffix(filter, "\"")
					affs = FilterAffordances(affs, filter)
					for i, aff := range affs {
						fmt.Printf("(%d)\t%s\n", i, aff.Description)
					}
				}
			}
		case AFF_TYP4:
			if mod, err := cxt.GetCurrentPackage(); err == nil {
				if strct, err := cxt.GetStruct(ident, mod.Name); err == nil {
					affs := strct.GetAffordances()
					filter := strings.TrimPrefix(lbl, "\"")
					filter = strings.TrimSuffix(filter, "\"")
					affs = FilterAffordances(affs, filter)
					affs[idx].ApplyAffordance()
				}
			}
		}
	case AFF_EXPR:
		switch affTyp {
		case AFF_TYP1:
			if mod, err := cxt.GetCurrentPackage(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					for _, expr := range fn.Expressions {
						if expr.Label == ident {
							PrintAffordances(expr.GetAffordances(nil))
							break
						}
					}
				}
			}
		case AFF_TYP2:
			if mod, err := cxt.GetCurrentPackage(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					for _, expr := range fn.Expressions {
						if expr.Label == ident {
							affs := expr.GetAffordances(nil)
							affs[idx].ApplyAffordance()
							break
						}
					}
				}
			}
		case AFF_TYP3:
			if mod, err := cxt.GetCurrentPackage(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					for _, expr := range fn.Expressions {
						if expr.Label == ident {
							affs := expr.GetAffordances(nil)
							filter := strings.TrimPrefix(lbl, "\"")
							filter = strings.TrimSuffix(filter, "\"")
							PrintAffordances(FilterAffordances(affs, filter))
							break
						}
					}
				}
			}
		case AFF_TYP4:
			if mod, err := cxt.GetCurrentPackage(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					for _, expr := range fn.Expressions {
						if expr.Label == ident {
							affs := expr.GetAffordances(nil)
							filter := strings.TrimPrefix(lbl, "\"")
							filter = strings.TrimSuffix(filter, "\"")
							affs = FilterAffordances(affs, filter)
							affs[idx].ApplyAffordance()
							break
						}
					}
				}
			}
		}
	}
}

func Stepping (steps int, delay int, withDelay bool) {
	if withDelay {
		if steps == 0 {
			// Maybe nothing for now
		} else {
			if steps < 0 {
				nCalls := steps * -1
				for i := 0; i < nCalls; i++ {
					time.Sleep(time.Duration(int32(delay)) * time.Millisecond)
					cxt.UnRun(1)
				}
			} else {

				for i := 0; i < steps; i++ {
					time.Sleep(time.Duration(int32(delay)) * time.Millisecond)
					err := cxt.Run(dStack, 1)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	} else {
		if steps == 0 {
			// we run until halt or end of program;
			if err := cxt.Run(dStack, -1); err != nil {
				fmt.Println(err)
			}
		} else {
			if steps < 0 {
				nCalls := steps * -1
				cxt.UnRun(int(nCalls))
			} else {
				//fmt.Println(cxt.Run(dStack, int(steps)))

				err := cxt.Run(dStack, int(steps))
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func DebugState () {
	if len(cxt.CallStack) > 0 {
		if len(cxt.CallStack[len(cxt.CallStack) - 1].State) > 0 {
			for _, def := range cxt.CallStack[len(cxt.CallStack) - 1].State {
				var isNonAssign bool
				if len(def.Name) > len(NON_ASSIGN_PREFIX) && def.Name[:len(NON_ASSIGN_PREFIX)] == NON_ASSIGN_PREFIX {
					isNonAssign = true
				}

				if !isNonAssign {
					if IsBasicType(def.Typ) {
						fmt.Printf("%s:\t\t%s\n", def.Name, PrintValue(def.Name, def.Value, def.Typ, cxt))
					} else {
						fmt.Println(def.Name)
						PrintValue(def.Name, def.Value, def.Typ, cxt)
					}
				}
			}
		}
	}
}

func DebugStack () {
	if dStack {
		dStack = false
		fmt.Println("* printing stack: false")
	} else {
		dStack = true
		fmt.Println("* printing stack: true")
	}
}

const (
	REM_TYP_FUNC = iota
	REM_TYP_PKG
	REM_TYP_GLBL
	REM_TYP_STRCT
	REM_TYP_IMP
	REM_TYP_EXPR
	REM_TYP_FLD
	REM_TYP_INPUT
	REM_TYP_OUTPUT
)

func Remover (remTyp int, fstIdent string, sndIdent string) {
	switch remTyp {
	case REM_TYP_FUNC:
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			mod.RemoveFunction(fstIdent)
		}
	case REM_TYP_PKG:
		cxt.RemoveModule(fstIdent)
	case REM_TYP_GLBL:
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			mod.RemoveGlobal(fstIdent)
		}
	case REM_TYP_STRCT:
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			mod.RemoveStruct(fstIdent)
		}
	case REM_TYP_IMP:
		impName := strings.TrimPrefix(fstIdent, "\"")
		impName = strings.TrimSuffix(impName, "\"")
		
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			mod.RemoveImport(impName)
		}
	case REM_TYP_EXPR:
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			if fn, err := mod.Program.GetFunction(sndIdent, mod.Name); err == nil {
				for i, expr := range fn.Expressions {
					if expr.Label == fstIdent {
						fn.RemoveExpression(i)
					}
				}
			}
		}
	case REM_TYP_FLD:
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			if strct, err := cxt.GetStruct(sndIdent, mod.Name); err == nil {
				strct.RemoveField(fstIdent)
			}
			
		}
	case REM_TYP_INPUT:
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			if fn, err := mod.Program.GetFunction(sndIdent, mod.Name); err == nil {
				fn.RemoveInput(fstIdent)
			}
		}
	case REM_TYP_OUTPUT:
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			if fn, err := mod.Program.GetFunction(sndIdent, mod.Name); err == nil {
				fn.RemoveOutput(fstIdent)
			}
		}
	}
}

const (
	// type of selector
	SELECT_TYP_PKG = iota
	SELECT_TYP_FUNC
	SELECT_TYP_STRCT
)

func SelectorFields (flds []*CXArgument) bool {
	if strct, err := cxt.GetCurrentStruct(); err == nil {
		for _, fld := range flds {
			fldFromParam := MakeField(fld.Name, fld.Typ)
			strct.AddField(fldFromParam)
		}
	}
	return true
}

func Selector (ident string, selTyp int) string {
	switch selTyp {
	case SELECT_TYP_PKG:
		var previousModule *CXPackage
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			previousModule = mod
		} else {
			fmt.Println("A current module does not exist")
		}
		if _, err := cxt.SelectModule(ident); err == nil {
			//fmt.Println(fmt.Sprintf("== Changed to package '%s' ==", mod.Name))
		} else {
			fmt.Println(err)
		}

		replTargetMod = ident
		replTargetStrct = ""
		replTargetFn = ""
		
		return previousModule.Name
	case SELECT_TYP_FUNC:
		var previousFunction *CXFunction
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			previousFunction = fn
		} else {
			fmt.Println("A current function does not exist")
		}
		if _, err := cxt.SelectFunction(ident); err == nil {
			//fmt.Println(fmt.Sprintf("== Changed to function '%s' ==", fn.Name))
		} else {
			fmt.Println(err)
		}

		replTargetMod = ""
		replTargetStrct = ""
		replTargetFn = ident
		
		return previousFunction.Name
	case SELECT_TYP_STRCT:
		var previousStruct *CXStruct
		if fn, err := cxt.GetCurrentStruct(); err == nil {
			previousStruct = fn
		} else {
			fmt.Println("A current struct does not exist")
		}
		if _, err := cxt.SelectStruct(ident); err == nil {
			//fmt.Println(fmt.Sprintf("== Changed to struct '%s' ==", fn.Name))
		} else {
			fmt.Println(err)
		}

		replTargetStrct = ident
		replTargetMod = ""
		replTargetFn = ""
		
		return previousStruct.Name
	}

	panic("")
	
}

func GlobalDeclaration (isBasic bool, ident string, typ string, assignment *CXArgument, line int) {
	if isBasic {
		if assignment != nil {
			if typ != assignment.Typ {
				panic(fmt.Sprintf("%s: %d: variable of type '%s' cannot be initialized with value of type '%s'", fileName, line, typ, assignment.Typ))
			}
		}

		if mod, err := cxt.GetCurrentPackage(); err == nil {
			var val *CXArgument;
			if assignment == nil {
				val = MakeArgument(MakeDefaultValue(typ), typ)
			} else {
				switch typ {
				case "byte":
					// var ds int32
					// encoder.DeserializeRaw(*assignment.Value, &ds)

					//fmt.Println("here", assignment.Value)
					
					//new := []byte{byte(ds)}
					//val = MakeArgument(&new, "byte")
					val = MakeArgument(assignment.Value, "byte")
				case "i64":
					// var ds int32
					// encoder.DeserializeRaw(*assignment.Value, &ds)
					// new := encoder.Serialize(int64(ds))
					// val = MakeArgument(&new, "i64")

					val = MakeArgument(assignment.Value, "i64")
				case "f64":
					// var ds float32
					// encoder.DeserializeRaw(*assignment.Value, &ds)
					// new := encoder.Serialize(float64(ds))
					// val = MakeArgument(&new, "f64")

					val = MakeArgument(assignment.Value, "f64")
				default:
					val = assignment
				}
			}

			mod.AddDefinition(MakeDefinition(ident, val.Value, typ))
		}
	} else {
		// we have to initialize all the fields
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			if zeroVal, err := ResolveStruct(typ, cxt); err == nil {
				mod.AddDefinition(MakeDefinition(ident, &zeroVal, typ))
			} else {
				fmt.Println(fmt.Sprintf("%s: %d: definition declaration: %s", fileName, line, err))
			}
		}
	}
}

func StructDeclaration (ident string, line int) {
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		strct := MakeStruct(ident)
		mod.AddStruct(strct)


		// creating manipulation functions for this type a la common lisp
		// append
		fn := MakeFunction(fmt.Sprintf("[]%s.append", ident))
		fn.AddInput(MakeParameter("arr", fmt.Sprintf("[]%s", ident)))
		fn.AddInput(MakeParameter("strctInst", ident))
		fn.AddOutput(MakeParameter("_arr", fmt.Sprintf("[]%s", ident)))
		mod.AddFunction(fn)

		if op, err := cxt.GetFunction("cstm.append", CORE_MODULE); err == nil {
			expr := MakeExpression(op)
			if !replMode {
				expr.FileLine = line
				expr.FileName = fileName
			}
			sArr := encoder.Serialize("arr")
			arrArg := MakeArgument(&sArr, "str")
			sStrctInst := encoder.Serialize("strctInst")
			strctInstArg := MakeArgument(&sStrctInst, "str")
			expr.AddInput(arrArg)
			expr.AddInput(strctInstArg)
			expr.AddOutputName("_arr")
			
			
			fn.AddExpression(expr)
		} else {
			fmt.Println(err)
		}

		// serialize
		fn = MakeFunction(fmt.Sprintf("%s.serialize", ident))
		fn.AddInput(MakeParameter("strctInst", ident))
		fn.AddOutput(MakeParameter("byts", "[]byte"))
		mod.AddFunction(fn)

		if op, err := cxt.GetFunction("cstm.serialize", CORE_MODULE); err == nil {
			expr := MakeExpression(op)
			if !replMode {
				expr.FileLine = line
				expr.FileName = fileName
			}
	
			sStrctInst := encoder.Serialize("strctInst")
			strctInstArg := MakeArgument(&sStrctInst, "str")
			expr.AddInput(strctInstArg)
			expr.AddOutputName("byts")
			fn.AddExpression(expr)
		} else {
			fmt.Println(err)
		}



		// deserialize
		fn = MakeFunction(fmt.Sprintf("%s.deserialize", ident))
		fn.AddInput(MakeParameter("byts", "[]byte"))
		fn.AddOutput(MakeParameter("strctInst", ident))
		mod.AddFunction(fn)

		if op, err := cxt.GetFunction("cstm.deserialize", CORE_MODULE); err == nil {
			expr := MakeExpression(op)
			if !replMode {
				expr.FileLine = line
				expr.FileName = fileName
			}

			sByts := encoder.Serialize("byts")
			sBytsArg := MakeArgument(&sByts, "str")

			sTyp := encoder.Serialize(ident)
			sTypArg := MakeArgument(&sTyp, "str")
			
			expr.AddInput(sBytsArg)
			expr.AddInput(sTypArg)
			expr.AddOutputName("strctInst")
			
			fn.AddExpression(expr)
		} else {
			fmt.Println(err)
		}

		
		// read
		fn = MakeFunction(fmt.Sprintf("[]%s.read", ident))
		fn.AddInput(MakeParameter("arr", fmt.Sprintf("[]%s", ident)))
		fn.AddInput(MakeParameter("index", "i32"))
		fn.AddOutput(MakeParameter("strctInst", ident))
		mod.AddFunction(fn)

		if op, err := cxt.GetFunction("cstm.read", CORE_MODULE); err == nil {
			expr := MakeExpression(op)
			if !replMode {
				expr.FileLine = line
				expr.FileName = fileName
			}
			sArr := encoder.Serialize("arr")
			arrArg := MakeArgument(&sArr, "str")
			sIndex := encoder.Serialize("index")
			indexArg := MakeArgument(&sIndex, "ident")
			expr.AddInput(arrArg)
			expr.AddInput(indexArg)
			expr.AddOutputName("strctInst")
			fn.AddExpression(expr)
		} else {
			fmt.Println(err)
		}
		// write
		fn = MakeFunction(fmt.Sprintf("[]%s.write", ident))
		fn.AddInput(MakeParameter("arr", fmt.Sprintf("[]%s", ident)))
		fn.AddInput(MakeParameter("index", "i32"))
		fn.AddInput(MakeParameter("inst", ident))
		fn.AddOutput(MakeParameter("_arr", fmt.Sprintf("[]%s", ident)))
		mod.AddFunction(fn)

		if op, err := cxt.GetFunction("cstm.write", CORE_MODULE); err == nil {
			expr := MakeExpression(op)
			if !replMode {
				expr.FileLine = line
				expr.FileName = fileName
			}
			sArr := encoder.Serialize("arr")
			arrArg := MakeArgument(&sArr, "str")
			sIndex := encoder.Serialize("index")
			indexArg := MakeArgument(&sIndex, "ident")
			sInst := encoder.Serialize("inst")
			instArg := MakeArgument(&sInst, "str")
			expr.AddInput(arrArg)
			expr.AddInput(indexArg)
			expr.AddInput(instArg)
			expr.AddOutputName("_arr")
			fn.AddExpression(expr)
		} else {
			fmt.Println(err)
		}
		// len
		fn = MakeFunction(fmt.Sprintf("[]%s.len", ident))
		fn.AddInput(MakeParameter("arr", fmt.Sprintf("[]%s", ident)))
		fn.AddOutput(MakeParameter("len", "i32"))
		mod.AddFunction(fn)

		if op, err := cxt.GetFunction("cstm.len", CORE_MODULE); err == nil {
			expr := MakeExpression(op)
			if !replMode {
				expr.FileLine = line
				expr.FileName = fileName
			}
			sArr := encoder.Serialize("arr")
			arrArg := MakeArgument(&sArr, "str")
			expr.AddInput(arrArg)
			expr.AddOutputName("len")
			fn.AddExpression(expr)
		} else {
			fmt.Println(err)
		}
		
		// make
		fn = MakeFunction(fmt.Sprintf("[]%s.make", ident))
		fn.AddInput(MakeParameter("len", "i32"))
		fn.AddOutput(MakeParameter("arr", fmt.Sprintf("[]%s", ident)))
		mod.AddFunction(fn)

		if op, err := cxt.GetFunction("cstm.make", CORE_MODULE); err == nil {
			expr := MakeExpression(op)
			if !replMode {
				expr.FileLine = line
				expr.FileName = fileName
			}
			sLen := encoder.Serialize("len")
			sTyp := encoder.Serialize(fmt.Sprintf("[]%s", ident))
			lenArg := MakeArgument(&sLen, "ident")
			typArg := MakeArgument(&sTyp, "str")
			expr.AddInput(lenArg)
			expr.AddInput(typArg)
			expr.AddOutputName("arr")
			fn.AddExpression(expr)
		} else {
			fmt.Println(err)
		}
	}
}

func StructDeclarationFields (flds []*CXArgument) {
	if strct, err := cxt.GetCurrentStruct(); err == nil {
		for _, fld := range flds {
			fldFromParam := MakeField(fld.Name, fld.Typ)
			strct.AddField(fldFromParam)
		}
	}
}

const (
	METHOD_INP = iota
	METHOD_INP_OUT
	FUNC_INP
	FUNC_INP_OUT
)

func FunctionDeclarationHeader (typFunc int, ident string, receiver []*CXArgument, inputs []*CXArgument, outputs []*CXArgument, line int) {
	switch typFunc {
	case METHOD_INP_OUT:
		if len(receiver) > 1 {
			panic(fmt.Sprintf("%s: %d: method '%s' has multiple receivers", fileName, line, ident))
		}

		if mod, err := cxt.GetCurrentPackage(); err == nil {
			if IsBasicType(receiver[0].Typ) {
				panic(fmt.Sprintf("%s: %d: cannot define methods on basic type %s", fileName, line, receiver[0].Typ))
			}
			
			inFn = true
			fn := MakeFunction(fmt.Sprintf("%s.%s", receiver[0].Typ, ident))
			mod.AddFunction(fn)
			if fn, err := mod.GetCurrentFunction(); err == nil {

				//checking if there are duplicate parameters
				dups := append(inputs, outputs...)
				dups = append(dups, receiver...)
				for _, param := range dups {
					for _, dup := range dups {
						if param.Name == dup.Name && param != dup {
							panic(fmt.Sprintf("%s: %d: duplicate receiver, input and/or output parameters in method '%s'", fileName, line, ident))
						}
					}
				}

				for _, rec := range receiver {
					fn.AddInput(rec)
				}
				for _, inp := range inputs {
					fn.AddInput(inp)
				}
				for _, out := range outputs {
					fn.AddOutput(out)
				}
			}
		}
	case METHOD_INP:
		if len(receiver) > 1 {
			panic(fmt.Sprintf("%s: %d: method '%s' has multiple receivers", fileName, line, ident))
		}
		
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			if IsBasicType(receiver[0].Typ) {
				panic(fmt.Sprintf("%s: %d: cannot define methods on basic type %s", fileName, line, receiver[0].Typ))
			}
			
			inFn = true
			fn := MakeFunction(fmt.Sprintf("%s.%s", receiver[0].Typ, ident))
			mod.AddFunction(fn)
			if fn, err := mod.GetCurrentFunction(); err == nil {

				//checking if there are duplicate parameters
				dups := append(receiver, inputs...)
				for _, param := range dups {
					for _, dup := range dups {
						if param.Name == dup.Name && param != dup {
							panic(fmt.Sprintf("%s: %d: duplicate receiver, input and/or output parameters in method '%s'", fileName, line, ident))
						}
					}
				}

				for _, rec := range receiver {
					fn.AddInput(rec)
				}
				for _, inp := range inputs {
					fn.AddInput(inp)
				}
			}
		}
	case FUNC_INP_OUT:
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			inFn = true
			fn := MakeFunction(ident)
			mod.AddFunction(fn)
			if fn, err := mod.GetCurrentFunction(); err == nil {

				//checking if there are duplicate parameters
				dups := append(inputs, outputs...)
				for _, param := range dups {
					for _, dup := range dups {
						if param.Name == dup.Name && param != dup {
							panic(fmt.Sprintf("%s: %d: duplicate input and/or output parameters in function '%s'", fileName, line, ident))
						}
					}
				}
				
				for _, inp := range inputs {
					fn.AddInput(inp)
				}
				for _, out := range outputs {
					fn.AddOutput(out)
				}
			}
		}
	case FUNC_INP:
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			inFn = true
			fn := MakeFunction(ident)
			mod.AddFunction(fn)
			if fn, err := mod.GetCurrentFunction(); err == nil {
				for _, inp := range inputs {
					fn.AddInput(inp)
				}
			}
		}
	}
}

func AssignBasicVar (ident string, typ string, initializer *CXArgument, line int) {
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			if initializer == nil {
				if op, err := cxt.GetFunction("initDef", mod.Name); err == nil {
					expr := MakeExpression(op)
					if !replMode {
						expr.FileLine = line
						expr.FileName = fileName
					}

					fn.AddExpression(expr)
					
					typ := encoder.Serialize(typ)
					arg := MakeArgument(&typ, "str")
					expr.AddInput(arg)
					expr.AddOutputName(ident)

					// if strct, err := cxt.GetStruct(typ, mod.Name); err == nil {
					// 	for _, fld := range strct.Fields {
					// 		expr := MakeExpression(op)
					// 		if !replMode {
					// 			expr.FileLine = line
					// 			expr.FileName = fileName
					// 		}
					// 		fn.AddExpression(expr)
					// 		typ := []byte(fld.Typ)
					// 		arg := MakeArgument(&typ, "str")
					// 		expr.AddInput(arg)
					// 		expr.AddOutputName(fmt.Sprintf("%s.%s", ident, fld.Name))
					// 	}
					// }
				}
			} else {
				switch typ {
				case "bool":
					var ds int32
					encoder.DeserializeRaw(*initializer.Value, &ds)
					new := encoder.SerializeAtomic(ds)
					val := MakeArgument(&new, "bool")
					
					if op, err := cxt.GetFunction("bool.id", mod.Name); err == nil {
						expr := MakeExpression(op)
						if !replMode {
							expr.FileLine = line
							expr.FileName = fileName
						}
						fn.AddExpression(expr)
						expr.AddInput(val)
						expr.AddOutputName(ident)
					}
				case "byte":
					var ds int32
					encoder.DeserializeRaw(*initializer.Value, &ds)
					new := []byte{byte(ds)}
					val := MakeArgument(&new, "byte")
					
					if op, err := cxt.GetFunction("byte.id", mod.Name); err == nil {
						expr := MakeExpression(op)
						if !replMode {
							expr.FileLine = line
							expr.FileName = fileName
						}
						fn.AddExpression(expr)
						expr.AddInput(val)
						expr.AddOutputName(ident)
					}
				case "i64":
					var ds int32
					encoder.DeserializeRaw(*initializer.Value, &ds)
					new := encoder.Serialize(int64(ds))
					val := MakeArgument(&new, "i64")

					if op, err := cxt.GetFunction("i64.id", mod.Name); err == nil {
						expr := MakeExpression(op)
						if !replMode {
							expr.FileLine = line
							expr.FileName = fileName
						}
						fn.AddExpression(expr)
						expr.AddInput(val)
						expr.AddOutputName(ident)
					}
				case "f64":
					var ds float32
					encoder.DeserializeRaw(*initializer.Value, &ds)
					new := encoder.Serialize(float64(ds))
					val := MakeArgument(&new, "f64")

					if op, err := cxt.GetFunction("f64.id", mod.Name); err == nil {
						expr := MakeExpression(op)
						if !replMode {
							expr.FileLine = line
							expr.FileName = fileName
						}
						fn.AddExpression(expr)
						expr.AddInput(val)
						expr.AddOutputName(ident)
					}
				default:
					val := initializer
					var getFn string
					switch typ {
					case "i32": getFn = "i32.id"
					case "f32": getFn = "f32.id"
					case "[]bool": getFn = "[]bool.id"
					case "[]byte": getFn = "[]byte.id"
					case "[]str": getFn = "[]str.id"
					case "[]i32": getFn = "[]i32.id"
					case "[]i64": getFn = "[]i64.id"
					case "[]f32": getFn = "[]f32.id"
					case "[]f64": getFn = "[]f64.id"
					}

					if op, err := cxt.GetFunction(getFn, mod.Name); err == nil {
						expr := MakeExpression(op)
						if !replMode {
							expr.FileLine = line
							expr.FileName = fileName
						}
						fn.AddExpression(expr)
						expr.AddInput(val)
						expr.AddOutputName(ident)
					}
				}
			}
		}
	}
}

func AssignCustomVar (ident string, typ string, line int) {
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			if op, err := cxt.GetFunction("initDef", mod.Name); err == nil {
				expr := MakeExpression(op)

				if !replMode {
					expr.FileLine = line
					expr.FileName = fileName
				}
				fn.AddExpression(expr)
				typ := encoder.Serialize(fmt.Sprintf("[]%s", typ))
				arg := MakeArgument(&typ, "str")
				expr.AddInput(arg)
				expr.AddOutputName(ident)
			}
		}
	}
}

func AssignExpression (to []*CXArgument, op string, from []*CXArgument, line int) {
	argsL := to
	argsR := from

	if len(argsL) > len(argsR) {
		panic(fmt.Sprintf("%s: %d: trying to assign values to variables using a function with no output parameters", fileName, line))
	}

	if fn, err := cxt.GetCurrentFunction(); err == nil {
		for i, argL := range argsL {
			if argsR[i] == nil {
				continue
			}
			// argL is going to be the output name
			typeParts := strings.Split(argsR[i].Typ, ".")

			var typ string
			var secondTyp string
			var idFn string
			var ptrs string

			for i, char := range typeParts[0] {
				if char != '*' {
					typeParts[0] = typeParts[0][i:]
					break
				} else {
					ptrs += "*"
				}
			}

			if len(typeParts) > 1 {
				typ = "str"
				secondTyp = strings.Join(typeParts[1:], ".")
			} else if typeParts[0] == "ident" {
				typ = "str"
				secondTyp = "ident"
			} else {
				typ = typeParts[0] // i32, f32, etc
			}

			if op == ":=" || op == "=" {
				if secondTyp == "" {
					idFn = MakeIdentityOpName(typ)
				} else {
					idFn = "identity"
				}

				if op, err := cxt.GetFunction(idFn, CORE_MODULE); err == nil {
					expr := MakeExpression(op)
					if !replMode {
						expr.FileLine = line
						expr.FileName = fileName
					}

					fn.AddExpression(expr)
					expr.AddLabel(tag)
					tag = ""

					var outName string
					encoder.DeserializeRaw(*argL.Value, &outName)

					// // checking if identifier was previously declared
					// if outType, err := GetIdentType(outName, line, fileName, cxt); err == nil {
					// 	if len(typeParts) > 1 {
					// 		if outType != secondTyp {
					// 			panic(fmt.Sprintf("%s: %d: identifier '%s' was previously declared as '%s'; cannot use type '%s' in assignment", fileName, line, outName, outType, secondTyp))
					// 		}
					// 	} else if typeParts[0] == "ident" {
					// 		var identName string
					// 		encoder.DeserializeRaw(*argsR[i].Value, &identName)
					// 		if rightTyp, err := GetIdentType(identName, line, fileName, cxt); err == nil {
					// 			if outType != ptrs + rightTyp {
					// 				panic(fmt.Sprintf("%s: %d: identifier '%s' was previously declared as '%s'; cannot use type '%s' in assignment", fileName, line, outName, outType, ptrs + rightTyp))
					// 			}
					// 		}
					// 	} else {
					// 		if outType != typ {
					// 			panic(fmt.Sprintf("%s: %d: identifier '%s' was previously declared as '%s'; cannot use type '%s' in assignment", fileName, line, outName, outType, ptrs + typ))
					// 		}
					// 	}
					// }

					if len(typeParts) > 1 || typeParts[0] == "ident" {
						var identName string
						encoder.DeserializeRaw(*argsR[i].Value, &identName)
						identName = ptrs + identName
						sIdentName := encoder.Serialize(identName)
						arg := MakeArgument(&sIdentName, typ)
						expr.AddInput(arg)
					} else {
						arg := MakeArgument(argsR[i].Value, typ)
						expr.AddInput(arg)
					}

					// arg := MakeArgument(argsR[i].Value, typ)
					// expr.AddInput(arg)
					
					expr.AddOutputName(outName)
				}
			} else {
				// +=, -=, *=, etc.
				var opName string
				var typName string

				if secondTyp == "ident" {
					var identName string
					encoder.DeserializeRaw(*argsR[i].Value, &identName)

					if argTyp, err := GetIdentType(identName, line, fileName, cxt); err == nil {
						typName = argTyp
					} else {
						panic(err)
					}
				} else if secondTyp == "" {
					typName = typ
				} else {
					typName = secondTyp
				}

				switch op {
				case "+=":
					opName = "add"
				case "-=":
					opName = "sub"
				case "*=":
					opName = "mul"
				case "/=":
					opName = "div"
				case "%=":
					opName = "mod"
				case "**=":
					opName = "pow"
				case "<<=":
					opName = "bitshl"
				case ">>=":
					opName = "bitshr"
				case "&=":
					opName = "bitand"
				case "^=":
					opName = "bitxor"
				case "|=":
					opName = "bitor"
				}

				if op, err := cxt.GetFunction(fmt.Sprintf("%s.%s", typName, opName), CORE_MODULE); err == nil {
					expr := MakeExpression(op)
					if !replMode {
						expr.FileLine = line
						expr.FileName = fileName
					}

					fn.AddExpression(expr)
					expr.AddLabel(tag)
					tag = ""

					var outName string
					encoder.DeserializeRaw(*argL.Value, &outName)

					// checking if identifier was previously declared
					if outType, err := GetIdentType(outName, line, fileName, cxt); err == nil {
						if len(typeParts) > 1 {
							if outType != secondTyp {
								panic(fmt.Sprintf("%s: %d: identifier '%s' was previously declared as '%s'; cannot use type '%s' in assignment", fileName, line, outName, outType, secondTyp))
							}
						} else if typeParts[0] == "ident" {
							var identName string
							encoder.DeserializeRaw(*argsR[i].Value, &identName)
							if rightTyp, err := GetIdentType(identName, line, fileName, cxt); err == nil {
								if outType != rightTyp {
									panic(fmt.Sprintf("%s: %d: identifier '%s' was previously declared as '%s'; cannot use type '%s' in assignment", fileName, line, outName, outType, rightTyp))
								}
							}
						} else {
							if outType != typ {
								panic(fmt.Sprintf("%s: %d: identifier '%s' was previously declared as '%s'; cannot use type '%s' in assignment", fileName, line, outName, outType, typ))
							}
						}
					}

					// needs to be in this order or addoutputname won't know the type when the operator is identity()
					expr.AddInput(argL)
					
					if len(typeParts) > 1 {
						expr.AddInput(MakeArgument(argsR[i].Value, "ident"))
					} else {
						expr.AddInput(MakeArgument(argsR[i].Value, typeParts[0]))
					}
					expr.AddOutputName(outName)
				}
			}
		}
	}
}

func NonAssignFunctionCall (ident string, args []*CXArgument, line int) []*CXArgument {
	var modName string
	var fnName string
	var err error
	var isMethod bool
	//var receiverType string
	identParts := strings.Split(ident, ".")
	
	if len(identParts) == 2 {
		mod, _ := cxt.GetCurrentPackage()
		if typ, err := GetIdentType(identParts[0], line, fileName, cxt); err == nil {
			// then it's a method call
			if IsStructInstance(typ, mod) {
				isMethod = true
				//receiverType = typ
				modName = mod.Name
				fnName = fmt.Sprintf("%s.%s", typ, identParts[1])
			}
		} else {
			// then it's a module
			modName = identParts[0]
			fnName = identParts[1]
		}
	} else {
		fnName = identParts[0]
		mod, e := cxt.GetCurrentPackage()
		modName = mod.Name
		err = e
	}

	found := false
	currModName := ""
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		currModName = mod.Name
		for _, imp := range mod.Imports {
			if modName == imp.Name {
				found = true
				break
			}
		}
	}

	isModule := false
	if _, err := cxt.GetModule(modName); err == nil {
		isModule = true
	}
	
	if !found && !IsNative(modName + "." + fnName) && modName != currModName && isModule {
		fmt.Printf("%s: %d: module '%s' was not imported or does not exist\n", fileName, line, modName)
	} else {
		if err == nil {
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				if op, err := cxt.GetFunction(fnName, modName); err == nil {
					expr := MakeExpression(op)
					if !replMode {
						expr.FileLine = line
						expr.FileName = fileName
					}
					fn.AddExpression(expr)
					expr.AddLabel(tag)
					tag = ""

					if isMethod {
						sIdent := encoder.Serialize(identParts[0])
						args = append([]*CXArgument{MakeArgument(&sIdent, "ident")}, args...)
					}
					
					for _, arg := range args {
						typeParts := strings.Split(arg.Typ, ".")

						arg.Typ = typeParts[0]
						expr.AddInput(arg)
					}

					lenOut := len(op.Outputs)
					outNames := make([]string, lenOut)
					args := make([]*CXArgument, lenOut)
					
					for i, out := range op.Outputs {
						outNames[i] = MakeGenSym(NON_ASSIGN_PREFIX)
						byteName := encoder.Serialize(outNames[i])
						args[i] = MakeArgument(&byteName, fmt.Sprintf("ident.%s", out.Typ))

						expr.AddOutputName(outNames[i])
					}
					
					return args
				} else {
					fmt.Printf("%s: %d: function '%s' not defined\n", fileName, line, ident)
				}
			}
		}
	}
	panic("")
}

func StatementReturn (retArg []*CXArgument, line int) {
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		if fn, err := mod.GetCurrentFunction(); err == nil {
			if len(retArg) > len(fn.Outputs) {
				panic(fmt.Sprintf("%s: %d: too many arguments to return", fileName, line))
			}
			if len(retArg) < len(fn.Outputs) {
				panic(fmt.Sprintf("%s: %d: not enough arguments to return", fileName, line))
			}
			if retArg != nil {
				for i, arg := range retArg {
					var typ string
					identParts := strings.Split(arg.Typ, ".")

					typ = identParts[0]

					// if len(identParts) > 1 {
					// 	typ = identParts[0]
					// } else {
					// 	typ = identParts[0]
					// }
					
					var idFn string
					if IsBasicType(typ) {
						idFn = MakeIdentityOpName(typ)
					} else {
						idFn = "identity"
					}

					if op, err := cxt.GetFunction(idFn, CORE_MODULE); err == nil {
						expr := MakeExpression(op)
						if !replMode {
							expr.FileLine = line
							expr.FileName = fileName
						}
						fn.AddExpression(expr)
						if idFn == "identity" {
							expr.AddInput(MakeArgument(arg.Value, "str"))
						} else {
							expr.AddInput(MakeArgument(arg.Value, typ))
						}

						var resolvedType string
						if typ == "ident" {
							var identName string
							encoder.DeserializeRaw(*arg.Value, &identName)
							if resolvedType, err = GetIdentType(identName, line, fileName, cxt); err != nil {
								panic(err)
							}
						} else {
							resolvedType = typ
						}

						if resolvedType != fn.Outputs[i].Typ {
							panic(fmt.Sprintf("%s: %d: wrong output type", fileName, line))
						}
						
						expr.AddOutputName(fn.Outputs[i].Name)
					}
				}
			}
			if goToFn, err := cxt.GetFunction("baseGoTo", CORE_MODULE); err == nil {
				expr := MakeExpression(goToFn)
				if !replMode {
					expr.FileLine = line
					expr.FileName = fileName
				}
				fn.AddExpression(expr)
				val := MakeDefaultValue("bool")
				expr.AddInput(MakeArgument(val, "bool"))
				lines := encoder.SerializeAtomic(int32(-len(fn.Expressions)))
				expr.AddInput(MakeArgument(&lines, "i32"))
				expr.AddInput(MakeArgument(&lines, "i32"))
			}
		}
	}
}

func StatementGoTo (ident string, line int) {
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		if fn, err := mod.GetCurrentFunction(); err == nil {
			// this one is goTo, not baseGoTo
			if goToFn, err := cxt.GetFunction("goTo", mod.Name); err == nil {
				expr := MakeExpression(goToFn)
				if !replMode {
					expr.FileLine = line
					expr.FileName = fileName
				}
				fn.AddExpression(expr)

				//label := []byte(ident)
				label := encoder.Serialize(ident)
				expr.AddInput(MakeArgument(&label, "str"))
			}
		}
	}
}

func StatementIfCondition (line int) {
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		if fn, err := mod.GetCurrentFunction(); err == nil {
			if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
				expr := MakeExpression(goToFn)
				if !replMode {
					expr.FileLine = line
					expr.FileName = fileName
				}
				fn.AddExpression(expr)
			}
		}
	}
}

func StatementIfElse (numStatements int, condStatement []*CXArgument, numElse int) {
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		if fn, err := mod.GetCurrentFunction(); err == nil {
			goToExpr := fn.Expressions[numStatements - 1]

			var elseLines []byte
			if numElse > 0 {
				elseLines = encoder.Serialize(int32(len(fn.Expressions) - numStatements - numElse + 1))
			} else {
				elseLines = encoder.Serialize(int32(len(fn.Expressions) - numStatements + 1))
			}
			
			thenLines := encoder.Serialize(int32(1))

			var typ string
			if len(condStatement[0].Typ) > len("ident.") && condStatement[0].Typ[:len("ident.")] == "ident." {
				typ = "ident"
			} else {
				typ = condStatement[0].Typ
			}

			//goToExpr.AddInput(MakeArgument(predVal, "ident"))
			goToExpr.AddInput(MakeArgument(condStatement[0].Value, typ))
			goToExpr.AddInput(MakeArgument(&thenLines, "i32"))
			goToExpr.AddInput(MakeArgument(&elseLines, "i32"))
		}
	}
}

func StatementForCondExpression (line int) {
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		if fn, err := mod.GetCurrentFunction(); err == nil {
			if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
				expr := MakeExpression(goToFn)
				if !replMode {
					expr.FileLine = line
					expr.FileName = fileName
				}
				fn.AddExpression(expr)
			}
		}
	}
}

func StatementForFinalizer (preNumExprs int, cond []*CXArgument, fnNumExprs int, isExpression bool, line int) {
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		if fn, err := mod.GetCurrentFunction(); err == nil {
			goToExpr := fn.Expressions[fnNumExprs]
			elseLines := encoder.Serialize(int32(len(fn.Expressions) - fnNumExprs + 1))
			thenLines := encoder.Serialize(int32(1))

			if isExpression {
				predVal := cond[0].Value
				goToExpr.AddInput(MakeArgument(predVal, "ident"))
			} else {
				goToExpr.AddInput(cond[0])
			}
			//if multiple value return, take first one for condition
			
			goToExpr.AddInput(MakeArgument(&thenLines, "i32"))
			goToExpr.AddInput(MakeArgument(&elseLines, "i32"))
			
			if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
				goToExpr := MakeExpression(goToFn)
				if !replMode {
					goToExpr.FileLine = line
					goToExpr.FileName = fileName
				}
				fn.AddExpression(goToExpr)

				elseLines := encoder.Serialize(int32(0))
				thenLines := encoder.Serialize(int32(-len(fn.Expressions) + preNumExprs + 1))

				alwaysTrue := encoder.Serialize(int32(1))

				goToExpr.AddInput(MakeArgument(&alwaysTrue, "bool"))
				goToExpr.AddInput(MakeArgument(&thenLines, "i32"))
				goToExpr.AddInput(MakeArgument(&elseLines, "i32"))
			}
			
		}
	}
}

func StatementForCondLenExpressions (line int) int {
	if fn, err := cxt.GetCurrentFunction(); err == nil {
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			if fn, err := mod.GetCurrentFunction(); err == nil {
				if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
					expr := MakeExpression(goToFn)
					if !replMode {
						expr.FileLine = line
						expr.FileName = fileName
					}
					fn.AddExpression(expr)
				}
			}
		} else {
			fmt.Println(err)
		}

		return len(fn.Expressions)
	}
	panic("")
}

func StatementForLoopAssignLenExpressions (condControl []*CXArgument, condLenExprs int, assignExpr bool, line int) int {
	if fn, err := cxt.GetCurrentFunction(); err == nil {
		goToExpr := fn.Expressions[condLenExprs - 1]
		if assignExpr {
			if mod, err := cxt.GetCurrentPackage(); err == nil {
				if fn, err := mod.GetCurrentFunction(); err == nil {
					if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
						expr := MakeExpression(goToFn)
						if !replMode {
							expr.FileLine = line
							expr.FileName = fileName
						}
						fn.AddExpression(expr)
					}
				}
			}

			thenLines := encoder.Serialize(int32(len(fn.Expressions) - condLenExprs + 1))
			// elseLines := encoder.Serialize(int32(0)) // this is added later in basicTyp2

			predVal := condControl[0].Value
			
			goToExpr.AddInput(MakeArgument(predVal, "ident"))
			goToExpr.AddInput(MakeArgument(&thenLines, "i32"))
			// goToExpr.AddInput(MakeArgument(&elseLines, "i32"))
		}
		return len(fn.Expressions)
	}
	panic("")
}

func StatementForThreePartsFinalizer (condControl []*CXArgument, condLenExprs int, lenBeforeLoop int, assignLenExprs int, assignExpr bool, line int) {
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		if fn, err := mod.GetCurrentFunction(); err == nil {
			goToExpr := fn.Expressions[assignLenExprs - 1]

			if assignExpr {
				predVal := condControl[0].Value

				thenLines := encoder.Serialize(int32(-(assignLenExprs - lenBeforeLoop) + 1))
				elseLines := encoder.Serialize(int32(0))

				goToExpr.AddInput(MakeArgument(predVal, "bool"))
				goToExpr.AddInput(MakeArgument(&thenLines, "i32"))
				goToExpr.AddInput(MakeArgument(&elseLines, "i32"))

				if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
					goToExpr := MakeExpression(goToFn)
					if !replMode {
						goToExpr.FileLine = line
						goToExpr.FileName = fileName
					}
					fn.AddExpression(goToExpr)

					alwaysTrue := encoder.Serialize(int32(1))

					thenLines := encoder.Serialize(int32(-len(fn.Expressions) + condLenExprs) + 1)
					elseLines := encoder.Serialize(int32(0))

					goToExpr.AddInput(MakeArgument(&alwaysTrue, "bool"))
					goToExpr.AddInput(MakeArgument(&thenLines, "i32"))
					goToExpr.AddInput(MakeArgument(&elseLines, "i32"))

					condGoToExpr := fn.Expressions[condLenExprs - 1]

					condThenLines := encoder.Serialize(int32(len(fn.Expressions) - condLenExprs + 1))
					
					condGoToExpr.AddInput(MakeArgument(&condThenLines, "i32"))
				}
			} else {
				predVal := condControl[0].Value

				thenLines := encoder.Serialize(int32(1))
				elseLines := encoder.Serialize(int32(len(fn.Expressions) - condLenExprs + 2))
				
				goToExpr.AddInput(MakeArgument(predVal, "ident"))
				goToExpr.AddInput(MakeArgument(&thenLines, "i32"))
				goToExpr.AddInput(MakeArgument(&elseLines, "i32"))

				if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
					goToExpr := MakeExpression(goToFn)
					if !replMode {
						goToExpr.FileLine = line
						goToExpr.FileName = fileName
					}
					fn.AddExpression(goToExpr)
					
					alwaysTrue := encoder.Serialize(int32(1))

					thenLines := encoder.Serialize(int32(-len(fn.Expressions) + lenBeforeLoop + 1))
					elseLines := encoder.Serialize(int32(0))

					goToExpr.AddInput(MakeArgument(&alwaysTrue, "bool"))
					goToExpr.AddInput(MakeArgument(&thenLines, "i32"))
					goToExpr.AddInput(MakeArgument(&elseLines, "i32"))
				}
			}
		}
	}
}

func VariableDeclaration (varName string, typName string, line int) {
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			if op, err := cxt.GetFunction("initDef", mod.Name); err == nil {
				expr := MakeExpression(op)
				if !replMode {
					expr.FileLine = line
					expr.FileName = fileName
				}

				fn.AddExpression(expr)

				typ := encoder.Serialize(typName)
				arg := MakeArgument(&typ, "str")
				expr.AddInput(arg)
				expr.AddOutputName(varName)
			}
		}
	}
}

func ElseStatementInitializer (line int) {
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		if fn, err := mod.GetCurrentFunction(); err == nil {
			if goToFn, err := cxt.GetFunction("baseGoTo", mod.Name); err == nil {
				expr := MakeExpression(goToFn)
				if !replMode {
					expr.FileLine = line
					expr.FileName = fileName
				}
				fn.AddExpression(expr)
			}
		}
	}
}

func ElseStatementFinalizer (beforeElseLenExprs int) int {
	if mod, err := cxt.GetCurrentPackage(); err == nil {
		if fn, err := mod.GetCurrentFunction(); err == nil {
			goToExpr := fn.Expressions[beforeElseLenExprs - 1]
			
			elseLines := encoder.Serialize(int32(0))
			thenLines := encoder.Serialize(int32(len(fn.Expressions) - beforeElseLenExprs + 1))

			alwaysTrue := encoder.Serialize(int32(1))

			goToExpr.AddInput(MakeArgument(&alwaysTrue, "bool"))
			goToExpr.AddInput(MakeArgument(&thenLines, "i32"))
			goToExpr.AddInput(MakeArgument(&elseLines, "i32"))

			return len(fn.Expressions) - beforeElseLenExprs
		}
	}
	panic("")
}

func UnaryPrefixOp (arg *CXArgument, nonAssignExpr []*CXArgument, isArgument bool, line int) *CXArgument {
	if isArgument {
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			if op, err := cxt.GetFunction("not", CORE_MODULE); err == nil {
				expr := MakeExpression(op)
				if !replMode {
					expr.FileLine = line
					expr.FileName = fileName
				}
				fn.AddExpression(expr)
				expr.AddLabel(tag)
				tag = ""
				expr.AddInput(arg)

				outName := MakeGenSym(NON_ASSIGN_PREFIX)
				byteName := encoder.Serialize(outName)
				
				expr.AddOutputName(outName)
				return MakeArgument(&byteName, "ident")
			}
		}
	} else {
		if fn, err := cxt.GetCurrentFunction(); err == nil {
			if op, err := cxt.GetFunction("not", CORE_MODULE); err == nil {
				expr := MakeExpression(op)
				if !replMode {
					expr.FileLine = line
					expr.FileName = fileName
				}
				fn.AddExpression(expr)
				expr.AddLabel(tag)
				tag = ""

				if (len(nonAssignExpr[0].Typ) > len("ident.") && nonAssignExpr[0].Typ[:len("ident.")] == "ident.") {
					nonAssignExpr[0].Typ = "ident"
				}
				
				expr.AddInput(nonAssignExpr[0])

				outName := MakeGenSym(NON_ASSIGN_PREFIX)
				byteName := encoder.Serialize(outName)
				
				expr.AddOutputName(outName)
				return MakeArgument(&byteName, "ident")
			}
		}
	}

	panic("")
}

func StructLiteralDeclaration (ident string, structFlds []*CXArgument, line int) *CXArgument {
	var result *CXArgument
	val := encoder.Serialize(ident)
	
	if len(structFlds) < 1 {
		result = MakeArgument(&val, "ident")
	} else {
		// then it's a struct literal
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			if fn, err := cxt.GetCurrentFunction(); err == nil {
				if op, err := cxt.GetFunction("initDef", mod.Name); err == nil {
					expr := MakeExpression(op)
					if !replMode {
						expr.FileLine = line
						expr.FileName = fileName
					}
					fn.AddExpression(expr)

					outName := MakeGenSym(NON_ASSIGN_PREFIX)
					sOutName := encoder.Serialize(outName)

					typ := encoder.Serialize(ident)
					expr.AddInput(MakeArgument(&typ, "str"))
					expr.AddOutputName(outName)

					result = MakeArgument(&sOutName, fmt.Sprintf("ident.%s", ident))
					for _, def := range structFlds {
						typeParts := strings.Split(def.Typ, ".")

						var typ string
						var secondTyp string
						var idFn string

						if len(typeParts) > 1 {
							typ = "str"
							secondTyp = typeParts[1] // i32, f32, etc
						} else if typeParts[0] == "ident" {
							typ = "str"
							secondTyp = "ident"
						} else {
							typ = typeParts[0] // i32, f32, etc
						}

						if secondTyp == "" {
							idFn = MakeIdentityOpName(typ)
						} else {
							idFn = "identity"
						}
						
						if op, err := cxt.GetFunction(idFn, CORE_MODULE); err == nil {
							expr := MakeExpression(op)
							if !replMode {
								expr.FileLine = line
								expr.FileName = fileName
							}
							fn.AddExpression(expr)
							expr.AddLabel(tag)
							tag = ""

							outName := fmt.Sprintf("%s.%s", outName, def.Name)
							arg := MakeArgument(def.Value, typ)
							expr.AddInput(arg)
							expr.AddOutputName(outName)
						}
					}
				}
			}
		}
	}

	return result
}

func BasicArrayLiteralDeclaration (basicTyp string, elts []*CXArgument, line int, isEmpty bool) *CXArgument {
	if isEmpty {
		switch basicTyp {
		case "[]bool":
			vals := make([]int32, 0)
			sVal := encoder.Serialize(vals)

			return MakeArgument(&sVal, "[]bool")
		case "[]byte":
			vals := make([]byte, 0)
			sVal := encoder.Serialize(vals)
			return MakeArgument(&sVal, "[]byte")
		case "[]str":
			vals := make([]string, 0)
			sVal := encoder.Serialize(vals)
			return MakeArgument(&sVal, "[]str")
		case "[]i32":
			vals := make([]int32, 0)
			sVal := encoder.Serialize(vals)
			return MakeArgument(&sVal, "[]i32")
		case "[]i64":
			vals := make([]int64, 0)
			sVal := encoder.Serialize(vals)
			return MakeArgument(&sVal, "[]i64")
		case "[]f32":
			vals := make([]float32, 0)
			sVal := encoder.Serialize(vals)
			return MakeArgument(&sVal, "[]f32")
		case "[]f64":
			vals := make([]float64, 0)
			sVal := encoder.Serialize(vals)
			return MakeArgument(&sVal, "[]f64")
		}
	} else {
		if mod, err := cxt.GetCurrentPackage(); err == nil {
			if fn, err := cxt.GetCurrentFunction(); err == nil && inFn {
				if op, err := cxt.GetFunction(INIT_FN, mod.Name); err == nil {
					expr := MakeExpression(op)
					if !replMode {
						expr.FileLine = line
						expr.FileName = fileName
					}

					outName := MakeGenSym(NON_ASSIGN_PREFIX)
					sOutName := encoder.Serialize(outName)

					fn.AddExpression(expr)

					var appendFnTyp string
					var ptrs string
					if basicTyp[0] == '*' {
						for i, char := range basicTyp {
							if char != '*' {
								appendFnTyp = basicTyp[i:]
								break
							} else {
								ptrs += "*"
							}
						}
					} else {
						appendFnTyp = basicTyp
					}

					typ := encoder.Serialize(appendFnTyp)
					arg := MakeArgument(&typ, "str")
					expr.AddInput(arg)
					expr.AddOutputName(outName)
					
					if op, err := cxt.GetFunction(fmt.Sprintf("%s.append", appendFnTyp), mod.Name); err == nil {
						for _, arg := range elts {
							typeParts := strings.Split(arg.Typ, ".")
							arg.Typ = typeParts[0]
							expr := MakeExpression(op)
							fn.AddExpression(expr)
							expr.AddInput(MakeArgument(&sOutName, "ident"))
							expr.AddOutputName(outName)
							expr.AddInput(CastArgumentForArray(appendFnTyp, arg))
						}
					}

					return MakeArgument(&sOutName, ptrs + "ident")
				}
			} else {
				// then it's for a global definition
				switch basicTyp {
				case "[]str":
					vals := make([]string, len(elts))
					for i, arg := range elts {
						var val string
						encoder.DeserializeRaw(*arg.Value, &val)
						vals[i] = val
					}
					sVal := encoder.Serialize(vals)
					return MakeArgument(&sVal, "[]str")
				case "[]bool":
					vals := make([]int32, len(elts))
					for i, arg := range elts {
						var val int32
						encoder.DeserializeRaw(*arg.Value, &val)
						vals[i] = val
					}
					sVal := encoder.Serialize(vals)
					return MakeArgument(&sVal, "[]bool")
				case "[]byte":
					vals := make([]byte, len(elts))
					for i, arg := range elts {
						var val int32
						encoder.DeserializeRaw(*arg.Value, &val)
						vals[i] = byte(val)
					}
					sVal := encoder.Serialize(vals)
					return MakeArgument(&sVal, "[]byte")
				case "[]i32":
					vals := make([]int32, len(elts))
					for i, arg := range elts {
						var val int32
						encoder.DeserializeRaw(*arg.Value, &val)
						vals[i] = val
					}
					sVal := encoder.Serialize(vals)
					return MakeArgument(&sVal, "[]i32")
				case "[]i64":
					vals := make([]int64, len(elts))
					for i, arg := range elts {
						var val int32
						encoder.DeserializeRaw(*arg.Value, &val)
						vals[i] = int64(val)
					}
					sVal := encoder.Serialize(vals)
					return MakeArgument(&sVal, "[]i64")
				case "[]f32":
					vals := make([]float32, len(elts))
					for i, arg := range elts {
						var val float32
						encoder.DeserializeRaw(*arg.Value, &val)
						vals[i] = val
					}
					sVal := encoder.Serialize(vals)
					return MakeArgument(&sVal, "[]f32")
				case "[]f64":
					vals := make([]float64, len(elts))
					for i, arg := range elts {
						var val float32
						encoder.DeserializeRaw(*arg.Value, &val)
						vals[i] = float64(val)
					}
					sVal := encoder.Serialize(vals)
					return MakeArgument(&sVal, "[]f64")
				}
			}
		}
	}

	panic("")
}
