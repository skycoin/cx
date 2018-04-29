package base

import (
	"time"
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

		assignOutput(0, val, "bool", expr, call)
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

		assignOutput(0, val, "bool", expr, call)
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

		assignOutput(0, val, "bool", expr, call)
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
			if expr.Label == tg {
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

func identity (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	var name string
	encoder.DeserializeRaw(*arg.Value, &name)

	if arg, err := resolveIdent(name, call); err == nil {
		assignOutput(0, *arg.Value, arg.Typ, expr, call)
		return nil
	} else {
		return err
	}
}

func initDef (arg1 *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("initDef", "str", arg1); err == nil {
		var typName string
		encoder.DeserializeRaw(*arg1.Value, &typName)

		var zeroVal []byte
		if IsBasicType(typName) {
			zeroVal = *MakeDefaultValue(typName)
		} else {
			call.Program.SelectModule(expr.Package.Name)
			if byts, err := ResolveStruct(typName, call.Program); err == nil {
				zeroVal = byts
			} else {
				return err
			}
		}

		assignOutput(0, zeroVal, typName, expr, call)
		return nil
	} else {
		return err
	}
}

func serialize_program (expr *CXExpression, call *CXCall) error {
	// val := Serialize(call.Program)

	// remove this once Serialize is fixed
	val := []byte{0}

	assignOutput(0, val, "[]byte", expr, call)
	return nil
}

// multi dimensional arrays functions

func mdim_append (arr, elt *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("mdim.append", "str", "str", arr, elt); err == nil {
		// we receive the identifiers of both variables
		var _arr string
		var _elt string

		encoder.DeserializeRaw(*arr.Value, &_arr)
		encoder.DeserializeRaw(*elt.Value, &_elt)
		
		if rArr, err := resolveIdent(_arr, call); err == nil {
			if rElt, err := resolveIdent(_elt, call); err == nil {
				// checking that the second argument's type is similar to the first argument's type
				if rArr.Typ[2:] == rElt.Typ {
					var arrSize int32
					encoder.DeserializeRaw((*rArr.Value)[:4], &arrSize)

					firstChunk := make([]byte, len((*rArr.Value)[4:]))
					secondChunk := make([]byte, len(*rElt.Value))

					copy(firstChunk, (*rArr.Value)[4:])
					copy(secondChunk, *rElt.Value)

					final := append(encoder.Serialize(arrSize + 1), firstChunk...)
					final = append(final, secondChunk...)

					assignOutput(0, final, rArr.Typ, expr, call)
					return nil
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

func mdim_read (arr, index *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("mdim.read", "str", "i32", arr, index); err == nil {
		var _arr string
		var _index int32

		encoder.DeserializeRaw(*arr.Value, &_arr)
		encoder.DeserializeAtomic(*index.Value, &_index)

		if rArr, err := resolveIdent(_arr, call); err == nil {
			if array, err, _, _ := GetArrayFromArray(*rArr.Value, rArr.Typ, _index); err == nil {
				assignOutput(0, array, rArr.Typ[2:], expr, call)
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

func mdim_write (arr, index, instance *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkThreeTypes("mdim.write", "str", "i32", "str", arr, index, instance); err == nil {
		var _arr string
		var _index int32
		var _instance string

		encoder.DeserializeRaw(*arr.Value, &_arr)
		encoder.DeserializeAtomic(*index.Value, &_index)
		encoder.DeserializeRaw(*instance.Value, &_instance)

		if rArr, err := resolveIdent(_arr, call); err == nil {
			if rInst, err := resolveIdent(_instance, call); err == nil {
				if _, err, offset, size := GetArrayFromArray(*rArr.Value, rArr.Typ, _index); err == nil {
					firstChunk := make([]byte, offset)
					secondChunk := make([]byte, len(*rArr.Value) - int((offset + size)))

					copy(firstChunk, (*rArr.Value)[:offset])
					copy(secondChunk, (*rArr.Value)[offset+size:])

					final := append(firstChunk, *rInst.Value...)
					final = append(final, secondChunk...)

					assignOutput(0, final, rArr.Typ, expr, call)
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

func mdim_len (arr *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("mdim.len", "str", arr); err == nil {
		var _arr string
		var len int32

		encoder.DeserializeRaw(*arr.Value, &_arr)
		
		if rArr, err := resolveIdent(_arr, call); err == nil {
			encoder.DeserializeAtomic((*rArr.Value)[:4], &len)
			output := encoder.Serialize(len)
			assignOutput(0, output, "i32", expr, call)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func mdim_make (length, typ *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("mdim.make", "i32", "str", length, typ); err == nil {
		var _len int32
		var _typ string

		encoder.DeserializeRaw(*length.Value, &_len)
		encoder.DeserializeRaw(*typ.Value, &_typ)

		if _len == 0 {
			output := []byte{0, 0, 0, 0}
			assignOutput(0, output, _typ, expr, call)
			return nil
		}
		
		var typSize int
		switch _typ[len(_typ)-4:] {
		case "]i64", "]f64":
			typSize = 8
		case "bool", "]i32", "]f32":
			typSize = 4
		case "byte", "]str":
			typSize = 1
		}

		instances := append(*length.Value, make([]byte, int(_len) * typSize)...)
		
		assignOutput(0, instances, _typ, expr, call)
		return nil
	} else {
		return err
	}
}

// custom type arrays functions

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
					
					assignOutput(0, final, rArr.Typ, expr, call)
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
				assignOutput(0, instance, rArr.Typ[2:], expr, call)
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
					// finalSize := int(offset) + (len(*rArr.Value) - int((offset + size))) + len(*rInst.Value)
					// final := make([]byte, finalSize, finalSize)


					// firstChunkSize := int(offset)
					// newValSize := len(*rInst.Value)

					// for c := 0; c < len(final); c++ {
					// 	if c < firstChunkSize {
					// 		final[c] = (*rArr.Value)[c]
					// 	} else if c >= firstChunkSize && c < firstChunkSize + newValSize {
					// 		final[c] = (*rInst.Value)[c - (firstChunkSize)]
					// 	} else if c >= firstChunkSize + newValSize {
					// 		final[c] = (*rArr.Value)[c - (firstChunkSize + newValSize)]
					// 	}
					// }
					

					firstChunk := make([]byte, offset)
					secondChunk := make([]byte, len(*rArr.Value) - int((offset + size)))

					copy(firstChunk, (*rArr.Value)[:offset])
					copy(secondChunk, (*rArr.Value)[offset+size:])



					// firstChunk := (*rArr.Value)[:offset]
					// secondChunk := (*rArr.Value)[offset+size:]

					final := append(firstChunk, *rInst.Value...)
					final = append(final, secondChunk...)

					// final := append((*rArr.Value)[:offset], *rInst.Value...)
					// final = append(final, (*rArr.Value)[offset+size:]...)

					//*rArr.Value = final
					assignOutput(0, final, rArr.Typ, expr, call)
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
			assignOutput(0, output, "i32", expr, call)
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
			assignOutput(0, output, _typ, expr, call)
			return nil
		}

		var instances []byte
		if oneInst, err := ResolveStruct(_typ[2:], call.Program); err == nil {
			for c := int32(0); c < _len; c++ {
				//var another []byte
				another := make([]byte, len(oneInst))
				copy(another, oneInst)
				instances = append(instances, another...)
			}
		}

		instances = append(*length.Value, instances...)
		
		assignOutput(0, instances, _typ, expr, call)
		return nil
	} else {
		return err
	}
}

func cstm_serialize (instance *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkType("cstm.serialize", "str", instance); err == nil {
		var _instance string
		encoder.DeserializeRaw(*instance.Value, &_instance)
		
		if rInst, err := resolveIdent(_instance, call); err == nil {
			sInst := encoder.Serialize(*rInst.Value)
			assignOutput(0, sInst, "[]byte", expr, call)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func cstm_deserialize (byts, typ *CXArgument, expr *CXExpression, call *CXCall) error {
	if err := checkTwoTypes("cstm.deserialize", "str", "str", byts, typ); err == nil {
		var _byts string
		var _typ string
		encoder.DeserializeRaw(*byts.Value, &_byts)
		encoder.DeserializeRaw(*typ.Value, &_typ)
		
		if rByts, err := resolveIdent(_byts, call); err == nil {
			var dsStrct []byte
			encoder.DeserializeRaw(*rByts.Value, &dsStrct)

			assignOutput(0, dsStrct, _typ, expr, call)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}
