package base

import (
	"fmt"
	"regexp"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// This method removes 1 or more expressions
// Then adds the equivalent of removed expressions

func (expr *cxExpression) BuildExpression (outNames []string) *cxExpression {
	numInps := len(expr.Operator.Inputs)
	numOuts := len(expr.Operator.Outputs)

	// for _, aff := range FilterAffordances(expr.GetAffordances(), "Argument") {
	// 	fmt.Println(aff.Description)
	// }

	for i := 0; i < numInps; i++ {
		affs := FilterAffordances(expr.GetAffordances(), "Argument")
		//fmt.Printf("Number affs %d\n", len(affs))
		r := random(0, len(affs))
		//fmt.Println(affs[r].Description)
		affs[r].ApplyAffordance()
		
	}

	// for _, arg := range expr.Arguments {
	// 	fmt.Println(string(*arg.Value))
	// }
	
	//fmt.Println()

	if len(outNames) == 0 {
		for i := 0; i < numOuts; i++ {
			affs := FilterAffordances(expr.GetAffordances(), "OutputName")
			r := random(0, len(affs))
			affs[r].ApplyAffordance()
		}
	} else {
		//fmt.Println(outNames)
		expr.OutputNames = outNames
		// for _, outName := range outNames {
		// 	fmt.Println("THIS")
		// 	fmt.Println(FilterAffordances(expr.GetAffordances(), "OutputName")[0].Description)
			
		// 	affs := FilterAffordances(expr.GetAffordances(), "OutputName", outName)
		// 	r := random(0, len(affs))
		// 	affs[r].ApplyAffordance()
		// }
	}
	
	return expr
}

func (cxt *cxContext) MutateSolution (numberExprs int) {
	cxt.SelectFunction("solution")
	if fn, err := cxt.GetCurrentFunction(); err == nil {
		// removing (we also need to remove those expressions which contain the output var of this expr)
		removeIndex := random(0, len(fn.Expressions))
		//removedVar := fn.Expressions[removeIndex].OutputName

		removedVars := fn.Expressions[removeIndex].OutputNames
		
		indexesToRemove := make([]int, 0)

		for i, expr := range fn.Expressions {
			for _, arg := range expr.Arguments {
				if i == removeIndex {
					indexesToRemove = append([]int{i}, indexesToRemove...)
					break
				}

				argIdentRemoved := false
				for _, removedVar := range removedVars {
					if i > removeIndex && string(*arg.Value) == removedVar {
						indexesToRemove = append([]int{i}, indexesToRemove...)
						argIdentRemoved = true
					}
				}
				if argIdentRemoved {
					break
				}
				

				// if, for example, we removed var25 = add(var1, var2) before
				// we also need to check for all the subsequent expressions
				// that could be using var25, and remove them too
				for _, arg := range fn.Expressions[i].Arguments {
					broke := false
					for _, j := range indexesToRemove {
						argIsOutput := false
						for _, outName := range fn.Expressions[j].OutputNames {
							if string(*arg.Value) == outName {
								indexesToRemove = append([]int{i}, indexesToRemove...)
								broke = true
								argIsOutput = true
							}
						}
						if argIsOutput {
							break
						}
					}
					if broke {
						break
					}
				}
			}
		}

		indexesToRemove = removeDuplicatesInt(indexesToRemove)

		//fmt.Println(indexesToRemove)

		for _, i := range indexesToRemove {
			if i != len(fn.Expressions) {
				fn.Expressions = append(fn.Expressions[:i], fn.Expressions[i+1:]...)
			} else {
				fn.Expressions = fn.Expressions[:i]
			}
		}

		// re-indexing expression line numbers
		for i, expr := range fn.Expressions {
			expr.Line = i
		}

		// adding
		// we need to make sure that at least one expression throws a result to the output var





		
		hasOutputs := false
		fnOutputs := fn.Outputs
		for _, expr := range fn.Expressions {
			exprOutNames := expr.OutputNames

			if len(exprOutNames) != len(fnOutputs) {
				break
			}

			sameOutputs := true
			
			for i, out := range fnOutputs {
				if out.Name != exprOutNames[i] {
					sameOutputs = false
				}
			}

			if sameOutputs {
				hasOutputs = true
				break
			}
		}

		lenFnExprs := len(fn.Expressions)
		for i := 0; len(fn.Expressions) < numberExprs; i++ {
			affs := FilterAffordances(fn.GetAffordances(), "Expression")

			// excluding array operations
			re := regexp.MustCompile("readAByte|writeAByte")
			filteredAffs := make([]*cxAffordance, 0)
			for _, aff := range affs {
				if re.FindString(aff.Description) == "" {
					filteredAffs = append(filteredAffs, aff)
				}
			}
			affs = filteredAffs
			
			r := random(0, len(affs))
			affs[r].ApplyAffordance()

			if expr, err := fn.GetCurrentExpression(); err == nil {
				if !hasOutputs && i == numberExprs - lenFnExprs - 1 {
					outNames := make([]string, len(fn.Outputs))
					for i, out := range fn.Outputs {
						outNames[i] = out.Name
					}
					expr.BuildExpression(outNames)
				} else {
					expr.BuildExpression(nil)
				}
			}
		}
	}
}

func (cxt *cxContext) EvolveSolution (inputs, outputs []int32, numberExprs, iterations int) *cxContext {
	cxt.SelectFunction("solution")

	// Initializing expressions
	if fn, err := cxt.GetCurrentFunction(); err == nil {
		for i := 0; i < numberExprs; i++ {

			if i != numberExprs - 1 {
				affs := FilterAffordances(fn.GetAffordances(), "Expression")
				
				// excluding array operations
				re := regexp.MustCompile("readAByte|writeAByte")
				filteredAffs := make([]*cxAffordance, 0)
				for _, aff := range affs {
					if re.FindString(aff.Description) == "" {
						filteredAffs = append(filteredAffs, aff)
					}
				}
				affs = filteredAffs
				
				r := random(0, len(affs))
				affs[r].ApplyAffordance()

				if expr, err := fn.GetCurrentExpression(); err == nil {
					expr.BuildExpression(nil)
				}
			} else {
				
				fnOutTypNames := make([]string, 0)
				for _, out := range fn.Outputs {
					fnOutTypNames = append(fnOutTypNames, out.Typ.Name)
				}

				//possibleOps := make([]byte, 0)
				var possibleOps string
				for _, op := range fn.Module.Functions {
					possibleOp := true
					for i, out := range op.Outputs {
						if len(op.Outputs) != len(fnOutTypNames) || out.Typ.Name != fnOutTypNames[i] {
							possibleOp = false
							break
						}
					}
					
					if possibleOp {
						possibleOps = concat(possibleOps, concat(regexp.QuoteMeta(op.Name), "|"))
					}
				}

				possibleOps = string([]byte(possibleOps)[:len(possibleOps)-1])

				affs := FilterAffordances(fn.GetAffordances(), "Expression", possibleOps)

				// excluding array operations
				re := regexp.MustCompile("readAByte|writeAByte")
				filteredAffs := make([]*cxAffordance, 0)
				for _, aff := range affs {
					if re.FindString(aff.Description) == "" {
						filteredAffs = append(filteredAffs, aff)
					}
				}
				affs = filteredAffs
				
				r := random(0, len(affs))
				affs[r].ApplyAffordance()
				
				// making sure last expression assigns to output
				outNames := make([]string, len(fn.Outputs))
				for i, out := range fn.Outputs {
					outNames[i] = out.Name
				}

				if expr, err := fn.GetCurrentExpression(); err == nil {
					expr.BuildExpression(outNames)
				} else {
					fmt.Println(err)
				}
			}
		}
	}

	best := cxt
	// printing initial solution
	//best.PrintProgram(false)
	
	if fn, err := cxt.GetCurrentFunction(); err == nil {
		if len(fn.Inputs) > 0 && len(fn.Outputs) > 0 {
			for i := 0; i < iterations; i++ {
				//fmt.Printf("Iteration:%d\n", i)
				// getting 4 copies of the parent (the best solution)
				// let's create an array holding these programs
				programs := make([]*cxContext, 5) // it's always 5, because of the 1+4 strategy
				errors := make([]float64, 5)
				programs[0] = best // the best solution will always be at index 0

				for i := 1; i < 5; i++ {
					programs[i] = MakeContextCopy(programs[0], 0)
					// we need to mutate these 4
					programs[i].MutateSolution(numberExprs)
					//programs[i].PrintProgram(false)
				}

				// we need to evaluate all of them with each of the inputs and get the error
				for i, program := range programs {
					//fmt.Printf("Program:%d\n", i)
					//program.PrintProgram(false)
					var error float64 = 0
					for i, inp := range inputs {
						// the input is always going to be num1 for now
						num1 := encoder.SerializeAtomic(inp)
						if def, err := program.GetDefinition("num1"); err == nil {
							def.Value = &num1
						} else {
							fmt.Println(err)
						}

						program.Reset()
						program.Run(false, -1)

						// getting the simulated output
						var result int32
						// We'll always take the first value
						// The algorithm shouldn't work with multiple value returns yet
						output := program.Outputs[0].Value
						encoder.DeserializeAtomic(*output, &result)

						// I don't want to import Math, so I will hardcode abs
						diff := float64(result - outputs[i])
						//fmt.Println(diff)
						if diff >= 0 {
							error += diff
						} else {
							error += diff * -1
						}
						//fmt.Println(error)

					}
					
					//fmt.Println(len(program.CallStack[len(program.CallStack) - 1].State))
					errors[i] = error / float64(len(inputs))
				}

				//fmt.Println(errors)
				
				// the program with the lowest error becomes the best
				bestIndex := 0
				for i, _ := range programs {
					if errors[i] <= errors[bestIndex] && errors[i] >= 0 {
						bestIndex = i
					}
				}
				//fmt.Println(errors)
				//fmt.Println(bestIndex)

				// print error each iteration
				fmt.Println(errors[bestIndex])
				
				//best.PrintProgram(false)
				best = programs[bestIndex]
				// we can't get any lower if error == 0
				if errors[bestIndex] == 0 {
					break
				}
			}
		}
	}

	return best
}

func RandomProgram (numberAffordances int) *cxContext {
	cxt := MakeContext()
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
			mod, err := cxt.GetCurrentModule()
			if err == nil {
				affs := make([]*cxAffordance, 0)
				if mod != nil {
					affs = mod.GetAffordances()
				}
				if len(affs) > 0 {
					affs[random(0, len(affs))].ApplyAffordance()
				}
			}
		case randomCase >= 15 && randomCase < 30:
			fn, err := cxt.GetCurrentFunction()
			if err == nil {
				affs := make([]*cxAffordance, 0)
				if fn != nil {
					affs = fn.GetAffordances()
				}
				if len(affs) > 0 {
					affs[random(0, len(affs))].ApplyAffordance()
				}
			}
		case randomCase >= 50 && randomCase < 60:
			strct, err := cxt.GetCurrentStruct()
			if err == nil {
				affs := make([]*cxAffordance, 0)
				if strct != nil {
					affs = strct.GetAffordances()
				}
				
				if len(affs) > 0 {
					affs[random(0, len(affs))].ApplyAffordance()
				}
			}
		case randomCase >= 60 && randomCase < 100:
			expr, err := cxt.GetCurrentExpression()
			if err == nil {
				affs := make([]*cxAffordance, 0)
				if expr != nil {
					affs = expr.GetAffordances()
				}
				if len(affs) > 0 {
					affs[random(0, len(affs))].ApplyAffordance()
				}
			}
		}
	}

	return cxt
}
