package base

import (
	"fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// This method removes 1 or more expressions
// Then adds the equivalent of removed expressions
func (cxt *cxContext) MutateSolution (numberExprs int) {
	cxt.SelectFunction("solution")
	if fn, err := cxt.GetCurrentFunction(); err == nil {
		// removing (we also need to remove those expressions which contain the output var of this expr)
		removeIndex := random(0, len(fn.Expressions))
		removedVar := fn.Expressions[removeIndex].OutputName
		indexesToRemove := make([]int, 0)

		for i, expr := range fn.Expressions {
			for _, arg := range expr.Arguments {
				if i == removeIndex {
					indexesToRemove = append([]int{i}, indexesToRemove...)
					break
				}

				if i > removeIndex && string(*arg.Value) == removedVar {
					indexesToRemove = append([]int{i}, indexesToRemove...)
					break
				}

				// if, for example, we removed var25 = add(var1, var2) before
				// we also need to check for all the subsequent expressions
				// that could be using var25, and remove them too
				for _, arg := range fn.Expressions[i].Arguments {
					broke := false
					for _, j := range indexesToRemove {
						if string(*arg.Value) == fn.Expressions[j].OutputName {
							indexesToRemove = append([]int{i}, indexesToRemove...)
							broke = true
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
		// we need to make sure at least one expression throws a result to the output var

		hasOutput := false
		for _, expr := range fn.Expressions {
			if expr.OutputName == fn.Output.Name {
				hasOutput = true
				break
			}
		}

		for i := 0; len(fn.Expressions) < numberExprs; i++ {
			//affs := FilterAffordances(fn.GetAffordances(), "Expression")
			affs := make([]*cxAffordance, 0)
			
			if !hasOutput && i == numberExprs - 1 {
				// making sure last expression assigns to output
				affs = FilterAffordances(fn.GetAffordances(), "Expression", concat(fn.Output.Name, " ="))
			} else {
				affs = FilterAffordances(fn.GetAffordances(), "Expression")
			}
			
			affs[random(0, len(affs))].ApplyAffordance()
		}
	}
}

func (cxt *cxContext) EvolveSolution (inputs, outputs []int32, numberExprs, iterations int) *cxContext {
	cxt.SelectFunction("solution")

	// Initializing expressions
	if fn, err := cxt.GetCurrentFunction(); err == nil {
		for i := 0; i < numberExprs; i++ {
			affs := make([]*cxAffordance, 0)
			if i != numberExprs - 1 {
				affs = FilterAffordances(fn.GetAffordances(), "Expression")
			} else {
				// making sure last expression assigns to output
				affs = FilterAffordances(fn.GetAffordances(), "Expression", concat(fn.Output.Name, " ="))
				// for _, aff := range affs {
				// 	fmt.Println(aff.Description)
				// }
			}
			r := random(0, len(affs))
			//fmt.Println(r)
			affs[r].ApplyAffordance()
		}
	}

	best := cxt
	// printing initial solution
	//best.PrintProgram(false)
	
	if fn, err := cxt.GetCurrentFunction(); err == nil {
		if len(fn.Inputs) > 0 && fn.Output != nil {
			for i := 0; i < iterations; i++ {
				//fmt.Printf("Iteration:%d\n", i)
				// getting 4 copies of the parent (the best solution)
				// let's create an array holding these programs
				programs := make([]*cxContext, 5) // it's always 5, because of the 1+4 strategy
				errors := make([]float64, 5)
				programs[0] = best // the best solution will always be at index 0

				for i := 1; i < 5; i++ {
					programs[i] = MakeContextCopy(programs[0], 0)
					//fmt.Printf("%p ", programs[i])
					// we need to mutate these 4
					programs[i].MutateSolution(numberExprs)
					//fmt.Println(programs[i].CurrentModule.CurrentFunction.Expressions[0].OutputName)
					//programs[i].PrintProgram(false)
					
				}

				//fmt.Println()

				// we need to evaluate all of them with each of the inputs and get the error
				for i, program := range programs {
					//fmt.Printf("Program:%d\n", i)					
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
						output := program.CallStack[0].State["outMain"].Value
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
