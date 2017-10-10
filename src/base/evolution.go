package base

import (
	"fmt"
	"regexp"
	"bytes"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// This method removes 1 or more expressions
// Then adds the equivalent of removed expressions

func (expr *CXExpression) BuildExpression (outNames []*CXDefinition) *CXExpression {
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
		expr.OutputNames = outNames
	}
	
	return expr
}

func (cxt *CXProgram) MutateSolution (solutionName string, fnBag string, numberExprs int) {
	cxt.SelectFunction(solutionName)
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

				// // testing if always removing the output expression gives better results
				// if expr.OutputNames[0] == fn.Outputs[0].Name {
				// 	indexesToRemove = append([]int{i}, indexesToRemove...)
				// 	break
				// }

				argIdentRemoved := false
				for _, removedVar := range removedVars {
					if i > removeIndex && string(*arg.Value) == removedVar.Name {
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
							if string(*arg.Value) == outName.Name {
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
				if out.Name != exprOutNames[i].Name {
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
			var affs []*CXAffordance
			if !hasOutputs && i == numberExprs - lenFnExprs - 1 {

				var outs bytes.Buffer
				for j, out := range fn.Outputs {
					if j == len(fn.Outputs) - 1 {
						outs.WriteString(concat(out.Typ))
					} else {
						outs.WriteString(concat(out.Typ, ", "))
					}
				}

				affs = FilterAffordances(fn.GetAffordances(), "Expression", concat("\\(", regexp.QuoteMeta(outs.String()), "\\)$"), fnBag)
			} else {
				affs = FilterAffordances(fn.GetAffordances(), "Expression", fnBag)
			}
			

			// excluding array operations
			// re := regexp.MustCompile("readByteA|writeByteA|evolve")
			// filteredAffs := make([]*CXAffordance, 0)
			// for _, aff := range affs {
			// 	if re.FindString(aff.Description) == "" {
			// 		filteredAffs = append(filteredAffs, aff)
			// 	}
			// }
			// affs = filteredAffs
			
			r := random(0, len(affs))
			affs[r].ApplyAffordance()

			if expr, err := fn.GetCurrentExpression(); err == nil {
				if !hasOutputs && i == numberExprs - lenFnExprs - 1 {
					outNames := make([]*CXDefinition, len(fn.Outputs))
					for i, out := range fn.Outputs {
						outDef := MakeDefinition(
							out.Name,
							MakeDefaultValue(out.Typ),
							out.Typ)
						outNames[i] = outDef
					}
					expr.BuildExpression(outNames)
				} else {
					expr.BuildExpression(nil)
				}
			}
		}
	}
}

func (cxt *CXProgram) adaptPreEvolution (solutionName string) {
	if mod, err := cxt.GetModule("main"); err == nil {
		if _, err := cxt.GetFunction("main", "main"); err == nil {
			//delete(mod.Functions, "main")
			idx := -1
			for i, fn := range mod.Functions {
				if fn.Name == "main" {
					idx = i
					break
				}
			}
			mod.Functions = append(mod.Functions[:idx], mod.Functions[idx+1:]...)
		} else {
			panic(err)
		}

		if solMod, err := cxt.GetCurrentModule(); err == nil {
			if mod, err := cxt.GetModule("main"); err == nil {
				mod.AddFunction(MakeFunction("main"))
				if fn, err := cxt.GetCurrentFunction(); err == nil {
					fn.AddOutput(MakeParameter("out", "f64"))
					if sol, err := cxt.GetFunction(solutionName, solMod.Name); err == nil {
						fn.AddExpression(MakeExpression(sol))
						
						if expr, err := cxt.GetCurrentExpression(); err == nil {
							expr.AddOutputName("out")
							//tmpVal := encoder.SerializeAtomic(float64(0))
							tmpVal := encoder.Serialize(float64(0))
							expr.AddArgument(MakeArgument(&tmpVal, "f64"))
						}
					}
				}
			}
		}
	}
}

func (cxt *CXProgram) adaptInput (testValue float64) {
	if fn, err := cxt.GetFunction("main", "main"); err == nil {
		//val := encoder.SerializeAtomic(testValue)
		val := encoder.Serialize(testValue)
		arg := MakeArgument(&val, "f64")
		// arg.Offset = -1
		// arg.Size = -1

		fn.Expressions[0].Arguments[0] = arg
	}
}

func (fromCxt *CXProgram) transferSolution (solutionName string, toCxt *CXProgram) {
	if fromMod, err := fromCxt.GetCurrentModule(); err == nil {
		if fromFn, err := fromCxt.GetFunction(solutionName, fromMod.Name); err == nil {
			if toMod, err := toCxt.GetCurrentModule(); err == nil {
				//delete(toMod.Functions, solutionName)
				//idx := -1
				for _, fn := range toMod.Functions {
					if fn.Name == solutionName {
						//idx = i
						fn.Expressions = fromFn.Expressions
						break
					}
				}
				//toMod.Functions = append(toMod.Functions[:idx], toMod.Functions[idx+1:]...)
				//toMod.AddFunction(MakeFunctionCopy(fromFn, toMod, toCxt))
			}
		}
	}
}

func (cxt *CXProgram) Evolve (solutionName string, fnBag string, inputs, outputs []float64, numberExprs, iterations int, epsilon float64, expr *CXExpression, call *CXCall) error {
	cxt.SelectFunction(solutionName)

	//cxt.PrintProgram(false)

	// Initializing expressions
	if fn, err := cxt.GetCurrentFunction(); err == nil {
		preExistingExpressions := len(fn.Expressions)
		for i := 0; i < numberExprs - preExistingExpressions; i++ {

			if i != numberExprs - preExistingExpressions - 1 {
				affs := FilterAffordances(fn.GetAffordances(), "Expression", fnBag)
				
				r := random(0, len(affs))
				affs[r].ApplyAffordance()

				if expr, err := fn.GetCurrentExpression(); err == nil {
					expr.BuildExpression(nil)
				}
			} else {
				
				fnOutTypNames := make([]string, 0)
				for _, out := range fn.Outputs {
					fnOutTypNames = append(fnOutTypNames, out.Typ)
				}

				var possibleOps string
				coreModule, _ := cxt.GetModule("core")
				for _, op := range coreModule.Functions {
					possibleOp := true
					for i, out := range op.Outputs {
						if len(op.Outputs) != len(fnOutTypNames) || out.Typ != fnOutTypNames[i] {
							possibleOp = false
							break
						}
					}
					
					if possibleOp {
						possibleOps = concat(possibleOps, concat(regexp.QuoteMeta(op.Name), "|"))
					}
				}
				for _, op := range fn.Module.Functions {
					possibleOp := true
					for i, out := range op.Outputs {
						if len(op.Outputs) != len(fnOutTypNames) || out.Typ != fnOutTypNames[i] {
							possibleOp = false
							break
						}
					}
					
					if possibleOp {
						possibleOps = concat(possibleOps, concat(regexp.QuoteMeta(op.Name), "|"))
					}
				}

				possibleOps = string([]byte(possibleOps)[:len(possibleOps)-1])

				affs := FilterAffordances(fn.GetAffordances(), "Expression", possibleOps, fnBag)

				r := random(0, len(affs))
				affs[r].ApplyAffordance()
				
				// making sure last expression assigns to output
				outNames := make([]*CXDefinition, len(fn.Outputs))
				for i, out := range fn.Outputs {
					outDef := MakeDefinition(
						out.Name,
						MakeDefaultValue(out.Typ),
						out.Typ)
					outNames[i] = outDef
				}

				if expr, err := fn.GetCurrentExpression(); err == nil {
					expr.BuildExpression(outNames)
				} else {
					fmt.Println(err)
				}
			}
		}
	}

	cxtCopy := MakeContextCopy(cxt, -1)
	
	best := cxtCopy
	best.adaptPreEvolution(solutionName)
	// printing initial solution
	//best.PrintProgram(false)

	var finalError float64
	
	if mod, err := cxt.GetCurrentModule(); err == nil {
		if fn, err := cxt.GetFunction(solutionName, mod.Name); err == nil {
			if len(fn.Inputs) > 0 && len(fn.Outputs) > 0 {
				fmt.Printf("Evolving function '%s'\n", solutionName)
				for i := 0; i < iterations; i++ {
					// fmt.Printf("Iteration:%d\n", i)
					// getting 4 copies of the parent (the best solution)
					// let's create an array holding these programs
					programs := make([]*CXProgram, 5) // it's always 5, because of the 1+4 strategy
					errors := make([]float64, 5)
					programs[0] = best // the best solution will always be at index 0

					for i := 1; i < 5; i++ {
						programs[i] = MakeContextCopy(programs[0], 0)
						// we need to mutate these 4
						programs[i].MutateSolution(solutionName, fnBag, numberExprs)
						//programs[i].PrintProgram(false)
					}

					// we need to evaluate all of them with each of the inputs and get the error
					for i, program := range programs {
						//fmt.Printf("Program:%d\n", i)
						var error float64 = 0
						for i, inp := range inputs {
							program.adaptInput(inp)
							//fmt.Println(program.Modules[0].Functions[0].Expressions[0].Arguments[0].Value)

							//program.PrintProgram(false)
							//fmt.Println(inp)

							program.Reset()
							program.Run(false, -1)

							// getting the simulated output
							var result float64
							// We'll always take the first value
							// The algorithm shouldn't work with multiple value returns yet

							output := *program.Outputs[0].Value
							encoder.DeserializeRaw(output, &result)

							//fmt.Println(program.Outputs[0].Value)

							// I don't want to import Math, so I will hardcode abs
							diff := float64(result - outputs[i])

							// fmt.Println(result)
							//fmt.Println(program.Outputs[0].Value)
							// fmt.Println(outputs[i])
							// fmt.Println()
							
							if diff >= 0 {
								error += diff
							} else {
								error += (diff * float64(-1))
							}
							//fmt.Println(error)

						}
						
						//fmt.Println(len(program.CallStack[len(program.CallStack) - 1].State))
						//fmt.Println(error / float64(len(inputs)))
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
					//fmt.Println(errors)
					
					//best.PrintProgram(false)
					best = programs[bestIndex]
					//best.PrintProgram(false)
					finalError = errors[bestIndex]

					if errors[bestIndex] < epsilon {
						break
					}
				}
				fmt.Printf("Finished evolving function '%s'\n", solutionName)
			}
		}
	}

	best.transferSolution(solutionName, cxt)

	sFinalError := encoder.Serialize(finalError)
	
	for _, def := range call.State {
		if def.Name == expr.OutputNames[0].Name {
			def.Value = &sFinalError
			return nil
		}
	}
	
	call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &sFinalError, "f64"))
	return nil
}

func RandomProgram (numberAffordances int) *CXProgram {
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
				affs := make([]*CXAffordance, 0)
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
				affs := make([]*CXAffordance, 0)
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
				affs := make([]*CXAffordance, 0)
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
				affs := make([]*CXAffordance, 0)
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
