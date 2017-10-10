package base

import (
	"fmt"
	"time"
	"errors"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

/*
  Logical Operators
*/

func and (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("and", "bool", "bool", arg1, arg2); err == nil {
		var c1 int32
		var c2 int32
		encoder.DeserializeRaw(*arg1.Value, &c1)
		encoder.DeserializeRaw(*arg2.Value, &c2)

		var val []byte
		
		if c1 == 1 && c2 == 1 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &val
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, "bool"))
		return nil
	} else {
		return err
	}
}

func or (arg1 *CXArgument, arg2 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("or", "bool", "bool", arg1, arg2); err == nil {
		var c1 int32
		var c2 int32
		encoder.DeserializeRaw(*arg1.Value, &c1)
		encoder.DeserializeRaw(*arg2.Value, &c2)

		var val []byte
		
		if c1 == 1 || c2 == 1 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &val
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, "bool"))
		return nil
	} else {
		return err
	}
}

func not (arg1 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("not", "bool", arg1); err == nil {
		var c1 int32
		encoder.DeserializeRaw(*arg1.Value, &c1)

		var val []byte

		if c1 == 0 {
			val = encoder.Serialize(int32(1))
		} else {
			val = encoder.Serialize(int32(0))
		}

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &val
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, "bool"))
		return nil
	} else {
		return err
	}
}

// goTo increments/decrements the call.Line to the desired expression line.
// Used for if/else and loop statements.
func baseGoTo (call *CXCall, predicate *CXArgument, thenLine *CXArgument, elseLine *CXArgument) error {
	if err := checkThreeTypes("baseGoTo", "bool", "i32", "i32", predicate, thenLine, elseLine); err == nil {
		var isFalse bool

		// var pred int32
		// encoder.DeserializeRaw(*predicate.Value, &pred)

		//if pred == 0 {}
		if (*predicate.Value)[0] == 0 {
			isFalse = true
		} else {
			isFalse = false
		}

		var thenLineNo int32
		var elseLineNo int32

		encoder.DeserializeAtomic(*thenLine.Value, &thenLineNo)
		encoder.DeserializeAtomic(*elseLine.Value, &elseLineNo)

		if isFalse {
			call.Line = call.Line + int(elseLineNo) - 1
		} else {
			call.Line = call.Line + int(thenLineNo) - 1
		}

		return nil
	} else {
		return err
	}
}

func goTo (call *CXCall, tag *CXArgument) error {
	if err := checkType("goTo", "str", tag); err == nil {
		tg := string(*tag.Value)

		for _, expr := range call.Operator.Expressions {
			if expr.Tag == tg {
				call.Line = expr.Line - 1
				break
			}
		}

		return nil
	} else {
		return err
	}
}

/*
  Time functions
*/

func sleep (ms *CXArgument) error {
	if err := checkType("sleep", "i32", ms); err == nil {
		var duration int32
		encoder.DeserializeRaw(*ms.Value, &duration)

		time.Sleep(time.Duration(duration) * time.Millisecond)

		return nil
	} else {
		return err
	}
}


/*
  Prolog functions
*/

func setClauses (clss *CXArgument, mod *CXModule) error {
	if err := checkType("setClauses", "str", clss); err == nil {
		clauses := string(*clss.Value)
		mod.AddClauses(clauses)

		return nil
	} else {
		return err
	}
}

func addObject (obj *CXArgument, mod *CXModule) error {
	if err := checkType("addObject", "str", obj); err == nil {
		mod.AddObject(MakeObject(string(*obj.Value)))

		return nil
	} else {
		return err
	}
}

func setQuery (qry *CXArgument, mod *CXModule) error {
	if err := checkType("setQuery", "str", qry); err == nil {
		query := string(*qry.Value)
		mod.AddQuery(query)

		return nil
	} else {
		return err
	}
}

func remObject (obj *CXArgument, mod *CXModule) error {
	if err := checkType("remObject", "str", obj); err == nil {
		object := string(*obj.Value)
		mod.RemoveObject(object)

		return nil
	} else {
		return err
	}
}

func remObjects (mod *CXModule) error {
	mod.RemoveObjects()

	return nil
}

/*
  Meta-programming functions
*/

func remArg (tag *CXArgument, caller *CXFunction) error {
	if err := checkType("remArg", "str", tag); err == nil {
		for _, expr := range caller.Expressions {
			if expr.Tag == string(*tag.Value) {
				expr.RemoveArgument()
				return nil
			}
		}
	} else {
		return err
	}
	return errors.New(fmt.Sprintf("remArg: no expression with tag '%s' was found", string(*tag.Value)))
}

// func addArg (tag *CXArgument, ident *CXArgument, caller *CXFunction) error {
// 	if err := checkTwoTypes("addArg", "str", "str", tag, ident); err == nil {
// 		for _, expr := range caller.Expressions {
// 			if expr.Tag == string(*tag.Value) {
// 				expr.AddArgument(MakeArgument(ident.Value, "ident"))
// 				val := encoder.Serialize(int32(0))
// 				return MakeArgument(&val, "bool"), nil
// 			}
// 		}
// 	} else {
// 		return err
// 	}
// 	return errors.New(fmt.Sprintf("remArg: no expression with tag '%s' was found", string(*tag.Value)))
// }

func addExpr (tag *CXArgument, fnName *CXArgument, caller *CXFunction, line int) error {
	if err := checkType("addExpr", "str", fnName); err == nil {
		mod := caller.Module
		
		opName := string(*fnName.Value)
		if fn, err := mod.Context.GetFunction(opName, mod.Name); err == nil {
			expr := MakeExpression(fn)
			expr.AddTag(string(*tag.Value))
			
			caller.AddExpression(expr)

			//caller.Expressions = append(caller.Expressions[:line], expr, caller.Expressions[line:(len(caller.Expressions)-2)]...)

			// re-indexing expression line numbers
			// for i, expr := range caller.Expressions {
			// 	expr.Line = i
			// }
			
			//val := encoder.Serialize(int32(0))
			//return MakeArgument(&val, "bool"), nil
			return nil
		} else {
			//val := encoder.Serialize(int32(1))
			//return MakeArgument(&val, "bool"), nil
			return nil
		}
	} else {
		return err
	}
}

func remExpr (tag *CXArgument, caller *CXFunction) error {
	if err := checkType("remExpr", "str", tag); err == nil {
		for i, expr := range caller.Expressions {
			if expr.Tag == string(*tag.Value) {
				caller.RemoveExpression(i)
				//val := encoder.Serialize(int32(0))
				//return MakeArgument(&val, "bool"), nil
				return nil
			}
		}
	} else {
		return err
	}
	return errors.New(fmt.Sprintf("remExpr: no expression with tag '%s' was found", string(*tag.Value)))
}

//func affFn (filter *CXArgument, )

func affExpr (tag *CXArgument, filter *CXArgument, idx *CXArgument, caller *CXFunction, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("affExpr", "str", "str", "i32", tag, filter, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		if index == -1 {
			for _, ex := range caller.Expressions {
				if ex.Tag == string(*tag.Value) {
					affs := FilterAffordances(ex.GetAffordances(), string(*filter.Value))
					PrintAffordances(affs)
					val := encoder.Serialize(int32(len(affs)))

					for _, def := range call.State {
						if def.Name == expr.OutputNames[0].Name {
							def.Value = &val
							return nil
						}
					}

					call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, "i32"))
					return nil
				}
			}
		} else if index < -1 {
			for _, ex := range caller.Expressions {
				if ex.Tag == string(*tag.Value) {
					affs := FilterAffordances(ex.GetAffordances(), string(*filter.Value))
					val := encoder.Serialize(int32(len(affs)))

					for _, def := range call.State {
						if def.Name == expr.OutputNames[0].Name {
							def.Value = &val
							return nil
						}
					}
					
					call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, "i32"))
					return nil
				}
			}
		} else {
			for _, ex := range caller.Expressions {
				if ex.Tag == string(*tag.Value) {
					affs := FilterAffordances(ex.GetAffordances(), string(*filter.Value))
					affs[index].ApplyAffordance()
					val := encoder.Serialize(int32(len(affs)))

					for _, def := range call.State {
						if def.Name == expr.OutputNames[0].Name {
							def.Value = &val
							return nil
						}
					}

					if len(expr.OutputNames) > 0 {
						call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, "i32"))
					}
					return nil
				}
			}
		}
	} else {
		return err
	}
	return errors.New(fmt.Sprintf("affExpr: no expression with tag '%s' was found", string(*tag.Value)))
}

func initDef (arg1 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("initDef", "str", arg1); err == nil {
		typName := string(*arg1.Value)

		var zeroVal []byte
		switch  typName {
		case "bool": zeroVal = encoder.Serialize(int32(0))
		case "byte": zeroVal = []byte{byte(0)}
		case "i32": zeroVal = encoder.Serialize(int32(0))
		case "i64": zeroVal = encoder.Serialize(int64(0))
		case "f32": zeroVal = encoder.Serialize(float32(0))
		case "f64": zeroVal = encoder.Serialize(float64(0))
		case "[]bool": zeroVal = encoder.Serialize([]int32{0})
		case "[]byte": zeroVal = []byte{byte(0)}
		case "[]i32": zeroVal = encoder.Serialize([]int32{0})
		case "[]i64": zeroVal = encoder.Serialize([]int64{0})
		case "[]f32": zeroVal = encoder.Serialize([]float32{0})
		case "[]f64": zeroVal = encoder.Serialize([]float64{0})
		}

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &zeroVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &zeroVal, typName))
		return nil
	} else {
		return err
	}
}

/*
  Make Array
*/

func makeArray (typ string, size *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("makeArray", "i32", size); err == nil {
		var len int32
		encoder.DeserializeRaw(*size.Value, &len)

		switch typ {
		case "[]bool":
			arr := make([]bool, len)
			val := encoder.Serialize(arr)

			for _, def := range call.State {
				if def.Name == expr.OutputNames[0].Name {
					def.Value = &val
					return nil
				}
			}
			
			call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, typ))
			return nil
		case "[]byte":
			arr := make([]byte, len)

			for _, def := range call.State {
				if def.Name == expr.OutputNames[0].Name {
					def.Value = &arr
					return nil
				}
			}
			
			call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &arr, typ))
			return nil
		case "[]i32":
			arr := make([]int32, len)
			val := encoder.Serialize(arr)

			for _, def := range call.State {
				if def.Name == expr.OutputNames[0].Name {
					def.Value = &val
					return nil
				}
			}
			
			call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, typ))
			return nil
		case "[]i64":
			arr := make([]int64, len)
			val := encoder.Serialize(arr)

			for _, def := range call.State {
				if def.Name == expr.OutputNames[0].Name {
					def.Value = &val
					return nil
				}
			}
			
			call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, typ))
			return nil
		case "[]f32":
			arr := make([]float32, len)
			val := encoder.Serialize(arr)

			for _, def := range call.State {
				if def.Name == expr.OutputNames[0].Name {
					def.Value = &val
					return nil
				}
			}
			
			call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, typ))
			return nil
		case "[]f64":
			arr := make([]float64, len)
			val := encoder.Serialize(arr)

			for _, def := range call.State {
				if def.Name == expr.OutputNames[0].Name {
					def.Value = &val
					return nil
				}
			}
			
			call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &val, typ))
			return nil
		case "default":
			return errors.New(fmt.Sprintf("makeArray: argument 1 is type '%s'; expected type 'i32'", size.Typ))
		}
		return nil
	} else {
		return err
	}
}

func serialize_program (expr *CXExpression, call *CXCall) error {
	val := Serialize(call.Context)
	
	for _, def := range call.State {
		if def.Name == expr.OutputNames[0].Name {
			def.Value = val
			return nil
		}
	}
	
	call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, val, "[]byte"))
	return nil
}
