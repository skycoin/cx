package base

import (
	"fmt"
	"time"
	"errors"
	"strings"
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

		assignOutput(&val, "bool", expr, call)
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

		assignOutput(&val, "bool", expr, call)
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

		assignOutput(&val, "bool", expr, call)
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
			return nil
		} else {
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
				return nil
			}
		}
	} else {
		return err
	}
	return errors.New(fmt.Sprintf("remExpr: no expression with tag '%s' was found", string(*tag.Value)))
}

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

					assignOutput(&val, "i32", expr, call)
					return nil
				}
			}
		} else if index < -1 {
			for _, ex := range caller.Expressions {
				if ex.Tag == string(*tag.Value) {
					affs := FilterAffordances(ex.GetAffordances(), string(*filter.Value))
					val := encoder.Serialize(int32(len(affs)))

					assignOutput(&val, "i32", expr, call)
					return nil
				}
			}
		} else {
			for _, ex := range caller.Expressions {
				if ex.Tag == string(*tag.Value) {
					affs := FilterAffordances(ex.GetAffordances(), string(*filter.Value))
					affs[index].ApplyAffordance()
					val := encoder.Serialize(int32(len(affs)))

					if len(expr.OutputNames) > 0 {
						assignOutput(&val, "i32", expr, call)
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

func ResolveStruct (typ string, cxt *CXProgram) ([]byte, error) {
	var bs []byte

	found := false
	if mod, err := cxt.GetCurrentModule(); err == nil {
		var foundStrct *CXStruct
		for _, strct := range mod.Structs {
			if strct.Name == typ {
				found = true
				foundStrct = strct
				break
			}
		}
		if !found {
			for _, imp := range mod.Imports {
				for _, strct := range imp.Structs {
					if strct.Name == typ {
						found = true
						foundStrct = strct
						break
					}
				}
			}
		}

		if !found {
			return nil, errors.New(fmt.Sprintf("type '%s' not defined\n", typ))
		}
		
		for _, fld := range foundStrct.Fields {
			isBasic := false
			for _, basic := range BASIC_TYPES {
				if fld.Typ == basic {
					isBasic = true
					bs = append(bs, *MakeDefaultValue(basic)...)
					break
				}
			}

			if !isBasic {
				if byts, err := ResolveStruct(fld.Typ, cxt); err == nil {
					bs = append(bs, byts...)
				} else {
					return nil, err
				}
			}
		}
	}
	return bs, nil
}

func identity (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	found := false
	name := string(*arg.Value)
	for _, def := range call.State {
		if def.Name == name {
			found = true
			assignOutput(def.Value, def.Typ, expr, call)
			return nil
		}
	}
	if !found {
		// then it can be a global
		identParts := strings.Split(name, ".")
		if def, err := expr.Module.GetDefinition(identParts[0]); err == nil {
			
			if len(identParts) > 1 {
				if strct, err := expr.Context.GetStruct(def.Typ, expr.Module.Name); err == nil {
					byts, typ := resolveStructField(identParts[1], def.Value, strct)
					
					assignOutput(&byts, typ, expr, call)
					return nil
				}
			} else {
				assignOutput(def.Value, def.Typ, expr, call)
			}
		}
	}
	return errors.New(fmt.Sprintf("identity: identifier '%s' not found", string(*arg.Value)))
}

func initDef (arg1 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("initDef", "str", arg1); err == nil {
		typName := string(*arg1.Value)

		isBasic := false
		for _, basic := range BASIC_TYPES {
			if basic == typName {
				isBasic = true
				break
			}
		}

		var zeroVal []byte
		if isBasic {
			zeroVal = *MakeDefaultValue(typName)
		} else {
			if byts, err := ResolveStruct(typName, call.Context); err == nil {
				zeroVal = byts
			} else {
				return err
			}
		}

		assignOutput(&zeroVal, typName, expr, call)
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

			assignOutput(&val, typ, expr, call)
			return nil
		case "[]byte":
			arr := make([]byte, len)

			assignOutput(&arr, typ, expr, call)
			return nil
		case "[]i32":
			arr := make([]int32, len)
			val := encoder.Serialize(arr)

			assignOutput(&val, typ, expr, call)
			return nil
		case "[]i64":
			arr := make([]int64, len)
			val := encoder.Serialize(arr)

			assignOutput(&val, typ, expr, call)
			return nil
		case "[]f32":
			arr := make([]float32, len)
			val := encoder.Serialize(arr)

			assignOutput(&val, typ, expr, call)
			return nil
		case "[]f64":
			arr := make([]float64, len)
			val := encoder.Serialize(arr)

			assignOutput(&val, typ, expr, call)
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

	assignOutput(val, "[]byte", expr, call)
	return nil
}
