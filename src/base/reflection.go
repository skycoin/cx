package base

import (
	//"github.com/skycoin/skycoin/src/cipher/encoder"
	"fmt"
	"bytes"
	"regexp"
	"strconv"
	"math/rand"
	"time"
)

func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max - min) + min
}

func RandomProgram (numberAffordances int) *cxContext {
	cxt := MakeContext()

	// Basic initialization to not waste affordances in these
	// FilterAffordances(cxt.GetAffordances(),
	// 	"AddModule")[0].ApplyAffordance()
	// FilterAffordances(cxt.GetCurrentModule().GetAffordances(),
	// 	"AddFunction")[0].ApplyAffordance()
	// FilterAffordances(cxt.GetCurrentModule().GetAffordances(),
	// 	"AddStruct")[0].ApplyAffordance()
	// FilterAffordances(cxt.GetCurrentModule().GetAffordances(),
	// 	"AddDefinition")[0].ApplyAffordance()

	// We could could randomly choose among the select operators
	// Then we apply a random affordance to that selection
	
	for i := 0; i < numberAffordances; i++ {
		randomCase := random(0, 100)
		// 0: Affordances on Context // Merge case 0 and 1
		// 1: Affordances on Module
		// 2: Affordances on Function
		// 3: Affordances on Struct
		// 4: Affordances on Expression

		// let's give different weights to the options
		
		switch {
		case randomCase >= 0 && randomCase < 5:
			affs := cxt.GetAffordances()
			if len(affs) > 0 {
				affs[random(0, len(affs))].ApplyAffordance()
			}
		case randomCase >= 5 && randomCase < 15 :
			mod := cxt.GetCurrentModule()
			affs := make([]*cxAffordance, 0)
			if mod != nil {
				affs = mod.GetAffordances()
			}
			if len(affs) > 0 {
				affs[random(0, len(affs))].ApplyAffordance()
			}
		case randomCase >= 15 && randomCase < 30:
			fn := cxt.GetCurrentFunction()
			affs := make([]*cxAffordance, 0)
			if fn != nil {
				affs = fn.GetAffordances()
			}
			if len(affs) > 0 {
				affs[random(0, len(affs))].ApplyAffordance()
			}
		case randomCase >= 50 && randomCase < 60:
			strct := cxt.GetCurrentStruct()
			affs := make([]*cxAffordance, 0)
			if strct != nil {
				affs = strct.GetAffordances()
			}
			
			if len(affs) > 0 {
				affs[random(0, len(affs))].ApplyAffordance()
			}
		case randomCase >= 60 && randomCase < 100:
			expr := cxt.GetCurrentExpression()
			affs := make([]*cxAffordance, 0)
			if expr != nil {
				affs = expr.GetAffordances()
			}
			if len(affs) > 0 {
				affs[random(0, len(affs))].ApplyAffordance()
			}
		}
	}

	return cxt
}

func FilterAffordances(affs []*cxAffordance, filters ...string) []*cxAffordance {
	filteredAffs := make([]*cxAffordance, 0)
	for _, filter := range filters {
		re := regexp.MustCompile(filter)
		for _, aff := range affs {
			if re.FindString(aff.Description) != "" {
				filteredAffs = append(filteredAffs, aff)
			}
		}
		affs = filteredAffs
		filteredAffs = make([]*cxAffordance, 0)
	}
	return affs
}

func concat (strs ...string) string {
	var buffer bytes.Buffer
	
	for i := 0; i < len(strs); i++ {
		buffer.WriteString(strs[i])
	}
	
	return buffer.String()
}

func PrintAffordances (affs []*cxAffordance) {
	for _, aff := range affs {
		fmt.Println(aff.Description)
	}
}

func (aff *cxAffordance) ApplyAffordance () {
	aff.Action()
}

// Just a function to debug for myself. Not meant to be used in production nor meant to be efficient/maintanable
func (cxt *cxContext) PrintProgram(withAffs bool) {

	fmt.Println("Context")
	if withAffs {
		for i, aff := range cxt.GetAffordances() {
			fmt.Printf(" * %d.- %s\n", i, aff.Description)
		}
	}
	
	for i, mod := range cxt.Modules {
		fmt.Printf("%d.- Module: %s\n", i, mod.Name)

		if withAffs {
			for i, aff := range mod.GetAffordances() {
				fmt.Printf("\t * %d.- %s\n", i, aff.Description)
			}
		}

		if len(mod.Definitions) > 0 {
			fmt.Println("\tDefinitions")
		}
		
		j := 0
		for _, v := range mod.Definitions {
			fmt.Printf("\t\t%d.- Definition: %s %s\n", j, v.Name, v.Typ.Name)
			j++
		}

		if len(mod.Structs) > 0 {
			fmt.Println("\tStructs")
		}

		j = 0
		for _, strct := range mod.Structs {
			fmt.Printf("\t\t%d.- Struct: %s\n", j, strct.Name)

			if withAffs {
				for i, aff := range strct.GetAffordances() {
					fmt.Printf("\t\t * %d.- %s\n", i, aff.Description)
				}
			}

			for k, fld := range strct.Fields {
				fmt.Printf("\t\t\t%d.- Field: %s %s\n",
					k, fld.Name, fld.Typ.Name)
			}
			
			j++
		}

		if len(mod.Functions) > 0 {
			fmt.Println("\tFunctions")
		}

		j = 0
		for _, fn := range mod.Functions {

			inOuts := make(map[string]string)
			for _, in := range fn.Inputs {
				inOuts[in.Name] = in.Typ.Name
			}
			
			
			var inps bytes.Buffer
			//inps.WriteString(" ")
			for i, inp := range fn.Inputs {
				if i == len(fn.Inputs) - 1 {
					inps.WriteString(concat(inp.Name, " ", inp.Typ.Name))
				} else {
					inps.WriteString(concat(inp.Name, " ", inp.Typ.Name, ", "))
				}
				
			}

			out := ""
			if fn.Output != nil {
				if (fn.Output.Name != "") {
					out = concat(fn.Output.Name, " ", fn.Output.Typ.Name)
					inOuts[fn.Output.Name] = fn.Output.Typ.Name
				} else {
					out = fn.Output.Typ.Name
				}
			}
			
			fmt.Printf("\t\t%d.- Function: %s (%s) %s\n",
				j, fn.Name, inps.String(), out)

			if withAffs {
				for i, aff := range fn.GetAffordances() {
					fmt.Printf("\t\t * %d.- %s\n", i, aff.Description)
				}
			}

			k := 0
			for _, expr := range fn.Expressions {
				//Arguments
				var args bytes.Buffer

				for i, arg := range expr.Arguments {
					typ := ""
					if arg.Typ.Name == "ident" {
						if arg.Typ != nil &&
							inOuts[string(*arg.Value)] != "" {
							typ = inOuts[string(*arg.Value)]
						} else if arg.Value != nil &&
							mod.Definitions[string(*arg.Value)] != nil &&
							mod.Definitions[string(*arg.Value)].Typ.Name != "" {
							typ = mod.Definitions[string(*arg.Value)].Typ.Name
						} else {
							typ = arg.Typ.Name
						}
					}

					if i == len(expr.Arguments) - 1 {
						args.WriteString(concat(string(*arg.Value), " ", typ))
					} else {
						args.WriteString(concat(string(*arg.Value), " ", typ, ", "))
					}
					
				}

				fmt.Printf("\t\t\t%d.- Expression: %s(%s)\n",
					k, expr.Operator.Name, args.String())

				if withAffs {
					for i, aff := range expr.GetAffordances() {
						fmt.Printf("\t\t\t * %d.- %s\n", i, aff.Description)
					}
				}
				
				k++
			}
			j++
		}
	}
}

func (strct *cxStruct) GetAffordances() []*cxAffordance {
	affs := make([]*cxAffordance, 0)
	mod := strct.Module

	types := make([]string, len(basicTypes))
	copy(types, basicTypes)
	
	for name, _ := range mod.Structs {
		types = append(types, name)
	}

	// Getting types from imported modules
	for impName, imp := range mod.Imports {
		for _, strct := range imp.Structs {
			types = append(types, concat(impName, ".", strct.Name))
		}
	}

	// definitions for each available type
	for _, typ := range types {
		fldGensym := MakeGenSym("fld")
		fldType := MakeType(typ)
		
		affs = append(affs, &cxAffordance{
			Description: concat("AddField ", fldGensym, " ", typ),
			Action: func() {
				strct.AddField(MakeField(fldGensym, fldType))
			}})
	}

	return affs
}

func (expr *cxExpression) GetAffordances() []*cxAffordance {
	op := expr.Operator
	affs := make([]*cxAffordance, 0)

	// The operator for this function doesn't require arguments
	if len(op.Inputs) < 1 {
		return affs
	}
	if len(expr.Arguments) >= len(op.Inputs) {
		return affs
	}
	
	fn := expr.Function
	mod := expr.Module
	reqType := op.Inputs[len(expr.Arguments)].Typ.Name // Required type for the current op's input
	defsTypes := make([]string, 0)
	args := make([]*cxArgument, 0)
	identType := MakeType("ident")

	inOutNames := make([]string, len(fn.Inputs) + 1)
	
	// Adding inputs and outputs as definitions
	for i, param := range fn.Inputs {
		if reqType == param.Typ.Name {
			inOutNames[i] = param.Name
			defsTypes = append(defsTypes, param.Typ.Name)
			identName := []byte(param.Name)
			args = append(args, &cxArgument{
				Typ: identType,
				Value: &identName})
		}
	}

	if fn.Output != nil &&
		fn.Output.Name != "" &&
		fn.Output.Typ.Name == reqType {

		inOutNames[len(inOutNames)] = fn.Output.Name
		defsTypes = append(defsTypes, fn.Output.Typ.Name)
		identName := []byte(fn.Output.Name)
		args = append(args, &cxArgument{
			Typ: identType,
			Value: &identName})
	}

	// Adding definitions (global vars)
	for _, def := range mod.Definitions {
		if reqType == def.Typ.Name {
			// we could have a var with the same name and type in global and local
			// contexts. We only want to show 1 affordance for this name
			notDuplicated := true
			for _, name := range inOutNames {
				if name == def.Name {
					notDuplicated = false
					break
				}
			}
			
			if notDuplicated {
				defsTypes = append(defsTypes, def.Typ.Name)
				identName := []byte(def.Name)
				args = append(args, &cxArgument{
					Typ: identType,
					Value: &identName})
			}
		}
	}

	for i, arg := range args {
		affs = append(affs, &cxAffordance{
			Description: concat("AddArgument ", string(*arg.Value), " ", defsTypes[i]),
			Action: func() {
				expr.AddArgument(arg)
			}})
	}

	return affs
}

func (fn *cxFunction) GetAffordances() []*cxAffordance {
	affs := make([]*cxAffordance, 0)
	mod := fn.Module
	opsNames := make([]string, 0)
	ops := make([]*cxFunction, 0)

	types := make([]string, len(basicTypes))
	copy(types, basicTypes)
	for name, _ := range mod.Structs {
		types = append(types, name)
	}

	// Getting types from imported modules
	for impName, imp := range mod.Imports {
		for _, strct := range imp.Structs {
			types = append(types, concat(impName, ".", strct.Name))
		}
	}

	// Getting operators from current module
	for opName, op := range mod.Functions {
		if fn.Name != opName {
			ops = append(ops, op)
			opsNames = append(opsNames, opName)
		}
	}

	// Getting operators from imported modules
	for impName, imp := range mod.Imports {
		for opName, op := range imp.Functions {
			ops = append(ops, op)
			opsNames = append(opsNames, concat(impName, ".", opName))
		}
	}

	// Inputs
	for _, typ := range types {
		affs = append(affs, &cxAffordance{
			Description: concat("AddInput ", typ),
			Action: func() {
				fn.AddInput(MakeParameter(MakeGenSym("in"), MakeType(typ)))
		}})
	}

	// Output. We can only add one output
	if fn.Output == nil {
		for _, typ := range types {
			affs = append(affs, &cxAffordance{
				Description: concat("AddOutput ", typ),
				Action: func() {
					fn.AddInput(MakeParameter(MakeGenSym("in"), MakeType(typ)))
				}})
		}
	}

	// Expressions
	for i, op := range ops {
		theOp := op // or will keep reference to last op
		affs = append(affs, &cxAffordance{
			Description: concat("AddExpression ", opsNames[i]),
			Action: func() {
				fn.AddExpression(MakeExpression(theOp))
		}})
	}

	return affs
}

func (mod *cxModule) GetAffordances() []*cxAffordance {
	affs := make([]*cxAffordance, 0)
	types := make([]string, len(basicTypes))
	copy(types, basicTypes)

	if len(mod.Structs) > 0 {
		for name, _ := range mod.Structs {
			types = append(types, name)
		}
	}

	// Getting types from imported modules
	for impName, imp := range mod.Imports {
		for _, strct := range imp.Structs {
			types = append(types, concat(impName, ".", strct.Name))
		}
	}

	// definitions for each available type
	for _, typ := range types {
		defGensym := MakeGenSym("def")
		defType := MakeType(typ)
		value := []byte{}
		
		affs = append(affs, &cxAffordance{
			Description: concat("AddDefinition ", defGensym, " ", typ),
			Action: func() {
				mod.AddDefinition(MakeDefinition(defGensym, &value, defType))
			}})
	}

	// add imports
	for _, imp := range mod.Context.Modules {
		if imp.Name != mod.Name {
			affs = append(affs, &cxAffordance{
				Description: concat("AddImport ", imp.Name),
				Action: func() {
					mod.AddImport(imp)
				}})
		}
	}
	
	// add function
	fnGensym := MakeGenSym("fn")
	affs = append(affs, &cxAffordance{
		Description: concat("AddFunction ", fnGensym),
		Action: func() {
			mod.AddFunction(MakeFunction(fnGensym))
		}})

	// add structure
	strctGensym := MakeGenSym("strct")
	affs = append(affs, &cxAffordance{
		Description: concat("AddStruct ", strctGensym),
		Action: func() {
			mod.AddStruct(MakeStruct(strctGensym))
		}})
	
	return affs
}

func (cxt *cxContext) GetAffordances() []*cxAffordance {
	affs := make([]*cxAffordance, 0)
	modGensym := MakeGenSym("mod")
	
	affs = append(affs, &cxAffordance {
		Description: concat("AddModule ", modGensym),
		Action: func() {
			cxt.AddModule(MakeModule(modGensym))
		}})

	// Select module
	for _, mod := range cxt.Modules {
		modName := mod.Name
		affs = append(affs, &cxAffordance {
			Description: concat("SelectModule ", modName),
			Action: func() {
				cxt.SelectModule(modName)
			}})
	}

	// Select function from current module
	if cxt.CurrentModule != nil {
		for _, fn := range cxt.CurrentModule.Functions {
			fnName := fn.Name
			affs = append(affs, &cxAffordance {
				Description: concat("SelectFunction ", fnName),
				Action: func() {
					cxt.SelectFunction(fnName)
				}})
		}
	}

	// Select struct from current module
	if cxt.CurrentModule != nil {
		for _, strct := range cxt.CurrentModule.Structs {
			strctName := strct.Name
			affs = append(affs, &cxAffordance {
				Description: concat("SelectStruct ", strctName),
				Action: func() {
					cxt.SelectStruct(strctName)
				}})
		}
	}

	// Select expression from current function
	if cxt.CurrentModule != nil && cxt.CurrentModule.CurrentFunction != nil {
		for _, expr := range cxt.CurrentModule.CurrentFunction.Expressions {
			lineNumber := expr.Line
			line := strconv.Itoa(expr.Line)
			
			affs = append(affs, &cxAffordance {
				Description: concat("SelectExpression Line # ", line),
				Action: func() {
					cxt.SelectExpression(lineNumber)
				}})
		}
	}
	
	return affs
}

// These would be part of the "functions_of"

func (strct *cxStruct) GetFields() []*cxField {
	return strct.Fields
}

func (mod *cxModule) GetFunctions() []*cxFunction {
	funcs := make([]*cxFunction, len(mod.Functions))
	i := 0
	for _, v := range mod.Functions {
		funcs[i] = v
		i++
	}
	return funcs
}
