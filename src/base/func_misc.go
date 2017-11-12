package base

import (
	"fmt"
	"time"
	"errors"
	//"strings"
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
		var tg string
		encoder.DeserializeRaw(*tag.Value, &tg)

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
		var clauses string
		encoder.DeserializeRaw(*clss.Value, &clauses)
		
		mod.AddClauses(clauses)

		return nil
	} else {
		return err
	}
}

func addObject (obj *CXArgument, mod *CXModule) error {
	if err := checkType("addObject", "str", obj); err == nil {
		var object string
		encoder.DeserializeRaw(*obj.Value, &object)
		mod.AddObject(MakeObject(object))

		return nil
	} else {
		return err
	}
}

func setQuery (qry *CXArgument, mod *CXModule) error {
	if err := checkType("setQuery", "str", qry); err == nil {
		var query string
		encoder.DeserializeRaw(*qry.Value, &query)
		mod.AddQuery(query)

		return nil
	} else {
		return err
	}
}

func remObject (obj *CXArgument, mod *CXModule) error {
	if err := checkType("remObject", "str", obj); err == nil {
		var object string
		encoder.DeserializeRaw(*obj.Value, &object)
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
	var tg string
	encoder.DeserializeRaw(*tag.Value, &tg)
	if err := checkType("remArg", "str", tag); err == nil {
		for _, expr := range caller.Expressions {
			if expr.Tag == tg {
				expr.RemoveArgument()
				return nil
			}
		}
	} else {
		return err
	}
	return errors.New(fmt.Sprintf("remArg: no expression with tag '%s' was found", tg))
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

		var opName string
		var tg string
		encoder.DeserializeRaw(*fnName.Value, &opName)
		encoder.DeserializeRaw(*tag.Value, &tg)
		
		if fn, err := mod.Context.GetFunction(opName, mod.Name); err == nil {
			expr := MakeExpression(fn)
			expr.AddTag(tg)
			
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
	var tg string
	encoder.DeserializeRaw(*tag.Value, &tg)
	if err := checkType("remExpr", "str", tag); err == nil {
		for i, expr := range caller.Expressions {
			if expr.Tag == tg {
				caller.RemoveExpression(i)
				return nil
			}
		}
	} else {
		return err
	}
	return errors.New(fmt.Sprintf("remExpr: no expression with tag '%s' was found", tg))
}

func affExpr (tag *CXArgument, filter *CXArgument, idx *CXArgument, caller *CXFunction, expr *CXExpression, call *CXCall) error {
	var tg string
	var _filter string
	encoder.DeserializeRaw(*tag.Value, &tg)
	encoder.DeserializeRaw(*filter.Value, &_filter)
	
	if err := checkThreeTypes("affExpr", "str", "str", "i32", tag, filter, idx); err == nil {
		var index int32
		encoder.DeserializeRaw(*idx.Value, &index)

		if index == -1 {
			for _, ex := range caller.Expressions {
				if ex.Tag == tg {
					affs := FilterAffordances(ex.GetAffordances(), _filter)
					PrintAffordances(affs)
					val := encoder.Serialize(int32(len(affs)))

					assignOutput(&val, "i32", expr, call)
					return nil
				}
			}
		} else if index < -1 {
			for _, ex := range caller.Expressions {
				if ex.Tag == tg {
					affs := FilterAffordances(ex.GetAffordances(), _filter)
					val := encoder.Serialize(int32(len(affs)))

					assignOutput(&val, "i32", expr, call)
					return nil
				}
			}
		} else {
			for _, ex := range caller.Expressions {
				if ex.Tag == tg {
					affs := FilterAffordances(ex.GetAffordances(), _filter)
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
	return errors.New(fmt.Sprintf("affExpr: no expression with tag '%s' was found", tg))
}

func ResolveStruct (typ string, cxt *CXProgram) ([]byte, error) {
	var bs []byte

	found := false
	if mod, err := cxt.GetCurrentModule(); err == nil {
		var foundStrct *CXStruct

		if typ[:2] == "[]" {
			// empty serialized struct array
			return []byte{0, 0, 0, 0}, nil
		}
		
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
	//found := false
	var name string
	encoder.DeserializeRaw(*arg.Value, &name)
	
	//var foundDef *CXDefinition

	//fmt.Println(name)

	if arg, err := resolveIdent(name, call); err == nil {
		assignOutput(arg.Value, arg.Typ, expr, call)
		return nil
	} else {
		return err
	}

	// for _, def := range call.State {
	// 	if def.Name == name {
	// 		found = true
	// 		foundDef = def
	// 		//return nil // we want to grab the last instance of that variable
	// 	}
	// }

	// if found && foundDef != nil {
	// 	assignOutput(foundDef.Value, foundDef.Typ, expr, call)
	// 	return nil
	// }
	
	// if !found {
	// 	// then it can be a global
	// 	identParts := strings.Split(name, ".")
	// 	if def, err := expr.Module.GetDefinition(identParts[0]); err == nil {
			
	// 		if len(identParts) > 1 {
	// 			if strct, err := expr.Context.GetStruct(def.Typ, expr.Module.Name); err == nil {
	// 				byts, typ, _, _ := resolveStructField(identParts[1], def.Value, strct)
					
	// 				assignOutput(&byts, typ, expr, call)
	// 				return nil
	// 			}
	// 		} else {
	// 			assignOutput(def.Value, def.Typ, expr, call)
	// 		}
	// 	}
	// }
	// return errors.New(fmt.Sprintf("identity: identifier '%s' not found", name))
}

func initDef (arg1 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("initDef", "str", arg1); err == nil {
		var typName string
		encoder.DeserializeRaw(*arg1.Value, &typName)

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
		var _len int32
		encoder.DeserializeRaw(*size.Value, &_len)

		switch typ {
		case "[]bool":
			arr := make([]int32, _len)
			val := encoder.Serialize(arr)

			assignOutput(&val, typ, expr, call)
			return nil
		case "[]byte":
			arr := make([]byte, _len)
			val := encoder.Serialize(arr)

			assignOutput(&val, typ, expr, call)
			return nil
		case "[]str":
			arr := make([]string, _len)
			val := encoder.Serialize(arr)

			assignOutput(&val, typ, expr, call)
			return nil
		case "[]i32":
			
			arr := make([]int32, _len)
			val := encoder.Serialize(arr)
			//fmt.Println("hohoho", len(val))
			
			assignOutput(&val, typ, expr, call)
			return nil
		case "[]i64":
			arr := make([]int64, _len)
			val := encoder.Serialize(arr)

			assignOutput(&val, typ, expr, call)
			return nil
		case "[]f32":
			arr := make([]float32, _len)
			val := encoder.Serialize(arr)

			assignOutput(&val, typ, expr, call)
			return nil
		case "[]f64":
			arr := make([]float64, _len)
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

// custom type arrays functions

func getStrctFromArray (arr *CXArgument, index int32, expr *CXExpression, call *CXCall) ([]byte, error, int32, int32) {
	var arrSize int32
	encoder.DeserializeAtomic((*arr.Value)[:4], &arrSize)

	if index < 0 {
		return nil, errors.New(fmt.Sprintf("%s.read: negative index %d", arr.Typ, index)), 0, 0
	}

	if index >= arrSize {
		return nil, errors.New(fmt.Sprintf("%s.read: index %d exceeds array of length %d", arr.Typ, index, arrSize)), 0, 0
	}

	if strct, err := call.Context.GetStruct(arr.Typ[2:], expr.Module.Name); err == nil {
		instances := (*arr.Value)[4:]
		lastFld := strct.Fields[len(strct.Fields) - 1]
		
		var lowerBound int32
		var upperBound int32
		
		var size int32
		encoder.DeserializeAtomic((*arr.Value)[:4], &size)

		// in here we can use <=. we can't do this in resolveStrctField
		for c := int32(0); c <= index; c++ {
			subArray := instances[upperBound:]
			_, _, off, size := resolveStructField(lastFld.Name, &subArray, strct)

			lowerBound = upperBound
			upperBound = upperBound + off + size
		}

		output := instances[lowerBound:upperBound]
		return output, nil, lowerBound + 4, upperBound - lowerBound
	} else {
		return nil, err, 0, 0
	}
}

func cstm_append (arr, strctInst *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("cstm.append", "str", "str", arr, strctInst); err == nil {
		// we receive the identifiers of both variables
		var _arr string
		var _strctInst string

		encoder.DeserializeRaw(*arr.Value, &_arr)
		encoder.DeserializeRaw(*strctInst.Value, &_strctInst)
		
		if rArr, err := resolveIdent(_arr, call); err == nil {
			if rStrctInst, err := resolveIdent(_strctInst, call); err == nil {
				// checking that the second argument's type is similar to the first argument's type
				if rArr.Typ[2:] == rStrctInst.Typ {
					var arrSize int32
					encoder.DeserializeRaw((*rArr.Value)[:4], &arrSize)

					// output := append(encoder.Serialize(arrSize + 1), (*rArr.Value)[4:]...)
					// *rArr.Value = append(output, *rStrctInst.Value...)

					firstChunk := make([]byte, len((*rArr.Value)[4:]))
					secondChunk := make([]byte, len(*rStrctInst.Value))

					copy(firstChunk, (*rArr.Value)[4:])
					copy(secondChunk, *rStrctInst.Value)

					final := append(encoder.Serialize(arrSize + 1), firstChunk...)
					final = append(final, secondChunk...)

					//fmt.Println(final)
					
					assignOutput(&final, rArr.Typ, expr, call)
					//*rArr.Value = final
					
					//fmt.Println("arrSize", arrSize, rStrctInst.Typ, len(output))
				}
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
	return nil
}

func cstm_read (arr, index *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("cstm.read", "str", "i32", arr, index); err == nil {
		var _arr string
		var _index int32

		encoder.DeserializeRaw(*arr.Value, &_arr)
		encoder.DeserializeAtomic(*index.Value, &_index)

		if rArr, err := resolveIdent(_arr, call); err == nil {
			if instance, err, _, _ := getStrctFromArray(rArr, _index, expr, call); err == nil {
				output := make([]byte, len(instance))
				copy(output, instance)
				assignOutput(&output, rArr.Typ[2:], expr, call)
			} else {
				return err
			}
			
		} else {
			return err
		}
	} else {
		return err
	}
	return nil
}

func cstm_write (arr, index, instance *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("cstm.write", "str", "i32", "str", arr, index, instance); err == nil {
		var _arr string
		var _index int32
		var _instance string

		encoder.DeserializeRaw(*arr.Value, &_arr)
		encoder.DeserializeAtomic(*index.Value, &_index)
		encoder.DeserializeRaw(*instance.Value, &_instance)

		if rArr, err := resolveIdent(_arr, call); err == nil {
			if rInst, err := resolveIdent(_instance, call); err == nil {
				if _, err, offset, size := getStrctFromArray(rArr, _index, expr, call); err == nil {
					firstChunk := make([]byte, offset)
					secondChunk := make([]byte, len(*rArr.Value) - int((offset + size)))

					copy(firstChunk, (*rArr.Value)[:offset])
					copy(secondChunk, (*rArr.Value)[offset+size:])

					final := append(firstChunk, *rInst.Value...)
					final = append(final, secondChunk...)
					
					// final := append((*rArr.Value)[:offset], *rInst.Value...)
					// final = append(final, (*rArr.Value)[offset+size:]...)

					*rArr.Value = final
					//assignOutput(&final, rArr.Typ, expr, call)
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
		return err
	}
	return nil
}

func cstm_len (arr *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("cstm.len", "str", arr); err == nil {
		var _arr string
		var len int32

		encoder.DeserializeRaw(*arr.Value, &_arr)
		
		if rArr, err := resolveIdent(_arr, call); err == nil {
			encoder.DeserializeAtomic((*rArr.Value)[:4], &len)
			output := encoder.Serialize(len)
			assignOutput(&output, "i32", expr, call)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func cstm_make (length, typ *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("cstm.make", "i32", "str", length, typ); err == nil {
		var _len int32
		var _typ string

		encoder.DeserializeRaw(*length.Value, &_len)
		encoder.DeserializeRaw(*typ.Value, &_typ)

		if _len == 0 {
			output := []byte{0, 0, 0, 0}
			assignOutput(&output, _typ, expr, call)
			return nil
		}

		var instances []byte
		if oneInst, err := ResolveStruct(_typ[2:], call.Context); err == nil {
			for c := int32(0); c < _len; c++ {
				//var another []byte
				another := make([]byte, len(oneInst))
				copy(another, oneInst)
				instances = append(instances, another...)
			}
		}

		instances = append(*length.Value, instances...)
		
		assignOutput(&instances, _typ, expr, call)
		return nil
	} else {
		return err
	}
}

// test functions

func test_error (message *CXArgument, isErrorPresent bool, expr *CXExpression) error {
	if !isErrorPresent {
		var _message string
		encoder.DeserializeRaw(*message.Value, &_message)
		if _message == "" {
			fmt.Println(fmt.Sprintf("%d: an error was expected and did not occur", expr.FileLine))
		} else {
			fmt.Println(fmt.Sprintf("%d: an error was expected and did not occur; %s", expr.FileLine, _message))
		}
		
		return nil
	} else {
		return nil
	}
}

func test_value (result *CXArgument, expected *CXArgument, message *CXArgument, expr *CXExpression) error {
	if result.Typ != expected.Typ {
		fmt.Println(fmt.Sprintf("%d: result and expected value are not of the same type", expr.FileLine))
		return nil
	}
	
	equal := true
	var _message string
	encoder.DeserializeRaw(*message.Value, &_message)

	if len(*result.Value) != len(*expected.Value) {
		equal = false
	}

	// fmt.Println("")
	// fmt.Println(*result.Value)
	// fmt.Println(*expected.Value)

	if equal {
		for i, byt := range *result.Value {
			if byt != (*expected.Value)[i] {
				equal = false
				break
			}
		}
	}

	if !equal {
		if _message == "" {
			fmt.Println(fmt.Sprintf("%d: result was not equal to the expected value", expr.FileLine))
		} else {
			fmt.Println(fmt.Sprintf("%d: result was not equal to the expected value; %s", expr.FileLine, _message))
		}
	}
	
	return nil
}
