package base

import (
	"fmt"
	"errors"
	"strconv"
	"strings"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func condOperation (operator string, stack []string, affs []*CXAffordance, exec *[]bool, call *CXCall) error {
	obj1 := strings.Split(stack[len(stack) - 2], ".")
	obj2 := strings.Split(stack[len(stack) - 1], ".")
	
	oneIsX := false
	obj1IsX := false
	obj2IsX := false
	var toCompare float64
	bothX := false
	if obj1[0] == "x" && obj2[0] == "x" && len(obj1) == 1 && len(obj2) == 1 {
		bothX = true
	}
	
	if obj1[0] == "x" {
		oneIsX = true
		obj1IsX = true
	} else {
		if f, err := strconv.ParseFloat(obj1[0], 64); err == nil {
			toCompare = f
		} else {
			if f, err := strconv.ParseFloat(obj1[0], 64); err == nil {
				toCompare = f
			} else {
				// then it's an identifier
				if arg, err := resolveIdent(obj1[0], call); err == nil {
					if len(obj1) > 1 {
						// then it's a struct
						if mod, err := call.Context.GetCurrentModule(); err == nil {
							if strct, err := call.Context.GetStruct(arg.Typ, mod.Name); err == nil {
								id1Val, id1Typ, _, _ := resolveStructField(obj1[1], arg.Value, strct)
								switch id1Typ {
								case "i32":
									var val int32
									encoder.DeserializeAtomic(id1Val, &val)
									toCompare = float64(val)
								}
							}
						}
					} else {
						// then it's an atomic
						switch arg.Typ {
						case "i32":
							var val int32
							encoder.DeserializeAtomic(*arg.Value, &val)
							toCompare = float64(val)
						}
					}
				}
			}
		}
	}
	if obj2[0] == "x" {
		oneIsX = true
		obj2IsX = true
	} else {
		if f, err := strconv.ParseFloat(obj2[0], 64); err == nil {
			toCompare = f
		} else {
			// then it's an identifier
			if arg, err := resolveIdent(obj2[0], call); err == nil {
				if len(obj2) > 1 {
					// then it's a struct
					if mod, err := call.Context.GetCurrentModule(); err == nil {
						if strct, err := call.Context.GetStruct(arg.Typ, mod.Name); err == nil {
							id2Val, id2Typ, _, _ := resolveStructField(obj2[1], arg.Value, strct)
							switch id2Typ {
								case "i32":
									var val int32
									encoder.DeserializeAtomic(id2Val, &val)
									toCompare = float64(val)
								}
						}
					}
				} else {
					switch arg.Typ {
					case "i32":
						var val int32
						encoder.DeserializeAtomic(*arg.Value, &val)
						toCompare = float64(val)
					}
				}
			}
		}
	}
	
	if !oneIsX {
		return errors.New("aff.query: malformed then block")
	}

	//resolving the identifier
	for _, aff := range affs {
		if bothX {
			*exec = append(*exec, true)
			continue
		}
		// we need to check that it isn't basic
		isBasic := false
		for _, basic := range BASIC_TYPES {
			if basic == aff.Typ {
				isBasic = true
			}
		}

		var obj1Val []byte
		var obj1Typ string
		var obj2Val []byte
		var obj2Typ string

		//fmt.Println(aff.Typ)

		if !isBasic {
			var typ string
			if aff.Typ[:2] == "[]" {
				typ = aff.Typ[2:]
			} else {
				typ = aff.Typ
			}
			if mod, err := call.Context.GetCurrentModule(); err == nil {
				if strct, err := call.Context.GetStruct(typ, mod.Name); err == nil {
					if arg, err := resolveIdent(aff.Name, call); err == nil {
						if aff.Index != "" {
							if i, err := strconv.ParseInt(aff.Index, 10, 64); err == nil {
								if expr, err := call.Context.GetCurrentExpression(); err == nil {
									if val, err, _, _ := getStrctFromArray(arg, int32(i), expr, call); err == nil {
										if obj1IsX {
											obj1Val, obj1Typ, _, _ = resolveStructField(obj1[1], &val, strct)
										}
										if obj2IsX {
											obj2Val, obj2Typ, _, _ = resolveStructField(obj2[1], &val, strct)
										}
									}
								} else {
									return err
								}
							} else {
								return err
							}
							
						} else {
							if obj1IsX {
								obj1Val, obj1Typ, _, _ = resolveStructField(obj1[1], arg.Value, strct)
							}
							if obj2IsX {
								obj2Val, obj2Typ, _, _ = resolveStructField(obj2[1], arg.Value, strct)
							}
						}
					}
				}
			}
		} else {
			if arg, err := resolveIdent(aff.Name, call); err == nil {
				if aff.Index != "" {
					if i, err := strconv.ParseInt(aff.Index, 10, 64); err == nil {
						if obj1IsX {
							if val, err := getValueFromArray(arg, int32(i)); err == nil {
								obj1Val = val
								obj1Typ = arg.Typ[2:]
							} else {
								return err
							}
						}
						if obj2IsX {
							if val, err := getValueFromArray(arg, int32(i)); err == nil {
								obj2Val = val
								obj2Typ = arg.Typ[2:]
							} else {
								return err
							}
						}
					} else {
						return err
					}
				} else {
					if obj1IsX {
						obj1Val = *arg.Value
						obj1Typ = arg.Typ
					}
					if obj2IsX {
						obj2Val = *arg.Value
						obj2Typ = arg.Typ
					}
				}
			}
		}

		// if we're comparing strings, we just compare the values byte by byte
		// and == is the only available operator
		//fmt.Println(obj1Typ)
		if obj1IsX && obj1Typ == "str" {
			sObj2 := encoder.Serialize(obj2[0])
			
			isEqual := true
			if len(obj1Val) != len(sObj2) {
				isEqual = false
			} else {
				for i, byt := range obj1Val {
					if byt != sObj2[i] {
						isEqual = false
						break
					}
				}
			}
			
			if isEqual {
				*exec = append(*exec, true)
				continue
			} else {
				*exec = append(*exec, false)
				continue
			}
		}
		if obj2IsX && obj2Typ == "str" {
			sObj1 := encoder.Serialize(obj1[0])
			
			isEqual := true
			if len(obj2Val) != len(sObj1) {
				isEqual = false
			} else {
				for i, byt := range obj1Val {
					if byt != sObj1[i] {
						isEqual = false
						break
					}
				}
			}
			if isEqual {
				*exec = append(*exec, true)
				continue
			} else {
				*exec = append(*exec, false)
				continue
			}
		}

		var pv float64
		var pv2 float64
		
		if !obj1IsX {
			switch obj2Typ {
			case "i32", "bool":
				var v int32
				encoder.DeserializeRaw(obj2Val, &v)
				pv = float64(v)
				
			}
		}
		if !obj2IsX {
			switch obj1Typ {
			case "i32", "bool":
				var v int32
				encoder.DeserializeRaw(obj1Val, &v)
				pv = float64(v)
			}
		}
		if obj1IsX && obj2IsX {
			switch obj1Typ { // both must be of same type
			case "i32", "bool":
				var v1 int32
				var v2 int32
				encoder.DeserializeRaw(obj1Val, &v1)
				encoder.DeserializeRaw(obj2Val, &v2)
				pv = float64(v1)
				pv2 = float64(v2)
			}
		}

		isTrue := false
		if !obj1IsX {
			switch operator {
			case ">": isTrue = toCompare > pv && !bothX
			case "<": isTrue = toCompare < pv && !bothX
			case "==": isTrue = toCompare == pv || bothX
			}
		}
		if !obj2IsX {
			switch operator {
			case ">": isTrue = pv > toCompare && !bothX
			case "<": isTrue = pv < toCompare && !bothX
			case "==": isTrue = pv == toCompare || bothX
			}
		}
		if obj1IsX && obj2IsX {
			switch operator {
			case ">": isTrue = pv > pv2 && !bothX
			case "<": isTrue = pv < pv2 && !bothX
			case "==": isTrue = pv == pv2 || bothX
			}
		}

		if isTrue {
			*exec = append(*exec, true)
		} else {
			*exec = append(*exec, false)
		}
	}
	return nil
}

func aff_query (target, objects, rules *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("aff.query", "[]str", "[]str", "[]str", target, objects, rules); err == nil {
		var _target []string
		var _objects []string
		var _rules []string

		encoder.DeserializeRaw(*target.Value, &_target)
		encoder.DeserializeRaw(*objects.Value, &_objects)
		encoder.DeserializeRaw(*rules.Value, &_rules)

		pObjs := make([]string, 0)
		pWeights := make([]float64, 0)

		objs := make([]string, 0)
		weights := make([]float64, 0)
		stack := make([]string, 0)

		// parsing _objects
		for i, obj := range _objects {
			if obj == "weight" {
				w := _objects[i - 1]
				obj1 := _objects[i - 2]
				obj2 := strings.Split(w, ".")
				if f, err := strconv.ParseFloat(w, 64); err == nil {
					pObjs = append(pObjs, obj1)
					pWeights = append(pWeights, f)
				} else {
					// then it's an identifier
					if arg, err := resolveIdent(obj2[0], call); err == nil {
						if len(obj2) > 1 {
							// then it's a struct
							if mod, err := call.Context.GetCurrentModule(); err == nil {
								if strct, err := call.Context.GetStruct(arg.Typ, mod.Name); err == nil {
									id2Val, id2Typ, _, _ := resolveStructField(obj2[1], arg.Value, strct)
									switch id2Typ {
									case "i32":
										var val int32
										encoder.DeserializeAtomic(id2Val, &val)

										pObjs = append(pObjs, obj1)
										pWeights = append(pWeights, float64(val))
									case "f32":
										var val float32
										encoder.DeserializeRaw(id2Val, &val)

										pObjs = append(pObjs, obj1)
										pWeights = append(pWeights, float64(val))
									}
								}
							}
						} else {
							switch arg.Typ {
							case "i32":
								var val int32
								encoder.DeserializeAtomic(*arg.Value, &val)

								pObjs = append(pObjs, obj1)
								pWeights = append(pWeights, float64(val))
							case "f32":
								var val float32
								encoder.DeserializeRaw(*arg.Value, &val)

								pObjs = append(pObjs, obj1)
								pWeights = append(pWeights, float64(val))
							}
						}
					}
					//return errors.New("aff.query: weight badly formatted in objects list")
				}
			}
		}

		var prevMod string
		var prevFn string
		var prevStrct string
		
		if mod, err := call.Context.GetCurrentModule(); err == nil {
			prevMod = mod.Name
			if fn, err := mod.GetCurrentFunction(); err == nil {
				prevFn = fn.Name
			}
			if strct, err := mod.GetCurrentStruct(); err == nil {
				prevStrct = strct.Name
			}
		}

		stack = nil
		var targetTyp string // strct, expr, etc
		for _, t := range _target {
			switch t {
			case "pkg":
				obj := stack[len(stack) - 1]
				call.Context.SelectModule(obj)
				targetTyp = "pkg"
				stack = stack[:len(stack) - 1]
			case "fn":
				obj := stack[len(stack) - 1]
				call.Context.SelectFunction(obj)
				targetTyp = "fn"
				stack = stack[:len(stack) - 1]
			case "strct":
				obj := stack[len(stack) - 1]
				call.Context.SelectStruct(obj)
				targetTyp = "strct"
				stack = stack[:len(stack) - 1]
			case "exp":
				obj := stack[len(stack) - 1]
				
				if fn, err := call.Context.GetCurrentFunction(); err == nil {
					found := false
					for _, expr := range fn.Expressions {
						if expr.Tag == obj {
							fn.SelectExpression(expr.Line)
							found = true
							targetTyp = "exp"
							break
						}
					}

					if !found {
						return errors.New(fmt.Sprintf("aff.query: no expression with tag '%s' was found", obj))
					}
				}
				stack = stack[:len(stack) - 1]
			default:
				stack = append(stack, t)
			}
		}

		var affs []*CXAffordance
		var filteredAffs []*CXAffordance
		
		// now getting the affordances
		switch targetTyp {
		case "pkg":
			if mod, err := call.Context.GetCurrentModule(); err == nil {
				affs = mod.GetAffordances()
				//PrintAffordances(affs)
			}
		case "fn":
			if fn, err := call.Context.GetCurrentFunction(); err == nil {
				affs = fn.GetAffordances()
				//PrintAffordances(affs)
			}
		case "strct":
			if strct, err := call.Context.GetCurrentStruct(); err == nil {
				affs = strct.GetAffordances()
				//PrintAffordances(affs)
			}
		case "exp":
			if expr, err := call.Context.GetCurrentExpression(); err == nil {
				lastArg := expr.Arguments[len(expr.Arguments) - 1]
				expr.RemoveArgument()
				affs = expr.GetAffordances()
				expr.AddArgument(lastArg)
				//PrintAffordances(affs)
			}
		default:
			return errors.New("aff.query: no target was specified")
		}

		falseRule := false
		var exec []bool // result of >, <, ==, etc

		for _, rule := range _rules {
			switch rule {
			case "weight":
				if falseRule {
					continue
				}
				w := stack[len(stack) - 1]
				obj1 := stack[len(stack) - 2]
				obj2 := strings.Split(w, ".")

				if f, err := strconv.ParseFloat(w, 64); err == nil {
					objs = append(objs, obj1)
					weights = append(weights, f)
				} else {
					// then it's an identifier
					if arg, err := resolveIdent(obj2[0], call); err == nil {
						if len(obj2) > 1 {
							// then it's a struct
							if mod, err := call.Context.GetCurrentModule(); err == nil {
								if strct, err := call.Context.GetStruct(arg.Typ, mod.Name); err == nil {
									id2Val, id2Typ, _, _ := resolveStructField(obj2[1], arg.Value, strct)
									switch id2Typ {
									case "i32":
										var val int32
										encoder.DeserializeAtomic(id2Val, &val)

										objs = append(objs, obj1)
										weights = append(weights, float64(val))
									}
								}
							}
						} else {
							switch arg.Typ {
							case "i32":
								var val int32
								encoder.DeserializeAtomic(*arg.Value, &val)

								objs = append(objs, obj1)
								weights = append(weights, float64(val))
							}
						}
					}
					//return errors.New("aff.query: weight badly formatted in rules list")
				}
			case "single":
				if falseRule {
					continue
				}
				found := false
				obj2 := objs[len(objs) - 1]

				w2 := weights[len(weights) - 1]

				if obj2 == "true" {
					found = true
				} else {
					for i, pObj := range pObjs {
						if (obj2 == pObj && pWeights[i] >= w2) {
							found = true
							break
						}
					}
				}

				objs = objs[:len(objs) - 1]
				weights = weights[:len(weights) - 1]
				
				if found {
					objs = append(objs, "true")
				} else {
					objs = append(objs, "false")
				}
				weights = append(weights, w2)
			case "obj":
				if falseRule {
					continue
				}
				obj1 := objs[len(objs) - 1]
				w1 := weights[len(weights) - 1]

				found := false
				for i, obj := range pObjs {
					// if already present, just replace weight
					if obj == obj1 {
						pWeights[i] = w1
						found = true
					}
				}
				if !found {
					pObjs = append(pObjs, obj1)
					pWeights = append(pWeights, w1)
				}
			case "or":
				if falseRule {
					continue
				}
				found := false
				obj1 := objs[len(objs) - 2]
				obj2 := objs[len(objs) - 1]

				w1 := weights[len(weights) - 2]
				w2 := weights[len(weights) - 1]

				if obj1 == "true" || obj2 == "true" {
					found = true
				} else {
					for i, pObj := range pObjs {
						if (obj1 == pObj && pWeights[i] >= w1) || (obj2 == pObj && pWeights[i] >= w2) {
							found = true
							break
						}
					}
				}

				objs = objs[:len(objs) - 2]
				weights = weights[:len(weights) - 2]

				var orWeight float64
				if w1 > w2 {
					orWeight = w1
				} else {
					orWeight = w2
				}
				
				if found {
					objs = append(objs, "true")
				} else {
					objs = append(objs, "false")
				}
				weights = append(weights, orWeight)
			case "and":
				if falseRule {
					continue
				}
				
				found1 := false
				found2 := false
				
				obj1 := objs[len(objs) - 2]
				obj2 := objs[len(objs) - 1]

				w1 := weights[len(weights) - 2]
				w2 := weights[len(weights) - 1]

				// I need to check if either obj1 or obj2 are present in
				if obj1 == "true" && obj2 == "true" {
					found1 = true
					found2 = true
				} else {
					for i, pObj := range pObjs {
						if (obj1 == pObj || obj1 == "true") && pWeights[i] >= w1 {
							found1 = true
						}
						if (obj2 == pObj || obj2 == "true") && pWeights[i] >= w2 {
							found2 = true
						}

						if found1 && found2 {
							break
						}
					}
				}
				
				objs = objs[:len(objs) - 2]
				weights = weights[:len(weights) - 2]

				var andWeight float64
				if w1 < w2 {
					andWeight = w1
				} else {
					andWeight = w2
				}
				
				if found1 && found2 {
					objs = append(objs, "true")
				} else {
					objs = append(objs, "false")
				}
				weights = append(weights, andWeight)
			case "x":
				if falseRule {
					continue
				}
				stack = append(stack, rule)
			case ">":
				if falseRule {
					continue
				}
				if err := condOperation(rule, stack, affs, &exec, call); err != nil {
					return err
				}
			case "<":
				if falseRule {
					continue
				}
				if err := condOperation(rule, stack, affs, &exec, call); err != nil {
					return err
				}
			case "==":
				if falseRule {
					continue
				}
				if err := condOperation(rule, stack, affs, &exec, call); err != nil {
					return err
				}
			case "allow":
				if falseRule {
					continue
				}
				for i, aff := range affs {
					if exec[i] {
						alreadyThere := false
						for _, fAff := range filteredAffs {
							if fAff == aff {
								alreadyThere = true
								break
							}
						}

						if !alreadyThere {
							filteredAffs = append(filteredAffs, aff)
						}
					}
				}
				exec = nil
			case "reject":
				if falseRule {
					continue
				}
				for i, aff := range affs {
					if exec[i] {
						for i, fAff := range filteredAffs {
							if fAff == aff {
								
								if i != len(filteredAffs) - 1 {
									filteredAffs = append(filteredAffs[:i], filteredAffs[i+1:]...)
								} else {
									filteredAffs = filteredAffs[:i]
								}
								break
							}
						}
					}
				}

				exec = nil
			case "if":
				if falseRule {
					continue
				}
				if len(objs) > 2 {
					return errors.New("aff.query: malformed predicate")
				}
				pred := objs[len(objs) - 1]

				if pred != "true" {
					falseRule = true
				}

				stack = nil
				objs = nil
			case "endif":
				// fmt.Println(pObjs)
				// fmt.Println(pWeights)
				// fmt.Println()
				// PrintAffordances(filteredAffs)
				stack = nil
				objs = nil
				falseRule = false
			default:
				if falseRule {
					continue
				}
				stack = append(stack, rule)
			}
		}

		// restoring previous selected cx objects

		if prevMod != "" {
			call.Context.SelectModule(prevMod)
		}
		if prevFn != "" {
			call.Context.SelectFunction(prevFn)
		}
		if prevStrct != "" {
			call.Context.SelectStruct(prevStrct)
		}

		// making commands
		var commands []string
		for _, aff := range filteredAffs {
			op := aff.Operator
			name := aff.Name
			index := aff.Index
			//typ := aff.Typ

			cmd := []string{
				"startcmd", op, "op", name, "name", index, "index", "endcmd",
			}
			commands = append(commands, cmd...)
		}

		output := encoder.Serialize(commands)
		assignOutput(&output, "[]str", expr, call)
		return nil
	} else {
		return err
	}
}

func aff_execute (target, commands, index *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("aff.execute", "[]str", "[]str", "i32", target, commands, index); err == nil {
		var _target []string
		var _commands []string
		var _index int32

		encoder.DeserializeRaw(*target.Value, &_target)
		encoder.DeserializeRaw(*commands.Value, &_commands)
		encoder.DeserializeAtomic(*index.Value, &_index)
		
		var op string
		var name string
		var index string
		var counter int

		isSkip := true
		for i, cmd := range _commands {
			switch cmd {
			case "startcmd":
				if int(_index) == counter {
					isSkip = false
				} else {
					isSkip = true
				}
				counter++
			case "op":
				if !isSkip {
					op = _commands[i - 1]
				}
			case "name":
				if !isSkip {
					name = _commands[i - 1]
				}
			case "index":
				if !isSkip {
					index = _commands[i - 1]
				}
			default:
				
			}
		}

		// now execute the command
		var prevMod string
		var prevFn string
		var prevStrct string
		
		if mod, err := call.Context.GetCurrentModule(); err == nil {
			prevMod = mod.Name
			if fn, err := mod.GetCurrentFunction(); err == nil {
				prevFn = fn.Name
			}
			if strct, err := mod.GetCurrentStruct(); err == nil {
				prevStrct = strct.Name
			}
		}

		stack := make([]string, 0)
		var targetTyp string // strct, expr, etc
		for _, t := range _target {
			switch t {
			case "pkg":
				obj := stack[len(stack) - 1]
				call.Context.SelectModule(obj)
				targetTyp = "pkg"
				stack = stack[:len(stack) - 1]
			case "fn":
				obj := stack[len(stack) - 1]
				call.Context.SelectFunction(obj)
				targetTyp = "fn"
				stack = stack[:len(stack) - 1]
			case "strct":
				obj := stack[len(stack) - 1]
				call.Context.SelectStruct(obj)
				targetTyp = "strct"
				stack = stack[:len(stack) - 1]
			case "exp":
				obj := stack[len(stack) - 1]
				
				if fn, err := call.Context.GetCurrentFunction(); err == nil {
					found := false
					for _, expr := range fn.Expressions {
						if expr.Tag == obj {
							fn.SelectExpression(expr.Line)
							found = true
							targetTyp = "exp"
							break
						}
					}

					if !found {
						return errors.New(fmt.Sprintf("aff.query: no expression with tag '%s' was found", obj))
					}
				}
				stack = stack[:len(stack) - 1]
			default:
				stack = append(stack, t)
			}
		}

		// now getting the affordances
		switch targetTyp {
		case "pkg":
			// if mod, err := call.Context.GetCurrentModule(); err == nil {
			
			// }
		case "fn":
			// if fn, err := call.Context.GetCurrentFunction(); err == nil {
			
			// }
		case "strct":
			// if strct, err := call.Context.GetCurrentStruct(); err == nil {
			
			// }
		case "exp":
			if expr, err := call.Context.GetCurrentExpression(); err == nil {
				// one CX object can have different types of affordances
				switch op {
				case "AddArgument":
					if index != "" {
						if arr, err := resolveIdent(name, call); err == nil {
							if i, err := strconv.ParseInt(index, 10, 64); err == nil {

								var val []byte
								var err error
								if isBasicType(arr.Typ) {
									val, err = getValueFromArray(arr, int32(i))
								} else {
									val, err, _, _ = getStrctFromArray(arr, int32(i), expr, call)
								}
								
								if err == nil {
									expr.RemoveArgument()
									expr.AddArgument(MakeArgument(&val, arr.Typ[2:]))
								} else {
									return err
								}
							} else {
								return err
							}
						} else {
							return err
						}
					} else {
						sName := encoder.Serialize(name)
						expr.RemoveArgument()
						expr.AddArgument(MakeArgument(&sName, "ident"))
					}
				}
			}
		default:
			return errors.New("aff.execute: no target was specified")
		}


		if prevMod != "" {
			call.Context.SelectModule(prevMod)
		}
		if prevFn != "" {
			call.Context.SelectFunction(prevFn)
		}
		if prevStrct != "" {
			call.Context.SelectStruct(prevStrct)
		}

		return nil
	} else {
		return err
	}
}

// prints affordances in a human readable format
func aff_print (commands *CXArgument, call *CXCall) error {
	if err := checkType("aff.print", "[]str", commands); err == nil {
		var _commands []string
		encoder.DeserializeRaw(*commands.Value, &_commands)

		var counter int
		for i, cmd := range _commands {
			switch cmd {
			case "op":
				fmt.Printf("(%d)\tOperator: %s\t", counter, _commands[i - 1])
				counter++
			case "name":
				fmt.Printf("Name: %s\t",  _commands[i - 1])
			case "index":
				if _commands[i - 1] != "" {
					fmt.Printf("Index: %s\t",  _commands[i - 1])
				}
			case "endcmd":
				fmt.Println()
			default:
			}
		}
		
		return nil
	} else {
		return err
	}
}

func aff_len (commands *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("aff.len", "[]str", commands); err == nil {
		var _commands []string
		encoder.DeserializeRaw(*commands.Value, &_commands)

		var counter int
		for _, cmd := range _commands {
			if cmd == "op" {
				counter++
			}
		}
		
		output := encoder.Serialize(int32(counter))
		assignOutput(&output, "i32", expr, call)
		return nil
	} else {
		return err
	}
}
