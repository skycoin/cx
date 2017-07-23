package base

import (
	"fmt"
	"math/rand"
	"time"
	"bytes"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

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

// Just a function to debug for myself. Not meant to be used in production nor meant to be efficient/maintanable
func (cxt *cxContext) PrintProgram(withAffs bool) {

	fmt.Println("Context")
	if withAffs {
		for i, aff := range cxt.GetAffordances() {
			fmt.Printf(" * %d.- %s\n", i, aff.Description)
		}
	}

	i := 0
	for _, mod := range cxt.Modules {
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

					if arg.Offset > -1 {
						offset := arg.Offset
						size := arg.Size
						var val []byte
						encoder.DeserializeRaw((*cxt.Heap)[offset:offset+size], &val)
						arg.Value = &val
					}

					if i == len(expr.Arguments) - 1 {
						args.WriteString(concat(string(*arg.Value), " ", typ))
					} else {
						args.WriteString(concat(string(*arg.Value), " ", typ, ", "))
					}
				}

				fmt.Printf("\t\t\t%d.- Expression: %s = %s(%s)\n",
					k,
					expr.OutputName,
					expr.Operator.Name,
					args.String())

				if withAffs {
					for i, aff := range expr.GetAffordances() {
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
