package base

import (
	"fmt"
	"math/rand"
	"time"
	"bytes"
	"strings"
	"strconv"
	"errors"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func assignOutput (output *[]byte, typ string, expr *CXExpression, call *CXCall) {
	outName := expr.OutputNames[0].Name
	expr.OutputNames[0].Typ = typ

	for _, char := range outName {
		if char == '.' {
			identParts := strings.Split(outName, ".")

			if def, err := expr.Module.GetDefinition(identParts[0]); err == nil {
				if strct, err := call.Context.GetStruct(def.Typ, expr.Module.Name); err == nil {
					_, _, offset, size := resolveStructField(identParts[1], def.Value, strct)

					firstChunk := make([]byte, offset)
					secondChunk := make([]byte, len(*def.Value) - int(offset + size))

					copy(firstChunk, (*def.Value)[:offset])
					copy(secondChunk, (*def.Value)[offset+size:])

					final := append(firstChunk, *output...)
					final = append(final, secondChunk...)
					
					*def.Value = final
					return
				}
			}

			for _, def := range call.State {
				if def.Name == identParts[0] {
					if strct, err := call.Context.GetStruct(def.Typ, expr.Module.Name); err == nil {
						byts, typ, offset, size := resolveStructField(identParts[1], def.Value, strct)

						isBasic := false
						for _, basic := range BASIC_TYPES {
							if basic == typ {
								isBasic = true
								break
							}
						}
						
						if true || typ == "str" || typ == "[]str" || typ == "[]bool" ||
							typ == "[]byte" || typ == "[]i32" ||
							typ == "[]i64" || typ == "[]f32" || typ == "[]f64" || !isBasic {

							firstChunk := make([]byte, offset)
							secondChunk := make([]byte, len(*def.Value) - int(offset + size))

							copy(firstChunk, (*def.Value)[:offset])
							copy(secondChunk, (*def.Value)[offset+size:])

							final := append(firstChunk, *output...)
							final = append(final, secondChunk...)


							// final := append((*def.Value)[:offset], *output...)
							// final = append(final, (*def.Value)[offset+size:]...)

							// if identParts[1] == "appearance" || identParts[1] == "displacement" {
							// 	fmt.Println()
							// 	fmt.Println(outName)
							// 	fmt.Println("def.Value", len(*def.Value))
							// 	fmt.Println("final", len(final))
							// }
							
							*def.Value = final
							//fmt.Println("after", def.Name, len(*def.Value))
							return
						} else {
							
							for c := 0; c < len(byts); c++ {
								byts[c] = (*output)[c]
							}
						}
					} else {
						panic(err)
					}
					return
				}
			}
			break
		}
		if char == '[' {
			identParts := strings.Split(outName, "[")

			for _, def := range call.State {
				if def.Name == identParts[0] {
					idx, _ := strconv.ParseInt(identParts[1], 10, 64)
					for c := 0; c < len(*output); c++ {
						if typ == "i64" || typ == "f64" {
							(*def.Value)[(int(idx)*8) + 4 + c] = (*output)[c]
						} else if typ == "byte" {
							(*def.Value)[int(idx) + c] = (*output)[c]
						} else {
							(*def.Value)[(int(idx)*4) + 4 + c] = (*output)[c]
						}
					}
					return
				}
			}
			break
		}
	}

	if def, err := expr.Module.GetDefinition(outName); err == nil {
		def.Value = output
		return
	}

	for _, def := range call.State {
		if def.Name == outName {
			def.Value = output
			return
		}
	}

	call.State = append(call.State, MakeDefinition(outName, output, typ))
	return
}

func checkType (fnName string, typ string, arg *CXArgument) error {
	if arg.Typ != typ {
		return errors.New(fmt.Sprintf("%s: argument 1 is type '%s'; expected type '%s'", fnName, arg.Typ, typ))
	}
	return nil
}

func checkTwoTypes (fnName string, typ1 string, typ2 string, arg1 *CXArgument, arg2 *CXArgument) error {
	if arg1.Typ != typ1 || arg2.Typ != typ2 {
		if arg1.Typ != typ1 {
			return errors.New(fmt.Sprintf("%s: argument 1 is type '%s'; expected type '%s'", fnName, arg1.Typ, typ1))
		}
		return errors.New(fmt.Sprintf("%s: argument 2 is type '%s'; expected type '%s'", fnName, arg2.Typ, typ2))
	}
	return nil
}

func checkThreeTypes (fnName string, typ1 string, typ2 string, typ3 string, arg1 *CXArgument, arg2 *CXArgument, arg3 *CXArgument) error {
	if arg1.Typ != typ1 || arg2.Typ != typ2 || arg3.Typ != typ3 {
		if arg1.Typ != typ1 {
			return errors.New(fmt.Sprintf("%s: argument 1 is type '%s'; expected type '%s'", fnName, arg1.Typ, typ1))
		} else if arg2.Typ != typ2 {
			return errors.New(fmt.Sprintf("%s: argument 2 is type '%s'; expected type '%s'", fnName, arg2.Typ, typ2))
		}
		return errors.New(fmt.Sprintf("%s: argument 3 is type '%s'; expected type '%s'", fnName, arg3.Typ, typ3))
	}
	return nil
}

func checkFourTypes (fnName, typ1, typ2, typ3, typ4 string, arg1, arg2, arg3, arg4 *CXArgument) error {
	if arg1.Typ != typ1 || arg2.Typ != typ2 || arg3.Typ != typ3 || arg4.Typ != typ4 {
		if arg1.Typ != typ1 {
			return errors.New(fmt.Sprintf("%s: argument 1 is type '%s'; expected type '%s'", fnName, arg1.Typ, typ1))
		} else if arg2.Typ != typ2 {
			return errors.New(fmt.Sprintf("%s: argument 2 is type '%s'; expected type '%s'", fnName, arg2.Typ, typ2))
		} else if arg3.Typ != typ3 {
			return errors.New(fmt.Sprintf("%s: argument 3 is type '%s'; expected type '%s'", fnName, arg3.Typ, typ3))
		}
		return errors.New(fmt.Sprintf("%s: argument 4 is type '%s'; expected type '%s'", fnName, arg4.Typ, typ4))
	}
	return nil
}

func checkFiveTypes (fnName, typ1, typ2, typ3, typ4, typ5 string, arg1, arg2, arg3, arg4, arg5 *CXArgument) error {
	if arg1.Typ != typ1 || arg2.Typ != typ2 || arg3.Typ != typ3 || arg4.Typ != typ4 || arg5.Typ != typ5 {
		if arg1.Typ != typ1 {
			return errors.New(fmt.Sprintf("%s: argument 1 is type '%s'; expected type '%s'", fnName, arg1.Typ, typ1))
		} else if arg2.Typ != typ2 {
			return errors.New(fmt.Sprintf("%s: argument 2 is type '%s'; expected type '%s'", fnName, arg2.Typ, typ2))
		} else if arg3.Typ != typ3 {
			return errors.New(fmt.Sprintf("%s: argument 3 is type '%s'; expected type '%s'", fnName, arg3.Typ, typ3))
		} else if arg4.Typ != typ4 {
			return errors.New(fmt.Sprintf("%s: argument 4 is type '%s'; expected type '%s'", fnName, arg4.Typ, typ4))
		}
		return errors.New(fmt.Sprintf("%s: argument 5 is type '%s'; expected type '%s'", fnName, arg5.Typ, typ5))
	}
	return nil
}

func checkSixTypes (fnName, typ1, typ2, typ3, typ4, typ5, typ6 string, arg1, arg2, arg3, arg4, arg5, arg6 *CXArgument) error {
	if arg1.Typ != typ1 || arg2.Typ != typ2 || arg3.Typ != typ3 || arg4.Typ != typ4 || arg5.Typ != typ5 || arg6.Typ != typ6 {
		if arg1.Typ != typ1 {
			return errors.New(fmt.Sprintf("%s: argument 1 is type '%s'; expected type '%s'", fnName, arg1.Typ, typ1))
		} else if arg2.Typ != typ2 {
			return errors.New(fmt.Sprintf("%s: argument 2 is type '%s'; expected type '%s'", fnName, arg2.Typ, typ2))
		} else if arg3.Typ != typ3 {
			return errors.New(fmt.Sprintf("%s: argument 3 is type '%s'; expected type '%s'", fnName, arg3.Typ, typ3))
		} else if arg4.Typ != typ4 {
			return errors.New(fmt.Sprintf("%s: argument 4 is type '%s'; expected type '%s'", fnName, arg4.Typ, typ4))
		} else if arg5.Typ != typ5 {
			return errors.New(fmt.Sprintf("%s: argument 5 is type '%s'; expected type '%s'", fnName, arg5.Typ, typ5))
		}
		return errors.New(fmt.Sprintf("%s: argument 6 is type '%s'; expected type '%s'", fnName, arg6.Typ, typ6))
	}
	return nil
}

func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max - min) + min
}

func removeDuplicatesInt(elements []int) []int {
	// Use map to record duplicates as we find them.
	encountered := map[int]bool{}
	result := []int{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func removeDuplicates(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func concat (strs ...string) string {
	var buffer bytes.Buffer
	
	for i := 0; i < len(strs); i++ {
		buffer.WriteString(strs[i])
	}
	
	return buffer.String()
}



// (modules (main (definitions (str ))
//                (functions (= double (addF64 ((f64 n) (f64 n))))
//                           (= triple (subF64 ((f64 n) (f64 n))))))
// 	 (math (definitions)))

// (modules (module main
//                  (definitions
//                      (definition greeting str "hello")
//                      (definition ten i32 10)
//                    (definition ten i32 10)
//                      )))


func PrintValue (identName string, value *[]byte, typName string, cxt *CXProgram) string {
	var argName string
	switch typName {
	case "str":
		var val string
		encoder.DeserializeRaw(*value, &val)
		argName = fmt.Sprintf("\"%s\"", val)
	case "bool":
		var val int32
		encoder.DeserializeRaw(*value, &val)
		if val == 0 {
			argName = "false"
		} else {
			fmt.Printf("true")
			argName = "true"
		}
	case "byte":
		var val []byte
		encoder.DeserializeRaw(*value, &val)
		argName = fmt.Sprintf("%#v", val)
	case "i32":
		var val int32
		encoder.DeserializeRaw(*value, &val)
		argName = fmt.Sprintf("%#v", val)
	case "i64":
		var val int64
		encoder.DeserializeRaw(*value, &val)
		argName = fmt.Sprintf("%#v", val)
	case "f32":
		var val float32
		encoder.DeserializeRaw(*value, &val)
		argName = fmt.Sprintf("%#v", val)
	case "f64":
		var val float64
		encoder.DeserializeRaw(*value, &val)
		argName = fmt.Sprintf("%#v", val)
	case "[]byte":
		var val []byte
		encoder.DeserializeRaw(*value, &val)
		argName = fmt.Sprintf("%#v", val)
	case "[]str":
		var val []string
		encoder.DeserializeRaw(*value, &val)
		argName = fmt.Sprintf("%#v", val)
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*value, &val)
		argName = fmt.Sprintf("%#v", val)
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*value, &val)
		argName = fmt.Sprintf("%#v", val)
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*value, &val)
		argName = fmt.Sprintf("%#v", val)
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*value, &val)
		argName = fmt.Sprintf("%#v", val)
	default:
		// struct or custom type
		if mod, err := cxt.GetCurrentModule(); err == nil {
			if strct, err := cxt.GetStruct(typName, mod.Name); err == nil {
				for _, fld := range strct.Fields {
					val, typ, _, _ := resolveStructField(fld.Name, value, strct)
					fmt.Printf("\t%s.%s:\t\t%s\n", identName, fld.Name, PrintValue("", &val, typ, cxt))
				}
			}
		}


		return ""
	}

	return argName
}

// func PrintValue (value *[]byte, typName string) string {
// 	return ""
// }

func rep (str string, n int) string {
	return strings.Repeat(str, n)
}

// func (cxt *CXProgram) PrintProgram (isCompressed bool) {
// 	tab := "\t"
// 	nl := "\n"
// 	if isCompressed {
// 		tab = ""
// 		nl = ""
// 	}
	
// 	fmt.Println()
// 	fmt.Printf("(Modules %s", nl)
// 	for _, mod := range cxt.Modules {
// 		if mod.Name == CORE_MODULE {
// 			continue
// 		}
// 		fmt.Printf("%s(Module %s %s", rep(tab, 1), mod.Name, nl)

// 		fmt.Printf("%s(Imports %s", rep(tab, 2), nl)
// 		fmt.Printf("%s)%s", rep(tab, 2), nl) // imports
		
// 		fmt.Printf("%s(Definitions %s", rep(tab, 2), nl)

// 		for _, def := range mod.Definitions {
// 			fmt.Printf("%s(Definition %s %s %s)%s",
// 				rep(tab, 3),
// 				def.Name,
// 				def.Typ,
// 				PrintValue(def.Value, def.Typ),
// 				nl)
// 		}
		
// 		fmt.Printf("%s)%s", rep(tab, 2), nl) // definitions

// 		fmt.Printf("%s(Structs %s", rep(tab, 2), nl)

// 		for _, strct := range mod.Structs {
// 			fmt.Printf("%s(Struct %s", rep(tab, 3), nl)

// 			for _, fld := range strct.Fields {
// 				fmt.Printf("%s%s %s%s", rep(tab, 4), fld.Name, fld.Typ, nl)
// 			}
			
// 			fmt.Printf("%s)%s", rep(tab, 3), nl) // structs
// 		}
		
// 		fmt.Printf("%s)%s", rep(tab, 2), nl) // structs

// 		fmt.Printf("%s(Functions %s", rep(tab, 2), nl)

// 		for _, fn := range mod.Functions {
// 			fmt.Printf("%s(Function %s%s", rep(tab, 3), fn.Name, nl)

// 			fmt.Printf("%s(Inputs %s", rep(tab, 4), nl)
// 			for _, inp := range fn.Inputs {
// 				fmt.Printf("%s(Input %s %s)%s", rep(tab, 5), inp.Name, inp.Typ, nl)
// 			}
// 			fmt.Printf("%s)%s", rep(tab, 4), nl) // inputs

// 			fmt.Printf("%s(Outputs %s", rep(tab, 4), nl)
// 			for _, out := range fn.Outputs {
// 				fmt.Printf("%s(Output %s %s)%s", rep(tab, 5), out.Name, out.Typ, nl)
// 			}
// 			fmt.Printf("%s)%s", rep(tab, 4), nl) // outputs

// 			fmt.Printf("%s(Expressions %s", rep(tab, 4), nl)
// 			for _, expr := range fn.Expressions {
// 				_ = expr
// 				fmt.Printf("%s(Expression %s", rep(tab, 5), nl)

// 				fmt.Printf("%s(Operator %s)%s", rep(tab, 6), expr.Operator.Name, nl)
				
// 				fmt.Printf("%s(OutputNames %s", rep(tab, 6), nl)
// 				for _, outName := range expr.OutputNames {
// 					fmt.Printf("%s(OutputName %s)%s", rep(tab, 7), outName.Name, nl)
// 				}
// 				fmt.Printf("%s)%s", rep(tab, 6), nl)
				
// 				fmt.Printf("%s(Arguments %s", rep(tab, 6), nl)
// 				for _, arg := range expr.Arguments {
// 					fmt.Printf("%s(Argument %s %s)%s", rep(tab, 7), PrintValue(arg.Value, arg.Typ), arg.Typ, nl)
// 				}
// 				fmt.Printf("%s)%s", rep(tab, 6), nl)
				
// 				fmt.Printf("%s)%s", rep(tab, 5), nl)
// 			}
// 			fmt.Printf("%s)%s", rep(tab, 4), nl) // expressions
			
// 			fmt.Printf("%s)%s", rep(tab, 3), nl) // function
// 		}
		
// 		fmt.Printf("%s)%s", rep(tab, 2), nl) // functions
		
// 		fmt.Printf("%s)%s", rep(tab, 1), nl) // modules
// 	}
// 	fmt.Printf(")")
// 	fmt.Println()
// }

func CastArgumentForArray (typ string, arg *CXArgument) *CXArgument {
	switch typ {
	case "[]bool":
		return MakeArgument(arg.Value, "bool")
	case "[]byte":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		sVal := encoder.Serialize(byte(val))
		return MakeArgument(&sVal, "byte")
	case "[]str":
		return arg
	case "[]i32":
		return arg
	case "[]i64":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		sVal := encoder.Serialize(int64(val))
		return MakeArgument(&sVal, "i64")
	case "[]f32":
		return arg
	case "[]f64":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		sVal := encoder.Serialize(float64(val))
		return MakeArgument(&sVal, "f64")
	default:
		return arg
	}
	//panic(fmt.Sprintf("CastArgumentForArray: unrecognized type %s", typ))
}

func isBasicType (typ string) bool {
	for _, basic := range BASIC_TYPES {
		if basic == typ {
			return true
		}
	}
	return false
}

func IsNative (fnName string) bool {
	if _, ok := NATIVE_FUNCTIONS[fnName]; ok {
		return true
	}
	if _, ok := NATIVE_FUNCTIONS[strings.Split(fnName, ".")[1]]; ok {
		return true
	}
	return false
}
func IsArray (typ string) bool {
	if len(typ) > 2 && typ[:2] == "[]" {
		return true
	}
	return false
}
func IsStructInstance (typ string, mod *CXModule) bool {
	if _, err := mod.Context.GetStruct(typ, mod.Name); err == nil {
		return true
	} else {
		return false
	}
}
func IsLocal (identName string, call *CXCall) bool {
	for _, def := range call.State {
		if def.Name == identName {
			return true
		}
	}
	return false
}
func IsGlobal (identName string, mod *CXModule) bool {
	for _, def := range mod.Definitions {
		if def.Name == identName {
			return true
		}
	}
	for _, imp := range mod.Imports {
		for _, def := range imp.Definitions {
			if def.Name == identName {
				return true
			}
		}
	}
	return false
}

func (cxt *CXProgram) PrintProgram(withAffs bool) {

	fmt.Println("Program")
	if withAffs {
		for i, aff := range cxt.GetAffordances() {
			fmt.Printf(" * %d.- %s\n", i, aff.Description)
		}
	}

	i := 0
	for _, mod := range cxt.Modules {
		if mod.Name == CORE_MODULE || mod.Name == "glfw" || mod.Name == "gl" {
			continue
		}

		fmt.Printf("%d.- Module: %s\n", i, mod.Name)

		if withAffs {
			for i, aff := range mod.GetAffordances() {
				fmt.Printf("\t * %d.- %s\n", i, aff.Description)
			}
		}

		if len(mod.Imports) > 0 {
			fmt.Println("\tImports")
		}

		j := 0
		for _, imp := range mod.Imports {
			fmt.Printf("\t\t%d.- Import: %s\n", j, imp.Name)
			j++
		}

		if len(mod.Definitions) > 0 {
			fmt.Println("\tDefinitions")
		}

		j = 0
		for _, v := range mod.Definitions {
			fmt.Printf("\t\t%d.- Definition: %s %s\n", j, v.Name, v.Typ)
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
					k, fld.Name, fld.Typ)
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
				inOuts[in.Name] = in.Typ
			}
			
			
			var inps bytes.Buffer
			for i, inp := range fn.Inputs {
				if i == len(fn.Inputs) - 1 {
					inps.WriteString(concat(inp.Name, " ", inp.Typ))
				} else {
					inps.WriteString(concat(inp.Name, " ", inp.Typ, ", "))
				}
			}

			var outs bytes.Buffer
			for i, out := range fn.Outputs {
				if i == len(fn.Outputs) - 1 {
					outs.WriteString(concat(out.Name, " ", out.Typ))
				} else {
					outs.WriteString(concat(out.Name, " ", out.Typ, ", "))
				}
			}

			fmt.Printf("\t\t%d.- Function: %s (%s) (%s)\n",
				j, fn.Name, inps.String(), outs.String())

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
					if arg.Typ == "ident" {
						var id string
						encoder.DeserializeRaw(*arg.Value, &id)
						if arg.Typ != "" &&
							inOuts[id] != "" {
							typ = inOuts[id]
						} else if arg.Value != nil {
							var found *CXDefinition
							for _, def := range mod.Definitions {
								if def.Name == id {
									found = def
									break
								}
							}
							if found != nil && found.Typ != "" {
								typ = found.Typ
							}
						} else {
							typ = arg.Typ
						}
					} else {
						typ = arg.Typ
					}

					var argName string
					encoder.DeserializeRaw(*arg.Value, &argName)

					if arg.Typ != "ident" {
						switch typ {
						case "str":
							argName = fmt.Sprintf("%#v", argName)
						case "bool":
							var val int32
							encoder.DeserializeRaw(*arg.Value, &val)
							if val == 0 {
								argName = "false"
							} else {
								argName = "true"
							}
						case "byte":
							argName = fmt.Sprintf("%#v", *arg.Value)
						case "i32":
							var val int32
							encoder.DeserializeRaw(*arg.Value, &val)
							argName = fmt.Sprintf("%#v", val)
						case "i64":
							var val int64
							encoder.DeserializeRaw(*arg.Value, &val)
							argName = fmt.Sprintf("%#v", val)
						case "f32":
							var val float32
							encoder.DeserializeRaw(*arg.Value, &val)
							argName = fmt.Sprintf("%#v", val)
						case "f64":
							var val float64
							encoder.DeserializeRaw(*arg.Value, &val)
							argName = fmt.Sprintf("%#v", val)
						case "[]bool":
							var val []bool
							encoder.DeserializeRaw(*arg.Value, &val)
							argName = fmt.Sprintf("%#v", val)
						case "[]byte":
							var val []byte
							encoder.DeserializeRaw(*arg.Value, &val)
							argName = fmt.Sprintf("%#v", val)
						case "[]str":
							var val []string
							encoder.DeserializeRaw(*arg.Value, &val)
							argName = fmt.Sprintf("%#v", val)
						case "[]i32":
							var val []int32
							encoder.DeserializeRaw(*arg.Value, &val)
							argName = fmt.Sprintf("%#v", val)
						case "[]i64":
							var val []int64
							encoder.DeserializeRaw(*arg.Value, &val)
							argName = fmt.Sprintf("%#v", val)
						case "[]f32":
							var val []float32
							encoder.DeserializeRaw(*arg.Value, &val)
							argName = fmt.Sprintf("%#v", val)
						case "[]f64":
							var val []float64
							encoder.DeserializeRaw(*arg.Value, &val)
							argName = fmt.Sprintf("%#v", val)
						default:
							argName = string(*arg.Value)
						}
					}

					// if arg.Offset > -1 {
					// 	offset := arg.Offset
					// 	size := arg.Size
					// 	var val []byte
					// 	encoder.DeserializeRaw((*cxt.Heap)[offset:offset+size], &val)
					// 	arg.Value = &val
					// }

					if i == len(expr.Arguments) - 1 {
						args.WriteString(concat(argName, " ", typ))
					} else {
						args.WriteString(concat(argName, " ", typ, ", "))
					}
				}

				if len(expr.OutputNames) > 0 {
					var outNames bytes.Buffer
					for i, outName := range expr.OutputNames {
						if i == len(expr.OutputNames) - 1 {
							outNames.WriteString(outName.Name)
						} else {
							outNames.WriteString(concat(outName.Name, ", "))
						}
					}

					var exprTag string
					if expr.Tag != "" {
						exprTag = fmt.Sprintf(" :tag %s", expr.Tag)
					}

					fmt.Printf("\t\t\t%d.- Expression: %s = %s(%s)%s\n",
						k,
						outNames.String(),
						expr.Operator.Name,
						args.String(),
						exprTag)
				} else {
					var exprTag string
					if expr.Tag != "" {
						exprTag = fmt.Sprintf(" :tag %s", expr.Tag)
					}
					
					fmt.Printf("\t\t\t%d.- Expression: %s(%s)%s\n",
						k,
						expr.Operator.Name,
						args.String(),
						exprTag)
				}

				

				if withAffs {
					for i, aff := range expr.GetAffordances(nil) {
						fmt.Printf("\t\t\t * %d.- %s\n", i, aff.Description)
					}
				}
				
				k++
			}
			j++
		}
		i++
	}
}
